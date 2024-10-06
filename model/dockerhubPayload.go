package model

// DockerHubPayload represents the payload sent by Docker Hub
type DockerHubPayload struct {
	CallbackURL string `json:"callback_url"`
	PushData    struct {
		PushedAt int64  `json:"pushed_at"`
		Pusher   string `json:"pusher"`
		Tag      string `json:"tag"`
	} `json:"push_data"`
	Repository struct {
		CommentCount int    `json:"comment_count"`
		DateCreated  int    `json:"date_created"`
		Description  string `json:"description"`
		IsOfficial   bool   `json:"is_official"`
		IsPrivate    bool   `json:"is_private"`
		IsTrusted    bool   `json:"is_trusted"`
		Name         string `json:"name"`
		Namespace    string `json:"namespace"`
		Owner        string `json:"owner"`
		RepoName     string `json:"repo_name"`
		RepoURL      string `json:"repo_url"`
		StarCount    int    `json:"star_count"`
		Status       string `json:"status"`
	} `json:"repository"`
}
