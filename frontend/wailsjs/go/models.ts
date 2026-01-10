export namespace domain {
	
	export class User {
	    login: string;
	    id: number;
	    node_id: string;
	    avatar_url: string;
	    gravatar_id: string;
	    url: string;
	    html_url: string;
	    followers_url: string;
	    following_url: string;
	    gists_url: string;
	    starred_url: string;
	    subscriptions_url: string;
	    organizations_url: string;
	    repos_url: string;
	    events_url: string;
	    received_events_url: string;
	    type: string;
	    user_view_type: string;
	    site_admin: boolean;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.login = source["login"];
	        this.id = source["id"];
	        this.node_id = source["node_id"];
	        this.avatar_url = source["avatar_url"];
	        this.gravatar_id = source["gravatar_id"];
	        this.url = source["url"];
	        this.html_url = source["html_url"];
	        this.followers_url = source["followers_url"];
	        this.following_url = source["following_url"];
	        this.gists_url = source["gists_url"];
	        this.starred_url = source["starred_url"];
	        this.subscriptions_url = source["subscriptions_url"];
	        this.organizations_url = source["organizations_url"];
	        this.repos_url = source["repos_url"];
	        this.events_url = source["events_url"];
	        this.received_events_url = source["received_events_url"];
	        this.type = source["type"];
	        this.user_view_type = source["user_view_type"];
	        this.site_admin = source["site_admin"];
	    }
	}
	export class Asset {
	    url: string;
	    id: number;
	    node_id: string;
	    name: string;
	    label: string;
	    uploader: User;
	    content_type: string;
	    state: string;
	    size: number;
	    digest: string;
	    download_count: number;
	    created_at: string;
	    updated_at: string;
	    browser_download_url: string;
	
	    static createFrom(source: any = {}) {
	        return new Asset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.id = source["id"];
	        this.node_id = source["node_id"];
	        this.name = source["name"];
	        this.label = source["label"];
	        this.uploader = this.convertValues(source["uploader"], User);
	        this.content_type = source["content_type"];
	        this.state = source["state"];
	        this.size = source["size"];
	        this.digest = source["digest"];
	        this.download_count = source["download_count"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.browser_download_url = source["browser_download_url"];
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
	export class PluginReplaceRule {
	    match: string;
	    replace: string;
	
	    static createFrom(source: any = {}) {
	        return new PluginReplaceRule(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.match = source["match"];
	        this.replace = source["replace"];
	    }
	}
	export class PluginToPatch {
	    plugin: string;
	    parametersPatchScript: string;
	    replaceRules: PluginReplaceRule[];
	
	    static createFrom(source: any = {}) {
	        return new PluginToPatch(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.plugin = source["plugin"];
	        this.parametersPatchScript = source["parametersPatchScript"];
	        this.replaceRules = this.convertValues(source["replaceRules"], PluginReplaceRule);
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
	    pluginsToPatch: PluginToPatch[];
	    creditsLocation: string;
	    dynamicWrapWidth: boolean;
	    locale: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.variablesToPatch = source["variablesToPatch"];
	        this.wrapWidth = source["wrapWidth"];
	        this.version = source["version"];
	        this.parametersToPatch = this.convertValues(source["parametersToPatch"], ParameterToPatch);
	        this.pluginsToPatch = this.convertValues(source["pluginsToPatch"], PluginToPatch);
	        this.creditsLocation = source["creditsLocation"];
	        this.dynamicWrapWidth = source["dynamicWrapWidth"];
	        this.locale = source["locale"];
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
	    imgPath: string;
	    gameTitle: string;
	
	    static createFrom(source: any = {}) {
	        return new GameInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.gameDir = source["gameDir"];
	        this.exePath = source["exePath"];
	        this.dataPath = source["dataPath"];
	        this.jsPath = source["jsPath"];
	        this.imgPath = source["imgPath"];
	        this.gameTitle = source["gameTitle"];
	    }
	}
	export class LocatedGame {
	    id: string;
	    gameDir: string;
	    exePath: string;
	    rjCode: string;
	    friendlyName: string;
	    tags: string[];
	    translated: boolean;
	    pinned: boolean;
	    playStatus: string;
	
	    static createFrom(source: any = {}) {
	        return new LocatedGame(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.gameDir = source["gameDir"];
	        this.exePath = source["exePath"];
	        this.rjCode = source["rjCode"];
	        this.friendlyName = source["friendlyName"];
	        this.tags = source["tags"];
	        this.translated = source["translated"];
	        this.pinned = source["pinned"];
	        this.playStatus = source["playStatus"];
	    }
	}
	
	
	export class PatchEntry {
	    title: string;
	    rjCode: string;
	    storeLink: string;
	    releaseDate: string;
	    systemGameTitle: string;
	    patchDownloadId: string;
	
	    static createFrom(source: any = {}) {
	        return new PatchEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.rjCode = source["rjCode"];
	        this.storeLink = source["storeLink"];
	        this.releaseDate = source["releaseDate"];
	        this.systemGameTitle = source["systemGameTitle"];
	        this.patchDownloadId = source["patchDownloadId"];
	    }
	}
	export class PatchInfo {
	    patchPath: string;
	    dictionary: Record<string, string>;
	    overrides: string[];
	    config?: Config;
	
	    static createFrom(source: any = {}) {
	        return new PatchInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.patchPath = source["patchPath"];
	        this.dictionary = source["dictionary"];
	        this.overrides = source["overrides"];
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
	
	
	export class ReleaseInfo {
	    url: string;
	    assets_url: string;
	    upload_url: string;
	    html_url: string;
	    id: number;
	    author: User;
	    node_id: string;
	    tag_name: string;
	    target_commitish: string;
	    name: string;
	    draft: boolean;
	    immutable: boolean;
	    prerelease: boolean;
	    created_at: string;
	    updated_at: string;
	    published_at: string;
	    assets: Asset[];
	    tarball_url: string;
	    zipball_url: string;
	    body: string;
	
	    static createFrom(source: any = {}) {
	        return new ReleaseInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.assets_url = source["assets_url"];
	        this.upload_url = source["upload_url"];
	        this.html_url = source["html_url"];
	        this.id = source["id"];
	        this.author = this.convertValues(source["author"], User);
	        this.node_id = source["node_id"];
	        this.tag_name = source["tag_name"];
	        this.target_commitish = source["target_commitish"];
	        this.name = source["name"];
	        this.draft = source["draft"];
	        this.immutable = source["immutable"];
	        this.prerelease = source["prerelease"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.published_at = source["published_at"];
	        this.assets = this.convertValues(source["assets"], Asset);
	        this.tarball_url = source["tarball_url"];
	        this.zipball_url = source["zipball_url"];
	        this.body = source["body"];
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

