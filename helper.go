package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/Suy56/ProofChain/internal/crypto/keyUtils"
	"github.com/Suy56/ProofChain/internal/crypto/zkp"
	pb "github.com/Suy56/ProofChain/internal/crypto/zkp/rpc/proto"
	"github.com/Suy56/ProofChain/internal/ingest"
	mo "github.com/Suy56/ProofChain/internal/models"
	"github.com/Suy56/ProofChain/internal/users"
	"github.com/Suy56/ProofChain/internal/utils"
	"github.com/Suy56/ProofChain/storage/models"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)
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
	log.Println("public key of ins: ",pub)
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
        log.Println("Error decrypting:", err)
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
    for k, v := range utils.Walk(certificate) {
        field := vCert.FieldByName(k)
        if !field.IsValid() || !field.CanSet() {
            continue
        }
        if val, err := utils.CoerceToInt(fmt.Sprint(v)); err == nil {
            field.Set(reflect.ValueOf(val))
        } else {
            field.Set(reflect.ValueOf(v))
        }
    }
    return typedCert
}

type Proof interface {
    BuildProofRequest(
		constraints []string, 
		actualValue any, 
		salt string,
		siblings []string,
		
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
	siblings []string,
)*pb.ProofRequest {
	var ok bool
	r.actualValue,ok=(actualValue).(int); if !ok{
		log.Println("error asserting actual value to string for membership proof")
	}
	
	return &pb.ProofRequest{
		ProofData: &pb.ProofRequest_Range{
			Range: &pb.RangeRequest{
				ActualValue: uint32(r.actualValue),
				ActualSalt:  salt,
				Siblings:    siblings, 
				LowerBound:  uint32(r.lower),
				UpperBound:  uint32(r.upper),
				PublicRoot:  "", 
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
	siblings []string,
)*pb.ProofRequest {
	var ok bool
	m.actualValue,ok=(actualValue).(string); if !ok{
		log.Println("error asserting actual value to string for membership proof")
	}
	return &pb.ProofRequest{
		ProofData: &pb.ProofRequest_Membership{
			Membership: &pb.MembershipRequest{
				ActualValue: m.actualValue,
				ActualSalt:  salt,
				Siblings:    siblings, 
				PublicList:  m.members,
				PublicRoot:  "",
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

func readProofFile(attribute string, constraints []string, finder ingest.Finder)(Proof, []byte,error){
	attribute,val:=extractProofValues(constraints)
	switch val := val.(type){
	case Range:
		rangeVal:=val
		bytes,err:=finder.Discover(context.Background(),"",ingest.RangeBuilder(attribute,rangeVal.lower,rangeVal.upper))
		return val,bytes,err 
	case Membership:
		membershipVal:=val
		targetMembers := make([]any, len(membershipVal.members))
		for i, member := range membershipVal.members {
			targetMembers[i] = member
		}
		bytes,err:=finder.Discover(context.Background(),"",ingest.Membershipbuilder(attribute,targetMembers))
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