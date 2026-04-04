    // SPDX-License-Identifier: MPL 2.0
pragma solidity ^0.8.0;
contract Verification{
    address owner;
    constructor(){
        owner=msg.sender;
    }
    enum DocStatus{
        accepted,
        rejected,
        pending
    }
    struct User{
        string publicKey;   
    }
    struct Institution{
        address publicAddr;
        string  publicKey;
        string  name;
        bool    approved;
    }
    struct Document{
        address     requester;
        address     verifiedBy;
        string      institution; 
        DocStatus   status;
        uint        index;
    }
    //each institution can only have one verifier, at least for now. 
    mapping(string=>Institution) public institutions;
    mapping(address=>User) private users;
    //sha3 to map documents and verify the document's integrity 
    mapping(string=>uint) private documentList;
    mapping(address=>bool) userList;
    
    address[]   requesters;
    address[]   verifiedBy;
    //documentOrCertificateHash is used to track documents AND certificate issued
    //by the institute uniquely. 
    //Initially, when user uploads a document, the document is identified by it's shahash
    //When institute ISSUES a digital copy of document uploaded by the user, the shahash of user's
    //document gets replaced by the Merkle Root HASH.  
    string[]    documentOrCertificateHash;
    string[]    institution;
    string[]    descriptions;
    DocStatus[] status;
    //all the above arrays depends on  'docIndexCounter' variable
    //potential improvement: use mapping(user=>[]docs) and mapping(institute=>[]docs)
    //that way there won't be a need to rely and keep track of all these arrays separately
    
    //Though I am not sure if this is cost efficient, since there will be a need to iterate 
    //over  both, user and institute's arrays and update them separately
    uint docIndexCounter=0;
    function registerAsUser(string calldata _publicKey) public{
        users[msg.sender]=User({
            publicKey:  _publicKey
        });
        userList[msg.sender]=true;
    }
    
    function registerInstitution(string memory _publicKey, string memory _name) public{
        require(institutions[_name].publicAddr == address(0), "Institution already registered");
        institutions[_name]=Institution({
            publicAddr: msg.sender,
            publicKey:  _publicKey,
            name:       _name,
            approved:   false
        });
    }

    function getInstituePublicKey(string memory _name) public view returns(string memory pubKey){
        return institutions[_name].publicKey;
    }
    function getUserPublicKey(address userAddr)public view returns (string memory){
        return users[userAddr].publicKey;        
    }
    function approveVerifier(string memory _name)public{
        require(msg.sender==owner,"Only admin can perfom this action");
        institutions[_name].approved=true;
    }
    function addDocument(string memory shaHash, string memory _institute) public{
        require( 
            userList[msg.sender] ||
            msg.sender==institutions[_institute].publicAddr,"register first to upload"
        );
        documentList[shaHash]=docIndexCounter;
        requesters.push(msg.sender);
        verifiedBy.push(address(0));
        institution.push(_institute);
        status.push(DocStatus.pending);
        //the proof hash is initially empty as it's not approved
        //by the institution yet
        documentOrCertificateHash.push(shaHash);
        docIndexCounter++;
    }

    function addCertificate(string memory _hash, string memory _institute, address _requestor)public{
        require( 
            msg.sender==institutions[_institute].publicAddr && institutions[_institute].approved==true,
            "You're not approved to issue certificate for this institution"
        );

        documentList[_hash]=docIndexCounter;
        requesters.push(_requestor);
        verifiedBy.push(msg.sender);
        institution.push(_institute);
        status.push(DocStatus.accepted);
        //the proof hash is initially empty as it's not approved
        //by the institution yet
        documentOrCertificateHash.push(_hash);
        docIndexCounter++;

    }
    //returns all the documents
    function getDocuments()public view returns(
        address[] memory requester ,
        address[] memory verifer ,
        string[] memory institute,
        string[] memory hash,
        DocStatus[] memory stats
    ){
        return (requesters,verifiedBy,institution,documentOrCertificateHash,status);
    }
    
    function verifyDocument(
        string memory shaHash,
        string memory _institute,  
        DocStatus _status, 
        string memory _proofHash
    ) public payable {
        require(institutions[_institute].approved==true && institutions[_institute].publicAddr==msg.sender);
        uint index=documentList[shaHash];
        status[index]=_status;
        documentList[_proofHash]=documentList[shaHash];
        documentOrCertificateHash[index] = _proofHash;
        delete documentList[shaHash];
        verifiedBy[index]=msg.sender;
    }

    function getDocIndexCounter()public view returns(uint) {
        return docIndexCounter;
    }   
    function isApprovedInstitute(string memory name)public view returns (bool){
        return institutions[name].approved==true;
    }
}
