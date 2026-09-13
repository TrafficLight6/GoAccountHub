package router

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// frontendRoutes are the client side routes declared in
// GAHFrontend/src/router/index.js. They are served with HTTP 200; every other
// path that is neither an API route nor a real file is a 404.
// Keep this list in sync with the frontend router.
var frontendRoutes = map[string]struct{}{
	"/":               {},
	"/component":      {},
	"/main":           {},
	"/main/home":      {},
	"/main/admin":     {},
	"/main/user":      {},
	"/main/character": {},
	"/main/key":       {},
}

// registerFrontend serves the embedded frontend for every path that does not
// match an API route. The frontend routes above are answered with HTTP 200,
// unknown paths are answered with HTTP 404; both return index.html so that the
// frontend can render the page or its own 404 page.
func registerFrontend(router *gin.Engine, frontendFS fs.FS) {
	if frontendFS == nil {
		router.NoRoute(frontendMissing())
		return
	}
	indexHTML, err := fs.ReadFile(frontendFS, "index.html")
	if err != nil {
		router.NoRoute(frontendMissing())
		return
	}
	fileServer := http.FileServer(http.FS(frontendFS))
	router.NoRoute(func(c *gin.Context) {
		//Keep API 404 as JSON instead of handing out the SPA
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "error": "Not Found"})
			return
		}
		//index.html is returned directly because http.FileServer redirects
		//"/index.html" back to "/"
		serveIndex := func(status int) {
			c.Data(status, "text/html; charset=utf-8", indexHTML)
		}
		frontendPath := path.Clean(c.Request.URL.Path)
		name := strings.TrimPrefix(frontendPath, "/")
		if name != "" && fs.ValidPath(name) {
			if info, err := fs.Stat(frontendFS, name); err == nil && !info.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}
		if _, ok := frontendRoutes[frontendPath]; ok {
			serveIndex(http.StatusOK)
			return
		}
		//Unknown client side route: let the frontend show its own 404 page
		serveIndex(http.StatusNotFound)
	})
}

// frontendMissing explains how to embed the frontend when it was not built
// before the server was compiled.
func frontendMissing() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":  http.StatusNotFound,
			"error": "Frontend is not embedded, run `npm run build` in GAHFrontend and rebuild the server",
		})
	}
}
