package internal

import "github.com/spf13/viper"

type Config struct {
	CfAccountID             string `mapstructure:"CF_ACCOUNT_ID"`
	CfImagesStreamReadWrite string `mapstructure:"CF_IMAGES_STREAM_READ_WRITE"`
	CfImagesStreamReadOnly  string `mapstructure:"CF_IMAGES_STREAM_READ_ONLY"`
}

func LoadConfig() (config Config, err error) {
	viper.SetConfigName("ems")
	viper.SetConfigType("env")
	viper.AddConfigPath("$HOME/.config/ems")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
