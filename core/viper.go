package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"idea_server/global"
	"idea_server/utils/constant"
)

func Viper() *viper.Viper {
	v := viper.New()
	v.SetConfigFile(constant.ConfigFile)
	v.SetConfigType("ini")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&global.IDEA_CONFIG); err != nil {
			fmt.Println(err)
		}
	})
	if err := v.Unmarshal(&global.IDEA_CONFIG); err != nil {
		fmt.Println(err)
	}
	return v
}
