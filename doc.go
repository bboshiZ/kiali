package main

import (
	"github.com/gogo/protobuf/types"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/status"
)

/////////////////////
// SWAGGER PARAMETERS - GENERAL
// - keep this alphabetized
/////////////////////

// swagger:parameters graphAppVersion
type AppVersionParam struct {
	// The app version (label value).
	//
	// in: path
	// required: false
	Name string `json:"version"`
}

// swagger:parameters graphAggregate graphAggregateByService graphApp graphAppVersion graphService graphWorkload
type ClusterParam struct {
	// The cluster name. If not supplied queries/results will not be constrained by cluster.
	//
	// in: query
	// required: false
	Name string `json:"container"`
}

// swagger:parameters podLogs
type ContainerParam struct {
	// The pod container name. Optional for single-container pod. Otherwise required.
	//
	// in: query
	// required: false
	Name string `json:"container"`
}

// swagger:parameters istioConfigList workloadList workloadDetails workloadUpdate serviceDetails serviceUpdate appSpans serviceSpans workloadSpans appTraces serviceTraces workloadTraces errorTraces workloadValidations appList serviceMetrics aggregateMetrics appMetrics workloadMetrics istioConfigDetails istioConfigDetailsSubtype istioConfigDelete istioConfigDeleteSubtype istioConfigUpdate istioConfigUpdateSubtype serviceList appDetails graphAggregate graphAggregateByService graphApp graphAppVersion graphNamespace graphService graphWorkload namespaceMetrics customDashboard appDashboard serviceDashboard workloadDashboard istioConfigCreate istioConfigCreateSubtype namespaceUpdate namespaceTls podDetails podLogs namespaceValidations getIter8Experiments postIter8Experiments patchIter8Experiments deleteIter8Experiments podProxyDump podProxyResource istioVirtualserviceCreate istioDestinationruleCreate
type NamespaceParam struct {
	// k8s 命令空间
	//
	// in: path
	// required: true
	Name string `json:"namespace"`
}

// swagger:parameters getIter8Experiments patchIter8Experiments deleteIter8Experiments
type NameParam struct {
	// The name param
	//
	// in: path
	// required: true
	Name string `json:"name"`
}

// swagger:parameters istioConfigDetails istioConfigDetailsSubtype istioConfigDelete istioConfigDeleteSubtype istioConfigUpdate istioConfigUpdateSubtype
type ObjectNameParam struct {
	// The Istio object name.
	//
	// in: path
	// required: true
	Name string `json:"object"`
}

// swagger:parameters istioConfigDetails istioConfigDetailsSubtype istioConfigDelete istioConfigDeleteSubtype istioConfigUpdate istioConfigUpdateSubtype istioConfigCreate istioConfigCreateSubtype
type ObjectTypeParam struct {
	// The Istio object type.
	//
	// in: path
	// required: true
	// pattern: ^(gateways|virtualservices|destinationrules|serviceentries|rules|quotaspecs|quotaspecbindings)$
	Name string `json:"object_type"`
}

// swagger:parameters podDetails podLogs podProxyDump podProxyResource
type PodParam struct {
	// The pod name.
	//
	// in: path
	// required: true
	Name string `json:"pod"`
}

// swagger:parameters podProxyResource
type ResourceParam struct {
	// The discovery service resource
	//
	// in: path
	// required: true
	Name string `json:"resource"`
}

// swagger:parameters serviceDetails serviceUpdate serviceMetrics graphService graphAggregateByService serviceDashboard serviceSpans serviceTraces
type ServiceParam struct {
	// The service name.
	//
	// in: path
	// required: true
	Name string `json:"service"`
}

// swagger:parameters podLogs
type SinceTimeParam struct {
	// The start time for fetching logs. UNIX time in seconds. Default is all logs.
	//
	// in: query
	// required: false
	Name string `json:"sinceTime"`
}

// swagger:parameters podLogs
type DurationLogParam struct {
	// Query time-range duration (Golang string duration). Duration starts on
	// `sinceTime` if set, or the time for the first log message if not set.
	//
	// in: query
	// required: false
	Name string `json:"duration"`
}

