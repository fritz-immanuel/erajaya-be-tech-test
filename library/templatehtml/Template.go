package templatehtml

import (
	"io"
	"log"
	"os"
)

func TemplateDefault() string {
	file, err := os.Open("html/DefaultTemplate.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	str, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	return string(str)
}
