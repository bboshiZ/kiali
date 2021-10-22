


# IstioConfig
IstioConfig project, observability for the Istio service mesh
  

## Informations

### Version

_

## Content negotiation

### URI Schemes
  * http
  * https

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  istio管理

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /istio/api/namespaces/{namespace}/config | [istio config list](#istio-config-list) |  |
| POST | /istio/api/namespaces/{namespace}/istio/destinationrules | [istio destination create](#istio-destination-create) |  |
| DELETE | /istio/api/namespaces/{namespace}/istio/destinationrules/{object} | [istio destination delete](#istio-destination-delete) |  |
| PUT | /istio/api/namespaces/{namespace}/istio/destinationrules/{object} | [istio destination update](#istio-destination-update) |  |
| POST | /istio/api/namespaces/{namespace}/istio/virtualservices | [istio virtual service create](#istio-virtual-service-create) |  |
| DELETE | /istio/api/namespaces/{namespace}/istio/virtualservices/{object} | [istio virtual service delete](#istio-virtual-service-delete) |  |
| PUT | /istio/api/namespaces/{namespace}/istio/virtualservices/{object} | [istio virtual service update](#istio-virtual-service-update) |  |
  


###  k8s服务

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /istio/api/api/namespaces/{namespace}/services/{service}/inject | [istio inject](#istio-inject) |  |
| POST | /istio/api/api/namespaces/{namespace}/services/{service}/unInject | [istio un inject](#istio-un-inject) |  |
| GET | /istio/api/namespaces/{namespace}/services/{service} | [service detail](#service-detail) |  |
| GET | /istio/api/namespaces/{namespace}/services | [service list](#service-list) |  |
  


###  流量图

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /istio/api/namespaces/{namespace}/services/{service}/graph | [graph service smple](#graph-service-smple) |  |
  


###  链路追踪

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /istio/api/namespaces/{namespace}/services/{service}/traces | [service traces](#service-traces) |  |
  


## Paths

### <span id="graph-service-smple"></span> graph service smple (*graphServiceSmple*)

```
GET /istio/api/namespaces/{namespace}/services/{service}/graph
```

服务流量图

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| service | `path` | string | `string` |  | ✓ |  | 服务名称 |
| duration | `query` | string | `string` |  |  | `"10m"` | 持续时间 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#graph-service-smple-200) | OK | 服务流量图请求返回值 |  | [schema](#graph-service-smple-200-schema) |

#### Responses


##### <span id="graph-service-smple-200"></span> 200 - 服务流量图请求返回值
Status: OK

###### <span id="graph-service-smple-200-schema"></span> Schema
   
  

[GraphServiceSmpleOKBody](#graph-service-smple-o-k-body)

###### Inlined models

**<span id="graph-service-smple-o-k-body"></span> GraphServiceSmpleOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int32 (formatted integer)| `int32` |  | `503`| HTTP status code | `503` |
| Message | string| `string` |  | |  |  |
| result | [Config](#config)| `models.Config` |  | |  |  |



### <span id="istio-config-list"></span> istio config list (*istioConfigList*)

```
GET /istio/api/namespaces/{namespace}/config
```

获取virtualservices destinationrules 接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| objects | `query` | string | `string` |  | ✓ |  | 对象名称, 可选值[virtualservices,destinationrules] |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-config-list-200) | OK | HTTP status code 200 and IstioConfigList model in data |  | [schema](#istio-config-list-200-schema) |

#### Responses


##### <span id="istio-config-list-200"></span> 200 - HTTP status code 200 and IstioConfigList model in data
Status: OK

###### <span id="istio-config-list-200-schema"></span> Schema
   
  

[IstioConfigListOKBody](#istio-config-list-o-k-body)

###### Inlined models

**<span id="istio-config-list-o-k-body"></span> IstioConfigListOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int32 (formatted integer)| `int32` |  | `503`| HTTP status code | `503` |
| Message | string| `string` |  | |  |  |
| result | [IstioCfgList](#istio-cfg-list)| `models.IstioCfgList` |  | |  |  |



### <span id="istio-destination-create"></span> istio destination create (*istioDestinationCreate*)

```
POST /istio/api/namespaces/{namespace}/istio/destinationrules
```

创建destinationrules接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| Body | `body` | [DestinationRule](#destination-rule) | `models.DestinationRule` | | ✓ | | 流量策略配置 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-destination-create-200) | OK | 操作结果返回 |  | [schema](#istio-destination-create-200-schema) |

#### Responses


##### <span id="istio-destination-create-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-destination-create-200-schema"></span> Schema
   
  

[IstioDestinationCreateOKBody](#istio-destination-create-o-k-body)

###### Inlined models

**<span id="istio-destination-create-o-k-body"></span> IstioDestinationCreateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-destination-delete"></span> istio destination delete (*istioDestinationDelete*)

```
DELETE /istio/api/namespaces/{namespace}/istio/destinationrules/{object}
```

删除destinationrules接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| object | `path` | string | `string` |  | ✓ |  | istio 流量规则名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-destination-delete-200) | OK | 操作结果返回 |  | [schema](#istio-destination-delete-200-schema) |

#### Responses


##### <span id="istio-destination-delete-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-destination-delete-200-schema"></span> Schema
   
  

[IstioDestinationDeleteOKBody](#istio-destination-delete-o-k-body)

###### Inlined models

**<span id="istio-destination-delete-o-k-body"></span> IstioDestinationDeleteOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-destination-update"></span> istio destination update (*istioDestinationUpdate*)

```
PUT /istio/api/namespaces/{namespace}/istio/destinationrules/{object}
```

修改destinationrules接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| object | `path` | string | `string` |  | ✓ |  | istio 流量规则名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-destination-update-200) | OK | 操作结果返回 |  | [schema](#istio-destination-update-200-schema) |

#### Responses


##### <span id="istio-destination-update-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-destination-update-200-schema"></span> Schema
   
  

[IstioDestinationUpdateOKBody](#istio-destination-update-o-k-body)

###### Inlined models

**<span id="istio-destination-update-o-k-body"></span> IstioDestinationUpdateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-inject"></span> istio inject (*istioInject*)

```
POST /istio/api/api/namespaces/{namespace}/services/{service}/inject
```

开启服务的serviceMesh

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| service | `path` | string | `string` |  | ✓ |  | 服务名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-inject-200) | OK | 操作结果返回 |  | [schema](#istio-inject-200-schema) |

#### Responses


##### <span id="istio-inject-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-inject-200-schema"></span> Schema
   
  

[IstioInjectOKBody](#istio-inject-o-k-body)

###### Inlined models

**<span id="istio-inject-o-k-body"></span> IstioInjectOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-un-inject"></span> istio un inject (*istioUnInject*)

```
POST /istio/api/api/namespaces/{namespace}/services/{service}/unInject
```

取消服务的serviceMesh

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| service | `path` | string | `string` |  | ✓ |  | 服务名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-un-inject-200) | OK | 操作结果返回 |  | [schema](#istio-un-inject-200-schema) |

#### Responses


##### <span id="istio-un-inject-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-un-inject-200-schema"></span> Schema
   
  

[IstioUnInjectOKBody](#istio-un-inject-o-k-body)

###### Inlined models

**<span id="istio-un-inject-o-k-body"></span> IstioUnInjectOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-virtual-service-create"></span> istio virtual service create (*istioVirtualServiceCreate*)

```
POST /istio/api/namespaces/{namespace}/istio/virtualservices
```

创建virtualservices接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| Body | `body` | [VirtualService](#virtual-service) | `models.VirtualService` | | ✓ | | 服务路由规则配置 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-virtual-service-create-200) | OK | 操作结果返回 |  | [schema](#istio-virtual-service-create-200-schema) |

#### Responses


##### <span id="istio-virtual-service-create-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-virtual-service-create-200-schema"></span> Schema
   
  

[IstioVirtualServiceCreateOKBody](#istio-virtual-service-create-o-k-body)

###### Inlined models

**<span id="istio-virtual-service-create-o-k-body"></span> IstioVirtualServiceCreateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-virtual-service-delete"></span> istio virtual service delete (*istioVirtualServiceDelete*)

```
DELETE /istio/api/namespaces/{namespace}/istio/virtualservices/{object}
```

删除virtualservices接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| object | `path` | string | `string` |  | ✓ |  | istio 流量规则名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-virtual-service-delete-200) | OK | 操作结果返回 |  | [schema](#istio-virtual-service-delete-200-schema) |

#### Responses


##### <span id="istio-virtual-service-delete-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-virtual-service-delete-200-schema"></span> Schema
   
  

[IstioVirtualServiceDeleteOKBody](#istio-virtual-service-delete-o-k-body)

###### Inlined models

**<span id="istio-virtual-service-delete-o-k-body"></span> IstioVirtualServiceDeleteOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-virtual-service-update"></span> istio virtual service update (*istioVirtualServiceUpdate*)

```
PUT /istio/api/namespaces/{namespace}/istio/virtualservices/{object}
```

修改virtualservices接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| object | `path` | string | `string` |  | ✓ |  | istio 流量规则名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#istio-virtual-service-update-200) | OK | 操作结果返回 |  | [schema](#istio-virtual-service-update-200-schema) |

#### Responses


##### <span id="istio-virtual-service-update-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-virtual-service-update-200-schema"></span> Schema
   
  

[IstioVirtualServiceUpdateOKBody](#istio-virtual-service-update-o-k-body)

###### Inlined models

**<span id="istio-virtual-service-update-o-k-body"></span> IstioVirtualServiceUpdateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="service-detail"></span> service detail (*serviceDetail*)

```
GET /istio/api/namespaces/{namespace}/services/{service}
```

Endpoint to get the details of a given service

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| service | `path` | string | `string` |  | ✓ |  | 服务名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#service-detail-200) | OK | Listing all the information related to a workload |  | [schema](#service-detail-200-schema) |

#### Responses


##### <span id="service-detail-200"></span> 200 - Listing all the information related to a workload
Status: OK

###### <span id="service-detail-200-schema"></span> Schema
   
  

[ServiceDetailOKBody](#service-detail-o-k-body)

###### Inlined models

**<span id="service-detail-o-k-body"></span> ServiceDetailOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int32 (formatted integer)| `int32` |  | `200`| HTTP status code | `200` |
| Message | string| `string` |  | |  |  |
| result | [ServiceDetails](#service-details)| `models.ServiceDetails` |  | |  |  |



### <span id="service-list"></span> service list (*serviceList*)

```
GET /istio/api/namespaces/{namespace}/services
```

获取service接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#service-list-200) | OK | 获取服务列表 |  | [schema](#service-list-200-schema) |

#### Responses


##### <span id="service-list-200"></span> 200 - 获取服务列表
Status: OK

###### <span id="service-list-200-schema"></span> Schema
   
  

[ServiceListOKBody](#service-list-o-k-body)

###### Inlined models

**<span id="service-list-o-k-body"></span> ServiceListOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int32 (formatted integer)| `int32` |  | `200`| HTTP status code | `200` |
| Message | string| `string` |  | |  |  |
| result | [ServiceListInfo](#service-list-info)| `models.ServiceListInfo` |  | |  |  |



### <span id="service-traces"></span> service traces (*serviceTraces*)

```
GET /istio/api/namespaces/{namespace}/services/{service}/traces
```

获取服务链路追踪信息接口

#### URI Schemes
  * http
  * https

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| namespace | `path` | string | `string` |  | ✓ |  | k8s 命令空间 |
| service | `path` | string | `string` |  | ✓ |  | 服务名称 |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#service-traces-200) | OK | 链路追踪请求返回值 |  | [schema](#service-traces-200-schema) |

#### Responses


##### <span id="service-traces-200"></span> 200 - 链路追踪请求返回值
Status: OK

###### <span id="service-traces-200-schema"></span> Schema
   
  

[ServiceTracesOKBody](#service-traces-o-k-body)

###### Inlined models

**<span id="service-traces-o-k-body"></span> ServiceTracesOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int32 (formatted integer)| `int32` |  | `503`| HTTP status code | `503` |
| Message | string| `string` |  | |  |  |
| Result | [][Trace](#trace)| `[]*models.Trace` |  | |  |  |



## Models

### <span id="additional-item"></span> AdditionalItem


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Icon | string| `string` |  | |  |  |
| Title | string| `string` |  | |  |  |
| Value | string| `string` |  | |  |  |



### <span id="address"></span> Address


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| IP | string| `string` |  | |  |  |
| Kind | string| `string` |  | |  |  |
| Name | string| `string` |  | |  |  |



### <span id="addresses"></span> Addresses


  

[][Address](#address)

### <span id="component-status"></span> ComponentStatus


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| IsCore | boolean| `bool` | ✓ | | When true, the component is necessary for Istio to function. Otherwise, it is an addon | `true` |
| Name | string| `string` | ✓ | | The app label value of the Istio component | `istio-ingressgateway` |
| Status | string| `string` | ✓ | | The status of a Istio component | `Not Found` |



### <span id="config"></span> Config


> json example: `{"timestamp":1634869647,"duration":300,"graphType":"workload","elements":{"nodes":[{"data":{"id":"39494bdc71b74a2fc2df7d288549e51a","nodeType":"service","cluster":"sgt-mesh-sg2-prod","namespace":"sample","app":"helloworld","service":"helloworld","destServices":[{"cluster":"sgt-mesh-sg2-prod","namespace":"sample","name":"helloworld"}],"traffic":[{"protocol":"http","rates":{"httpIn":"0.04","httpOut":"0.04"}}],"hasRequestRouting":true,"hasTrafficShifting":true,"hasVS":{"hostnames":["helloworld.sample.svc.cluster.local"]}}},{"data":{"id":"3f774a4cd6c4b0b08b132c1e7e5fdbf6","nodeType":"workload","cluster":"sgt-mesh-sg2-prod","namespace":"sample","workload":"helloworld-v1","app":"helloworld","version":"v1","destServices":[{"cluster":"sgt-mesh-sg2-prod","namespace":"sample","name":"helloworld"}],"traffic":[{"protocol":"http","rates":{"httpIn":"0.02"}}]}},{"data":{"id":"b8fd6a26306e650545e765c613d8a4f3","nodeType":"workload","cluster":"sgt-mesh-sg2-prod","namespace":"sample","workload":"helloworld-v2","app":"helloworld","version":"v2","destServices":[{"cluster":"sgt-mesh-sg2-prod","namespace":"sample","name":"helloworld"}],"traffic":[{"protocol":"http","rates":{"httpIn":"0.01"}}]}},{"data":{"id":"48b437feb3a74839e75f8a20382734af","nodeType":"workload","cluster":"sgt-mesh-sg2-prod","namespace":"sample","workload":"sleep","app":"sleep","version":"latest","traffic":[{"protocol":"http","rates":{"httpOut":"0.04"}}],"isRoot":true}}],"edges":[{"data":{"id":"38bddbc3f3461a5b74d8fcfa7d854524","source":"39494bdc71b74a2fc2df7d288549e51a","target":"3f774a4cd6c4b0b08b132c1e7e5fdbf6","destPrincipal":"spiffe://cluster.local/ns/sample/sa/default","isMTLS":"100","responseTime":"237","sourcePrincipal":"spiffe://cluster.local/ns/sample/sa/sleep","throughput":"12","traffic":{"protocol":"http","rates":{"http":"0.02","httpPercentReq":"63.9"},"responses":{"200":{"flags":{"-":"100.0"},"hosts":{"helloworld.sample.svc.cluster.local":"100.0"}}}}}},{"data":{"id":"8363b30ead90a12d1ad0fc6ff5785381","source":"39494bdc71b74a2fc2df7d288549e51a","target":"b8fd6a26306e650545e765c613d8a4f3","destPrincipal":"spiffe://cluster.local/ns/sample/sa/default","isMTLS":"100","responseTime":"242","sourcePrincipal":"spiffe://cluster.local/ns/sample/sa/sleep","traffic":{"protocol":"http","rates":{"http":"0.01","httpPercentReq":"36.1"},"responses":{"200":{"flags":{"-":"100.0"},"hosts":{"helloworld.sample.svc.cluster.local":"100.0"}}}}}},{"data":{"id":"99cd00d641bcc52a9ab850a31c4789e1","source":"48b437feb3a74839e75f8a20382734af","target":"39494bdc71b74a2fc2df7d288549e51a","destPrincipal":"spiffe://cluster.local/ns/sample/sa/default","isMTLS":"100","sourcePrincipal":"spiffe://cluster.local/ns/sample/sa/sleep","traffic":{"protocol":"http","rates":{"http":"0.04","httpPercentReq":"100.0"},"responses":{"200":{"flags":{"-":"100.0"},"hosts":{"helloworld.sample.svc.cluster.local":"100.0"}}}}}}]}}`
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Duration | int64 (formatted integer)| `int64` |  | |  |  |
| GraphType | string| `string` |  | |  |  |
| Timestamp | int64 (formatted integer)| `int64` |  | |  |  |
| elements | [Elements](#elements)| `Elements` |  | |  |  |



### <span id="connection-pool-settings"></span> ConnectionPoolSettings


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| http | [ConnectionPoolSettingsHTTPSettings](#connection-pool-settings-http-settings)| `ConnectionPoolSettingsHTTPSettings` |  | |  |  |
| tcp | [ConnectionPoolSettingsTCPSettings](#connection-pool-settings-tcp-settings)| `ConnectionPoolSettingsTCPSettings` |  | |  |  |



### <span id="connection-pool-settings-http-settings"></span> ConnectionPoolSettings_HTTPSettings


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Http1MaxPendingRequests | int32 (formatted integer)| `int32` |  | | 最大等待HTTP请求数，默认1024 |  |
| Http2MaxRequests | int32 (formatted integer)| `int32` |  | | HTTP2最大连接数 |  |
| IdleTimeout | string| `string` |  | | 一个连接idle状态的最大时长,默认1h （1h/1m/10s） |  |
| MaxRequestsPerConnection | int32 (formatted integer)| `int32` |  | | 每个连接最大请求数 |  |
| MaxRetries | int32 (formatted integer)| `int32` |  | | 最大重试次数 |  |



### <span id="connection-pool-settings-tcp-settings"></span> ConnectionPoolSettings_TCPSettings


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ConnectTimeout | string| `string` |  | | TCP连接超时时间 （1h/1m/1s/1ms. MUST BE >=1ms. Default is 10s） |  |
| MaxConnections | int32 (formatted integer)| `int32` |  | | Envoy为上游集群建立的最大连接数 |  |
| tcp_keepalive | [ConnectionPoolSettingsTCPSettingsTCPKeepalive](#connection-pool-settings-tcp-settings-tcp-keepalive)| `ConnectionPoolSettingsTCPSettingsTCPKeepalive` |  | |  |  |



### <span id="connection-pool-settings-tcp-settings-tcp-keepalive"></span> ConnectionPoolSettings_TCPSettings_TcpKeepalive


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Interval | string| `string` |  | | The time duration between keep-alive probes.
Default is to use the OS level configuration
(unless overridden, Linux defaults to 75s.) |  |
| Probes | uint32 (formatted integer)| `uint32` |  | | Maximum number of keepalive probes to send without response before
deciding the connection is dead. Default is to use the OS level configuration
(unless overridden, Linux defaults to 9.) |  |
| Time | string| `string` |  | | The time duration a connection needs to be idle before keep-alive
probes start being sent. Default is to use the OS level configuration
(unless overridden, Linux defaults to 7200s (ie 2 hours.) |  |



### <span id="destination"></span> Destination


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Host | string| `string` |  | | 服务名 |  |
| Subset | string| `string` |  | | 服务版本 |  |



### <span id="destination-rule"></span> DestinationRule


> json example:`{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"host":"helloworld.sample.svc.cluster.local","subsets":[{"name":"v1","labels":{"version":"v1"}},{"name":"v2","labels":{"version":"v2"}}],"trafficPolicy":{"loadBalancer":{"simple":null,"consistentHash":{"httpHeaderName":"xiaoming","httpCookie":{"name":"xiaoming","ttl":"10s"},"useSourceIp":true}},"connectionPool":{"tcp":{"maxConnections":123},"http":{"http1MaxPendingRequests":123}},"outlierDetection":{"consecutiveErrors":111}}}}`
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` | ✓ | |  |  |
| spec | [DestinationSpec](#destination-spec)| `DestinationSpec` |  | |  |  |



### <span id="destination-rules"></span> DestinationRules


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Items | [][DestinationRule](#destination-rule)| `[]*DestinationRule` |  | |  |  |
| permissions | [ResourcePermissions](#resource-permissions)| `ResourcePermissions` |  | |  |  |



### <span id="destination-spec"></span> DestinationSpec


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Host | string| `string` |  | | 目的主机名 |  |
| Subsets | [][Subset](#subset)| `[]*Subset` |  | | 服务版本 |  |
| traffic_policy | [TrafficPolicy](#traffic-policy)| `TrafficPolicy` |  | |  |  |



### <span id="edge-data"></span> EdgeData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| DestPrincipal | string| `string` |  | | App Fields (not required by Cytoscape) |  |
| ID | string| `string` |  | | Cytoscape Fields |  |
| IsMTLS | string| `string` |  | |  |  |
| ResponseTime | string| `string` |  | |  |  |
| Source | string| `string` |  | |  |  |
| SourcePrincipal | string| `string` |  | |  |  |
| Target | string| `string` |  | |  |  |
| Throughput | string| `string` |  | |  |  |
| traffic | [ProtocolTraffic](#protocol-traffic)| `ProtocolTraffic` |  | |  |  |



### <span id="edge-wrapper"></span> EdgeWrapper


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [EdgeData](#edge-data)| `EdgeData` |  | |  |  |



### <span id="elements"></span> Elements


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Edges | [][EdgeWrapper](#edge-wrapper)| `[]*EdgeWrapper` |  | |  |  |
| Nodes | [][NodeWrapper](#node-wrapper)| `[]*NodeWrapper` |  | |  |  |



### <span id="endpoint"></span> Endpoint


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| addresses | [Addresses](#addresses)| `Addresses` |  | |  |  |
| ports | [Ports](#ports)| `Ports` |  | |  |  |



### <span id="endpoints"></span> Endpoints


  

[][Endpoint](#endpoint)

### <span id="external-service-info"></span> ExternalServiceInfo


> This is used for returning a response of Kiali Status
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` | ✓ | | The name of the service | `Istio` |
| Url | string| `string` |  | | The service url | `jaeger-query-istio-system.127.0.0.1.nip.io` |
| Version | string| `string` |  | | The installed version of the service | `0.8.0` |



### <span id="g-w-info"></span> GWInfo


> GWInfo contains the resolved gateway configuration if the node represents an Istio gateway
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ingressInfo | [GWInfoIngress](#g-w-info-ingress)| `GWInfoIngress` |  | |  |  |



### <span id="g-w-info-ingress"></span> GWInfoIngress


> GWInfoIngress contains the resolved gateway configuration if the node represents an Istio ingress gateway
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Hostnames | []string| `[]string` |  | | Hostnames is the list of hosts being served by the associated Istio gateways. |  |



### <span id="http-fault-injection"></span> HTTPFaultInjection


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| abort | [Abort](#abort)| `Abort` |  | |  |  |
| delay | [Delay](#delay)| `Delay` |  | |  |  |



#### Inlined models

**<span id="abort"></span> Abort**


> 流量丢弃配置
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| HttpStatus | int32 (formatted integer)| `int32` | ✓ | | http 响应码 | `404` |
| percentage | [Percentage](#percentage)| `Percentage` | ✓ | |  |  |



**<span id="percentage"></span> Percentage**


> 流量百分值（1-100）
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Value | int64 (formatted integer)| `int64` | ✓ | | 流量百分值（1-100） | `40` |



**<span id="delay"></span> Delay**


> 请求延时响应配置
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| FixedDelay | string| `string` | ✓ | | 延时响应时间（5s） |  |
| percentage | [Percentage](#percentage)| `Percentage` | ✓ | |  |  |



**<span id="percentage"></span> Percentage**


> 流量百分值（1-100）
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Value | int64 (formatted integer)| `int64` | ✓ | | 流量百分值（1-100） | `40` |



### <span id="http-match-request"></span> HTTPMatchRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Headers | map of [map[string]string](#map-string-string)| `map[string]map[string]string` |  | | http请求头匹配规则，只能选prefix，exact，regex中一个 | `{"user":{"exact":"xiaoming","prefix":"xiao","regex":"^.*$"}}` |
| Uri | map of string| `map[string]string` |  | | uri匹配规则, key只能为prefix，exact，regex中一个 | `{"exact":"/v2/user","prefix":"/v1","regex":"/v3/.*?/user"}` |



### <span id="http-retry"></span> HTTPRetry


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Attempts | int32 (formatted integer)| `int32` |  | | 重试次数 |  |
| PerTryTimeout | string| `string` |  | | 每次重试超时时间 |  |
| RetryOn | string| `string` |  | | 重试条件，可选5xx，gateway-error，reset，connect-failure |  |



### <span id="http-route"></span> HTTPRoute


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Match | [][HTTPMatchRequest](#http-match-request)| `[]*HTTPMatchRequest` |  | | 匹配规则 |  |
| Route | [][HTTPRouteDestination](#http-route-destination)| `[]*HTTPRouteDestination` |  | | 路由配置 |  |
| Timeout | string| `string` |  | | 请求超时配置 | `10s` |
| fault | [HTTPFaultInjection](#http-fault-injection)| `HTTPFaultInjection` |  | |  |  |
| retries | [HTTPRetry](#http-retry)| `HTTPRetry` |  | |  |  |



### <span id="http-route-destination"></span> HTTPRouteDestination


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Weight | int32 (formatted integer)| `int32` |  | | 流量权重 (0-100) |  |
| destination | [Destination](#destination)| `Destination` |  | |  |  |



### <span id="health-config"></span> HealthConfig


> HealthConfig maps annotations information for health
  



[HealthConfig](#health-config)

### <span id="istio-cfg-list"></span> IstioCfgList


> This type is used for returning a response of IstioConfigList
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| destinationRules | [DestinationRules](#destination-rules)| `DestinationRules` |  | |  |  |
| namespace | [Namespace](#namespace)| `Namespace` | ✓ | |  |  |
| virtualServices | [VirtualServices](#virtual-services)| `VirtualServices` |  | |  |  |



### <span id="istio-component-status"></span> IstioComponentStatus


  

[][ComponentStatus](#component-status)

### <span id="istio-validation-key"></span> IstioValidationKey


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | |  |  |
| Namespace | string| `string` |  | |  |  |
| ObjectType | string| `string` |  | |  |  |



### <span id="istio-validations"></span> IstioValidations


  

[interface{}](#interface)

### <span id="key-value"></span> KeyValue


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Key | string| `string` |  | |  |  |
| Value | [interface{}](#interface)| `interface{}` |  | |  |  |
| type | [ValueType](#value-type)| `ValueType` |  | |  |  |



### <span id="load-balancer-settings"></span> LoadBalancerSettings


> 负载均衡策略，simple和consistentHash只能选择一个
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Simple | string| `string` |  | |  | `可选值ROUND_ROBIN LEAST_CONN RANDOM PASSTHROUGH` |
| consistent_hash | [LoadBalancerSettingsConsistentHashLB](#load-balancer-settings-consistent-hash-l-b)| `LoadBalancerSettingsConsistentHashLB` |  | |  |  |



### <span id="load-balancer-settings-consistent-hash-l-b"></span> LoadBalancerSettings_ConsistentHashLB


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| HttpHeaderName | string| `string` |  | |  |  |
| HttpQueryParameterName | string| `string` |  | |  |  |
| MinimumRingSize | uint64 (formatted integer)| `uint64` |  | |  |  |
| UseSourceIp | boolean| `bool` |  | |  |  |
| httpCookie | [HTTPCookie](#http-cookie)| `HTTPCookie` |  | |  |  |



#### Inlined models

**<span id="http-cookie"></span> HttpCookie**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | |  |  |
| Path | string| `string` |  | |  |  |
| Ttl | string| `string` |  | |  |  |



### <span id="log"></span> Log


> Log is a log emitted in a span
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Fields | [][KeyValue](#key-value)| `[]*KeyValue` |  | |  |  |
| Timestamp | uint64 (formatted integer)| `uint64` |  | |  |  |



### <span id="m-tls-status"></span> MTLSStatus


> MTLSStatus describes the current mTLS status of a mesh entity
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Status | string| `string` | ✓ | | mTLS status: MTLS_ENABLED, MTLS_PARTIALLY_ENABLED, MTLS_NOT_ENABLED | `MTLS_ENABLED` |



### <span id="namespace"></span> Namespace


> A Namespace provide a scope for names
This type is used to describe a set of objects.
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Annotations | map of string| `map[string]string` |  | | Specific annotations used in Kiali |  |
| Labels | map of string| `map[string]string` |  | | Labels for Namespace |  |
| Name | string| `string` | ✓ | | The id of the namespace. | `istio-system` |



### <span id="node-data"></span> NodeData


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Aggregate | string| `string` |  | |  |  |
| App | string| `string` |  | |  |  |
| Cluster | string| `string` |  | |  |  |
| DestServices | [][ServiceName](#service-name)| `[]*ServiceName` |  | |  |  |
| HasCB | boolean| `bool` |  | |  |  |
| HasFaultInjection | boolean| `bool` |  | |  |  |
| HasMissingSC | boolean| `bool` |  | |  |  |
| HasRequestRouting | boolean| `bool` |  | |  |  |
| HasRequestTimeout | boolean| `bool` |  | |  |  |
| HasTCPTrafficShifting | boolean| `bool` |  | |  |  |
| HasTrafficShifting | boolean| `bool` |  | |  |  |
| ID | string| `string` |  | | Cytoscape Fields |  |
| IsBox | string| `string` |  | |  |  |
| IsDead | boolean| `bool` |  | |  |  |
| IsIdle | boolean| `bool` |  | |  |  |
| IsInaccessible | boolean| `bool` |  | |  |  |
| IsOutside | boolean| `bool` |  | |  |  |
| IsRoot | boolean| `bool` |  | |  |  |
| Namespace | string| `string` |  | |  |  |
| NodeType | string| `string` |  | | App Fields (not required by Cytoscape) |  |
| Parent | string| `string` |  | |  |  |
| Service | string| `string` |  | |  |  |
| Traffic | [][ProtocolTraffic](#protocol-traffic)| `[]*ProtocolTraffic` |  | |  |  |
| Version | string| `string` |  | |  |  |
| Workload | string| `string` |  | |  |  |
| hasHealthConfig | [HealthConfig](#health-config)| `HealthConfig` |  | |  |  |
| hasVS | [VSInfo](#v-s-info)| `VSInfo` |  | |  |  |
| isGateway | [GWInfo](#g-w-info)| `GWInfo` |  | |  |  |
| isServiceEntry | [SEInfo](#s-e-info)| `SEInfo` |  | |  |  |



### <span id="node-wrapper"></span> NodeWrapper


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [NodeData](#node-data)| `NodeData` |  | |  |  |



### <span id="object-meta"></span> ObjectMeta


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | | 名称 | `helloworld` |
| Namespace | string| `string` |  | | k8s 命名空间 | `sgt` |



### <span id="outlier-detection"></span> OutlierDetection


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| BaseEjectionTime | string| `string` |  | `"30s"`| 最小驱逐时间。驱逐时间会随着错误次数增加而增加。即错误次数*最小驱逐时间 | `1h/1m/1s/1ms` |
| ConsecutiveErrors | int32 (formatted integer)| `int32` |  | | 连续错误次数。对于HTTP服务，502、503、504会被认为异常，TPC服务，连接超时即异常 | `5` |
| Interval | string| `string` |  | `"10s"`| 扫描分析周期，（1h/1m/1s/1ms） |  |
| MaxEjectionPercent | int32 (formatted integer)| `int32` |  | `10`| 负载均衡池中可以被驱逐的实例的最大比例。以免某个接口瞬时不可用，导致太多实例被驱逐，进而导致服务整体全部不可用。 | `50` |
| MinHealthPercent | int32 (formatted integer)| `int32` |  | | 最小健康实例比例 | `50` |
| consecutive_5xx_errors | [UInt32Value](#u-int32-value)| `UInt32Value` |  | |  |  |
| consecutive_gateway_errors | [UInt32Value](#u-int32-value)| `UInt32Value` |  | |  |  |



### <span id="port"></span> Port


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | |  |  |
| Port | int32 (formatted integer)| `int32` |  | |  |  |
| Protocol | string| `string` |  | |  |  |



### <span id="ports"></span> Ports


  

[][Port](#port)

### <span id="process"></span> Process


> Process is the process emitting a set of spans
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ServiceName | string| `string` |  | |  |  |
| Tags | [][KeyValue](#key-value)| `[]*KeyValue` |  | |  |  |



### <span id="process-id"></span> ProcessID


  

| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| ProcessID | string| string | |  |  |



### <span id="protocol-traffic"></span> ProtocolTraffic


> ProtocolTraffic supplies all of the traffic information for a single protocol
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Protocol | string| `string` |  | |  |  |
| Rates | map of string| `map[string]string` |  | |  |  |
| responses | [Responses](#responses)| `Responses` |  | |  |  |



### <span id="reference"></span> Reference


> Reference is a reference from one span to another
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| refType | [ReferenceType](#reference-type)| `ReferenceType` |  | |  |  |
| spanID | [SpanID](#span-id)| `SpanID` |  | |  |  |
| traceID | [TraceID](#trace-id)| `TraceID` |  | |  |  |



### <span id="reference-type"></span> ReferenceType


> ReferenceType is the reference type of one span to another
  



| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| ReferenceType | string| string | | ReferenceType is the reference type of one span to another |  |



### <span id="request-health"></span> RequestHealth


> RequestHealth holds several stats about recent request errors
Inbound//Outbound are the rates of requests by protocol and status_code.
Example:   Inbound: { "http": {"200": 1.5, "400": 2.3}, "grpc": {"1": 1.2} }
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| HealthAnnotations | map of string| `map[string]string` |  | |  |  |
| Inbound | map of [map[string]float64](#map-string-float64)| `map[string]map[string]float64` |  | |  |  |
| Outbound | map of [map[string]float64](#map-string-float64)| `map[string]map[string]float64` |  | |  |  |



### <span id="resource-permissions"></span> ResourcePermissions


> ResourcePermissions holds permission flags for an object type
True means allowed.
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Create | boolean| `bool` |  | |  |  |
| Delete | boolean| `bool` |  | |  |  |
| Update | boolean| `bool` |  | |  |  |



### <span id="response-detail"></span> ResponseDetail


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| flags | [ResponseFlags](#response-flags)| `ResponseFlags` |  | |  |  |
| hosts | [ResponseHosts](#response-hosts)| `ResponseHosts` |  | |  |  |



### <span id="response-flags"></span> ResponseFlags


> "200" : {
"-"     : "80.0",
"DC"    : "10.0",
"FI,FD" : "10.0"
}, ...
  



[ResponseFlags](#response-flags)

### <span id="response-hosts"></span> ResponseHosts


> "200" : {
"www.google.com" : "80.0",
"www.yahoo.com"  : "20.0"
}, ...
  



[ResponseHosts](#response-hosts)

### <span id="responses"></span> Responses


> Responses maps responseCodes to detailed information for that code
  



[Responses](#responses)

### <span id="s-e-info"></span> SEInfo


> SEInfo provides static information about the service entry
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Hosts | []string| `[]string` |  | |  |  |
| Location | string| `string` |  | |  |  |
| Namespace | string| `string` |  | |  |  |



### <span id="service"></span> Service


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| CreatedAt | string| `string` |  | |  |  |
| ExternalName | string| `string` |  | |  |  |
| HealthAnnotations | map of string| `map[string]string` |  | |  |  |
| Ip | string| `string` |  | |  |  |
| Labels | map of string| `map[string]string` |  | |  |  |
| Name | string| `string` |  | |  |  |
| ResourceVersion | string| `string` |  | |  |  |
| Selectors | map of string| `map[string]string` |  | |  |  |
| Type | string| `string` |  | |  |  |
| namespace | [Namespace](#namespace)| `Namespace` |  | |  |  |
| ports | [Ports](#ports)| `Ports` |  | |  |  |



### <span id="service-details"></span> ServiceDetails


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| AdditionalDetails | [][AdditionalItem](#additional-item)| `[]*AdditionalItem` |  | |  |  |
| IstioSidecar | boolean| `bool` |  | |  |  |
| destinationRules | [DestinationRules](#destination-rules)| `DestinationRules` |  | |  |  |
| endpoints | [Endpoints](#endpoints)| `Endpoints` |  | |  |  |
| health | [ServiceHealth](#service-health)| `ServiceHealth` |  | |  |  |
| namespaceMTLS | [MTLSStatus](#m-tls-status)| `MTLSStatus` |  | |  |  |
| service | [Service](#service)| `Service` |  | |  |  |
| validations | [IstioValidations](#istio-validations)| `IstioValidations` |  | |  |  |
| virtualServices | [VirtualServices](#virtual-services)| `VirtualServices` |  | |  |  |
| workloads | [WorkloadOverviews](#workload-overviews)| `WorkloadOverviews` |  | |  |  |



### <span id="service-health"></span> ServiceHealth


> ServiceHealth contains aggregated health from various sources, for a given service
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| requests | [RequestHealth](#request-health)| `RequestHealth` |  | |  |  |



### <span id="service-list-info"></span> ServiceListInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Services | [][ServiceOverview](#service-overview)| `[]*ServiceOverview` |  | |  |  |
| namespace | [Namespace](#namespace)| `Namespace` |  | |  |  |



### <span id="service-name"></span> ServiceName


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Cluster | string| `string` |  | |  |  |
| Name | string| `string` |  | |  |  |
| Namespace | string| `string` |  | |  |  |



### <span id="service-overview"></span> ServiceOverview


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| AppLabel | boolean| `bool` | ✓ | | 服务标签 | `true` |
| IstioSidecar | boolean| `bool` | ✓ | | 服务是否开启serviceMesh | `true` |
| Name | string| `string` | ✓ | | 服务名称 | `reviews-v1` |



### <span id="span"></span> Span


> Span is a span denoting a piece of work in some infrastructure
When converting to UI model, ParentSpanID and Process should be dereferenced into
References and ProcessID, respectively.
When converting to ES model, ProcessID and Warnings should be omitted. Even if
included, ES with dynamic settings off will automatically ignore unneeded fields.
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Duration | uint64 (formatted integer)| `uint64` |  | |  |  |
| Flags | uint32 (formatted integer)| `uint32` |  | |  |  |
| Logs | [][Log](#log)| `[]*Log` |  | |  |  |
| OperationName | string| `string` |  | |  |  |
| References | [][Reference](#reference)| `[]*Reference` |  | |  |  |
| StartTime | uint64 (formatted integer)| `uint64` |  | |  |  |
| Tags | [][KeyValue](#key-value)| `[]*KeyValue` |  | |  |  |
| Warnings | []string| `[]string` |  | |  |  |
| parentSpanID | [SpanID](#span-id)| `SpanID` |  | |  |  |
| process | [Process](#process)| `Process` |  | |  |  |
| processID | [ProcessID](#process-id)| `ProcessID` |  | |  |  |
| spanID | [SpanID](#span-id)| `SpanID` |  | |  |  |
| traceID | [TraceID](#trace-id)| `TraceID` |  | |  |  |



### <span id="span-id"></span> SpanID


> SpanID is the id of a span
  



| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| SpanID | string| string | | SpanID is the id of a span |  |



### <span id="status-info"></span> StatusInfo


> This is used for returning a response of Kiali Status
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ExternalServices | [][ExternalServiceInfo](#external-service-info)| `[]*ExternalServiceInfo` | ✓ | | An array of external services installed |  |
| Status | map of string| `map[string]string` | ✓ | | The state of Kiali
A hash of key,values with versions of Kiali and state |  |
| WarningMessages | []string| `[]string` |  | | An array of warningMessages |  |



### <span id="subset"></span> Subset


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | |  |  |
| traffic_policy | [TrafficPolicy](#traffic-policy)| `TrafficPolicy` |  | |  |  |



### <span id="trace"></span> Trace


> json example:
`{"data":[{"traceID":"b8aa44f34c54ed762f14d3676255b035","spans":[{"traceID":"b8aa44f34c54ed762f14d3676255b035","spanID":"2f14d3676255b035","operationName":"helloworld:5000/*","references":[],"startTime":1634871946870049,"duration":104484,"tags":[{"key":"http.method","type":"string","value":"GET"},{"key":"request_size","type":"string","value":"0"},{"key":"http.url","type":"string","value":"http://helloworld:5000/hello"},{"key":"upstream_cluster","type":"string","value":"outbound|5000|v1|helloworld.sample.svc.cluster.local"},{"key":"user_agent","type":"string","value":"curl/7.79.1-DEV"},{"key":"component","type":"string","value":"proxy"},{"key":"downstream_cluster","type":"string","value":"-"},{"key":"guid:x-request-id","type":"string","value":"88985a37-0ed9-98f1-9f9d-604f14e4a508"},{"key":"response_size","type":"string","value":"60"},{"key":"response_flags","type":"string","value":"-"},{"key":"peer.address","type":"string","value":"10.116.3.57"},{"key":"node_id","type":"string","value":"sidecar~10.116.3.57~sleep-f7c4bdcff-79sx7.sample~sample.svc.cluster.local"},{"key":"upstream_cluster.name","type":"string","value":"outbound|5000|v1|helloworld.sample.svc.cluster.local"},{"key":"http.protocol","type":"string","value":"HTTP/1.1"},{"key":"istio.canonical_service","type":"string","value":"sleep"},{"key":"istio.canonical_revision","type":"string","value":"latest"},{"key":"istio.mesh_id","type":"string","value":"mesh1"},{"key":"istio.namespace","type":"string","value":"sample"},{"key":"http.status_code","type":"string","value":"200"},{"key":"span.kind","type":"string","value":"client"},{"key":"internal.span.format","type":"string","value":"zipkin"}],"logs":[],"processID":"p1","warnings":null},{"traceID":"b8aa44f34c54ed762f14d3676255b035","spanID":"c5f5080266d129ef","operationName":"helloworld:5000/*","references":[{"refType":"CHILD_OF","traceID":"b8aa44f34c54ed762f14d3676255b035","spanID":"2f14d3676255b035"}],"startTime":1634871946870374,"duration":103742,"tags":[{"key":"node_id","type":"string","value":"sidecar~10.116.3.56~helloworld-v1-578dd69f69-6djbt.sample~sample.svc.cluster.local"},{"key":"request_size","type":"string","value":"0"},{"key":"component","type":"string","value":"proxy"},{"key":"user_agent","type":"string","value":"curl/7.79.1-DEV"},{"key":"istio.mesh_id","type":"string","value":"mesh1"},{"key":"istio.canonical_service","type":"string","value":"helloworld"},{"key":"istio.canonical_revision","type":"string","value":"v1"},{"key":"http.status_code","type":"string","value":"200"},{"key":"istio.namespace","type":"string","value":"sample"},{"key":"upstream_cluster","type":"string","value":"inbound|5000||"},{"key":"http.method","type":"string","value":"GET"},{"key":"http.protocol","type":"string","value":"HTTP/1.1"},{"key":"peer.address","type":"string","value":"10.116.3.57"},{"key":"downstream_cluster","type":"string","value":"-"},{"key":"response_size","type":"string","value":"60"},{"key":"guid:x-request-id","type":"string","value":"88985a37-0ed9-98f1-9f9d-604f14e4a508"},{"key":"http.url","type":"string","value":"http://helloworld:5000/hello"},{"key":"upstream_cluster.name","type":"string","value":"inbound|5000||"},{"key":"response_flags","type":"string","value":"-"},{"key":"span.kind","type":"string","value":"server"},{"key":"internal.span.format","type":"string","value":"zipkin"}],"logs":[],"processID":"p2","warnings":null}],"processes":{"p1":{"serviceName":"sleep.sample","tags":[{"key":"ip","type":"string","value":"10.116.3.57"}]},"p2":{"serviceName":"helloworld.sample","tags":[{"key":"ip","type":"string","value":"10.116.3.56"}]}},"warnings":null}],"errors":null,"jaegerServiceName":"helloworld.sample"}`
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Processes | map of [Process](#process)| `map[string]Process` |  | |  |  |
| Spans | [][Span](#span)| `[]*Span` |  | |  |  |
| Warnings | []string| `[]string` |  | |  |  |
| traceID | [TraceID](#trace-id)| `TraceID` |  | |  |  |



### <span id="trace-id"></span> TraceID


  

| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| TraceID | string| string | |  |  |



### <span id="traffic-policy"></span> TrafficPolicy


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| connection_pool | [ConnectionPoolSettings](#connection-pool-settings)| `ConnectionPoolSettings` |  | |  |  |
| load_balancer | [LoadBalancerSettings](#load-balancer-settings)| `LoadBalancerSettings` |  | |  |  |
| outlier_detection | [OutlierDetection](#outlier-detection)| `OutlierDetection` |  | |  |  |



### <span id="u-int32-value"></span> UInt32Value


> The JSON representation for `UInt32Value` is JSON number.
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Value | uint32 (formatted integer)| `uint32` |  | | The uint32 value. |  |



### <span id="v-s-info"></span> VSInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Hostnames | []string| `[]string` |  | | Hostnames is the list of hostnames configured in the associated VSs |  |



### <span id="value-type"></span> ValueType


  

| Name | Type | Go type | Default | Description | Example |
|------|------|---------| ------- |-------------|---------|
| ValueType | string| string | |  |  |



### <span id="virtual-service"></span> VirtualService


> json example: `{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"hosts":["helloworld.sample.svc.cluster.local"],"http":[{"route":[{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v1"},"weight":50},{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v2"},"weight":50}],"match":[{"headers":{"aabb":{"regex":"^.*$"}},"uri":{"prefix":"/api/v1"}}]}],"fault":{"delay":{"percentage":{"value":100},"fixedDelay":"5s"},"abort":{"percentage":{"value":11},"httpStatus":503}},"timeout":"2s","retries":{"attempts":3,"perTryTimeout":"2s","retryOn":"gateway-error,connect-failure,refused-stream"},"gateways":null}}`
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` | ✓ | |  |  |
| spec | [VirtualServiceSpec](#virtual-service-spec)| `VirtualServiceSpec` |  | |  |  |



### <span id="virtual-service-spec"></span> VirtualServiceSpec


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Hosts | []string| `[]string` |  | | 域名 |  |
| Http | [][HTTPRoute](#http-route)| `[]*HTTPRoute` |  | | 路由配置 |  |



### <span id="virtual-services"></span> VirtualServices


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Items | [][VirtualService](#virtual-service)| `[]*VirtualService` |  | |  |  |
| permissions | [ResourcePermissions](#resource-permissions)| `ResourcePermissions` |  | |  |  |



### <span id="workload-list-item"></span> WorkloadListItem


> WorkloadListItem has the necessary information to display the console workload list
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| AppLabel | boolean| `bool` | ✓ | | Define if Pods related to this Workload has the label App | `true` |
| CreatedAt | string| `string` | ✓ | | Creation timestamp (in RFC3339 format) | `2018-07-31T12:24:17Z` |
| DashboardAnnotations | map of string| `map[string]string` |  | | Dashboard annotations |  |
| HealthAnnotations | map of string| `map[string]string` |  | | HealthAnnotations |  |
| IstioInjectionAnnotation | boolean| `bool` |  | | Define if Workload has an explicit Istio policy annotation
It's mapped as a pointer to show three values nil, true, false |  |
| IstioReferences | [][IstioValidationKey](#istio-validation-key)| `[]*IstioValidationKey` |  | | Istio References |  |
| IstioSidecar | boolean| `bool` | ✓ | | Define if Pods related to this Workload has an IstioSidecar deployed | `true` |
| Labels | map of string| `map[string]string` |  | | Workload labels |  |
| Name | string| `string` | ✓ | | Name of the workload | `reviews-v1` |
| PodCount | int64 (formatted integer)| `int64` | ✓ | | Number of current workload pods | `1` |
| ResourceVersion | string| `string` | ✓ | | Kubernetes ResourceVersion | `192892127` |
| Type | string| `string` | ✓ | | Type of the workload | `deployment` |
| VersionLabel | boolean| `bool` | ✓ | | Define if Pods related to this Workload has the label Version | `true` |
| additionalDetailSample | [AdditionalItem](#additional-item)| `AdditionalItem` |  | |  |  |



### <span id="workload-overviews"></span> WorkloadOverviews


  

[][WorkloadListItem](#workload-list-item)
