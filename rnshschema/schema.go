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
				Type:        graphql.String,
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

var RootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RNSHMutations",
	Fields: graphql.Fields{
		"createPatient": &graphql.Field{
			Type:        patientSchemaType,
			Description: "Create a new patient.",
			Args: graphql.FieldConfigArgument{
				"firstname": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Patient first name",
				},
				"surname": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Patient surname",
				},
				"gender": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Patient Gender, either MALE or FEMALE",
				},
				"dob": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Patient Date of birth",
				},
				"address": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Patients main contact address",
				},
				"mrn": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Medical Record Number used for patient identification at RNSH",
				},
				"tumorType": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Tumour Type Either Prostate, Breast or CNS",
				},
				"surgical": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "If a patient has had surgery on their tumour then true otherwise false",
				},
				"phone": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Patients phone number",
				},
				"email": &graphql.ArgumentConfig{
					Type:        graphql.NewNonNull(graphql.String),
					Description: "Patients email address",
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// marshall and cast the argument values
				firstName, _ := params.Args["firstname"].(string)
				surname, _ := params.Args["surname"].(string)
				gender, _ := params.Args["gender"].(string)
				dob, _ := params.Args["dob"].(string)
				address, _ := params.Args["address"].(string)
				mrn, _ := params.Args["mrn"].(string)
				tumorType, _ := params.Args["tumorType"].(string)
				surgical, _ := params.Args["surgical"].(string)
				phone, _ := params.Args["phone"].(string)
				email, _ := params.Args["email"].(string)

				patientService := *services.GetPatientService()

				return patientService.CreatePatient(firstName, surname, gender, dob, address, mrn, tumorType, surgical, phone, email), nil
			},
		},
	},
})

var RnshSchema = graphql.SchemaConfig{
	Query:    graphql.NewObject(rootQuery),
	Mutation: RootMutation,
}
