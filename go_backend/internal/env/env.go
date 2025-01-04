package env

import (
	"fmt"
	"os"
	"strconv"
	"github.com/joho/godotenv"

)

type ServerConfig struct {
	ServerURL string
	DBUrl string

}

func InitConfig() ServerConfig {
	godotenv.Load("secrets.env")

	dbUname := GetString("DB_USERNAME", "root")
	dbUrl := GetString("DB_URL", "URL")
	dbPort := GetString("DB_PORT", "5432")
	dbName := GetString("DB_NAME", "postgres")
	dbPassword := GetString("DB_PASSWORD", "admin")

	dbConnectionLink := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUname, dbPassword, dbUrl, dbPort, dbName)

	serverHost := GetString("HOST", "http://localhost")
	serverPort :=  GetString("PORT", "8000")

	return ServerConfig{
		ServerURL: fmt.Sprintf("%s:%s", serverHost, serverPort),
		DBUrl: dbConnectionLink,
	}
}

func GetString(key string, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return intVal
}