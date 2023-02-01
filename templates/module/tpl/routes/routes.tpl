package routes

import (
	"atom/contracts"
	"atom/providers/http"
)

type Route struct {
	svc *http.Service
}

func NewRoute(svc *http.Service) contracts.Route {
	return &Route{svc: svc}
}

func (r *Route) Register() {
	// r.svc.Engine.GET("/", gen.DataFunc(r.captcha.GetName))
}
