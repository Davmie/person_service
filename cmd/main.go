package main

import (
	"fmt"
	"github.com/Davmie/person_service/cmd/server"
	personDel "github.com/Davmie/person_service/internal/person/delivery"
	pgPerson "github.com/Davmie/person_service/internal/person/repository/postgres"
	personUseCase "github.com/Davmie/person_service/internal/person/usecase"
	"github.com/Davmie/person_service/pkg/middleware"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var prodCfgPg = postgres.Config{DSN: "host=postgres user=program password=test dbname=persons port=5432"}

func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	db, err := gorm.Open(postgres.New(prodCfgPg), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	personHandler := personDel.PersonHandler{
		PersonUseCase: personUseCase.New(pgPerson.New(logger, db)),
		Logger:        logger,
	}

	r := http.NewServeMux()

	r.Handle("GET /api/v1/persons/{personId}", http.HandlerFunc(personHandler.Get))
	r.Handle("GET /api/v1/persons", http.HandlerFunc(personHandler.GetAll))
	r.Handle("POST /api/v1/persons", http.HandlerFunc(personHandler.Create))
	r.Handle("PATCH /api/v1/persons/{personId}", http.HandlerFunc(personHandler.Update))
	r.Handle("DELETE /api/v1/persons/{personId}", http.HandlerFunc(personHandler.Delete))

	router := middleware.AccessLog(logger, r)
	router = middleware.Panic(logger, router)

	s := server.NewServer(router)
	if err := s.Start(); err != nil {
		logger.Fatal(err)
	}

	err = zapLogger.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
