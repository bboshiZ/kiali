####开发环境后端接口地址
http://10.23.5.212:20001

####通用响应body
```
{
    "code":100200,
    "message":"status ok",
    "result":{}
}
```

#####请求获取服务列表serviceList

```
curl --location --request GET 'http://10.23.5.212:20001/istio/api/namespaces/sample/services?cluster=shareit-cce-testlimit=10&page=1'
```
```
{
    "code":100200,
    "message":"status ok",
    "result":{
        "namespace":{
            "name":"sample",
            "labels":null,
            "annotations":null
        },
        "services":[
            {
                "name":"helloworld",
                "istioSidecar":true,
                "appLabel":true,
                "additionalDetailSample":null,
                "healthAnnotations":{

                },
                "labels":{
                    "app":"helloworld",
                    "service":"helloworld"
                },
                "istioReferences":[
                    {
                        "objectType":"VirtualService",
                        "name":"helloworld",
                        "namespace":"sample"
                    },
                    {
                        "objectType":"DestinationRule",
                        "name":"helloworld",
                        "namespace":"sample"
                    }
                ],
                "kialiWizard":"request_routing"
            },
            {
                "name":"httpbin",
                "istioSidecar":false,
                "appLabel":true,
                "additionalDetailSample":null,
                "healthAnnotations":{

                },
                "labels":{
                    "app":"httpbin",
                    "service":"httpbin"
                },
                "istioReferences":[

                ],
                "kialiWizard":""
            },
            {
                "name":"sleep",
                "istioSidecar":true,
                "appLabel":true,
                "additionalDetailSample":null,
                "healthAnnotations":{

                },
                "labels":{
                    "app":"sleep",
                    "service":"sleep"
                },
                "istioReferences":[
                    {
                        "objectType":"VirtualService",
                        "name":"sleep",
                        "namespace":"sample"
                    },
                    {
                        "objectType":"DestinationRule",
                        "name":"sleep",
                        "namespace":"sample"
                    }
                ],
                "kialiWizard":"request_routing"
            }
        ],
        "validations":{
            "service":{
                "helloworld":{
                    "name":"helloworld",
                    "objectType":"service",
                    "valid":true,
                    "checks":[

                    ],
                    "references":null
                },
                "httpbin":{
                    "name":"httpbin",
                    "objectType":"service",
                    "valid":true,
                    "checks":[

                    ],
                    "references":null
                },
                "sleep":{
                    "name":"sleep",
                    "objectType":"service",
                    "valid":true,
                    "checks":[

                    ],
                    "references":null
                }
            }
        }
    }
}
```

