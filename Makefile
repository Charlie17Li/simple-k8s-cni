
#path=`pwd`
path=$(shell pwd)
etcd:
	ETCDCTL_API=3 etcdctl --endpoints https://172.18.0.4:2379 --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/healthcheck-client.crt --key /etc/kubernetes/pki/etcd/healthcheck-client.key get /testcni/ipam/10.244.0.0/16/cni-control-plane --prefix --keys-only
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
clean:
	kind delete clusters cni
init:
	# 修改映射路径
	sed -i "s#- hostPath:.*$$$$#- hostPath: $(path)#g" kind.yaml
	# 创建集群
	kind create cluster --config kind.yaml
	# 获取到 etcd的IP地址
	kubectl get node cni-control-plane -ojson | jq '.status.addresses | .[] | select(.type=="InternalIP").address' > etcd.conf
	# 构建main
	docker exec cni-control-plane bash /root/cni/env_init.sh
	# 将main复制到/opt/cin/bin目录下
	docker cp main cni-worker2:/opt/cni/bin/testcni && \
	docker cp main cni-control-plane:/opt/cni/bin/testcni && \
	docker cp main cni-worker:/opt/cni/bin/testcni
	# 创建日志文件（没有创建貌似有点问题）
	docker exec cni-worker2 touch /root/test-cni.log /root/log.error.txt && \
	docker exec cni-control-plane touch /root/test-cni.log /root/log.error.txt && \
	docker exec cni-worker touch /root/test-cni.log /root/log.error.txt
	# 同步证书
	rm -rf ./etcdca
	docker cp -a cni-control-plane:/etc/kubernetes/pki/etcd/ ./etcdca
	docker cp ./etcdca/ cni-worker:/etc/kubernetes/pki/etcd
	docker cp ./etcdca/ cni-worker2:/etc/kubernetes/pki/etcd
	# 同步CNI配置文件
	docker cp testcni.conf cni-worker2:/etc/cni/net.d/testcni.conf && \
	docker cp testcni.conf cni-control-plane:/etc/cni/net.d/testcni.conf && \
	docker cp testcni.conf cni-worker:/etc/cni/net.d/testcni.conf
ci:




