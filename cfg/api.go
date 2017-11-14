package cfg

import (
	"github.com/spf13/viper"
)

//DBPath returns the path of database file
func DBPath() string {
	return viper.GetString("db_file")
}

//FragBasePath returns the path of the directoty that is used to store fragments
func FragBasePath() string {
	return viper.GetString("frag_base")
}
