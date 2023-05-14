package tools

import (
	"chatserver-api/pkg/logger"
	"os"
)

func CreatePath(path ...string) {
	for _, v := range path {
		if _, err := os.Stat(v); err != nil {
			err := os.MkdirAll(v, 0711)
			if err != nil {
				logger.Errorf("Error creating directory")
				return
			}
		}
	}
}
