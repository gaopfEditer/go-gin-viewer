package str

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

const InvitationCryptKey = "iEA8wIl8J6vA5OVWrSt6UNiXVLaLOykH5YQKNI"

// generateKey 生成一个基于密码和盐的 AES 密钥
func generateKey(salt string) []byte {
	key := sha256.Sum256([]byte(salt))
	return key[:]
}

// customBase64Encode 使用自定义字符集编码以避免 `_` 字符
func customBase64Encode(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "/", ".")
	encoded = strings.ReplaceAll(encoded, "=", "")
	return encoded
}

// customBase64Decode 使用自定义字符集解码
func customBase64Decode(encoded string) ([]byte, error) {
	encoded = strings.ReplaceAll(encoded, "-", "+")
	encoded = strings.ReplaceAll(encoded, ".", "/")
	if len(encoded)%4 != 0 {
		encoded += strings.Repeat("=", 4-(len(encoded)%4))
	}
	return base64.StdEncoding.DecodeString(encoded)
}

// Encrypt 加密给定的字符串
func Encrypt(plainText, salt string) (string, error) {
	key := generateKey(salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return customBase64Encode(cipherText), nil
}

// Decrypt 解密给定的字符串
func Decrypt(encryptedText, salt string) (string, error) {
	key := generateKey(salt)
	cipherText, err := customBase64Decode(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
