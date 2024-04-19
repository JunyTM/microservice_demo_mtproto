package controller

import (
	"ms_gmail/infrastructure"
	"ms_gmail/model"
	"net/http"

	"github.com/go-chi/render"
)

var serverHost string = infrastructure.GetServerHost()

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Bad Request", http.StatusBadRequest)
	w.WriteHeader(http.StatusBadRequest)
	res := &model.Response{
		Data:    nil,
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, res)
}

func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, "Not Found", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	res := &model.Response{
		Data:    nil,
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, res)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	// http.Error(w, "Internal Server Error\n", http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	res := &model.Response{
		Data:    nil,
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, res)
}
