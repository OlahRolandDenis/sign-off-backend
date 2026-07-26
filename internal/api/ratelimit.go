package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/httprate"
)

func RateLimitMiddleware() func(http.Handler) http.Handler {
	return httprate.Limit(
		50,             // limita de request-uri
		15*time.Minute, // fereastra de timp
		httprate.WithKeyFuncs(func(r *http.Request) (string, error) {
			agencyID, ok := r.Context().Value(AgencyIDKey).(int64)
			if ok {
				return fmt.Sprintf("agency:%d", agencyID), nil
			}
			ip := strings.Split(r.RemoteAddr, ":")[0]
			return fmt.Sprintf("ip:%s", ip), nil
		}),
	)
}
