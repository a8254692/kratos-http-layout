package monitor

import (
	"net/http"
)

func init() {
	HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("UP"))
	})
}
