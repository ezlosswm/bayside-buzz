package lib

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetId(r *http.Request) (int32, error) {
	vars := mux.Vars(r)
	idParam := vars["id"]

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return 0, nil
	}

	return int32(id), nil
}