package domain

// GameInfo represents information about an RPG Maker game
type GameInfo struct {
	GameDir   string `json:"gameDir"`
	ExePath   string `json:"exePath"`
	DataPath  string `json:"dataPath"`
	JsPath    string `json:"jsPath"`
	ImgPath   string `json:"imgPath"`
	GameTitle string `json:"gameTitle"`
}

// LocatedGame represents a game stored in the user's collection
type LocatedGame struct {
	Id           string   `json:"id"`
	GameDir      string   `json:"gameDir"`
	ExePath      string   `json:"exePath"`
	RJCode       string   `json:"rjCode"`
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
	Title           string `json:"title"`
	RJCode          string `json:"rjCode"`
	StoreLink       string `json:"storeLink"`
	ReleaseDate     string `json:"releaseDate"`
	SystemGameTitle string `json:"systemGameTitle"`
	PatchDownloadId string `json:"patchDownloadId"`
}

// PatchSummary records which files were patched during the patch process
type PatchSummary struct {
	PatchedAt    string   `json:"patchedAt"`    // ISO timestamp
	PatchedFiles []string `json:"patchedFiles"` // Relative paths from game directory
}
