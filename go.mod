module github.com/blocktree/algorand-adapter

go 1.12

require (
	github.com/algorand/go-algorand-sdk v0.0.0-20190615134606-39000aad5c8e
	github.com/algorand/go-codec/codec v1.1.5-pre // indirect
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.11.1
	github.com/blocktree/go-owcdrivers v1.0.14 // indirect
	github.com/blocktree/go-owcrypt v1.0.1
	github.com/blocktree/openwallet v1.4.5
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/nullstyle/go-xdr v0.0.0-20180726165426-f4c839f75077 // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/stellar/go v0.0.0-20190619212043-88133e09f280
	github.com/stellar/go-xdr v0.0.0-20180917104419-0bc96f33a18e // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	golang.org/x/net v0.0.0-20190110200230-915654e7eabc // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2 // indirect
)

replace github.com/algorand/go-codec/codec v1.1.5-pre => github.com/blocktree/go-codec/codec v1.1.5
