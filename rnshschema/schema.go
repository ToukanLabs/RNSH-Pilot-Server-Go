package rnshschema

import (
	"github.com/FiviumAustralia/RNSH-Pilot-Server-Go/models"
	"github.com/FiviumAustralia/RNSH-Pilot-Server-Go/services"
	"github.com/graphql-go/graphql"
)

// Schema
var allergiesSchemaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Allergies",
		Description: "Represent the allergies of a patient",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type:        graphql.Int,
				Description: "The name of the Allergy the patient has",
			},
			"date": &graphql.Field{
				Type:        graphql.String,
				Description: "Then date the patient identified the allergy",
			},
		},
	},
)

var patientSchemaType = graphql.NewObject(
	graphql.ObjectConfig{
		Name:        "Patient",
		Description: "Represent the type of a Patient in an openEHR party",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.Int,
				Description: "openEHR Party Id of the Patient",
			},
			"mrn": &graphql.Field{
				Type:        graphql.String,
				Description: "Medical Record Number used for patient identification at RNSH",
			},
			"ehrId": &graphql.Field{
				Type:        graphql.String,
				Description: "openEHR electronic health record identifier",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					patientService := *services.GetPatientService()
					var mrn string
					if patient, ok := p.Source.(models.Patient); ok {
						mrn = patient.Mrn
					}
					return patientService.GetEhrId(mrn), nil
				},
			},
			"dob": &graphql.Field{
				Type:        graphql.String,
				Description: "Patient Date of birth",
			},
			"firstname": &graphql.Field{
				Type:        graphql.String,
				Description: "Patient first name",
			},
			"surname": &graphql.Field{
				Type:        graphql.String,
				Description: "Patient surname",
			},
			"address": &graphql.Field{
				Type:        graphql.String,
				Description: "Patients main contact address",
			},
			"phone": &graphql.Field{
				Type:        graphql.String,
				Description: "Patients phone number",
			},
			"email": &graphql.Field{
				Type:        graphql.String,
				Description: "Patients email address",
			},
			"gender": &graphql.Field{
				Type:        graphql.String,
				Description: "Patient Gender, either MALE or FEMALE",
			},
			"tumorType": &graphql.Field{
				Type:        graphql.String,
				Description: "Tumour Type Either Prostate, Breast or CNS",
			},
			"surgical": &graphql.Field{
				Type:        graphql.String,
				Description: "If a patient has had surgery on their tumour then true otherwise false",
			},
			"allergies": &graphql.Field{
				Type:        graphql.NewList(allergiesSchemaType),
				Description: "A list of allergies that the patient may have",
			},
		},
	},
)

var fields = graphql.Fields{
	"patients": &graphql.Field{
		Type: graphql.NewList(patientSchemaType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			patientService := *services.GetPatientService()
			return patientService.GetAllPatients(), nil
		},
	},
	"patient": &graphql.Field{
		Type: patientSchemaType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			patientService := *services.GetPatientService()
			id, _ := p.Args["id"].(int)
			return patientService.GetPatient(id), nil
		},
	},
}
var rootQuery = graphql.ObjectConfig{Name: "RootQuery", Fields: fields}

var RnshSchema = graphql.SchemaConfig{
	Query: graphql.NewObject(rootQuery),
}
