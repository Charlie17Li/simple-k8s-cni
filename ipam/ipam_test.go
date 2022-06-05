package ipam

import (
	"fmt"
	"testing"
)

//func TestIpam(t *testing.T) {
//	test := assert.New(t)
//
//	Init("10.244.0.0", "16")
//	is, err := GetIpamService()
//	if err != nil {
//		fmt.Println("ipam 初始化失败: ", err.Error())
//		return
//	}
//
//	fmt.Println("成功: ", is.MaskIP)
//	test.Equal(is.MaskIP, "255.255.0.0")
//
//	fmt.Println("成功: ", is.MaskIP)
//	test.Equal(is.MaskIP, "255.255.0.0")
//	cidr, _ := is.Get().CIDR("cni-control-plane")
//	test.Equal(cidr, "10.244.0.0/24")
//	cidr, _ = is.Get().CIDR("cni-worker")
//	test.Equal(cidr, "10.244.1.0/24")
//	cidr, _ = is.Get().CIDR("cni-worker2")
//	test.Equal(cidr, "10.244.2.0/24")
//
//	names, err := is.Get().NodeNames()
//	if err != nil {
//		fmt.Println("这里的 err 是: ", err.Error())
//		return
//	}
//
//	test.Equal(len(names), 3)
//
//	for _, name := range names {
//		ip, err := is.Get().NodeIp(name)
//		if err != nil {
//			fmt.Println("这里的 err 是: ", err.Error())
//			return
//		}
//		fmt.Println("这里的 ip 是: ", ip)
//	}
//
//	nets, err := is.Get().AllHostNetwork()
//	if err != nil {
//		fmt.Println("这里的 err 是: ", err.Error())
//		return
//	}
//	fmt.Println("集群全部网络信息是: ", nets)
//
//	for _, net := range nets {
//		fmt.Println("其他主机的网络信息是: ", net)
//	}
//
//	currentNet, err := is.Get().HostNetwork()
//	if err != nil {
//		fmt.Println("这里的 err 是: ", err.Error())
//		return
//	}
//	fmt.Println("本机的网络信息是: ", currentNet)
//}

func TestGet_CIDR(t *testing.T) {
	var ip string
	Init("10.244.0.0", "16")
	is, err := GetIpamService()
	if err != nil {
		t.Errorf("ipam 初始化失败: %v", err.Error())
		return
	}
	if ip, err = is.Get().CIDR("cni-control-plane", is); err != nil {
		t.Errorf("获取CIDR失败:%v", err)
	}

	fmt.Println("ip是", ip)
}

func TestGet_HostNetwork(t *testing.T) {
	var network *Network
	Init("10.244.0.0", "16")
	is, err := GetIpamService()
	if err != nil {
		t.Errorf("ipam 初始化失败: %v", err.Error())
	}

	if network, err = is.Get().HostNetwork(); err != nil {
		t.Errorf("获取本机网卡信息失败")
	}
	t.Log(network)
}
