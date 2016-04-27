package models

type Allergy struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type Patient struct {
	Id        string `json:"id"`
	Mrn       string `json:"mrn"`
	EhrId     string `json:"ehrId"`
	Dob       string `json:"dob"`
	Firstname string `json:"firstname"`
	Surname   string `json:"surname"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	TumorType string `json:"tumorType"`
	Surgical  string `json:"surgical"`
	Allergies []Allergy
}
