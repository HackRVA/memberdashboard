package routes

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/HackRVA/memberserver/pkg/membermgr/ui"
	"github.com/sirupsen/logrus"
)

type spaRouter struct {
	fs.ReadFileFS
}

func (s spaRouter) isStaticFile(path string) bool {
	web, _ := fs.Sub(s.ReadFileFS, "web")
	p := strings.TrimLeft(filepath.Clean(path), "/")
	_, err := fs.Stat(web, p)

	// todo: determine actual err
	return err == nil
}

func (s spaRouter) notFoundHandler(w http.ResponseWriter, _ *http.Request) {
	index, _ := s.ReadFileFS.ReadFile("web/index.html")
	if _, err := w.Write(index); err != nil {
		logrus.Error(err)
	}
}

func (s spaRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	web, _ := fs.Sub(ui.UI, "web")

	if !s.isStaticFile(r.URL.Path) {
		s.notFoundHandler(w, r)
		return
	}

	http.FileServer(http.FS(web)).ServeHTTP(w, r)
}
