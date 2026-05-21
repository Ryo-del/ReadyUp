package config

import "time"

type Config struct {
	DB  DBConfig  `mapstructure:"DB"`
	JWT JWTConfig `mapstructure:"jwt"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

func (c JWTConfig) ExpireDuration() time.Duration {
	if c.ExpireHours <= 0 {
		return 24 * time.Hour
	}
	return time.Duration(c.ExpireHours) * time.Hour
}
