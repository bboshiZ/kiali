


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

###  config

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /api/namespaces/{namespace}/istio/destinationrules | [istio destinationrule create](#istio-destinationrule-create) |  |
| POST | /api/namespaces/{namespace}/istio/virtualservices | [istio virtualservice create](#istio-virtualservice-create) |  |
  


###  services

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /api/namespaces/{namespace}/services | [service list](#service-list) |  |
  


## Paths

### <span id="istio-destinationrule-create"></span> istio destinationrule create (*istioDestinationruleCreate*)

```
POST /api/namespaces/{namespace}/istio/destinationrules
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
| [200](#istio-destinationrule-create-200) | OK | 操作结果返回 |  | [schema](#istio-destinationrule-create-200-schema) |

#### Responses


##### <span id="istio-destinationrule-create-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-destinationrule-create-200-schema"></span> Schema
   
  

[IstioDestinationruleCreateOKBody](#istio-destinationrule-create-o-k-body)

###### Inlined models

**<span id="istio-destinationrule-create-o-k-body"></span> IstioDestinationruleCreateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="istio-virtualservice-create"></span> istio virtualservice create (*istioVirtualserviceCreate*)

```
POST /api/namespaces/{namespace}/istio/virtualservices
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
| [200](#istio-virtualservice-create-200) | OK | 操作结果返回 |  | [schema](#istio-virtualservice-create-200-schema) |

#### Responses


##### <span id="istio-virtualservice-create-200"></span> 200 - 操作结果返回
Status: OK

###### <span id="istio-virtualservice-create-200-schema"></span> Schema
   
  

[IstioVirtualserviceCreateOKBody](#istio-virtualservice-create-o-k-body)

###### Inlined models

**<span id="istio-virtualservice-create-o-k-body"></span> IstioVirtualserviceCreateOKBody**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Code | int64 (formatted integer)| `int64` |  | | 自定义状态码 |  |
| Message | string| `string` |  | | 操作信息 |  |



### <span id="service-list"></span> service list (*serviceList*)

```
GET /api/namespaces/{namespace}/services
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

### <span id="authorization-policies"></span> AuthorizationPolicies


> This is used for returning an array of AuthorizationPolicies
  



[][AuthorizationPolicy](#authorization-policy)

### <span id="authorization-policy"></span> AuthorizationPolicy


> This is used for returning an AuthorizationPolicy
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Action | [interface{}](#interface)| `interface{}` |  | |  |  |
| Rules | [interface{}](#interface)| `interface{}` |  | |  |  |
| Selector | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="component-status"></span> ComponentStatus


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| IsCore | boolean| `bool` | ✓ | | When true, the component is necessary for Istio to function. Otherwise, it is an addon | `true` |
| Name | string| `string` | ✓ | | The app label value of the Istio component | `istio-ingressgateway` |
| Status | string| `string` | ✓ | | The status of a Istio component | `Not Found` |



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


> example:{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"host":"helloworld.sample.svc.cluster.local","subsets":[{"name":"v1","labels":{"version":"v1"}},{"name":"v2","labels":{"version":"v2"}}],"trafficPolicy":{"loadBalancer":{"simple":null,"consistentHash":{"httpHeaderName":"xiaoming","httpCookie":{"name":"xiaoming","ttl":"10s"},"useSourceIp":true}},"connectionPool":{"tcp":{"maxConnections":123},"http":{"http1MaxPendingRequests":123}},"outlierDetection":{"consecutiveErrors":111}}}}
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
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



### <span id="endpoint"></span> Endpoint


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| addresses | [Addresses](#addresses)| `Addresses` |  | |  |  |
| ports | [Ports](#ports)| `Ports` |  | |  |  |



### <span id="endpoints"></span> Endpoints


  

[][Endpoint](#endpoint)

### <span id="envoy-filter"></span> EnvoyFilter


> This is used for returning an EnvoyFilter
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| ConfigPatches | [interface{}](#interface)| `interface{}` |  | |  |  |
| WorkloadSelector | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="envoy-filters"></span> EnvoyFilters


> This is used for returning an array of EnvoyFilter
  



[][EnvoyFilter](#envoy-filter)

### <span id="external-service-info"></span> ExternalServiceInfo


> This is used for returning a response of Kiali Status
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` | ✓ | | The name of the service | `Istio` |
| Url | string| `string` |  | | The service url | `jaeger-query-istio-system.127.0.0.1.nip.io` |
| Version | string| `string` |  | | The installed version of the service | `0.8.0` |



### <span id="gateway"></span> Gateway


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Selector | map of string| `map[string]string` |  | |  |  |
| Servers | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="gateways"></span> Gateways


  

[][Gateway](#gateway)

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
| Headers | map of [map[string]string](#map-string-string)| `map[string]map[string]string` |  | | http请求头匹配规则 | `{"user":{"exact":"xiaoming","prefix":"xiao","regex":"^.*$"}}` |
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



### <span id="istio-component-status"></span> IstioComponentStatus


  

[][ComponentStatus](#component-status)

### <span id="istio-config-list"></span> IstioConfigList


> This type is used for returning a response of IstioConfigList
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| authorizationPolicies | [AuthorizationPolicies](#authorization-policies)| `AuthorizationPolicies` |  | |  |  |
| destinationRules | [DestinationRules](#destination-rules)| `DestinationRules` |  | |  |  |
| envoyFilters | [EnvoyFilters](#envoy-filters)| `EnvoyFilters` |  | |  |  |
| gateways | [Gateways](#gateways)| `Gateways` |  | |  |  |
| namespace | [Namespace](#namespace)| `Namespace` | ✓ | |  |  |
| peerAuthentications | [PeerAuthentications](#peer-authentications)| `PeerAuthentications` |  | |  |  |
| requestAuthentications | [RequestAuthentications](#request-authentications)| `RequestAuthentications` |  | |  |  |
| serviceEntries | [ServiceEntries](#service-entries)| `ServiceEntries` |  | |  |  |
| sidecars | [Sidecars](#sidecars)| `Sidecars` |  | |  |  |
| validations | [IstioValidations](#istio-validations)| `IstioValidations` |  | |  |  |
| virtualServices | [VirtualServices](#virtual-services)| `VirtualServices` |  | |  |  |
| workloadEntries | [WorkloadEntries](#workload-entries)| `WorkloadEntries` |  | |  |  |
| workloadGroups | [WorkloadGroups](#workload-groups)| `WorkloadGroups` |  | |  |  |



### <span id="istio-validation-key"></span> IstioValidationKey


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | |  |  |
| Namespace | string| `string` |  | |  |  |
| ObjectType | string| `string` |  | |  |  |



### <span id="istio-validations"></span> IstioValidations


  

[interface{}](#interface)

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



### <span id="object-meta"></span> ObjectMeta


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | | 名称 |  |
| Namespace | string| `string` |  | | k8s 命名空间 |  |



### <span id="outlier-detection"></span> OutlierDetection


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| BaseEjectionTime | string| `string` |  | | 最小驱逐时间。驱逐时间会随着错误次数增加而增加。即错误次数*最小驱逐时间 |  |
| ConsecutiveErrors | int32 (formatted integer)| `int32` |  | | 连续错误次数。对于HTTP服务，502、503、504会被认为异常，TPC服务，连接超时即异常 |  |
| Interval | string| `string` |  | | 扫描分析周期，（1h/1m/1s/1ms） |  |
| MaxEjectionPercent | int32 (formatted integer)| `int32` |  | | 负载均衡池中可以被驱逐的实例的最大比例。以免某个接口瞬时不可用，导致太多实例被驱逐，进而导致服务整体全部不可用。 |  |
| MinHealthPercent | int32 (formatted integer)| `int32` |  | |  |  |
| consecutive_5xx_errors | [UInt32Value](#u-int32-value)| `UInt32Value` |  | |  |  |
| consecutive_gateway_errors | [UInt32Value](#u-int32-value)| `UInt32Value` |  | |  |  |



### <span id="peer-authentication"></span> PeerAuthentication


> This is used for returning an PeerAuthentication
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Mtls | [interface{}](#interface)| `interface{}` |  | |  |  |
| PortLevelMtls | [interface{}](#interface)| `interface{}` |  | |  |  |
| Selector | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="peer-authentications"></span> PeerAuthentications


> This is used for returning an array of PeerAuthentication
  



[][PeerAuthentication](#peer-authentication)

### <span id="port"></span> Port


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Name | string| `string` |  | |  |  |
| Port | int32 (formatted integer)| `int32` |  | |  |  |
| Protocol | string| `string` |  | |  |  |



### <span id="ports"></span> Ports


  

[][Port](#port)

### <span id="request-authentication"></span> RequestAuthentication


> This is used for returning an RequestAuthentication
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| JwtRules | [interface{}](#interface)| `interface{}` |  | |  |  |
| Selector | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="request-authentications"></span> RequestAuthentications


> This is used for returning an array of RequestAuthentication
  



[][RequestAuthentication](#request-authentication)

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



### <span id="service-entries"></span> ServiceEntries


  

[][ServiceEntry](#service-entry)

### <span id="service-entry"></span> ServiceEntry


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Addresses | [interface{}](#interface)| `interface{}` |  | |  |  |
| Endpoints | [interface{}](#interface)| `interface{}` |  | |  |  |
| ExportTo | [interface{}](#interface)| `interface{}` |  | |  |  |
| Hosts | [interface{}](#interface)| `interface{}` |  | |  |  |
| Location | [interface{}](#interface)| `interface{}` |  | |  |  |
| Ports | [interface{}](#interface)| `interface{}` |  | |  |  |
| Resolution | [interface{}](#interface)| `interface{}` |  | |  |  |
| SubjectAltNames | [interface{}](#interface)| `interface{}` |  | |  |  |
| WorkloadSelector | [interface{}](#interface)| `interface{}` |  | |  |  |



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



### <span id="service-overview"></span> ServiceOverview


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| AppLabel | boolean| `bool` | ✓ | | Has label app | `true` |
| IstioSidecar | boolean| `bool` | ✓ | | Define if Pods related to this Service has an IstioSidecar deployed | `true` |
| Name | string| `string` | ✓ | | Name of the Service | `reviews-v1` |



### <span id="sidecar"></span> Sidecar


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Egress | [interface{}](#interface)| `interface{}` |  | |  |  |
| Ingress | [interface{}](#interface)| `interface{}` |  | |  |  |
| Localhost | [interface{}](#interface)| `interface{}` |  | |  |  |
| OutboundTrafficPolicy | [interface{}](#interface)| `interface{}` |  | |  |  |
| WorkloadSelector | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="sidecars"></span> Sidecars


  

[][Sidecar](#sidecar)

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



### <span id="virtual-service"></span> VirtualService


> example: `{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"hosts":["helloworld.sample.svc.cluster.local"],"http":[{"route":[{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v1"},"weight":50},{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v2"},"weight":50}],"match":[{"headers":{"aabb":{"regex":"^.*$"}},"uri":{"prefix":"/api/v1"}}]}],"fault":{"delay":{"percentage":{"value":100},"fixedDelay":"5s"},"abort":{"percentage":{"value":11},"httpStatus":503}},"timeout":"2s","retries":{"attempts":3,"perTryTimeout":"2s","retryOn":"gateway-error,connect-failure,refused-stream"},"gateways":null}}`
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
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



### <span id="workload-entries"></span> WorkloadEntries


  

[][WorkloadEntry](#workload-entry)

### <span id="workload-entry"></span> WorkloadEntry


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Address | [interface{}](#interface)| `interface{}` |  | |  |  |
| Labels | [interface{}](#interface)| `interface{}` |  | |  |  |
| Locality | [interface{}](#interface)| `interface{}` |  | |  |  |
| Network | [interface{}](#interface)| `interface{}` |  | |  |  |
| Ports | [interface{}](#interface)| `interface{}` |  | |  |  |
| ServiceAccount | [interface{}](#interface)| `interface{}` |  | |  |  |
| Weight | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="workload-group"></span> WorkloadGroup


> This is used for returning a WorkloadGroup
  





**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| APIVersion | string| `string` |  | | APIVersion defines the versioned schema of this representation of an object.
Servers should convert recognized schemas to the latest internal value, and
may reject unrecognized values.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
+optional |  |
| Kind | string| `string` |  | | Kind is a string value representing the REST resource this object represents.
Servers may infer this from the endpoint the client submits requests to.
Cannot be updated.
In CamelCase.
More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
+optional |  |
| Status | map of any | `map[string]interface{}` |  | |  |  |
| metadata | [ObjectMeta](#object-meta)| `ObjectMeta` |  | |  |  |
| spec | [Spec](#spec)| `Spec` |  | |  |  |



#### Inlined models

**<span id="spec"></span> Spec**


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| Metadata | [interface{}](#interface)| `interface{}` |  | | This is not an error, the WorkloadGroup has a Metadata inside the Spec
https://istio.io/latest/docs/reference/config/networking/workload-group/#WorkloadGroup |  |
| Probe | [interface{}](#interface)| `interface{}` |  | |  |  |
| Template | [interface{}](#interface)| `interface{}` |  | |  |  |



### <span id="workload-groups"></span> WorkloadGroups


> This is used for returning an array of WorkloadGroup
  



[][WorkloadGroup](#workload-group)

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
