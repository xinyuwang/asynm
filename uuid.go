package asynm

import (
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
)

// check later
// use -1 to return full uuid
func GenerateShortUuid(l int) (string, error) {
	if uuid, err := uuid.NewRandom(); err != nil {
		return "", err
	} else {
		escaper := strings.NewReplacer("9", "99", "-", "90", "_", "91")
		shortUuid := escaper.Replace(base64.RawURLEncoding.EncodeToString([]byte(uuid.String())))
		if l < 0 || l > len(shortUuid) {
			l = len(shortUuid)
		}
		return shortUuid[:l], nil
	}
}

func GenerateUuid() (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

func IsValidUuid(keyId string) bool {
	if len(keyId) == 0 {
		return false
	}
	_, err := uuid.Parse(keyId)
	if err != nil {
		return false
	}
	return true
}

func IsValidUuids(ids []string) bool {
	if len(ids) == 0 {
		return false
	}
	for _, id := range ids {
		if IsValidUuid(id) == false {
			return false
		}
	}
	return true
}
