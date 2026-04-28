package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	blockchain "github.com/Suy56/ProofChain/chaincore/core"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func FilterDocument(docs []blockchain.VerificationDocument, condition func(blockchain.VerificationDocument)bool)[]blockchain.VerificationDocument{
	var userDocs []blockchain.VerificationDocument
	for _,doc :=range docs{
		if(condition(doc)){	
			userDocs=append(userDocs,doc)
		}
	}
	return userDocs
}



func Walk[S any](s S) func(yield func(string, any) bool) {
	v := reflect.ValueOf(s)

	// Dereference pointer if needed
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return func(yield func(string, any) bool) {}
	}

	return func(yield func(string, any) bool) {
		t := v.Type()
		numFields := v.NumField()

		for i := range numFields {
			field := t.Field(i)
			value := v.Field(i)

			if !field.IsExported() {
				continue
			}

			switch value.Kind() {
			case reflect.Map:
				// Iterate map keys
				for _, key := range value.MapKeys() {
					val := value.MapIndex(key)
					if !yield(fmt.Sprint(key.Interface()), val.Interface()) {
						return
					}
				}
			default:
				// Use the struct field name as the attribute key
				if !yield(field.Name, value.Interface()) {
					return
				}
			}
		}
	}
}



func Keccak256File(path string) (string, error) {
	file, err := os.Open(path);if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	hasher := sha3.New256()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %v", err)
	}

	hashBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString, nil
}

func FormatDateToInt(date string) (int, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, date)
	if err != nil {
		return 0, err
	}
	res,err:=strconv.Atoi(t.Format("20060102"))
	if err!=nil {
		return -1, fmt.Errorf("Error converting date to int : ",err)
	}
	return res, nil
}

func isCurrency(input string) bool {
	// Pattern: Starts with 1-3 digits, followed by zero or more groups of (comma + 3 digits)
	// and an optional decimal suffix (.xx)
	const pattern = `^\d{1,3}(,\d{3})*(\.\d{2})?$`
	
	// Compile the regex
	re := regexp.MustCompile(pattern)
	
	return re.MatchString(input)
}
// It supports direct numeric strings and YYYY-MM-DD date formats.
func CoerceToInt(s string) (int, error) {
	s = strings.TrimSpace(s) // Clean up any stray whitespace first

	// 1. Try direct conversion (e.g., "123" -> 123)
	if val, err := strconv.Atoi(s); err == nil {
		return val, nil
	}

	// 2. Try date conversion (e.g., "2026-03-26" -> 20260326)
	// Using the Go reference date: Jan 2, 2006
	if t, err := FormatDateToInt(s); err == nil {
		// Format to "20060102" then convert that string to int
		return t,nil
	}
	// 3. Check if it's a currency format (e.g., "1,234" -> 1234)
	if isCurrency(s) {
		cleaned := strings.ReplaceAll(s, ",", "")
		if val, err := strconv.Atoi(cleaned); err == nil {
			return val, nil
		}
	}

	return -1, fmt.Errorf("value %q cannot be coerced to int", s)
}

func GetAttributeValue(obj any, fields ...string) (any, error) {
	val := reflect.ValueOf(obj)

	for _, field := range fields {
		// 1. Handle Pointers: Dereference until we find the actual value
		for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
			if val.IsNil() {
				return nil, fmt.Errorf("encountered nil pointer while looking for field '%s'", field)
			}
			val = val.Elem()
		}

		// 2. Ensure we are looking at a Struct
		if val.Kind() != reflect.Struct {
			return nil, fmt.Errorf("cannot get field '%s' from non-struct type %s", field, val.Type())
		}

		// 3. Find the field
		val = val.FieldByName(field)

		// 4. Check if field exists
		if !val.IsValid() {
			return nil, fmt.Errorf("field '%s' does not exist in struct", field)
		}

		// 5. Check for unexported fields (private fields)
		// val.Interface() panics if the field is not exported.
		if !val.CanInterface() {
			return nil, fmt.Errorf("field '%s' is unexported and cannot be accessed", field)
		}
	}

	return val.Interface(), nil
}

func GetDirPath(name string)(string,error){
	var dir string
	cmd:=exec.Command("xdg-user-dir", name)
	var out bytes.Buffer

	if err:=cmd.Run();err==nil{
		dir = strings.TrimSpace(out.String())
	}

	if dir==""{
		home,err:=os.UserHomeDir()
		if err!=nil{
			return "",fmt.Errorf("failed to detect home directory: %w", err)
		}
		dir=filepath.Join(home,name)
	}

	finalPath:=filepath.Join(dir,"ProofChain")
	if err:=os.MkdirAll(finalPath,0755);err!=nil{
		return "",fmt.Errorf("failed to ensure ProofChain directory: %w", err)
	}
	return finalPath,nil
}
