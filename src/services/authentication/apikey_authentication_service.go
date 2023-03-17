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
	crypt_util "github.com/brutalzinn/api-task-list/utils/crypt"
	"github.com/google/uuid"
)

func CreateUUID() string {
	uuid := uuid.New().String()
	return uuid
}
func CreateRandomFactor() (result string) {
	b := make([]byte, 4)
	rand.Read(b)
	result = hex.EncodeToString(b)
	return result
}
func CreateApiKeyCrypt(keyId string, randomFactor string) (keyhash string, err error) {
	keyhash, err = crypt_util.Encrypt(fmt.Sprintf("%s#%s", keyId, randomFactor))
	if err != nil {
		return "", err
	}
	return keyhash, nil
}
func CreateApiPrefix(apiKeyCrypt string, appName string) string {
	return fmt.Sprintf("%s-%s", appName, apiKeyCrypt)
}

func VerifyApiKey(apiKeyHeaderValue string) (*database_entities.ApiKey, error) {
	apiKeyCrypt, err := removeApiPrefix(apiKeyHeaderValue)
	if err != nil {
		return nil, errors.New("Api key invalid")
	}
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
	isKeyInvalid := crypt_util.CheckPasswordHash(apiKeyHeaderValue, apiKey.ApiKey)
	if isKeyInvalid == false {
		return nil, errors.New("Api key invalid")
	}
	return &apiKey, err
}
func getApiKeyInfo(apiKeyDescrypted string) (keyId string, err error) {
	apikeyformat := strings.Split(apiKeyDescrypted, "#")
	if len(apikeyformat) != 2 {
		return "", nil
	}
	keyId = apikeyformat[0]
	return keyId, nil
}

func removeApiPrefix(apiKeyCrypt string) (string, error) {
	apikeyformat := strings.Split(apiKeyCrypt, "-")
	if len(apikeyformat) != 2 {
		return "", nil
	}
	apiKey := apikeyformat[1]
	return apiKey, nil
}

func isKeyExpired(expireAt time.Time) bool {
	currentDateTime := time.Now()
	if expireAt.After(currentDateTime) {
		return false
	}
	return true
}
