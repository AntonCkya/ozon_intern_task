package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/AntonCkya/ozon_habr/graph"
	"github.com/AntonCkya/ozon_habr/internal/auth"
	"github.com/AntonCkya/ozon_habr/internal/db"
	rest_handler "github.com/AntonCkya/ozon_habr/internal/handler"
	"github.com/AntonCkya/ozon_habr/internal/mem_repository"
	"github.com/AntonCkya/ozon_habr/internal/pg_repository"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	deployType := flag.String("d", "", "deploy type (d (in Docker) or n (native))")
	storageType := flag.String("s", "", "storage type (m (in memory) or p (postgres))")

	flag.Parse()
	if *storageType != "m" && *storageType != "p" {
		fmt.Println("-s flag must be 'm' or 'p'")
		flag.Usage()
		os.Exit(1)
	}
	if *deployType != "d" && *deployType != "n" {
		fmt.Println("-d flag must be 'd' or 'n'")
		flag.Usage()
		os.Exit(1)
	}

	var userRepo graph.UserRepoInterface
	var resolver *graph.Resolver
	var Host string
	if *deployType == "d" {
		Host = "db"
	} else {
		Host = "localhost"
	}

	if *storageType == "p" {
		pg, err := db.InitDB(db.DBConfig{
			Host:     Host,
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DBName:   "ozon_habr",
			SSLMode:  "disable",
		})

		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.CloseDB(pg)

		resolver = graph.NewPgResolver(pg)
		userRepo = pg_repository.NewUserRepository(pg)
	}
	if *storageType == "m" {
		resolver = graph.NewMemResolver()
		userRepo = mem_repository.NewUserRepository()
	}

	c := graph.Config{Resolvers: resolver}
	c.Directives.IsAuthenticated = auth.AuthMiddleware
	srv := handler.New(graph.NewExecutableSchema(c))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	})

	corsMiddleware := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", corsMiddleware.Handler(auth.Middleware(srv)))

	authHandler := rest_handler.NewAuthHandler(userRepo)
	http.Handle("/auth/register", http.HandlerFunc(authHandler.Register))
	http.Handle("/auth/login", http.HandlerFunc(authHandler.Login))
	http.Handle("/auth/me", http.HandlerFunc(authHandler.Me))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
