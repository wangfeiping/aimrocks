package main

import (
	"context"
	"fmt"

	keplerkey "github.com/QOSGroup/kepler/server/handler/key"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/viper"
	kepler "github.com/wangfeiping/aimrocks/kepler/client"
	"github.com/wangfeiping/aimrocks/kepler/client/key"
	"github.com/wangfeiping/aimrocks/log"
)

var chainNodeInit = func() (context.CancelFunc, error) {
	cdc := star.MakeCodec()
	log.Infof("chain node init... kepler:\t%s", viper.GetString("kepler"))

	// key generate
	privKey, pubKey, err := genKey(cdc)
	if err != nil {
		log.Errorf("GET /key/gen calling error: %v", err)
		return nil, nil
	}
	log.Debugf("private key: %s", privKey)
	log.Debugf("public key: %v", pubKey)
	return nil, nil
}

func genKey(cdc *wire.Codec) (privKey, pubKey string, err error) {
	var keyGen *key.GetKeyGenOK
	keyGen, err = kepler.Default.Key.GetKeyGen(key.NewGetKeyGenParams())
	if err != nil {
		return
	}
	// log.Debugf("GET /key/gen response: %d, %v",
	// keyGen.Payload.Code, keyGen.Payload.Data)
	priv := &keplerkey.KeyValue{}
	pub := &keplerkey.KeyValue{}
	if keys, ok := keyGen.Payload.Data.(map[string]interface{}); ok {
		privKey, err = parseKey(priv, keys["privKey"], cdc)
		if err != nil {
			return
		}
		pubKey, err = parseKey(pub, keys["pubKey"], cdc)
		if err != nil {
			return
		}
	}
	return
}

func parseKey(key *keplerkey.KeyValue,
	data interface{}, cdc *wire.Codec) (string, error) {
	var ok bool
	var model map[string]interface{}
	if model, ok = data.(map[string]interface{}); !ok {
		return "", fmt.Errorf("Can not parse data to KeyValue: %v", data)
	}
	t := model["type"]
	if key.Type, ok = t.(string); !ok {
		return "", fmt.Errorf("Can not convert to string: %v", t)
	}
	v := model["value"]
	if key.Value, ok = v.(string); !ok {
		return "", fmt.Errorf("Can not convert to string: %v", v)
	}
	var bytes []byte
	var err error
	if bytes, err = cdc.MarshalJSON(key); err != nil {
		return "", err
	}
	return string(bytes), nil
}
