package main

import (
	"os"
	"runtime"
	"strings"

	"code.google.com/p/gcfg"
	log "github.com/Sirupsen/logrus"
)

type GeneralConf struct {
	Debug bool
}

type Config struct {
	General  GeneralConf
	Mqtt     MqttConf     `gcfg:"mqforward-mqtt"`
	InfluxDB InfluxDBConf `gcfg:"mqforward-influxdb"`
}

func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func LoadConf(path string) (MqttConf, InfluxDBConf, error) {
	home := UserHomeDir()
	path = strings.Replace(path, "~", home, 1)

	var cfg Config
	err := gcfg.ReadFileInto(&cfg, path)
	if err != nil {
		return MqttConf{}, InfluxDBConf{}, err
	}

	em := MqttConf{}
	ei := InfluxDBConf{}

	if cfg.Mqtt == em {
		log.Fatal(`empty "mqforward-mqtt" configuration`)
	}
	if cfg.InfluxDB == ei {
		log.Fatal(`empty "mqforward-influxdb" configuration`)
	}

	if cfg.General.Debug {
		log.SetLevel(log.DebugLevel)
	}

	return cfg.Mqtt, cfg.InfluxDB, nil
}
