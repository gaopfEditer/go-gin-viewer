package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
	"testing"
)

func TestGenerateRSA(t *testing.T) {
	// 生成 RSA 私钥 (4096 位)
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatalf("生成私钥失败: %v", err)
	}

	// 将私钥保存为 PKCS#1 格式的 PEM 文件
	if err := saveKeyToFile("private_key.pem", "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey)); err != nil {
		log.Fatalf("保存私钥失败: %v", err)
	}
	log.Println("私钥已保存到 private_key.pem")

	// 提取公钥并保存为 PKCS#1 格式的 PEM 文件
	publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	if err := saveKeyToFile("public_key.pem", "RSA PUBLIC KEY", publicKeyBytes); err != nil {
		log.Fatalf("保存公钥失败: %v", err)
	}
	log.Println("公钥已保存到 public_key.pem")
}

// saveKeyToFile 将密钥保存为 PEM 格式的文件
func saveKeyToFile(filename, pemType string, keyBytes []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	pemBlock := &pem.Block{
		Type:  pemType,
		Bytes: keyBytes,
	}

	if err := pem.Encode(file, pemBlock); err != nil {
		return err
	}

	return nil
}
