package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

// ErrHMACNotEqual ошибка, возвращаемая сервисом, если  hmacs не равны
var ErrHMACNotEqual = errors.New("hmacs not equal")

// Sign подпись данных data ключом key с использованием hmac sha256
func Sign(data []byte, key string) []byte {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write(data)
	signedData := hash.Sum(nil)
	return signedData
}

// SignToString подпись данных service.Sign() и преобразование её в строку
func SignToString(data []byte, key string) string {
	signedData := Sign(data, key)
	return hex.EncodeToString(signedData)
}

// Equal проверка на совпадение подписей
func Equal(msg string, bytes []byte, key string) error {
	msgBytes, err := hex.DecodeString(msg)
	if err != nil {
		return fmt.Errorf("failed to hex decode: %w", err)
	}

	signedBytes := Sign(bytes, key)

	if !hmac.Equal(msgBytes, signedBytes) {
		return ErrHMACNotEqual
	}
	return nil
}
