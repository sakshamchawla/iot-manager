package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iot-manager/config"
	"iot-manager/models"
	"net/http"
	"strings"
	"sync"
)

func sendRequest(url string, commands map[string][]models.IOTRGBFlashLights, login_config models.Login) {

	body, _ := json.Marshal(commands)
	// strbody := string(body)
	// body = string(body)[1:len(string(body))-1]
	// body = bytes(body)
	t := string(body)

	fmt.Printf("Before: %s\n", t)
	t = strings.Replace(t, "\\n", "", -1)
	t = strings.Replace(t, "\\", "", -1)
	body_t := []byte(t)
	fmt.Printf("After: %s\n", body_t)
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(body_t))
	request.Header.Set("Content-Type", login_config.Headers.Content_type)
	// request.Header.Set("client_id", login_config.(models.Login).Headers.Client_id)
	// request.Header.Set("t", login_config.(models.Login).Headers.T)
	request.Header.Set("mode", login_config.Headers.Mode)

	request.Header.Set("access_token", Token)
	BuildHeader(request, body_t)
	// request.Header.Set("sign", login_config.(models.Login).Headers.Sign)
	// request.Header.Set("access_token", login_config.(models.Login).Headers.Access_token)
	fmt.Printf("%s", request.Header)
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	fmt.Println(string(b))
}

func turnOnEmergencyLights(emergencylights models.EmergencyLights) {
	emergencylights_config := config.Read_Config("Light_Modes")
	var commands map[string][]models.IOTRGBFlashLights
	commands = make(map[string][]models.IOTRGBFlashLights, 1)
	var emergency_lights_obj models.IOTRGBFlashLights
	_ = json.Unmarshal([]byte(emergencylights_config.(models.Light_Modes).Modes["emergency_lights"]), &emergency_lights_obj)
	fmt.Printf("Obj type: %T\n", emergency_lights_obj)
	commands["commands"] = append(commands["commands"], emergency_lights_obj)
	// fmt.Printf("commands type %T\n", commands)
	light_devices := config.Read_Config("Devices").(models.Devices).Lights
	login_config := config.Read_Config("Login")
	base_url := login_config.(models.Login).Host + login_config.(models.Login).Device_path
	Host = login_config.(models.Login).Host
	ClientID = login_config.(models.Login).Headers.Client_id
	Secret = login_config.(models.Login).Headers.Secret
	_ = GetToken()
	var wg sync.WaitGroup
	for _, light_device := range light_devices {
		url := base_url + light_device + "/commands"
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendRequest(url, commands, login_config.(models.Login))
		}()

	}

	// fmt.Printf("%s", emergencylights_config.(models.Light_Modes).Modes["emergency_lights"])

}

func EmergencyLights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var emergencylights models.EmergencyLights
	json.NewDecoder(r.Body).Decode(&emergencylights)
	// fmt.Printf("%s", emergencylights)
	message := models.Status{Status: "accepted"}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		turnOnEmergencyLights(emergencylights)
	}()

	// wg.Wait()

	json.NewEncoder(w).Encode(message)
}
