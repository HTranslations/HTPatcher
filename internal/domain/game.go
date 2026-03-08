package domain

// GameInfo represents information about an RPG Maker game
type GameInfo struct {
	GameDir     string `json:"gameDir"`
	ExePath     string `json:"exePath"`
	DataPath    string `json:"dataPath"`
	JsPath      string `json:"jsPath"`
	ImgPath     string `json:"imgPath"`
	GameTitle   string `json:"gameTitle"`
	GameVersion string `json:"gameVersion"` // "mv", "mz", "vxace", or "" for unknown
}

// LocatedGame represents a game stored in the user's collection
type LocatedGame struct {
	Id           string   `json:"id"`
	GameDir      string   `json:"gameDir"`
	ExePath      string   `json:"exePath"`
	RJCode       string   `json:"rjCode"`
	StoreCode    string   `json:"storeCode"`
	FriendlyName string   `json:"friendlyName"`
	Tags         []string `json:"tags"`
	Translated   bool     `json:"translated"`
	Pinned       bool     `json:"pinned"`
	PlayStatus   string   `json:"playStatus"` // "unplayed", "playing", "on-hold", "finished", "given-up"
}

// PersistentData holds user's persistent application data
type PersistentData struct {
	LocatedGames []LocatedGame `json:"locatedGames"`
	GamesPerRow  int           `json:"gamesPerRow"` // 3 or 4, defaults to 3
}

// PatchEntry represents a patch available for download
type PatchEntry struct {
	Store             string         `json:"store"`
	StoreCode         string         `json:"storeCode"`
	Title             string         `json:"title"`
	Slug              string         `json:"slug"`
	ThumbnailId       string         `json:"thumbnailId"`
	ThumbnailFileName string         `json:"thumbnailFileName"`
	StoreLink         string         `json:"storeLink"`
	ReleaseDate       string         `json:"releaseDate"`
	Patch             *PatchDownload `json:"patch"`
}

// PatchDownload represents the downloadable patch file
type PatchDownload struct {
	Url      string `json:"url"`
	FileName string `json:"fileName"`
	FileSize int    `json:"fileSize"`
	Version  string `json:"version"`
}

// GamePatchInfo represents the response from GET /api/patches/{storeCode}
type GamePatchInfo struct {
	Store     string           `json:"store"`
	StoreCode string          `json:"storeCode"`
	Status    string           `json:"status"`
	Title     string           `json:"title"`
	Slug      string           `json:"slug"`
	StoreLink string           `json:"storeLink"`
	Patches   []GamePatchEntry `json:"patches"`
}

// GamePatchEntry represents a single patch version in the GamePatchInfo response
type GamePatchEntry struct {
	Version      string         `json:"version"`
	ReleaseNotes []string       `json:"releaseNotes"`
	Download     *PatchDownload `json:"download"`
}

// PatchSummary records which files were patched during the patch process
type PatchSummary struct {
	PatchedAt    string   `json:"patchedAt"`    // ISO timestamp
	PatchedFiles []string `json:"patchedFiles"` // Relative paths from game directory
}
