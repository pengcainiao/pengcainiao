package models

import "github.com/pengcainiao/zero/core/discov"

func EtcdCli() *discov.EtcdClient {
	return discov.Etcd()
}
