package monitor

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"
)

var (
	// DefaultServeMux ...
	DefaultServeMux = mux.NewRouter()
)

func init() {
	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/allocs", pprof.Index)
	HandleFunc("/debug/pprof/block", pprof.Index)
	HandleFunc("/debug/pprof/goroutine", pprof.Index)
	HandleFunc("/debug/pprof/heap", pprof.Index)
	HandleFunc("/debug/pprof/mutex", pprof.Index)
	HandleFunc("/debug/pprof/threadcreate", pprof.Index)

	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// HandleFunc ...
func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	DefaultServeMux.HandleFunc(pattern, handler)
}

// Handle ...
func Handle(pattern string, handler http.Handler) {
	DefaultServeMux.Handle(pattern, handler)
}

// ListenAndServe ...
func ListenAndServe(addr string) error {
	svr := &http.Server{Handler: DefaultServeMux, Addr: addr}
	return svr.ListenAndServe()
}
