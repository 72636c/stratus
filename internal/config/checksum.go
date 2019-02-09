package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func CalculateChecksum(model interface{}) (string, error) {
	data, err := json.Marshal(model)
	if err != nil {
		return "", err
	}

	checksum := fmt.Sprintf("%x", sha256.Sum256(data))

	return checksum, nil
}
