package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	keplerkey "github.com/QOSGroup/kepler/server/handler/key"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/config"
	kepler "github.com/wangfeiping/aimrocks/kepler/client"
	"github.com/wangfeiping/aimrocks/kepler/client/key"
	"github.com/wangfeiping/aimrocks/kepler/client/qcp"
	"github.com/wangfeiping/aimrocks/log"
)

var chainNodeInit = func() (context.CancelFunc, error) {
	cdc := star.MakeCodec()
	log.Infof("chain node init... kepler:\t%s", viper.GetString("kepler"))

	// QOS testnet aquarius-1000
	// init QCP
	// http://docs.qoschain.info/qos/command/qoscli.html
	// >> qoscli tx init-qcp

	client := newKeplerClient()

	// key generate
	_, pubKey, err := genKey(client, cdc)
	if err != nil {
		log.Errorf("GET /key/gen calling error: %v", err)
		return nil, nil
	}

	// apply QCP certificate
	var applyID int64
	applyID, err = applyCert(client, pubKey, cdc)
	if err != nil {
		log.Errorf("POST /qcp/apply calling error: %v", err)
		return nil, nil
	}
	// log.Debugf("apply id: %d", applyID)

	// issue certificate
	p := qcp.NewPutQcpApplyParams()
	p.SetID(applyID)
	p.SetStatus(1)
	resp, err := client.Qcp.PutQcpApply(p)
	if err != nil {
		log.Errorf("PUT /qcp/apply calling error: %v", err)
		return nil, nil
	}
	log.Debugf("issue apply: %v", resp)

	// get QCP certificate

	return nil, nil
}

func applyCert(client *kepler.Kepler, pubKey string,
	cdc *wire.Codec) (id int64, err error) {
	p := qcp.NewPostQcpApplyParams()
	p.SetPhone(viper.GetString("phone"))
	p.SetEmail(viper.GetString("email"))
	p.SetInfo(viper.GetString("info"))
	p.SetQosChainID(viper.GetString("qos_chain_id"))
	p.SetQcpChainID(viper.GetString("qsc_chain_id"))
	p.SetQcpPub(pubKey)
	var resp *qcp.PostQcpApplyOK
	resp, err = client.Qcp.PostQcpApply(p)
	if err != nil {
		log.Errorf("apply cert error: %v", err)
		return
	}
	if apply, ok := resp.Payload.Data.(map[string]interface{}); ok {
		v := apply["id"]
		// log.Debugf("apply id: %v, %T", v, v)
		var idStr json.Number
		if idStr, ok = v.(json.Number); ok {
			id, err = idStr.Int64()
			return
		}
	}
	err = errors.New("Can not parse apply response")
	return
}

func genKey(client *kepler.Kepler,
	cdc *wire.Codec) (privKey, pubKey string, err error) {
	var resp *key.GetKeyGenOK
	resp, err = client.Key.GetKeyGen(nil)
	if err != nil {
		return
	}
	// log.Debugf("GET /key/gen response: %d, %v",
	// keyGen.Payload.Code, keyGen.Payload.Data)
	priv := &keplerkey.KeyValue{}
	pub := &keplerkey.KeyValue{}
	if keys, ok := resp.Payload.Data.(map[string]interface{}); ok {
		privKey, err = parseKey(priv, keys["privKey"], cdc)
		if err != nil {
			return
		}
		pubKey, err = parseKey(pub, keys["pubKey"], cdc)
		if err != nil {
			return
		}
	}
	log.Debugf("public key: %v", pubKey)

	home := viper.GetString(commands.FlagHome)
	keyFile := config.GetKeyFilePath(home,
		fmt.Sprintf("%s.pri", viper.GetString("qsc_chain_id")))
	cmn.MustWriteFile(keyFile, []byte(privKey), 0644)
	keyFile = config.GetKeyFilePath(home,
		fmt.Sprintf("%s.pub", viper.GetString("qsc_chain_id")))
	cmn.MustWriteFile(keyFile, []byte(pubKey), 0644)
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

func newKeplerClient() *kepler.Kepler {
	keplerAPI := viper.GetString("kepler")
	u, err := url.Parse(keplerAPI)
	if err != nil {
		panic(err)
	}
	cfg := &kepler.TransportConfig{
		Host:     u.Host,
		BasePath: u.Path,
		Schemes:  []string{u.Scheme},
	}
	log.Debugf("parse kepler api: %v", cfg)
	return kepler.NewHTTPClientWithConfig(nil, cfg)
}
