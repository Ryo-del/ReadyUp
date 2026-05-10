package config

type Config struct {
	DB DBConfig `mapstructure:"DB"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}
