<a name="readme-top"></a>
<br />
<div align="center">

  <h1 align="center">ProofChain</h1>

  <p align="center">
</div>



<!-- TABLE OF CONTENTS -->
<!-- <details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#development">Development</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#Se">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#UML">UML</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details> -->



<!-- ABOUT THE PROJECT -->
## About The Project
Proofchain is a document verficationa and certficate issuance platform, allowing users to selectively disclose their identity to third party verifiers without exposing any extra Personally Identifiable Information
* Digital certificates and Digital copies of documents issued by authorized institutions are structured as Merkle Trees where only the root is stored on-chain. This allows users to provide   cryptographic proofs for individual fields that verifiers can validate against the Ethereum ledger.
* The public ECDH keys of institutions and requestors are stored on blockchain
* The digital certificates and documents are encrypted using ECDH for key exchange and AES-256 for encryption and stored off-chain on mongodb, ensuring only the requestor and issuing institution can view the document
* Third-party verifiers can recompute the Merkle tree from shared fields and confirm authenticity by comparing the result against the issuer’s on-chain root, verifying specific data points without accessing unrevealed fields.
<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

[![Go][Go]][Go-url]
[![React][React.js]][React-url]
[![rust][rust]][rust-url]
[![wails][wails]][wails-url]
[![risc0][risc0]][risc0-url]
[![Ethereum][Ethereum]][Ethereum-url]
[![mongodb][mongodb]][mongodb-url]


<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Development

### Prerequisites

1. Ganache
    ```sh
    npm install ganache --global
    ```
2. Wails
    ```sh
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    ```
3. ProofChain-store 
    ```sh
    git clone https://github.com/Raaffs/ProofChain-Store.git
    ```
4. Clone Repository
    ```sh
    git clone https://github.com/Raaffs/ProofChain.git
    ```
  ---
### Set Up

1. Set up Envihttps://img.shields.io/badge/RISC0-FFC700?style=for-the-badgeroment variables in `.env` file
    ```
    cd ProofChain
    PRIVATE_KEY=YOUR_PRIVATE_KEY
    ```
2. Set up config
   ```sh
   cp .config.example.json .config.json
   ```
3. Deploy Contract
   ```sh
   go test -v ./test -run TestDeploy
   ```
4. Add contract address to .config.json
   ```sh
   .config.json
   ...
       "services": {
        "STORAGE": "localhost:8754",      
        "CONTRACT_ADDR": "CONTRACT_ADDR" , #edit this
        "RPC_PROVIDERS_URLS": {
    ...
   ```
---
### Set Up storage service
1. Go to the directory where you install ProofChain-Store
2. Set up .env
   ```sh
   MONGO_URI=your_mongo_url
   MONGO_DB=ProofChain
   MONGO_COLLECTION_DOCUMENTS=Documents
   MONGO_COLLECTION_INSTITUTES=institute
   # Application Secret Key to access secure routes and perform sensitive operations
   SECRET_KEY=secret
   ```
4. Install dependencies
   ```sh
   go mod download
   ```
3. Run the storage service
   ```bash
   go run .
   ```
  Storage service should be up on port 8754
  
  ___Note__: If you are running storage service on some other port, make sure to edit .config.json in proofchain to that specific port_
  
  ---
### Build & Run the app
Make sure you've ganache & storage service up and running
```bash
wails build
```
```bash
./build/bin/ProofChain
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->



<!-- UML -->
## UML
 ### 1. Uploading Documents
<img width="1759" height="1335" alt="Sequence diagram(1)" src="https://github.com/user-attachments/assets/5cddf928-6f11-494f-bc87-1eecbf4f1ef9" />

---
 ### 2. Issue certificate or create digital copy
 <img width="1404" height="1310" alt="issue" src="https://github.com/user-attachments/assets/cc3359db-2b01-48b0-a3fb-94910cbb9ffc" />

---
 ### 3. Third Party Verification
 <img width="1474" height="1047" alt="verify" src="https://github.com/user-attachments/assets/d3a5d068-31ff-4a17-a50b-aa2cca800dfa" />

---
<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTRIBUTING -->
## Contributing


If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!
<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTACT -->
## Contact

Suyash - suyashsaraf5@gmail.com

---
Thank You!

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/
[React.js]: https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB
[React-url]: https://reactjs.org/
[Ethereum]: https://img.shields.io/badge/Ethereum-3C3C3D?style=for-the-badge&logo=Ethereum&logoColor=white
[Ethereum-url]: https://ethereum.org/
[mongodb]: https://img.shields.io/badge/-MongoDB-13aa52?style=for-the-badge&logo=mongodb&logoColor=white
[mongodb-url]: https://www.mongodb.com/
[wails]: https://img.shields.io/badge/wails-red?style=for-the-badge&logo=wails
[wails-url]: https://wails.io
[risc0]: https://img.shields.io/badge/RISC0-FFC700?style=for-the-badge
[risc0-url]: https://risczero.com/
[rust]: https://img.shields.io/badge/Rust-E57324?style=for-the-badge&logo=rust&logoColor=white
[rust-url]: https://rust-lang.org/
