package config

import "github.com/dovics/wx-demo/util/config"

func init() {
	config.Add("wx", config.StrMap{
		"appid":  config.Env("WX_APPID", ""),
		"secret": config.Env("WX_SECRET", ""),
	})
}
