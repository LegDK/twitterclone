package config

import (
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

type database struct {
	URL string
}

type Config struct {
	Database database
	JWT      jwt
}

type jwt struct {
	Secret string
	Issuer string
}

func LoadEnv(fileName string) {
	re := regexp.MustCompile(`^(.*` + "twitterclone" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + "/" + fileName)
	if err != nil {
		godotenv.Load()
	}

}

func New() *Config {
	return &Config{
		Database: database{
			URL: os.Getenv("DATABASE_URL"),
		},
		JWT: jwt{
			Secret: os.Getenv("JWT_SECRET"),
			Issuer: os.Getenv("DOMAIN"),
		},
	}
}
