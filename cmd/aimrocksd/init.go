package main

import (
	"context"
	"fmt"
	"net/url"

	keplerkey "github.com/QOSGroup/kepler/server/handler/key"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/viper"
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
	privKey, pubKey, err := genKey(client, cdc)
	if err != nil {
		log.Errorf("GET /key/gen calling error: %v", err)
		return nil, nil
	}
	log.Debugf("private key: %s", privKey)
	log.Debugf("public key: %v", pubKey)
	// pubKey := "{\"type\":\"tendermint/PubKeyEd25519\",\"value\":\"tnQTl6SjqSIpT5DZDcq8L+penU4y08ddRxLyHUp9/Ws=\"}"

	// apply QCP certificate
	err = applyCert(client, pubKey, cdc)
	if err != nil {
		log.Errorf("POST /qcp/apply calling error: %v", err)
		return nil, nil
	}

	// issue certificate

	// get QCP certificate

	return nil, nil
}

func applyCert(client *kepler.Kepler, pubKey string,
	cdc *wire.Codec) (err error) {
	p := qcp.NewPostQcpApplyParams()
	p.SetPhone(viper.GetString("phone"))
	p.SetEmail(viper.GetString("email"))
	p.SetInfo(viper.GetString("info"))
	p.SetQosChainID(viper.GetString("qos_chain_id"))
	p.SetQcpChainID(viper.GetString("qsc_chain_id"))
	p.SetQcpPub(pubKey)
	var resp *qcp.PostQcpApplyOK
	resp, err = client.Qcp.PostQcpApply(p)
	if err == nil {
		log.Debug("QCP apply: ", resp)
		log.Debug("QCP apply: ", resp.Payload.Data)
	}
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
