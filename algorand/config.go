package algorand

import (
	"path/filepath"
	"strings"

	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/openwallet/common/file"
)

const (
	//币种
	Symbol              = "ALGO"
	CurveType           = owcrypt.ECC_CURVE_ED25519
	ALGO_SEQUENCEID_KEY = "ALGO_SEQUENCEID_KEY"

	//默认配置内容
	defaultConfig = `

# RPC api url
serverAPI = ""
`
)

type WalletConfig struct {

	//币种
	Symbol string
	//配置文件路径
	configFilePath string
	//配置文件名
	configFileName string
	//区块链数据文件
	BlockchainFile string
	//本地数据库文件路径
	dbPath string
	//钱包服务API
	ServerAPI string
	//algod service token
	ServerToken string
	//默认配置内容
	DefaultConfig string
	//曲线类型
	CurveType uint32
	//链ID
	NetworkID string
}

func NewConfig(symbol string) *WalletConfig {

	c := WalletConfig{}

	//币种
	c.Symbol = symbol
	c.CurveType = CurveType

	//区块链数据
	//blockchainDir = filepath.Join("data", strings.ToLower(Symbol), "blockchain")
	//配置文件路径
	c.configFilePath = filepath.Join("conf")
	//配置文件名
	c.configFileName = c.Symbol + ".ini"
	//区块链数据文件
	c.BlockchainFile = "blockchain.db"
	//本地数据库文件路径
	c.dbPath = filepath.Join("data", strings.ToLower(c.Symbol), "db")
	//钱包服务algod API
	c.ServerAPI = ""
	//algod token
	c.ServerToken = ""

	//创建目录
	file.MkdirAll(c.dbPath)

	return &c
}
