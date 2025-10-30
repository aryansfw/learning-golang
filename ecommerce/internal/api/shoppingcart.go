package api

import (
	"context"
	"net/http"
)

func (a *api) addCartHandler(w http.ResponseWriter, r *http.Request) {
	_, cancel := context.WithCancel(r.Context())
	defer cancel()
}
