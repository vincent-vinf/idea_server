package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"idea_server/config"
)

var (
	IDEA_VP  *viper.Viper
	IDEA_DB     *gorm.DB
	IDEA_REDIS  *redis.Client
	IDEA_CONFIG config.Server
	IDEA_LOG    *zap.Logger
)
