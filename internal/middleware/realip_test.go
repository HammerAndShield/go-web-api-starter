package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXForwardForIp(t *testing.T) {
	xForwardedForIps := []string{
		"100.100.100.100",
		"100.100.100.100, 200.200.200.200",
		"100.100.100.100,200.200.200.200",
	}

	for _, v := range xForwardedForIps {
		req, _ := http.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("X-Forwarded-For", v)

		w := httptest.NewRecorder()

		mux := http.NewServeMux()

		realIp := ""
		mux.HandleFunc("GET /", func(writer http.ResponseWriter, request *http.Request) {
			realIp = request.RemoteAddr
			w.WriteHeader(200)
		})
		var handler http.Handler = RealIP(mux)
		handler.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Fatal("should have 200 code")
		}

		if realIp != "100.100.100.100" {
			t.Fatal("real ip is not correct")
		}
	}
}
