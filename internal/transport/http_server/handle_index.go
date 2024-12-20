package http_server

import (
	"fmt"
	"io"

	"net/http"
	// "internal/app"
	// "internal/ports/storage"
)

const indexPageTpl = `<html>
	<head>
		<title>%s</title>
		<style>
	  		table, td, th {
	    		border: 1px solid black;
	    		border-spacing: 0px;
	  		}
	  		td, th {
	    		padding: 5px;
	  		}
		</style>
	</head>
	<body>
	   	%s
	</body>
	</html>`

const indexTableTpl = "<table>%s</table>"
const indexTableHeaderTpl = "<tr><th>%s</th><th>%v</th></tr>"
const indexTableRowTpl = "<tr><td>%s</td><td style=\"text-align: right;\">%v</td></tr>"

// index page handler
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	//check for malformed requests - only exact root path accepted
	//Important: covered by tests, removal will bring tests to fail
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// set correct data type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	htmlBody := fmt.Sprintf(indexTableHeaderTpl, "Metric", "Value")

	// metrics := app.Stor.GetMetrics()
	// for _, key := range storage.SortKeys(metrics) {
	// 	metric, _ := app.Stor.GetMetric(key)
	// 	htmlBody += fmt.Sprintf(indexTableRowTpl, key, metric.GetValue())
	// }
	htmlBody = fmt.Sprintf(indexTableTpl, htmlBody)

	io.WriteString(w, fmt.Sprintf(indexPageTpl, "Metrics", htmlBody))
}
