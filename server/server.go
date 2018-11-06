package main

import (
	log "log"
	http "net/http"
	os "os"

	handler "github.com/99designs/gqlgen/handler"
	todoql "github.com/NoahOrberg/todoql"
	"github.com/NoahOrberg/todoql/loader"
	"github.com/NoahOrberg/todoql/repository"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	repo, err := repository.New()
	if err != nil {
		panic(err) // TODO: error handling
	}
	resolver, err := todoql.NewResolver(repo)
	if err != nil {
		panic(err)
	}
	http.Handle("/query", func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			l := loader.New(repo)
			ctx := l.Attach(r.Context())
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
		}
	}(handler.GraphQL(todoql.NewExecutableSchema(todoql.Config{Resolvers: resolver}))))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
