package utils

import (
	"log"
	"os"
)

func ChangeWorkingDir() {
	if len(os.Args) >= 2 {
		if err := os.Chdir(os.Args[1]); err != nil {
			log.Fatalln("cannot change dir", err)
		}
	}
}

func WriteFile(path string, content string) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer closeFile(f)

	_, err = f.WriteString(content)
	if err != nil {
		log.Fatalln(err)
	}
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalln(err)
	}
}
