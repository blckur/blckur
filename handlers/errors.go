package handlers

type ErrorData struct {
	Error string `json:"error" binding:"required"`
	Message string `json:"error_msg" binding:"required"`
}
