package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"iot-manager/controller"
	"iot-manager/models"
	"net/http"
)

func IOTRoutes() *mux.Router {
	var router = mux.NewRouter()
	router = mux.NewRouter().StrictSlash(true)

	//Home Toute
	router.HandleFunc("/iotapi/", func(rw http.ResponseWriter, r *http.Request) {
		// handle home route here
		message := models.Message{
			Message: "Root",
		}
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(message)
	})
	router.HandleFunc("/iotapi/emergencylights", controller.EmergencyLights).Methods(http.MethodPost)
	router.HandleFunc("/iotapi/partylights", controller.PartyLights).Methods(http.MethodPost)
	router.HandleFunc("/iotapi/whitelights", controller.WhiteLights).Methods(http.MethodPost)
	return router
}
