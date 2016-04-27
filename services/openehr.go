package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/FiviumAustralia/RNSH-Pilot-Server-Go/models"
)

type OpenEHRPatientService struct {
	baseUrl          string
	subjectNamespace string
	username         string
	password         string
}

type partyInfoType struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AddressType struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

type Party struct {
	Id                  string          `json:"id"`
	FirstNames          string          `json:"firstNames"`
	LastNames           string          `json:"lastNames"`
	Gender              string          `json:"gender"`
	DateOfBirth         string          `json:"dateOfBirth"`
	Address             AddressType     `json:"address"`
	PartyAdditionalInfo []partyInfoType `json:"partyAdditionalInfo"`
}

type PartiesType struct {
	Parties []Party `json:"parties"`
}

type PartyType struct {
	Party Party `json:"party"`
}

type ehrType struct {
	EhrId string `json:"ehrId"`
}

func (oeps OpenEHRPatientService) getAuthorizationHeader() string {
	s := fmt.Sprintf("%s:%s", oeps.username, oeps.password)
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("Basic %s", encoded)
}

func (oeps OpenEHRPatientService) GetOpenEhr(url string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", oeps.getAuthorizationHeader())
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body
}

func partyToPatient(party *Party, patient *models.Patient) {
	patient.Id = party.Id
	patient.Dob = party.DateOfBirth
	patient.Firstname = party.FirstNames
	patient.Surname = party.LastNames
	patient.Address = party.Address.Address
	patient.Gender = party.Gender

	for _, a := range party.PartyAdditionalInfo {
		switch a.Key {
		case "rnsh.mrn":
			patient.Mrn = a.Value
		case "tumorType":
			patient.TumorType = a.Value
		case "email":
			patient.Email = a.Value
		case "phone":
			patient.Phone = a.Value
		case "surgical":
			patient.Surgical = a.Value
		}
	}
}

// patient interface (service)
func (oeps OpenEHRPatientService) GetAllPatients() []models.Patient {
	url := fmt.Sprintf("%sdemographics/party/query?lastNames=*&rnsh.mrn=*", oeps.baseUrl)
	body := oeps.GetOpenEhr(url)
	var parties PartiesType
	_ = json.Unmarshal(body, &parties)
	var patients []models.Patient

	for _, p := range parties.Parties {
		var patient models.Patient
		partyToPatient(&p, &patient)

		patients = append(patients, patient)
	}

	return patients
}

func (oeps OpenEHRPatientService) GetPatient(id int) models.Patient {
	url := fmt.Sprintf("%sdemographics/party/%v", oeps.baseUrl, id)
	body := oeps.GetOpenEhr(url)
	var party PartyType
	_ = json.Unmarshal(body, &party)
	var patient models.Patient
	partyToPatient(&party.Party, &patient)
	return patient
}

func (oeps OpenEHRPatientService) GetEhrId(mrn string) string {
	url := fmt.Sprintf("%sehr/?subjectId=%s&subjectNamespace=%s", oeps.baseUrl, mrn, oeps.subjectNamespace)
	body := oeps.GetOpenEhr(url)
	var ehr ehrType
	err := json.Unmarshal(body, &ehr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", string(body))
	fmt.Println(ehr)
	fmt.Println(ehr.EhrId)
	return "%TESTING%"
}
