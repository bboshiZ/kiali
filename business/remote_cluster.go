package business

import (
	"fmt"
	"time"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/kubernetes/cache"
	"github.com/kiali/kiali/log"
	"k8s.io/client-go/tools/clientcmd/api"
)

type RemoteCluster struct {
	// ApiEndpoint is the URL where the Kubernetes/Cluster API Server can be contacted
	ApiEndpoint string `json:"apiEndpoint"`

	// IsKialiHome specifies if this cluster is hosting this Kiali instance (and the observed Mesh Control Plane)
	IsKialiHome bool `json:"isKialiHome"`

	// KialiInstances is the list of Kialis discovered in the cluster.
	KialiInstances []KialiInstance `json:"kialiInstances"`

	// Name specifies the CLUSTER_ID as known by the Control Plane
	Name string `json:"name"`

	// Network specifies the logical NETWORK_ID as known by the Control Plane
	Network string `json:"network"`

	// SecretName is the name of the kubernetes "remote secret" where data of this cluster was resolved
	SecretName string `json:"secretName"`

	K8s kubernetes.ClientInterface
}

var (
	remoteIstioClusters = map[string]RemoteCluster{}
	defaultClusterID    = "shareit-cce-test"
)

func InitRemoteCluster() {
	saToken, _ := kubernetes.GetKialiToken()
	authInfo := &api.AuthInfo{Token: saToken}
	business, _ := Get(authInfo)
	in := business.Mesh
	// business.Mesh.InitRemoteCluster()

	conf := config.Get()
	secrets, err := in.k8s.GetSecrets(conf.IstioNamespace, "istio/remoteKiali=true")
	if err != nil {
		return
	}

	if len(secrets) == 0 {
		return
	}

	// clusters := make([]RemoteCluster, 0, len(secrets))

	// Inspect the secret to extract the cluster_id and api_endpoint of each remote cluster.
	for _, secret := range secrets {
		clusterName, ok := secret.Annotations["networking.istio.io/cluster"]
		if !ok {
			clusterName = DefaultClusterID
		}

		kubeconfigFile, ok := secret.Data[clusterName]
		if !ok {
			continue
		}

		remoteSecret, parseErr := kubernetes.ParseRemoteSecretBytes(kubeconfigFile)
		if parseErr != nil {
			continue
		}

		if len(remoteSecret.Clusters) != 1 {
			continue
		}

		restConfig, restConfigErr := kubernetes.UseRemoteCreds(remoteSecret)
		if restConfigErr != nil {
			log.Errorf("Error using remote creds: %v", restConfigErr)
			continue
		}

		restConfig.Timeout = 15 * time.Second
		restConfig.BearerToken = remoteSecret.Users[0].User.Token

		// kubeconfig := "/root/code/kkk/kiali/kubeconfig"
		// restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)

		// remoteClientSet, clientSetErr := clientFactory.GetClient(&api.AuthInfo{Token: remoteSecret.Users[0].User.Token})

		remoteClientSet, clientSetErr := kubernetes.NewClientFromConfig(restConfig)
		// remoteClientSet, clientSetErr := in.newRemoteClient(restConfig)
		if clientSetErr != nil {
			log.Errorf("Error creating client set: %v", clientSetErr)
			continue
		}

		// res, err := remoteClientSet.ForwardGetRequest("sample", "sleep-754d6588c4-cc2vz", 1234, 8080, "/")
		// log.Errorf("Error creating client set: %v,%v", string(res), err)

		meshCluster := RemoteCluster{
			Name:        clusterName,
			SecretName:  secret.Name,
			ApiEndpoint: remoteSecret.Clusters[0].Cluster.Server,
			K8s:         remoteClientSet,
		}
		k8sClient, _ := kubernetes.NewClientFromConfig(restConfig)
		cache.RemoteK8S = k8sClient

		networkName := in.resolveNetwork(clusterName, remoteSecret)
		if len(networkName) != 0 {
			meshCluster.Network = networkName
		}

		meshCluster.KialiInstances = in.findRemoteKiali(clusterName, remoteSecret)

		defaultClusterID = secret.Annotations["istio/remoteCluster"]
		remoteIstioClusters[secret.Annotations["istio/remoteCluster"]] = meshCluster
	}

	once.Do(initKialiCache)

	for k, c := range remoteIstioClusters {
		fmt.Printf("new remote k8s cluster:%s,%+v\n", k, c)
	}

}
