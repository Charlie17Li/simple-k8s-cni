ETCDCTL_API=3 etcdctl --endpoints https://172.18.0.3:2379 --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/healthcheck-client.crt --key /etc/kubernetes/pki/etcd/healthcheck-client.key get / --prefix --keys-only


cp main /opt/cni/bin/testcni

systemctl start containerd  && systemctl start kubelet

systemctl stop kubelet && systemctl stop containerd

journalctl -exu containerd -f

ip link del testbr0 && ip link del veth9e409db6

docker cp main 95cc9b36316d:/opt/cni/bin/testcni && \
docker cp main 585965504cf3:/opt/cni/bin/testcni && \
docker cp main bb45332b0571:/opt/cni/bin/testcni