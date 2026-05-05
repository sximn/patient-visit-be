package patient_visit

type UserDocument struct {
	ID    string `bson:"id"`
	Name  string `bson:"name"`
	Email string `bson:"email"`
	Role  string `bson:"role"`
}

func (u UserDocument) GetID() string {
	return u.ID
}
