package i2p

import (
	"fmt"
	checksam "github.com/eyedeekay/checki2cp/samcheck"
	"github.com/eyedeekay/i2pkeys"
	sam "github.com/eyedeekay/sam3/helper"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func save(c *config.Cfg) error {
	log.Printf("Save I2P Config: %v", c.I2PCfg)
	err := data.SaveI2pConfig(c.I2PCfg)
	if err != nil {
		return err
	}
	return nil
}

func InitI2P() (net.Listener, error) {
	if config.Conf.I2PCfg.Enabled {
		log.Println("Starting I2P Mode")
		if err := os.Setenv("NO_PROXY", "127.0.0.1:7672"); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 3)

		for !checksam.CheckSAMAvailable("127.0.0.1:7656") {
			log.Println("Checking SAM")
			time.Sleep(time.Second * 15)
			log.Println("Waiting for SAM")
		}
		log.Println("SAM is available")

		if status, faddr, err := portCheck(config.Conf.I2PCfg.HttpHostAndPort); err == nil {
			if status {
				log.Fatal(err, faddr)
				return nil, err
			}
		} else {
			log.Fatal(err)
		}
		log.Println("Starting I2P server")

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
		log.Fatal("Invalid address")
	}
	if host == "" {
		host = "127.0.0.1"
	}

	conn, err := net.DialTimeout("tcp", config.Conf.I2PCfg.HttpHostAndPort, time.Second*5)
	if err != nil {
		if strings.Contains(err.Error(), "refused") {
			err = nil
		}
		log.Println("Connecting error:", err)
	}
	if conn != nil {
		defer conn.Close()
		status = true
		faddr = net.JoinHostPort(host, port)
		log.Println("Opened", net.JoinHostPort(host, port))
	}
	return
}
func waitPass(afterName string) (bool, net.Listener, error) {

	fmt.Println("User exists, ready to go.")
	listener, err := sam.I2PListener(config.Conf.AppName+afterName, "127.0.0.1:7656", config.Conf.AppName+afterName)
	if err != nil {
		panic(err)
	}

	config.Conf.I2PCfg.Host = strings.Split(listener.Addr().(i2pkeys.I2PAddr).Base32(), ":")[0]
	if !strings.HasSuffix(config.Conf.I2PCfg.HttpsUrl, "i2p") {
		config.Conf.I2PCfg.HttpsUrl = "https://" + listener.Addr().(i2pkeys.I2PAddr).Base32()
		//log.Println(domainhelp)
	}
	if !strings.HasSuffix(config.Conf.I2PCfg.Url, "i2p") {
		config.Conf.I2PCfg.Url = "https://" + listener.Addr().(i2pkeys.I2PAddr).Base32()
		//log.Println(domainhelp)
	}
	config.Conf.I2PCfg.Url = "http://" + listener.Addr().(i2pkeys.I2PAddr).Base32()
	if err := save(config.Conf); err != nil {
		log.Println(err)
	}
	return true, listener, err
}
