package handlers

import (
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
	jsonPatch := string(`{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject":"true"}}}}}`)
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
	jsonPatch := string(`{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject":"false"}}}}}`)
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

// ServiceList is the API handler to fetch the list of services in a given namespace
func ServiceList(w http.ResponseWriter, r *http.Request) {
	resp := RespList{
		Data: []interface{}{},
	}
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
	if _, ok := business.ClusterMap[cluster]; !ok {
		RespondWithJSON(w, http.StatusOK, resp)
		return
	}

	params := mux.Vars(r)
	namespace := params["namespace"]

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

		mirrorName := geneMirrorName(svc.Name, namespace)
		result, err = business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, mirrorName)
		if err == nil && result.EnvoyFilter != nil {
			istioConfigStatus["mirror"] = true
		} else {
			result, err = business.IstioConfig.GetIstioConfigDetails(ISTIO_SYSTEM_NAMESPACE, kubernetes.EnvoyFilters, mirrorName)
			if err == nil && result.EnvoyFilter != nil {
				istioConfigStatus["mirror"] = true
			}
		}

		serviceList.Services[i].IstioConfigStatus = istioConfigStatus

	}

	resp.Data = serviceList

	RespondWithJSON(w, http.StatusOK, resp)
}

// ServiceDetails is the API handler to fetch full details of an specific service
func ServiceDetails(w http.ResponseWriter, r *http.Request) {
	cluster := r.URL.Query().Get("cluster")
	if _, ok := business.ClusterMap[cluster]; !ok {
		RespondWithJSON(w, http.StatusOK, "")
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
