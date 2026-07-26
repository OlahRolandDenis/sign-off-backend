package api

import (
	"Desktop/signoff/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	Router *chi.Mux
	Config *config.Config
}

func New(cfg *config.Config) *API {
	var api API
	api.Router = chi.NewRouter()
	api.Config = cfg

	api.setUpMiddleware()
	api.setUpRoutes()

	return &api
}

func (api *API) setUpMiddleware() {
	api.Router.Use(middleware.Logger)
	api.Router.Use(middleware.Recoverer)
	api.Router.Use(corsMiddleware)
}

func (api *API) setUpRoutes() {
	api.Router.Get("/health", HealthHandler)
	api.Router.Post("/api/register", Register)
	api.Router.Post("/api/login", Login)

	api.Router.Get("/approve/{token}", ShowApprovePage)
	api.Router.Post("/approve/{token}/decision", ProcessDecision)

	api.Router.With(AuthMiddleware).Post("/api/keys", CreateAPIKeyHandler)
	api.Router.With(AuthMiddleware).Get("/api/keys", ListAPIKeysHandler)
	api.Router.With(AuthMiddleware).Delete("/api/keys/{id}", DeleteAPIKeyHandler)

	api.Router.With(AuthMiddleware).Post("/webhook", WebhookHandler)
}
