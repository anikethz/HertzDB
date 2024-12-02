package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ConvertStoJ[T interface{}](serializedJson string) T {
	var result T
	err := json.Unmarshal([]byte(serializedJson), &result)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
	}
	return result
}

func ConvertStoJList[T interface{}](serializedJson []string) []T {

	result := Map(serializedJson, ConvertStoJ[T])
	fmt.Printf("StoJ: %v ", result[0])
	return result
}

func ConvertStoMap(serializedJson string) map[string]interface{} {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(serializedJson), &result)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
	}
	return result
}

func ConvertStoMapList(serializedJson []string) []map[string]interface{} {

	result := Map(serializedJson, ConvertStoMap)
	fmt.Printf("StoM: %v ", result[0])
	return result
}

func RespondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if strPayload, ok := payload.(string); ok {
		w.Write([]byte(strPayload))
	} else {
		data, err := json.Marshal(payload)

		if err != nil {
			log.Printf("Failed to marshal JSON response: %v", payload)
			w.WriteHeader(500)
			return
		}
		w.Write(data)
	}
}

func ResponseWithError(w http.ResponseWriter, code int, msg string) {

	if code > 499 {
		log.Println("Failed with 5XX: ", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJson(w, code, errorResponse{
		Error: msg,
	})
}
