package greetings

import (
	"testing"
	"regexp"
)

func TestHelloName(t *testing.T) {
	InitSeed()

	name := "Flo"
	excpected := regexp.MustCompile(`\b` + name)
	msg, err := Hello(name)

	if err != nil {
		t.Fatalf("Hello(Flo) errored: %v", err)
	}

	if !excpected.MatchString(msg) {
		t.Fatalf("Hello(Flo) is %q, wanted match for %#q", msg, excpected)
	}
}

func TestHelloEmpty(t *testing.T) {
	InitSeed()

	msg, err := Hello("")

	if nil == err {
		t.Fatalf("Hello('') should error, got %q", msg)
	}
}
