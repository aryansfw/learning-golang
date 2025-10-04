package main

import (
	"encoding/json"
	"net/http"

	"github.com/patrickmn/go-cache"
)

func WeatherHandler(c *cache.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		result, found := c.Get("weather")
		if !found {
			client := new(http.Client)

			request, err := http.NewRequest(http.MethodGet,
				"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/-7.2615005,112.7799616?key=8P9L7NM5KA9FPWLYQWK2M3LDK", nil)
			if err != nil {
				http.Error(w, "Creating Request Error", http.StatusBadRequest)
				return
			}

			res, err := client.Do(request)

			if err != nil {
				http.Error(w, "Fetching Request Error", http.StatusBadRequest)
				return
			}
			defer res.Body.Close()

			if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
				http.Error(w, "Decode Error", http.StatusBadRequest)
				return
			}
			c.Set("weather", result, cache.DefaultExpiration)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(result); err != nil {
			http.Error(w, "Encode Error", http.StatusBadRequest)
			return
		}
	}
}
