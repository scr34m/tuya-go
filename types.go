package tuya

import (
	"encoding/json"
)

type ResultType struct {
	T       int    `json:"t"`
	Success bool   `json:"success"`
	Status  string `json:"status"`

	Error        string `json:"error"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMsg"`

	A string `json:"a"`
	V string `json:"v"`

	result     map[string]interface{}
	resultList []interface{}
}

func NewResultType(b []byte) ResultType {
	var res map[string]interface{}
	json.Unmarshal(b, &res)

	t := ResultType{}
	if val, ok := res["a"]; ok {
		t.A = val.(string)
	}

	if val, ok := res["t"]; ok {
		t.T = int(val.(float64))
	}

	if val, ok := res["success"]; ok {
		t.Success = val.(bool)
	}

	if val, ok := res["v"]; ok {
		t.V = val.(string)
	}

	if val, ok := res["status"]; ok {
		t.Status = val.(string)
	}

	if val, ok := res["error"]; ok {
		t.Error = val.(string)
	}

	if val, ok := res["errorCode"]; ok {
		t.ErrorCode = val.(string)
	}

	if val, ok := res["errorMsg"]; ok {
		t.ErrorMessage = val.(string)
	}

	if val, ok := res["result"]; ok {
		switch val.(type) {
		case map[string]interface{}:
			t.result = val.(map[string]interface{})
		case []interface{}:
			t.resultList = val.([]interface{})
		}
	}

	return t
}

type RsaType struct {
	PublicKey string `json:"publicKey"`
	Token     string `json:"token"`
	Exponent  string `json:"exponent"`
}

func (t *ResultType) GetRsaType() RsaType {
	var rsa RsaType
	b, _ := json.Marshal(t.result)
	json.Unmarshal(b, &rsa)
	return rsa
}

type LoginType struct {
	Timezone string `json:"timezone"`
	TempUnit int    `json:"tempUnit"`
	// "extras": {"developer": 0},
	SID                string `json:"sid"`
	UID                string `json:"uid"`
	Nickname           string `json:"nickname"`
	PhoneCode          string `json:"phoneCode"`
	Attribute          int    `json:"attribute"`
	Email              string `json:"email"`
	ImproveCompanyInfo bool   `json:"improveCompanyInfo"`
	SnsNickname        string `json:"snsNickname"`
	Receiver           string `json:"receiver"`
	DataVersion        int    `json:"dataVersion"`
	AccountType        int    `json:"accountType"`
	Sex                int    `json:"sex"`
	Mobile             string `json:"mobile"`
	HeadPic            string `json:"headPic"`
	Ecode              string `json:"ecode"`
	RegFrom            int    `json:"regFrom"`
	/*
		"domain": {
			"aispeechHttpsUrl": "https://aispeech.tuyaeu.com",
			"aispeechQuicUrl": "https://i1.tuyaeu.com",
			"deviceHttpUrl": "http://a.tuyaeu.com",
			"deviceHttpsPskUrl": "https://a3.tuyaeu.com",
			"deviceHttpsUrl": "https://a2.tuyaeu.com",
			"deviceMediaMqttUrl": "s.tuyaeu.com",
			"deviceMediaMqttsUrl": "ms.tuyaeu.com",
			"deviceMqttsPskUrl": "m2.tuyaeu.com",
			"deviceMqttsUrl": "m2.tuyaeu.com",
			"gwApiUrl": "http://a.gw.tuyaeu.com/gw.json",
			"gwMqttUrl": "mq.gw.tuyaeu.com",
			"httpPort": 80,
			"httpsPort": 443,
			"httpsPskPort": 443,
			"mobileApiUrl": "https://a1.tuyaeu.com",
			"mobileMediaMqttUrl": "s.tuyaeu.com",
			"mobileMqttUrl": "mq.mb.tuyaeu.com",
			"mobileMqttsUrl": "m1.tuyaeu.com",
			"mobileQuicUrl": "https://u1.tuyaeu.com",
			"mqttPort": 1883,
			"mqttQuicUrl": "q1.tuyaeu.com",
			"mqttsPort": 8883,
			"mqttsPskPort": 8886,
			"pxApiUrl": "http://px.tuyaeu.com",
			"regionCode": "EU"
		}
	*/
	TimezoneID      string `json:"timezoneId"`
	UserType        int    `json:"userType"`
	PartnerIdentity string `json:"partnerIdentity"`
	Username        string `json:"username"`
}

