package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LongTermDatabase LongTermDatabase `yaml:"LongTermDatabase"`
	GeneralDatabase  GeneralDatabase  `yaml:"GeneralDatabase"`
	ShardingDatabase ShardingDatabase `yaml:"ShardingDatabase"`
	Microservice     Microservice     `yaml:"Microservice"`
	Redis            Redis            `yaml:"Redis"`
	Scheduler        Scheduler        `yaml:"Scheduler"`
	UploadFolderPath string
}

func Load(path string) (config Config, err error) {
	fileName, err := filepath.Abs(path)
	if err != nil {
		return config, err
	}

	fs, err := os.ReadFile(fileName)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(fs, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
