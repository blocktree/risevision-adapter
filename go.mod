module github.com/blocktree/risevision-adapter

go 1.12

require (
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.11.1
	github.com/blocktree/go-owcdrivers v1.0.12
	github.com/blocktree/go-owcrypt v1.0.1
	github.com/blocktree/lisk-adapter v1.1.5
	github.com/blocktree/openwallet v1.4.3
	github.com/eoscanada/eos-go v0.8.10
	github.com/imroc/req v0.2.3
	github.com/kevinburke/nacl v0.0.0-20190102041437-38063847a0c9
	github.com/liskascend/lisk-go v0.0.0-20180425000847-f432f55b9ba8
	github.com/pkg/errors v0.8.1
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/tidwall/gjson v1.2.1
	github.com/tidwall/sjson v1.0.4 // indirect
	go.uber.org/atomic v1.3.2 // indirect
	go.uber.org/multierr v1.1.0 // indirect
	go.uber.org/zap v1.9.1 // indirect
	golang.org/x/crypto v0.0.0-20190513172903-22d7a77e9e5f
	gopkg.in/resty.v1 v1.12.0
)

//replace github.com/blocktree/openwallet => ../../openwallet
