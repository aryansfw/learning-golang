package api

import (
	"context"
	"ecommerce/internal/domain"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type createdProductResponse struct {
	ID int64 `json:"id"`
}

func (a *api) createProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	userID := ctx.Value(domain.UserContextKey).(int64)

	prd := &domain.Product{
		UserID: userID,
	}

	if err := json.NewDecoder(r.Body).Decode(&prd); err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := a.productRepo.Create(ctx, prd); err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	a.successResponse(w, createdProductResponse{ID: prd.ID})
}

func (a *api) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/product/"), 10, 64)
	if err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := ctx.Value(domain.UserContextKey).(int64)

	prd := &domain.Product{
		ID:     id,
		UserID: userID,
	}

	if err := json.NewDecoder(r.Body).Decode(&prd); err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := a.productRepo.Update(ctx, prd); err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	a.successResponse(w, createdProductResponse{ID: prd.ID})
}

func (a *api) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/product/"), 10, 64)
	if err != nil {
		a.errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID := ctx.Value(domain.UserContextKey).(int64)

	prd, err := a.productRepo.GetByID(ctx, id)
	if err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	prd.UserID = userID
	if err := a.productRepo.Delete(ctx, prd); err != nil {
		a.errorResponse(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	a.successResponse(w, nil)
}
