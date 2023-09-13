package infra

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"go.uber.org/fx"
	"net/http"
	"time"
)

type (
	HTTPHandler interface {
		Method() string
		Pattern() string

		http.Handler
	}

	RouteParams struct {
		fx.In
		Handlers []HTTPHandler `group:"handlers"`
	}
)

func NewHTTPRouter(params RouteParams) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})

	router.Route("/generator-pdf", func(r chi.Router) {
		for _, handler := range params.Handlers {
			r.Route(handler.Pattern(), func(subRoute chi.Router) {
				subRoute.Method(handler.Method(), "/", handler)
			})
		}
	})

	return router
}

func StartHttpServer(lifecycle fx.Lifecycle, router *chi.Mux) {
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5000 * time.Second,
		WriteTimeout: 5000 * time.Second,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			fmt.Println("Starting server...")
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func Serve() {
	ServerDependencies := fx.Provide(
		NewHTTPRouter)

	app := fx.New(
		fx.Options(
			ServerDependencies,
		),
		fx.Invoke(StartHttpServer),
	)

	app.Run()
}
