package traefik_api_key_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	ERROR_RESPONSE_INVALID_API_KEY = "invalid api key. provide a valid api key using header %s: %s"
)

type Response struct {
	Message string `json:"message"`
}

type APIKeyAuth struct {
	next                  http.Handler
	headerName            string
	keys                  []string
	removeHeaderOnSuccess bool
}

func (a *APIKeyAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if Contains(a.keys, req.Header.Get(a.headerName)) {
		a.RemoveHeaderIfAllowed(req)
		a.next.ServeHTTP(rw, req)
		return
	}

	response := new(Response)
	response.Message = fmt.Sprintf(ERROR_RESPONSE_INVALID_API_KEY, a.headerName, "<key>")
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(rw).Encode(response); err != nil {
		fmt.Printf("Error when sending response to an invalid key: %s", err.Error())
	}
}

func (a *APIKeyAuth) RemoveHeaderIfAllowed(r *http.Request) {
	if a.removeHeaderOnSuccess {
		r.Header.Del(a.headerName)
	}
	return
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Printf("Creating plugin: %s instance: %+v, ctx: %+v\n", name, *config, ctx)
	if len(config.Keys) == 0 {
		return nil, fmt.Errorf("please provide keys for api")
	}
	if config.HeaderName == "" {
		config.HeaderName = "X-Api-Key"
	}

	return &APIKeyAuth{
		next:                  next,
		headerName:            config.HeaderName,
		keys:                  config.Keys,
		removeHeaderOnSuccess: config.RemoveHeaderOnSuccess,
	}, nil

}
