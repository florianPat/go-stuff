package greetings

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func Hello(name string) (string, error) {
	if "" == name {
		return "", errors.New("empty name given")
	}

	message := fmt.Sprintln(randomFormat(), name)

	return message, nil
}

var randInstance *rand.Rand

func InitSeed() {
	randInstance = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randomFormat() string {
	if nil == randInstance {
		panic(errors.New("random device was not initialized"))
	}

	formats := []string{
		"Hello %v",
		"Monsieur %v",
		"Moin %v",
	}

	return formats[randInstance.Intn(len(formats))]
}

func Hellos(names []string) (map[string]string, error) {
	messages := make(map[string]string)

	for _, name := range names {
		message, err := Hello(name)
		if nil != err {
			return nil, err
		}

		messages[name] = message
	}

	return messages, nil
}