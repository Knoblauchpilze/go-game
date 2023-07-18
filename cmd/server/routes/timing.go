package routes

import (
	"net/http"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

type timingLogEntry struct {
	start time.Time
	host  string
	req   *http.Request
}

func (t *timingLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	if status == http.StatusOK {
		logger.ScopedInfof(t.req.Context(), "%v %v from %v, elapsed: %v", status, t.req.Method, t.host, elapsed)
	} else {
		logger.ScopedWarnf(t.req.Context(), "%v %v from %v, elapsed: %v", status, t.req.Method, t.host, elapsed)
	}
}

func (t *timingLogEntry) Panic(v interface{}, stack []byte) {
	logger.ScopedErrorf(t.req.Context(), "FAIL %v from %v, elapsed: %v (err: %+v)", t.req.Method, t.host, time.Since(t.start), stack)
}

type timingLogFormatter struct{}

func (t timingLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &timingLogEntry{
		start: time.Now(),
		host:  r.Host + r.RequestURI,
		req:   r,
	}
}
