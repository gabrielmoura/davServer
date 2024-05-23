package i2p

import (
	"encoding/base64"
	"fmt"
	"github.com/eyedeekay/i2pkeys"
	"github.com/eyedeekay/sam3"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/log"
	"net"
)

func NetListener(name, samaddr string) (net.Listener, error) {
	return I2PListener(name, sam3.SAMDefaultAddr(samaddr))
}

// I2PListener is a convenience function which takes a SAM tunnel name, a SAM address and a filename.
// If the file contains I2P keys, it will create a service using that address. If the file does not
// exist, keys will be generated and stored in that file.
func I2PListener(name, samaddr string) (*sam3.StreamListener, error) {

	log.Logger.Info("Starting and registering I2P service, please wait a couple of minutes...")
	listener, err := I2PStreamSession(name, sam3.SAMDefaultAddr(samaddr))
	if err != nil {
		return nil, err
	}
	//if keyspath != "" {
	//	err = ioutil.WriteFile(keyspath+".i2p.public.txt", []byte(listener.Keys().Addr().Base32()), 0644)
	//	if err != nil {
	//		log.Fatalf("error storing I2P base32 address in adjacent text file, %s", err)
	//	}
	//}

	log.Logger.Info(fmt.Sprintf("Listening on: %s", listener.Addr().Base32()))
	return listener.Listen()
}

// I2PStreamSession is a convenience function which returns a sam3.StreamSession instead
// of a sam3.StreamListener. It also takes care of setting a persisitent key on behalf
// of the user.
func I2PStreamSession(name, samaddr string) (*sam3.StreamSession, error) {

	log.Logger.Info("Starting and registering I2P session...")
	sam, err := sam3.NewSAM(sam3.SAMDefaultAddr(samaddr))
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("error connecting to SAM to %s: %s", sam3.SAMDefaultAddr(samaddr), err))
	}
	keys, err := GenerateOrLoadKeys(sam)
	if err != nil {
		return nil, err
	}
	stream, err := sam.NewStreamSession(name, *keys, sam3.Options_Medium)
	return stream, err
}

// I2PDataGramsession is a convenience function which returns a sam3.DatagramSession.
// It also takes care of setting a persisitent key on behalf of the user.
func I2PDatagramSession(name, samaddr string) (*sam3.DatagramSession, error) {
	log.Logger.Info("Starting and registering I2P session...")
	sam, err := sam3.NewSAM(sam3.SAMDefaultAddr(samaddr))
	if err != nil {

		log.Logger.Fatal(fmt.Sprintf("error connecting to SAM to %s: %s", sam3.SAMDefaultAddr(samaddr), err))
	}
	keys, err := GenerateOrLoadKeys(sam)
	if err != nil {
		return nil, err
	}
	gram, err := sam.NewDatagramSession(name, *keys, sam3.Options_Medium, 0)
	return gram, err
}

// I2PPrimarySession is a convenience function which returns a sam3.PrimarySession.
// It also takes care of setting a persisitent key on behalf of the user.
func I2PPrimarySession(name, samaddr string) (*sam3.PrimarySession, error) {

	log.Logger.Info("Starting and registering I2P session...")
	sam, err := sam3.NewSAM(sam3.SAMDefaultAddr(samaddr))
	if err != nil {
		log.Logger.Fatal(fmt.Sprintf("error connecting to SAM to %s: %s", sam3.SAMDefaultAddr(samaddr), err))
	}
	keys, err := GenerateOrLoadKeys(sam)
	if err != nil {
		return nil, err
	}
	gram, err := sam.NewPrimarySession(name, *keys, sam3.Options_Medium)
	return gram, err
}

// GenerateOrLoadKeys is a convenience function which takes a filename and a SAM session.
// if the SAM session is nil, a new one will be created with the defaults.
// The keyspath must be the path to a place to store I2P keys. The keyspath will be suffixed with
// .i2p.private for the private keys, and public.txt for the b32 addresses.
// If the keyspath.i2p.private file does not exist, keys will be generated and stored in that file.
// if the keyspath.i2p.private does exist, keys will be loaded from that location and returned
func GenerateOrLoadKeys(sam *sam3.SAM) (keys *i2pkeys.I2PKeys, err error) {
	if sam == nil {
		sam, err = sam3.NewSAM(sam3.SAMDefaultAddr("127.0.0.1:7656"))
		if err != nil {
			return nil, err
		}
	}
	i2pKey, err := data.ReadI2pKey()
	if err != nil {
		// Criar Chave
		tkeys, err := sam.NewKeys()
		if err != nil {
			log.Logger.Fatal(fmt.Sprintf("error generating I2P keys: %s", err))
		}
		keys = &tkeys
		err = data.WriteI2pKey(&data.I2PKey{
			Private: keys.Both,
			Public:  base64.StdEncoding.EncodeToString([]byte(keys.Address)),
			Url:     keys.Address.Base32(),
		})
		if err != nil {
			log.Logger.Fatal(fmt.Sprintf("error saving I2P keys: %s", err))
		}

	} else {
		dataAddr, err := base64.StdEncoding.DecodeString(i2pKey.Public)
		if err != nil {
			log.Logger.Fatal(fmt.Sprintf("error decoding I2P keys: %s", err))
		}
		// Carregar Chave
		keys = &i2pkeys.I2PKeys{
			Address: i2pkeys.I2PAddr(dataAddr),
			Both:    i2pKey.Private,
		}
	}
	return keys, nil
}

// GenerateKeys is a shorter version of GenerateOrLoadKeys which generates keys and stores them in a file.
// it always uses a new default SAM session.
func GenerateKeys() (keys *i2pkeys.I2PKeys, err error) {
	return GenerateOrLoadKeys(nil)
}
