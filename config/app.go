package config

import "github.com/dovics/wx-demo/util/config"

func init() {
	config.Add("app", config.StrMap{
		"name": config.Env("APP_NAME", "github.com/dovics/wx-demo"),
		// multiple environments
		"env": config.Env("APP_ENV", "production"),
		// debug mode
		"debug": config.Env("APP_DEBUG", false),
		// port
		"port": config.Env("APP_PORT", "3000"),
		// BaseUrl
		"url": config.Env("APP_URL", "http://localhost:3000"),
	})
}
