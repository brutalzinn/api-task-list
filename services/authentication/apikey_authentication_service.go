package authentication_service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	converter_util "github.com/brutalzinn/api-task-list/services/utils/converter"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
	"github.com/google/uuid"
)

func CreateUUID() string {
	uuid := uuid.New().String()
	uuidNormalized := strings.Replace(uuid, "-", "", -1)
	return uuidNormalized
}
func CreateApiHash(user_id int64, appName string, uuid string, expireAt string) (keyhash string, err error) {
	keyhash, err = crypt_util.Encrypt(fmt.Sprintf("%d#%s#%s#%s", user_id, appName, uuid, expireAt))
	if err != nil {
		return "", err
	}
	return keyhash, nil
}
func getApiKeyInfo(apiKeyDescrypted string) (user_id int64, appName string, expireAt string, err error) {
	apikeyformat := strings.Split(apiKeyDescrypted, "#")
	if len(apikeyformat) != 4 {
		return 0, "", "", nil
	}
	user_id, _ = strconv.ParseInt(apikeyformat[0], 10, 64)
	appName = apikeyformat[1]
	expireAt = apikeyformat[3]
	return user_id, appName, expireAt, nil
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
