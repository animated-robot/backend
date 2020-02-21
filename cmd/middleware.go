package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Middleware func(next http.Handler) http.Handler


/// MiddlewarePipeliner ///

type MiddlewarePipeFier interface {
	Pipeline() Middleware
}

type MiddlewarePipeliner struct {}

func NewMiddlewarePipeliner() MiddlewarePipeliner {
	return MiddlewarePipeliner{}
}

func (m MiddlewarePipeliner) Pipeline(ms ...Middleware) Middleware{
	return pipefy(ms)
}

func pipefy(ms []Middleware) Middleware {
	return func(h http.Handler) http.Handler {
		current := h
		last := len(ms) - 1
		for i := last; i >= 0; i-- {
			current = ms[i](current)
		}
		return current
	}
}

/////////////////////

type Middlewares struct {
	log *logrus.Logger
}

func NewMiddlewares(logger *logrus.Logger) Middlewares {
	return Middlewares{
		log: logger,
	}
}

func NewMiddlewarePipeline(logger *logrus.Logger) Middleware{
	ms := NewMiddlewares(logger)

	mps := NewMiddlewarePipeliner()
	return mps.Pipeline(/*ms.LogRequest, */ms.CorsMiddleware)
}

/// Middlewares ///
func (m Middlewares) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		r.Header.Del("Origin")

		next.ServeHTTP(w, r)
	})
}

func (m Middlewares) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqJson, err := m.parseRequestToJson(r)
		if err != nil {
			m.log.WithFields(logrus.Fields{
				"middleware": "LogRequest",
			}).Trace("Error during request parsing")
		} else {
			m.log.WithFields(logrus.Fields{
				"middleware": "LogRequest",
				"request": reqJson,
			}).Trace("Log Request Middleware")
		}
		next.ServeHTTP(w, r)
	})
}

type Request struct {
	Url     string   `json:"url"`
	Headers map[string][]string `json:"headers"`
	Body    string `json:"body"`
}

func (m Middlewares) parseRequestToJson(r *http.Request) (string, error) {
	body, err := m.parseBody(r)
	if err != nil {
		return "", err
	}

	req := Request{
		Url:     r.Method + " => " + r.Host + r.URL.String(),
		Headers: r.Header,
		Body:    body,
	}
	bs, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	str := string(bs)
	return str, nil
}

func (m Middlewares) parseBody(r *http.Request) (string, error) {
	switch r.Body {
	case http.NoBody:
		return "", nil
	default:
		r2 := r.Clone(r.Context())

		var body []byte
		if _, err := r2.Body.Read(body); err != nil {
			return "", err
		}
		return string(body), nil
	}
}