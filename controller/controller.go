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

var message models.Status

func sendRequest(url string, body []byte, login_config models.Login) {
	// setAuth()
	request, error := http.NewRequest("POST", url, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", login_config.Headers.Content_type)
	request.Header.Set("mode", login_config.Headers.Mode)
	request.Header.Set("access_token", Token)
	BuildHeader(request, body)
	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		message = models.Status{Status: "accepted"}
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
	setAuth()
	login_config, lights_config, devices := getConfig()
	light_devices := devices.EmergencyLights
	base_url := login_config.Host + login_config.Device_path
	var commands map[string][]models.IOTStringLights
	commands = make(map[string][]models.IOTStringLights, 1)
	var white_lights_obj models.IOTStringLights
	_ = json.Unmarshal([]byte(lights_config.Modes["white_lights"]), &white_lights_obj)
	commands["commands"] = append(commands["commands"], white_lights_obj)
	body, _ := json.Marshal(commands)
	setDevices(light_devices, base_url, body, login_config)
}

func switchAllLights(status bool) {
	setAuth()
	login_config, lights_config, devices := getConfig()
	light_devices := devices.EmergencyLights
	base_url := login_config.Host + login_config.Device_path
	var commands map[string][]models.IOTSwitchDevices
	commands = make(map[string][]models.IOTSwitchDevices, 1)
	var switch_lights_obj models.IOTSwitchDevices
	_ = json.Unmarshal([]byte(lights_config.Modes["switch"]), &switch_lights_obj)
	switch_lights_obj.Value = status
	commands["commands"] = append(commands["commands"], switch_lights_obj)
	body, _ := json.Marshal(commands)
	setDevices(light_devices, base_url, body, login_config)
}

func setDevices(light_devices []string, base_url string, body []byte, login_config models.Login) {
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

func turnOnEmergencyLights(timer int64, ch chan bool) {
	lights_config := config.Read_Config("Light_Modes")
	login_config := config.Read_Config("Login")
	light_devices := config.Read_Config("Devices").(models.Devices).Lights
	var commands map[string][]models.IOTRGBFlashLights
	commands = make(map[string][]models.IOTRGBFlashLights, 1)
	var emergency_lights_obj models.IOTRGBFlashLights
	_ = json.Unmarshal([]byte(lights_config.(models.Light_Modes).Modes["emergency_lights"]), &emergency_lights_obj)
	commands["commands"] = append(commands["commands"], emergency_lights_obj)
	base_url := login_config.(models.Login).Host + login_config.(models.Login).Device_path
	emegergency_lights_body, err := json.Marshal(commands)
	if err != nil {
		panic(err)
		message = models.Status{Status: "failed"}
	}
	switchAllLights(true)
	setDevices(light_devices, base_url, emegergency_lights_body, login_config.(models.Login))

	ch <- true
	<-time.After(time.Duration(timer) * time.Second)
	setAllWhiteLights()
	switchAllLights(false)
}

func turnOnPartyLights(timer int64, frequency int, ch chan bool) {
	lights_config := config.Read_Config("Light_Modes")
	login_config := config.Read_Config("Login")
	light_devices := config.Read_Config("Devices").(models.Devices).Lights
	var commands map[string][]models.IOTRGBFlashLights
	commands = make(map[string][]models.IOTRGBFlashLights, 1)
	var party_lights_obj models.IOTRGBFlashLights
	_ = json.Unmarshal([]byte(lights_config.(models.Light_Modes).Modes["party_lights"]), &party_lights_obj)
	party_lights_obj.Value.Frequency = frequency
	commands["commands"] = append(commands["commands"], party_lights_obj)
	fmt.Printf("%s", commands)
	base_url := login_config.(models.Login).Host + login_config.(models.Login).Device_path
	party_lights_body, err := json.Marshal(commands)
	if err != nil {
		panic(err)
		message = models.Status{Status: "failed"}
	}
	switchAllLights(true)
	setDevices(light_devices, base_url, party_lights_body, login_config.(models.Login))

	ch <- true
	<-time.After(time.Duration(timer) * time.Second)
	setAllWhiteLights()
	switchAllLights(false)
}

func handleEmergencyLights(emergencylights models.EmergencyLights, ch chan bool) {
	switch emergencylights.Status {
	case "On":
		{
			if emergencylights.Timer < 1 {
				emergencylights.Timer = 15
			}
			go func() {
				turnOnEmergencyLights(emergencylights.Timer, ch)
			}()
			message = models.Status{Status: "accepted"}
			select {
			case <-ch:
				ch <- true
			}
		}

	case "Off":
		// TODO: Needs better abort system to stop On goroutine
		message = models.Status{Status: "accepted"}
		setAllWhiteLights()
		switchAllLights(false)
	}
}

func handlePartyLights(partylights models.PartyLights, ch chan bool) {
	switch partylights.Status {
	case "On":
		{
			if partylights.Timer < 1 {
				partylights.Timer = 15
			}
			go func() {
				turnOnPartyLights(partylights.Timer, partylights.Frequency, ch)
			}()
			message = models.Status{Status: "accepted"}
			select {
			case <-ch:
				ch <- true
			}
		}

	case "Off":
		// TODO: Needs better abort system to stop On goroutine
		message = models.Status{Status: "accepted"}
		setAllWhiteLights()
		switchAllLights(false)
	}
}

func handleWhiteLights(whitelights models.WhiteLights) {
	switch whitelights.Status {
	case "On":
		{
			go func() {
				setAllWhiteLights()
			}()
			message = models.Status{Status: "accepted"}
		}

	case "Off":
		message = models.Status{Status: "accepted"}

		switchAllLights(false)
	}
}

func PartyLights(w http.ResponseWriter, r *http.Request) {
	ch := make(chan bool, 1)
	w.Header().Set("Content-Type", "application/json")
	var partylights models.PartyLights
	err := json.NewDecoder(r.Body).Decode(&partylights)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handlePartyLights(partylights, ch)
	}()
	wg.Wait()
	json.NewEncoder(w).Encode(message)
}

func WhiteLights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var whitelights models.WhiteLights
	err := json.NewDecoder(r.Body).Decode(&whitelights)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleWhiteLights(whitelights)
	}()
	wg.Wait()
	json.NewEncoder(w).Encode(message)
}

func EmergencyLights(w http.ResponseWriter, r *http.Request) {
	ch := make(chan bool, 1)
	w.Header().Set("Content-Type", "application/json")
	var emergencylights models.EmergencyLights
	err := json.NewDecoder(r.Body).Decode(&emergencylights)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleEmergencyLights(emergencylights, ch)
	}()
	wg.Wait()
	json.NewEncoder(w).Encode(message)
}

func TurnOffAllDevices(){
	switchAllLights(false)
	login_config, lights_config, devices := getConfig()
	light_devices := devices.Sockets
	base_url := login_config.Host + login_config.Device_path
	var commands map[string][]models.IOTSwitchDevices
	commands = make(map[string][]models.IOTSwitchDevices, 1)
	var switch_lights_obj models.IOTSwitchDevices
	_ = json.Unmarshal([]byte(lights_config.Modes["switch"]), &switch_lights_obj)
	switch_lights_obj.Code = "switch"
	switch_lights_obj.Value = false
	commands["commands"] = append(commands["commands"], switch_lights_obj)
	body, _ := json.Marshal(commands)
	setDevices(light_devices, base_url, body, login_config)
}
