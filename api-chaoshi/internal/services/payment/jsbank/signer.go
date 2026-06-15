package jsbank

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/crypto/pkcs12"
)

func resolveFilePath(pathValue string) string {
	trimmedPath := strings.TrimSpace(pathValue)
	if trimmedPath == "" {
		return ""
	}
	if filepath.IsAbs(trimmedPath) {
		return trimmedPath
	}

	workDir, err := os.Getwd()
	if err != nil {
		return trimmedPath
	}
	return filepath.Join(workDir, trimmedPath)
}

func loadPrivateKeyFromPFX(pathValue, password string) (*rsa.PrivateKey, error) {
	filePath := resolveFilePath(pathValue)
	pfxData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取 PFX 失败: %w", err)
	}

	privateKey, _, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, fmt.Errorf("解析 PFX 失败: %w", err)
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("PFX 私钥类型不是 RSA")
	}

	return rsaKey, nil
}

func loadPublicKeyFromCert(pathValue string) (*rsa.PublicKey, error) {
	filePath := resolveFilePath(pathValue)
	certData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取公钥证书失败: %w", err)
	}

	block, _ := pem.Decode(certData)
	if block == nil {
		return nil, fmt.Errorf("解析公钥证书失败")
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("读取公钥证书内容失败: %w", err)
	}

	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥证书类型不是 RSA")
	}

	return publicKey, nil
}

func buildSigningString(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		if key == "sign" || key == "signType" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	pairs := make([]string, 0, len(keys))
	for _, key := range keys {
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, params[key]))
	}

	return strings.Join(pairs, "&")
}

func signPayload(params map[string]string, privateKey *rsa.PrivateKey) (string, error) {
	signingString := buildSigningString(params)
	hash := sha1.Sum([]byte(signingString))

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return "", fmt.Errorf("签名失败: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func verifyPayload(params map[string]string, publicKey *rsa.PublicKey) error {
	signatureBase64 := strings.TrimSpace(params["sign"])
	if signatureBase64 == "" {
		return fmt.Errorf("缺少签名字段")
	}

	signingString := buildSigningString(params)
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("解码签名失败: %w", err)
	}

	hash := sha1.Sum([]byte(signingString))
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hash[:], signature); err != nil {
		return fmt.Errorf("验签失败: %w", err)
	}

	return nil
}
