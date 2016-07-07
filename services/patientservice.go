package services

import "github.com/FiviumAustralia/RNSH-Pilot-Server-Go/models"

type PatientService interface {
	GetAllPatients() []models.Patient
	GetPatient(id int) models.Patient
	GetEhrId(mrn string) string
	CreatePatient(firstName string, surname string, gender string, dob string, address string, mrn string,
		tumorType string, surgical string, phone string, email string) models.Patient
}

var currentPatientService PatientService

func init() {
	if currentPatientService == nil {
		rmps := RabbitMQPatientService{
			username: "go-graphql-service",
			password: "go-graphql-service",
			host:     "localhost",
			port:     "5672",
		}
		currentPatientService = rmps
	}
}

func GetPatientService() *PatientService {
	return &currentPatientService
}
