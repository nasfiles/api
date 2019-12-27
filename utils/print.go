package utils

import (
	"encoding/json"
	"net/http"
)

//JSONPrint will marshal any struct into the response body
func JSONPrint(w http.ResponseWriter, v interface{}) (int, error) {
	data, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	w.Write(data)
	return 0, nil
}
