package patient_visit

import "strings"

func ToUserDoc(dto User) UserDocument {
	return UserDocument{
		ID:    dto.Id,
		Name:  dto.Name,
		Email: dto.Email,
		Role:  string(dto.Role),
	}
}

func ToUserDTO(doc UserDocument) User {
	return User{
		Id:    doc.ID,
		Name:  doc.Name,
		Email: doc.Email,
		Role:  Role(doc.Role),
	}
}

func ToVisitDoc(dto PatientVisit) PatientVisitDocument {
	return PatientVisitDocument{
		ID:         dto.Id,
		PatientID:  dto.PatientId,
		DoctorID:   dto.DoctorId,
		VisitDate:  dto.VisitDate,
		Anamnesis:  dto.Anamnesis,
		Findings:   dto.Findings,
		Conclusion: dto.Conclusion,
		Notes:      strings.TrimSpace(dto.Notes),
		Status:     string(dto.Status),
		CreatedAt:  dto.CreatedAt,
		UpdatedAt:  dto.CreatedAt,
	}
}

func ToVisitDTO(doc PatientVisitDocument, patient UserDocument, doctor UserDocument) PatientVisit {
	return PatientVisit{
		Id:          doc.ID,
		PatientId:   doc.PatientID,
		PatientName: patient.Name,
		DoctorId:    doc.DoctorID,
		DoctorName:  doctor.Name,
		VisitDate:   doc.VisitDate,
		Anamnesis:   doc.Anamnesis,
		Findings:    doc.Findings,
		Conclusion:  doc.Conclusion,
		Notes:       doc.Notes,
		Status:      VisitStatus(doc.Status),
		CreatedAt:   doc.CreatedAt,
		UpdatedAt:   doc.UpdatedAt,
	}
}
