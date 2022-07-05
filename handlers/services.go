package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/util"
)

func ServiceInject(w http.ResponseWriter, r *http.Request) {
	cluster := r.URL.Query().Get("cluster")
	if _, ok := business.ClusterMap[cluster]; !ok {
		RespondWithJSON(w, http.StatusOK, "")
		return
	}

	params := mux.Vars(r)
	service := params["service"]
	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Workloads initialization error: "+err.Error())
		return
	}

	namespace := params["namespace"]
	// workloadType := "Deployment"

	serviceDetails, err := business.Svc.GetService(cluster, namespace, service, defaultHealthRateInterval, util.Clock.Now())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
	}

	type InjectData struct {
		ProxyCPU         float64 `json:"proxyCPU"`
		ProxyCPULimit    float64 `json:"proxyCPULimit"`
		ProxyMemory      int     `json:"proxyMemory"`
		ProxyMemoryLimit int     `json:"proxyMemoryLimit"`
		Concurrency      int     `json:"concurrency"`
	}

	var data InjectData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}
	if data.ProxyCPU == 0 {
		data.ProxyCPU = 0.5
	}

	if data.ProxyCPULimit == 0 {
		data.ProxyCPULimit = 2
	}

	if data.ProxyMemory == 0 {
		data.ProxyMemory = 128
	}

	if data.ProxyMemoryLimit == 0 {
		data.ProxyMemoryLimit = 1024
	}
	if data.Concurrency == 0 {
		data.Concurrency = 2
	}

	jsonPatch := string(`{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject":"true","sidecar.istio.io/concurrency":"%d","sidecar.istio.io/proxyCPU":"%dm","sidecar.istio.io/proxyCPULimit":"%dm","sidecar.istio.io/proxyMemory":"%dMi","sidecar.istio.io/proxyMemoryLimit":"%dMi"}}}}}`)

	jsonPatch = fmt.Sprintf(jsonPatch, data.Concurrency, int(data.ProxyCPU*1000), int(data.ProxyCPULimit*1000), data.ProxyMemory, data.ProxyMemoryLimit)

	// jsonPatch := string(`{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject":"true","sidecar.istio.io/concurrency":"2","sidecar.istio.io/proxyCPU":"1000m","sidecar.istio.io/proxyCPULimit":"1000m","sidecar.istio.io/proxyMemory":"256Mi","sidecar.istio.io/proxyMemoryLimit":"256Mi"}}}}}`)
	for _, deploy := range serviceDetails.Workloads {
		err1 := business.Workload.UpdateRemoteWorkload(cluster, namespace, deploy.Name, deploy.Type, true, jsonPatch)
		if err1 != nil {
			err = err1
		}
	}

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	audit(r, "UPDATE on Namespace: "+namespace+" Service name: "+service+" Patch: "+jsonPatch)
	RespondWithJSON(w, http.StatusOK, "")
}

func ServiceUnInject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cluster := r.URL.Query().Get("cluster")
	if _, ok := business.ClusterMap[cluster]; !ok {
		RespondWithJSON(w, http.StatusOK, "")
		return
	}

	service := params["service"]
	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Workloads initialization error: "+err.Error())
		return
	}

	namespace := params["namespace"]
	// workloadType := "Deployment"

	serviceDetails, err := business.Svc.GetService(cluster, namespace, service, defaultHealthRateInterval, util.Clock.Now())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	// fmt.Println(cluster, namespace, service, defaultHealthRateInterval, util.Clock.Now())
	// fmt.Println(serviceDetails.Workloads)
	jsonPatch := string(`{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject":"false"}}}}}`)
	for _, deploy := range serviceDetails.Workloads {
		// fmt.Printf("deploy:xxx:%+v\n", deploy)
		err1 := business.Workload.UpdateRemoteWorkload(cluster, namespace, deploy.Name, deploy.Type, true, jsonPatch)
		if err1 != nil {
			err = err1
		}
	}

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	audit(r, "UPDATE on Namespace: "+namespace+" Service name: "+service+" Patch: "+jsonPatch)
	RespondWithJSON(w, http.StatusOK, "")
}

