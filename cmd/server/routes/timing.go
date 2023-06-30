package routes

import (
	"net/http"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
)

type timingLogEntry struct {
	start time.Time
	verb  string
	host  string
}

func (t *timingLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	if status == http.StatusOK {
		logger.Infof("%v %v from %v, elapsed: %v", status, t.verb, t.host, elapsed)
	} else {
		logger.Warnf("%v %v from %v, elapsed: %v", status, t.verb, t.host, elapsed)
	}
}

func (t *timingLogEntry) Panic(v interface{}, stack []byte) {
	logger.Errorf("FAIL %v from %v, elapsed: %v (err: %+v)", t.verb, t.host, time.Since(t.start), stack)
}

type TimingLogFormatter struct{}

func (t TimingLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &timingLogEntry{
		start: time.Now(),
		verb:  r.Method,
		host:  r.Host + r.RequestURI,
	}
}
