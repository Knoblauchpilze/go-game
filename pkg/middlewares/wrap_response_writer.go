package middlewares

import "net/http"

type wrapResponseWriter interface {
	http.ResponseWriter

	Status() int
	BytesWritten() int
}

type wrapResponseWriterImpl struct {
	wrapped http.ResponseWriter

	bytesWritten int
	writeErr     error
	code         int
}

func wrap(w http.ResponseWriter) wrapResponseWriter {
	return &wrapResponseWriterImpl{
		wrapped: w,
	}
}

func (w *wrapResponseWriterImpl) Header() http.Header {
	return w.wrapped.Header()
}

func (w *wrapResponseWriterImpl) Write(in []byte) (int, error) {
	w.bytesWritten, w.writeErr = w.wrapped.Write(in)
	return w.bytesWritten, w.writeErr
}

func (w *wrapResponseWriterImpl) WriteHeader(statusCode int) {
	w.code = statusCode
	w.wrapped.WriteHeader(statusCode)
}

func (w *wrapResponseWriterImpl) Status() int {
	return w.code
}

func (w *wrapResponseWriterImpl) BytesWritten() int {
	return w.bytesWritten
}
