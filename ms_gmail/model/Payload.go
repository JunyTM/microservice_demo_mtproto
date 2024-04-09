package model

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegistPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

type UserData struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type MTprotoPayload struct {
	Message string `json:"message"`
}