func (t *ResultType) GetLoginType() LoginType {
	var login LoginType
	b, _ := json.Marshal(t.result)
	json.Unmarshal(b, &login)
	return login
}

type LocationType struct {
	GeoName      string        `json:"geoName"`
	Rooms        []interface{} `json:"rooms"`
	GmtModified  int           `json:"gmtModified"`
	Role         int           `json:"role"`
	Gid          int           `json:"gid"`
	GroupID      int           `json:"groupId"`
	DisplayOrder int           `json:"displayOrder"`
	Admin        bool          `json:"admin"`
	Lon          string        `json:"lon"`
	DealStatus   int           `json:"dealStatus"`
	GmtCreate    int           `json:"gmtCreate"`
	OwnerID      string        `json:"ownerId"`
	UID          string        `json:"uid"`
	GroupUserID  int           `json:"groupUserId"`
	Background   string        `json:"background"`
	Name         string        `json:"name"`
	ID           int           `json:"id"`
	Lat          string        `json:"lat"`
	Status       bool          `json:"status"`
}

func (t *ResultType) GetLocationList() []LocationType {
	var _l []LocationType
	for _, r := range t.resultList {
		var _t LocationType
		b, _ := json.Marshal(r)
		json.Unmarshal(b, &_t)
		_l = append(_l, _t)
	}
	return _l
}

type DeviceModuleWifiType struct {
	UpgradeStatus int    `json:"upgradeStatus"`
	Cdv           string `json:"cdv"`
	Bv            string `json:"bv"`
	Pv            string `json:"pv"`
	VerSw         string `json:"verSw"`
	IsOnline      bool   `json:"isOnline"`
	ID            int    `json:"id"`
	Cadv          string `json:"cadv"`
}

type DeviceModuleMcuType struct {
	UpgradeStatus int    `json:"upgradeStatus"`
	Cdv           string `json:"cdv"`
	VerSw         string `json:"verSw"`
	IsOnline      bool   `json:"isOnline"`
	ID            int    `json:"id"`
	Cadv          string `json:"cadv"`
}

type DeviceModuleMapType struct {
	Wifi DeviceModuleWifiType `json:"wifi"`
	Mcu  DeviceModuleMcuType  `json:"mcu"`
}

type DeviceType struct {
	Virtual      bool                   `json:"virtual"`
	DpName       map[string]interface{} `json:"dpName"`
	Lon          string                 `json:"lon"`
	UUID         string                 `json:"uuid"`
	Mac          string                 `json:"mac"`
	IconURL      string                 `json:"iconUrl"`
	RuntimeEnv   string                 `json:"runtimeEnv"`
	Lat          string                 `json:"lat"`
	DevID        string                 `json:"devId"`
	DpMaxTime    int                    `json:"dpMaxTime"`
	ProductID    string                 `json:"productId"`
	Dps          map[string]interface{} `json:"dps"`
	IP           string                 `json:"ip"`
	ActiveTime   int                    `json:"activeTime"`
	CategoryCode string                 `json:"categoryCode"`
	ModuleMap    DeviceModuleMapType    `json:"moduleMap"`
	SevAttribute string                 `json:"devAttribute"`
	Name         string                 `json:"name"`
	TimezoneID   string                 `json:"timezoneId"`
	Category     string                 `json:"category"`
	LocalKey     string                 `json:"localKey"`
}

func (t *ResultType) GetDeviceList() []DeviceType {
	var _l []DeviceType
	for _, r := range t.resultList {
		var _t DeviceType
		b, _ := json.Marshal(r)
		json.Unmarshal(b, &_t)
		_l = append(_l, _t)
	}
	return _l
}

type UDPMessage struct {
	IP         string `json:"ip"`
	GwID       string `json:"gwId"`
	Active     int    `json:"active"`
	Ability    int    `json:"ability"`
	Mode       int    `json:"mode"`
	Encrypt    bool   `json:"encrypt"`
	ProductKey string `json:"productKey"`
	Version    string `json:"version"`
}
