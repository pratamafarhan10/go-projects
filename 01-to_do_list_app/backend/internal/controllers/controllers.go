package controllers

type ErrorResponse struct {
	Errors map[string][]string
}

type SuccessResponse struct {
	Data interface{}
}
