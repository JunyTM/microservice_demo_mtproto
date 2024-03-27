package utils

import "github.com/gofrs/uuid"

func GenAuthKey() (string, error) {
	key, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return key.String(), nil
}
