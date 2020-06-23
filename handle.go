package gsd

import (
	"io"
	"net/http"
	"strings"

	"google.golang.org/api/googleapi"
)

var (
	blacklistedHeaders = make(map[string]bool)
)

func init() {
	blacklistedHeaders["host"] = true
	blacklistedHeaders["accept-encoding"] = true
}

// HandleFunc initiates srv and serves.
func HandleFunc(pattern, credentialsFile, tokenFile string) {
	initSrv(credentialsFile, tokenFile)
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET", "HEAD":
		default:
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		paths := strings.Split(r.URL.Path, "/")
		if len(paths) == 1 {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if paths[0] == "" {
			paths = paths[1:]
		}
		handle(w, paths[0], r.Method, r.Header)
	})
}

func handle(w http.ResponseWriter, id, method string, rheader http.Header) {
	fileCall := srv.Files.Get(id).SupportsAllDrives(true)
	for k, vs := range rheader {
		if blacklistedHeaders[strings.ToLower(k)] {
			continue
		}
		for _, v := range vs {
			fileCall.Header().Add(k, v)
		}
	}
	resp, err := fileCall.Download()
	if err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			http.Error(w, e.Message, e.Code)
		} else {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		return
	}
	for k, vs := range resp.Header {
		if strings.ToLower(k) == "host" {
			continue
		}
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	defer resp.Body.Close()
	if method == "HEAD" {
		return
	}

	_, _ = io.Copy(w, resp.Body)
}
