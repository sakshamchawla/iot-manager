package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"iot-manager/config"
	"iot-manager/models"
	"net/http"
	"sync"
	"time"
)

func sendRequest(url string, body []byte, login_config models.Login) {

	request, error := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", login_config.Headers.Content_type)
	// request.Header.Set("client_id", login_config.(models.Login).Headers.Client_id)
	// request.Header.Set("t", login_config.(models.Login).Headers.T)
	request.Header.Set("mode", login_config.Headers.Mode)

	request.Header.Set("access_token", Token)
	BuildHeader(request, body)
	// request.Header.Set("sign", login_config.(models.Login).Headers.Sign)
	// request.Header.Set("access_token", login_config.(models.Login).Headers.Access_token)
	// fmt.Printf("%s", request.Header)
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	b, _ := io.ReadAll(response.Body)

	fmt.Println(string(b))
}

func setAuth() {
	login_config := config.Read_Config("Login")
	Host = login_config.(models.Login).Host
	ClientID = login_config.(models.Login).Headers.Client_id
	Secret = login_config.(models.Login).Headers.Secret
	_ = GetToken()
}

func getConfig() (models.Login, models.Light_Modes, models.Devices) {
	return config.Read_Config("Login").(models.Login), config.Read_Config("Light_Modes").(models.Light_Modes), config.Read_Config("Devices").(models.Devices)
}

func setAllWhiteLights() {
	login_config, lights_config, devices := getConfig()
	light_devices := devices.Lights
	base_url := login_config.Host + login_config.Device_path
	var commands map[string][]models.IOTStringLights
	commands = make(map[string][]models.IOTStringLights, 1)
	var white_lights_obj models.IOTStringLights
	_ = json.Unmarshal([]byte(lights_config.Modes["white_lights"]), &white_lights_obj)
	commands["commands"] = append(commands["commands"], white_lights_obj)
	body, _ := json.Marshal(commands)
	setLight(light_devices, base_url, body, login_config)
}

func switchAllLights(status bool) {
	login_config, lights_config, devices := getConfig()
	light_devices := devices.Lights
	base_url := login_config.Host + login_config.Device_path
	var commands map[string][]models.IOTSwitchLights
	commands = make(map[string][]models.IOTSwitchLights, 1)
	var switch_lights_obj models.IOTSwitchLights
	_ = json.Unmarshal([]byte(lights_config.Modes["switch"]), &switch_lights_obj)
	switch_lights_obj.Value = status
	commands["commands"] = append(commands["commands"], switch_lights_obj)
	body, _ := json.Marshal(commands)
	setLight(light_devices, base_url, body, login_config)
}

func setLight(light_devices []string, base_url string, body []byte, login_config models.Login) {
	var wg sync.WaitGroup
	for _, light_device := range light_devices {
		url := base_url + light_device + "/commands"
		wg.Add(1)
		go func() {
			defer wg.Done()
			sendRequest(url, body, login_config)
		}()
	}
	wg.Wait()
}

func turnOnEmergencyLights(emergencylights models.EmergencyLights) {
	lights_config := config.Read_Config("Light_Modes")
	login_config := config.Read_Config("Login")
	light_devices := config.Read_Config("Devices").(models.Devices).Lights
	var commands map[string][]models.IOTRGBFlashLights
	commands = make(map[string][]models.IOTRGBFlashLights, 1)
	var emergency_lights_obj models.IOTRGBFlashLights
	_ = json.Unmarshal([]byte(lights_config.(models.Light_Modes).Modes["emergency_lights"]), &emergency_lights_obj)
	fmt.Printf("Obj type: %T\n", emergency_lights_obj)
	commands["commands"] = append(commands["commands"], emergency_lights_obj)
	// fmt.Printf("commands type %T\n", commands)

	base_url := login_config.(models.Login).Host + login_config.(models.Login).Device_path
	body, _ := json.Marshal(commands)
	switchAllLights(true)
	setLight(light_devices, base_url, body, login_config.(models.Login))
	time.Sleep(10 * time.Second)
	setAllWhiteLights()
	switchAllLights(false)
}

func EmergencyLights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var emergencylights models.EmergencyLights
	json.NewDecoder(r.Body).Decode(&emergencylights)
	// fmt.Printf("%s", emergencylights)
	message := models.Status{Status: "accepted"}
	var wg sync.WaitGroup
	wg.Add(1)
	setAuth()
	go func() {
		defer wg.Done()
		turnOnEmergencyLights(emergencylights)
	}()

	json.NewEncoder(w).Encode(message)
}
