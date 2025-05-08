package config

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

const projectDirName = "7-solutions-backend-challenge"

func loadEnv() error {
	var (
		err      error
		re       = regexp.MustCompile(`^(.*` + projectDirName + `)`)
		cwd, _   = os.Getwd()
		rootPath = string(re.Find([]byte(cwd)))
	)

	if testFile(".env") {
		log.Println("Test file pass: Using environment from .env")
		err = godotenv.Load(".env")
	} else {
		f := buildFilePath(rootPath, ".env")
		log.Println("Test file failed: Using environment from", f)
		err = godotenv.Load(f)

		return err
	}

	return err
}

func testFile(p string) bool {
	info, err := os.Stat(p)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func buildFilePath(rootPath string, filename string) string {
	return fmt.Sprintf("%s/%s", rootPath, filename)
}
