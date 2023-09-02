package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Mode int

const (
	ModeDev Mode = iota + 1
	ModeProd
)

type Config struct {
	Mode Mode

	DBPath string

	BEHost   string
	BEPort   string
	BESecret []byte
}

func New() Config {
	var mode Mode
	switch str := os.Getenv("MODE"); str {
	case "dev":
		mode = ModeDev
	case "prod":
		mode = ModeProd
	case "":
		mode = ModeDev
	default:
		log.Fatalf("ERROR: config: expected MODE=dev|prod; got MODE=%s\n", str)
	}

	envMap := map[string]string{}
	switch mode {
	case ModeDev:
		var err error
		envMap, err = godotenv.Read(".env")
		if err != nil {
			log.Fatalln("ERROR: config:", err)
		}
	case ModeProd:
		for _, env := range os.Environ() {
			parts := strings.SplitN(env, "=", 2)
			envMap[parts[0]] = parts[1]
		}
	}

	return Config{
		Mode:     mode,
		DBPath:   getenvOrExit(envMap, "DB_PATH"),
		BEHost:   getenvOrExit(envMap, "BE_HOST"),
		BEPort:   getenvOrExit(envMap, "BE_PORT"),
		BESecret: []byte(getenvOrExit(envMap, "BE_SECRET")),
	}
}

func getenvOrExit(envMap map[string]string, key string) string {
	env, ok := envMap[key]
	if !ok {
		log.Fatalf("ERROR: config: cannot find variable '%s' in the environment\n", key)
	}

	return env
}
