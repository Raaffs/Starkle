package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"

	blockchain "github.com/Suy56/ProofChain/chaincore/core"
	"github.com/Suy56/ProofChain/internal/crypto/keyUtils"
	"github.com/Suy56/ProofChain/internal/crypto/zkp"
	"github.com/Suy56/ProofChain/internal/download"
	"github.com/Suy56/ProofChain/internal/users"
	"github.com/Suy56/ProofChain/storage/models"
	storageclient "github.com/Suy56/ProofChain/storage/storage_client"
	"github.com/Suy56/ProofChain/internal/wallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"golang.org/x/sync/errgroup"
)

const (
	GANACHE    = "Ganache"
	INFURA     = "Infura"
	CLOUDFLARE = "CloudFlare"
	DRPC       = "dRPC"
)

// App struct
type App struct {
	ctx     context.Context
	account users.User
	keys    *keyUtils.ECKeys
	envMap  map[string]any
	storage *storageclient.Client
	proof   zkp.ZKProof
	config  Config
	logger  *slog.Logger
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (app *App) startup(ctx context.Context) {
	app.ctx = ctx
	keyMap, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("Error loading .env : ", err)
	}
	for key, val := range keyMap {
		app.envMap[key] = val
	}
	if err := app.config.Load(); err != nil {
		log.Fatalf("Fatal error: loading config failed\n%v", err)
	}
	app.proof = zkp.NewMerkleProof()
	app.storage = storageclient.New(app.config.Services.STORAGE)
	app.logger = NewLogger(os.Stdout)
}

func (app *App) Login(username string, password string) error {
	c := &blockchain.ClientConnection{}
	i := &blockchain.ContractVerifyOperations{}
	g, _ := errgroup.WithContext(context.Background())
	profile, ok := app.config.Profiles[username]
	if !ok {
		return fmt.Errorf("Profile doesn't exist")
	}
	g.Go(func() error {
		if err := app.keys.OnLogin(username, password, profile.KeyPath); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		privateKey, err := wallet.RetriveAccount(
			username,
			password,
			profile.AccountPath,
		)
		if err != nil {
			app.logger.Error(
				"Error retrieving user's wallet",
				"user", username,
				"path", profile.AccountPath,
				"err", err,
			)
			return fmt.Errorf("error retrieving account. Make sure the credentials are correct")
		}
		log.Println(privateKey)
		if err := blockchain.Init(
			c,
			i,
			privateKey,
			app.config.Services.CONTRACT_ADDR,
			app.config.Services.RPC_PROVIDERS_URLS.Local[GANACHE],
		); err != nil {
			log.Println(err,privateKey,
			)
			app.logger.Error(
				"Error connecting to the blockchain",
				"endpoint", app.config.Services.RPC_PROVIDERS_URLS.Local,
				"contract address", app.config.Services.CONTRACT_ADDR,
				"err", err,
			)
			return fmt.Errorf("error connecting to the blockchain")
		}
		approved, err := i.Instance.IsApprovedInstitute(c.CallOpts, username)
		if err != nil {
			app.logger.Error(
				"Error getting the account verification status",
				"username", username,
				"is approved", approved,
				"err", err,
			)
			return fmt.Errorf("error getting the account verification status")
		}
		app.logger.Info("Is authorized institution? ", "status", approved)
		if approved {
			app.account = &users.Verifier{Conn: c, Instance: i, Name: username}
		} else {
			app.account = &users.Requester{Conn: c, Instance: i}
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		return err
	}
	app.account.SetName(username)
	return nil
}

func (app *App) Logout() {
	app = &App{}
}

func (app *App) IsLoggedIn() bool {
	return app.account.GetTxOpts() != nil
}

func (app *App) Register(privateKeyString, name, password string, isInstitute bool) error {
	if len(privateKeyString) < 64 {
		return fmt.Errorf("invalid private key")
	}
	if _, exist := app.config.Profiles[name]; exist {
		return fmt.Errorf("Profile already exist, please use different name or private key")
	}
	c := &blockchain.ClientConnection{}
	i := &blockchain.ContractVerifyOperations{}
	var (
		publicKey    string
		accountPath  string
		keyPath      string
		identityPath string
	)
	if err := blockchain.Init(
		c,
		i,
		privateKeyString[2:],
		app.config.Services.CONTRACT_ADDR,
		app.config.Services.RPC_PROVIDERS_URLS.Local[GANACHE],
	); err != nil {
		app.logger.Error(
			"Error connecting to the blockchain",
			"endpoint", app.config.Services.RPC_PROVIDERS_URLS.Local,
			"contract address", app.config.Services.CONTRACT_ADDR,
			"err", err,
		)
		return fmt.Errorf("error connecting to the blockchain")
	}

	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		pub, path, err := app.keys.OnRegister(password, app.config.Dirs.Key)
		if err != nil {
			return err
		}
		publicKey = pub
		keyPath = path
		return nil
	})

	g.Go(func() error {
		path, err := wallet.NewWallet(
			privateKeyString[2:],
			name, password,
			app.config.Dirs.Account,
		)
		if err != nil {
			
			return err
		}
		accountPath = path
		return nil
	})
	if err := g.Wait(); err != nil {
		app.logger.Error("Error registering user", "err", err)

		return fmt.Errorf("error connecting to blockchain")
	}

	if err := app.config.AddProfile(name, accountPath, keyPath, identityPath); err != nil {
		app.logger.Error("Error adding profile","err", err)
		return fmt.Errorf("Failed to create profile")
	}

	if isInstitute {
		verifier := &users.Verifier{Conn: c, Instance: i, Name: name}
		app.account = verifier
		if err := app.account.Register(publicKey, name); err != nil {
			app.logger.Error(
				"Error registering the institute",
				"name", name,
				"err", err,
			)
			return fmt.Errorf("error registering institution")
		}
		app.account.SetName(name)
	} else {
		requester := &users.Requester{Conn: c, Instance: i}
		app.account = requester
		if err := app.account.Register(publicKey, name); err != nil {
			app.logger.Error("error registering requester","err", err)
			return fmt.Errorf("error registering institution")
		}
		app.account.SetName(name)
	}
	return nil
}

