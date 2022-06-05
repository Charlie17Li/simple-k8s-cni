package etcd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testcni/utils"
	"time"

	"go.etcd.io/etcd/client/pkg/v3/transport"
	etcd "go.etcd.io/etcd/client/v3"
	"k8s.io/kubectl/pkg/scheme"
)

type EtcdConfig struct {
	EtcdScheme       string `json:"etcdScheme" envconfig:"APIV1_ETCD_SCHEME" default:""`
	EtcdAuthority    string `json:"etcdAuthority" envconfig:"APIV1_ETCD_AUTHORITY" default:""`
	EtcdEndpoints    string `json:"etcdEndpoints" envconfig:"APIV1_ETCD_ENDPOINTS"`
	EtcdDiscoverySrv string `json:"etcdDiscoverySrv" envconfig:"APIV1_ETCD_DISCOVERY_SRV"`
	EtcdUsername     string `json:"etcdUsername" envconfig:"APIV1_ETCD_USERNAME"`
	EtcdPassword     string `json:"etcdPassword" envconfig:"APIV1_ETCD_PASSWORD"`
	EtcdKeyFile      string `json:"etcdKeyFile" envconfig:"APIV1_ETCD_KEY_FILE"`
	EtcdCertFile     string `json:"etcdCertFile" envconfig:"APIV1_ETCD_CERT_FILE"`
	EtcdCACertFile   string `json:"etcdCACertFile" envconfig:"APIV1_ETCD_CA_CERT_FILE"`
}

type EtcdClient struct {
	client  *etcd.Client
	Version string
}

const (
	clientTimeout = 30 * time.Second
	etcdTimeout   = 2 * time.Second
	confPath      = "/root/cni/etcd.conf"
)

var (
	globalClient *EtcdClient
)

func newEtcdClient(config *EtcdConfig) (*etcd.Client, error) {
	var etcdLocation []string
	if config.EtcdAuthority != "" {
		etcdLocation = []string{config.EtcdScheme + "://" + config.EtcdAuthority}
	}
	if config.EtcdEndpoints != "" {
		etcdLocation = strings.Split(config.EtcdEndpoints, ",")
	}

	if len(etcdLocation) == 0 {
		return nil, errors.New("找不到 etcd")
	}

	tlsInfo := transport.TLSInfo{
		CertFile:      config.EtcdCertFile,
		KeyFile:       config.EtcdKeyFile,
		TrustedCAFile: config.EtcdCACertFile,
	}

	tlsConfig, err := tlsInfo.ClientConfig()

	client, err := etcd.New(etcd.Config{
		Endpoints:   etcdLocation,
		TLS:         tlsConfig,
		DialTimeout: clientTimeout,
	})

	if err != nil {
		return nil, err
	}

	return client, nil
}

var __GetEtcdClient func() (*EtcdClient, error)

func GetEtcdClient() (*EtcdClient, error) {
	return _GetEtcdClient()()
}

