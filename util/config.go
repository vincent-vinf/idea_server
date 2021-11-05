package util

import (
	"gopkg.in/ini.v1"
	"log"
	"strconv"
	"sync"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Passwd   string
}

type JWTConfig struct {
	SecretKey []byte
}

type MysqlConfig struct {
	Source string
}

type RedisConfig struct {
	Address string
	Passwd  string
	DB      int
}

type MqConfig struct {
	Url string
}

var (
	emailCfg  *EmailConfig
	emailOnce sync.Once

	jwtCfg  *JWTConfig
	jwtOnce sync.Once

	mysqlCfg  *MysqlConfig
	mysqlOnce sync.Once

	redisCfg  *RedisConfig
	redisOnce sync.Once

	mqCfg  *MqConfig
	mqOnce sync.Once

	cfg        *ini.File
	configOnce sync.Once
)

func loadConfig() {
	if cfg == nil {
		configOnce.Do(func() {
			var err error
			cfg, err = ini.Load("config.ini")
			if err != nil {
				panic(err)
			}
		})
	}
}

func LoadEmailCfg() *EmailConfig {
	if emailCfg == nil {
		loadConfig()
		emailOnce.Do(func() {
			section, err := cfg.GetSection("email")
			if err != nil {
				panic(err)
			}
			port, err := strconv.Atoi(section.Key("port").Value())
			if err != nil {
				panic(err)
			}
			emailCfg = &EmailConfig{
				Host:     section.Key("host").Value(),
				Port:     port,
				Username: section.Key("username").Value(),
				Passwd:   section.Key("passwd").Value(),
			}
		})
	}
	return emailCfg
}

func LoadJWTCfg() *JWTConfig {
	if jwtCfg == nil {
		loadConfig()
		jwtOnce.Do(func() {
			section, err := cfg.GetSection("jwt")
			if err != nil {
				panic(err)
			}
			jwtCfg = &JWTConfig{
				SecretKey: []byte(section.Key("key").Value()),
			}
		})
	}
	return jwtCfg
}

func LoadMysqlCfg() *MysqlConfig {
	if mysqlCfg == nil {
		loadConfig()
		mysqlOnce.Do(func() {
			section, err := cfg.GetSection("mysql")
			if err != nil {
				panic(err)
			}
			mysqlCfg = &MysqlConfig{
				Source: section.Key("source").Value(),
			}
		})
	}
	return mysqlCfg
}

func LoadRedisCfg() *RedisConfig {
	if redisCfg == nil {
		loadConfig()
		redisOnce.Do(func() {
			section, err := cfg.GetSection("redis")
			if err != nil {
				panic(err)
			}
			db, err := strconv.Atoi(section.Key("db").Value())
			if err != nil {
				panic(err)
			}
			redisCfg = &RedisConfig{
				Address: section.Key("address").Value(),
				Passwd:  section.Key("passwd").Value(),
				DB:      db,
			}
		})
	}
	return redisCfg
}

func LoadMqCfg() *MqConfig {
	if mqCfg == nil {
		loadConfig()
		mqOnce.Do(func() {
			section, err := cfg.GetSection("mq")
			user, err := section.GetKey("user")
			handleErr(err)
			mqUser := user.Value()

			pass, err := section.GetKey("password")
			handleErr(err)
			mqPass := pass.Value()

			hostKey, err := section.GetKey("host")
			handleErr(err)
			host := hostKey.Value()

			portKey, err := section.GetKey("port")
			handleErr(err)
			port := portKey.Value()

			routeKey, err := section.GetKey("route")
			handleErr(err)
			route := routeKey.Value()
			if err != nil {
				panic(err)
			}
			mqCfg = &MqConfig{
				Url: "amqp://" + mqUser + ":" + mqPass + "@" + host + ":" + port + route,
			}
		})
	}
	return mqCfg
}

func handleErr(err error) {
	if err != nil {
		log.Fatal("load config error: ", err)
	}
}
