package handler

import "log"

// BadRequest HTTP 400
type BadRequest struct {
	Message string
}

func (b BadRequest) Error() string {
	return b.Message
}

// InternalServer HTTP 500
type InternalServer struct {
	Status int
	Message string
}

func (i InternalServer) Error() string {
	log.Print(i.Message)
	return "internal server error"
}
