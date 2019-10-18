package app

import (
	"bytes"
	"encoding/json"

	"github.com/QOSGroup/qbase/baseabci"
	"github.com/tendermint/go-amino"
)

// Codec for amino codec to marshal/unmarshal
type Codec = amino.Codec

// MakeCodec returns a well-setting cdc via the funcation for both client and server
func MakeCodec() *Codec {
	cdc := baseabci.MakeQBaseCodec()
	return cdc
}

// MarshalJSONIndent attempt to make some pretty json
func MarshalJSONIndent(cdc *Codec, obj interface{}) ([]byte, error) {
	bz, err := cdc.MarshalJSON(obj)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = json.Indent(&out, bz, "", "  ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