func (app *App) UploadDocument(institute, name, description string) error {
	var document models.Document
	if err := users.UpdateNonce(app.account); err != nil {
		app.logger.Error(
			"Invalid transaction nonce",
			"nonce",app.account.GetClient().TxOpts.Nonce,
			"err", err,
		)
		return fmt.Errorf("invalid transaction nonce")
	}
	file, path, err := app.GetFileAndPath()
	if err != nil {
		app.logger.Error("Error uploading File","err", err)
		return fmt.Errorf("Error uploading file")
	}

	encryptedDocument, err := app.Encrypt(file, institute)
	if err != nil {
		app.logger.Error("An error occurred while encrypting document","err", err)
		return fmt.Errorf("An error occurred while encrypting document")
	}

	shaHash, err := Keccak256File(path)
	if err != nil {
		app.logger.Error("Error hashing file","err", err,)
		return fmt.Errorf("Error uploading file")
	}
	document.EncryptedDocument = encryptedDocument
	document.Shahash = shaHash
	document.PublicAddress = app.account.GetPublicAddress()
	if err := app.storage.UploadDocument(document); err != nil {
		app.logger.Error(
			"Error uploading file to mongodb",
			"storage endpoint", app.config.Services.STORAGE,
			"err", err,
		)
		return fmt.Errorf("Error uploading file")
	}
	if account, ok := app.account.(*users.Requester); ok {
		if err := account.Instance.AddDocument(app.account.GetTxOpts(), shaHash, institute); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) GetAcceptedDocs() ([]blockchain.VerificationDocument, error) {
	docs, err := app.account.GetDocuments()
	if err != nil {
		app.logger.Error("Error retrieving accepted documents","err", err)
		return nil, fmt.Errorf("Error retrieving accepted documents")
	}
	verifiedDocs := app.account.GetAcceptedDocuments(docs)
	return verifiedDocs, nil
}

func (app *App) GetRejectedDocuments() ([]blockchain.VerificationDocument, error) {
	docs, err := app.account.GetDocuments()
	if err != nil {
		app.logger.Error("Error retrieving documents for rejection check","err", err)
		return nil, err
	}
	rejectedDocs := app.account.GetRejectedDocuments(docs)
	return rejectedDocs, nil
}

func (app *App) GetPendingDocuments() ([]blockchain.VerificationDocument, error) {
	docs, err := app.account.GetDocuments()
	if err != nil {
		app.logger.Error("Error retrieving documents for pending check","err", err,)
		return nil, err
	}
	pendingDocs := app.account.GetPendingDocuments(docs)
	return pendingDocs, nil
}
func (app *App) CreateDigitalCopy(status int, hash string, certificate models.CertificateData) error {
	if err := users.UpdateNonce(app.account); err != nil {
		app.logger.Error(
			"Invalid transaction nonce",
			"nonce",app.account.GetClient().TxOpts.Nonce,
			"err", err,
		)
		return err
	}
	_, ok := app.account.(*users.Verifier)
	if !ok {
		return fmt.Errorf("You're not approved to perform this action")
	}

	switch status {
	case blockchain.Rejected:
		if _, err :=
			app.account.GetInstance().Instance.VerifyDocument(
				app.account.GetTxOpts(),
				hash,
				app.account.GetName(),
				uint8(status),
				hash,
			); err != nil {
			app.logger.Error("Error approving document (rejection path)","err", err)
			return fmt.Errorf("An error occurred ")
		}
		return nil

	case blockchain.Pending:
		return nil
	}
	doc, publicCommit, err := app.PrepareDigitalCopy(certificate)
	if err != nil {
		app.logger.Error("Error preparing digital copy","err", err)
		return fmt.Errorf("An error occurred while issuing document")
	}

	if err := app.storage.UploadDocument(doc); err != nil {
		app.logger.Error("Error uploading certificate to storage","err", err)
		return fmt.Errorf("Error creating certificate")
	}

	if _, err := app.account.GetInstance().Instance.VerifyDocument(
		app.account.GetTxOpts(),
		hash,
		app.account.GetName(),
		0,
		publicCommit,
	); err != nil {
		app.logger.Error("Error verifying document on blockchain","err", err)
		return nil
	}
	return nil
}

func (app *App) IssueCertificate(certificate models.CertificateData) error {
    if err := users.UpdateNonce(app.account); err != nil {
        app.logger.Error(
            "Invalid transaction nonce",
			"nonce",app.account.GetTxOpts().Nonce,
            "err", err,
        )
        return err
    }
    doc, publicCommit, err := app.PrepareDigitalCopy(certificate)
    if err != nil {
        app.logger.Error("Error preparing digital copy","err", err)
        return fmt.Errorf("An error occurred while issuing certificate")
    }
    if _, err := app.account.GetInstance().Instance.AddCertificate(
        app.account.GetTxOpts(),
        publicCommit,
        app.account.GetName(),
        common.HexToAddress(certificate.PublicAddress),
    ); err != nil {
        app.logger.Error("Error adding certificate to blockchain","err", err)
        return fmt.Errorf("an error occurred while issuing certificate")
    }
    if err := app.storage.UploadDocument(doc); err != nil {
        app.logger.Error("Error uploading document to storage","err", err)
        return fmt.Errorf("Error issuing certificate")
    }
    return nil
}

func (app *App) ViewDocument(shahash, instituteName, requesterAddress string) (string, error) {
    encryptedDocument, err := app.storage.RetrieveDocument(shahash)
    if err != nil {
        app.logger.Error(
            "Error retrieving document",
            "hash", shahash,
            "err", err,
        )
        return "", fmt.Errorf("Error retrieving document")
    }
    decryptedDoc, err := app.TryDecrypt(encryptedDocument.EncryptedDocument, instituteName, requesterAddress)
    if err != nil {
        app.logger.Error(
            "Error decrypting document",
            "hash", shahash,
            "requester", requesterAddress,
            "err", err,
        )
        return "", fmt.Errorf("Error decrypting document")
    }
    encodedDocument := base64.StdEncoding.EncodeToString(decryptedDoc)
    return encodedDocument, nil
}

func (app *App) ViewDigitalCertificate(hash, instituteName, requesterAddress string) (map[string]any, error) {
    var cert map[string]any
    decryptedCert, err := app.getDecryptedCertificate(hash, instituteName, requesterAddress)
    if err != nil {
        return nil, err
    }

    if err := json.Unmarshal(decryptedCert, &cert); err != nil {
        app.logger.Error(
            "Error unmarshaling decrypted certificate",
            "hash", hash,
			"institute",instituteName,
			"req addr",requesterAddress,
            "err", err,
        )
		return nil,fmt.Errorf("an error occurred")
    }
    return cert, nil
}
func (app *App) Download(hash, instituteName, requesterAddress string) (string, error) {
	decryptedCert, err := app.getDecryptedCertificate(hash, instituteName, requesterAddress)
    if err != nil {
        app.logger.Error(
            "An error occurred while downloading/decrypting certificate",
            "hash", hash,
            "err", err,
        )
        return "", fmt.Errorf("an error occurred while downloading")
    }

	downloader, err := download.NewDownloader(decryptedCert, NewLogger(os.Stdout))
	if err != nil {
		app.logger.Error("error creating new downloader","err", err)
		return "", fmt.Errorf("an error occurred while downloading")
	}

	if err := downloader.Exec(); err != nil {
		app.logger.Error("error downloading", "err", err)
	}
	return "Downloaded successfully", nil
}

func (app *App) GenerateZKP(hash, instituteName, requesterAddress string, publicConstraint []any )(string, error){
	decryptedCert, err := app.getDecryptedCertificate(hash, instituteName, requesterAddress)
    if err != nil {
        app.logger.Error(
            "An error occurred while downloading/decrypting certificate",
            "hash", hash,
            "err", err,
        )
        return "", fmt.Errorf("an error occurred while downloading")
    }
	log.Println(decryptedCert)
	return "proof generated successfully",nil
}