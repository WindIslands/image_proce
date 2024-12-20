export namespace app {
	
	export class FileInfo {
	    id: string;
	    name: string;
	    size: number;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.path = source["path"];
	    }
	}
	export class ImageTask {
	    name: string;
	    path: string;
	    outputFormat: string;
	    quality: number;
	    convertFormat: boolean;
	    width: number;
	    height: number;
	    keepOriginalSize: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImageTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.outputFormat = source["outputFormat"];
	        this.quality = source["quality"];
	        this.convertFormat = source["convertFormat"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.keepOriginalSize = source["keepOriginalSize"];
	    }
	}
	export class ProcessImagesOptions {
	    tasks: ImageTask[];
	    outputPath: string;
	    workers: number;
	
	    static createFrom(source: any = {}) {
	        return new ProcessImagesOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tasks = this.convertValues(source["tasks"], ImageTask);
	        this.outputPath = source["outputPath"];
	        this.workers = source["workers"];
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

}

