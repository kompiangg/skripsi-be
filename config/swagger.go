package config

type Swagger struct {
	Title       string   `yaml:"Title" validate:"required,notblank"`
	Hostname    string   `yaml:"Hostname" validate:"required,notblank"`
	Port        string   `yaml:"Port" validate:"required,notblank,numeric"`
	Version     string   `yaml:"Version" validate:"required,notblank"`
	Description string   `yaml:"Description" validate:"required,notblank"`
	Schemes     []string `yaml:"Schemes" validate:"required,notblank,dive,notblank"`
}
