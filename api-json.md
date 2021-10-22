####通用响应body
```
{
    "code":100200,
    "message":"status ok",
    "result":{}
}
```

#####创建或修改请求virtualService json body示例
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
        "retries":{  //必填，默认 attempts写0，其他字段不写
            "attempts":3,
            "perTryTimeout":"2s",
            "retryOn":"gateway-error,connect-failure,refused-stream" //多选，用","组合成字符串
        },
        "gateways":null
    }
}
```


#####创建或修改请求destaination json body示例
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
            "outlierDetection":{
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

#####请求获取服务列表serviceList json响应示例

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
#####请求获取服务详情serviceDetail json响应示例
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
                    }
                }
            ]
        }
    }
}



```

