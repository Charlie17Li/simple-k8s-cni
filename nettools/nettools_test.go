package nettools

import (
	"fmt"
	"github.com/containernetworking/plugins/pkg/ns"
	// "github.com/vishvananda/netlink"
	"testcni/ipam"
	"testing"
)

var (
	brName = "testbr0"
	cidr   = "10.244.0.1/16"
	ifName = "eth0"
	podIP  = "10.244.0.2/24"
	mtu    = 1500
	nspath = "/run/netns/test.net.1"
)

func clean() {
	//删除veth対

	//删除bridge
}

func TestNettools(t *testing.T) {

	// err := ns.IsNSorErr(nspath)
	// if err != nil {
	// 	fmt.Println("IsNSorErr 失败: ", err.Error())
	// }

	// _, err = os.Open(nspath)
	// if err != nil {
	// 	fmt.Println("Open 失败: ", err.Error())
	// }
	netns, err := ns.GetNS(nspath)
	if err != nil {
		fmt.Println("获取 ns 失败: ", err.Error())
		return
	}

	// _nettools(brName, cidr, ifName, podIP, mtu, netns)
	if err = CreateBridgeAndCreateVethAndSetNetworkDeviceStatusAndSetVethMaster(brName, cidr, ifName, podIP, mtu, netns); err != nil {
		t.Errorf("失败：%v", err)
		return
	}

	// brName = "testbr0"
	// podIP = "10.244.1.3/24"
	// mtu = 1500
	// netns, err = ns.GetNS("/run/netns/test.net.2")
	// if err != nil {
	// 	fmt.Println("获取 ns 失败: ", err.Error())
	// 	return
	// }
	// _nettools(brName, cidr, ifName, podIP, mtu, netns)

	// 目前同一台主机上的 pod 可以 ping 通了
	// 接下来要让不同节点上的 pod 互相通信了
	/**
	 * 手动操作
	 * 	1. 主机上添加路由规则: ip route add 10.244.2.0/24 via 192.168.98.144 dev ens33
	 *  2. 对方主机也添加
	 *  3. 将双方主机上的网卡添加进网桥: brctl addif testbr0 ens33
	 * 以上手动操作可成功
	 * TODO: 接下来要给它转化成代码
	 */

	ipam.Init("10.244.0.0", "16")
	is, err := ipam.GetIpamService()
	if err != nil {
		fmt.Println("ipam 初始化失败: ", err.Error())
		return
	}

	fmt.Println("成功: ", is.MaskIP)
	// t.Equal(is.MaskIP, "255.255.0.0")

	names, err := is.Get().NodeNames()
	if err != nil {
		fmt.Println("这里的 err 是: ", err.Error())
		return
	}

	// t.Equal(len(names), 3)

	for _, name := range names {
		fmt.Println("这里的 name 是: ", name)
		ip, err := is.Get().NodeIp(name)
		if err != nil {
			fmt.Println("这里的 err 是: ", err.Error())
			return
		}
		fmt.Println("这里的 ip 是: ", ip)
	}

}