func getEtcdIp(path string) string {
	if file, err := os.Open(path); err == nil {
		buf := bufio.NewReader(file)
		if line, err := buf.ReadString('\n'); err == nil {
			return regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+`).FindString(line)
		}
	}
	return ""
}

func _GetEtcdClient() func() (*EtcdClient, error) {

	return func() (*EtcdClient, error) {
		if globalClient != nil {
			return globalClient, nil
		} else {
			// ETCDCTL_API=3 etcdctl --endpoints https://192.168.98.143:2379:2379 --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/healthcheck-client.crt --key /etc/kubernetes/pki/etcd/healthcheck-client.key get / --prefix --keys-only
			ip := getEtcdIp(confPath)
			if ip == "" {
				return nil, fmt.Errorf("ETCDIP找不到")
			}

			// TODO: 这里暂时把 etcd 的地址写死了
			client, err := newEtcdClient(&EtcdConfig{
				EtcdEndpoints:  "https://" + ip + ":2379",
				EtcdCertFile:   "/etc/kubernetes/pki/etcd/healthcheck-client.crt",
				EtcdKeyFile:    "/etc/kubernetes/pki/etcd/healthcheck-client.key",
				EtcdCACertFile: "/etc/kubernetes/pki/etcd/ca.crt",
			})

			if err != nil {
				utils.WriteLog("newEtcdClient 失败, err:", err.Error())
				return nil, err
			} else {
				utils.WriteLog("newEtcdClient 成功")
			}

			status, err := client.Status(context.TODO(), client.Endpoints()[0])

			if err != nil {
				utils.WriteLog("无法获取ECTD的状态")
			} else {
				utils.WriteLog("ETCD的Version为：", status.Version)
			}

			if client != nil {
				globalClient = &EtcdClient{
					client: client,
				}

				if status != nil && status.Version != "" {
					globalClient.Version = status.Version
				}
				utils.WriteLog("客户端初始化成功", status.Version)
				return globalClient, nil
			}
		}
		return nil, errors.New("初始化 etcd client 失败")
	}
}

func Init() {
	if __GetEtcdClient == nil {
		__GetEtcdClient = _GetEtcdClient()
	}
}

func (c *EtcdClient) Set(key, value string) error {
	_, err := c.client.Put(context.TODO(), key, value)

	if err != nil {
		utils.WriteLog("Set失败, key", key, ", value:", value)
		return err
	}
	return err
}

func (c *EtcdClient) Del(key string) error {
	_, err := c.client.Delete(context.TODO(), key)

	if err != nil {
		utils.WriteLog("Del失败, key:", key)
		return err
	}
	return err
}

func (c *EtcdClient) Get(key string, opts ...etcd.OpOption) (string, error) {
	resp, err := c.client.Get(context.TODO(), key, opts...)
	if err != nil {
		utils.WriteLog("Get失败, key:", key)
		return "", err
	}

	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[len(resp.Kvs)-1:][0].Value), nil
	}
	return "", nil
}

func (c *EtcdClient) GetObj(key string, opts ...etcd.OpOption) (interface{}, error) {
	resp, err := c.client.Get(context.TODO(), key, opts...)
	if err != nil {
		utils.WriteLog("GetObj失败, key:", key)
		return nil, err
	}
	decoder := scheme.Codecs.UniversalDeserializer()
	obj, _, err := decoder.Decode(resp.Kvs[len(resp.Kvs)-1:][0].Value, nil, nil)
	if err != nil {
		utils.WriteLog("Decode失败")
		return nil, err
	}
	return obj, nil
}

func (c *EtcdClient) GetKey(key string, opts ...etcd.OpOption) (string, error) {
	resp, err := c.client.Get(context.TODO(), key, opts...)
	if err != nil {
		utils.WriteLog("GetKey失败, key:", key)
		return "", err
	}

	// for _, ev := range resp.Kvs {
	// 	fmt.Println("这里的 ev 是: ", ev)
	// 	fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	// }

	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[len(resp.Kvs)-1:][0].Key), nil
	}
	return "", nil
}

func (c *EtcdClient) GetAll(key string, opts ...etcd.OpOption) ([]string, error) {
	resp, err := c.client.Get(context.TODO(), key, opts...)
	if err != nil {
		utils.WriteLog("GetALL失败, key:", key)
		return nil, err
	}

	var res []string

	for _, ev := range resp.Kvs {
		res = append(res, string(ev.Value))
	}

	return res, nil
}

func (c *EtcdClient) GetAllKey(key string, opts ...etcd.OpOption) ([]string, error) {
	resp, err := c.client.Get(context.TODO(), key, opts...)
	if err != nil {
		utils.WriteLog("GetAllKey失败, key:", key)
		return nil, err
	}

	var res []string

	for _, ev := range resp.Kvs {
		// fmt.Println("这里的 ev 是: ", ev)
		// fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		res = append(res, string(ev.Key))
	}

	return res, nil
}
