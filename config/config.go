package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"iot-manager/models"
)

func Read_Config(config_type string) interface{} {
	switch config_type {
	case "Devices":
		{
			data := readFile("config/devices.json")
			var obj models.Devices
			err := json.Unmarshal(data, &obj)
			if err != nil {
				fmt.Println("error:", err)
			}
			return obj
		}
	case "Login":
		{
			data := readFile("config/login.json")
			var obj models.Login
			err := json.Unmarshal(data, &obj)
			if err != nil {
				fmt.Println("error:", err)
			}
			return obj
		}
	case "Light_Modes":
		{
			var obj models.Light_Modes
			obj.Modes = make(map[string]string)
			data := readFile("config/light_modes.json")
			var arr []string
			err := json.Unmarshal([]byte(data), &arr)
			for _, element := range arr {
				obj.Modes[element] = string(readFile("config/light_modes/" + element + ".json"))
			}
			if err != nil {
				fmt.Println("error:", err)
			}
			return obj
		}
	}
	return nil
}

func readFile(filepath string) []byte {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}
	return data
}

func Read_devices_config() models.Devices {
	devices, err := ioutil.ReadFile("config/devices.json")

	var obj models.Devices

	// unmarshall it
	err = json.Unmarshal(devices, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("%s", obj)
	return obj
}

func Read_login_config() models.Login {
	login, err := ioutil.ReadFile("config/login.json")
	if err != nil {
		fmt.Print(err)
	}

	var obj models.Login

	// unmarshall it
	err = json.Unmarshal(login, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	// fmt.Printf("%s", obj.Headers.Mode)
	return obj
}

// func Read_Config() (models.Login, models.Devices) {
// 	return Read_login_config(), Read_devices_config()
// }
