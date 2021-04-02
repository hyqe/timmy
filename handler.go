package timmy

import (
	"net/http"
	"strings"
)

func Handler(tables ...Tabler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var matched bool
		for _, t := range tables {
			if strings.HasPrefix(r.URL.String(), t.Name()) {
				matched = true
				switch r.Method {
				case http.MethodPut:
				case http.MethodPatch:
				case http.MethodPost:
				case http.MethodGet:
				case http.MethodDelete:
				default:
				}
				break
			}
		}
		if !matched {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}
