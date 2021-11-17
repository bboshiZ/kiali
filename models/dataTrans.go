package models

import "github.com/gogo/protobuf/types"

type ClusterM struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

//json example:`{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"host":"helloworld.sample.svc.cluster.local","subsets":[{"name":"v1","labels":{"version":"v1"}},{"name":"v2","labels":{"version":"v2"}}],"trafficPolicy":{"loadBalancer":{"simple":null,"consistentHash":{"httpHeaderName":"xiaoming","httpCookie":{"name":"xiaoming","ttl":"10s"},"useSourceIp":true}},"connectionPool":{"tcp":{"maxConnections":123},"http":{"http1MaxPendingRequests":123}},"outlierDetection":{"consecutiveErrors":111}}}}`
type DestinationRuleM struct {
	IstioBase
	Spec DestinationSpec `json:"spec"`
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
type Subset struct {
	Name          string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Labels        map[string]string `protobuf:"bytes,2,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	TrafficPolicy *TrafficPolicy    `protobuf:"bytes,3,opt,name=traffic_policy,json=trafficPolicy,proto3" json:"traffic_policy,omitempty"`
}

//负载均衡策略，simple和consistentHash只能选择一个
type LoadBalancerSettings struct {
	//example:可选值ROUND_ROBIN LEAST_CONN RANDOM PASSTHROUGH
	Simple         string                                 `protobuf:"varint,1,opt,name=simple,proto3,enum=istio.networking.v1alpha3.LoadBalancerSettings_SimpleLB,oneof" json:"simple,omitempty"`
	ConsistentHash *LoadBalancerSettings_ConsistentHashLB `protobuf:"bytes,2,opt,name=consistent_hash,json=consistentHash,proto3,oneof" json:"consistent_hash,omitempty"`
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

type OutlierDetection struct {
	//连续错误次数。对于HTTP服务，502、503、504会被认为异常，TPC服务，连接超时即异常
	//example: 5
	ConsecutiveErrors int32 `protobuf:"varint,1,opt,name=consecutive_errors,json=consecutiveErrors,proto3" json:"consecutive_errors,omitempty"` // Deprecated: Do not use.

	//连续网关故障
	ConsecutiveGatewayErrors *types.UInt32Value `protobuf:"bytes,6,opt,name=consecutive_gateway_errors,json=consecutiveGatewayErrors,proto3" json:"consecutive_gateway_errors,omitempty"`
	//连续 5xx 响应
	Consecutive_5XxErrors *types.UInt32Value `protobuf:"bytes,7,opt,name=consecutive_5xx_errors,json=consecutive5xxErrors,proto3" json:"consecutive_5xx_errors,omitempty"`

	//扫描分析周期，（1h/1m/1s/1ms）
	//default:  10s
	Interval string `protobuf:"bytes,2,opt,name=interval,proto3" json:"interval,omitempty"`

	//最小驱逐时间。驱逐时间会随着错误次数增加而增加。即错误次数*最小驱逐时间
	//example:  1h/1m/1s/1ms
	//default:  30s
	BaseEjectionTime string `protobuf:"bytes,3,opt,name=base_ejection_time,json=baseEjectionTime,proto3" json:"base_ejection_time,omitempty"`

	//负载均衡池中可以被驱逐的实例的最大比例。以免某个接口瞬时不可用，导致太多实例被驱逐，进而导致服务整体全部不可用。
	//default: 10
	//example: 50
	MaxEjectionPercent int32 `protobuf:"varint,4,opt,name=max_ejection_percent,json=maxEjectionPercent,proto3" json:"max_ejection_percent,omitempty"`
	//最小健康实例比例
	//default: 0
	//example: 50
	MinHealthPercent int32 `protobuf:"varint,5,opt,name=min_health_percent,json=minHealthPercent,proto3" json:"min_health_percent,omitempty"`
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

//json example: `{"metadata":{"namespace":"sample","name":"helloworld"},"spec":{"hosts":["helloworld.sample.svc.cluster.local"],"http":[{"route":[{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v1"},"weight":50},{"destination":{"host":"helloworld.sample.svc.cluster.local","subset":"v2"},"weight":50}],"match":[{"headers":{"aabb":{"regex":"^.*$"}},"uri":{"prefix":"/api/v1"}}]}],"fault":{"delay":{"percentage":{"value":100},"fixedDelay":"5s"},"abort":{"percentage":{"value":11},"httpStatus":503}},"timeout":"2s","retries":{"attempts":3,"perTryTimeout":"2s","retryOn":"gateway-error,connect-failure,refused-stream"},"gateways":null}}`
type VirtualServiceM struct {
	IstioBase
	Spec VirtualServiceSpec `json:"spec"`
}

type VirtualServiceSpec struct {
	//域名
	Hosts []string `json:"hosts,omitempty"`
	//路由配置
	Http []HTTPRoute `json:"http,omitempty"`
}

type HTTPRoute struct {
	// 匹配规则
	// Match []*HTTPMatchRequest `protobuf:"bytes,1,rep,name=match,proto3" json:"match,omitempty"`
	Match interface{} `protobuf:"bytes,1,rep,name=match,proto3" json:"match,omitempty"`

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

type HTTPRouteDestination struct {
	//目的服务
	Destination *Destination `protobuf:"bytes,1,opt,name=destination,proto3" json:"destination,omitempty"`
	//流量权重 (0-100)
	Weight int32 `protobuf:"varint,2,opt,name=weight,proto3" json:"weight,omitempty"`
	// Headers              *Headers `protobuf:"bytes,7,opt,name=headers,proto3" json:"headers,omitempty"`

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
	Delay *struct {
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
	Abort *struct {
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
