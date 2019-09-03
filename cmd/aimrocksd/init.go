package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	keplerkey "github.com/QOSGroup/kepler/server/handler/key"
	keplermodule "github.com/QOSGroup/kepler/server/module"
	"github.com/QOSGroup/qstars/star"
	"github.com/QOSGroup/qstars/wire"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/wangfeiping/aimrocks/commands"
	"github.com/wangfeiping/aimrocks/config"
	kepler "github.com/wangfeiping/aimrocks/kepler/client"
	"github.com/wangfeiping/aimrocks/kepler/client/key"
	"github.com/wangfeiping/aimrocks/kepler/client/qcp"
	"github.com/wangfeiping/aimrocks/kepler/models"
	"github.com/wangfeiping/aimrocks/log"
)

const (
	defaultQstarsConfFile = "qstarsconf.toml"
	defaultQstarsTemplate = `# this is QCP privatekey
QStarsPrivateKey = "%s"
QOSChainName = "%s"
community = ""

`
)

var chainNodeInit = func() (context.CancelFunc, error) {
	cdc := star.MakeCodec()
	log.Infof("chain node init... kepler:\t%s", viper.GetString("kepler"))

	// QOS testnet aquarius-1000
	// init QCP
	// http://docs.qoschain.info/qos/command/qoscli.html
	// >> qoscli tx init-qcp

	client := newKeplerClient()
	qcpChainID := viper.GetString("qsc_chain_id")
	qosChainID := viper.GetString("qos_chain_id")

	// key generate
	privKey, pubKey, err :=
		genKey(client, cdc, qcpChainID)
	if err != nil {
		log.Errorf("generate qcp keys error: %v", err)
		return nil, nil
	}

	// cassini key generate
	_, _, err =
		genKey(client, cdc, "cassini")
	if err != nil {
		log.Errorf("generate cassini keys error: %v", err)
		return nil, nil
	}

	// create private key config file
	toml := fmt.Sprintf(defaultQstarsTemplate,
		privKey.Value,
		qosChainID)
	home := viper.GetString(commands.FlagHome)
	tomlFile := config.GetConfigFilePath(home,
		defaultQstarsConfFile)
	cmn.MustWriteFile(tomlFile, []byte(toml), 0644)

	// apply QCP certificate
	var applyID int64
	applyID, err = applyCert(client, cdc,
		pubKey, qcpChainID, qosChainID)
	if err != nil {
		log.Errorf("apply cert error: %v", err)
		return nil, nil
	}
	// log.Debugf("apply id: %d", applyID)

	// issue certificate
	err = issueCert(client, applyID)
	if err != nil {
		log.Errorf("issue cert (applyId: %d) error: %v",
			applyID, err)
		return nil, nil
	}

	// get QCP certificate
	err = getQcpCert(client, applyID, qcpChainID)
	if err != nil {
		log.Errorf("get cert (applyId: %d) error: %v",
			applyID, err)
		return nil, nil
	}
	return nil, nil
}

func getQcpCert(client *kepler.Kepler,
	applyID int64, name string) (err error) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println("recover err: ", err)
	// 	}
	// }()
	p := qcp.NewGetQcpCaApplyIDParams()
	p.SetApplyID(applyID)
	var resp *qcp.GetQcpCaApplyIDOK
	resp, err = client.Qcp.GetQcpCaApplyID(p)
	if err != nil {
		return
	}
	// log.Debugf("QCP cert: %v", resp)
	// log.Debugf("QCP cert: %v", resp.Payload.Data)
	ca := &keplermodule.CaQcp{}
	if err = parseResponse(resp.GetPayload(), ca); err != nil {
		return
	}
	var bytesCa []byte
	bytesCa, err = json.Marshal(ca)
	if err != nil {
		return
	}
	home := viper.GetString(commands.FlagHome)
	crtFile := config.GetKeyFilePath(home,
		fmt.Sprintf("%s.crt", name))
	cmn.MustWriteFile(crtFile, bytesCa, 0644)
	log.Infof("get QCP cert ok: %d", applyID)
	return
}

func issueCert(client *kepler.Kepler,
	applyID int64) (err error) {
	p := qcp.NewPutQcpApplyParams()
	p.SetID(applyID)
	p.SetStatus(1)
	var resp *qcp.PutQcpApplyOK
	resp, err = client.Qcp.PutQcpApply(p)
	if err != nil {
		return
	}
	if resp.Payload.Code != 0 {
		err = fmt.Errorf("failed response: %v", resp)
		return
	}
	return
}

func applyCert(client *kepler.Kepler, cdc *wire.Codec,
	pubKey string, qcpChainID, qosChainID string) (id int64, err error) {
	p := qcp.NewPostQcpApplyParams()
	p.SetPhone(viper.GetString("phone"))
	p.SetEmail(viper.GetString("email"))
	p.SetInfo(viper.GetString("info"))
	p.SetQcpChainID(qcpChainID)
	p.SetQosChainID(qosChainID)
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

func genKey(client *kepler.Kepler, cdc *wire.Codec,
	name string) (priv *keplerkey.KeyValue,
	pubKey string, err error) {
	var resp *key.GetKeyGenOK
	resp, err = client.Key.GetKeyGen(nil)
	if err != nil {
		return
	}
	// log.Debugf("GET /key/gen response: %d, %v",
	// keyGen.Payload.Code, keyGen.Payload.Data)
	priv = &keplerkey.KeyValue{}
	pub := &keplerkey.KeyValue{}
	var privKey string
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
		fmt.Sprintf("%s.pri", name))
	cmn.MustWriteFile(keyFile, []byte(privKey), 0644)
	keyFile = config.GetKeyFilePath(home,
		fmt.Sprintf("%s.pub", name))
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

func parseResponse(payload *models.TypesResult, obj interface{}) (err error) {
	if payload == nil {
		return errors.New("payload is nil")
	}
	var data map[string]interface{}
	var ok bool
	if data, ok = payload.Data.(map[string]interface{}); !ok {
		err = fmt.Errorf("can not parse response data: %v", payload)
		return
	}
	dc := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   obj,
		DecodeHook: mapstructure.StringToTimeHookFunc(
			"2006-01-02T15:04:05-07:00")}
	var decoder *mapstructure.Decoder
	decoder, err = mapstructure.NewDecoder(dc)
	if err != nil {
		return
	}
	err = decoder.Decode(data)
	return
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
