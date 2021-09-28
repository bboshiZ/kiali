package handlers

import "github.com/kiali/kiali/business"

func InitRemoteCluster() {
	business.InitRemoteCluster()
	// saToken, _ := kubernetes.GetKialiToken()
	// authInfo := &api.AuthInfo{Token: saToken}
	// business, _ := business.Get(authInfo)
	// business.Mesh.InitRemoteCluster()

	// r, _ := http.NewRequest("GET", "https://www.baidu.com", nil)
	// business, err := getBusiness(r)

}
