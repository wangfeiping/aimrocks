package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	keplerkey "github.com/QOSGroup/kepler/server/handler/key"
	keplermodule "github.com/QOSGroup/kepler/server/module"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/baseapp"
	"github.com/QOSGroup/qstars/slim"
	"github.com/QOSGroup/qstars/star"
	sdk "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/wire"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	go_amino "github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	cmn "github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
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
	baseapp.InitApp()
	ctx := baseapp.GetServerContext().ServerContext
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
		log.Errorf("generate keys error: %v", err)
		return nil, nil
	}
	var privJSON string
	var pubJSON string
	privJSON, err = marshalKey(cdc, privKey)
	if err != nil {
		log.Errorf("marshal keys error: %v", err)
		return nil, nil
	}
	pubJSON, err = marshalKey(cdc, pubKey)
	if err != nil {
		log.Errorf("marshal keys error: %v", err)
		return nil, nil
	}
	saveKey(qcpChainID, privJSON, pubJSON)

	// cassini key generate
	cassiniPrivKey, cassiniPubKey, err :=
		genKey(client, cdc, "cassini")
	if err != nil {
		log.Errorf("generate keys error: %v", err)
		return nil, nil
	}
	cassiniPrivJSON, err := marshalKey(cdc, cassiniPrivKey)
	if err != nil {
		log.Errorf("marshal keys error: %v", err)
		return nil, nil
	}
	cassiniPubJSON, err := marshalKey(cdc, cassiniPubKey)
	if err != nil {
		log.Errorf("marshal keys error: %v", err)
		return nil, nil
	}
	saveKey("cassini", cassiniPrivJSON, cassiniPubJSON)

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
		pubJSON, qcpChainID, qosChainID)
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

	// gen genesis.json
	// server.InitCmd(ctx, cdc, genBaseCoindGenesisDoc, rootDir)
	err = initGenesisJSON(ctx, cdc,
		qosChainID, cassiniPubKey, createGenesis)
	if err != nil {
		log.Errorf("init genesis (chainId: %s) error: %v", qcpChainID, err)
	}
	return nil, nil
}

func initGenesisJSON(ctx *server.Context, cdc *wire.Codec,
	chainID string, cassiniPubKey *keplerkey.KeyValue,
	genGenesis funcCreateGenesisDoc) error {
	config := ctx.Config
	config.SetRoot(viper.GetString(cli.HomeFlag))

	// chainID := viper.GetString(flagChainID)
	// if chainID == "" {
	// 	chainID = fmt.Sprintf("test-chain-%v", common.RandStr(6))
	// }

	nodeID, valPubkey, err := server.InitializeNodeValidatorFiles(config)
	if err != nil {
		return err
	}

	config.Moniker = "test" // viper.GetString(flagMoniker)

	genFile := config.GenesisFile()

	// overwrite := viper.GetBool(flagOverwrite)
	overwrite := false
	if !overwrite && common.FileExists(genFile) {
		return fmt.Errorf("genesis.json file already exists: %v", genFile)
	}

	genesisDoc, err := genGenesis(ctx, cdc, chainID, cassiniPubKey, valPubkey)
	if err != nil {
		return err
	}

	if err = server.SaveGenDoc(genFile, genesisDoc); err != nil {
		return err
	}

	toPrint := newPrintInfo(config.Moniker, chainID, nodeID, "", genesisDoc.AppState)

	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	return displayInfo(cdc, toPrint)
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
	// bytesCa, err = json.Marshal(ca.Crt)
	// if err != nil {
	// 	return
	// }
	bytesCa = []byte(ca.Crt)
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
	apply := &keplermodule.ApplyQcp{}
	if err = parseResponse(resp.GetPayload(), apply); err != nil {
		return
	}
	id = apply.Id
	return
}

