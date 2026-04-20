package patient_visit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implPatientsAPI struct {
}

func NewPatientsApi() PatientsAPI {
	return &implPatientsAPI{}
}

func (o implPatientsAPI) GetPatients(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientsAPI) GetPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientsAPI) GetPatientVisitsByPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
