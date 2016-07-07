package serviceshelper

import "encoding/json"

const (
	RPC_METHOD_GET_ALL_PATIENTS = "GetAllPatients"
	RPC_METHOD_GET_PATIENT      = "GetPatient"
	RPC_METHOD_GET_EHR_ID       = "GetEhrId"
	RPC_METHOD_CREATE_PATIENT   = "CreatePatient"
)

type RPCClientRequest struct {
	Method string           `json:"method"`
	Params *json.RawMessage `json:"params,omitempty"`
}

type RPCServerResponse struct {
	Result *json.RawMessage `json:"result,omitempty"`
	Error  string           `json:"error,omitempty"`
}

type RPCParamsGetPatient struct {
	PatientId int `json:"id"`
}

type RPCParamsGetEhrId struct {
	MRN string `json:"mrn"`
}

type RPCResultGetEhrId struct {
	EhrId string `json:"ehrId"`
}
