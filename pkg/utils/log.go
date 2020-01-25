package utils

import (
	"fmt"
	"github.com/mamezou-tech/concourse-k8s-resource/pkg/models"
	"log"
	"os"
)

var debug *log.Logger

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	debug = log.New(os.Stderr, "[Debug] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Debug(s *models.Source, v ...interface{}) {
	if s.Debug {
		debug.Println(fmt.Sprint(v))
	}
}