// swagger:parameters workloadDetails workloadUpdate workloadValidations workloadMetrics graphWorkload workloadDashboard workloadSpans workloadTraces
type WorkloadParam struct {
	// The workload name.
	//
	// in: path
	// required: true
	Name string `json:"workload"`
}

/////////////////////
// SWAGGER PARAMETERS - GRAPH
// - keep this alphabetized
/////////////////////

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type AppendersParam struct {
	// Comma-separated list of Appenders to run. Available appenders: [aggregateNode, deadNode, healthConfig, idleNode, istio, responseTime, securityPolicy, serviceEntry, sidecarsCheck, throughput].
	//
	// in: query
	// required: false
	// default: run all appenders
	Name string `json:"appenders"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type BoxByParam struct {
	// Comma-separated list of desired node boxing. Available boxings: [app, cluster, namespace, none].
	//
	// in: query
	// required: false
	// default: none
	Name string `json:"boxBy"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type DurationGraphParam struct {
	// Query time-range duration (Golang string duration).
	//
	// in: query
	// required: false
	// default: 10m
	Name string `json:"duration"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type GraphTypeParam struct {
	// Graph type. Available graph types: [app, service, versionedApp, workload].
	//
	// in: query
	// required: false
	// default: workload
	Name string `json:"graphType"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphWorkload
type IncludeIdleEdges struct {
	// Flag for including edges that have no request traffic for the time period.
	//
	// in: query
	// required: false
	// default: false
	Name string `json:"includeIdleEdges"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphWorkload
type InjectServiceNodes struct {
	// Flag for injecting the requested service node between source and destination nodes.
	//
	// in: query
	// required: false
	// default: false
	Name string `json:"injectServiceNodes"`
}

// swagger:parameters graphNamespaces
type NamespacesParam struct {
	// Comma-separated list of namespaces to include in the graph. The namespaces must be accessible to the client.
	//
	// in: query
	// required: true
	Name string `json:"namespaces"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type QueryTimeParam struct {
	// Unix time (seconds) for query such that time range is [queryTime-duration..queryTime]. Default is now.
	//
	// in: query
	// required: false
	// default: now
	Name string `json:"queryTime"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type RateGrpcParam struct {
	// How to calculate gRPC traffic rate. One of: none | received (i.e. response_messages) | requests | sent (i.e. request_messages) | total (i.e. sent+received).
	//
	// in: query
	// required: false
	// default: requests
	Name string `json:"rateGrpc"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type RateHttpParam struct {
	// How to calculate HTTP traffic rate. One of: none | requests.
	//
	// in: query
	// required: false
	// default: requests
	Name string `json:"rateHttp"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type RateTcpParam struct {
	// How to calculate TCP traffic rate. One of: none | received (i.e. received_bytes) | sent (i.e. sent_bytes) | total (i.e. sent+received).
	//
	// in: query
	// required: false
	// default: sent
	Name string `json:"rateTcp"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type ResponseTimeParam struct {
	// Used only with responseTime appender. One of: avg | 50 | 95 | 99.
	//
	// in: query
	// required: false
	// default: 95
	Name string `json:"responseTime"`
}

// swagger:parameters graphApp graphAppVersion graphNamespaces graphService graphWorkload
type ThroughputParam struct {
	// Used only with throughput appender. One of: request | response.
	//
	// in: query
	// required: false
	// default: request
	Name string `json:"throughput"`
}

/////////////////////
// SWAGGER PARAMETERS - METRICS
// - keep this alphabetized
/////////////////////

// swagger:parameters customDashboard
type AdditionalLabelsParam struct {
	// In custom dashboards, additional labels that are made available for grouping in the UI, regardless which aggregations are defined in the MonitoringDashboard CR
	//
	// in: query
	// required: false
	Name string `json:"additionalLabels"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics customDashboard appDashboard serviceDashboard workloadDashboard
type AvgParam struct {
	// Flag for fetching histogram average. Default is true.
	//
	// in: query
	// required: false
	// default: true
	Name bool `json:"avg"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics customDashboard appDashboard serviceDashboard workloadDashboard
type ByLabelsParam struct {
	// List of labels to use for grouping metrics (via Prometheus 'by' clause).
	//
	// in: query
	// required: false
	// default: []
	Name []string `json:"byLabels[]"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics appDashboard serviceDashboard workloadDashboard
type DirectionParam struct {
	// Traffic direction: 'inbound' or 'outbound'.
	//
	// in: query
	// required: false
	// default: outbound
	Name string `json:"direction"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics customDashboard appDashboard serviceDashboard workloadDashboard
type DurationParam struct {
	// Duration of the query period, in seconds.
	//
	// in: query
	// required: false
	// default: 1800
	Name int `json:"duration"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics
type FiltersParam struct {
	// List of metrics to fetch. Fetch all metrics when empty. List entries are Kiali internal metric names.
	//
	// in: query
	// required: false
	// default: []
	Name []string `json:"filters[]"`
}

// swagger:parameters customDashboard
type LabelsFiltersParam struct {
	// In custom dashboards, labels filters to use when fetching metrics, formatted as key:value pairs. Ex: "app:foo,version:bar".
	//
	// in: query
	// required: false
	//
	Name string `json:"labelsFilters"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics customDashboard appDashboard serviceDashboard workloadDashboard
type QuantilesParam struct {
	// List of quantiles to fetch. Fetch no quantiles when empty. Ex: [0.5, 0.95, 0.99].
	//
	// in: query
	// required: false
	// default: []
	Name []string `json:"quantiles[]"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics customDashboard appDashboard serviceDashboard workloadDashboard
type RateFuncParam struct {
	// Prometheus function used to calculate rate: 'rate' or 'irate'.
	//
	// in: query
	// required: false
	// default: rate
	Name string `json:"rateFunc"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics customDashboard appDashboard serviceDashboard workloadDashboard
type RateIntervalParam struct {
	// Interval used for rate and histogram calculation.
	//
	// in: query
	// required: false
	// default: 1m
	Name string `json:"rateInterval"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics appDashboard serviceDashboard workloadDashboard
type RequestProtocolParam struct {
	// Desired request protocol for the telemetry: For example, 'http' or 'grpc'.
	//
	// in: query
	// required: false
	// default: all protocols
	Name string `json:"requestProtocol"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics appDashboard serviceDashboard workloadDashboard
type ReporterParam struct {
	// Istio telemetry reporter: 'source' or 'destination'.
	//
	// in: query
	// required: false
	// default: source
	Name string `json:"reporter"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics customDashboard appDashboard serviceDashboard workloadDashboard
type StepParam struct {
	// Step between [graph] datapoints, in seconds.
	//
	// in: query
	// required: false
	// default: 15
	Name int `json:"step"`
}

// swagger:parameters serviceMetrics aggregateMetrics appMetrics workloadMetrics
type VersionParam struct {
	// Filters metrics by the specified version.
	//
	// in: query
	// required: false
	Name string `json:"version"`
}

/////////////////////
// SWAGGER RESPONSES
/////////////////////

// NoContent: the response is empty
// swagger:response noContent
type NoContent struct {
	// in: body
	Body struct {
		// HTTP status code
		// example: 204
		// default: 204
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// BadRequestError: the client request is incorrect
//
// swagger:response badRequestError
type BadRequestError struct {
	// in: body
	Body struct {
		// HTTP status code
		// example: 400
		// default: 400
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// A NotFoundError is the error message that is generated when server could not find what was requested.
//
// swagger:response notFoundError
type NotFoundError struct {
	// in: body
	Body struct {
		// HTTP status code
		// example: 404
		// default: 404
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// A NotAcceptable is the error message that means request can't be accepted
//
// swagger:response notAcceptableError
type NotAcceptableError struct {
	// in: body
	Body struct {
		// HTTP status code
		// example: 404
		// default: 404
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// 内部操作错误
//
// swagger:response internalError
type InternalError struct {
	// in: body
	Body struct {
		// HTTP status code
		// example: 500
		// default: 500
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// A Internal is the error message that means something has gone wrong
//
// swagger:response serviceUnavailableError
type serviceUnavailableError struct {
	// in: body
	Body struct {
		// HTTP status code
		// example: 503
		// default: 503
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// HTTP status code 200 and statusInfo model in data
// swagger:response statusInfo
type swaggStatusInfoResp struct {
	// in:body
	Body status.StatusInfo
}

// HTTP status code 200 and IstioConfigList model in data
// swagger:response istioConfigList
type IstioConfigResponse struct {
	// in:body
	Body models.IstioConfigList
}

// 获取服务列表
// swagger:response serviceListResponse
type ServiceListResponse struct {
	// in: body
	Body struct {
		// HTTP status code
		// example: 200
		// default: 200
		Code    int32           `json:"code"`
		Message error           `json:"message"`
		Result  ServiceListInfo `json:"result"`
	} `json:"body"`
}

// Listing all the information related to a workload
// swagger:response serviceDetailsResponse
type ServiceDetailsResponse struct {
	// in:body
	Body models.ServiceDetails
}

// 操作结果返回
// swagger:response commonResponse
type commonResponse struct {
	// in:body
	Body struct {
		//自定义状态码
		Code int `json:"code"`
		//操作信息
		Message string `json:"message"`
		// //响应结果
		// Result interface{} `json:"reslut"`
	}
}

// type Resp struct {
// 	//自定义状态码
// 	Code int `json:"code"`
// 	//操作信息
// 	Message string `json:"message"`
// 	// //响应结果
// 	// Result interface{} `json:"reslut"`
// }

// List of Namespaces
// swagger:response namespaceList
type NamespaceListResponse struct {
	// in:body
	Body []models.Namespace
}

//////////////////
// SWAGGER MODELS
//////////////////

// Return a list of Istio components along its status
// swagger:response istioStatusResponse
type IstioStatusResponse struct {
	// in: body
	Body business.IstioComponentStatus
}

// Posted parameters for virtualservice
// swagger:parameters istioVirtualserviceCreate
type VirtualServiceQueryBody struct {
	// 服务路由规则配置
	//
	// in: body
	// required: true
	Body VirtualService
}

// Posted parameters for destinationrule
// swagger:parameters istioDestinationruleCreate
type DestinationRuleQueryBody struct {
	// 流量策略配置
	//
	// in: body
	// required: true
	Body DestinationRule
}

//example:{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"host":"helloworld.sample.svc.cluster.local","subsets":[{"name":"v1","labels":{"version":"v1"}},{"name":"v2","labels":{"version":"v2"}}],"trafficPolicy":{"loadBalancer":{"simple":null,"consistentHash":{"httpHeaderName":"xiaoming","httpCookie":{"name":"xiaoming","ttl":"10s"},"useSourceIp":true}},"connectionPool":{"tcp":{"maxConnections":123},"http":{"http1MaxPendingRequests":123}},"outlierDetection":{"consecutiveErrors":111}}}}
type DestinationRule struct {
	IstioBase
	Spec DestinationSpec `json:"spec"`
}

// A Namespace provide a scope for names
// This type is used to describe a set of objects.
//
// swagger:model namespace
type Namespace struct {
	// The id of the namespace.
	//
	// example:  istio-system
	// required: true
	Name string `json:"name"`
}

type ServiceOverview struct {
	// Name of the Service
	// required: true
	// example: reviews-v1
	Name string `json:"name"`
	// Define if Pods related to this Service has an IstioSidecar deployed
	// required: true
	// example: true
	IstioSidecar bool `json:"istioSidecar"`
	// Has label app
	// required: true
	// example: true
	AppLabel bool `json:"appLabel"`
}

type ServiceListInfo struct {
	Namespace Namespace         `json:"namespace"`
	Services  []ServiceOverview `json:"services"`
}

type DestinationSpec struct {
	// 目的主机名
	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	// 流量熔断配置
	TrafficPolicy *TrafficPolicy `protobuf:"bytes,2,opt,name=traffic_policy,json=trafficPolicy,proto3" json:"traffic_policy,omitempty"`
	// 服务版本
	Subsets []*Subset `protobuf:"bytes,3,rep,name=subsets,proto3" json:"subsets,omitempty"`

	// ExportTo             []string `protobuf:"bytes,4,rep,name=export_to,json=exportTo,proto3" json:"export_to,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type TrafficPolicy struct {
	// 负载均衡策略
	LoadBalancer *LoadBalancerSettings `protobuf:"bytes,1,opt,name=loadBalancer,json=loadBalancer,proto3" json:"load_balancer,omitempty"`
	// 连接池管理
	ConnectionPool *ConnectionPoolSettings `protobuf:"bytes,2,opt,name=connectionPool,json=connectionPool,proto3" json:"connection_pool,omitempty"`
	// 异常检测
	OutlierDetection *OutlierDetection `protobuf:"bytes,3,opt,name=outlierDetection,json=outlierDetection,proto3" json:"outlier_detection,omitempty"`

	// PortLevelSettings    []*TrafficPolicy_PortTrafficPolicy `protobuf:"bytes,5,rep,name=port_level_settings,json=portLevelSettings,proto3" json:"port_level_settings,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type OutlierDetection struct {
	//连续错误次数。对于HTTP服务，502、503、504会被认为异常，TPC服务，连接超时即异常
	ConsecutiveErrors int32 `protobuf:"varint,1,opt,name=consecutive_errors,json=consecutiveErrors,proto3" json:"consecutive_errors,omitempty"` // Deprecated: Do not use.

	// SplitExternalLocalOriginErrors bool `protobuf:"varint,8,opt,name=split_external_local_origin_errors,json=splitExternalLocalOriginErrors,proto3" json:"split_external_local_origin_errors,omitempty"`

	// ConsecutiveLocalOriginFailures *types.UInt32Value `protobuf:"bytes,9,opt,name=consecutive_local_origin_failures,json=consecutiveLocalOriginFailures,proto3" json:"consecutive_local_origin_failures,omitempty"`
	//连续网关故障
	ConsecutiveGatewayErrors *types.UInt32Value `protobuf:"bytes,6,opt,name=consecutive_gateway_errors,json=consecutiveGatewayErrors,proto3" json:"consecutive_gateway_errors,omitempty"`
	//连续 5xx 响应
	Consecutive_5XxErrors *types.UInt32Value `protobuf:"bytes,7,opt,name=consecutive_5xx_errors,json=consecutive5xxErrors,proto3" json:"consecutive_5xx_errors,omitempty"`

	//扫描分析周期，（1h/1m/1s/1ms）
	Interval string `protobuf:"bytes,2,opt,name=interval,proto3" json:"interval,omitempty"`

	//最小驱逐时间。驱逐时间会随着错误次数增加而增加。即错误次数*最小驱逐时间
	BaseEjectionTime string `protobuf:"bytes,3,opt,name=base_ejection_time,json=baseEjectionTime,proto3" json:"base_ejection_time,omitempty"`

	//负载均衡池中可以被驱逐的实例的最大比例。以免某个接口瞬时不可用，导致太多实例被驱逐，进而导致服务整体全部不可用。
	MaxEjectionPercent int32 `protobuf:"varint,4,opt,name=max_ejection_percent,json=maxEjectionPercent,proto3" json:"max_ejection_percent,omitempty"`

	MinHealthPercent     int32    `protobuf:"varint,5,opt,name=min_health_percent,json=minHealthPercent,proto3" json:"min_health_percent,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type ConnectionPoolSettings struct {
	//  HTTP 和 TCP 连接池管理
	Tcp *ConnectionPoolSettings_TCPSettings `protobuf:"bytes,1,opt,name=tcp,proto3" json:"tcp,omitempty"`
	// HTTP 连接池管理
	Http                 *ConnectionPoolSettings_HTTPSettings `protobuf:"bytes,2,opt,name=http,proto3" json:"http,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                             `json:"-"`
	XXX_unrecognized     []byte                               `json:"-"`
	XXX_sizecache        int32                                `json:"-"`
}

// Settings applicable to HTTP1.1/HTTP2/GRPC connections.
type ConnectionPoolSettings_HTTPSettings struct {
	//最大等待HTTP请求数，默认1024
	Http1MaxPendingRequests int32 `protobuf:"varint,1,opt,name=http1_max_pending_requests,json=http1MaxPendingRequests,proto3" json:"http1_max_pending_requests,omitempty"`
	//HTTP2最大连接数
	Http2MaxRequests int32 `protobuf:"varint,2,opt,name=http2_max_requests,json=http2MaxRequests,proto3" json:"http2_max_requests,omitempty"`
	//每个连接最大请求数
	MaxRequestsPerConnection int32 `protobuf:"varint,3,opt,name=max_requests_per_connection,json=maxRequestsPerConnection,proto3" json:"max_requests_per_connection,omitempty"`

	//最大重试次数
	MaxRetries int32 `protobuf:"varint,4,opt,name=max_retries,json=maxRetries,proto3" json:"max_retries,omitempty"`
	//一个连接idle状态的最大时长,默认1h （1h/1m/10s）
	IdleTimeout string `protobuf:"bytes,5,opt,name=idle_timeout,json=idleTimeout,proto3" json:"idleTimeout,omitempty"`
	// H2UpgradePolicy ConnectionPoolSettings_HTTPSettings_H2UpgradePolicy `protobuf:"varint,6,opt,name=h2_upgrade_policy,json=h2UpgradePolicy,proto3,enum=istio.networking.v1alpha3.ConnectionPoolSettings_HTTPSettings_H2UpgradePolicy" json:"h2_upgrade_policy,omitempty"`

	// UseClientProtocol    bool     `protobuf:"varint,7,opt,name=use_client_protocol,json=useClientProtocol,proto3" json:"use_client_protocol,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

// Settings common to both HTTP and TCP upstream connections.
type ConnectionPoolSettings_TCPSettings struct {
	// Envoy为上游集群建立的最大连接数
	MaxConnections int32 `protobuf:"varint,1,opt,name=max_connections,json=maxConnections,proto3" json:"max_connections,omitempty"`
	// TCP连接超时时间 （1h/1m/1s/1ms. MUST BE >=1ms. Default is 10s）
	ConnectTimeout string `protobuf:"bytes,2,opt,name=connect_timeout,json=connectTimeout,proto3" json:"connect_timeout,omitempty"`
	// If set then set SO_KEEPALIVE on the socket to enable TCP Keepalives.
	TcpKeepalive         *ConnectionPoolSettings_TCPSettings_TcpKeepalive `protobuf:"bytes,3,opt,name=tcp_keepalive,json=tcpKeepalive,proto3" json:"tcp_keepalive,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                                         `json:"-"`
	XXX_unrecognized     []byte                                           `json:"-"`
	XXX_sizecache        int32                                            `json:"-"`
}

// TCP keepalive.
type ConnectionPoolSettings_TCPSettings_TcpKeepalive struct {
	// Maximum number of keepalive probes to send without response before
	// deciding the connection is dead. Default is to use the OS level configuration
	// (unless overridden, Linux defaults to 9.)
	Probes uint32 `protobuf:"varint,1,opt,name=probes,proto3" json:"probes,omitempty"`
	// The time duration a connection needs to be idle before keep-alive
	// probes start being sent. Default is to use the OS level configuration
	// (unless overridden, Linux defaults to 7200s (ie 2 hours.)
	Time string `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	// The time duration between keep-alive probes.
	// Default is to use the OS level configuration
	// (unless overridden, Linux defaults to 75s.)
	Interval             string   `protobuf:"bytes,3,opt,name=interval,proto3" json:"interval,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type LoadBalancerSettings_ConsistentHash struct {
	ConsistentHash *LoadBalancerSettings_ConsistentHashLB `protobuf:"bytes,2,opt,name=consistent_hash,json=consistentHash,proto3,oneof" json:"consistent_hash,omitempty"`
}

type isLoadBalancerSettings_ConsistentHashLB_HashKey struct {
	LoadBalancerSettings_ConsistentHashLB_HttpHeaderName
	LoadBalancerSettings_ConsistentHashLB_HttpCookie
	LoadBalancerSettings_ConsistentHashLB_UseSourceIp
}

type LoadBalancerSettings_ConsistentHashLB_HttpHeaderName struct {
	HttpHeaderName string `protobuf:"bytes,1,opt,name=http_header_name,json=httpHeaderName,proto3,oneof" json:"http_header_name,omitempty"`
}
type LoadBalancerSettings_ConsistentHashLB_HttpCookie struct {
	HttpCookie *LoadBalancerSettings_ConsistentHashLB_HTTPCookie `protobuf:"bytes,2,opt,name=http_cookie,json=httpCookie,proto3,oneof" json:"http_cookie,omitempty"`
}

type LoadBalancerSettings_ConsistentHashLB_HTTPCookie struct {
	// Name of the cookie.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Path to set for the cookie.
	Path string `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	// Lifetime of the cookie.
	Ttl string `protobuf:"bytes,3,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

type LoadBalancerSettings_ConsistentHashLB_UseSourceIp struct {
	UseSourceIp bool `protobuf:"varint,3,opt,name=use_source_ip,json=useSourceIp,proto3,oneof" json:"use_source_ip,omitempty"`
}

type LoadBalancerSettings_ConsistentHashLB struct {
	HttpHeaderName string `json:"httpHeaderName"`
	HttpCookie     struct {
		Name string `json:"name"`
		Path string `json:"path"`
		Ttl  string `json:"ttl"`
	} `json:"httpCookie"`
	UseSourceIp            bool   `json:"useSourceIp"`
	HttpQueryParameterName string `json:"httpQueryParameterName"`
	MinimumRingSize        uint64 `json:"minimumRingSize"`
}

//负载均衡策略，simple和consistentHash只能选择一个
type LoadBalancerSettings struct {
	//example:可选值ROUND_ROBIN LEAST_CONN RANDOM PASSTHROUGH
	Simple         string                                 `protobuf:"varint,1,opt,name=simple,proto3,enum=istio.networking.v1alpha3.LoadBalancerSettings_SimpleLB,oneof" json:"simple,omitempty"`
	ConsistentHash *LoadBalancerSettings_ConsistentHashLB `protobuf:"bytes,2,opt,name=consistent_hash,json=consistentHash,proto3,oneof" json:"consistent_hash,omitempty"`
}

type Subset struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Labels               map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TrafficPolicy        *TrafficPolicy `protobuf:"bytes,3,opt,name=traffic_policy,json=trafficPolicy,proto3" json:"traffic_policy,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

type ObjectMeta struct {
	//名称
	Name string `json:"name"`
	//k8s 命名空间
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
	// Labels      map[string]string `json:"labels,omitempty" protobuf:"bytes,11,rep,name=labels"`
	// Annotations map[string]string `json:"annotations,omitempty" protobuf:"bytes,12,rep,name=annotations"`
}
type IstioBase struct {
	Metadata ObjectMeta `json:"metadata"`
}

//example: `{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"hosts":["helloworld.sample.svc.cluster.local"],"http":[{"route":[{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v1"},"weight":50},{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v2"},"weight":50}],"match":[{"headers":{"aabb":{"regex":"^.*$"}},"uri":{"prefix":"/api/v1"}}]}],"fault":{"delay":{"percentage":{"value":100},"fixedDelay":"5s"},"abort":{"percentage":{"value":11},"httpStatus":503}},"timeout":"2s","retries":{"attempts":3,"perTryTimeout":"2s","retryOn":"gateway-error,connect-failure,refused-stream"},"gateways":null}}`
type VirtualService struct {
	IstioBase
	Spec VirtualServiceSpec `json:"spec"`
}

type VirtualServiceSpec struct {
	//域名
	Hosts []string `json:"hosts,omitempty"`
	//路由配置
	Http []HTTPRoute `json:"http,omitempty"`
}

//uri 匹配，值可选prefix，exact,regex
type StringMatch_Prefix struct {
	Prefix string `json:"prefix,omitempty"`
	Exact  string `json:"exact,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

type HeaderMatch struct {
	//header 匹配，值可选prefix，exact，regex
	Prefix string `protobuf:"bytes,2,opt,name=prefix,proto3,oneof" json:"prefix/exact/regex,omitempty"`
}

type HTTPMatchRequest struct {
	//uri匹配规则, key只能为prefix，exact，regex中一个
	//example: {"prefix":"/v1","exact":"/v2/user","regex":"/v3/.*?/user"}
	Uri map[string]string ` json:"uri,omitempty"`
	// http请求头匹配规则
	//example:{"user":{"exact":"xiaoming","prefix":"xiao","regex":"^.*$"}}
	Headers map[string]map[string]string `protobuf:"bytes,5,rep,name=headers,proto3" json:"headers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

type HTTPRoute struct {
	// 匹配规则
	Match []*HTTPMatchRequest `protobuf:"bytes,1,rep,name=match,proto3" json:"match,omitempty"`
	// 路由配置
	Route []*HTTPRouteDestination `protobuf:"bytes,2,rep,name=route,proto3" json:"route,omitempty"`
	// 重试配置
	Retries *HTTPRetry `protobuf:"bytes,4,rep,name=retries,proto3" json:"retries,omitempty"`
	//错误注入配置
	//
	Fault *HTTPFaultInjection `json:"fault,omitempty"`
	// 请求超时配置
	// example: 10s
	Timeout string `json:"timeout,omitempty"`
}

type HTTPRetry struct {
	//重试次数
	Attempts int32 `protobuf:"varint,1,opt,name=attempts,proto3" json:"attempts,omitempty"`
	//每次重试超时时间
	PerTryTimeout string `json:"perTryTimeout,omitempty"`
	//重试条件，可选5xx，gateway-error，reset，connect-failure
	RetryOn string `protobuf:"bytes,3,opt,name=retry_on,json=retryOn,proto3" json:"retry_on,omitempty"`
}

type HTTPFaultInjection struct {
	//请求延时响应配置
	//required: false
	Delay struct {
		//流量百分值（1-100）
		// required: true
		Percentage struct {
			//流量百分值（1-100）
			// example: 40
			// required: true
			Value int `json:"value"`
		} `json:"percentage,omitempty"`
		//延时响应时间（5s）
		// required: true
		FixedDelay string `json:"fixedDelay,omitempty"`
	} `json:"delay,omitempty"`
	//流量丢弃配置
	//required: false
	Abort struct {
		//流量百分值（1-100）
		// required: true
		Percentage struct {
			//流量百分值（1-100）
			// example: 40
			// required: true
			Value int `json:"value"`
		} `json:"percentage,omitempty"`
		// http 响应码
		// required: true
		// example: 404
		HttpStatus int32 `json:"httpStatus"`
	} `json:"abort,omitempty"`
}

type HTTPRouteDestination struct {
	//目的服务
	Destination *Destination `protobuf:"bytes,1,opt,name=destination,proto3" json:"destination,omitempty"`
	//流量权重 (0-100)
	Weight int32 `protobuf:"varint,2,opt,name=weight,proto3" json:"weight,omitempty"`
	// Headers              *Headers `protobuf:"bytes,7,opt,name=headers,proto3" json:"headers,omitempty"`

}

type PortSelector struct {
	//端口号
	Number uint32 `protobuf:"varint,1,opt,name=number,proto3" json:"number,omitempty"`
}
type Destination struct {
	//服务名
	Host string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	//服务版本
	Subset string `protobuf:"bytes,2,opt,name=subset,proto3" json:"subset,omitempty"`
	// Port                 *PortSelector `protobuf:"bytes,3,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}
