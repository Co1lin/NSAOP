package utils

import (
	"fmt"
	"math/rand"
)

func AddIDToPath(path string, id interface{}) string {
	return fmt.Sprintf(path, fmt.Sprint(id))
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz~`!@#$%^&*()-_=+{[}];:?/><,."
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < l; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

func GetRandomStringWithLetters(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < l; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}