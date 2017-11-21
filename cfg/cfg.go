package cfg

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	defaultPeerDBFileName    = "peer.db"         //默认数据库文件名
	defaultFragManageBaseDir = "fragbase"        //默认段管理存储文件夹
	defaultMongoAddress      = "127.0.0.1:27017" //默认MongoDB数据库服务器地址

	fyerConfigFileName = "config.yaml"
)

func installDefaults() {
	curExecPath, err := os.Executable()
	if err != nil {
		panic("cannot get executable's path")
	}
	curExecPath, err = filepath.EvalSymlinks(curExecPath)
	if err != nil {
		panic("cannot resolve symbolic link on the executable's path")
	}

	//set default database file path
	viper.SetDefault("db_file", filepath.Join(curExecPath, defaultPeerDBFileName))
	//set default fragment manager's base path
	viper.SetDefault("frag_base", filepath.Join(curExecPath, defaultFragManageBaseDir))
	//set default mongodb server address
	viper.SetDefault("mongo_address", defaultMongoAddress)
}

func setupFileConfig() {
	viper.SetConfigFile(fyerConfigFileName)
	viper.AddConfigPath(".")
}

func setupEnvConfig() {
	viper.SetEnvPrefix("FYER")
	viper.AutomaticEnv()
}

//Init initializes the config modules, read all configs and be ready to provide config
//through its APIs
func Init() error {
	setupEnvConfig()
	setupFileConfig()
	installDefaults()

	return viper.ReadInConfig()
}
