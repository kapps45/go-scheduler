package env

import "os"

// CheckEnv возвращает значение переменной окружения или defaultValue, если она не задана.
func CheckEnv(envString, defaultValue string) string {
	res := os.Getenv(envString)
	if res == "" {
		return defaultValue
	}
	return res
}
