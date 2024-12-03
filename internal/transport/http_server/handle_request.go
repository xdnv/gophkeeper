package http_server

import (
	"net/http"
)

// HTTP single metric request v2 processing
func HandleRequestMetricV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// data, hs := app.RequestMetricV2(r.Body)
	// if hs.Err != nil {
	// 	logger.Error("handleRequestMetricV2: " + hs.Message)
	// 	http.Error(w, hs.Message, hs.HTTPStatus)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// w.Write(*data)
}
