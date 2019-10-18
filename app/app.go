package app

import (
	"github.com/QOSGroup/qbase/server"
	"github.com/tendermint/tendermint/crypto"
)

// AimContext wraps server.Context
type AimContext struct {
	ServerContext      *server.Context
	QStarsSignerPriv   crypto.PrivKey
	QStarsTransactions []string
}

var ctx *AimContext

// GetServerContext returns app context
func GetServerContext() *AimContext {
	return ctx
}

// InitApp init app
func InitApp() *AimContext {
	ctx = &AimContext{
		ServerContext: server.NewDefaultContext(),
	}
	return ctx
}
