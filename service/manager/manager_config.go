package manager

import (
	gnfd "github.com/bnb-chain/greenfield-storage-provider/pkg/greenfield"
	"github.com/bnb-chain/greenfield-storage-provider/store/config"
)

// ManagerConfig defines manager service config
type ManagerConfig struct {
	SpOperatorAddress string
	ChainConfig       *gnfd.GreenfieldChainConfig
	SpDBConfig        *config.SQLDBConfig
}
