package patient_visit

import (
	"testing"
	"time"
)

func TestToUserDoc(t *testing.T) {
	dto := User{
		Id:    "u1",
		Name:  "Jozo Mrkvak",
		Email: "jozo@mail.com",
		Role:  Role("admin"),
	}

	doc := ToUserDoc(dto)

	if doc.ID != dto.Id {
		t.Errorf("ID: expected %q, got %q", dto.Id, doc.ID)
	}
	if doc.Name != dto.Name {
		t.Errorf("Name: expected %q, got %q", dto.Name, doc.Name)
	}
	if doc.Email != dto.Email {
		t.Errorf("Email: expected %q, got %q", dto.Email, doc.Email)
	}
	if doc.Role != string(dto.Role) {
		t.Errorf("Role: expected %q, got %q", string(dto.Role), doc.Role)
	}
}

func TestToUserDoc_EmptyFields(t *testing.T) {
	doc := ToUserDoc(User{})

	if doc.ID != "" || doc.Name != "" || doc.Email != "" || doc.Role != "" {
		t.Error("expected all fields to be empty for zero-value User")
	}
}

func TestToUserDTO(t *testing.T) {
	doc := UserDocument{
		ID:    "u2",
		Name:  "Miro Petarda",
		Email: "miro@mail.com",
		Role:  "doctor",
	}

	dto := ToUserDTO(doc)

	if dto.Id != doc.ID {
		t.Errorf("Id: expected %q, got %q", doc.ID, dto.Id)
	}
	if dto.Name != doc.Name {
		t.Errorf("Name: expected %q, got %q", doc.Name, dto.Name)
	}
	if dto.Email != doc.Email {
		t.Errorf("Email: expected %q, got %q", doc.Email, dto.Email)
	}
	if dto.Role != Role(doc.Role) {
		t.Errorf("Role: expected %q, got %q", Role(doc.Role), dto.Role)
	}
}

func TestToUserDTO_EmptyFields(t *testing.T) {
	dto := ToUserDTO(UserDocument{})

	if dto.Id != "" || dto.Name != "" || dto.Email != "" || dto.Role != Role("") {
		t.Error("expected all fields to be empty for zero-value UserDocument")
	}
}

func TestUserRoundtrip_DTOToDocToDTO(t *testing.T) {
	original := User{
		Id:    "u3",
		Name:  "Rudo Trulo",
		Email: "rudenko@mail.com",
		Role:  Role("patient"),
	}

	result := ToUserDTO(ToUserDoc(original))

	if result != original {
		t.Errorf("roundtrip mismatch: expected %+v, got %+v", original, result)
	}
}

func TestUserRoundtrip_DocToDTOToDoc(t *testing.T) {
	original := UserDocument{
		ID:    "u4",
		Name:  "Dano Drevo",
		Email: "danielito@mail.com",
		Role:  "patient",
	}

	result := ToUserDoc(ToUserDTO(original))

	if result != original {
		t.Errorf("roundtrip mismatch: expected %+v, got %+v", original, result)
	}
}

