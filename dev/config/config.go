package config

import (
	"time"

	"github.com/spf13/viper"
)

const ConfigName = "config"
const ConfigType = "yaml"

var Configuration Config
var WorkDir string

type Config struct {
	Server struct {
		Mode            string        `mapstructure:"mode"`
		Port            int           `mapstructure:"port"`
		ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
		SheetName       string        `mapstructure:"sheet_name"`
		UploadDir       string        `mapstructure:"upload_dir"`
		Endpoint        struct {
			ExcelUpload string `mapstructure:"excel_upload"`
		} `mapstructure:"endpoint"`
	} `mapstructure:"server"`

	Artemis struct {
		Host          string `mapstructure:"host"`
		Port          int    `mapstructure:"port"`
		Username      string `mapstructure:"username"`
		Password      string `mapstructure:"password"`
		Address       string `mapstructure:"address"`
		ReportAddress string `mapstructure:"report_address"`
	} `mapstructure:"artemis"`

	Mysql struct {
		Host     string   `mapstructure:"host"`
		Port     int      `mapstructure:"port"`
		Database string   `mapstructure:"database"`
		Username string   `mapstructure:"username"`
		Password string   `mapstructure:"password"`
		Options  []string `mapstructure:"options"`
	} `mapstructure:"mysql"`

	Logger struct {
		Dir        string `mapstructure:"dir"`
		FileName   string `mapstructure:"file_name"`
		MaxBackups int    `mapstructure:"max_backups"`
		MaxSize    int    `mapstructure:"max_size"`
		MaxAge     int    `mapstructure:"max_age"`
		Compress   bool   `mapstructure:"compress"`
		LocalTime  bool   `mapstructure:"local_time"`
	} `mapstructure:"logger"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	var config Config
	err = viper.Unmarshal(&config)
	Configuration = config
	return
}
