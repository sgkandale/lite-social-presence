package server

type HealthResponse struct {
	Status string `json:"status"`
}

type GeneralResponse struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	Name string `json:"name"`
}

type LoginRequest struct {
	Name string `json:"name"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
