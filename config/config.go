package config

import (
	"fmt"
	"github.com/btstar/shorturl/models"
	"github.com/jinzhu/configor"
)

func LoadConfig() models.Config {
	cfg := models.Config{}
	err := configor.Load(&cfg, "config.yml")
	fmt.Printf("err:%v, config: %#v\n", err, cfg)
	return cfg
}