#####请求获取服务详情serviceDetail
```curl --location --request GET 'http://10.23.5.212:20001/istio/api/namespaces/sample/services/helloworld?cluster=shareit-cce-test'
```
```
{
    "code":100200,
    "message":"status ok",
    "result":{
        "service":{  //service 信息
            "name":"helloworld",
            "namespace":{
                "name":"sample",
                "labels":null,
                "annotations":null
            },
            "labels":{
                "app":"helloworld",
                "service":"helloworld"
            },
            "selectors":{
                "app":"helloworld"
            },
            "type":"ClusterIP",
            "ports":[
                {
                    "name":"http",
                    "protocol":"TCP",
                    "port":5000
                }
            ]
        },
        "istioSidecar":true, //service 是否启用服务网格
        "workloads":[  //此service下各个版本的Deployment
            {
                "name":"helloworld-v1",
                "type":"Deployment",
                "createdAt":"2021-10-13T02:59:42Z",
                "resourceVersion":"178060969",
                "istioSidecar":true,
                "additionalDetailSample":null,
                "labels":{
                    "app":"helloworld",
                    "version":"v1"
                },
                "appLabel":true,
                "versionLabel":true,
                "podCount":1,
                "healthAnnotations":{

                },
                "istioReferences":[

                ],
                "dashboardAnnotations":null
            },
            {
                "name":"helloworld-v2",
                "type":"Deployment",
                "createdAt":"2021-10-13T02:59:42Z",
                "resourceVersion":"178060974",
                "istioSidecar":true,
                "additionalDetailSample":null,
                "labels":{
                    "app":"helloworld",
                    "version":"v2"
                },
                "appLabel":true,
                "versionLabel":true,
                "podCount":1,
                "healthAnnotations":{

                },
                "istioReferences":[

                ],
                "dashboardAnnotations":null
            }
        ],
        "virtualServices":{ //此service对应的virtualServices
            "items":[
                {
                    "kind":"VirtualService",
                    "apiVersion":"networking.istio.io/v1alpha3",
                    "metadata":{
                        "name":"helloworld",
                        "namespace":"sample",
                        "selfLink":"/apis/networking.istio.io/v1alpha3/namespaces/sample/virtualservices/helloworld",
                        "uid":"3b2de387-cd61-40af-970e-5909a4e5925c",
                        "resourceVersion":"177514764",
                        "generation":1,
                        "creationTimestamp":"2021-10-20T08:33:47Z",
                        "labels":{
                            "kiali_wizard":"request_routing"
                        }
                    }
              
                }
            ]
        },
        "destinationRules":{ //此service对应的destinationRules
            "items":[
                {
                    "kind":"DestinationRule",
                    "apiVersion":"networking.istio.io/v1alpha3",
                    "metadata":{
                        "name":"helloworld",
                        "namespace":"sample",
                        "selfLink":"/apis/networking.istio.io/v1alpha3/namespaces/sample/destinationrules/helloworld",
                        "uid":"4cfe71d9-6045-4b2e-aadf-23ed7d6c7fe6",
                        "resourceVersion":"177518259",
                        "generation":3,
                        "creationTimestamp":"2021-10-20T08:33:47Z",
                        "labels":{
                            "kiali_wizard":"request_routing"
                        }
                    },
                    "spec": {
                        "host": "helloworld.sample.svc.cluster.local",
                        "trafficPolicy": { // 流量配置
                            "connectionPool": { //连接池
                                "http": {
                                    "http1MaxPendingRequests": 100
                                },
                                "tcp": {
                                    "maxConnections": 100
                                }
                            },
                            "loadBalancer": { //负载均衡配置
                                "consistentHash": {
                                    "httpHeaderName": "user"
                                }
                            },
                            "outlierDetection": { //
                                "consecutiveErrors": 10
                                "consecutiveGatewayErrors":5,
                                "consecutive_5XxErrors":5,
                                "Interval":"2m",
                                "maxEjectionPercent":10,
                                "minHealthPercent":10
                            }
                        },
                        "subsets": [
                            {
                                "labels": {
                                    "version": "v1"
                                },
                                "name": "v1"
                            },
                            {
                                "labels": {
                                    "version": "v2"
                                },
                                "name": "v2"
                            }
                        ]
                    }
                }
            ]
        }
    }
}



```

#####服务开启serviceMesh
```
curl --location --request POST 'http://10.23.5.212:20001/istio/api/namespaces/sample/services/helloworld/inject?cluster=shareit-cce-test'
```
```
{
    "code":100200,
    "message":"status ok",
    "result":{}
}
```

#####服务关闭serviceMesh
```
curl --location --request POST 'http://10.23.5.212:20001/istio/api/namespaces/sample/services/helloworld/unInject?cluster=shareit-cce-test'
```
```
{
    "code":100200,
    "message":"status ok",
    "result":{}
}
```


