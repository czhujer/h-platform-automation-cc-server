package calculoid

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var receivedData []byte

type CalculoidData struct {
	CalculatorId        string               `json:"calculatorId"`
	Email               string               `json:"email"`
	Id                  string               `json:"id"`
	CalculoidDataParams CalculoidDataParams  `json:"params"`
	UserSignature       string               `json:"userSignature"`
	FromEmail           string               `json:"fromEmail"`
	Fields              CalculoidDataField   `json:"fields"`
	Payment             CalculoidDataPayment `json:"Payment"`
}

type Handler struct {
}

type CalculoidDataParams struct {
}

type CalculoidDataField map[string]CalculoidDataFields

type CalculoidDataFields struct {
	CalculatorFieldId string `json:"calculatorFieldId"`
	Name              string `json:"name"`
}

type CalculoidDataPayment struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

func (c *Handler) queryParams(w http.ResponseWriter, r *http.Request) {
	var err error
	receivedData, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	//log.Printf("received data: %s \n", receivedData)
}

func (c *Handler) CalculoidWebhook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		code := http.StatusBadRequest
		if r.Method == "GET" {
			code = http.StatusOK
			result := fmt.Sprintln("CalculoidWebhook")
			_, err = w.Write([]byte(result))
			if err != nil {
				log.Fatal(err)
			}
		} else if r.Method == "POST" {
			code = http.StatusOK
			result := fmt.Sprintln("CalculoidWebhook")
			_, err = w.Write([]byte(result))
			if err != nil {
				log.Fatal(err)
			}
			c.queryParams(w, r)
			c.calculoidWebhookParser()
		} else {
			w.WriteHeader(code)
		}
	}
}

func (c *Handler) calculoidWebhookParser() {
	var parsedData CalculoidData

	//log.Printf("received data for parsing: %s \n", receivedData)

	err := json.Unmarshal([]byte(receivedData), &parsedData)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("received data after parsing:")

	fmt.Printf("\ncalculatorId: %v\n", parsedData.CalculatorId)
	fmt.Printf("email: %s\n", parsedData.Email)

	fmt.Printf("fromEmail: %s\n", parsedData.FromEmail)

	for key := range parsedData.Fields {
		fmt.Printf("Fields: Key: %s Values: ", key)
		fmt.Printf("calculatorFieldId: %s ", parsedData.Fields[key].CalculatorFieldId)
		fmt.Printf("calculatorFieldId: %s ", parsedData.Fields[key].Name)

		fmt.Println("")
	}

	fmt.Printf("Payment ID: %v and Status: %v \n", parsedData.Payment.Id, parsedData.Payment.Status)

	fmt.Println("")
}
