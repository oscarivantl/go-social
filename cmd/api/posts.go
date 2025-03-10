package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/oscarivantl/go-social/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title" validator:"required,max=100"`
	Content string   `json:"content" validator:"required,max=1000"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		//writeJSONError(w, http.StatusBadRequest, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// TODO: Change after auth
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		//writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}

	if err := writeJSON(w, http.StatusCreated, post); err != nil {
		//writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		//writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	post, err := app.store.Posts.GetById(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			//writeJSONError(w, http.StatusNotFound, err.Error())
			app.notFoundResponse(w, r, err)
		default:
			//writeJSONError(w, http.StatusInternalServerError, err.Error())
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := writeJSON(w, http.StatusOK, post); err != nil {
		//writeJSONError(w, http.StatusInternalServerError, err.Error())
		app.internalServerError(w, r, err)
		return
	}
}
