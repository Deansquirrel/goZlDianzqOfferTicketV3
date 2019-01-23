package object

import "encoding/json"

type SysConfig struct {
	Total   total   `toml:"total"`
	PeiZhDb peiZhDb `toml:"peiZhDb"`
}

type total struct {
	IsDebug bool `toml:"isDebug"`
	Port    int  `toml:"port"`

	MaxTicketNum int    `toml:"maxTicketNum"`
	SnoWorkerId  int    `toml:"snoWorkerId"`
	AppId        string `toml:"appid"`
	JPeiZh       string `toml:"jpeizh"`
}

type peiZhDb struct {
	Server   string `toml:"server"`
	Port     int    `toml:"port"`
	DbName   string `toml:"dbName"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}

func (sc *SysConfig) GetConfigStr() (string, error) {
	sConfig, err := json.Marshal(sc)
	if err != nil {
		return "", err
	}
	return string(sConfig), nil
}
