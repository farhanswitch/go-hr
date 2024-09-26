package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SimpleValidationError(errorName string, errorData error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	log.Printf("Error %s: %s\n", errorName, errorData.Error())
	errorObject := map[string][1]string{
		"errors": {errorData.Error()},
	}
	errorString, _ := json.Marshal(errorObject)
	w.Write([]byte(fmt.Sprintf(`{"errors":%s}`, errorString)))
}
