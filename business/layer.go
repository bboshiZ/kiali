package business

import (
	"crypto/md5"
	"errors"
	"sync"

	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/jaeger"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/kubernetes/cache"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/prometheus"
)

// Layer is a container for fast access to inner services
type Layer struct {
	App            AppService
	Health         HealthService
	IstioConfig    IstioConfigService
	IstioStatus    IstioStatusService
	Iter8          Iter8Service
	Jaeger         JaegerService
	k8s            kubernetes.ClientInterface
	Mesh           MeshService
	Namespace      NamespaceService
	OpenshiftOAuth OpenshiftOAuthService
	ProxyStatus    ProxyStatusService
	RegistryStatus RegistryStatusService
	Svc            SvcService
	TLS            TLSService
	TokenReview    TokenReviewService
	Validations    IstioValidationsService
	Workload       WorkloadService
}

// Global clientfactory and prometheus clients.
var clientFactory kubernetes.ClientFactory
var clientFactoryMap = map[string]kubernetes.ClientFactory{}

var prometheusClient prometheus.ClientInterface
var once sync.Once
var kialiCache cache.KialiCache
var kialiRemoteCache = map[string]cache.KialiCache{}

func initKialiCache(remoteClusters []string) {
	if config.Get().KubernetesConfig.CacheEnabled {
		if cache, err := cache.NewKialiCache(""); err != nil {
			log.Errorf("Error initializing Kiali Cache. Details: %s", err)
		} else {
			kialiCache = cache
		}

		for _, c := range remoteClusters {
			if cache, err := cache.NewKialiCache(c); err != nil {
				log.Errorf("Error initializing Kiali Cache. Details: %s", err)
			} else {
				kialiRemoteCache[c] = cache
			}
		}

	}

	if excludedWorkloads == nil {
		excludedWorkloads = make(map[string]bool)
		for _, w := range config.Get().KubernetesConfig.ExcludeWorkloads {
			excludedWorkloads[w] = true
		}
	}
}

func IsNamespaceCached(cluster, namespace string) bool {
	if cluster == "" {
		return false
	}
	if cluster == "" {
		ok := kialiCache != nil && kialiCache.CheckNamespace(namespace)
		return ok
	}
	ok := kialiRemoteCache[cluster] != nil && kialiRemoteCache[cluster].CheckNamespace(namespace)

	return ok
}

func IsResourceCached(cluster, namespace string, resource string) bool {
	ok := IsNamespaceCached(cluster, namespace)
	if ok && resource != "" {
		if cluster == "" {
			ok = kialiCache.CheckIstioResource(resource)
			return ok
		}
		ok = kialiRemoteCache[cluster].CheckIstioResource(resource)

	}
	return ok
}

func getTokenHash(authInfo *api.AuthInfo) string {
	tokenData := authInfo.Token

	if authInfo.Impersonate != "" {
		tokenData += authInfo.Impersonate
	}

	if authInfo.ImpersonateGroups != nil {
		for _, group := range authInfo.ImpersonateGroups {
			tokenData += group
		}
	}

	if authInfo.ImpersonateUserExtra != nil {
		for key, element := range authInfo.ImpersonateUserExtra {
			for _, userExtra := range element {
				tokenData += key + userExtra
			}
		}

	}

	h := md5.New()
	_, err := h.Write([]byte(tokenData))
	if err != nil {
		// errcheck linter want us to check for the error returned by h.Write.
		// However, docs of md5 say that this Writer never returns an error.
		// See: https://golang.org/pkg/hash/#Hash
		// So, let's check the error, and panic. Per the docs, this panic should
		// never be reached.
		panic("md5.Write returned error.")
	}
	return string(h.Sum(nil))

}

// Get the business.Layer
func GetByCluster(cluster string) (*Layer, error) {

	if len(clientFactoryMap) == 0 {
		userClient, err := kubernetes.GetClientFactoryMap()
		if err != nil {
			return nil, err
		}
		for c, _ := range userClient {
			clientFactoryMap[c] = userClient[c]

		}

	}

	if clientFac, ok := clientFactoryMap[cluster]; ok {
		// Creates a new k8s client based on the current users token

		token := kubernetes.IstioPrimaryResctConfig[cluster].BearerToken
		k8s, err := clientFac.GetClient(&api.AuthInfo{Token: token})
		if err != nil {
			return nil, err
		}

		// Use an existing Prometheus client if it exists, otherwise create and use in the future
		if prometheusClient == nil {
			prom, err := prometheus.NewClient()
			if err != nil {
				return nil, err
			}
			prometheusClient = prom
		}

		// Create Jaeger client
		jaegerLoader := func() (jaeger.ClientInterface, error) {
			return jaeger.NewClient("")
		}

		return NewWithBackends(k8s, prometheusClient, jaegerLoader), nil

	} else {
		return nil, errors.New("clientFactory not found")
	}

}

// Get the business.Layer
func Get(authInfo *api.AuthInfo) (*Layer, error) {
	// Kiali Cache will be initialized once at first use of Business layer
	// once.Do(initKialiCache)

	// Use an existing client factory if it exists, otherwise create and use in the future
	if clientFactory == nil {
		userClient, err := kubernetes.GetClientFactory()
		if err != nil {
			return nil, err
		}
		clientFactory = userClient
	}
	// Creates a new k8s client based on the current users token
	k8s, err := clientFactory.GetClient(authInfo)
	if err != nil {
		return nil, err
	}

	// Use an existing Prometheus client if it exists, otherwise create and use in the future
	if prometheusClient == nil {
		prom, err := prometheus.NewClient()
		if err != nil {
			return nil, err
		}
		prometheusClient = prom
	}

	// Create Jaeger client
	jaegerLoader := func() (jaeger.ClientInterface, error) {
		return jaeger.NewClient(authInfo.Token)
	}

	return NewWithBackends(k8s, prometheusClient, jaegerLoader), nil
}

// SetWithBackends allows for specifying the ClientFactory and Prometheus clients to be used.
// Mock friendly. Used only with tests.
func SetWithBackends(cf kubernetes.ClientFactory, prom prometheus.ClientInterface) {
	clientFactory = cf
	prometheusClient = prom
}

// NewWithBackends creates the business layer using the passed k8s and prom clients
func NewWithBackends(k8s kubernetes.ClientInterface, prom prometheus.ClientInterface, jaegerClient JaegerLoader) *Layer {
	temporaryLayer := &Layer{}
	temporaryLayer.App = AppService{prom: prom, k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.Health = HealthService{prom: prom, k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.IstioConfig = IstioConfigService{k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.IstioStatus = IstioStatusService{k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.Iter8 = Iter8Service{k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.Jaeger = JaegerService{loader: jaegerClient, businessLayer: temporaryLayer}
	temporaryLayer.k8s = k8s
	temporaryLayer.Mesh = NewMeshService(k8s, temporaryLayer, nil)
	temporaryLayer.Namespace = NewNamespaceService(k8s)
	temporaryLayer.OpenshiftOAuth = OpenshiftOAuthService{k8s: k8s}
	temporaryLayer.ProxyStatus = ProxyStatusService{k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.RegistryStatus = RegistryStatusService{k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.Svc = SvcService{prom: prom, k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.TLS = TLSService{k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.TokenReview = NewTokenReview(k8s)
	temporaryLayer.Validations = IstioValidationsService{k8s: k8s, businessLayer: temporaryLayer}
	temporaryLayer.Workload = WorkloadService{k8s: k8s, prom: prom, businessLayer: temporaryLayer}

	return temporaryLayer
}

func Stop() {
	if kialiCache != nil {
		kialiCache.Stop()
	}
}
