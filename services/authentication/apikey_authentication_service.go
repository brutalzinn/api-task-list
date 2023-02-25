package authentication_service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	database_entities "github.com/brutalzinn/api-task-list/models/database"
	apikey_service "github.com/brutalzinn/api-task-list/services/database/apikey"
	converter_util "github.com/brutalzinn/api-task-list/services/utils/converter"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
	"github.com/google/uuid"
)

func CreateUUID() string {
	uuid := uuid.New().String()
	return uuid
}
func CreateRandomFactor() (result string) {
	b := make([]byte, 4) //equals 8 characters
	rand.Read(b)
	result = hex.EncodeToString(b)
	return result
}
func CreateApiHash(keyId string) (keyhash string, err error) {
	keyhash, err = crypt_util.Encrypt(fmt.Sprintf("%s", keyId))
	if err != nil {
		return "", err
	}
	return keyhash, nil
}
func VerifyApiKey(apiKeyCrypt string) (*database_entities.ApiKey, error) {
	decrypt, err := crypt_util.Decrypt(apiKeyCrypt)
	if err != nil {
		return nil, errors.New("Api key invalid")
	}
	keyId, err := getApiKeyInfo(decrypt)
	apiKey, err := apikey_service.Get(keyId)
	if err != nil {
		return nil, errors.New("Api key invalid")
	}
	isKeyExpired := isKeyExpired(apiKey.ExpireAt)
	if isKeyExpired {
		return nil, errors.New("Api key expired")
	}
	isKeyInvalid := crypt_util.CheckPasswordHash(apiKeyCrypt, apiKey.ApiKey)
	if isKeyInvalid == false {
		return nil, errors.New("Api key invalid")
	}
	return &apiKey, err
}
func getApiKeyInfo(apiKeyDescrypted string) (keyId string, err error) {
	apikeyformat := strings.Split(apiKeyDescrypted, "#")
	if len(apikeyformat) != 4 {
		return "", nil
	}
	keyId = apikeyformat[0]
	return keyId, nil
}

func isKeyExpired(expireAt string) bool {
	date, err := converter_util.ToDateTime(expireAt)
	if err != nil {
		fmt.Println(err)
		return true
	}
	currentDateTime := time.Now()
	if date.After(currentDateTime) {
		return false
	}
	return true
}
