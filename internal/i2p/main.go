package i2p

import (
	"fmt"
	checksam "github.com/eyedeekay/checki2cp/samcheck"
	"github.com/eyedeekay/i2pkeys"
	sam "github.com/eyedeekay/sam3/helper"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/log"
	"go.uber.org/zap"
	"net"
	"os"
	"strings"
	"time"
)

func save(c *config.Cfg) error {
	fmt.Println("Save I2P Config:")
	fmt.Printf("HTTP_HOST_AND_PORT: %s\n", c.I2PCfg.HttpHostAndPort)
	fmt.Printf("HOST: %s\n", c.I2PCfg.Host)
	fmt.Printf("URL: %s\n", c.I2PCfg.Url)
	err := data.SaveI2pConfig(c.I2PCfg)
	if err != nil {
		return err
	}
	return nil
}

func InitI2P() (net.Listener, error) {
	if config.Conf.I2PCfg.Enabled {
		log.Logger.Info("Starting I2P Mode")
		if err := os.Setenv("NO_PROXY", "127.0.0.1:7672"); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 3)

		for !checksam.CheckSAMAvailable(config.Conf.I2PCfg.SAMAddr) {
			log.Logger.Info("Checking SAM")
			time.Sleep(time.Second * 15)
			log.Logger.Info("Waiting for SAM")
		}
		log.Logger.Info("SAM is available")

		if status, faddr, err := portCheck(config.Conf.I2PCfg.HttpHostAndPort); err == nil {
			if status {
				log.Logger.Fatal(faddr, zap.Error(err))
				return nil, err
			}
		} else {
			log.Logger.Fatal(err.Error())
		}
		log.Logger.Info("Starting I2P server")

		_, listener, err := waitPass("")
		if err != nil {
			panic(err)
		}
		return listener, nil
	}
	return nil, nil
}

func portCheck(addr string) (status bool, faddr string, err error) {
	host, port, err := net.SplitHostPort(addr)

	config.Conf.I2PCfg.Host = host
	if err != nil {
		log.Logger.Fatal("Invalid address")
	}
	if host == "" {
		host = "127.0.0.1"
	}

	conn, err := net.DialTimeout("tcp", config.Conf.I2PCfg.HttpHostAndPort, time.Second*5)
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			err = nil
		}
		log.Logger.Error("Connecting error:", zap.Error(err))
	}
	if conn != nil {
		defer conn.Close()
		status = true
		faddr = net.JoinHostPort(host, port)
		log.Logger.Info(fmt.Sprintf("Opened %s", net.JoinHostPort(host, port)))
	}
	return
}
func waitPass(afterName string) (bool, net.Listener, error) {
	listener, err := sam.I2PListener(config.Conf.AppName+afterName, "127.0.0.1:7656", config.Conf.GetI2pPath(afterName))
	if err != nil {
		panic(err)
	}

	config.Conf.I2PCfg.Host = strings.Split(listener.Addr().(i2pkeys.I2PAddr).Base32(), ":")[0]
	if !strings.HasSuffix(config.Conf.I2PCfg.HttpsUrl, "i2p") {
		config.Conf.I2PCfg.HttpsUrl = "https://" + listener.Addr().(i2pkeys.I2PAddr).Base32()
	}
	if !strings.HasSuffix(config.Conf.I2PCfg.Url, "i2p") {
		config.Conf.I2PCfg.Url = "https://" + listener.Addr().(i2pkeys.I2PAddr).Base32()
	}
	config.Conf.I2PCfg.Url = "http://" + listener.Addr().(i2pkeys.I2PAddr).Base32()
	if err := save(config.Conf); err != nil {
		log.Logger.Error(err.Error())
	}
	return true, listener, err
}
