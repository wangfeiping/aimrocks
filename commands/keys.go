package commands

import (
	"github.com/QOSGroup/qbase/client/keys"
	"github.com/spf13/cobra"
	go_amino "github.com/tendermint/go-amino"
)

// NewKeysCommand returns key command
func NewKeysCommand(cdc *go_amino.Codec) *cobra.Command {
	return keys.KeysCommand(cdc)
}
