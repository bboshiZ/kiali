make build
cp /root/go/bin/kiali ./
docker-remote.sh build -f test-dockerfile -t registry.ushareit.me/sgt/istio-config:v0.0.1 .
docker-remote.sh push registry.ushareit.me/sgt/istio-config:v0.0.1
