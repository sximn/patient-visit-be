package patient_visit

import "time"

type PatientVisitDocument struct {
	ID        string `bson:"id"`
	PatientID string `bson:"patientId"`
	DoctorID  string `bson:"doctorId"`

	VisitDate string `bson:"visitDate"`

	Anamnesis  string `bson:"anamnesis"`
	Findings   string `bson:"findings"`
	Conclusion string `bson:"conclusion"`
	Notes      string `bson:"notes"`

	Status string `bson:"status"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

func (u PatientVisitDocument) GetID() string {
	return u.ID
}
