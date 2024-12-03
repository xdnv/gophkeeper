package http_server

import (
	"net/http"
)

// HTTP single metric update v2 processing
func HandleUpdateMetricV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// data, hs := app.UpdateMetricV2(r.Body)
	// if hs.Err != nil {
	// 	logger.Error("handleUpdateMetricV2: " + hs.Message)
	// 	http.Error(w, hs.Message, hs.HTTPStatus)
	// 	return
	// }

	// w.WriteHeader(http.StatusOK)
	// w.Write(*data)
}
