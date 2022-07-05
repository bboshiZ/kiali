package handlers

import (
	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/graph/telemetry/istio"
)

func InitRemoteCluster() {
	istio.RemoteCluster = business.InitRemoteCluster()

	// saToken, _ := kubernetes.GetKialiToken()
	// authInfo := &api.AuthInfo{Token: saToken}
	// business, _ := business.Get(authInfo)
	// business.Mesh.InitRemoteCluster()

	// r, _ := http.NewRequest("GET", "https://www.baidu.com", nil)
	// business, err := getBusiness(r)

}
