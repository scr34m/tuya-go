package tuya

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"
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

var key rsa.PublicKey

func fromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("Bad number: %s" + base10)
	}
	return i
}

func importRSAKey(n, e string) {
	ei, _ := strconv.Atoi(e)

	key = rsa.PublicKey{
		N: fromBase10(n),
		E: ei,
	}
}

func computeRSA(data string) string {
	encrypted := new(big.Int)
	e := big.NewInt(int64(key.E))
	payload := new(big.Int).SetBytes([]byte(data))
	encrypted.Exp(payload, e, key.N)
	b := encrypted.Bytes()
	return strings.Repeat("0", 160-len(b)) + hex.EncodeToString(b)
}

func computeMd5(message string) string {
	m := md5.New()
	m.Write([]byte(message))
	return hex.EncodeToString(m.Sum(nil))
}

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
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
