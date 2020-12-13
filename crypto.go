package tuya

import (
	"crypto/aes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"
	"strings"
)

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

func aesDecrypt(data []byte, key []byte) ([]byte, error) {
	n := len(data)
	block, er2 := aes.NewCipher(key)
	if er2 != nil {
		return []byte{}, er2
	}
	bs := block.BlockSize()
	if n%bs != 0 && n < 16 {
		return []byte{}, errors.New("Bad ciphertext len")
	}
	cleartext := make([]byte, n)
	for i := 0; i < n; i = i + bs {
		block.Decrypt(cleartext[i:i+bs], data[i:i+bs])
	}
	// remove padding
	p := int(cleartext[n-1])
	if p < 0 || p > bs {
		return []byte{}, errors.New("Bad padding")
	}
	return cleartext[:n-p], nil
}

func aesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	remain := len(data) % bs
	if remain == 0 {
		remain = bs
	}
	padd := make([]byte, bs-remain)
	for i := range padd {
		padd[i] = byte(bs - remain)
	}
	data = append(data, padd...)
	ciphertext := make([]byte, len(data))
	for i := 0; i < len(data); i = i + bs {
		block.Encrypt(ciphertext[i:i+bs], data[i:i+bs])
	}
	return ciphertext, nil
}
