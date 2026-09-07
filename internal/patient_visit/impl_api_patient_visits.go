package patient_visit

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type implPatientVisitsAPI struct {
}

func NewPatientVisitsApi() PatientVisitsAPI {
	return &implPatientVisitsAPI{}
}

func indexUsers(users []UserDocument) map[string]UserDocument {
	result := make(map[string]UserDocument, len(users))
	for _, user := range users {
		result[user.ID] = user
	}

	return result
}

func buildVisitDTOs(visits []PatientVisitDocument, users []UserDocument) []PatientVisit {
	userMap := indexUsers(users)

	result := make([]PatientVisit, 0)
	for _, v := range visits {
		result = append(result, ToVisitDTO(
			v,
			userMap[v.PatientID],
			userMap[v.DoctorID],
		))
	}

	return result
}

func (o implPatientVisitsAPI) GetPatientVisits(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	visits, err := s.Visits.FindMany(c, bson.M{})
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

	userIds := make([]string, 0)
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

	result := buildVisitDTOs(visits, users)

	c.JSON(http.StatusOK, result)
}

func (o implPatientVisitsAPI) GetPatientVisit(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	id := c.Param("id")

	visit, err := s.Visits.FindDocument(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	users, err := s.Users.FindMany(c, bson.M{
		"id": bson.M{"$in": []string{visit.PatientID, visit.DoctorID}},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	userMap := indexUsers(users)

	dto := ToVisitDTO(
		visit,
		userMap[visit.PatientID],
		userMap[visit.DoctorID],
	)

	c.JSON(http.StatusOK, dto)
}

func (o implPatientVisitsAPI) DeletePatientVisit(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	id := c.Param("id")

	err := s.Visits.DeleteDocument(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

func (o implPatientVisitsAPI) UpdatePatientVisit(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	id := c.Param("id")

	var input PatientVisitCreateUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	existing, err := s.Visits.FindDocument(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	existing.PatientID = input.PatientId
	existing.DoctorID = input.DoctorId
	existing.VisitDate = input.VisitDate
	existing.Anamnesis = input.Anamnesis
	existing.Findings = input.Findings
	existing.Conclusion = input.Conclusion
	existing.Notes = input.Notes
	existing.Status = string(input.Status)
	existing.UpdatedAt = time.Now().UTC()

	if err := s.Visits.UpdateDocument(c, existing); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	users, _ := s.Users.FindMany(c, bson.M{
		"id": bson.M{"$in": []string{existing.PatientID, existing.DoctorID}},
	})

	userMap := indexUsers(users)

	dto := ToVisitDTO(existing, userMap[existing.PatientID], userMap[existing.DoctorID])

	c.JSON(http.StatusOK, dto)
}

func (o implPatientVisitsAPI) CreatePatientVisit(c *gin.Context) {
	s := c.MustGet("services").(*Services)

	var input PatientVisitCreateUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, err.Error())
		return
	}

	now := time.Now().UTC()

	doc := PatientVisitDocument{
		ID:         uuid.NewString(),
		PatientID:  input.PatientId,
		DoctorID:   input.DoctorId,
		VisitDate:  input.VisitDate,
		Anamnesis:  input.Anamnesis,
		Findings:   input.Findings,
		Conclusion: input.Conclusion,
		Notes:      input.Notes,
		Status:     string(input.Status),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := s.Visits.CreateDocument(c, doc); err != nil {
		c.JSON(409, err.Error())
		return
	}

	// fetch users to use in the response (we need their names)
	users, _ := s.Users.FindMany(c, bson.M{
		"id": bson.M{"$in": []string{doc.PatientID, doc.DoctorID}},
	})

	userMap := indexUsers(users)

	dto := ToVisitDTO(doc, userMap[doc.PatientID], userMap[doc.DoctorID])

	c.JSON(200, dto)
}
