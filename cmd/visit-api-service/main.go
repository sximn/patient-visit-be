package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sximn/patient-visit-be/api"
	"github.com/sximn/patient-visit-be/internal/db_service"
	"github.com/sximn/patient-visit-be/internal/patient_visit"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	userService := db_service.NewMongoService[patient_visit.UserDocument](
		db_service.MongoServiceConfig{
			Collection: "users",
		},
	)
	defer userService.Disconnect(context.Background())

	visitService := db_service.NewMongoService[patient_visit.PatientVisitDocument](
		db_service.MongoServiceConfig{
			Collection: "patient_visits",
		},
	)
	defer visitService.Disconnect(context.Background())

	services := &patient_visit.Services{
		Users:  userService,
		Visits: visitService,
	}
	engine.Use(func(c *gin.Context) {
		c.Set("services", services)
		c.Next()
	})

	// request routings
	handleFunctions := &patient_visit.ApiHandleFunctions{
		PatientsAPI:      patient_visit.NewPatientsApi(),
		DoctorsAPI:       patient_visit.NewDoctorsApi(),
		PatientVisitsAPI: patient_visit.NewPatientVisitsApi(),
	}
	patient_visit.NewRouterWithGinEngine(engine, *handleFunctions)
	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
