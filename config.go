package tjmsync

import (
	"github.com/BurntSushi/toml"
	"log"
)

func ParseConfig(data []byte) (*Config, error) {
	conf := &Config{}
	_, err := toml.Decode(string(data), conf)
	if err != nil {
		log.Println("Conf Parse Decode Error: ", err)
		return nil, err
	}
	return conf, nil
}

type Config struct {
	Jobs         []Job              `toml:"job"`
	Log          ConfigLog          `toml:"job_log"`
	StatusServer ConfigStatusServer `toml:"status_server"`
}

type ConfigLog struct {
	Type          CONFIG_LOG_TYPE `toml:"type"`
	RedisAddr     string          `toml:"redis_addr"`
	RedisPassword string          `toml:"redis_password"`
	RedisDB       int             `toml:"redis_db"`
	RedisKey      string          `toml:"redis_key"`
}

type CONFIG_LOG_TYPE string

const (
	CONFIG_LOG_TYPE_REDIS = "redis"
)

type ConfigStatusServer struct {
	Listen string `toml:"listen"`
	Port   string `toml:"port"`
}

func (c *ConfigStatusServer) GetAddr() string {
	return c.Listen + ":" + c.Port
}
