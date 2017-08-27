package main

import (
	"crypto/rand"
	"fmt"
)

func newUUID() (string, error) {
	uuid := make([]byte, 16)

	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	uuid[8] = uuid[8]&0xbf | 0x80
	uuid[6] = uuid[6]&0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
