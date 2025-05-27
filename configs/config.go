package configs

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

const (
	androidAppMinimumVersion = "ANDROID_APP_MINIMUM_VERSION"
	iosAppMinimumVersion     = "IOS_APP_MINIMUM_VERSION"

	externalURL         = "EXTERNAL_URL"
	externalToken       = "EXTERNAL_TOKEN"
	externalAccessToken = "EXTERNAL_ACCESS_TOKEN"

	appUrl             = "APP_URL"
	serverName         = "SERVER_NAME"
	dbConnectionString = "DB_CONNECTION_STRING"

	portApps = "PORT_APPS"

	redisAddr     = "REDIS_ADDR"
	redisDB       = "REDIS_DB"
	redisPassword = "REDIS_PASSWORD"
	redisTimeOut  = "REDIS_TIME_OUT"

	sendWhatsappAPI   = "SEND_WHATSAPP_API"
	sendWhatsappToken = "SEND_WHATSAPP_TOKEN"

	teleBotToken = "TELE_BOT_TOKEN"
	teleGroupID  = "TELE_GROUP_ID"

	firebaseStorageBucketURL = "FIREBASE_BUCKET_URL"
	firebaseAuthFilePath     = "FIREBASE_AUTH_FILE_PATH"

	whitelistedIps = "WHITELISTED_IPS"
)

// TODO check mana yg masih dipakai
var (
	JwtActiveToken *string
)

// Config contains application configuration
type Config struct {
	// Minimum App versions
	AndroidAppMinimumVersion string
	IosAppMinimumVersion     string

	ExternalURL         string
	ExternalToken       string
	ExternalAccessToken string

	// DB
	DBConnectionString string

	// Misc
	PortApps       string
	ServerName     string
	WhitelistedIps string

	// Redis
	RedisAddr     string
	RedisDB       int
	RedisPassword string
	RedisTimeOut  int

	// WA
	SendWhatsappAPI   string
	SendWhatsappToken string

	// Telegram
	TeleBotToken string
	TeleGroupID  string

	// Firebase
	FirebaseStorageBucketURL string
	FirebaseAuthFilePath     string

	AppURL string
}

var config *Config

func getEnvOrDefault(env string, defaultVal string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaultVal
	}

	return e
}

// GetConfiguration , get application configuration based on set environment
func GetConfiguration() (*Config, error) {
	if config != nil {
		return config, nil
	}

	var dataENV []byte
	var err error

	dataENV, err = os.ReadFile(".env")
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, fmt.Errorf("failed to locate env file: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(dataENV, &result)

	redisDBi, err := strconv.Atoi(result[redisDB].(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis db: %v", err)
	}

	redisTimeOut, err := strconv.Atoi(result[redisTimeOut].(string)) // 3 days
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis timeout: %v", err)
	}

	config := &Config{
		AndroidAppMinimumVersion: result[androidAppMinimumVersion].(string),
		IosAppMinimumVersion:     result[iosAppMinimumVersion].(string),

		ExternalURL:         result[externalURL].(string),
		ExternalToken:       result[externalToken].(string),
		ExternalAccessToken: result[externalAccessToken].(string),

		DBConnectionString: result[dbConnectionString].(string),

		PortApps: result[portApps].(string),

		RedisAddr:     result[redisAddr].(string),
		RedisDB:       redisDBi,
		RedisPassword: result[redisPassword].(string),
		RedisTimeOut:  redisTimeOut,

		SendWhatsappAPI:   result[sendWhatsappAPI].(string),
		SendWhatsappToken: result[sendWhatsappToken].(string),

		// TELEGRAM
		TeleBotToken: result[teleBotToken].(string),
		TeleGroupID:  result[teleGroupID].(string),

		FirebaseAuthFilePath:     result[firebaseAuthFilePath].(string),
		FirebaseStorageBucketURL: result[firebaseStorageBucketURL].(string),

		AppURL:         result[appUrl].(string),
		WhitelistedIps: result[whitelistedIps].(string),
		ServerName:     result[serverName].(string),
	}

	return config, nil
}
