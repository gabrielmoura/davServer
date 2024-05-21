package config

import (
	"errors"
	"flag"
	"github.com/spf13/viper"
	"log"
)

type Cfg struct {
	AppName      string  `mapstructure:"APP_NAME"`
	DBDir        string  `mapstructure:"DB_DIR"`
	Port         int     `mapstructure:"PORT"`
	ShareRootDir string  `mapstructure:"SHARE_ROOT_DIR"`
	TimeFormat   string  `mapstructure:"TIME_FORMAT"`
	TimeZone     string  `mapstructure:"TIME_ZONE"`
	GlobalToken  string  `mapstructure:"GLOBAL_TOKEN"`
	I2PCfg       *I2PCfg `mapstructure:"I2P_CFG"`
}
type I2PCfg struct {
	Enabled         bool   `mapstructure:"ENABLED"`
	HttpHostAndPort string `mapstructure:"HTTP_HOST_AND_PORT"`
	Host            string `mapstructure:"HOST"`
	Url             string `mapstructure:"URL"`
	HttpsUrl        string `mapstructure:"HTTPS_URL"`
}

var (
	enabledConfig = flag.Bool("config", false, "Enable config file")
	enabledI2P    = flag.Bool("i2p", false, "Enable I2P")
	rootDirectory = flag.String("root", "./root", "Diretório raiz do servidor WebDAV")
	globalToken   = flag.String("token", "123456", "Token de autenticação")
	port          = flag.Int("port", 8080, "Server Port")

	Conf *Cfg
)

func loadByFlag() error {
	cfg := &Cfg{
		Port:         *port,
		ShareRootDir: *rootDirectory,
		GlobalToken:  *globalToken,
		AppName:      "DavServer",
		TimeFormat:   "02-Jan-2006",
		TimeZone:     "America/Sao_Paulo",
		DBDir:        "/tmp/badgerDB",
		I2PCfg: &I2PCfg{
			Enabled:         *enabledI2P,
			HttpHostAndPort: "",
			Host:            "",
			Url:             "",
			HttpsUrl:        "",
		},
	}
	// Atualiza a variável global Conf
	Conf = cfg
	return nil
}
func loadByConfigFile() error {
	var cfg Cfg
	vip := viper.New()

	// Definindo valores padrão
	vip.SetDefault("PORT", 8080)
	vip.SetDefault("SHARE_ROOT_DIR", "./root")
	vip.SetDefault("GLOBAL_TOKEN", "123456")
	vip.SetDefault("DB_DIR", "/tmp/badgerDB")
	vip.SetDefault("APP_NAME", "DavServer")
	vip.SetDefault("TIME_FORMAT", "02-Jan-2006")
	vip.SetDefault("TIME_ZONE", "America/Sao_Paulo")
	vip.SetDefault("I2P_CFG.ENABLED", false)

	// Lendo o arquivo de configuração conf.yml
	vip.SetConfigName("conf")
	vip.SetConfigType("yml")
	vip.AddConfigPath(".")
	vip.AddConfigPath("/opt/davSrv")
	vip.AddConfigPath("/etc/davSrv")

	// Lendo as configurações do arquivo conf.yml
	if err := vip.ReadInConfig(); err != nil {
		// Se o arquivo conf.yml não for encontrado, continue sem erro
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// Atribua as configurações ao cfg
	if err := vip.Unmarshal(&cfg); err != nil {
		return err
	}

	// Atualiza a variável global Conf
	Conf = &cfg

	return nil
}
func LoadConfig() error {
	flag.Parse()
	if *enabledConfig {
		log.Printf("Carregando configurações do arquivo")
		return loadByConfigFile()
	} else {
		log.Printf("Carregando configurações por flag")
		return loadByFlag()
	}
}
