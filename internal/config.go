package internal

type AppConfig struct {
	Version   string `yaml:"version"`
	Port      int    `yaml:"port"`
	Env       string `yaml:"env"`
	Debug     bool   `yaml:"debug"`
	LogToFile bool   `yaml:"logToFile"`
}
