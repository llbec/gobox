package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/pbkdf2"
)

// EncryptPrivateKey 使用 AES-GCM + PBKDF2 加密私钥
func EncryptPrivateKey(privateKey string, password string) (string, error) {
	// 生成密钥
	key := pbkdf2.Key([]byte(password), []byte("salecontract_salt"), 4096, 32, sha256.New)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := aesgcm.Seal(nonce, nonce, []byte(privateKey), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey 解密并返回明文私钥字符串
func DecryptPrivateKey(encrypted string, password string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	key := pbkdf2.Key([]byte(password), []byte("salecontract_salt"), 4096, 32, sha256.New)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(data) < aesgcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce := data[:aesgcm.NonceSize()]
	ct := data[aesgcm.NonceSize():]

	plain, err := aesgcm.Open(nil, nonce, ct, nil)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}
