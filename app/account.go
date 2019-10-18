package app

type ResultCreateAccount struct {
	PubKey   string `json:"pubKey"`
	PrivKey  string `json:"privKey"`
	Addr     string `json:"addr"`
	Mnemonic string `json:"mnemonic"`
	Type     string `json:"type"`
}

// MockAccount returns a mock account
func MockAccount() *ResultCreateAccount {
	result := &ResultCreateAccount{}
	result.PubKey = "pubkeyAminoStr"
	result.PrivKey = "privkeyAminoStr"
	result.Addr = "bech32Addr"
	result.Mnemonic = "mnemonic"
	result.Type = "local"

	return result
}
