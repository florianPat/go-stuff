package hello

import (
	"fmt"
	"log"
	"example/greetings"
)

func main() {
	log.SetPrefix("greetings: ")
	log.SetFlags(0)
	greetings.InitSeed()

	message, err := greetings.Hello("Flo")
	if nil != err {
		log.Fatal(err)
	}

	fmt.Print(message)
}
