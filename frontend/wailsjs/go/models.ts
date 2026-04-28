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
	
	export class LeafFields {
	    hash: string;
	    key: string;
	    value: any;
	    salt: string;
	
	    static createFrom(source: any = {}) {
	        return new LeafFields(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.key = source["key"];
	        this.value = source["value"];
	        this.salt = source["salt"];
	    }
	}
	export class CertificateBase_github_com_Suy56_ProofChain_internal_models_LeafFields_ {
	    address: LeafFields;
	    age: LeafFields;
	    birthDate: LeafFields;
	    certificateName: LeafFields;
	    name: LeafFields;
	    publicAddress: LeafFields;
	    uniqueId: LeafFields;
	    extra: Record<string, LeafFields>;
	
	    static createFrom(source: any = {}) {
	        return new CertificateBase_github_com_Suy56_ProofChain_internal_models_LeafFields_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.address = this.convertValues(source["address"], LeafFields);
	        this.age = this.convertValues(source["age"], LeafFields);
	        this.birthDate = this.convertValues(source["birthDate"], LeafFields);
	        this.certificateName = this.convertValues(source["certificateName"], LeafFields);
	        this.name = this.convertValues(source["name"], LeafFields);
	        this.publicAddress = this.convertValues(source["publicAddress"], LeafFields);
	        this.uniqueId = this.convertValues(source["uniqueId"], LeafFields);
	        this.extra = this.convertValues(source["extra"], LeafFields, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
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

