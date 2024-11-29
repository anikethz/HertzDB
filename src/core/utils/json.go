package utils

import (
	"encoding/json"
	"fmt"
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
