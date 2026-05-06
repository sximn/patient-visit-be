package patient_visit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sximn/patient-visit-be/internal/db_service"
	"go.mongodb.org/mongo-driver/bson"
)

type implPatientsAPI struct {
}

func NewPatientsApi() PatientsAPI {
	return &implPatientsAPI{}
}

func (o implPatientsAPI) GetPatients(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	patients, err := s.Users.FindMany(c, bson.M{
		"role": "patient",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result := make([]User, 0)
	for _, doc := range patients {
		result = append(result, ToUserDTO(doc))
	}

	c.JSON(http.StatusOK, result)
}

func (o implPatientsAPI) GetPatient(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	id := c.Param("id")

	doc, err := s.Users.FindDocument(c, id)
	if err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	if doc.Role != "patient" {
		c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		return
	}

	c.JSON(http.StatusOK, ToUserDTO(doc))
}

func (o implPatientsAPI) GetPatientVisitsByPatient(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	patientId := c.Param("id")

	visits, err := s.Visits.FindMany(c, bson.M{
		"patientId": patientId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// collect all unique userIDs (using a set) that we need to fetch for their names
	userIdsSet := make(map[string]struct{})
	for _, v := range visits {
		userIdsSet[v.PatientID] = struct{}{}
		userIdsSet[v.DoctorID] = struct{}{}
	}
	// convert to a slice
	userIds := make([]string, 0)
	for id := range userIdsSet {
		userIds = append(userIds, id)
	}
	users, err := s.Users.FindMany(c, bson.M{
		"id": bson.M{"$in": userIds},
	})

	userMap := make(map[string]UserDocument)
	for _, u := range users {
		userMap[u.ID] = u
	}

	var result []PatientVisit

	for _, v := range visits {
		patient := userMap[v.PatientID]
		doctor := userMap[v.DoctorID]

		result = append(result, ToVisitDTO(v, patient, doctor))
	}

	c.JSON(http.StatusOK, result)
}
