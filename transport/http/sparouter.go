package http

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type spaRouter struct {
	fs.FS
}

func (s spaRouter) isStaticFile(path string) bool {
	p := strings.TrimLeft(filepath.Clean(path), "/")
	_, err := fs.Stat(s.FS, p)

	// todo: determine actual err
	return err == nil
}

func (s spaRouter) notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	index, err := fs.ReadFile(s.FS, "index.html")
	if err != nil {
		logrus.Error("Failed to read index.html:", err)
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(index)
	if err != nil {
		logrus.Error(err)
	}
}

func (s spaRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.isStaticFile(r.URL.Path) {
		s.notFoundHandler(w, r)
		return
	}

	http.FileServer(http.FS(s.FS)).ServeHTTP(w, r)
}
