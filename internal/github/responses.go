package github

type AppDetails struct {
	AppId        int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type AppInstallation struct {
	Id              int    `json:"id"`
	TargetId        int    `json:"target_id"`
	TargetType      string `json:"target_type"`
	AccessTokensUrl string `json:"access_tokens_url"`
	Account         struct {
		Id    int    `json:"id"`
		Login string `json:"login"`
		Type  string `json:"type"`
	} `json:"account"`
}

type AppToken struct {
	Token       string            `json:"token"`
	ExpiresAt   string            `json:"expires_at"`
	Permissions map[string]string `json:"permissions"`
}
