etcd:
	ETCDCTL_API=3 etcdctl --endpoints https://172.18.0.4:2379 --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/healthcheck-client.crt --key /etc/kubernetes/pki/etcd/healthcheck-client.key get / --prefix --keys-only
cp:
	cp main /opt/cni/bin/testcni
start:
	systemctl start containerd  && systemctl start kubelet
stop:
	systemctl stop kubelet && systemctl stop containerd
log:
	journalctl -exu containerd -f
link:
	ip link del testbr0 && ip link del veth9e409db6

init:
	# 获取到 etcd的IP地址
	kubectl get node cni-control-plane -ojson | jq '.status.addresses | .[] | select(.type=="InternalIP").address' > etcd.conf
	docker cp main cni-worker2:/opt/cni/bin/testcni && \
	docker cp main cni-control-plane:/opt/cni/bin/testcni && \
	docker cp main cni-worker:/opt/cni/bin/testcni

	docker exec cni-worker2 touch /root/test-cni.log /root/log.error.txt && \
	docker exec cni-control-plane touch /root/test-cni.log /root/log.error.txt && \
	docker exec cni-worker touch /root/test-cni.log /root/log.error.txt

	rm -rf ./etcdca
	docker cp -a cni-control-plane:/etc/kubernetes/pki/etcd/ ./etcdca
	docker cp ./etcdca/ cni-worker:/etc/kubernetes/pki/etcd
	docker cp ./etcdca/ cni-worker2:/etc/kubernetes/pki/etcd
#
#	docker cp testcni.conf cni-worker2:/etc/cni/net.d/testcni.conf && \
#	docker cp testcni.conf cni-control-plane:/etc/cni/net.d/testcni.conf && \
#	docker cp testcni.conf cni-worker:/etc/cni/net.d/testcni.conf



