module github.com/blocktree/algorand-adapter

go 1.12

require (
	github.com/algorand/go-algorand-sdk v0.0.0-20190615134606-39000aad5c8e
	github.com/algorand/go-codec/codec v1.1.5-pre // indirect
	github.com/asdine/storm v2.1.2+incompatible
	github.com/astaxie/beego v1.12.2
	github.com/blocktree/go-owcdrivers v1.2.23 // indirect
	github.com/blocktree/go-owcrypt v1.1.9
	github.com/blocktree/openwallet v1.5.4
	github.com/nullstyle/go-xdr v0.0.0-20180726165426-f4c839f75077 // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/stellar/go v0.0.0-20190619212043-88133e09f280
	github.com/stellar/go-xdr v0.0.0-20180917104419-0bc96f33a18e // indirect
)

replace github.com/algorand/go-codec/codec v1.1.5-pre => github.com/blocktree/go-codec/codec v1.1.5

replace golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190621222207-cc06ce4a13d4
