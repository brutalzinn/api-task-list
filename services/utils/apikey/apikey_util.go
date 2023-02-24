package apikey_util

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	converter_util "github.com/brutalzinn/api-task-list/services/utils/converter"
	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
	"github.com/google/uuid"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func CreateUUID() string {
	uuid := uuid.New().String()
	uuidNormalized := strings.Replace(uuid, "-", "", -1)
	return uuidNormalized
}
func CreateApiHash(user_id int64, appName string, uuid string, expireAt string) (string, error) {
	keyhash, err := crypt_util.Encrypt(fmt.Sprintf("%d#%s#%s#%s", user_id, appName, uuid, expireAt))
	if err != nil {
		return "", err
	}
	return keyhash, nil
}
func GetApiKeyInfo(apiKeyDescrypted string) (user_id int64, appName string, expireAt string, err error) {
	apikeyformat := strings.Split(apiKeyDescrypted, "#")
	if len(apikeyformat) != 4 {
		return 0, "", "", nil
	}
	user_id, _ = strconv.ParseInt(apikeyformat[0], 10, 64)
	appName = apikeyformat[1]
	expireAt = apikeyformat[3]
	return user_id, appName, expireAt, nil
}

func IsKeyExpired(expireAt string) bool {
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

// thahnks to https://gist.github.com/micheltlutz/8398f801abf3ac61ff002dd5be7caeb4
func Normalize(oldText string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(IsMn), norm.NFC)
	newText, _, err := transform.String(t, oldText)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(newText))
	return strings.ToLower(newText)
}

func IsMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) || unicode.IsSpace(r)
}
