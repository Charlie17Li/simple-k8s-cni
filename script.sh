ETCDCTL_API=3 etcdctl --endpoints https://172.18.0.3:2379 --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/healthcheck-client.crt --key /etc/kubernetes/pki/etcd/healthcheck-client.key get / --prefix --keys-only


cp main /opt/cni/bin/testcni

systemctl start containerd  && systemctl start kubelet

systemctl stop kubelet && systemctl stop containerd

ip link del testbr0 && ip link del veth1f92a1c4