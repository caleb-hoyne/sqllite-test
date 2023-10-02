package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/caleb-hoyne/slogctx"
	db "github.com/caleb-hoyne/sqllite-test/repository"
	"log/slog"
	"net/http"
	"strconv"
)

type GetNameResp struct {
	Name string `json:"name"`
}

type PostNameReq struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Repository interface {
	GetNameByID(id int) (string, error)
	StoreUser(id int, name string) error
}

type RequestHandler struct {
	R Repository
}

func (r *RequestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.handleGet(req, w)
	case http.MethodPost:
		r.handlePost(req, w)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *RequestHandler) handlePost(req *http.Request, w http.ResponseWriter) {
	var postReq PostNameReq
	err := json.NewDecoder(req.Body).Decode(&postReq)
	if err != nil {
		r.handleError(req.Context(), err, http.StatusBadRequest, w)
		return
	}
	ctx := slogctx.AddValues(req.Context(), slog.Int("id", postReq.ID))
	err = r.R.StoreUser(postReq.ID, postReq.Name)
	if err != nil {
		if errors.Is(err, db.ErrIDAlreadyExists) {
			r.handleError(ctx, err, http.StatusConflict, w)
			return
		}
		r.handleError(ctx, err, http.StatusInternalServerError, w)
		return
	}
	slog.InfoContext(ctx, "Stored name")
	w.WriteHeader(http.StatusCreated)
}

func (r *RequestHandler) handleGet(req *http.Request, w http.ResponseWriter) {
	id := req.URL.Query().Get("id")
	idStr, err := strconv.Atoi(id)
	ctx := req.Context()
	if err != nil {
		r.handleError(ctx, err, http.StatusBadRequest, w)
		return
	}
	ctx = slogctx.AddValues(req.Context(), slog.String("id", id))
	name, err := r.R.GetNameByID(idStr)
	if err != nil {
		r.handleError(ctx, err, http.StatusInternalServerError, w)
		return
	}
	if name == "" {
		r.handleError(ctx, errors.New("id not found"), http.StatusNotFound, w)
		return
	}
	slog.InfoContext(ctx, "Got name")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(name)
	if err != nil {
		panic(err)
	}
}

func (r *RequestHandler) handleError(ctx context.Context, err error, code int, w http.ResponseWriter) {
	slog.ErrorContext(ctx, err.Error())
	w.WriteHeader(code)
}
