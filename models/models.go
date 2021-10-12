package models

type Login_headers struct {
	Sign_method  string `json:"sign_method"`
	Client_id    string `json:"client_id"`
	T            string `json:"t"`
	Mode         string `json:"mode"`
	Content_type string `json:"Content-Type"`
	Sign         string `json:"sign"`
	Access_token string `json:"access_token"`
	Secret       string `json:"secret"`
}

type Login struct {
	Host        string        `json:"host"`
	Device_path string        `json:"device_path"`
	Headers     Login_headers `json:"headers"`
}

type Light_Modes struct {
	Modes map[string]string
}

type Devices struct {
	Lights          []string `json:"lights"`
	EmergencyLights []string `json:"emergencylights"`
}

type Message struct {
	Message string `json:"message"`
}

type Status struct {
	Status string `json:"status"`
}

type EmergencyLights struct {
	Status string `json: "status"`
	Code   string `json: "code"`
	Timer  int64  `json: "timer"`
}

type HSVValue struct {
	H int `json:"h"`
	S int `json:"s"`
	V int `json:"v"`
}

type IOTRGBFlashValue struct {
	Bright      int        `'json: "bright"'`
	Frequency   int        `json: "frequency"`
	HSV         []HSVValue `json: "hsv"`
	Temperature int        `json: "temperature"`
}

type IOTRGBFlashLights struct {
	Code  string           `json: "code"`
	Value IOTRGBFlashValue `json: "value"`
}

type IOTSwitchLights struct {
	Code  string `json: "code"`
	Value bool   `json: value`
}

type IOTStringLights struct {
	Code  string `json: "code"`
	Value string `json: "value"`
}
