package services

import (
	"fmt"
	"log"
	"math/rand"

	"encoding/json"

	"github.com/FiviumAustralia/RNSH-Pilot-Server-Go/models"
	"github.com/FiviumAustralia/RNSH-Pilot-Server-Go/serviceshelper"
	"github.com/streadway/amqp"
)

type RabbitMQPatientService struct {
	username string
	password string
	host     string
	port     string
}

func (rmps RabbitMQPatientService) GetAllPatients() []models.Patient {
	res := rmps.MakeRPC(serviceshelper.RPC_METHOD_GET_ALL_PATIENTS, nil)

	patients := []models.Patient{}
	err := json.Unmarshal(res, &patients)
	if err != nil {
		log.Println("Unable to unmarshall GetAllPatients response:")
		log.Println(string(res))
	}

	return patients
}

func (rmps RabbitMQPatientService) GetPatient(id int) models.Patient {
	res := rmps.MakeRPC(serviceshelper.RPC_METHOD_GET_PATIENT, &serviceshelper.RPCParamsGetPatient{
		PatientId: id,
	})

	patient := models.Patient{}
	err := json.Unmarshal(res, &patient)
	if err != nil {
		log.Println("Unable to unmarshall GetPatient response:")
		log.Println(string(res))
	}

	return patient
}

func (rmps RabbitMQPatientService) GetEhrId(mrn string) string {
	res := rmps.MakeRPC(serviceshelper.RPC_METHOD_GET_EHR_ID, &serviceshelper.RPCParamsGetEhrId{
		MRN: mrn,
	})

	rrgei := serviceshelper.RPCResultGetEhrId{}
	err := json.Unmarshal(res, &rrgei)
	if err != nil {
		log.Println("Unable to unmarshall GetEhrId response:")
		log.Println(string(res))
	}

	return rrgei.EhrId
}

func (rmps RabbitMQPatientService) CreatePatient(firstName string, surname string, gender string, dob string, address string, mrn string, tumorType string, surgical string, phone string, email string) models.Patient {
	patient := models.Patient{
		Mrn:       mrn,
		Dob:       dob,
		Firstname: firstName,
		Surname:   surname,
		Address:   address,
		Phone:     phone,
		Email:     email,
		Gender:    gender,
		TumorType: tumorType,
		Surgical:  surgical,
	}

	res := rmps.MakeRPC(serviceshelper.RPC_METHOD_CREATE_PATIENT, &patient)

	newPatient := models.Patient{}
	err := json.Unmarshal(res, &newPatient)
	if err != nil {
		log.Println("Unable to unmarshall CreatePatient response:")
		log.Println(string(res))
	}

	return newPatient
}

func (rmps RabbitMQPatientService) MakeRPC(methodName string, params interface{}) []byte {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%s/", rmps.username, rmps.password, rmps.host, rmps.port)

	log.Printf("[ START ] Requesting %s", methodName)

	conn, err := amqp.Dial(connString)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var requestBody serviceshelper.RPCClientRequest
	if params != nil {
		paramsJson, err := json.Marshal(params)
		if err != nil {
			log.Println("Unable to marshall params.")
		}

		rawParams := json.RawMessage(paramsJson)

		requestBody = serviceshelper.RPCClientRequest{
			Method: methodName,
			Params: &rawParams,
		}
	} else {
		requestBody = serviceshelper.RPCClientRequest{
			Method: methodName,
		}
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("Unable to marshall JSON.")
	}

	corrId := randomString(32)

	err = ch.Publish(
		"",
		"patient_queue", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          requestBodyBytes,
		})
	if err != nil {
		log.Println("Unable to publish.", methodName)
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			log.Printf("[  END  ] Requesting %s", methodName)
			return d.Body
		}
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
