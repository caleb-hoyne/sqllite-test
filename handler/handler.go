package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type GetNameResp struct {
	Name string `json:"name"`
}

type Repository interface {
	GetNameByID(id int) (string, error)
}

type RequestHandler struct {
	R Repository
}

func (r *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.handleGet(req, w)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (r *RequestHandler) handleGet(req *http.Request, w http.ResponseWriter) {
	id := req.URL.Query().Get("id")
	idStr, err := strconv.Atoi(id)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	name, err := r.R.GetNameByID(idStr)
	if err != nil {
		slog.ErrorContext(req.Context(), err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(name)
	if name == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(name)
	if err != nil {
		panic(err)
	}
}
