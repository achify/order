package metrics

import (
    "encoding/json"
    "expvar"
    "net/http"
    "sync/atomic"
    "runtime"
)

var requestCount uint64

// Middleware counts requests
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        atomic.AddUint64(&requestCount, 1)
        next.ServeHTTP(w, r)
    })
}

// Handler exposes simple metrics about memory and request rate
func Handler(w http.ResponseWriter, r *http.Request) {
    m := &runtime.MemStats{}
    runtime.ReadMemStats(m)
    data := map[string]interface{}{
        "requests_total": atomic.LoadUint64(&requestCount),
        "alloc":          m.Alloc,
        "sys":            m.Sys,
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// Expose expvar metrics as well
func init() {
    expvar.Publish("requests_total", expvar.Func(func() interface{} { return atomic.LoadUint64(&requestCount) }))
}
