package router

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
)

// registerFrontend serves the embedded frontend for every path that does not
// match an API route. Unknown paths fall back to index.html, so client side
// routes such as /main/key survive a page refresh.
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
		name := strings.TrimPrefix(path.Clean(c.Request.URL.Path), "/")
		if name != "" && fs.ValidPath(name) {
			if info, err := fs.Stat(frontendFS, name); err == nil && !info.IsDir() {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
		}
		//SPA fallback (index.html is returned directly because http.FileServer
		//redirects "/index.html" back to "/")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
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
