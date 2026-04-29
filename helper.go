package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"

	"github.com/Suy56/ProofChain/internal/crypto/keyUtils"
	"github.com/Suy56/ProofChain/internal/crypto/zkp"
	pb "github.com/Suy56/ProofChain/internal/crypto/zkp/rpc/proto"
	"github.com/Suy56/ProofChain/internal/ingest"
	mo "github.com/Suy56/ProofChain/internal/models"
	"github.com/Suy56/ProofChain/internal/users"
	"github.com/Suy56/ProofChain/internal/utils"
	"github.com/Suy56/ProofChain/storage/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getClient() (pb.ProverServiceClient, func(),error) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil,nil,fmt.Errorf("Failed to connect: %v", err)
	}
	client := pb.NewProverServiceClient(conn)
	return client, func() { conn.Close() },nil
}

func (app *App)Encrypt(file []byte, entity string)([]byte,error){
	pubKey,err:=app.account.GetPublicKeys(entity);if err!=nil{
		return nil,err
	}
	if pubKey==""{
		log.Println("error retrieving public keys")
		return nil,fmt.Errorf("invalid institution")
	}
	if err:=app.keys.SetMultiSigKey(pubKey);err!=nil{
		return nil,err
	}
	secretKey,err:=app.keys.GenerateSecret();if err!=nil{
		return nil,err
	}
	encryptedDocument,err:=keyUtils.Encrypt(secretKey,file);if err!=nil{
		return nil,err
	}
	return encryptedDocument,nil
}	

func(app *App)TryDecrypt(encryptedDocument []byte,institute string,user string)([]byte,error){
	var targetEntity string
	if _,ok:=app.account.(*users.Requester); ok{
		targetEntity=institute
	}
	if _,ok:=app.account.(*users.Verifier); ok{
		targetEntity=user
	}
	pub,err:=app.account.GetPublicKeys(targetEntity); if err!=nil{
		return nil,fmt.Errorf("helper.go: error retrieving public keys %w",err)
	}
	if err:=app.keys.SetMultiSigKey(pub);err!=nil{
		log.Println("error setting multisigkey: ",err)
		return nil,fmt.Errorf("Error retrieving multi-sig keys")
	}
	sec,err:=app.keys.GenerateSecret();if err!=nil{
		log.Println("error generating secret: ",err)
		return nil,fmt.Errorf("Error generating secret key")
	}
	document,err:=keyUtils.Decrypt(sec,encryptedDocument);if err!=nil{
		log.Println("error decrypting ipfs hash here: ",err)
		return nil,fmt.Errorf("You're not authorized")
	}
	return document,nil
}

func (app *App)GetFileAndPath()([]byte, string, error){
	filePath,err:=runtime.OpenFileDialog(app.ctx,runtime.OpenDialogOptions{
		Title: "Select Document",
		Filters: []runtime.FileFilter{	
			{
				DisplayName: "Documents (*.pdf; *.jpg; *.png)",
				Pattern: "*.pdf;*.jpg;*.png",
			},
		},
	})
	if err!=nil{
		return nil,"",err
	}
	file,err:=os.ReadFile(filePath);if err!=nil{
		log.Println("Error reading file : ",err)
		return nil,"",err
	}
	return file,filePath,nil
}

func (app *App)IsApprovedInstitute()bool{
	approved,err:=app.account.GetApprovalStatus();if err!=nil{
		log.Println("Error getting approval status : ",err)
		return false
	}
	return approved
}

func (app *App)PrepareDigitalCopy(certificate models.CertificateData)(models.Document,string,error){
	proof:=zkp.NewMerkleProof()
	typedCert := mapToCertificateBase[string,any](certificate)
	publicCommit,saltedCertificate,err:=proof.GenerateRootProof(typedCert);if err!=nil{
		app.logger.Error(
			"error generating proof",
			"err",err.Error(),
		)
		return models.Document{},"",fmt.Errorf("an error occurred while issuing certificate")
	}
	json, err := json.Marshal(saltedCertificate);if err!=nil{
		app.logger.Error("Error Marshalling salted certificate", "err", err)
		 log.Println(err)

		return models.Document{},"",fmt.Errorf("invalid certificate format")
	}
	encryptedCertificate,err:=app.Encrypt(json,certificate.PublicAddress);if err!=nil{
		return models.Document{},"",fmt.Errorf("error encrypting document %w",err)
	}
	doc:=models.Document{
		Shahash: string(publicCommit),
		EncryptedDocument: encryptedCertificate,
		PublicAddress: certificate.PublicAddress,
	}
	return doc,string(publicCommit),nil
}