func genKey(client *kepler.Kepler, cdc *wire.Codec,
	name string) (priv *keplerkey.KeyValue,
	pub *keplerkey.KeyValue, err error) {
	var resp *key.GetKeyGenOK
	resp, err = client.Key.GetKeyGen(nil)
	if err != nil {
		return
	}
	// log.Debugf("GET /key/gen response: %d, %v",
	// keyGen.Payload.Code, keyGen.Payload.Data)
	keys := &keplerkey.KeyData{}
	if err = parseResponse(resp.GetPayload(), keys); err != nil {
		return
	}
	priv = &keys.PrivKey
	pub = &keys.PubKey
	return
}

func marshalKey(cdc *go_amino.Codec,
	keyValue *keplerkey.KeyValue) (keyJSON string, err error) {
	var bytes []byte
	if bytes, err = cdc.MarshalJSON(keyValue); err == nil {
		keyJSON = string(bytes)
	}
	return
}

func saveKey(name string,
	privJSON string, pubJSON string) {
	home := viper.GetString(commands.FlagHome)
	keyFile := config.GetKeyFilePath(home,
		fmt.Sprintf("%s.pri", name))
	cmn.MustWriteFile(keyFile, []byte(privJSON), 0644)
	keyFile = config.GetKeyFilePath(home,
		fmt.Sprintf("%s.pub", name))
	cmn.MustWriteFile(keyFile, []byte(pubJSON), 0644)
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

type printInfo struct {
	Moniker    string          `json:"moniker"`
	ChainID    string          `json:"chain_id"`
	NodeID     string          `json:"node_id"`
	GenTxsDir  string          `json:"gentxs_dir"`
	AppMessage json.RawMessage `json:"app_message"`
}

// nolint: errcheck
func displayInfo(cdc *go_amino.Codec, info printInfo) error {
	out, err := cdc.MarshalJSONIndent(info, "", " ")
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "%s\n", string(out))
	return nil
}

func newPrintInfo(moniker, chainID, nodeID, genTxsDir string,
	appMessage json.RawMessage) printInfo {

	return printInfo{
		Moniker:    moniker,
		ChainID:    chainID,
		NodeID:     nodeID,
		GenTxsDir:  genTxsDir,
		AppMessage: appMessage,
	}
}

type funcCreateGenesisDoc func(ctx *server.Context, cdc *go_amino.Codec,
	chainID string, cassiniPubKey *keplerkey.KeyValue,
	nodeValidatorPubKey crypto.PubKey) (tmtypes.GenesisDoc, error)

func createGenesis(ctx *server.Context, cdc *go_amino.Codec,
	chainID string, cassiniPubKey *keplerkey.KeyValue,
	nodeValidatorPubKey crypto.PubKey) (tmtypes.GenesisDoc, error) {

	validator := tmtypes.GenesisValidator{
		PubKey: nodeValidatorPubKey,
		Power:  10,
	}

	//addr, _, err := types.GenerateCoinKey(cdc, types.DefaultCLIHome)
	//if err != nil {
	//	return tmtypes.GenesisDoc{}, err
	//}

	acc := slim.AccountCreate("")

	output, err := wire.MarshalJSONIndent(cdc, acc)
	if err != nil {
		return tmtypes.GenesisDoc{}, err
	}

	fmt.Println(string(output))
	addr, _ := sdk.AccAddressFromBech32(acc.Addr)

	appState, err := genAppState(cdc, chainID, cassiniPubKey, addr)
	if err != nil {
		return tmtypes.GenesisDoc{}, err
	}

	return tmtypes.GenesisDoc{
		ChainID:    chainID,
		Validators: []tmtypes.GenesisValidator{validator},
		AppState:   appState,
	}, nil

}

func genAppState(cdc *go_amino.Codec, chainID string,
	cassiniPubKey *keplerkey.KeyValue, addr types.Address) (
	appState json.RawMessage, err error) {
	appState = json.RawMessage(fmt.Sprintf(`{
		"qcps":[{
			"name": "qos",
			"chain_id": "%s",
			"pub_key":{
        		"type": "tendermint/PubKeyEd25519",
        		"value": "%s"
			}
		}],
  		"accounts": [{
    		"address": "%s",
    		"coins": [
      			{
        			"coin_name":"qstar",
        			"amount":"100000000"
      			}
			]
  		}]
	}`, chainID, cassiniPubKey.Value, addr))
	return
}
