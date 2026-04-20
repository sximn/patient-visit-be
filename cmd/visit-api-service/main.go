package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sximn/patient-visit-be/api"
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