#####创建或修改请求destaination
```
curl --location --request POST 'http://10.23.5.212:20001/istio/api/namespaces/sample/destinationrules' \
--header 'Content-Type: application/json' \
--data-raw ''
```
```
{
    "metadata":{
        "namespace":"sample",  //必填，服务所在k8s命名空间
        "name":"helloworld"   //必填，服务名
    },
    "spec":{
        "host":"helloworld.sample.svc.cluster.local", //格式为："服务名"+"."+"k8s命名空间"+"svc.cluster.local"
        "subsets":[   //非必填，根据从service接口获取的数据，选择填写
            {
                "name":"v1",
                "labels":{
                    "version":"v1"
                }
            },
            {
                "name":"v2",
                "labels":{
                    "version":"v2"
                }
            }
        ],
        "trafficPolicy":{ //非必填， 流量策略
            "loadBalancer":{ //负载均衡策略，simple和consistentHash只能二选一
                "simple":null,
                "consistentHash":{ //三选一
                    "httpHeaderName":"xiaoming",
                    "httpCookie":{
                        "name":"xiaoming",
                        "ttl":"10s"
                    },
                    "useSourceIp":true
                }
            },
            "connectionPool":{ //非必填，连接池管理
                "tcp":{
                    "maxConnections":123
                    "connectTimeout":10s  
                },
                "http":{
                    "http1MaxPendingRequests":123
                    "http2MaxRequests":1000,
                    "maxRequestsPerConnection":10000
                    "idleTimeout":"1h",
                }
            },
            "outlierDetection":{ //异常检测
                "consecutiveErrors":5,
                "consecutiveGatewayErrors":5,
                "consecutive_5XxErrors":5,
                "Interval":"2m",
                "maxEjectionPercent":10,
                "minHealthPercent":10
            }
        }
    }
}
```

#####删除destination
```
curl --location --request DELETE 'http://10.23.5.212:20001/istio/api/namespaces/sample/destinationrules/helloworld'
```

#####获取destionation详情
```
curl --location --request GET 'http://10.23.5.212:20001/istio/api/namespaces/sample/destinationrules/helloworld'
```

```
{
    "code": 200,
    "message": "status ok",
    "result": {
        "namespace": {
            "name": "sample",
            "labels": null,
            "annotations": null
        },
        "objectType": "destinationrules",
        "gateway": null,
        "virtualService": null,
        "destinationRule": {
            "kind": "DestinationRule",
            "apiVersion": "networking.istio.io/v1alpha3",
            "metadata": {
                "name": "helloworld",
                "namespace": "sample",
                "selfLink": "/apis/networking.istio.io/v1alpha3/namespaces/sample/destinationrules/helloworld",
                "uid": "5c8210dd-900e-4666-a83f-b54361157d67",
                "resourceVersion": "178562970",
                "generation": 7,
                "creationTimestamp": "2021-10-22T09:09:58Z"
            },
            "spec": {
                "host": "helloworld.sample.svc.cluster.local",
                "trafficPolicy": {
                    "loadBalancer": {
                        "consistentHash": {
                            "httpHeaderName": "abcd",
                            "useSourceIp": true
                        }
                    }
                },
                "subsets": [
                    {
                        "labels": {
                            "version": "v1"
                        },
                        "name": "v1"
                    },
                    {
                        "labels": {
                            "version": "v2"
                        },
                        "name": "v2"
                    }
                ]
            }
        }
    }
}
```

#####创建或修改请求virtualService
```
curl --location --request POST 'http://10.23.5.212:20001/istio/api/namespaces/sample/virtualservices' \
--header 'Content-Type: application/json' \
--data-raw ''
```
```
{
    "metadata":{
        "namespace":"sample",  //必填，服务名
        "name":"helloworld"    //必填，k8s命名空间
    },
    "spec":{
        "hosts":[  //必填
            "helloworld.sample.svc.cluster.local"  //默认加上此值，格式为："服务名"+"."+"k8s命名空间"+"svc.cluster.local"
            "scmp-test.ushareit.me"  //选填
        ],
        "http":[ //流量规则
            {
                "route":[
                    {
                        "destination":{
                            "host":"helloworld.sample.svc.cluster.local",
                            "subset":"v1"
                        },
                        "weight":50
                    },
                    {
                        "destination":{
                            "host":"helloworld.sample.svc.cluster.local",
                            "subset":"v2"
                        },
                        "weight":50
                    }
                ],
                "match":[
                    {
                        "headers":{
                            "aabb":{
                                "regex":"^.*$"
                            }
                        },
                        "uri":{
                            "prefix":"/api/v1"
                        }
                    }
                ]
            }
        ],
        "fault":{  //选填，错误注入
            "delay":{
                "percentage":{
                    "value":100
                },
                "fixedDelay":"5s"
            },
            "abort":{
                "percentage":{
                    "value":11
                },
                "httpStatus":503
            }
        },
        "timeout":"2s",  //选填，超时控制
        "retries":{  //必填，默认值 attempts写0，其他字段不写
            "attempts":3,
            "perTryTimeout":"2s",
            "retryOn":"gateway-error,connect-failure,refused-stream" //多选，用","组合成字符串
        },
        "gateways":null
    }
}
```

