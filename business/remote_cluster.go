package business

import (
	"fmt"
	"time"

	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
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
	Token      string `json:"token"`
	K8s        kubernetes.ClientInterface
}

var (
	remoteIstioClusters = map[string]RemoteCluster{}
	defaultClusterID    = "shareit-cce-test"
	secretDataName      = "istio-config-manager"
	ClusterMap          = map[string]bool{}
	remoteClusters      []string

	IstioPrimary = map[string]string{}
)

func InitRemoteCluster() (clusterID string) {
	saToken, err := kubernetes.GetKialiToken()
	if err != nil {
		log.Errorf("Error GetKialiToken: %v", err)
		return
	}
	fmt.Printf("saToken:%+v\n", string(saToken))
	business, err := Get(&api.AuthInfo{Token: saToken})
	if err != nil {
		log.Errorf("Error Get business: %v", err)
		return
	}
	in := business.Mesh

	conf := config.Get()
	secrets, err := in.k8s.GetSecrets(conf.IstioNamespace, "istio-config-manager=true")
	if err != nil {
		log.Errorf("Error GetSecrets: %v", err)
		return
	}

	// fmt.Printf("secrets:%+v\n", secrets)

	if len(secrets) == 0 {
		return
	}

	for _, secret := range secrets {
		// log.Errorf("GetSecrets-11111: %+v", secret.Labels)
		// log.Errorf("GetSecrets-11111: %+v", (secret.Data))

		clusterName, ok := secret.Labels["istio-k8s-cluster"]
		if !ok {
			continue
		}

		kubeconfigFile, ok := secret.Data[secretDataName]
		// log.Errorf("GetSecrets-22: %+v", string(kubeconfigFile))
		// log.Errorf("GetSecrets-22: %+v", ok)
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

		IstioPrimary[clusterName] = secret.Labels["istio-primary-cluster"]

		restConfig, restConfigErr := kubernetes.UseRemoteCreds(remoteSecret)
		if restConfigErr != nil {
			log.Errorf("Error using remote creds: %v", restConfigErr)
			continue
		}

		if _, ok := secret.Labels["istio-primary"]; ok {
			kubernetes.IstioPrimaryResctConfig[clusterName] = restConfig
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
			Token:       remoteSecret.Users[0].User.Token,
		}

		networkName := in.resolveNetwork(clusterName, remoteSecret)
		if len(networkName) != 0 {
			meshCluster.Network = networkName
		}

		meshCluster.KialiInstances = in.findRemoteKiali(clusterName, remoteSecret)
		ClusterMap[clusterName] = true
		remoteIstioClusters[clusterName] = meshCluster

		// k8sClient, err := kubernetes.NewClientFromConfig(restConfig)
		// if err != nil {
		// 	log.Errorf("Error NewClientFromConfig: %v", err)
		// 	continue
		// }
		// cache.RemoteK8S = k8sClient
		// cache.RemoteClusters[remoteC] = k8sClient
		// remoteClusters = append(remoteClusters, remoteC)

	}

	initKialiCache(remoteClusters)

	for k, c := range remoteIstioClusters {
		fmt.Printf("new remote k8s cluster:%s,%+v\n", k, c)
	}

	fmt.Printf("new IstioPrimary:%s\n", IstioPrimary)

	return
}
