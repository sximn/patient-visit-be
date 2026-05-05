package patient_visit

import (
	"github.com/sximn/patient-visit-be/internal/db_service"
)

type Services struct {
	Users  db_service.DbService[UserDocument]
	Visits db_service.DbService[PatientVisitDocument]
}