func TestToVisitDoc(t *testing.T) {
	now := time.Now()
	dto := PatientVisit{
		Id:         "v1",
		PatientId:  "p1",
		DoctorId:   "d1",
		VisitDate:  now.Format("2006-01-02"),
		Anamnesis:  "headache",
		Findings:   "normal",
		Conclusion: "migraine",
		Notes:      "rest advised",
		Status:     VisitStatus("open"),
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	doc := ToVisitDoc(dto)

	if doc.ID != dto.Id {
		t.Errorf("ID: expected %q, got %q", dto.Id, doc.ID)
	}
	if doc.PatientID != dto.PatientId {
		t.Errorf("PatientID: expected %q, got %q", dto.PatientId, doc.PatientID)
	}
	if doc.DoctorID != dto.DoctorId {
		t.Errorf("DoctorID: expected %q, got %q", dto.DoctorId, doc.DoctorID)
	}
	if doc.VisitDate != dto.VisitDate {
		t.Errorf("VisitDate: expected %v, got %v", dto.VisitDate, doc.VisitDate)
	}
	if doc.Anamnesis != dto.Anamnesis {
		t.Errorf("Anamnesis: expected %q, got %q", dto.Anamnesis, doc.Anamnesis)
	}
	if doc.Findings != dto.Findings {
		t.Errorf("Findings: expected %q, got %q", dto.Findings, doc.Findings)
	}
	if doc.Conclusion != dto.Conclusion {
		t.Errorf("Conclusion: expected %q, got %q", dto.Conclusion, doc.Conclusion)
	}
	if doc.Notes != dto.Notes {
		t.Errorf("Notes: expected %q, got %q", dto.Notes, doc.Notes)
	}
	if doc.Status != string(dto.Status) {
		t.Errorf("Status: expected %q, got %q", string(dto.Status), doc.Status)
	}
	if !doc.CreatedAt.Equal(dto.CreatedAt) {
		t.Errorf("CreatedAt: expected %v, got %v", dto.CreatedAt, doc.CreatedAt)
	}
	if !doc.UpdatedAt.Equal(dto.UpdatedAt) {
		t.Errorf("UpdatedAt: expected %v, got %v", dto.UpdatedAt, doc.UpdatedAt)
	}
}

func TestToVisitDoc_EmptyFields(t *testing.T) {
	doc := ToVisitDoc(PatientVisit{})

	if doc.ID != "" || doc.PatientID != "" || doc.DoctorID != "" || doc.Status != "" {
		t.Error("expected string fields to be empty for zero-value PatientVisit")
	}
}

// --- ToVisitDTO ---

func TestToVisitDTO(t *testing.T) {
	now := time.Now()
	visitDoc := PatientVisitDocument{
		ID:         "v2",
		PatientID:  "p2",
		DoctorID:   "d2",
		VisitDate:  now.Format("2006-01-02"),
		Anamnesis:  "fever",
		Findings:   "elevated temp",
		Conclusion: "flu",
		Notes:      "fluids recommended",
		Status:     "archived",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	patientDoc := UserDocument{ID: "p2", Name: "Eva Greenova"}
	doctorDoc := UserDocument{ID: "d2", Name: "Dr. House"}

	dto := ToVisitDTO(visitDoc, patientDoc, doctorDoc)

	if dto.Id != visitDoc.ID {
		t.Errorf("Id: expected %q, got %q", visitDoc.ID, dto.Id)
	}
	if dto.PatientId != visitDoc.PatientID {
		t.Errorf("PatientId: expected %q, got %q", visitDoc.PatientID, dto.PatientId)
	}
	if dto.PatientName != patientDoc.Name {
		t.Errorf("PatientName: expected %q, got %q", patientDoc.Name, dto.PatientName)
	}
	if dto.DoctorId != visitDoc.DoctorID {
		t.Errorf("DoctorId: expected %q, got %q", visitDoc.DoctorID, dto.DoctorId)
	}
	if dto.DoctorName != doctorDoc.Name {
		t.Errorf("DoctorName: expected %q, got %q", doctorDoc.Name, dto.DoctorName)
	}
	if dto.Status != VisitStatus(visitDoc.Status) {
		t.Errorf("Status: expected %q, got %q", VisitStatus(visitDoc.Status), dto.Status)
	}
	if dto.VisitDate != visitDoc.VisitDate {
		t.Errorf("VisitDate: expected %v, got %v", visitDoc.VisitDate, dto.VisitDate)
	}
	if !dto.CreatedAt.Equal(visitDoc.CreatedAt) {
		t.Errorf("CreatedAt: expected %v, got %v", visitDoc.CreatedAt, dto.CreatedAt)
	}
	if !dto.UpdatedAt.Equal(visitDoc.UpdatedAt) {
		t.Errorf("UpdatedAt: expected %v, got %v", visitDoc.UpdatedAt, dto.UpdatedAt)
	}
}

func TestToVisitDTO_NamesResolvedFromUserDocs(t *testing.T) {
	dto := ToVisitDTO(
		PatientVisitDocument{},
		UserDocument{Name: "Patient One"},
		UserDocument{Name: "Doctor Two"},
	)

	if dto.PatientName != "Patient One" {
		t.Errorf("PatientName: expected %q, got %q", "Patient One", dto.PatientName)
	}
	if dto.DoctorName != "Doctor Two" {
		t.Errorf("DoctorName: expected %q, got %q", "Doctor Two", dto.DoctorName)
	}
}

// --- ToVisitDoc / ToVisitDTO roundtrip ---

func TestVisitRoundtrip_DocToDTOToDoc(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	original := PatientVisitDocument{
		ID:         "v3",
		PatientID:  "p3",
		DoctorID:   "d3",
		VisitDate:  now.Format("2006-01-02"),
		Anamnesis:  "cough",
		Findings:   "clear lungs",
		Conclusion: "viral",
		Notes:      "monitor",
		Status:     "open",
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	patientDoc := UserDocument{ID: "p3", Name: "Frank"}
	doctorDoc := UserDocument{ID: "d3", Name: "MuDr. Slivon"}

	result := ToVisitDoc(ToVisitDTO(original, patientDoc, doctorDoc))

	if result.ID != original.ID {
		t.Errorf("ID: expected %q, got %q", original.ID, result.ID)
	}
	if result.Status != original.Status {
		t.Errorf("Status: expected %q, got %q", original.Status, result.Status)
	}
	if result.VisitDate != original.VisitDate {
		t.Errorf("VisitDate: expected %v, got %v", original.VisitDate, result.VisitDate)
	}
}
