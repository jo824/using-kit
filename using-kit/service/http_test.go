package service

import (
	"encoding/json"
	"github.com/go-kit/kit/log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHTTP(t *testing.T) {
	svc := NewThingSvc()
	h := BuildHTTPHandler(svc, log.NewNopLogger())
	testSrv := httptest.NewServer(h)
	defer testSrv.Close()

	addReq := postThingRequest{
		ID:        "thing1",
		Available: true,
	}
	body, _ := json.Marshal(addReq)

	for _, tc := range []struct {
		m, url, b string
		want      int
	}{
		{"GET", testSrv.URL + "/thing/yik", "", http.StatusOK},
		{"GET", testSrv.URL + "/thing/exists", "", http.StatusNotFound},
		{"POST", testSrv.URL + "/thing", string(body), http.StatusOK},
		{"POST", testSrv.URL + "/thing", string(body), http.StatusBadRequest},
		{"POST", testSrv.URL + "/thing", "", http.StatusBadRequest},
	} {
		req, _ := http.NewRequest(tc.m, tc.url, strings.NewReader(tc.b))
		res, _ := http.DefaultClient.Do(req)
		if tc.want != res.StatusCode {
			t.Errorf("%s %s %s: want %d have %d", tc.m, tc.url, tc.b, tc.want, res.StatusCode)
		}
	}
}
