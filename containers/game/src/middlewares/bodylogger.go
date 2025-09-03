package middlewares

// https://developers.cyberagent.co.jp/blog/archives/50538/
// この記事を参考にしたレスポンスロガー

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"game/logger"
	"io"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RequestLogger() echo.MiddlewareFunc {
	var (
		reqBody []byte
		resBody *bytes.Buffer
	)

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURIPath: true, // リクエスト URI パスを指定
		LogMethod:  true, // リクエストメソッドを指定
		LogStatus:  true, // レスポンスステータスを指定
		BeforeNextFunc: func(c echo.Context) {
			// Request
			reqBody = []byte{}
			if c.Request().Body != nil { // Read
				reqBody, _ = io.ReadAll(c.Request().Body)
			}
			c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset

			// Response
			resBody = new(bytes.Buffer)
			mw := io.MultiWriter(c.Response().Writer, resBody)
			writer := &bodyDumpResponseWriter{Writer: mw, ResponseWriter: c.Response().Writer}
			c.Response().Writer = writer
		},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Println("レスポンスLogger", fmt.Sprintf("request: %s %s, request body: %v, response status: %d, response body: %v",
				v.Method, v.URIPath, string(reqBody), v.Status, resBody.String(),
			))
			return nil
		},
	})
}

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *bodyDumpResponseWriter) Flush() {
	err := http.NewResponseController(w.ResponseWriter).Flush()
	if err != nil && errors.Is(err, http.ErrNotSupported) {
		panic(errors.New("response writer flushing is not supported"))
	}
}

func (w *bodyDumpResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(w.ResponseWriter).Hijack()
}

func (w *bodyDumpResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}