// ServiceList is the API handler to fetch the list of services in a given namespace
func ServiceList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]

	// page := r.URL.Query().Get("page")
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "page error "+r.URL.Query().Get("page"))
		return
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "limit error "+r.URL.Query().Get("limit"))
		return
	}

	cluster := r.URL.Query().Get("cluster")
	cluster = strings.ToLower(cluster)
	resp := RespList{
		Data: []interface{}{},
	}
	if _, ok := business.ClusterMap[cluster]; !ok {
		RespondWithJSON(w, http.StatusOK, resp)
		return
	}

	// Get business layer
	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	// Fetch and build services
	serviceList, err := business.Svc.GetServiceList(cluster, namespace, false)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	searchName := r.URL.Query().Get("name")
	if len(searchName) > 0 {
		tmp := []models.ServiceOverview{}
		for i := range serviceList.Services {
			if strings.Contains(serviceList.Services[i].Name, searchName) {
				tmp = append(tmp, serviceList.Services[i])
			}
		}
		serviceList.Services = tmp
	}
	sort.SliceStable(serviceList.Services, func(i, j int) bool {
		return serviceList.Services[i].Name < serviceList.Services[j].Name
	})

	resp.CurrentPage = page
	resp.TotalCount = len(serviceList.Services)

	start, end, pageCount := SlicePage(page, limit, resp.TotalCount)

	resp.PageCount = pageCount
	resp.PageSize = limit
	serviceList.Services = serviceList.Services[start:end]

	for i := range serviceList.Services {
		svc := serviceList.Services[i]
		istioConfigStatus := map[string]bool{
			"ratelimit":         false,
			"mirror":            false,
			"localityLbsetting": false,
			"slowStart":         false,
			"outlierDetection":  false,
		}
		result, err := business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.DestinationRules, svc.Name)
		if err == nil {
			log.Infof("GetIstioConfigDetails:[%+v]", result)
			if result.DestinationRule != nil {
				if tp, ok := result.DestinationRule.Spec.TrafficPolicy.(map[string]interface{}); ok {
					if tp["outlierDetection"] != nil {
						istioConfigStatus["outlierDetection"] = true
					}

					if tp["loadBalancer"] != nil {
						if lb, ok := tp["loadBalancer"].(map[string]interface{}); ok {
							if lb["localityLbSetting"] != nil {
								istioConfigStatus["localityLbsetting"] = true
							}

							if lb["warmupDurationSecs"] != nil {
								istioConfigStatus["slowStart"] = true
							}
						}
					}
				}

			}

		}

		name := fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, svc.Name)
		result, err = business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, name)
		if err == nil {
			if result.EnvoyFilter != nil {
				istioConfigStatus["ratelimit"] = true
			}
		}

		mirrorLabel := geneMirrorLabelV1(svc.Name, namespace)
		selectLabel := "mirror=" + mirrorLabel
		allNs, _ := business.Namespace.GetNamespaces()
		for _, ns := range allNs {
			objects, err := business.IstioConfig.GetIstioObject(ns.Name, kubernetes.EnvoyFilters, selectLabel)
			if err == nil && len(objects) > 0 {
				istioConfigStatus["mirror"] = true
				break
			}
		}

		// mirrorName := geneMirrorName(svc.Name, namespace)
		// result, err = business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, mirrorName)
		// if err == nil && result.EnvoyFilter != nil {
		// 	istioConfigStatus["mirror"] = true
		// } else {
		// 	result, err = business.IstioConfig.GetIstioConfigDetails(ISTIO_SYSTEM_NAMESPACE, kubernetes.EnvoyFilters, mirrorName)
		// 	if err == nil && result.EnvoyFilter != nil {
		// 		istioConfigStatus["mirror"] = true
		// 	}
		// }

		serviceList.Services[i].IstioConfigStatus = istioConfigStatus

	}

	resp.Data = serviceList

	RespondWithJSON(w, http.StatusOK, resp)
}

// ServiceDetails is the API handler to fetch full details of an specific service
func ServiceDetails(w http.ResponseWriter, r *http.Request) {
	// cluster := r.URL.Query().Get("cluster")
	// if _, ok := business.ClusterMap[cluster]; !ok {
	// 	RespondWithJSON(w, http.StatusOK, "")
	// 	return
	// }

	idStr := r.Header.Get("cid")
	cid, _ := strconv.Atoi(idStr)
	if cid <= 0 {
		RespondWithError(w, http.StatusBadRequest, idStr)
		return
	}
	cluster, err := GetRemoteClustersByCid(cid)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Get business layer
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	// Rate interval is needed to fetch request rates based health
	queryParams := r.URL.Query()
	rateInterval := queryParams.Get("rateInterval")
	if rateInterval == "" {
		rateInterval = defaultHealthRateInterval
	}

	includeValidations := false
	if _, found := queryParams["validate"]; found {
		includeValidations = true
	}

	params := mux.Vars(r)
	namespace := params["namespace"]
	service := params["service"]
	queryTime := util.Clock.Now()
	rateInterval, err = adjustRateInterval(business, namespace, rateInterval, queryTime)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Adjust rate interval error: "+err.Error())
		return
	}

	var istioConfigValidations = models.IstioValidations{}
	var errValidations error

	wg := sync.WaitGroup{}
	if includeValidations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			istioConfigValidations, errValidations = business.Validations.GetValidations(namespace, service)
		}()
	}

	serviceDetails, err := business.Svc.GetService(cluster, namespace, service, rateInterval, queryTime)
	if includeValidations && err == nil {
		wg.Wait()
		serviceDetails.Validations = istioConfigValidations
		err = errValidations
	}

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, serviceDetails)
}