######删除virtualservice
```
curl --location --request DELETE 'http://10.23.5.212:20001/istio/api/namespaces/sample/virtualservices/helloworld'
```

######获取virtualservice详情
```
curl --location --request GET 'http://10.23.5.212:20001/istio/api/namespaces/sample/virtualservices/helloworld'
```
```
{
    "code": 200,
    "message": "status ok",
    "result": {
        "namespace": {
            "name": "sample",
            "labels": null,
            "annotations": null
        },
        "objectType": "virtualservices",
        "gateway": null,
        "virtualService": {
            "kind": "VirtualService",
            "apiVersion": "networking.istio.io/v1alpha3",
            "metadata": {
                "name": "helloworld",
                "namespace": "sample",
                "selfLink": "/apis/networking.istio.io/v1alpha3/namespaces/sample/virtualservices/helloworld",
                "uid": "24705e3f-b0d5-411f-b2a6-35449a684128",
                "resourceVersion": "499567607",
                "generation": 1,
                "creationTimestamp": "2021-10-25T09:41:33Z"
            },
            "spec": {
                "hosts": [
                    "helloworld.sample.svc.cluster.local"
                ],
                "http": [
                    {
                        "match": [
                            {
                                "uri": {
                                    "prefix": "/api/v1"
                                }
                            }
                        ],
                        "route": [
                            {
                                "destination": {
                                    "host": "helloworld.sample.svc.cluster.local",
                                    "subset": "v1"
                                },
                                "weight": 50
                            },
                            {
                                "destination": {
                                    "host": "helloworld.sample.svc.cluster.local",
                                    "subset": "v2"
                                },
                                "weight": 50
                            }
                        ]
                    },
                    {
                        "route": [
                            {
                                "destination": {
                                    "host": "helloworld.sample.svc.cluster.local",
                                    "subset": "v1"
                                },
                                "weight": 50
                            },
                            {
                                "destination": {
                                    "host": "helloworld.sample.svc.cluster.local",
                                    "subset": "v2"
                                },
                                "weight": 50
                            }
                        ]
                    }
                ]
            }
        },
        "destinationRule": null
    }
}
```
###获取virtualService列表
```
curl --location --request GET 'http://10.23.5.212:20001/istio/api/namespaces/sample/config?objects=virtualservices&limit=10&page=1'
```

```
{
    "code": 200,
    "message": "status ok",
    "result": {
        "total_count": 10,
        "page_count": 1,
        "current_page": 1,
        "page_size": 10,
        "data": {
            "namespace": {
                "name": "sample",
                "labels": null,
                "annotations": null
            },
            "gateways": [],
            "virtualServices": {
                "permissions": {
                    "create": false,
                    "update": false,
                    "delete": false
                },
                "items": [
                    {
                        "kind": "VirtualService",
                        "apiVersion": "networking.istio.io/v1alpha3",
                        "metadata": {
                            "name": "helloworld",
                            "namespace": "sample"
                        },
                        "spec": {
                            "hosts": [
                                "helloworld.sample.svc.cluster.local"
                            ],
                            "http": [
                                {
                                    "match": [
                                        {
                                            "uri": {
                                                "prefix": "/api/v1"
                                            }
                                        }
                                    ],
                                    "route": [
                                        {
                                            "destination": {
                                                "host": "helloworld.sample.svc.cluster.local",
                                                "subset": "v1"
                                            },
                                            "weight": 1
                                        },
                                        {
                                            "destination": {
                                                "host": "helloworld.sample.svc.cluster.local",
                                                "subset": "v2"
                                            },
                                            "weight": 99
                                        }
                                    ]
                                },
                                {
                                    "route": [
                                        {
                                            "destination": {
                                                "host": "helloworld.sample.svc.cluster.local",
                                                "subset": "v1"
                                            },
                                            "weight": 50
                                        },
                                        {
                                            "destination": {
                                                "host": "helloworld.sample.svc.cluster.local",
                                                "subset": "v2"
                                            },
                                            "weight": 50
                                        }
                                    ]
                                }
                            ]
                        }
                    }
                ]
            },
            "destinationRules": {
                "permissions": {
                    "create": false,
                    "update": false,
                    "delete": false
                },
                "items": [
                    {
                        "kind": "DestinationRule",
                        "apiVersion": "networking.istio.io/v1alpha3",
                        "metadata": {
                            "name": "helloworld",
                            "namespace": "sample"
                        },
                        "spec": {
                            "host": "helloworld.sample.svc.cluster.local",
                            "trafficPolicy": {
                                "loadBalancer": {
                                    "consistentHash": {
                                        "httpHeaderName": "abcd",
                                        "useSourceIp": true
                                    }
                                }
                            },
                            "subsets": [
                                {
                                    "labels": {
                                        "version": "v1"
                                    },
                                    "name": "v1"
                                },
                                {
                                    "labels": {
                                        "version": "v2"
                                    },
                                    "name": "v2"
                                }
                            ]
                        }
                    }
                ]
            }
        }
    }
}
```


