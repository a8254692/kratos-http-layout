package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
)

//func init() {
//	flag.Parse()
//}

// LoadConfig ...
func LoadConfig(flagPath *string) error {
	// 根据Flag的值设置Viper的配置文件路径
	if *flagPath != "" {
		viper.SetConfigFile(*flagPath)
	} else {
		// 如果没有提供Flag，则使用默认的配置文件路径
		viper.SetConfigName("config")    // name of config file (without extension)
		viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("./configs") // path to look for the config file in
	}

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return nil
}
