package apikey_util

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	crypt_util "github.com/brutalzinn/api-task-list/services/utils/crypt"
	"github.com/google/uuid"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func CreateApiKey(user_id int64, appName string) (string, error) {
	uuid := uuid.New().String()
	uuidNormalized := strings.Replace(uuid, "-", "", -1)
	keyhash, err := crypt_util.Encrypt(fmt.Sprintf("%d-%s-%s", user_id, appName, uuidNormalized))
	if err != nil {
		return "", err
	}
	return keyhash, nil
}
func GetApiKeyInfo(apiKeyDescrypted string) (int64, string, error) {
	apikeyformat := strings.Split(apiKeyDescrypted, "-")
	if len(apikeyformat) != 3 {
		return 0, "", nil
	}

	user_id, _ := strconv.ParseInt(apikeyformat[0], 10, 64)
	appName := apikeyformat[1]
	return user_id, appName, nil
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
