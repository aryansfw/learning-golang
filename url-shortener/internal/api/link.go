package api

import (
	"context"
	"encoding/json"
	"net/http"
	"url-shortener/internal/domain"
)

type createListRequest struct {
	ShortURL string `json:"short_url"`
	LongURL  string `json:"long_url"`
}

type listCreatedResponse struct {
	ID int64 `json:"id"`
}

func (a *api) redirectLinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	url := r.PathValue("url")

	link, err := a.linkRepo.GetByShortURL(ctx, url)
	if err != nil {
		a.errorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	http.Redirect(w, r, link.LongURL, http.StatusFound)
}

func (a *api) createLinkHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	lr := &createListRequest{}

	if err := json.NewDecoder(r.Body).Decode(&lr); err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	link := domain.Link{
		ShortURL: lr.ShortURL,
		LongURL:  lr.LongURL,
	}

	if err := a.linkRepo.Create(ctx, &link); err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.successResponse(w, listCreatedResponse{ID: link.ID})
}
