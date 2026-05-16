export namespace booklet {
	
	export class Options {
	    Input: string;
	    Output: string;
	    N: number;
	    FormSize: string;
	    Guides: boolean;
	    Margin: number;
	    Binding: string;
	    BType: string;
	    Multifolio: boolean;
	    FolioSize: number;
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Input = source["Input"];
	        this.Output = source["Output"];
	        this.N = source["N"];
	        this.FormSize = source["FormSize"];
	        this.Guides = source["Guides"];
	        this.Margin = source["Margin"];
	        this.Binding = source["Binding"];
	        this.BType = source["BType"];
	        this.Multifolio = source["Multifolio"];
	        this.FolioSize = source["FolioSize"];
	    }
	}

}

