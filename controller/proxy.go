package controller

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"

	"github.com/SpeedVan/go-common/app/rest"
	"github.com/SpeedVan/go-common/client/httpclient"
	"github.com/SpeedVan/go-common/config"
	"github.com/gorilla/mux"
)

// Proxy todo
type Proxy struct {
	rest.Controller
	HTTPClient *http.Client
}

// New todo
func New(cfg config.Config) *Proxy {
	c, _ := httpclient.New(cfg)
	return &Proxy{
		HTTPClient: c,
	}
}

// GetRoute todo
func (s *Proxy) GetRoute() rest.RouteMap {
	items := []*rest.RouteItem{
		{Path: "/{path:.*}", HandleFunc: s.CORS(s.Proxy)},
	}

	return rest.NewRouteMap(items...)
}

// Proxy todo
func (s *Proxy) Proxy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	path, err := base64.StdEncoding.DecodeString(vars["path"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	req, err := http.NewRequest(r.Method, string(path), r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	req.Header = r.Header

	res, err := s.HTTPClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer res.Body.Close()

	// w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	for k, v := range res.Header {
		w.Header().Set(k, strings.Join(v, ","))
	}
	io.Copy(w, res.Body)
}

// CORS todo
func (s *Proxy) CORS(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, private-token")
		if r.Method != "OPTIONS" {
			f(w, r)
		}
	}
}
