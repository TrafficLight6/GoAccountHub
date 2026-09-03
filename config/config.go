package config

import (
	"github.com/TrafficLight6/GoAccountHub/file"
	jsonOperate "github.com/TrafficLight6/GoAccountHub/json"
)

type Config struct {
	Port             string `json:"port"`
	DatabaseName     string `json:"database_name"`
	DatabaseHost     string `json:"database_host"`
	DatabasePort     string `json:"database_port"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`

	RootAdminPasswordHash string `json:"root_admin_password_hash"`
	RootAdminUUHash       string `json:"root_admin_uu_hash"`

	//Bool ONLY
	SwitchConfig SwitchConfig `json:"switch_config"`
}

type SwitchConfig struct {
	AllowMultiCharacter bool `json:"allow_multi_character"`
}

func GetConfig(path string) Config {
	config := Config{}
	jsonOperate.Unmarshal(file.Read(path), &config, " on "+path)
	return config
}

func SaveConfig(path string, config Config) {
	file.Write(path, jsonOperate.Marshal(config))
}
