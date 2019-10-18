module github.com/wangfeiping/aimrocks

go 1.12

require (
	github.com/QOSGroup/kepler v0.6.1-0.20190826125648-e0551d2f4e77
	github.com/QOSGroup/qbase v0.2.3
	github.com/QOSGroup/qos v0.0.7
	github.com/QOSGroup/qstars v0.5.0
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/cihub/seelog v0.0.0-20170130134532-f561c5e57575
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.4
	github.com/go-openapi/strfmt v0.19.2
	github.com/go-openapi/swag v0.19.5
	github.com/mitchellh/mapstructure v1.1.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.0
)

replace github.com/QOSGroup/qstars v0.5.0 => github.com/wangfeiping/qstars v0.4.2-0.20190909103442-d241979cce36
