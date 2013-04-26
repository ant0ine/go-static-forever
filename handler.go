// Serve files that never change
//
// If you can put a version string (commitish for instance) in the path of your
// static files, then the content served by the corresponding URLs is guaranteed
// to never change. A whole set of optimisations become possible.
//
// * If the request contains If-Modified-Since, return 304 without checking anything
//
// * Set the Expires and Last-Modified headers to <forever>
//
// * Set the Cache-Control header to "public; max-age=<forever>; s-maxage=<forever>"
//
// * Set the Etag header to the full file path ? TODO
//
// This handler is implemented as a wrapper around http.FileServer, and when the
// isDevelopment flag is set, http.FileServer is used directly.
//
// Example:
//
//        package main
//
//        import(
//                "github.com/ant0ine/go-static-forever"
//                "net/http"
//        )
//
//        handler := forever.NewStaticHandler(
//                http.Dir("/static/"),   // FileSytem to serve
//                "1234567"               // version string, like a commitish for instance
//                nil,                    // "forever duration" default to 100 years
//                false,                  // isDevelopement
//        )
//
//        http.ListenAndServe(":8080", &handler)
//
package forever

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type staticHandler struct {
	fileHandler     http.Handler
	versionPrefix   string
	foreverHttpDate string
	deltaSeconds    int
	isDevelopment   bool
}

// borrowed from net/http/server.go
const timeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

// foreverDuration default to 100 years.
func NewStaticHandler(
	root http.FileSystem,
	version string,
	foreverDuration *time.Duration,
	isDevelopment bool) http.Handler {

	// set the default
	if foreverDuration == nil {
		dur := time.Duration(100 * 365 * 86400 * time.Second)
		foreverDuration = &dur
	}

	deltaSeconds := int(foreverDuration.Seconds())
	forever := time.Now().Add(*foreverDuration)
	foreverHttpDate := forever.Format(timeFormat)

	prefix := "/" + version

	return &staticHandler{
		fileHandler:     http.FileServer(root),
		versionPrefix:   prefix,
		foreverHttpDate: foreverHttpDate,
		deltaSeconds:    deltaSeconds,
		isDevelopment:   isDevelopment,
	}
}

func (self *staticHandler) ServeHTTP(origWriter http.ResponseWriter, origRequest *http.Request) {

	if !strings.HasPrefix(origRequest.URL.Path, self.versionPrefix) {
		http.NotFound(origWriter, origRequest)
		return
	}

	origRequest.URL.Path = origRequest.URL.Path[len(self.versionPrefix):]

	if self.isDevelopment {
		self.fileHandler.ServeHTTP(origWriter, origRequest)
		return
	}

	// If the request contains If-Modified-Since, return 304 without checking anything
	if origRequest.Header.Get("If-Modified-Since") != "" {
		http.Error(origWriter, "Not Modified", http.StatusNotModified)
		return
	}

	// Provide writer wrapper to write the custom headers only when the response code is 200
	writer := &responseWriter{
		origWriter,
		self,
		false,
	}

	self.fileHandler.ServeHTTP(writer, origRequest)
}

// Inherit from an object implementing the http.ResponseWriter interface
type responseWriter struct {
	http.ResponseWriter
	handler     *staticHandler
	wroteHeader bool
}

// Overloading of the http.ResponseWriter method.
func (self *responseWriter) WriteHeader(code int) {
	if code == 200 {
		// Cache forever
		self.Header().Set("Expires", self.handler.foreverHttpDate)
		self.Header().Set("Last-Modified", self.handler.foreverHttpDate)
		self.Header().Set("Cache-Control", fmt.Sprintf(
			"public; max-age=%d; s-maxage=%d",
			self.handler.deltaSeconds,
		))
		// TODO Set the Etag header to the full file path ?
	}
	self.ResponseWriter.WriteHeader(code)
	self.wroteHeader = true
}

// Overloading of the http.ResponseWriter method.
func (self *responseWriter) Write(b []byte) (int, error) {
	if !self.wroteHeader {
		self.WriteHeader(http.StatusOK)
	}
	return self.ResponseWriter.Write(b)
}
