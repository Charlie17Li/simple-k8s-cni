package etcd

import (
	apiv1 "k8s.io/api/core/v1"
	"testing"
)

func TestGetObj(t *testing.T) {
	var client *EtcdClient
	var err error
	if client, err = GetEtcdClient(); err != nil || client == nil {
		t.Errorf("获取client失败, err:%v", err)
		return
	} else {
		t.Log("获取client成功")
	}

	if val, err := client.GetObj("/registry/minions/cni-control-plane"); err == nil {
		node := val.(*apiv1.Node)
		arr := node.Status.Addresses

		for _, val := range arr {
			if val.Type == "InternalIP" {
				t.Log(val.Address)
			}
		}

		return
	}

	t.Errorf("获取key失败")
}
