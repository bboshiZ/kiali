package routing

import (
	"net/http"

	"github.com/kiali/kiali/handlers"
)

// Route describes a single route
type Route struct {
	Name          string
	Method        string
	Pattern       string
	HandlerFunc   http.HandlerFunc
	Authenticated bool
}

// Routes holds an array of Route. A note on swagger documentation. The path variables and query parameters
// are defined in ../doc.go.  YOu need to manually associate params and routes.
type Routes struct {
	Routes []Route
}

// NewRoutes creates and returns all the API routes
func NewRoutes() (r *Routes) {
	r = new(Routes)

	r.Routes = []Route{

		{
			"IstioNetworkConfigUpdate",
			"POST",
			"/api/namespaces/{namespace}/{object_type}/{network_type}/{object}",
			handlers.IstioNetworkConfigUpdate,
			true,
		},

		{
			"IstioNetworkConfigDelete",
			"DELETE",
			"/api/namespaces/{namespace}/{object_type}/{network_type}/{object}",
			handlers.IstioNetworkConfigDelete,
			true,
		},

		{
			"IstioConfigDetails",
			"GET",
			"/api/namespaces/{namespace}/{object_type}/{object}",
			handlers.IstioConfigDetailsV2,
			true,
		},

		{
			"IstioMirrorClusterList",
			"GET",
			"/api/mirrorclient/namespaces/{namespace}/services/{object}",
			handlers.IstioMirrorClientDetail,
			true,
		},

		{
			"ClusterList",
			"GET",
			"/api/clusters",
			handlers.ClusterList,
			true,
		},

		{
			"MeshClusterList",
			"GET",
			"/api/mesh/clusters",
			handlers.MeshClusterList,
			true,
		},

		{
			"NamespaceList",
			"GET",
			"/api/namespaces",
			handlers.NamespaceList,
			true,
		},

		{
			"NamespaceUpdate",
			"PATCH",
			"/api/namespaces/{namespace}",
			handlers.NamespaceUpdate,
			true,
		},

		{
			"LocalityList",
			"GET",
			"/api/localities",
			handlers.LocalityList,
			true,
		},

		{
			"ServiceList",
			"GET",
			"/api/namespaces/{namespace}/services",
			handlers.ServiceList,
			true,
		},

		{
			"ServiceDetails",
			"GET",
			"/api/namespaces/{namespace}/services/{service}",
			handlers.ServiceDetails,
			true,
		},

		{
			"ServiceInject",
			"POST",
			"/api/namespaces/{namespace}/services/{service}/inject",
			handlers.ServiceInject,
			true,
		},

		{
			"ServiceDisInject",
			"POST",
			"/api/namespaces/{namespace}/services/{service}/unInject",
			handlers.ServiceUnInject,
			true,
		},

		{
			"IstioConfigCreate",
			"POST",
			"/api/namespaces/{namespace}/{object_type}",
			handlers.IstioConfigCreate,
			true,
		},
		{
			"IstioConfigDelete",
			"DELETE",
			"/api/namespaces/{namespace}/{object_type}/{object}",
			handlers.IstioConfigDelete,
			true,
		},

		{
			"IstioConfigUpdate",
			"PUT",
			"/api/namespaces/{namespace}/{object_type}/{object}",
			handlers.IstioConfigUpdate,
			true,
		},

		{
			"IstioConfigList",
			"GET",
			"/api/namespaces/{namespace}/config",
			handlers.IstioConfigList,
			true,
		},

		{
			"IstioConfigDetails",
			"GET",
			"/api/namespaces/{namespace}/istio/{object_type}/{object}",
			handlers.IstioConfigDetails,
			true,
		},

		{
			"ServiceTraces",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/traces",
			handlers.ServiceTraces,
			true,
		},

		{
			"GraphService",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/graph",
			handlers.GraphService,
			true,
		},

		{
			"Healthz",
			"GET",
			"/healthz",
			handlers.Healthz,
			false,
		},
		{
			"Root",
			"GET",
			"/api",
			handlers.Root,
			false,
		},

		{
			"Authenticate",
			"GET",
			"/api/authenticate",
			handlers.Authenticate,
			false,
		},

		{
			"OpenshiftCheckToken",
			"POST",
			"/api/authenticate",
			handlers.Authenticate,
			false,
		},

		{
			"Logout",
			"GET",
			"/api/logout",
			handlers.Logout,
			false,
		},

		{
			"AuthenticationInfo",
			"GET",
			"/api/auth/info",
			handlers.AuthenticationInfo,
			false,
		},

		{
			"AuthenticationInfo",
			"GET",
			"/api/auth/openid_redirect",
			handlers.OpenIdRedirect,
			false,
		},

		{
			"Status",
			"GET",
			"/api/status",
			handlers.Root,
			true,
		},

		{
			"Config",
			"GET",
			"/api/config",
			handlers.Config,
			true,
		},

		{
			"IstioConfigPermissions",
			"GET",
			"/api/istio/permissions",
			handlers.IstioConfigPermissions,
			true,
		},

		{
			"ServiceUpdate",
			"PATCH",
			"/api/namespaces/{namespace}/services/{service}",
			handlers.ServiceUpdate,
			true,
		},

		{
			"AppSpans",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/spans",
			handlers.AppSpans,
			true,
		},

		{
			"WorkloadSpans",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/spans",
			handlers.WorkloadSpans,
			true,
		},

		{
			"ServiceSpans",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/spans",
			handlers.ServiceSpans,
			true,
		},

		{
			"AppTraces",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/traces",
			handlers.AppTraces,
			true,
		},

		{
			"WorkloadTraces",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/traces",
			handlers.WorkloadTraces,
			true,
		},

		{
			"ErrorTraces",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/errortraces",
			handlers.ErrorTraces,
			true,
		},

		{
			"TracesDetails",
			"GET",
			"/api/traces/{traceID}",
			handlers.TraceDetails,
			true,
		},

		{
			"WorkloadList",
			"GET",
			"/api/namespaces/{namespace}/workloads",
			handlers.WorkloadList,
			true,
		},

		{
			"WorkloadDetails",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}",
			handlers.WorkloadDetails,
			true,
		},

		//
		{
			"WorkloadUpdate",
			"PATCH",
			"/api/namespaces/{namespace}/workloads/{workload}",
			handlers.WorkloadUpdate,
			true,
		},

		//
		{
			"AppList",
			"GET",
			"/api/namespaces/{namespace}/apps",
			handlers.AppList,
			true,
		},

		{
			"AppDetails",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}",
			handlers.AppDetails,
			true,
		},

		{
			"NamespaceUpdate",
			"PATCH",
			"/api/namespaces/{namespace}",
			handlers.NamespaceUpdate,
			true,
		},

		{
			"ServiceMetrics",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/metrics",
			handlers.ServiceMetrics,
			true,
		},

		{
			"AggregateMetrics",
			"GET",
			"/api/namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/metrics",
			handlers.AggregateMetrics,
			true,
		},

		{
			"AppMetrics",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/metrics",
			handlers.AppMetrics,
			true,
		},

		{
			"WorkloadMetrics",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/metrics",
			handlers.WorkloadMetrics,
			true,
		},

		{
			"ServiceDashboard",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/dashboard",
			handlers.ServiceDashboard,
			true,
		},

		{
			"AppDashboard",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/dashboard",
			handlers.AppDashboard,
			true,
		},

		{
			"WorkloadDashboard",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/dashboard",
			handlers.WorkloadDashboard,
			true,
		},

		{
			"CustomDashboard",
			"GET",
			"/api/namespaces/{namespace}/customdashboard/{dashboard}",
			handlers.CustomDashboard,
			true,
		},

		{
			"ServiceHealth",
			"GET",
			"/api/namespaces/{namespace}/services/{service}/health",
			handlers.ServiceHealth,
			true,
		},

		{
			"AppHealth",
			"GET",
			"/api/namespaces/{namespace}/apps/{app}/health",
			handlers.AppHealth,
			true,
		},

		{
			"WorkloadHealth",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/health",
			handlers.WorkloadHealth,
			true,
		},

		{
			"NamespaceMetrics",
			"GET",
			"/api/namespaces/{namespace}/metrics",
			handlers.NamespaceMetrics,
			true,
		},

		{
			"NamespaceHealth",
			"GET",
			"/api/namespaces/{namespace}/health",
			handlers.NamespaceHealth,
			true,
		},

		{
			"NamespaceValidationSummary",
			"GET",
			"/api/namespaces/{namespace}/validations",
			handlers.NamespaceValidationSummary,
			true,
		},

		{
			"NamespaceTls",
			"GET",
			"/api/mesh/tls",
			handlers.MeshTls,
			true,
		},

		{
			"NamespaceTls",
			"GET",
			"/api/namespaces/{namespace}/tls",
			handlers.NamespaceTls,
			true,
		},

		{
			"IstioStatus",
			"GET",
			"/api/istio/status",
			handlers.IstioStatus,
			true,
		},

		{
			"GraphNamespaces",
			"GET",
			"/api/namespaces/graph",
			handlers.GraphNamespaces,
			true,
		},

		{

			"GraphAggregate",
			"GET",
			"/api/namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/graph",
			handlers.GraphNode,
			true,
		},

		{

			"GraphAggregateByService",
			"GET",
			"/api/namespaces/{namespace}/aggregates/{aggregate}/{aggregateValue}/{service}/graph",
			handlers.GraphNode,
			true,
		},

		{

			"GraphAppVersion",
			"GET",
			"/api/namespaces/{namespace}/applications/{app}/versions/{version}/graph",
			handlers.GraphNode,
			true,
		},

		{
			"GraphApp",
			"GET",
			"/api/namespaces/{namespace}/applications/{app}/graph",
			handlers.GraphNode,
			true,
		},

		{
			"GraphWorkload",
			"GET",
			"/api/namespaces/{namespace}/workloads/{workload}/graph",
			handlers.GraphNode,
			true,
		},

		{
			"GrafanaURL",
			"GET",
			"/api/grafana",
			handlers.GetGrafanaInfo,
			true,
		},

		{
			"JaegerURL",
			"GET",
			"/api/jaeger",
			handlers.GetJaegerInfo,
			true,
		},

		{
			"PodDetails",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}",
			handlers.PodDetails,
			true,
		},

		{
			"PodLogs",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}/logs",
			handlers.PodLogs,
			true,
		},

		{
			"PodConfigDump",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}/config_dump",
			handlers.ConfigDump,
			true,
		},

		{
			"PodConfigDump",
			"GET",
			"/api/namespaces/{namespace}/pods/{pod}/config_dump/{resource}",
			handlers.ConfigDumpResourceEntries,
			true,
		},

		{
			"Iter8Info",
			"GET",
			"/api/iter8",
			handlers.Iter8Status,
			true,
		},

		{
			"Iter8ExperimentsByNamespaceAndName",
			"GET",
			"/api/iter8/namespaces/{namespace}/experiments/{name}",
			handlers.Iter8ExperimentGet,
			true,
		},

		{
			"Iter8Experiments",
			"GET",
			"/api/iter8/experiments",
			handlers.Iter8Experiments,
			true,
		},

		{
			Name:          "Iter8ExperimentCreate",
			Method:        "POST",
			Pattern:       "/api/iter8/namespaces/{namespace}/experiments",
			HandlerFunc:   handlers.Iter8ExperimentCreate,
			Authenticated: true,
		},

		{
			Name:          "Iter8ExperimentsUpdate",
			Method:        "PATCH",
			Pattern:       "/api/iter8/namespaces/{namespace}/experiments/{name}",
			HandlerFunc:   handlers.Iter8ExperimentUpdate,
			Authenticated: true,
		},

		{
			Name:          "Iter8ExperimentDelete",
			Method:        "DELETE",
			Pattern:       "/api/iter8/namespaces/{namespace}/experiments/{name}",
			HandlerFunc:   handlers.Iter8ExperimentDelete,
			Authenticated: true,
		},

		{
			"Iter8Metrics",
			"GET",
			"/api/iter8/metrics",
			handlers.Iter8Metrics,
			true,
		},

		{
			"Iter8ExperimentGetYaml",
			"GET",
			"/api/iter8/namespaces/{namespace}/experiments/{name}/yaml",
			handlers.Iter8ExperimentGetYaml,
			true,
		},

		{
			Name:          "MetricsStats",
			Method:        "POST",
			Pattern:       "/api/stats/metrics",
			HandlerFunc:   handlers.MetricsStats,
			Authenticated: true,
		},
		{
			"GetClusters",
			"GET",
			"/api/clusters",
			handlers.GetClusters,
			true,
		},
	}

	return
}
