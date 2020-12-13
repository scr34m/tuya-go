package tuya

import (
	"crypto/rand"
	"fmt"
)

func neetToSign(s string) bool {
	var keysToSign = []string{"a", "v", "lat", "lon", "lang", "deviceId", "imei", "imsi", "appVersion", "ttid", "isH5", "h5Token", "os", "clientId", "postData", "time", "requestId", "n4h5", "sid", "sp", "et"}
	for _, v := range keysToSign {
		if v == s {
			return true
		}
	}
	return false
}

func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
