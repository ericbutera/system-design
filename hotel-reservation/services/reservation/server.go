package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"

	// gMiddleware "github.com/ericbutera/system-design/hotel-reservation/services/reservation/middleware"
	"github.com/ericbutera/system-design/hotel-reservation/services/reservation/graph"
	"github.com/ericbutera/system-design/hotel-reservation/services/reservation/graph/auth"
	"github.com/ericbutera/system-design/hotel-reservation/services/reservation/internal/db"
	"github.com/ericbutera/system-design/hotel-reservation/services/reservation/internal/reservations"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(auth.Middleware())

	db := lo.Must(db.New(db.NewDefaultConfig()))
	// TODO: dataloader router.Use(gMiddleware.DataloaderMiddleware(db))

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Reservations: reservations.New(db),
	}}))
	srv.AddTransport(transport.POST{})
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
