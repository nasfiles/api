package http

import (
	"encoding/json"
	"net/http"
)

func jsonKeys(jsonBytes []byte) ([]string, error) {
	parsed := make(map[string]interface{})

	// parse JSON into map
	err := json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		return nil, err
	}

	// create string array with all the keys
	keys := make([]string, len(parsed))
	for key := range parsed {
		keys = append(keys, key)
	}

	return keys, nil
}

func jsonPrint(w http.ResponseWriter, v interface{}) (int, error) {
	// convert all the data to bytes
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// write the bytes into the response
	w.Write(jsonBytes)

	return 0, nil
}
