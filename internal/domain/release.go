package domain

type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Asset struct {
	URL                string `json:"url"`
	ID                 int    `json:"id"`
	NodeID             string `json:"node_id"`
	Name               string `json:"name"`
	Label              string `json:"label"`
	Uploader           User   `json:"uploader"`
	ContentType        string `json:"content_type"`
	State              string `json:"state"`
	Size               int    `json:"size"`
	Digest             string `json:"digest"`
	DownloadCount      int    `json:"download_count"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type ReleaseInfo struct {
	URL             string  `json:"url"`
	AssetsURL       string  `json:"assets_url"`
	UploadURL       string  `json:"upload_url"`
	HTMLURL         string  `json:"html_url"`
	ID              int     `json:"id"`
	Author          User    `json:"author"`
	NodeID          string  `json:"node_id"`
	TagName         string  `json:"tag_name"`
	TargetCommitish string  `json:"target_commitish"`
	Name            string  `json:"name"`
	Draft           bool    `json:"draft"`
	Immutable       bool    `json:"immutable"`
	Prerelease      bool    `json:"prerelease"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	PublishedAt     string  `json:"published_at"`
	Assets          []Asset `json:"assets"`
	TarballURL      string  `json:"tarball_url"`
	ZipballURL      string  `json:"zipball_url"`
	Body            string  `json:"body"`
}
