package kugou

import "testing"

func TestCryptoRSAEncryptPKCS1Hex(t *testing.T) {
	encrypted, err := CryptoRSAEncryptPKCS1Hex(map[string]string{"token": "test"}, PublicRASKey)
	if err != nil {
		t.Fatalf("CryptoRSAEncryptPKCS1Hex() error = %v", err)
	}
	if encrypted == "" {
		t.Fatal("CryptoRSAEncryptPKCS1Hex() returned an empty ciphertext")
	}
}