func ServiceUpdate(w http.ResponseWriter, r *http.Request) {
	// Get business layer
	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	// Rate interval is needed to fetch request rates based health
	queryParams := r.URL.Query()
	rateInterval := queryParams.Get("rateInterval")
	if rateInterval == "" {
		rateInterval = defaultHealthRateInterval
	}

	includeValidations := false
	if _, found := queryParams["validate"]; found {
		includeValidations = true
	}

	params := mux.Vars(r)
	namespace := params["namespace"]
	service := params["service"]
	queryTime := util.Clock.Now()
	rateInterval, err = adjustRateInterval(business, namespace, rateInterval, queryTime)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Adjust rate interval error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Update request with bad update patch: "+err.Error())
	}
	jsonPatch := string(body)
	var istioConfigValidations = models.IstioValidations{}
	var errValidations error

	wg := sync.WaitGroup{}
	if includeValidations {
		wg.Add(1)
		go func() {
			defer wg.Done()
			istioConfigValidations, errValidations = business.Validations.GetValidations(namespace, service)
		}()
	}

	serviceDetails, err := business.Svc.UpdateService(namespace, service, rateInterval, queryTime, jsonPatch)

	if includeValidations && err == nil {
		wg.Wait()
		serviceDetails.Validations = istioConfigValidations
		err = errValidations
	}

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	audit(r, "UPDATE on Namespace: "+namespace+" Service name: "+service+" Patch: "+jsonPatch)
	RespondWithJSON(w, http.StatusOK, serviceDetails)
}

func IstioMirrorClientDetail(w http.ResponseWriter, r *http.Request) {

	// params := mux.Vars(r)
	// namespace := params["namespace"]
	// object := params["object"]

	idStr := r.Header.Get("cid")
	cid, _ := strconv.Atoi(idStr)

	istioPrimaryRemote := business.IstioPrimary

	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	resp := map[string]interface{}{
		"cluster":   nil,
		"namespace": nil,
		"service":   nil,
	}

	hulkCluster, err := GetHulkClusters()
	if err != nil {
		log.Errorf("GetHulkClusters err:[%s]", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	getType := r.URL.Query().Get("type")
	if getType == "cluster" {

		cRes := []string{}

		var svcCluster string
		for _, c := range hulkCluster.Result {
			if c.Id == cid {
				svcCluster = c.Name
				break
			}
		}

		svcCluster = strings.ToLower(svcCluster)
		primaryCluster, ok := istioPrimaryRemote[svcCluster]
		if ok {
			tmp := []string{}

			for k, v := range istioPrimaryRemote {
				if primaryCluster == v {
					tmp = append(tmp, k)
				}
			}

			for _, c := range tmp {
				for _, hc := range hulkCluster.Result {
					if c == strings.ToLower(hc.Name) {
						cRes = append(cRes, hc.Name)
						break
					}
				}
			}

			resp["cluster"] = cRes
		}

	} else if getType == "namespace" {
		c := r.URL.Query().Get("clientCluster")
		c = strings.ToLower(c)
		allNs, _ := business.Namespace.GetRemoteNamespaces(c)
		nsRes := []string{}
		for _, ns := range allNs {
			nsRes = append(nsRes, ns.Name)
		}
		resp["namespace"] = nsRes

	} else if getType == "service" {

		c := r.URL.Query().Get("clientCluster")
		c = strings.ToLower(c)
		ns := r.URL.Query().Get("clientNamespace")

		sList, err := business.Svc.GetServiceList(c, ns, false)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "GetServiceList: "+err.Error())
			return
		}

		svcRes := []string{}
		for _, s := range sList.Services {
			if s.Name == "kubernetes" {
				continue
			}
			if s.IstioSidecar {
				svcRes = append(svcRes, s.Name)
			}
		}

		resp["service"] = svcRes

	}

	RespondWithJSON(w, http.StatusOK, resp)

}