func (app *App) FetchAndParseCertificate(hash, instituteName, requesterAddress string) (mo.CertificateBase[mo.LeafFields], error) {
	var tester map[string]any
    decryptedCert, err := app.getDecryptedCertificate(hash, instituteName, requesterAddress)
    if err != nil {
        app.logger.Error(
            "An error occurred while downloading/decrypting certificate",
            "hash", hash,
            "err", err,
        )
        return mo.CertificateBase[mo.LeafFields]{}, fmt.Errorf("an error occurred while downloading")
    }

    var wrapper struct {
        SaltedFields mo.CertificateBase[mo.LeafFields] `json:"salted_fields"`
    }

    if err := json.Unmarshal(decryptedCert, &wrapper); err != nil {
        app.logger.Error(
            "Error unmarshaling decrypted certificate",
            "hash", hash,
            "institute", instituteName,
            "req addr", requesterAddress,
            "err", err,
        )
        return mo.CertificateBase[mo.LeafFields]{}, fmt.Errorf("an error occurred while downloading")
    }
	json.Unmarshal(decryptedCert, &tester)
	log.Println(" tester decrypted cert: ",tester)
    return wrapper.SaltedFields, nil
}

func (app *App) getDecryptedCertificate(hash, instituteName, requesterAddress string) ([]byte, error) {
    encryptedCert, err := app.storage.RetrieveDocument(hash)
    if err != nil {
        app.logger.Error(
			"Error retrieving document",
			"storage endpoint",app.config.Services.STORAGE,
			"hash",hash,
			"err",err,
		)
        return nil, fmt.Errorf("error retrieving document")
    }

    decryptedCert, err := app.TryDecrypt(encryptedCert.EncryptedDocument, instituteName, requesterAddress)
    if err != nil {
        return nil, fmt.Errorf("error decrypting document")
    }

    return decryptedCert, nil
}


//this is insane, and should never be done again. NOT DRYing is fine if it avoids this monstrosity
//I just wanted to see how'd it work, so I did it but I'd never use this in actual prod

//We're initializing input and output of mo.CertificateBase[] differently because when the issuer
//inputs the fields it's of type mo.CertificateBase[string] but to generate merkle proof we need to 
//store the associate salt and hash in type of struct mo.LeafField, and all methods associated with proof
//generation requires mo.CertificateBase[mo.LeafField] or mo.CertificateBase[any]


//The attempted coercion of string to int is done so that any attributes 
//that needs zk-proof of range (such as date or salary) needs to be in format of LittleEndian UInt32
//but to keep schema flexible, issuer basically uses a map[string]string, so we need to dynamically figure out
//which attributes are compatible with zk-proof of range instead of needing the user to 
// specifiy which proof the want 
//to generate 
func mapToCertificateBase[T any, U any](certificate mo.CertificateBase[T]) mo.CertificateBase[U] {
    var typedCert mo.CertificateBase[U]
    vCert := reflect.ValueOf(&typedCert).Elem()
    
    // Identify if the struct has a Map field to collect 'extra' values
    var mapField reflect.Value
    for i := 0; i < vCert.NumField(); i++ {
        if vCert.Field(i).Kind() == reflect.Map {
            mapField = vCert.Field(i)
            // Initialize the map if it's nil
            if mapField.IsNil() {
                mapField.Set(reflect.MakeMap(mapField.Type()))
            }
            break
        }
    }

    for k, v := range utils.Walk(certificate) {
        field := vCert.FieldByName(k)
        
        // Coerce value before assignment
        var finalValue reflect.Value
        if val, err := utils.CoerceToInt(fmt.Sprint(v)); err == nil {
            finalValue = reflect.ValueOf(val)
        } else {
            finalValue = reflect.ValueOf(v)
        }
        
        if field.IsValid() && field.CanSet() {
            // Standard field assignment
            field.Set(finalValue)
        } else if mapField.IsValid() {
            // Insert coerced value into the struct's map field
            mapField.SetMapIndex(reflect.ValueOf(k), finalValue)
        }
    }
    return typedCert
}

