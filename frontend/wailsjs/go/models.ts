export namespace main {
	
	export class ParameterPathToPatch {
	    path: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new ParameterPathToPatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.type = source["type"];
	    }
	}
	export class ParameterToPatch {
	    plugin: string;
	    function: string;
	    rootType: string;
	    parameterPathsToPatch: ParameterPathToPatch[];
	
	    static createFrom(source: any = {}) {
	        return new ParameterToPatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.plugin = source["plugin"];
	        this.function = source["function"];
	        this.rootType = source["rootType"];
	        this.parameterPathsToPatch = this.convertValues(source["parameterPathsToPatch"], ParameterPathToPatch);
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
	export class Config {
	    variablesToPatch: number[];
	    wrapWidth: number;
	    version: number;
	    parametersToPatch: ParameterToPatch[];
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.variablesToPatch = source["variablesToPatch"];
	        this.wrapWidth = source["wrapWidth"];
	        this.version = source["version"];
	        this.parametersToPatch = this.convertValues(source["parametersToPatch"], ParameterToPatch);
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
	export class GameInfo {
	    gameDir: string;
	    exePath: string;
	    dataPath: string;
	    jsPath: string;
	
	    static createFrom(source: any = {}) {
	        return new GameInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gameDir = source["gameDir"];
	        this.exePath = source["exePath"];
	        this.dataPath = source["dataPath"];
	        this.jsPath = source["jsPath"];
	    }
	}
	
	
	export class PatchInfo {
	    patchPath: string;
	    dictionary: Record<string, string>;
	    config?: Config;
	
	    static createFrom(source: any = {}) {
	        return new PatchInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.patchPath = source["patchPath"];
	        this.dictionary = source["dictionary"];
	        this.config = this.convertValues(source["config"], Config);
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

