package api

import (
	"Desktop/signoff/internal/auth"
	"Desktop/signoff/internal/config"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

type API struct {
	Router *chi.Mux
	Config *config.Config
}

func New(cfg *config.Config) *API {
	var api API
	api.Router = chi.NewRouter()
	api.Config = cfg
	auth.JWTSecret = cfg.JWTSecret

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
	// 1. Public routes
	api.Router.Get("/health", HealthHandler)

	// 2. Approve routes (public)
	api.Router.Get("/approve/{token}", ShowApprovePage)
	api.Router.Post("/approve/{token}/decision", ProcessDecision)

	// 3. Rate limit pentru login/register (IP-based, 10 requests per minute)
	api.Router.With(httprate.LimitByIP(10, time.Minute)).Post("/api/register", Register)
	api.Router.With(httprate.LimitByIP(10, time.Minute)).Post("/api/login", Login)

	// 4. Protected routes (AuthMiddleware + RateLimitMiddleware)
	api.Router.Group(func(r chi.Router) {
		r.Use(AuthMiddleware)
		r.Use(RateLimitMiddleware()) // rate limit per agency (50 req / 15 min)

		r.Post("/api/logout", Logout)
		r.Get("/api/account", GetAccount)
		r.Delete("/api/account", DeleteAccount)
		r.Post("/api/keys", CreateAPIKeyHandler)
		r.Get("/api/keys", ListAPIKeysHandler)
		r.Delete("/api/keys/{id}", DeleteAPIKeyHandler)
		r.Post("/webhook", WebhookHandler)
	})
}