type Proof interface {
    BuildProofRequest(
		constraints []string, 
		actualValue any, 
		salt string,
		expectedRoot string,
		siblings []string,
		basePath string,
	)*pb.ProofRequest
}

type Range struct{
	upper int
	lower int
	actualValue int
}

func (r Range) BuildProofRequest(
	constraints []string, 
	actualValue any, 
	salt string,
	expectedRoot string,
	siblings []string,
	basePath string,
)*pb.ProofRequest {
	switch v := actualValue.(type) {
		case float64:
    		r.actualValue = int(v)
		case int:
   			 r.actualValue = v
		case int64:
    		r.actualValue = int(v)
	default:
	    log.Printf("Unsupported type %T for value %v", v, v)
}
	
	return &pb.ProofRequest{
		ProofData: &pb.ProofRequest_Range{
			Range: &pb.RangeRequest{
				ActualValue: uint32(r.actualValue),
				ActualSalt:  salt,
				Siblings:    siblings, 
				LowerBound:  uint32(r.lower),
				UpperBound:  uint32(r.upper),
				PublicRoot:  expectedRoot, 
				Path: basePath,
			},
		},
	}
}

type Membership struct{
	members []string
	actualValue string
}

func (m Membership) BuildProofRequest(
	constraints []string, 
	actualValue any, 
	salt string,
	expectedRoot string,
	siblings []string,
	basePath string,
)*pb.ProofRequest {
	var ok bool
	m.actualValue,ok=(actualValue).(string); if !ok{
		log.Println("error asserting actual value to string for membership proof", "actual value: ",actualValue,"constraints: ",constraints)
		log.Printf("Expected int, but got type %T with value %v", actualValue, actualValue)
	}
	return &pb.ProofRequest{
		ProofData: &pb.ProofRequest_Membership{
			Membership: &pb.MembershipRequest{
				ActualValue: m.actualValue,
				ActualSalt:  salt,
				Siblings:    siblings, 
				PublicList:  m.members,
				PublicRoot:  expectedRoot,
				Path: basePath,
			},
		},
	}
}


func resolveIngestionMode(mode string, maxWorkers int)ingest.Finder{
	if mode=="manual"{
		return ingest.NewSelection()
	}
	return ingest.NewParallel(maxWorkers)
}

func readProofFile(attribute string, constraints []string, path string, finder ingest.Finder)(Proof, []byte,error){
	attribute,val:=extractProofValues(constraints)
	switch val := val.(type){
	case Range:
		rangeVal:=val
		bytes,err:=finder.Discover(context.Background(),path, func(b []byte) bool {
			return rangeComparator(b, attribute, rangeVal.lower, rangeVal.upper)
		})
		return val,bytes,err 
	case Membership:
		membershipVal:=val
		targetMembers := make([]any, len(membershipVal.members))
		for i, member := range membershipVal.members {
			targetMembers[i] = member
		}
		bytes,err:=finder.Discover(context.Background(),path,func(b []byte) bool {
			return memerbshipComparator(b, attribute, targetMembers)
		})
		return val,bytes,err
	default:
		return nil,nil,fmt.Errorf("invalid proof type")
	}
}

func extractProofValues(constraints []string)(string,Proof){
	if len(constraints)==3{
		lower,errl:= utils.CoerceToInt(constraints[1]); 
		upper,erru:=utils.CoerceToInt(constraints[2]); 
		if errl==nil && erru==nil{
			return constraints[0],Range{lower: lower,upper: upper}
		}
	}
	return constraints[0],Membership{members: constraints[1:]}
}

var memerbshipComparator = func (rawJson []byte, field string, targets []any) bool {
	var data map[string]mo.LeafFields
	if err := json.Unmarshal(rawJson, &data); err != nil {
		return false
	}
	val, exists := data[field]
	if !exists {
		return false
	}
	return slices.Contains(targets, val.Value)
}

var rangeComparator = func (rawJson []byte, field string, lower int, upper int) bool {
	var data map[string]mo.LeafFields
	if err := json.Unmarshal(rawJson, &data); err != nil {
		return false
	}
	val, exists := data[field]
	if !exists {
		return false
	}
	// JSON numbers unmarshal as float64. We convert to float for the comparison
	// or cast to int if we're sure it's a whole number.
	num, ok := val.Value.(float64)
	if !ok {
		return false
	}

	intVal := int(num)
	return intVal >= lower && intVal <= upper
}
