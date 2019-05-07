package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	keyPrefix = "rpaasv2"
)

func keyName(name string) string {
	return fmt.Sprintf("%s-%s", keyPrefix, name)
}

func Value(key string) string {
	return os.Getenv(keyName(key))
}

func StringMap(key string) map[string]string {
	val, isSet := os.LookupEnv(keyName(key))
	if !isSet {
		return nil
	}
	var ret map[string]string
	json.Unmarshal([]byte(val), &ret)
	return ret
}