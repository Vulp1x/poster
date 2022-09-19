package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"strconv"
	"time"
)

func EncryptPassword(pass []byte, publicKey, publicKeyID string) (string, error) {
	sessionKey := make([]byte, 32)
	_, _ = rand.Read(sessionKey)

	rsaEncrypted, err := encryptSessionKey(publicKey, sessionKey)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher(sessionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create aes cipher block: %v", err)
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", fmt.Errorf("failed to create gcm: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to create nonce: %v", err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	additionalData := []byte(timestamp)
	tagSize := gcm.Overhead()
	var aesDst = make([]byte, 0, len(pass)+tagSize)
	aesEncrypted := gcm.Seal(aesDst, nonce, pass, additionalData)

	return combineEncrypted(timestamp, publicKeyID, pass, aesEncrypted, nonce, rsaEncrypted)
}

func combineEncrypted(
	timestamp, publicKeyID string,
	pass, aesEncrypted, nonce, rsaEncrypted []byte,
) (string, error) {
	tag := aesEncrypted[len(pass):]
	aesEncrypted = aesEncrypted[:len(pass)]

	buf := new(bytes.Buffer)

	buf.WriteString("#PWD_INSTAGRAM:4:")
	buf.WriteString(timestamp)
	buf.Write([]byte(":\x01"))

	err := binary.Write(buf, binary.BigEndian, publicKeyID)
	if err != nil {
		return "", fmt.Errorf("failed to write public key id: %v", err)
	}

	buf.Write(nonce)

	err = binary.Write(buf, binary.LittleEndian, len(rsaEncrypted))
	if err != nil {
		return "", fmt.Errorf("failed to write len of rsa encrypted: %v", err)
	}

	buf.Write(rsaEncrypted)
	buf.Write(tag)
	buf.Write(aesEncrypted)

	// 	cipher_aes = AES.new(session_key, AES.MODE_GCM, nonce)
	// 	cipher_aes.update(timestamp.encode())
	// 	aes_encrypted, tag = cipher_aes.encrypt_and_digest(password.encode("utf8"))
	// 	payload = base64.b64encode(b''.join([
	// 		b"\x01",
	// 		publickeyid.to_bytes(1, byteorder='big'),
	// 		nonce,
	// 		len(rsa_encrypted).to_bytes(2, byteorder='little')
	// 		rsa_encrypted,
	// 		tag,
	// 		aes_encrypted
	// ]))
	// 	return f"#PWD_INSTAGRAM:4:{timestamp}:{payload.decode()}"

	return buf.String(), nil
}

// encryptSessionKey зашифровывает sessionKey публичным ключом publicKey
// используется rsa.EncryptPKCS1v15
func encryptSessionKey(publicKey string, sessionKey []byte) ([]byte, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key '%s' : %v", publicKey, err)
	}

	block, _ := pem.Decode(decodedPublicKey)
	if block == nil {
		return nil, fmt.Errorf("failed to decode public key from '%s': %v", publicKey, err)
	}
	parseResult, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %v", err)
	}

	rsaKey, ok := parseResult.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed cast parsed key (%T) to rsa.PublicKey", parseResult)
	}

	rsaEncrypted, err := rsa.EncryptPKCS1v15(rand.Reader, rsaKey, sessionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt session key using rsa: %v", err)
	}
	return rsaEncrypted, nil
}
