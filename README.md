<a name="readme-top"></a>
<br />
<div align="center">

  <h1 align="center">Starkle</h1>

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
Starkle is a privacy-preserving identity platform that allows users to selectively disclose attributes or prove identity predicates (e.g., age range or set-membership) without exposing raw Personally Identifiable Information (PII). It utilizes Merkle Trees for point disclosure and RISC0 (STARKs) for verifiable computational logic.

* Digital certificates and Digital copies of documents issued by authorized institutions are structured as Merkle Trees where only the root is stored on-chain. This allows users to provide   cryptographic proofs for individual fields that verifiers can validate against the Ethereum ledger.
* Utilizes RISC0 zkVM to generate STARK receipts for identity predicates. This allows users to prove specific data conditions (like age ranges) without revealing the actual attribute values.
* A high-performance Go-Rust gRPC pipeline delivers public and private inputs alongside O(logn) sibling paths to the zkVM. The architecture is designed to minimize guest memory overhead during the proof generation process for documents containing large attribute sets.
* Documents are structured as Merkle Trees and encrypted off-chain in MongoDB using AES-256 and ECDH. Only the issuing institution and the owner via cryptographic signing can access or reveal the records.
* Verification is achieved either by validating STARK receipts where the guest code has committed the Merkle root or by providing a Merkle proof of disclosed values and sibling paths. In the latter case, the verifier recomputes the Merkle root and confirms authenticity by comparing the result against the issuer's on-chain root.
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
    git clone https://github.com/Raaffs/Starkle.git
    ```
  ---
### Set Up

1. Set up in `.env` file
    ```
    cd Starkle
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
5. Set up risc0
   ```sh
     cd Starkle/internal/crypto/zkp
   ```
   ```sh
   NVCC_APPEND_FLAGS="-arch=sm_86 -DNDEBUG" \
   RISC0_GPU_BACKEND=cuda \
   cargo build --release --bin host
   ```
6. Start prover
   ```sh
   ./target/release/host                
   ```
---
### Set Up storage service
1. Go to the directory where you install ProofChain-Store
2. Set up .env
   ```sh
   MONGO_URI=your_mongo_url
   MONGO_DB=Starkle
   MONGO_COLLECTION_DOCUMENTS=Documents
   MONGO_COLLECTION_INSTITUTES=institute
   # Application Secret Key to access secure routes and perform sensitive operations
   SECRET_KEY=secret
   ```
3 . Install dependencies
   ```sh
   go mod download
   ```
4. Run the storage service
```bash
  go run .
```
  Storage service should be up on port 8754
  
  ___Note__: If you are running storage service on some other port, make sure to edit .config.json in starkle to that specific port_
  
  ---
### Build & Run the app
Make sure you've ganache & storage service up and running
```bash
wails build
```
```bash
./build/bin/starkle
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->



<!-- UML -->
## UML
 ### 1. Systeml Workflow
<img width="1287" height="832" alt="system workflow" src="https://github.com/user-attachments/assets/b8966fab-716f-4829-bc83-e2778918af9c" />


 ### 2. Uploading Documents
<img width="1662" height="1459" alt="1" src="https://github.com/user-attachments/assets/2506121d-06b6-4b76-95db-a99c09eb29d2" />

---
 ### 3. Creating digital copy
<img width="1647" height="1312" alt="2" src="https://github.com/user-attachments/assets/44e06519-68e7-49c6-bf42-436fbf2e333d" />

---
 ### 4. Issuing digital certificate 
<img width="1373" height="1083" alt="4" src="https://github.com/user-attachments/assets/f9393a43-07fd-43cf-8712-c5df9643801a" />

---
 ### 5. Third Party Verification
<img width="1434" height="1074" alt="3" src="https://github.com/user-attachments/assets/98bde04d-f3d3-4a79-b9cd-aad0f8d5f797" />

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
