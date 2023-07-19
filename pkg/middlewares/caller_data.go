package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/KnoblauchPilze/go-game/pkg/logger"
)

var infoLog = logger.ScopedInfof
var warnLog = logger.ScopedWarnf

type callerData struct {
	start time.Time
	req   *http.Request
}

func (c callerData) host() string {
	if c.req == nil {
		return "N/A"
	}

	return c.req.Host + c.req.RequestURI
}

func (c callerData) write(w wrapResponseWriter) {
	str := c.serialize(w)

	if w.Status() == http.StatusOK {
		infoLog(c.req.Context(), str)
		return
	}

	warnLog(c.req.Context(), str)
}

func (c callerData) serialize(w wrapResponseWriter) string {
	return fmt.Sprintf("%v %v from %v, elapsed: %v", w.Status(), c.req.Method, c.host(), time.Since(c.start))
}
