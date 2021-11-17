make build
cp /root/go/bin/kiali ./
docker-remote.sh build -f test-dockerfile -t harbor.ushareit.me/sgt/istio-config:v0.0.1 .
docker-remote.sh push harbor.ushareit.me/sgt/istio-config:v0.0.1
