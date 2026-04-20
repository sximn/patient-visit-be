package patient_visit

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implDoctorsAPI struct {
}

func NewDoctorsApi() DoctorsAPI {
	return &implDoctorsAPI{}
}

func (o implDoctorsAPI) GetDoctors(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDoctorsAPI) GetDoctor(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDoctorsAPI) GetPatientVisitsByDoctor(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
