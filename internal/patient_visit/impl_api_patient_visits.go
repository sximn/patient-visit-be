package patient_visit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implPatientVisitsAPI struct {
}

func NewPatientVisitsApi() PatientVisitsAPI {
	return &implPatientVisitsAPI{}
}

func (o implPatientVisitsAPI) GetPatientVisits(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientVisitsAPI) GetPatientVisit(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientVisitsAPI) DeletePatientVisit(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientVisitsAPI) UpdatePatientVisit(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientVisitsAPI) CreatePatientVisit(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
