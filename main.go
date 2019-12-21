package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	mr "github.com/vrnvu/tiny-to-go/repositories/mongodb"
	rr "github.com/vrnvu/tiny-to-go/repositories/redis"
	h "github.com/vrnvu/tiny-to-go/webservice"

	"github.com/vrnvu/tiny-to-go/tiny"
)

// https://www.google.com -> 98sj1-293
// http://localhost:8000/98sj1-293 -> https://www.google.com

// repo <- service -> serializer  -> http

func main() {

	// os.Getenv failing in ubuntu 18.04
	os.Setenv("URL_DB", "mongo")
	os.Setenv("MONGO_URL", "mongodb://localhost/tiny")
	os.Setenv("MONGO_TIMEOUT", "30")
	os.Setenv("MONGO_DB", "tiny")

	//os.Setenv("URL_DB", "redis")
	//os.Setenv("REDIS_URL", "redis://localhost:6379")

	repo := chooseRepo()
	service := tiny.NewService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)

}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() tiny.Repository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	default:
		fmt.Println("nil")
	}

	return nil
}
