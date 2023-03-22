package domain

type AuthError struct {
	Messages []struct {
		Text     string `json:"text"`
		Severity string `json:"severity"`
	} `json:"messages"`
	FieldErrors []interface{} `json:"fieldErrors"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
