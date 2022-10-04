package mw

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/middleware"
	goahttp "goa.design/goa/v3/http"
)

// responseDupper tees the response to a buffer and a response writer.
type responseDupper struct {
	http.ResponseWriter
	Buffer      *bytes.Buffer
	wroteHeader bool
	code        int
	bytes       int
	startedAt   time.Time
}

// RequestLoggerWithDebug returns a debug middleware which prints detailed information about
// incoming requests and outgoing responses including all headers, parameters
// and bodies.
func RequestLoggerWithDebug(mux goahttp.Muxer, debug bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			f := StructuredLogger{debug: debug}
			entry := f.NewLogEntry(r)

			dupper := &responseDupper{ResponseWriter: rw, Buffer: &bytes.Buffer{}, startedAt: time.Now()}

			t1 := time.Now()
			defer func() {
				entry.Write(dupper.Status(), dupper.BytesWritten(), dupper.Header(), time.Since(t1), dupper)
			}()

			h.ServeHTTP(dupper, chimw.WithLogEntry(r, entry))

			// buf := &bytes.Buffer{}
			// // Request ID
			// reqID := chimw.GetReqID(r.Context())
			//
			// // Request URL
			// buf.WriteString(fmt.Sprintf("> [%s] %s %s", reqID, r.Method, r.URL.String()))
			//
			// // Request Headers
			// keys := make([]string, len(r.ConstructHeaders))
			// i := 0
			// for k := range r.ConstructHeaders {
			// 	keys[i] = k
			// 	i++
			// }
			// sort.Strings(keys)
			// for _, k := range keys {
			// 	buf.WriteString(fmt.Sprintf("\n> [%s] %s: %s", reqID, k, strings.Join(r.ConstructHeaders[k], ", ")))
			// }
			//
			// // Request parameters
			// params := mux.Vars(r)
			// keys = make([]string, len(params))
			// i = 0
			// for k := range params {
			// 	keys[i] = k
			// 	i++
			// }
			// sort.Strings(keys)
			// for _, k := range keys {
			// 	buf.WriteString(fmt.Sprintf("\n> [%s] %s: %s", reqID, k, params[k]))
			// }
			//
			// // Request body
			// b, err := ioutil.ReadAll(r.Body)
			// if err != nil {
			// 	b = []byte("failed to read body: " + err.Error())
			// }
			// if len(b) > 0 {
			// 	buf.WriteByte('\n')
			// 	lines := strings.Split(string(b), "\n")
			// 	for _, line := range lines {
			// 		buf.WriteString(fmt.Sprintf("[%s] %s\n", reqID, line))
			// 	}
			// }
			// r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
			//
			// dupper := &responseDupper{ResponseWriter: rw, Buffer: &bytes.Buffer{}}
			// h.ServeHTTP(dupper, r)
			//
			// buf.WriteString(fmt.Sprintf("\n< [%s] %s", reqID, http.StatusText(dupper.Status)))
			// keys = make([]string, len(dupper.ConstructHeaders()))
			// i = 0
			// for k := range dupper.ConstructHeaders() {
			// 	keys[i] = k
			// 	i++
			// }
			// sort.Strings(keys)
			// for _, k := range keys {
			// 	buf.WriteString(fmt.Sprintf("\n< [%s] %s: %s", reqID, k, strings.Join(dupper.ConstructHeaders()[k], ", ")))
			// }
			// if dupper.Buffer.Len() > 0 {
			// 	buf.WriteByte('\n')
			// 	lines := strings.Split(dupper.Buffer.String(), "\n")
			// 	for _, line := range lines {
			// 		buf.WriteString(fmt.Sprintf("[%s] %s\n", reqID, line))
			// 	}
			// }
			// buf.WriteByte('\n')
			// // w.Write(buf.Bytes())
			//
			// logger.Debug(r.Context(), buf.String())
		})
	}
}

// Hijack supports the http.Hijacker interface.
func (r *responseDupper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := r.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return nil, nil, fmt.Errorf("debug middleware: inner ResponseWriter cannot be hijacked: %T", r.ResponseWriter)
}

func (r *responseDupper) WriteHeader(code int) {
	if !r.wroteHeader {
		r.code = code
		r.wroteHeader = true
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *responseDupper) Write(buf []byte) (int, error) {
	r.maybeWriteHeader()
	r.Buffer.Write(buf)
	n, err := r.ResponseWriter.Write(buf)
	r.bytes += n
	return n, err
}

func (r *responseDupper) maybeWriteHeader() {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
}

func (r *responseDupper) Status() int {
	return r.code
}

func (r *responseDupper) BytesWritten() int {
	return r.bytes
}

// shortID produces a " unique" 6 bytes long string.
// Do not use as a reliable way to get unique IDs, instead use for things like logging.
func shortID() string {
	b := make([]byte, 6)
	io.ReadFull(rand.Reader, b)
	return base64.RawURLEncoding.EncodeToString(b)
}
