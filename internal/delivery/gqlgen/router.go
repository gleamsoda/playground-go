package gqlgen

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"playground/internal/delivery/gqlgen/gen"
	"playground/internal/delivery/gqlgen/resolver"
	"playground/internal/pkg/token"
)

func NewRouter(r *resolver.Resolver, tm token.Manager) *chi.Mux {
	h := handler.NewDefaultServer(gen.NewExecutableSchema(gen.Config{Resolvers: r}))

	router := chi.NewRouter()
	router.Use(
		middleware.Logger,
		middleware.Recoverer,
		MetadataMiddleware,
		AuthMiddlewareFunc(tm),
	)
	router.Get("/", playground.Handler("GraphQL playground", "/query"))
	router.Post("/query", h.ServeHTTP)
	return router
}
