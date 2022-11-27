package tuya

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Config struct {
	CertSign  string
	Secret    string
	AppSecret string
}

type Tuya struct {
	config   Config
	deviceId string
	clientId string

	sessionId string
	groupId   string

	endpoint string
}

func New(config Config, deviceId, clientId string) (*Tuya, error) {
	return &Tuya{
		config:   config,
		deviceId: deviceId,
		clientId: clientId,
		endpoint: "https://a1.tuyaeu.com/api.json",
	}, nil
}

func (tuya *Tuya) SessionIdSet(s string) {
	tuya.sessionId = s
}

func (tuya *Tuya) GroupIdSet(s string) {
	tuya.groupId = s
}

func (tuya *Tuya) Login(email string, password string, countyCode int) (string) {
	tokenData := struct {
		CountryCode int    `json:"countryCode"`
		Email       string `json:"email"`
	}{
		countyCode,
		email,
	}

	res := tuya.Call("tuya.m.user.email.token.create", "1.0", tokenData, false)
	if !res.Success {
		fmt.Println(res.ErrorMessage)
		return ""
	}

	rsa := res.GetRsaType()

	importRSAKey(rsa.PublicKey, rsa.Exponent)

	encryptedPass := computeRSA(computeMd5(password))

	type passwordDataOptions struct {
		Group int `json:"group"`
	}

	passwordData := struct {
		CountryCode int                 `json:"countryCode"`
		Email       string              `json:"email"`
		Ifencrypt   int                 `json:"ifencrypt"`
		Options     passwordDataOptions `json:"options"`
		Password    string              `json:"passwd"`
		Token       string              `json:"token"`
	}{
		countyCode,
		email,
		1,
		passwordDataOptions{Group: 1},
		encryptedPass,
		rsa.Token,
	}

	res = tuya.Call("tuya.m.user.email.password.login", "1.0", passwordData, false)
	if !res.Success {
		fmt.Println(res.ErrorMessage)
		return ""
	}

	login := res.GetLoginType()

	tuya.sessionId = login.SID

	return tuya.sessionId
}

func (tuya *Tuya) Call(action, version string, postData interface{}, signed bool) ResultType {

	apiEtVersion := true

	if signed && tuya.sessionId == "" {
		panic("Need login first")
	}

	data := map[string]string{}
	data["a"] = action
	data["deviceId"] = tuya.deviceId
	data["os"] = "Linux"
	data["lang"] = "en"
	data["clientId"] = tuya.clientId
	data["v"] = version
	data["time"] = fmt.Sprintf("%d", time.Now().Unix())

	if postData != nil {
		postDataJSON, err := json.Marshal(postData)
		if err != nil {
			panic(err)
		}
		data["postData"] = string(postDataJSON)
	}

	if tuya.groupId != "" {
		data["gid"] = tuya.groupId
	}

	if apiEtVersion {
		data["et"] = "0.0.1"
		data["ttid"] = "tuya"
		data["appVersion"] = "3.8.0"
	}

	if signed {
		data["sid"] = tuya.sessionId
	}

	data["requestId"] = generateUUID()

	data["sign"] = tuya.sign(data)

	post := url.Values{}
	if postData != nil {
		post.Add("postData", data["postData"])
	}

	params := url.Values{}
	for k, v := range data {
		params.Add(k, v)
	}

	resp, err := http.PostForm(tuya.endpoint+"?"+params.Encode(), post)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return NewResultType(body)
}

func (tuya *Tuya) sign(pairs map[string]string) string {
	keys := make([]string, 0, len(pairs))
	for k := range pairs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	data := ""
	for _, k := range keys {
		if pairs[k] == "" || !neetToSign(k) {
			continue
		}

		if len(data) > 0 {
			data += "||"
		}
		data += k
		data += "="

		if k == "postData" {
			data += tuya.hash(pairs[k])
		} else {
			data += pairs[k]
		}
	}

	return computeHmac256(data, tuya.config.CertSign+"_"+tuya.config.Secret+"_"+tuya.config.AppSecret)
}

func (tuya *Tuya) hash(data string) string {
	hash := computeMd5(data)
	return hash[8:16] + hash[0:8] + hash[24:32] + hash[16:24]
}
