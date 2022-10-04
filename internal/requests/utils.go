package requests

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/url"
	"strconv"

	"github.com/google/uuid"
)

const instHost = "i.instagram.com"

const defaultCountryCode = 1 // default is USA

func generateJazoest(phoneId uuid.UUID) string {
	var sum int32
	for _, s := range phoneId.String() {
		sum += s
	}

	return strconv.FormatInt(int64(sum), 10)
}

func generateSignature(data map[string]string) *bytes.Buffer {
	dataBytes, _ := json.Marshal(data)
	buff := &bytes.Buffer{}

	buff.WriteString("signed_body=SIGNATURE.")

	buff.WriteString(url.QueryEscape(string(dataBytes)))

	return buff
}

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateCSRFToken() string {
	buf := &bytes.Buffer{}
	for i := 0; i < 64; i++ {
		buf.WriteByte(alphabet[rand.Intn(len(alphabet))])
	}

	return buf.String()
}
