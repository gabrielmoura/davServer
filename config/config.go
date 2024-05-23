package config

import (
	"errors"
	"flag"
	"github.com/gabrielmoura/davServer/internal/log"
	"github.com/spf13/viper"
	"path/filepath"
	"regexp"
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
	SAMAddr         string `mapstructure:"SAM_ADDR"`
	KeyPath         string `mapstructure:"KEY_PATH"`
}

var (
	enabledConfig = flag.Bool("config", false, "Enable config file")
	enabledI2P    = flag.Bool("i2p", false, "Enable I2P")
	rootDirectory = flag.String("root", "./root", "Diretório raiz do servidor WebDAV")
	globalToken   = flag.String("token", "123456", "Token de autenticação")
	port          = flag.Int("port", 8080, "Server Port")
	ExportUsers   = flag.Bool("export", false, "Export Users")
	ImportUsers   = flag.Bool("import", false, "Import Users")

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
		DBDir:        "/tmp/DavServer",
		I2PCfg: &I2PCfg{
			Enabled:         *enabledI2P,
			HttpHostAndPort: "127.0.0.1:7672",
			Host:            "",
			Url:             "127.0.0.1:7672",
			HttpsUrl:        "",
			SAMAddr:         "127.0.0.1:7656",
			KeyPath:         "./",
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
	vip.SetDefault("DB_DIR", "/tmp/DavServer")
	vip.SetDefault("APP_NAME", "DavServer")
	vip.SetDefault("TIME_FORMAT", "02-Jan-2006")
	vip.SetDefault("TIME_ZONE", "America/Sao_Paulo")
	vip.SetDefault("I2P_CFG.ENABLED", false)
	vip.SetDefault("I2P_CFG.SAM_ADDR", "127.0.0.1:7656")
	vip.SetDefault("I2P_CFG.HTTP_HOST_AND_PORT", "127.0.0.1:7672")
	vip.SetDefault("I2P_CFG.KEY_PATH", "./")

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
	// Se APP_NAME não estiver definido no padrão, retorne um erro
	if vip.IsSet("APP_NAME") {
		regex := regexp.MustCompile("^[A-Za-z0-9]+$")
		if !regex.MatchString(vip.GetString("APP_NAME")) {
			return errors.New("APP_NAME só pode conter letras e números")
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
		log.Logger.Info("Carregando configurações do arquivo")
		return loadByConfigFile()
	} else {
		log.Logger.Info("Carregando configurações por flag")
		return loadByFlag()
	}
}

func (c *Cfg) GetI2pPath(afterName string) string {
	return filepath.Join(c.I2PCfg.KeyPath, c.AppName+afterName)
}
