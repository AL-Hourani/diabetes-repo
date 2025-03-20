package config

import (
	"fmt"
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost string
	Port       string

	DBUser     string
	Host       string
	DBPort     string
	DBHost     string
	DBPassword string
	DBAddress  string
	DBName     string
	JWTExpirationInSecond int64
	JWTSecret  string
}

var Envs = initConfig()

func initConfig() Config {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("error can load .env file so we use the loacl env")
	}
	return Config{
		PublicHost: GetEnv("PUBLIC_HOST", ""),
		DBPort:GetEnv("DB_PORT" ,""),
		DBHost:GetEnv("DB_HOST" ,""),
		DBUser: GetEnv("DB_USER",""),
		DBPassword: GetEnv("DB_PASSWORD", ""),
		DBAddress: fmt.Sprintf("%s:%s", GetEnv("DB_HOST","") , GetEnv("DB_PORT","")),
		DBName: GetEnv("DB_NAME",""),
		JWTSecret: GetEnv("JWT_SECRET",""),
		JWTExpirationInSecond: getEnvAsInt("JWT_EXP", 3600 * 24 * 7),
	}
}

func GetEnv(key, fallback string) string {
	if value , ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}


func getEnvAsInt(key string , fallback int64) int64{
	if value , ok := os.LookupEnv(key); ok {
		i , err := strconv.ParseInt(value , 10 ,64)
		if err != nil {
			return fallback
		}
		return i
	}

	return fallback
}