#####获取流量图serviceGraph json响应示例
[![5hxVjx.png](https://z3.ax1x.com/2021/10/25/5hxVjx.png)](https://imgtu.com/i/5hxVjx)

```
{
    "code": 100200,
    "message": "status ok",
    "result": {
        "timestamp": 1635142541,
        "duration": 600,
        "graphType": "versionedApp",
        "elements": {
            "nodes": [  //表示流量图节点信息
                {
                    "data": {
                        "id": "735840c21949d68930349c66b9723ea3", //唯一标识
                        "nodeType": "box", //box,app,service。 box代表一个框，框内含有service，app节点
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "app": "ratings",
                        "isBox": "app"
                    }
                },
                {
                    "data": {
                        "id": "e90757ecc78d09a4504d23b07048b54e",
                        "nodeType": "box",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "app": "reviews",
                        "isBox": "app"
                    }
                },
                {
                    "data": {
                        "id": "50849a4c963e5471a0027188f2c50bbd",
                        "nodeType": "app",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "workload": "productpage-v1",
                        "app": "productpage",
                        "version": "v1",
                        "traffic": [
                            {
                                "protocol": "http",
                                "rates": {
                                    "httpOut": "0.03"
                                }
                            }
                        ],
                        "isRoot": true  //表示流量的起点
                    }
                },
                {
                    "data": {
                        "id": "d647647cbac61c45875acc2a79dd09c3",
                        "parent": "735840c21949d68930349c66b9723ea3", //表示此节点属于的box的节点id
                        "nodeType": "service",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "app": "ratings",
                        "service": "ratings",
                        "destServices": [
                            {
                                "cluster": "sgt-mesh-sg2-prod",
                                "namespace": "sample",
                                "name": "ratings"
                            }
                        ],
                        "traffic": [
                            {
                                "protocol": "http",
                                "rates": {
                                    "httpIn": "0.02",
                                    "httpOut": "0.02"
                                }
                            }
                        ]
                    }
                },
                {
                    "data": {
                        "id": "c666285f732b48a35b63557349caf4ee",
                        "parent": "735840c21949d68930349c66b9723ea3",
                        "nodeType": "app",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "workload": "ratings-v1",
                        "app": "ratings",
                        "version": "v1",
                        "destServices": [
                            {
                                "cluster": "sgt-mesh-sg2-prod",
                                "namespace": "sample",
                                "name": "ratings"
                            }
                        ],
                        "traffic": [
                            {
                                "protocol": "http",
                                "rates": {
                                    "httpIn": "0.02"
                                }
                            }
                        ]
                    }
                },
                {
                    "data": {
                        "id": "1ae755149c28ac39f378dc975dbac95b",
                        "parent": "e90757ecc78d09a4504d23b07048b54e",
                        "nodeType": "service",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "app": "reviews",
                        "service": "reviews",
                        "destServices": [
                            {
                                "cluster": "sgt-mesh-sg2-prod",
                                "namespace": "sample",
                                "name": "reviews"
                            }
                        ],
                        "traffic": [
                            {
                                "protocol": "http",
                                "rates": {
                                    "httpIn": "0.03",
                                    "httpOut": "0.03"
                                }
                            }
                        ]
                    }
                },
                {
                    "data": {
                        "id": "f075cbcd822192f29dd54c392e637565",
                        "parent": "e90757ecc78d09a4504d23b07048b54e",
                        "nodeType": "app",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "workload": "reviews-v1",
                        "app": "reviews",
                        "version": "v1",
                        "destServices": [
                            {
                                "cluster": "sgt-mesh-sg2-prod",
                                "namespace": "sample",
                                "name": "reviews"
                            }
                        ],
                        "traffic": [
                            {
                                "protocol": "http",
                                "rates": {
                                    "httpIn": "0.01"
                                }
                            }
                        ]
                    }
                },
                {
                    "data": {
                        "id": "5bbc6cd1f63cbd88faf7ef4bc1ede6c7",
                        "parent": "e90757ecc78d09a4504d23b07048b54e",
                        "nodeType": "app",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "workload": "reviews-v2",
                        "app": "reviews",
                        "version": "v2",
                        "destServices": [
                            {
                                "cluster": "sgt-mesh-sg2-prod",
                                "namespace": "sample",
                                "name": "reviews"
                            }
                        ],
                        "traffic": [
                            {
                                "protocol": "http",
                                "rates": {
                                    "httpIn": "0.01",
                                    "httpOut": "0.01"
                                }
                            }
                        ]
                    }
                },
                {
                    "data": {
                        "id": "1b82aad95da1d48acbd4f9a2de873b2a",
                        "parent": "e90757ecc78d09a4504d23b07048b54e",
                        "nodeType": "app",
                        "cluster": "sgt-mesh-sg2-prod",
                        "namespace": "sample",
                        "workload": "reviews-v3",
                        "app": "reviews",
                        "version": "v3",
                        "destServices": [
                            {
                                "cluster": "sgt-mesh-sg2-prod",
                                "namespace": "sample",
                                "name": "reviews"
                            }
                        ],
                        "traffic": [
                            {
                                "protocol": "http",
                                "rates": {
                                    "httpIn": "0.01",
                                    "httpOut": "0.01"
                                }
                            }
                        ]
                    }
                }
            ],
            "edges": [ //流量图中边的信息
                {
                    "data": {
                        "id": "0286aec242a8d0a7508b6ed259e3dc00",
                        "source": "1ae755149c28ac39f378dc975dbac95b", //边的开始节点id
                        "target": "1b82aad95da1d48acbd4f9a2de873b2a", //边的结束节点id
                        "destPrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-reviews",
                        "isMTLS": "100",
                        "responseTime": "42", //95响应延时，单位ms
                        "sourcePrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-productpage",
                        "throughput": "14",  //每秒字节数bytes/sec
                        "traffic": {
                            "protocol": "http",
                            "rates": {
                                "http": "0.01", //每秒请求数，request/sec
                                "httpPercentReq": "29.4" //service下多个app,此app请求数占比
                            },
                            "responses": {
                                "200": {
                                    "flags": {
                                        "-": "100.0" //直接图中显示即可
                                    },
                                    "hosts": {
                                        "reviews.sample.svc.cluster.local": "100.0"
                                    }
                                }
                            }
                        }
                    }
                },
                {
                    "data": {
                        "id": "9358d6c427cdab1c39b47271173bf44a",
                        "source": "1ae755149c28ac39f378dc975dbac95b",
                        "target": "5bbc6cd1f63cbd88faf7ef4bc1ede6c7",
                        "destPrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-reviews",
                        "isMTLS": "100",
                        "responseTime": "41",
                        "sourcePrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-productpage",
                        "throughput": "17",
                        "traffic": {
                            "protocol": "http",
                            "rates": {
                                "http": "0.01",
                                "httpPercentReq": "35.3"
                            },
                            "responses": {
                                "200": {
                                    "flags": {
                                        "-": "100.0"
                                    },
                                    "hosts": {
                                        "reviews.sample.svc.cluster.local": "100.0"
                                    }
                                }
                            }
                        }
                    }
                },
                {
                    "data": {
                        "id": "e63e4456612c2326e96865644c8e053c",
                        "source": "1ae755149c28ac39f378dc975dbac95b",
                        "target": "f075cbcd822192f29dd54c392e637565",
                        "destPrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-reviews",
                        "isMTLS": "100",
                        "responseTime": "5",
                        "sourcePrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-productpage",
                        "throughput": "17",
                        "traffic": {
                            "protocol": "http",
                            "rates": {
                                "http": "0.01",
                                "httpPercentReq": "35.3"
                            },
                            "responses": {
                                "200": {
                                    "flags": {
                                        "-": "100.0"
                                    },
                                    "hosts": {
                                        "reviews.sample.svc.cluster.local": "100.0"
                                    }
                                }
                            }
                        }
                    }
                },
                {
                    "data": {
                        "id": "f2758b3da168400dbd17d0e512cc8c0e",
                        "source": "1b82aad95da1d48acbd4f9a2de873b2a",
                        "target": "d647647cbac61c45875acc2a79dd09c3",
                        "destPrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-ratings",
                        "isMTLS": "100",
                        "sourcePrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-reviews",
                        "traffic": {
                            "protocol": "http",
                            "rates": {
                                "http": "0.01",
                                "httpPercentReq": "100.0"
                            },
                            "responses": {
                                "200": {
                                    "flags": {
                                        "-": "100.0"
                                    },
                                    "hosts": {
                                        "ratings.sample.svc.cluster.local": "100.0"
                                    }
                                }
                            }
                        }
                    }
                },
                {
                    "data": {
                        "id": "bd46fd1d402624e5a28abdced3443f49",
                        "source": "50849a4c963e5471a0027188f2c50bbd",
                        "target": "1ae755149c28ac39f378dc975dbac95b",
                        "destPrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-reviews",
                        "isMTLS": "100",
                        "sourcePrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-productpage",
                        "traffic": {
                            "protocol": "http",
                            "rates": {
                                "http": "0.03",
                                "httpPercentReq": "100.0"
                            },
                            "responses": {
                                "200": {
                                    "flags": {
                                        "-": "100.0"
                                    },
                                    "hosts": {
                                        "reviews.sample.svc.cluster.local": "100.0"
                                    }
                                }
                            }
                        }
                    }
                },
                {
                    "data": {
                        "id": "f29b91dd192e526de4579c37db458c26",
                        "source": "5bbc6cd1f63cbd88faf7ef4bc1ede6c7",
                        "target": "d647647cbac61c45875acc2a79dd09c3",
                        "destPrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-ratings",
                        "isMTLS": "100",
                        "sourcePrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-reviews",
                        "traffic": {
                            "protocol": "http",
                            "rates": {
                                "http": "0.01",
                                "httpPercentReq": "100.0"
                            },
                            "responses": {
                                "200": {
                                    "flags": {
                                        "-": "100.0"
                                    },
                                    "hosts": {
                                        "ratings.sample.svc.cluster.local": "100.0"
                                    }
                                }
                            }
                        }
                    }
                },
                {
                    "data": {
                        "id": "4ec7f5e0edff2a19350e01279f94539c",
                        "source": "d647647cbac61c45875acc2a79dd09c3",
                        "target": "c666285f732b48a35b63557349caf4ee",
                        "destPrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-ratings",
                        "isMTLS": "100",
                        "responseTime": "5",
                        "sourcePrincipal": "spiffe://cluster.local/ns/sample/sa/bookinfo-reviews",
                        "throughput": "29",
                        "traffic": {
                            "protocol": "http",
                            "rates": {
                                "http": "0.02",
                                "httpPercentReq": "100.0"
                            },
                            "responses": {
                                "200": {
                                    "flags": {
                                        "-": "100.0"
                                    },
                                    "hosts": {
                                        "ratings.sample.svc.cluster.local": "100.0"
                                    }
                                }
                            }
                        }
                    }
                }
            ]
        }
    }
}

```