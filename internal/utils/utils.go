package utils

import (
	"crypto/rand"
	"crypto/sha3"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"reflect"
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

	return -1, fmt.Errorf("value %q cannot be coerced to int", s)
}

func GetAttributeValue(obj any, fields ...string) any {
	val := reflect.ValueOf(obj)

	for _, field := range fields {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		val = val.FieldByName(field)
		
		if !val.IsValid() {
			return nil
		}
	}
	return val.Interface()
}

