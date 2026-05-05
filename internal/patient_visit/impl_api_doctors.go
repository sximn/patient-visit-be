package patient_visit

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sximn/patient-visit-be/internal/db_service"
	"go.mongodb.org/mongo-driver/bson"
)

type implDoctorsAPI struct {
}

func NewDoctorsApi() DoctorsAPI {
	return &implDoctorsAPI{}
}

func (o implDoctorsAPI) GetDoctors(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	doctors, err := s.Users.FindMany(c, bson.M{
		"role": "doctor",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var result []User
	for _, doc := range doctors {
		result = append(result, ToUserDTO(doc))
	}

	c.JSON(http.StatusOK, result)
}

func (o implDoctorsAPI) GetDoctor(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	id := c.Param("id")

	doc, err := s.Users.FindDocument(c, id)
	if err != nil {
		if err == db_service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
		}
		return
	}

	if doc.Role != "doctor" {
		c.JSON(http.StatusNotFound, gin.H{"error": "doctor not found"})
		return
	}

	c.JSON(http.StatusOK, ToUserDTO(doc))
}

func (o implDoctorsAPI) GetPatientVisitsByDoctor(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	doctorId := c.Param("id")

	visits, err := s.Visits.FindMany(c, bson.M{
		"doctorId": doctorId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// collect user IDs
	userIdsSet := make(map[string]struct{})
	for _, v := range visits {
		userIdsSet[v.PatientID] = struct{}{}
		userIdsSet[v.DoctorID] = struct{}{}
	}

	var userIds []string
	for id := range userIdsSet {
		userIds = append(userIds, id)
	}

	users, err := s.Users.FindMany(c, bson.M{
		"id": bson.M{"$in": userIds},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	userMap := make(map[string]UserDocument)
	for _, u := range users {
		userMap[u.ID] = u
	}

	var result []PatientVisit
	for _, v := range visits {
		result = append(result, ToVisitDTO(v, userMap[v.PatientID], userMap[v.DoctorID]))
	}

	c.JSON(http.StatusOK, result)
}
