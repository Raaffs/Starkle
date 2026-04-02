export namespace blockchain {
	
	export class VerificationDocument {
	    ID: number;
	    Requester: string;
	    Verifier: string;
	    Institute: string;
	    ShaHash: string;
	    Stats: number;
	
	    static createFrom(source: any = {}) {
	        return new VerificationDocument(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Requester = source["Requester"];
	        this.Verifier = source["Verifier"];
	        this.Institute = source["Institute"];
	        this.ShaHash = source["ShaHash"];
	        this.Stats = source["Stats"];
	    }
	}

}

export namespace models {
	
	export class CertificateBase_string_ {
	    address: string;
	    age: string;
	    birthDate: string;
	    certificateName: string;
	    name: string;
	    publicAddress: string;
	    uniqueId: string;
	    extra: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new CertificateBase_string_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = source["address"];
	        this.age = source["age"];
	        this.birthDate = source["birthDate"];
	        this.certificateName = source["certificateName"];
	        this.name = source["name"];
	        this.publicAddress = source["publicAddress"];
	        this.uniqueId = source["uniqueId"];
	        this.extra = source["extra"];
	    }
	}
	export class Document {
	    shahash: string;
	    encryptedDocument: number[];
	    publicAddress: string;
	
	    static createFrom(source: any = {}) {
	        return new Document(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.shahash = source["shahash"];
	        this.encryptedDocument = source["encryptedDocument"];
	        this.publicAddress = source["publicAddress"];
	    }
	}

}

