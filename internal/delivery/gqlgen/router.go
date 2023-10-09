package gqlgen

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samber/do"

	"playground/internal/delivery/gqlgen/gen"
	"playground/internal/delivery/gqlgen/resolver"
	"playground/internal/pkg/token"
)

func NewRouter(i *do.Injector) (*chi.Mux, error) {
	r := do.MustInvoke[*resolver.Resolver](i)
	tm := do.MustInvoke[token.Manager](i)
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
	return router, nil
}
