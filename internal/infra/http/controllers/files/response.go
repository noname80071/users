package files

type UploadResponse struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
	Url     string `json:"url"`
}

type Response struct {
	Message string `json:"message"`
}

type GetSkinResponse struct {
	Skin string `json:"skin"`
}

type GetCloakResponse struct {
	Cloak string `json:"cloak"`
}
