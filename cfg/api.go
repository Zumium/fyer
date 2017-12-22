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

//MongoAddress returns the mongodb server address in the configuration
func MongoAddress() string {
	return viper.GetString("mongo_address")
}

//Port returns the RPC service port
func Port() int {
	return viper.GetInt("port")
}

//PeerRegisterPort returns the peer registration service port
func PeerRegisterPort() int {
	return viper.GetInt("peer_register_port")
}

//CenterAddress returns the center's address
func CenterAddress() string {
	return viper.GetString("center_address")
}

//PeerID returns peer's local id (label)
func PeerID() string {
	return viper.GetString("peer_id")
}

//MaxSendRecvMsgSize returns the max size of sent or received message
func MaxSendRecvMsgSize() int {
	return viper.GetInt("max_send_recv_msg_size")
}

//FragSize returns the size of max fragment
func FragSize() int64 {
	return viper.GetInt64("frag_size")
}

//Replica returns the replica of a frag
func Replica() int {
	return viper.GetInt("replica")
}
