make build
cp /root/go/bin/kiali ./
docker build -f test-dockerfile -t harbor.ushareit.me/sgt/kiali:v1.39 .
docker push harbor.ushareit.me/sgt/kiali:v1.39 
