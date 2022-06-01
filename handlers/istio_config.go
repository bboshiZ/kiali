package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"unicode"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/graph/api"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/util"
	"github.com/shopspring/decimal"
)

func LocalityList(w http.ResponseWriter, r *http.Request) {

	hw := []string{"ap-southeast-3/*", "ap-southeast-3/ap-southeast-3a/*", "ap-southeast-3/ap-southeast-3b/*", "ap-southeast-3/ap-southeast-3c/*"}
	aws := []string{"ap-southeast-1/*", "ap-southeast-1/ap-southeast-1a/*", "ap-southeast-1/ap-southeast-1b/*", "ap-southeast-1/ap-southeast-1c/*"}
	resp := map[string]interface{}{
		"locality": append(hw, aws...),
		// "epDistribution": map[string]int{},
		// "serviceRegion": []string{},
	}
	// typs Locality struct{
	// 	Region string
	// 	Zone []string
	// }
	// hw:=Locality{
	// 	Region:"ap-southeast-3",
	// 	Zone:[]string{}
	// }
	// Get business layer

	meshCluster := business.ClusterMap
	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	q := r.URL.Query()
	// cluster := q.Get("cluster")
	namespace := q.Get("namespace")
	service := q.Get("service")

	dist := map[string]int{}

	for c, _ := range meshCluster {
		serviceDetails, err := business.Svc.GetService(c, namespace, service, defaultHealthRateInterval, util.Clock.Now())
		if err != nil {
			log.Errorf("GetService from [%s] err:[%s]", c, err)
			continue
		}
		nodeList, err := business.Svc.GetAllNode()
		if err != nil {
			log.Errorf("GetAllNode err:[%s]", err)
			continue
		}
		for _, ep := range serviceDetails.Endpoints {
			for _, dpAddress := range ep.Addresses {
				for _, cn := range nodeList {
					for _, n := range cn.Items {
						if dpAddress.NodeName == n.Name {
							region := n.Labels["topology.kubernetes.io/region"]
							if region != "" {
								dist[region+"/*"]++
							}
							zone := n.Labels["topology.kubernetes.io/zone"]
							if zone != "" {
								dist[region+"/"+zone]++
							}
							// resp["serviceRegion"] = append(resp["serviceRegion"], n.Labels["topology.kubernetes.io/region"])
						}
					}

				}

			}

		}
	}
	// serviceDetails, err := business.Svc.GetService(cluster, namespace, service, defaultHealthRateInterval, queryTime)
	// nodeList, err := business.Svc.GetAllNode()

	// dist := map[string]int{}
	// for _, ep := range serviceDetails.Endpoints {
	// 	for _, dpAddress := range ep.Addresses {
	// 		for _, cn := range nodeList {
	// 			for _, n := range cn.Items {
	// 				if dpAddress.NodeName == n.Name {
	// 					region := n.Labels["topology.kubernetes.io/region"]
	// 					if region != "" {
	// 						dist[region+"/*"]++
	// 					}
	// 					zone := n.Labels["topology.kubernetes.io/zone"]
	// 					if zone != "" {
	// 						dist[region+"/"+zone]++
	// 					}
	// 					// resp["serviceRegion"] = append(resp["serviceRegion"], n.Labels["topology.kubernetes.io/region"])
	// 				}
	// 			}

	// 		}

	// 	}

	// }

	resp["distribution"] = dist

	// serviceList, err := business.Svc.GetNodeList()

	RespondWithJSON(w, http.StatusOK, resp)

}

func IstioConfigList(w http.ResponseWriter, r *http.Request) {
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

	params := mux.Vars(r)
	namespace := params["namespace"]
	query := r.URL.Query()
	objects := ""
	parsedTypes := make([]string, 0)
	if _, ok := query["objects"]; ok {
		objects = strings.ToLower(query.Get("objects"))
		if len(objects) > 0 {
			parsedTypes = strings.Split(objects, ",")
		}
	}

	includeValidations := false
	labelSelector := ""
	if _, found := query["labelSelector"]; found {
		labelSelector = query.Get("labelSelector")
	}

	workloadSelector := ""
	if _, found := query["workloadSelector"]; found {
		workloadSelector = query.Get("workloadSelector")
	}

	criteria := business.ParseIstioConfigCriteria(namespace, objects, labelSelector, workloadSelector)

	// Get business layer
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	var istioConfigValidations models.IstioValidations

	wg := sync.WaitGroup{}
	if includeValidations {
		wg.Add(1)
		go func(namespace string, istioConfigValidations *models.IstioValidations, err *error) {
			defer wg.Done()
			// We don't filter by objects when calling validations, because certain validations require fetching all types to get the correct errors
			istioConfigValidationResults, errValidations := business.Validations.GetValidations(namespace, "")
			if errValidations != nil && *err == nil {
				*err = errValidations
			} else {
				if len(parsedTypes) > 0 {
					istioConfigValidationResults = istioConfigValidationResults.FilterByTypes(parsedTypes)
				}
				*istioConfigValidations = istioConfigValidationResults
			}
		}(namespace, &istioConfigValidations, &err)
	}

	istioConfig, err := business.IstioConfig.GetIstioConfigList(criteria)
	if includeValidations {
		// Add validation results to the IstioConfigList once they're available (previously done in the UI layer)
		wg.Wait()
		istioConfig.IstioValidations = istioConfigValidations
	}

	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	resp := RespList{
		TotalCount:  10,
		PageCount:   1,
		CurrentPage: 1,
		PageSize:    10,
		Data:        istioConfig,
	}

	searchName := r.URL.Query().Get("name")
	if len(searchName) > 0 {
		if objects == "virtualservices" {
			tmp := []models.VirtualService{}
			for i := range istioConfig.VirtualServices.Items {
				if strings.HasPrefix(istioConfig.VirtualServices.Items[i].IstioBase.Metadata.Name, searchName) {
					tmp = append(tmp, istioConfig.VirtualServices.Items[i])
				}
			}
			istioConfig.VirtualServices.Items = tmp
		} else if objects == "destinationrules" {
			tmp := []models.DestinationRule{}
			for i := range istioConfig.DestinationRules.Items {
				if strings.HasPrefix(istioConfig.DestinationRules.Items[i].IstioBase.Metadata.Name, searchName) {
					tmp = append(tmp, istioConfig.DestinationRules.Items[i])
				}
			}
			istioConfig.DestinationRules.Items = tmp
		}

	}
	if objects == "virtualservices" {
		sort.SliceStable(istioConfig.VirtualServices.Items, func(i, j int) bool {
			return istioConfig.VirtualServices.Items[i].IstioBase.Metadata.Name < istioConfig.VirtualServices.Items[j].IstioBase.Metadata.Name
		})

		resp.CurrentPage = page
		resp.TotalCount = len(istioConfig.VirtualServices.Items)
		start, end, pageCount := SlicePage(page, limit, resp.TotalCount)
		resp.PageCount = pageCount
		resp.PageSize = limit
		istioConfig.VirtualServices.Items = istioConfig.VirtualServices.Items[start:end]
	} else if objects == "destinationrules" {
		sort.SliceStable(istioConfig.VirtualServices.Items, func(i, j int) bool {
			return istioConfig.DestinationRules.Items[i].IstioBase.Metadata.Name < istioConfig.DestinationRules.Items[j].IstioBase.Metadata.Name
		})

		resp.CurrentPage = page
		resp.TotalCount = len(istioConfig.VirtualServices.Items)
		start, end, pageCount := SlicePage(page, limit, resp.TotalCount)
		resp.PageCount = pageCount
		resp.PageSize = limit
		istioConfig.VirtualServices.Items = istioConfig.VirtualServices.Items[start:end]
	}

	resp.Data = istioConfig
	RespondWithJSON(w, http.StatusOK, resp)
}

type RespList struct {
	TotalCount  int         `json:"total_count"`
	PageCount   int         `json:"page_count"`
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	Data        interface{} `json:"data"`
}

func IstioConfigDetailsV2(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]
	clusterMap := business.ClusterMap
	// cidStr := r.URL.Query().Get("cid")
	// sourceCid, _ := strconv.Atoi(cidStr)

	idStr := r.Header.Get("cid")
	cid, _ := strconv.Atoi(idStr)

	if !checkObjectType(objectType) {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// Get business layer
	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	var svcCluster string

	HulkCluster, err := GetHulkClusters()
	if err != nil {
		log.Errorf("GetHulkClusters err:[%s]", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, c := range HulkCluster.Result {
		if c.Id == cid {
			svcCluster = c.Name
			break
		}
	}

	svcCluster = strings.ToLower(svcCluster)
	mirrorService, err := business.Svc.GetService(svcCluster, namespace, object, defaultHealthRateInterval, util.Clock.Now())
	if err != nil {
		log.Errorf("GetService err:[%s]", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	istioConfigDetails, err := business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		RespondWithError(w, http.StatusInternalServerError, "Services GetIstioConfigDetails error: "+err.Error())
		return
	}

	istioConfigDetails.IstioSidecar = mirrorService.IstioSidecar
	hCluster, err := GetHulkClusters()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services GetHulkClusters error: "+err.Error())
		return
	}
	if objectType == kubernetes.VirtualServices {
		if istioConfigDetails.VirtualService == nil {
			type CommonSpec struct {
				Hosts    []string    `json:"hosts,omitempty"`
				Gateways interface{} `json:"gateways,omitempty"`
				Http     interface{} `json:"http,omitempty"`
				Tcp      interface{} `json:"tcp,omitempty"`
				Tls      interface{} `json:"tls,omitempty"`
				ExportTo interface{} `json:"exportTo,omitempty"`
				QPS      float64     `json:"qps"`
			}

			istioConfigDetails.VirtualService = &models.VirtualService{
				IstioBase: models.IstioBase{
					Metadata: meta_v1.ObjectMeta{
						Name:      object,
						Namespace: namespace,
					},
				},
				Spec: CommonSpec{
					Hosts: []string{fmt.Sprintf("%s.%s.svc.cluster.local", object, namespace)},
					Http: []interface{}{
						map[string]interface{}{"mirror": nil},
					},
				},
			}

		}

		// hCluster, err := GetHulkClusters()
		// log.Errorf("GetHulkClusters-xxxx:[%+v]", hCluster)
		// log.Errorf("GetHulkClusters-xxxx:[%+v]", hCluster)
		var svcClusterName string
		for _, hc := range hCluster.Result {
			if cid == hc.Id {
				svcClusterName = hc.Name
				break
			}
		}
		if svcClusterName == "" {
			RespondWithError(w, http.StatusInternalServerError, "K8S cluster not found")
			return
		}
		svcClusterName = strings.ToLower(svcClusterName)
		query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="destination",destination_cluster="%s",destination_service_namespace="%s",destination_app="%s"}[1m]))`, svcClusterName, namespace, object)
		currentQPS, _ := api.GetMirrorQps(business, query)
		if currentQPS <= 0 {
			query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="source",destination_cluster="%s",destination_service_namespace="%s",destination_service_name="%s"}[1m]))`, svcClusterName, namespace, object)
			currentQPS, _ = api.GetMirrorQps(business, query)
		}

		istioConfigDetails.VirtualService.Spec.QPS, _ = decimal.NewFromFloat(currentQPS).Round(2).Float64()

		mirrorLabel := geneMirrorLabelV1(object, namespace)
		selectLabel := "mirror=" + mirrorLabel
		// log.Errorf("svcClusterName-xxxx:[%+v]", svcClusterName)

		allNs, err := business.Namespace.GetRemoteNamespaces(svcClusterName)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Services GetRemoteNamespaces error: "+err.Error())
			return
		}
		// log.Errorf("GetRemoteNamespaces-xxxx:[%+v]", err)
		// log.Errorf("GetRemoteNamespaces-xxxx:[%+v]", allNs)

		mirrorConfig := []interface{}{}

		for _, ns := range allNs {
			objects, _ := business.IstioConfig.GetIstioObject(ns.Name, kubernetes.EnvoyFilters, selectLabel)
			for _, ob := range objects {
				result, err := business.IstioConfig.GetIstioConfigDetails(ns.Name, kubernetes.EnvoyFilters, ob.GetObjectMeta().Name)
				if err == nil && result.EnvoyFilter != nil {
					clientInfo := result.EnvoyFilter.Metadata.Annotations["clientInfo"]
					// serviceInfo := result.EnvoyFilter.Metadata.Annotations["serviceInfo"]
					targetInfo := result.EnvoyFilter.Metadata.Annotations["targetInfo"]

					if configP, ok := result.EnvoyFilter.Spec.ConfigPatches.([]interface{}); ok {
						date, err := json.Marshal(configP[0])
						if err == nil {
							var patch MirrorConfigPatches
							err = json.Unmarshal(date, &patch)
							if err == nil {
								cid := 0
								for _, mirror := range patch.Patch.Value.Route.RequestMirrorPolicies {
									k8sClusterType := "outMesh"
									targetSvc := strings.Split(targetInfo, "|")
									if len(targetSvc) != 4 {
										log.Errorf("err:ServiceInfo format wrong [%s]", mirror.ServiceInfo)
										continue
									}
									for inC := range clusterMap {
										if targetSvc[3] == inC {
											k8sClusterType = "inMesh"
											break
										}
									}

									if targetSvc[3] == "shareit-cce-test" {
										k8sClusterType = "outMesh"
									}

									for _, hc := range hCluster.Result {
										if targetSvc[3] == hc.Name {
											cid = hc.Id
											break
										}
									}
									tPort, err := strconv.Atoi(targetSvc[0])
									if err != nil {
										log.Errorf("err:ServiceInfo format err:[%s]", err)
										continue
									}

									if cid == 0 || tPort == 0 {
										log.Errorf("err:cid tPort err:[%d],[%d]", cid, tPort)
										continue

									}

									mc := map[string]interface{}{
										"clientMirror":     false,
										"clusterType":      k8sClusterType,
										"cluster":          targetSvc[3],
										"cid":              cid,
										"namespace":        targetSvc[2],
										"service":          targetSvc[1],
										"targetPort":       tPort,
										"mirrorPercentage": float64(mirror.RuntimeFraction.DefaultValue.Numerator) / 10000,
									}

									clientSvc := strings.Split(clientInfo, "|")
									if len(clientSvc) == 3 {
										mc["clientMirror"] = true
										mc["clientService"] = clientSvc[0]
										mc["clientNamespace"] = clientSvc[1]
										mc["clientCluster"] = clientSvc[2]

									}

									mirrorConfig = append(mirrorConfig, mc)

								}
							}
						}
					}
				}

			}

		}

		if len(mirrorConfig) > 0 {
			if httpList, ok := istioConfigDetails.VirtualService.Spec.Http.([]interface{}); ok {
				var vsHttp map[string]interface{}
				if len(httpList) > 0 {
					vsHttp = httpList[0].(map[string]interface{})
				}
				vsHttp["mirror"] = mirrorConfig
				istioConfigDetails.VirtualService.Spec.Http = []interface{}{vsHttp}
			}
		}

	}

	if objectType == kubernetes.DestinationRules {
		type CommonSpec struct {
			Host          interface{} `json:"host,omitempty"`
			TrafficPolicy interface{} `json:"trafficPolicy,omitempty"`
			Subsets       interface{} `json:"subsets,omitempty"`
			ExportTo      interface{} `json:"exportTo,omitempty"`
		}

		if istioConfigDetails.DestinationRule == nil {
			err = nil
			istioConfigDetails.DestinationRule = &models.DestinationRule{
				IstioBase: models.IstioBase{
					Metadata: meta_v1.ObjectMeta{
						Name:      object,
						Namespace: namespace,
					},
				},
				Spec: CommonSpec{
					Host:          fmt.Sprintf("%s.%s.svc.cluster.local", object, namespace),
					TrafficPolicy: map[string]interface{}{},
				},
			}
		}

		rateLimitFilter := fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, object)
		result, err := business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, rateLimitFilter)
		if err == nil && result.EnvoyFilter != nil {
			if configP, ok := result.EnvoyFilter.Spec.ConfigPatches.([]interface{}); ok {
				date, err := json.Marshal(configP[0])
				if err == nil {
					var patch ConfigPatches
					err = json.Unmarshal(date, &patch)
					bucket := patch.Patch.Value.TypedConfig.Value.TokenBucket
					if err == nil {
						traffic := istioConfigDetails.DestinationRule.Spec.TrafficPolicy.(map[string]interface{})
						mToken, _ := strconv.Atoi(bucket.MaxTokens)
						tokenFil, _ := strconv.Atoi(bucket.TokensPerFill)

						traffic["rateLimit"] = TokenBucketInt{
							MaxTokens:     mToken,
							TokensPerFill: tokenFil,
							FillInterval:  bucket.FillInterval}
						istioConfigDetails.DestinationRule.Spec.TrafficPolicy = traffic
					}

				}
			}

		}

	}

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, istioConfigDetails)
}

func IstioConfigDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]
	clusterMap := business.ClusterMap
	// cidStr := r.URL.Query().Get("cid")
	// sourceCid, _ := strconv.Atoi(cidStr)

	idStr := r.Header.Get("cid")
	cid, _ := strconv.Atoi(idStr)

	includeValidations := false
	query := r.URL.Query()
	if _, found := query["validate"]; found {
		includeValidations = true
	}

	if !checkObjectType(objectType) {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// Get business layer
	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	var istioConfigValidations models.IstioValidations

	wg := sync.WaitGroup{}
	if includeValidations {
		wg.Add(1)
		go func(istioConfigValidations *models.IstioValidations, err *error) {
			defer wg.Done()
			istioConfigValidationResults, errValidations := business.Validations.GetIstioObjectValidations(namespace, objectType, object)
			if errValidations != nil && *err == nil {
				*err = errValidations
			} else {
				*istioConfigValidations = istioConfigValidationResults
			}
		}(&istioConfigValidations, &err)
	}

	istioConfigDetails, err := business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	hCluster, _ := GetHulkClusters()

	if objectType == kubernetes.VirtualServices {
		if istioConfigDetails.VirtualService == nil {
			err = nil

			type CommonSpec struct {
				Hosts    []string    `json:"hosts,omitempty"`
				Gateways interface{} `json:"gateways,omitempty"`
				Http     interface{} `json:"http,omitempty"`
				Tcp      interface{} `json:"tcp,omitempty"`
				Tls      interface{} `json:"tls,omitempty"`
				ExportTo interface{} `json:"exportTo,omitempty"`
				QPS      float64     `json:"qps"`
			}

			istioConfigDetails.VirtualService = &models.VirtualService{
				IstioBase: models.IstioBase{
					Metadata: meta_v1.ObjectMeta{
						Name:      object,
						Namespace: namespace,
					},
				},
				Spec: CommonSpec{
					Hosts: []string{fmt.Sprintf("%s.%s.svc.cluster.local", object, namespace)},
					Http: []interface{}{
						map[string]interface{}{"mirror": nil},
					},
				},
			}

		}

		// hCluster, err := GetHulkClusters()
		// log.Errorf("GetHulkClusters-xxxx:[%+v]", err)
		// log.Errorf("GetHulkClusters-xxxx:[%+v]", hCluster)
		if len(hCluster.Result) > 0 {
			var svcClusterName string
			for _, hc := range hCluster.Result {
				if cid == hc.Id {
					svcClusterName = hc.Name
					break
				}
			}
			svcClusterName = strings.ToLower(svcClusterName)
			query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="destination",destination_cluster="%s",destination_service_namespace="%s",destination_app="%s"}[1m]))`, svcClusterName, namespace, object)
			currentQPS, _ := api.GetMirrorQps(business, query)
			if currentQPS <= 0 {
				query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="source",destination_cluster="%s",destination_service_namespace="%s",destination_service_name="%s"}[1m]))`, svcClusterName, namespace, object)
				currentQPS, _ = api.GetMirrorQps(business, query)
			}

			istioConfigDetails.VirtualService.Spec.QPS, _ = decimal.NewFromFloat(currentQPS).Round(2).Float64()
		}

		// mirrorFilter := fmt.Sprintf("filter-mirror-%s", object)
		// mirrorFilter := geneMirrorName(object, namespace)

		result, err := business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, geneMirrorName(object, namespace))
		if err != nil {
			result, err = business.IstioConfig.GetIstioConfigDetails(ISTIO_SYSTEM_NAMESPACE, kubernetes.EnvoyFilters, geneMirrorName(object, namespace))
		}
		if err == nil && result.EnvoyFilter != nil {
			if configP, ok := result.EnvoyFilter.Spec.ConfigPatches.([]interface{}); ok {
				date, err := json.Marshal(configP[0])
				if err == nil {
					var patch MirrorConfigPatches
					err = json.Unmarshal(date, &patch)
					if err == nil {
						if httpList, ok := istioConfigDetails.VirtualService.Spec.Http.([]interface{}); ok {
							if len(httpList) > 0 {
								vsHttp := httpList[0].(map[string]interface{})
								// hCluster, err := GetHulkClusters()
								// if err != nil {
								// 	log.Errorf("GetHulkClusters err:[%s]", err)
								// 	RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
								// 	return
								// }

								// var svcClusterName string
								// for _, hc := range hCluster.Result {
								// 	if sourceCid == hc.Id {
								// 		svcClusterName = hc.Name
								// 		break
								// 	}
								// }

								// query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="destination",destination_cluster="%s",destination_service_namespace="%s",destination_app="%s"}[1m]))`, svcClusterName, namespace, object)

								// // query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="source",destination_cluster="%s",destination_workload_namespace="%s",destination_workload="%s"}[1m]))`, svcClusterName, namespace, object)

								// // query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="source",source_cluster="%s",source_workload_namespace="%s",source_workload="%s",destination_cluster="%s",destination_workload_namespace="%s",destination_workload="%s"}[1m]))`, sourceClusterName, namespace, object, s[3], s[2], s[1])
								// currentQPS, err := api.GetMirrorQps(business, query)
								// if err != nil {
								// 	log.Errorf("err:GetMirrorQps:[%+v]", err)
								// }

								// if currentQPS <= 0 {
								// 	query := fmt.Sprintf(`sum(rate(istio_requests_total{reporter="source",destination_cluster="%s",destination_service_namespace="%s",destination_service_name="%s"}[1m]))`, svcClusterName, namespace, object)
								// 	currentQPS, _ = api.GetMirrorQps(business, query)
								// }

								mirrorConfig := []interface{}{}
								cid := 0
								for _, mirror := range patch.Patch.Value.Route.RequestMirrorPolicies {
									k8sClusterType := "outMesh"
									s := strings.Split(mirror.ServiceInfo, "|")
									if len(s) != 4 {
										log.Errorf("err:ServiceInfo format wrong [%s]", mirror.ServiceInfo)
										continue
									}
									for inC := range clusterMap {
										if s[3] == inC {
											k8sClusterType = "inMesh"
										}
									}

									if s[3] == "shareit-cce-test" {
										k8sClusterType = "outMesh"
									}

									for _, hc := range hCluster.Result {
										if s[3] == hc.Name {
											cid = hc.Id
											break
										}
									}
									tPort, err := strconv.Atoi(s[0])
									if err != nil {
										log.Errorf("err:ServiceInfo format err:[%s]", err)
										continue
									}

									if cid == 0 || tPort == 0 {
										log.Errorf("err:cid tPort err:[%d],[%d]", cid, tPort)
										continue

									}

									// qps := float64(currentQPS) * (float64(mirror.RuntimeFraction.DefaultValue.Numerator) / 10000) / 100
									// qps, _ = decimal.NewFromFloat(qps).Round(2).Float64()
									mirrorConfig = append(mirrorConfig, map[string]interface{}{
										"clusterType":      k8sClusterType,
										"cluster":          s[3],
										"cid":              cid,
										"namespace":        s[2],
										"service":          s[1],
										"targetPort":       tPort,
										"mirrorPercentage": float64(mirror.RuntimeFraction.DefaultValue.Numerator) / 10000,
										// "currentQPS":       qps,
										// "currentQPS":       float64(currentQPS) * (float64(mirror.RuntimeFraction.DefaultValue.Numerator) / 10000) / 100,
									})

								}
								vsHttp["mirror"] = mirrorConfig

								istioConfigDetails.VirtualService.Spec.Http = []interface{}{vsHttp}
							}
						}
					}
				}

			}
		}
	}

	if objectType == kubernetes.DestinationRules {
		type CommonSpec struct {
			Host          interface{} `json:"host,omitempty"`
			TrafficPolicy interface{} `json:"trafficPolicy,omitempty"`
			Subsets       interface{} `json:"subsets,omitempty"`
			ExportTo      interface{} `json:"exportTo,omitempty"`
		}

		if istioConfigDetails.DestinationRule == nil {
			err = nil
			istioConfigDetails.DestinationRule = &models.DestinationRule{
				IstioBase: models.IstioBase{
					Metadata: meta_v1.ObjectMeta{
						Name:      object,
						Namespace: namespace,
					},
				},
				Spec: CommonSpec{
					Host:          fmt.Sprintf("%s.%s.svc.cluster.local", object, namespace),
					TrafficPolicy: map[string]interface{}{},
				},
			}
		}

		rateLimitFilter := fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, object)
		result, err := business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, rateLimitFilter)
		if err == nil && result.EnvoyFilter != nil {
			if configP, ok := result.EnvoyFilter.Spec.ConfigPatches.([]interface{}); ok {
				date, err := json.Marshal(configP[0])
				if err == nil {
					var patch ConfigPatches
					err = json.Unmarshal(date, &patch)
					bucket := patch.Patch.Value.TypedConfig.Value.TokenBucket
					if err == nil {
						traffic := istioConfigDetails.DestinationRule.Spec.TrafficPolicy.(map[string]interface{})
						mToken, _ := strconv.Atoi(bucket.MaxTokens)
						tokenFil, _ := strconv.Atoi(bucket.TokensPerFill)

						traffic["rateLimit"] = TokenBucketInt{
							MaxTokens:     mToken,
							TokensPerFill: tokenFil,
							FillInterval:  bucket.FillInterval}
						istioConfigDetails.DestinationRule.Spec.TrafficPolicy = traffic
					}

				}
			}

		}

	}
	if includeValidations && err == nil {
		wg.Wait()

		if validation, found := istioConfigValidations[models.IstioValidationKey{ObjectType: models.ObjectTypeSingular[objectType], Namespace: namespace, Name: object}]; found {
			istioConfigDetails.IstioValidation = validation
		}
	}

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, istioConfigDetails)
}

func IstioConfigDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]

	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// Get business layer
	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	err = business.IstioConfig.DeleteIstioConfigDetail(api, namespace, objectType, object)
	if err != nil {
		handleErrorResponse(w, err)
		return
	} else {
		audit(r, "DELETE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object)
		RespondWithCode(w, http.StatusOK)
	}
}

func IstioConfigUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]

	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// Get business layer
	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Update request with bad update patch: "+err.Error())
		return
	}

	newBody, err := fixRequestBody(objectType, body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "request body error: "+err.Error())
		return
	}

	jsonPatch := string(newBody)
	updatedConfigDetails, err := business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, jsonPatch)

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	audit(r, "UPDATE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object+" Patch: "+jsonPatch)
	RespondWithJSON(w, http.StatusOK, updatedConfigDetails)
}

func IstioDestinationruleCreate(w http.ResponseWriter, r *http.Request) {
	// Feels kinda replicated for multiple functions..
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := "destinationrules"

	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// Get business layer
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
	}

	createdConfigDetails, err := business.IstioConfig.CreateIstioConfigDetail(api, namespace, objectType, body)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	audit(r, "CREATE on Namespace: "+namespace+" Type: "+objectType+" Object: "+string(body))
	RespondWithJSON(w, http.StatusOK, createdConfigDetails)
}

var (
	tplFuncMap = template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}
)

const (
	TYPE_ROUTE                 = "route"
	TYPE_LOADBALANCE           = "load-balance"
	TYPE_CONNECTIONPOOL        = "connection-pool"
	TYPE_OUTLIERDETECTION      = "outlier-detection"
	TYPE_RETRY                 = "retry"
	TYPE_RATELIMIT             = "ratelimit"
	TYPE_MIRROR                = "mirror"
	TYPE_LOCALITYLB            = "locality-lbsetting"
	TYPE_FAULTINJECT           = "fault-inject"
	TYPE_SLOWSTART             = "slow-start"
	TYPE_SUBNET                = "subset"
	reteLimitEnvoyFilterPrefix = "filter-local-ratelimit-"

	ISTIO_SYSTEM_NAMESPACE = "istio-system"
	serviceEntryTpl        = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "ServiceEntry",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "istio-system",
		"labels": {
			"mirror": "{{.MirrorLabel}}",
			"targetCluster": "{{.TargetCluster}}",
			"targetService": "{{.TargetService}}",
			"targetNamespace": "{{.TargetNamespace}}"
		}
	},
	"spec": {
		"hosts": [
		"{{.Host}}"
		],
		"ports": [
		{
			"number": {{.Port}},
			"name": "{{.PortName}}",
			"protocol": "{{.AppProtocol}}"
		}
		],
		"resolution": "STATIC",
		"endpoints": [
		{{ $length := len .Address }}
		{{range $index, $element := .Address}}
		{
			"address": "{{.}}"
		}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
		{{end}}
		],
		"location": "MESH_EXTERNAL"
	}
	}	
	
`

	// 	serviceEntryTplClient = `
	// {
	// 	"apiVersion": "networking.istio.io/v1alpha3",
	// 	"kind": "ServiceEntry",
	// 	"metadata": {
	// 		"name": "{{.Name}}",
	// 		"namespace": "{{.Namespace}}",
	// 		"labels": {
	// 			"mirror": "{{.MirrorLabel}}"
	// 		}
	// 	},
	// 	"spec": {
	// 		"hosts": [
	// 		"{{.Host}}"
	// 		],
	// 		"ports": [
	// 		{
	// 			"number": {{.Port}},
	// 			"name": "http",
	// 			"protocol": "HTTP"
	// 		}
	// 		],
	// 		"resolution": "STATIC",
	// 		"endpoints": [
	// 		{{ $length := len .Address }}
	// 		{{range $index, $element := .Address}}
	// 		{
	// 			"address": "{{.}}"
	// 		}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
	// 		{{end}}
	// 		],
	// 		"location": "MESH_EXTERNAL"
	// 	}
	// 	}

	// `

	mirrorTplClientSideMultiHost = `
	{
		"apiVersion": "networking.istio.io/v1alpha3",
		"kind": "EnvoyFilter",
		"metadata": {
			"name": "{{.Name}}",
			"namespace": "istio-system"
		},
		"spec": {
			"configPatches": [
			{{ $hlength := len .HostArray }}
			{{range $hindex, $helement := .HostArray}}
			{
				"applyTo": "HTTP_ROUTE",
				"match": {
				"routeConfiguration": {
					"vhost": {
					"name": "{{.Host}}"
					}
				}
				},
				"patch": {
				"operation": "MERGE",
				"value": {
					"route": {
					"request_mirror_policies": [
						{{ $length := len .Array }}
						{{range $index, $element := .Array}}
						{
						"cluster": "{{.MirrorCluster}}",
						"service_info": "{{.ServiceInfo}}",
						"runtime_fraction": {
							"default_value": {
							"numerator": {{.Numerator}},
							"denominator": "MILLION"
							}
						}
						}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
						{{end}}
					]
					}
				}
				}
			}{{ $hindex := inc $hindex }}{{ if lt $hindex $hlength }},{{ end }}
			{{end}}
			]
		}
		}
		
	`

	mirrorTplClientSideMultiHostV2 = `
	{
		"apiVersion": "networking.istio.io/v1alpha3",
		"kind": "EnvoyFilter",
		"metadata": {
			"name": "{{.Name}}",
			"namespace": "{{.Namespace}}",
			"labels": {
				"mirror": "{{.MirrorLabel}}"
			},
			"annotations": {
				"clientInfo": "{{.ClientInfo}}",
				"serviceInfo": "{{.ServiceInfo}}",
				"targetInfo": "{{.TargetInfo}}"
			}
		},
		"spec": {
			"configPatches": [
			{{ $hlength := len .HostArray }}
			{{range $hindex, $helement := .HostArray}}
			{
				"applyTo": "HTTP_ROUTE",
				"match": {
				"routeConfiguration": {
					"vhost": {
					"name": "{{.Host}}"
					}
				}
				},
				"patch": {
				"operation": "MERGE",
				"value": {
					"route": {
					"request_mirror_policies": [
						{{ $length := len .Array }}
						{{range $index, $element := .Array}}
						{
						"cluster": "{{.MirrorCluster}}",
						"runtime_fraction": {
							"default_value": {
							"numerator": {{.Numerator}},
							"denominator": "MILLION"
							}
						}
						}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
						{{end}}
					]
					}
				}
				}
			}{{ $hindex := inc $hindex }}{{ if lt $hindex $hlength }},{{ end }}
			{{end}}
			]
		}
		}
		
	`

	mirrorTplClientSideV2 = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "{{.Namespace}}",
		"labels": {
			"mirror": "{{.MirrorLabel}}"
		},
		"annotations": {
			"clientInfo": "{{.ClientInfo}}",
			"serviceInfo": "{{.ServiceInfo}}",
			"targetInfo": "{{.TargetInfo}}"
		}
	},
	"spec": {
		"workloadSelector": {
			"labels": {
			  "app": "{{.WorkloadSelector}}"
			}
		  },
		"configPatches": [
		{
			"applyTo": "HTTP_ROUTE",
			"match": {
			"routeConfiguration": {
				"vhost": {
				"name": "{{.Host}}"
				}
			}
			},
			"patch": {
			"operation": "MERGE",
			"value": {
				"route": {
				"request_mirror_policies": [
					{{ $length := len .Array }}
					{{range $index, $element := .Array}}
					{
					"cluster": "{{.MirrorCluster}}",
					"runtime_fraction": {
						"default_value": {
						"numerator": {{.Numerator}},
						"denominator": "MILLION"
						}
					}
					}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
					{{end}}
				]
				}
			}
			}
		}
		]
	}
	}
	
`

	mirrorTplClientSideV1 = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "istio-system",
		"labels": {
			"mirror": "{{.MirrorLabel}}"
		},
		"annotations": {
			"clientInfo": "{{.ClientInfo}}",
			"serviceInfo": "{{.ServiceInfo}}",
			"targetInfo": "{{.TargetInfo}}"
		}
	},
	"spec": {
		"configPatches": [
		{
			"applyTo": "HTTP_ROUTE",
			"match": {
			"routeConfiguration": {
				"vhost": {
				"name": "{{.Host}}"
				}
			}
			},
			"patch": {
			"operation": "MERGE",
			"value": {
				"route": {
				"request_mirror_policies": [
					{{ $length := len .Array }}
					{{range $index, $element := .Array}}
					{
					"cluster": "{{.MirrorCluster}}",
					"service_info": "{{.ServiceInfo}}",
					"runtime_fraction": {
						"default_value": {
						"numerator": {{.Numerator}},
						"denominator": "MILLION"
						}
					}
					}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
					{{end}}
				]
				}
			}
			}
		}
		]
	}
	}
	
`

	mirrorTplClientSide = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "istio-system"
	},
	"spec": {
		"configPatches": [
		{
			"applyTo": "HTTP_ROUTE",
			"match": {
			"routeConfiguration": {
				"vhost": {
				"name": "{{.Host}}"
				}
			}
			},
			"patch": {
			"operation": "MERGE",
			"value": {
				"route": {
				"request_mirror_policies": [
					{{ $length := len .Array }}
					{{range $index, $element := .Array}}
					{
					"cluster": "{{.MirrorCluster}}",
					"service_info": "{{.ServiceInfo}}",
					"runtime_fraction": {
						"default_value": {
						"numerator": {{.Numerator}},
						"denominator": "MILLION"
						}
					}
					}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
					{{end}}
				]
				}
			}
			}
		}
		]
	}
	}
	
`

	mirrorTplServerSideV2 = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "{{.Namespace}}",
		"labels": {
			"mirror": "{{.MirrorLabel}}"
		},
		"annotations": {
			"clientInfo": "{{.ClientInfo}}",
			"serviceInfo": "{{.ServiceInfo}}",
			"targetInfo": "{{.TargetInfo}}"
		}
	},
	"spec": {
		"workloadSelector": {
			"labels": {
			  "app": "{{.WorkloadSelector}}"
			}
		  },
		"configPatches": [
		{
			"applyTo": "HTTP_ROUTE",
			"match": {
			"routeConfiguration": {
				"vhost": {
				"name": "{{.Host}}"
				}
			}
			},
			"patch": {
			"operation": "MERGE",
			"value": {
				"route": {
				"request_mirror_policies": [
					{{ $length := len .Array }}
					{{range $index, $element := .Array}}
					{
					"cluster": "{{.MirrorCluster}}",
					"service_info": "{{.ServiceInfo}}",
					"runtime_fraction": {
						"default_value": {
						"numerator": {{.Numerator}},
						"denominator": "MILLION"
						}
					}
					}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
					{{end}}
				]
				}
			}
			}
		}
		]
	}
	}
	
`

	mirrorTplServerSide = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "{{.Namespace}}"
	},
	"spec": {
		"workloadSelector": {
			"labels": {
			  "app": "{{.WorkloadSelector}}"
			}
		  },
		"configPatches": [
		{
			"applyTo": "HTTP_ROUTE",
			"match": {
			"routeConfiguration": {
				"vhost": {
				"name": "{{.Host}}"
				}
			}
			},
			"patch": {
			"operation": "MERGE",
			"value": {
				"route": {
				"request_mirror_policies": [
					{{ $length := len .Array }}
					{{range $index, $element := .Array}}
					{
					"cluster": "{{.MirrorCluster}}",
					"service_info": "{{.ServiceInfo}}",
					"runtime_fraction": {
						"default_value": {
						"numerator": {{.Numerator}},
						"denominator": "MILLION"
						}
					}
					}{{ $index := inc $index }}{{ if lt $index $length }},{{ end }}
					{{end}}
				]
				}
			}
			}
		}
		]
	}
	}
	
`

	rateLimtTpl = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "{{.Namespace}}"
	},
	"spec": {
		"workloadSelector": {
		"labels": {
			"app": "{{.WorkloadSelector}}"
		}
		},
		"configPatches": [
		{
			"applyTo": "HTTP_FILTER",
			"match": {
			"context": "SIDECAR_INBOUND",
			"listener": {
				"filterChain": {
				"filter": {
					"name": "envoy.filters.network.http_connection_manager"
				}
				}
			}
			},
			"patch": {
			"operation": "INSERT_BEFORE",
			"value": {
				"name": "envoy.filters.http.local_ratelimit",
				"typed_config": {
				"@type": "type.googleapis.com/udpa.type.v1.TypedStruct",
				"type_url": "type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit",
				"value": {
					"stat_prefix": "http_local_rate_limiter",
					"token_bucket": {
					"max_tokens": "{{.MaxTokens}}",
					"tokens_per_fill": "{{.TokensPerFill}}",
					"fill_interval": "{{.FillInterval}}"
					},
					"filter_enabled": {
					"runtime_key": "local_rate_limit_enabled",
					"default_value": {
						"numerator": 100,
						"denominator": "HUNDRED"
					}
					},
					"filter_enforced": {
					"runtime_key": "local_rate_limit_enforced",
					"default_value": {
						"numerator": 100,
						"denominator": "HUNDRED"
					}
					},
					"response_headers_to_add": [
					{
						"append": false,
						"header": {
						"key": "x-local-rate-limit",
						"value": "true"
						}
					}
					]
				}
				}
			}
			}
		}
	  ]
	}
}	
`
)

func geneMirrorLabel(name, namespace string) string {
	return fmt.Sprintf("mirror=mirror-%s-%s", name, namespace)
}

func geneMirrorLabelV1(name, namespace string) string {
	return fmt.Sprintf("mirror-%s-%s", name, namespace)
}

func geneMirrorName(name, namespace string) string {
	return fmt.Sprintf("filter-mirror-%s-%s", name, namespace)
}

func geneMirrorNameV2(cientSide bool, serviceName, serviceNs, mName string) string {
	if cientSide {
		return fmt.Sprintf("mirror-%s-%s-clientside-fromclient-%s", serviceName, serviceNs, mName)
	} else {
		return fmt.Sprintf("mirror-%s-%s-serverside-%s", serviceName, serviceNs, mName)
	}
}

func IstioNetworkConfigDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]
	networkType := params["network_type"]
	object := params["object"]

	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	idStr := r.Header.Get("cid")
	fmt.Println(r.Header.Get("cid"))

	cid, _ := strconv.Atoi(idStr)
	if cid == 0 {
		RespondWithError(w, http.StatusBadRequest, "cid value error:"+idStr)
		return
	}
	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	if networkType == TYPE_RATELIMIT {
		filterName := fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, object)
		err = business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.EnvoyFilters, filterName)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "delete  ratelimit error: "+err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, "")
		return
	}

	if networkType == TYPE_MIRROR {

		// hCluster, err := GetHulkClusters()
		// if err != nil {
		// 	RespondWithError(w, http.StatusInternalServerError, "GetHulkClusters error: "+err.Error())
		// 	return
		// }
		// var cName string
		// for _, hc := range hCluster.Result {
		// 	if cid == hc.Id {
		// 		cName = strings.ToLower(hc.Name)
		// 		break
		// 	}
		// }
		// if cName == "" {
		// 	RespondWithError(w, http.StatusInternalServerError, "K8S cluster not found")
		// 	return
		// }

		mirrorLabel := geneMirrorLabelV1(object, namespace)
		selectLabel := "mirror=" + mirrorLabel
		allNs, _ := business.Namespace.GetNamespaces()
		for _, ns := range allNs {
			objects, err := business.IstioConfig.GetIstioObject(ns.Name, kubernetes.EnvoyFilters, selectLabel)
			for _, ob := range objects {
				log.Infof("DeleteMirrorEnvoyFilters:[%+v]", ob)
				err = business.IstioConfig.DeleteIstioConfigDetail(api, ns.Name, kubernetes.EnvoyFilters, ob.GetObjectMeta().Name)
				if err != nil {
					log.Infof("DeleteIstioConfigDetail err:[%s]", err)
				}
			}

		}

		objects, err := business.IstioConfig.GetIstioObject(ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, selectLabel)
		for _, ob := range objects {
			log.Infof("DeleteMirrorServiceEntries:[%+v]", ob)
			err = business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, ob.GetObjectMeta().Name)
			if err != nil {
				log.Infof("DeleteIstioConfigDetail err:[%s]", err)
			}
		}

		// filterName := geneMirrorName(object, namespace)
		// business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.EnvoyFilters, filterName)

		// err = business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.EnvoyFilters, filterName)
		// if err != nil {
		// 	log.Infof("err:[%s]", err)
		// }

		// objects, err := business.IstioConfig.GetIstioObject(ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, geneMirrorLabel(object, namespace))
		// if err == nil {
		// 	for _, ob := range objects {
		// 		err = business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, ob.GetObjectMeta().Name)
		// 		if err != nil {
		// 			log.Infof("err:[%s]", err)
		// 		}
		// 	}
		// }

		RespondWithJSON(w, http.StatusOK, "")
		return
	}

	result, err := business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}
	tp := result.DestinationRule.Spec.TrafficPolicy.(map[string]interface{})
	if _, ok := tp["loadBalancer"]; !ok {
		tp["loadBalancer"] = models.LoadBalancerSettings{
			Simple: "ROUND_ROBIN",
		}
	}

	switch networkType {
	case TYPE_OUTLIERDETECTION:
		tp["outlierDetection"] = nil
	case TYPE_CONNECTIONPOOL:
		// tp["connectionPool"] = nil
	case TYPE_LOCALITYLB:
		lb := tp["loadBalancer"].(map[string]interface{})
		lb["localityLbSetting"] = nil
		tp["loadBalancer"] = lb
	// case TYPE_MIRROR:
	case TYPE_SLOWSTART:
		lb := tp["loadBalancer"].(map[string]interface{})
		lb["warmupDurationSecs"] = nil
		tp["loadBalancer"] = lb
	}

	result.DestinationRule.Spec.TrafficPolicy = tp
	jsonPatch, err := json.Marshal(result.DestinationRule)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}
	_, err = business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, string(jsonPatch))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "IstioNetworkConfigDelete error: "+err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, "")
}

func isDuration(d string) bool {
	unit := d[len(d)-1:]
	if unit != "ms" && unit != "s" && unit != "m" && unit != "h" {
		return false
	}

	for _, ch := range d[:len(d)-1] {
		if !unicode.IsDigit(ch) {
			return false
		}
	}

	// _, err := strconv.ParseFloat(d[:len(d)-1], 64)
	// return err == nil

	return true

}

func genMirrorCluster(port int32, svc, ns string) string {

	return fmt.Sprintf("outbound|%d||%s.%s.svc.cluster.local", port, svc, ns)

	// return fmt.Sprintf("outbound|%d||%s.%s.svc.cluster.local", port, svc, ns)

}

func genMirrorClusterClientSide(port int32, svc, ns string) string {

	// return fmt.Sprintf("outbound|%d||%s.%s.svc.cluster.local", port, svc, ns)

	return fmt.Sprintf("outbound|%d||%s.%s.svc.cluster.local", port, svc, ns)

}

func IstioNetworkConfigUpdateMirrorV2(w http.ResponseWriter, r *http.Request, params map[string]string) {
	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]
	clusterMap := business.ClusterMap

	idStr := r.Header.Get("cid")
	cid, _ := strconv.Atoi(idStr)

	api := business.GetIstioAPI(objectType)

	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
		return
	}

	objectType = kubernetes.EnvoyFilters
	dstRule := &models.VirtualServiceM{}
	err = json.Unmarshal(body, dstRule)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}

	for _, m := range dstRule.Spec.Http[0].Mirror {
		if m.Cid == cid && m.Namespace == namespace && m.Service == object {
			RespondWithError(w, http.StatusBadRequest, "源服务不能和目的服务相同")
			return
		}

		if m.TargetPort <= 0 {
			log.Errorf("m.TargetPort <= 0:[%s]", m.TargetPort)
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("目的服务端口不能小于0:%d", m.TargetPort))
			return
		}
	}

	type Service struct {
		MirrorCluster string `json:"mirror_cluster"`
		Numerator     string `json:"numerator"`
		ServiceInfo   string `json:"serviceInfo"`
	}
	type Mirror struct {
		Name             string `json:"name"`
		Namespace        string `json:"namespace"`
		MirrorLabel      string `json:"mirrorLabel"`
		Array            []Service
		Host             string `json:"host"`
		WorkloadSelector string `json:"workloadSelector"`
		ClientInfo       string `json:"clientInfo"`
		ServiceInfo      string `json:"serviceInfo"`
		TargetInfo       string `json:"targetInfo"`
	}

	var mirrorCluster string

	HulkCluster, err := GetHulkClusters()
	if err != nil {
		log.Errorf("GetHulkClusters err:[%s]", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, c := range HulkCluster.Result {
		if c.Id == cid {
			mirrorCluster = c.Name
			break
		}
	}

	mirrorCluster = strings.ToLower(mirrorCluster)
	mirrorService, err := business.Svc.GetService(mirrorCluster, dstRule.Metadata.Namespace, dstRule.Metadata.Name, defaultHealthRateInterval, util.Clock.Now())
	// for _, w := range targerService.Workloads {
	// 	log.Errorf("GetService xxxxxx:[%+v]", w)

	// }

	// log.Errorf("GetService xxxxxx:[%+v]", targerService)
	// log.Errorf("GetService xxxxxx:[%+s]", err)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}
	if len(mirrorService.Service.Ports) < 1 {
		RespondWithError(w, http.StatusBadRequest, "服务端口为空")
		return
	}

	mirrorLabel := geneMirrorLabelV1(object, namespace)
	selectLabel := "mirror=" + mirrorLabel
	allNs, _ := business.Namespace.GetRemoteNamespaces(mirrorCluster)
	for _, ns := range allNs {
		objects, err := business.IstioConfig.GetIstioObject(ns.Name, kubernetes.EnvoyFilters, selectLabel)
		for _, ob := range objects {
			log.Infof("DeleteMirrorEnvoyFilters:[%+v]", ob)
			err = business.IstioConfig.DeleteIstioConfigDetail(api, ns.Name, kubernetes.EnvoyFilters, ob.GetObjectMeta().Name)
			if err != nil {
				log.Infof("DeleteIstioConfigDetail err:[%s]", err)
			}
		}

	}

	objects, err := business.IstioConfig.GetIstioObject(ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, selectLabel)
	for _, ob := range objects {
		log.Infof("DeleteMirrorServiceEntries:[%+v]", ob)
		err = business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, ob.GetObjectMeta().Name)
		if err != nil {
			log.Infof("DeleteIstioConfigDetail err:[%s]", err)
		}
	}

	appProtocol := "HTTP"
	if strings.Contains(mirrorService.Service.Ports[0].Name, "grpc") || strings.Contains(mirrorService.Service.Ports[0].Name, "GRPC") {
		appProtocol = "GRPC"
	}
	var inMesh bool
	for _, m := range dstRule.Spec.Http[0].Mirror {
		var mConfig Mirror
		var mirrorTpl string
		if mirrorService.IstioSidecar {
			mirrorTpl = mirrorTplServerSideV2
			mConfig.Host = fmt.Sprintf("inbound|http|%d", mirrorService.Service.Ports[0].Port)
			mConfig.Namespace = namespace
			mConfig.WorkloadSelector = object
		} else {
			mirrorTpl = mirrorTplClientSideV1
			mConfig.Host = fmt.Sprintf("%s.%s.svc.cluster.local:%d", dstRule.Metadata.Name, dstRule.Metadata.Namespace, mirrorService.Service.Ports[0].Port)
			mConfig.Namespace = ISTIO_SYSTEM_NAMESPACE
		}

		mConfig.Name = geneMirrorNameV2(false, object, namespace, m.Service)
		mConfig.MirrorLabel = mirrorLabel

		mConfig.ServiceInfo = fmt.Sprintf("%d|%s|%s|%s", mirrorService.Service.Ports[0].Port, dstRule.Metadata.Name, dstRule.Metadata.Namespace, mirrorCluster)
		// mConfig.ClientInfo = fmt.Sprintf("%s|%s|%s", m.ClientService, m.ClientNamespace, m.ClientCluster)
		mConfig.TargetInfo = fmt.Sprintf("%d|%s|%s|%s", m.TargetPort, m.Service, m.Namespace, m.Cluster)

		if m.ClientMirror {
			mirrorTpl = mirrorTplClientSideV2
			mConfig.Host = fmt.Sprintf("%s.%s.svc.cluster.local:%d", dstRule.Metadata.Name, dstRule.Metadata.Namespace, mirrorService.Service.Ports[0].Port)
			mConfig.Namespace = m.ClientNamespace
			mConfig.WorkloadSelector = m.ClientService
			mConfig.ClientInfo = fmt.Sprintf("%s|%s|%s", m.ClientService, m.ClientNamespace, m.ClientCluster)

			fmt.Println("xxxxxclientinfo:", object, namespace, m.ClientService)

			mConfig.Name = geneMirrorNameV2(true, object, namespace, m.ClientService)
			fmt.Println("xxxxxclient:", mConfig.Name)

		}
		fmt.Println("xxxxxnow:", mConfig.Name)

		lowerCName := strings.ToLower(m.Cluster)
		inMesh = false
		for cName := range clusterMap {
			if lowerCName == "shareit-cce-test" {
				break
			}
			if lowerCName == cName {
				inMesh = true
				mirrorService, err := business.Svc.GetService(lowerCName, m.Namespace, m.Service, defaultHealthRateInterval, util.Clock.Now())
				if err != nil {
					log.Infof("err:[%s]", err)
					handleErrorResponse(w, err)
					return
				}
				mConfig.Array = append(mConfig.Array, Service{
					MirrorCluster: genMirrorCluster(mirrorService.Service.Ports[0].Port, m.Service, m.Namespace),
					// MirrorCluster: fmt.Sprintf("outbound|%d||%s.%s.svc.cluster.local", mirrorService.Service.Ports[0].Port, m.Service, m.Namespace),
					Numerator:   strconv.FormatInt(int64(m.MirrorPercentage*10000), 10),
					ServiceInfo: fmt.Sprintf("%d|%s|%s|%s", m.TargetPort, m.Service, m.Namespace, m.Cluster),
				})
			}
		}

		if !inMesh {
			var clusterFound bool
			for _, c := range HulkCluster.Result {
				clusterFound = false
				if m.Cluster == c.Name {
					podIps, err := GetHulkClusterEndpointsIps(m.Service, m.Namespace, m.Cid)
					if err != nil {
						log.Infof("GetHulkClusterEndpointsIps err:[%s]", err)
						RespondWithError(w, http.StatusInternalServerError, err.Error())
						return
					}

					type ServiceEn struct {
						Name            string   `json:"name"`
						Address         []string `json:"address"`
						Host            string   `json:"host"`
						Port            string   `json:"port"`
						MirrorLabel     string   `json:"mirrorLabel"`
						TargetService   string   `json:"TargetService"`
						TargetNamespace string   `json:"TargetNamespace"`
						TargetCluster   string   `json:"TargetCluster"`
						PortName        string   `json:"PortName"`
						AppProtocol     string   `json:"AppProtocol"`
					}

					snName := fmt.Sprintf("mirror-%s-%s-to-%s-%s", object, namespace, m.Service, m.Namespace)
					err = ValidateWildcardDomain(snName)
					if err != nil {
						snName = fmt.Sprintf("mirror-%s-%s-to-randomname-%s", object, namespace, RandStringBytesRmndr(5))
					}
					snHost := snName + ".ushareit"
					sn := ServiceEn{
						Name:        snName,
						Host:        snHost,
						Port:        strconv.Itoa(m.TargetPort),
						Address:     podIps,
						PortName:    mirrorService.Service.Ports[0].Name,
						AppProtocol: appProtocol,
						// MirrorLabel:     fmt.Sprintf("mirror-%s-%s", object, namespace),
						MirrorLabel:     mirrorLabel,
						TargetCluster:   m.Cluster,
						TargetService:   m.Service,
						TargetNamespace: m.Namespace,
					}

					t, err := template.New("test").Funcs(tplFuncMap).Parse(serviceEntryTpl)
					if err != nil {
						log.Infof("template:[%s]", err)
						handleErrorResponse(w, err)
						return
					}
					var buf bytes.Buffer
					err = t.Execute(&buf, sn)
					if err != nil {
						log.Errorf("err:[%s]", err)
						handleErrorResponse(w, err)
						return
					}

					body = buf.Bytes()
					business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, snName)
					_, err = business.IstioConfig.CreateIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, body)
					if err != nil {
						handleErrorResponse(w, err)
						return
					}

					mConfig.Array = append(mConfig.Array, Service{
						MirrorCluster: fmt.Sprintf("outbound|%d||%s", m.TargetPort, snHost),
						Numerator:     strconv.FormatInt(int64(m.MirrorPercentage*10000), 10),
						ServiceInfo:   fmt.Sprintf("%d|%s|%s|%s", m.TargetPort, m.Service, m.Namespace, m.Cluster),
					})

					clusterFound = true
					break
				}
			}

			if !clusterFound {
				log.Errorf("hulkCluster not found:[%s][%+v]", m.Cluster, HulkCluster.Result)
				RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("hulkCluster not found:%s", m.Cluster))
				return
			}
		}

		var multiHost bool

		if m.ClientMirror && len(mirrorService.VirtualServices.Items) > 0 {
			vs := mirrorService.VirtualServices.Items[0]
			if len(vs.Spec.Hosts) > 0 {
				multiHost = true
				type A struct {
					Name             string `json:"name"`
					Namespace        string `json:"namespace"`
					WorkloadSelector string `json:"workloadSelector"`
					MirrorLabel      string `json:"mirrorLabel"`
					ClientInfo       string `json:"clientInfo"`
					ServiceInfo      string `json:"serviceInfo"`
					TargetInfo       string `json:"targetInfo"`
					HostArray        []Mirror
				}
				mirrorTpl = mirrorTplClientSideMultiHostV2

				mirrorConfigs := A{
					Name:        mConfig.Name,
					Namespace:   m.ClientNamespace,
					HostArray:   []Mirror{mConfig},
					MirrorLabel: mConfig.MirrorLabel,
					ClientInfo:  mConfig.ClientInfo,
					ServiceInfo: mConfig.ServiceInfo,
					TargetInfo:  mConfig.TargetInfo,
				}
				// mirrorConfigs.Name = mConfig.Name
				// mirrorConfigs.HostArray = []Mirror{mConfig}
				for _, h := range vs.Spec.Hosts {
					if h == object || h == object+"."+namespace+".svc.cluster.local" {
						continue
					}
					hm := mConfig
					hm.Host = h + ":80"
					mirrorConfigs.HostArray = append(mirrorConfigs.HostArray, hm)
				}

				t, err := template.New("test").Funcs(tplFuncMap).Parse(mirrorTpl)
				if err != nil {
					log.Infof("template:[%s]", err)
					handleErrorResponse(w, err)
					return
				}
				var buf bytes.Buffer
				err = t.Execute(&buf, mirrorConfigs)
				if err != nil {
					log.Infof("template Execute:[%s]", err)
					handleErrorResponse(w, err)
					return
				}

				body = buf.Bytes()
				log.Infof("mirror filter:[%s]", string(body))
				object = geneMirrorName(object, namespace)
			}

		}
		if !multiHost {
			t, err := template.New("test").Funcs(tplFuncMap).Parse(mirrorTpl)
			if err != nil {
				log.Infof("template:[%s]", err)
				handleErrorResponse(w, err)
				return
			}
			var buf bytes.Buffer
			err = t.Execute(&buf, mConfig)
			if err != nil {
				log.Infof("template Execute:[%s]", err)
				handleErrorResponse(w, err)
				return
			}

			body = buf.Bytes()
			// log.Infof("mirror filter:[%s]", string(body))
			// object = geneMirrorName(object, namespace)
		}

		// _, err = business.IstioConfig.GetIstioConfigDetails(mConfig.Namespace, objectType, object)
		fmt.Println("create-mirror-xxx:", mConfig.Name, objectType, string(body))

		_, err := business.IstioConfig.CreateIstioConfigDetail(api, mConfig.Namespace, objectType, body)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "CREATE on Namespace: "+mConfig.Namespace+" Type: "+objectType+" Object: "+string(body))
	}

	RespondWithJSON(w, http.StatusOK, "")
}

func IstioNetworkConfigUpdateMirror(w http.ResponseWriter, r *http.Request, params map[string]string) {
	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]
	clusterMap := business.ClusterMap
	cidStr := r.URL.Query().Get("cid")
	sourceCid, _ := strconv.Atoi(cidStr)
	api := business.GetIstioAPI(objectType)

	business, err := getBusinessByCluster(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
		return
	}

	objectType = kubernetes.EnvoyFilters
	dstRule := &models.VirtualServiceM{}
	err = json.Unmarshal(body, dstRule)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}

	for _, m := range dstRule.Spec.Http[0].Mirror {
		if m.Cid == sourceCid && m.Namespace == namespace && m.Service == object {
			RespondWithError(w, http.StatusBadRequest, "源服务不能和目的服务相同")
			return
		}

		if m.TargetPort <= 0 {
			log.Errorf("m.TargetPort <= 0:[%s]", m.TargetPort)
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("目的服务端口不能小于0:%d", m.TargetPort))
			return
		}
	}

	type Service struct {
		MirrorCluster string `json:"mirror_cluster"`
		Numerator     string `json:"numerator"`
		ServiceInfo   string `json:"serviceInfo"`
	}
	type Mirror struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Array     []Service
		Host      string `json:"host"`
		// HostArray        []string `json:"HostArray"`
		WorkloadSelector string `json:workloadSelector`
		// MirrorCluster string `json:"mirror_cluster"`
		// Numerator     string `json:"numerator"`
	}

	var cName string

	HulkCluster, err := GetHulkClusters()
	if err != nil {
		log.Errorf("GetHulkClusters err:[%s]", err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, c := range HulkCluster.Result {
		if c.Id == sourceCid {
			cName = c.Name
			break
		}
	}

	cName = strings.ToLower(cName)
	targerService, err := business.Svc.GetService(cName, dstRule.Metadata.Namespace, dstRule.Metadata.Name, defaultHealthRateInterval, util.Clock.Now())
	// for _, w := range targerService.Workloads {
	// 	log.Errorf("GetService xxxxxx:[%+v]", w)

	// }

	// log.Errorf("GetService xxxxxx:[%+v]", targerService)
	// log.Errorf("GetService xxxxxx:[%+s]", err)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}
	if len(targerService.Service.Ports) < 1 {
		RespondWithError(w, http.StatusBadRequest, "服务端口为空")
		return
	}

	var mConfig = Mirror{
		Name: geneMirrorName(object, namespace),
	}

	objects, err := business.IstioConfig.GetIstioObject(ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, geneMirrorLabel(object, namespace))
	if err == nil {
		for _, ob := range objects {
			err = business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, ob.GetObjectMeta().Name)
			if err != nil {
				log.Infof("DeleteIstioConfigDetail err:[%s]", err)
			}
		}
	}

	business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.EnvoyFilters, mConfig.Name)
	business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.EnvoyFilters, mConfig.Name)

	appProtocol := "HTTP"
	if strings.Contains(targerService.Service.Ports[0].Name, "grpc") || strings.Contains(targerService.Service.Ports[0].Name, "GRPC") {
		appProtocol = "GRPC"
	}
	var mirrorTpl string
	if targerService.IstioSidecar {
		mirrorTpl = mirrorTplServerSide
		mConfig.Host = fmt.Sprintf("inbound|http|%d", targerService.Service.Ports[0].Port)
		mConfig.Namespace = namespace
		mConfig.WorkloadSelector = object
	} else {
		mirrorTpl = mirrorTplClientSide
		mConfig.Host = fmt.Sprintf("%s.%s.svc.cluster.local:%d", dstRule.Metadata.Name, dstRule.Metadata.Namespace, targerService.Service.Ports[0].Port)
		mConfig.Namespace = ISTIO_SYSTEM_NAMESPACE
	}

	var inMesh bool
	for _, m := range dstRule.Spec.Http[0].Mirror {

		// if m.ClientMirror {
		// 	mirrorTpl = mirrorTplClientSideV2

		// }
		lowerCName := strings.ToLower(m.Cluster)
		inMesh = false
		for cName := range clusterMap {
			if lowerCName == "shareit-cce-test" {
				break
			}
			if lowerCName == cName {
				inMesh = true
				mirrorService, err := business.Svc.GetService(lowerCName, m.Namespace, m.Service, defaultHealthRateInterval, util.Clock.Now())
				if err != nil {
					log.Infof("err:[%s]", err)
					handleErrorResponse(w, err)
					return
				}
				mConfig.Array = append(mConfig.Array, Service{
					MirrorCluster: genMirrorCluster(mirrorService.Service.Ports[0].Port, m.Service, m.Namespace),
					// MirrorCluster: fmt.Sprintf("outbound|%d||%s.%s.svc.cluster.local", mirrorService.Service.Ports[0].Port, m.Service, m.Namespace),
					Numerator:   strconv.FormatInt(int64(m.MirrorPercentage*10000), 10),
					ServiceInfo: fmt.Sprintf("%d|%s|%s|%s", m.TargetPort, m.Service, m.Namespace, m.Cluster),
				})

				if m.ClientMirror {
					mirrorTpl = mirrorTplClientSideV2
					mConfig.Host = fmt.Sprintf("%s.%s.svc.cluster.local:%d", dstRule.Metadata.Name, dstRule.Metadata.Namespace, targerService.Service.Ports[0].Port)
					mConfig.Namespace = m.ClientNamespace
					mConfig.WorkloadSelector = m.Service

				}

				inMesh = true
				break
			}
		}

		if inMesh {
			continue
		}

		var clusterFound bool
		for _, c := range HulkCluster.Result {
			clusterFound = false
			// lowerCName := strings.ToLower(m.Cluster)

			if m.Cluster == c.Name {
				podIps, err := GetHulkClusterEndpointsIps(m.Service, m.Namespace, m.Cid)
				if err != nil {
					log.Infof("GetHulkClusterEndpointsIps err:[%s]", err)
					RespondWithError(w, http.StatusInternalServerError, err.Error())
					return
				}

				type ServiceEn struct {
					Name string `json:"name"`
					// Namespace   string   `json:"namespace"`
					Address         []string `json:"address"`
					Host            string   `json:"host"`
					Port            string   `json:"port"`
					MirrorLabel     string   `json:"mirrorLabel"`
					TargetService   string   `json:"TargetService"`
					TargetNamespace string   `json:"TargetNamespace"`
					TargetCluster   string   `json:"TargetCluster"`
					PortName        string   `json:"PortName"`
					AppProtocol     string   `json:"AppProtocol"`
					// MirrorCluster string `json:"mirror_cluster"`
					// Numerator     string `json:"numerator"`
				}

				// snName := fmt.Sprintf("mirror-%s-to-%s", object, m.Service)
				snName := fmt.Sprintf("mirror-%s-%s-to-%s-%s", object, namespace, m.Service, m.Namespace)
				err = ValidateWildcardDomain(snName)
				if err != nil {
					snName = fmt.Sprintf("mirror-%s-%s-to-randomname-%s", object, namespace, RandStringBytesRmndr(5))
				}
				snHost := snName + ".ushareit"
				sn := ServiceEn{
					Name: snName,
					// Namespace:   namespace,
					Host:            snHost,
					Port:            strconv.Itoa(m.TargetPort),
					Address:         podIps,
					PortName:        targerService.Service.Ports[0].Name,
					AppProtocol:     appProtocol,
					MirrorLabel:     fmt.Sprintf("mirror-%s-%s", object, namespace),
					TargetCluster:   m.Cluster,
					TargetService:   m.Service,
					TargetNamespace: m.Namespace,
				}

				t, err := template.New("test").Funcs(tplFuncMap).Parse(serviceEntryTpl)
				if err != nil {
					log.Infof("template:[%s]", err)
					handleErrorResponse(w, err)
					return
				}
				var buf bytes.Buffer
				err = t.Execute(&buf, sn)
				if err != nil {
					log.Errorf("err:[%s]", err)
					handleErrorResponse(w, err)
					return
				}

				body = buf.Bytes()
				business.IstioConfig.DeleteIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, snName)
				_, err = business.IstioConfig.CreateIstioConfigDetail(api, ISTIO_SYSTEM_NAMESPACE, kubernetes.ServiceEntries, body)
				if err != nil {
					handleErrorResponse(w, err)
					return
				}

				//Request URL: https://scmp.ushareit.me/hulk/api/v2/apps/sprs/deployment/store-house/pods
				// podUrl := fmt.Sprintf("http://scmp-hulk.sgt:80/hulk/api/v2/apps/%s/%s/%s/pods", m.Namespace, appType, m.Service)

				mConfig.Array = append(mConfig.Array, Service{
					// MirrorCluster: genMirrorCluster(int32(m.TargetPort), object, m.Service),
					MirrorCluster: fmt.Sprintf("outbound|%d||%s", m.TargetPort, snHost),

					// MirrorCluster: fmt.Sprintf("outbound|%d||mirror-%s-to-%s.ushareit", m.TargetPort, object, m.Service),
					Numerator:   strconv.FormatInt(int64(m.MirrorPercentage*10000), 10),
					ServiceInfo: fmt.Sprintf("%d|%s|%s|%s", m.TargetPort, m.Service, m.Namespace, m.Cluster),
				})

				clusterFound = true
				break
			}
		}

		if !clusterFound {
			log.Errorf("hulkCluster not found:[%s][%+v]", m.Cluster, HulkCluster.Result)
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("hulkCluster not found:%s", m.Cluster))
			return
		}
	}

	// mirrorConfigs := []Mirror{mConfig}
	var multiHost bool
	if !targerService.IstioSidecar {
		if len(targerService.VirtualServices.Items) > 0 {
			vs := targerService.VirtualServices.Items[0]
			if len(vs.Spec.Hosts) > 0 {
				multiHost = true
				type A struct {
					Name      string `json:"name"`
					HostArray []Mirror
				}
				mirrorTpl = mirrorTplClientSideMultiHost

				var mirrorConfigs A
				mirrorConfigs.Name = mConfig.Name
				mirrorConfigs.HostArray = []Mirror{mConfig}
				for _, h := range vs.Spec.Hosts {
					if h == object || h == object+"."+namespace+".svc.cluster.local" {
						continue
					}
					hm := mConfig
					hm.Host = h + ":80"
					mirrorConfigs.HostArray = append(mirrorConfigs.HostArray, hm)
				}

				t, err := template.New("test").Funcs(tplFuncMap).Parse(mirrorTpl)
				if err != nil {
					log.Infof("template:[%s]", err)
					handleErrorResponse(w, err)
					return
				}
				var buf bytes.Buffer
				err = t.Execute(&buf, mirrorConfigs)
				if err != nil {
					log.Infof("template Execute:[%s]", err)
					handleErrorResponse(w, err)
					return
				}

				body = buf.Bytes()
				log.Infof("mirror filter:[%s]", string(body))
				object = geneMirrorName(object, namespace)
			}
		}
	}

	if !multiHost {
		t, err := template.New("test").Funcs(tplFuncMap).Parse(mirrorTpl)
		if err != nil {
			log.Infof("template:[%s]", err)
			handleErrorResponse(w, err)
			return
		}
		var buf bytes.Buffer
		err = t.Execute(&buf, mConfig)
		if err != nil {
			log.Infof("template Execute:[%s]", err)
			handleErrorResponse(w, err)
			return
		}

		body = buf.Bytes()
		log.Infof("mirror filter:[%s]", string(body))
		object = geneMirrorName(object, namespace)
	}

	_, err = business.IstioConfig.GetIstioConfigDetails(mConfig.Namespace, objectType, object)
	if err == nil {
		// fmt.Printf("IstioNetworkConfig-updatemirror-xxx:%+v\n", string(body))

		jsonPatch := string(body)
		updatedConfigDetails, err := business.IstioConfig.UpdateIstioConfigDetail(api, mConfig.Namespace, objectType, object, jsonPatch)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "UPDATE on Namespace: "+mConfig.Namespace+" Type: "+objectType+" Name: "+object+" Patch: "+jsonPatch)
		RespondWithJSON(w, http.StatusOK, updatedConfigDetails)
	} else {
		fmt.Println("create-xxx:", mConfig.Namespace, objectType, string(body))

		createdConfigDetails, err := business.IstioConfig.CreateIstioConfigDetail(api, mConfig.Namespace, objectType, body)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "CREATE on Namespace: "+mConfig.Namespace+" Type: "+objectType+" Object: "+string(body))
		RespondWithJSON(w, http.StatusOK, createdConfigDetails)
	}

}

func IstioNetworkConfigUpdateRateLimit(w http.ResponseWriter, r *http.Request, params map[string]string) {

	namespace := params["namespace"]
	objectType := params["object_type"]
	// object := params["object"]
	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
	}

	objectType = kubernetes.EnvoyFilters
	dstRule := &models.DestinationRuleM{}
	err = json.Unmarshal(body, dstRule)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}

	type limit struct {
		Name             string `json:"name"`
		WorkloadSelector string `json:"WorkloadSelector"`
		Namespace        string `json:"namespace"`
		MaxTokens        int    `json:"max_tokens"`
		TokensPerFill    int    `json:"tokens_per_fill"`
		FillInterval     string `json:"fill_interval"`
	}
	li := dstRule.Spec.TrafficPolicy.RateLimit
	object := fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, dstRule.Metadata.Name)
	p := limit{
		WorkloadSelector: dstRule.Metadata.Name,
		Name:             object,
		Namespace:        dstRule.Metadata.Namespace,
	}

	if fill, ok := li["fillInterval"].(string); ok {
		p.FillInterval = fill
	}

	if num, ok := li["maxTokens"].(float64); ok {
		p.MaxTokens = int(num)
	}
	if num, ok := li["tokensPerFill"].(float64); ok {
		p.TokensPerFill = int(num)
	}

	t, err := template.New("test").Parse(rateLimtTpl)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, p)
	if err != nil {
		log.Infof("err:[%s]", err)
		handleErrorResponse(w, err)
		return
	}

	body = buf.Bytes()
	_, err = business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	if err == nil {
		jsonPatch := string(body)
		updatedConfigDetails, err := business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, jsonPatch)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "UPDATE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object+" Patch: "+jsonPatch)
		RespondWithJSON(w, http.StatusOK, updatedConfigDetails)
	} else {
		createdConfigDetails, err := business.IstioConfig.CreateIstioConfigDetail(api, namespace, objectType, body)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "CREATE on Namespace: "+namespace+" Type: "+objectType+" Object: "+string(body))
		RespondWithJSON(w, http.StatusOK, createdConfigDetails)
	}

}

func IstioNetworkConfigUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	networkType := params["network_type"]

	if networkType == TYPE_MIRROR {
		IstioNetworkConfigUpdateMirrorV2(w, r, params)
		return
	}
	if networkType == TYPE_RATELIMIT {
		IstioNetworkConfigUpdateRateLimit(w, r, params)
		return
	}

	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]
	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// business, err := getBusiness(r)
	business, err := getBusinessByCluster(r)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
	}

	// if networkType == TYPE_RATELIMIT {
	// 	objectType = kubernetes.EnvoyFilters
	// 	dstRule := &models.DestinationRuleM{}
	// 	err = json.Unmarshal(body, dstRule)
	// 	if err != nil {
	// 		log.Infof("err:[%s]", err)
	// 		handleErrorResponse(w, err)
	// 		return
	// 	}

	// 	type limit struct {
	// 		Name          string `json:name`
	// 		Namespace     string `json:namespace`
	// 		MaxTokens     int    `json:max_tokens`
	// 		TokensPerFill int    `json:tokens_per_fill`
	// 		FillInterval  string `json:fill_interval`
	// 	}
	// 	// trafic := dstRule.Spec.TrafficPolicy
	// 	li := dstRule.Spec.TrafficPolicy.RateLimit
	// 	object = fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, dstRule.Metadata.Name)
	// 	p := limit{
	// 		Name:      object,
	// 		Namespace: dstRule.Metadata.Namespace,
	// 	}

	// 	if fill, ok := li["fillInterval"].(string); ok {
	// 		p.FillInterval = fill
	// 	}

	// 	if num, ok := li["maxTokens"].(float64); ok {
	// 		p.MaxTokens = int(num)
	// 	}
	// 	if num, ok := li["tokensPerFill"].(float64); ok {
	// 		p.TokensPerFill = int(num)
	// 	}

	// 	t, err := template.New("test").Parse(rateLimtTpl)
	// 	if err != nil {
	// 		log.Infof("err:[%s]", err)
	// 		handleErrorResponse(w, err)
	// 		return
	// 	}
	// 	var buf bytes.Buffer
	// 	err = t.Execute(&buf, p)
	// 	if err != nil {
	// 		log.Infof("err:[%s]", err)
	// 		handleErrorResponse(w, err)
	// 		return
	// 	}

	// 	body = buf.Bytes()

	// }

	if objectType == kubernetes.DestinationRules {
		if networkType == TYPE_SLOWSTART {
			result, err := business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.DestinationRules, object)
			if err == nil {
				if result.DestinationRule != nil {
					dstRule := result.DestinationRule
					if tp, ok := dstRule.Spec.TrafficPolicy.(map[string]interface{}); ok {
						if lb, ok := tp["loadBalancer"].(map[string]interface{}); ok {
							if m, ok := lb["simple"]; ok {
								if m != "LEAST_REQUEST" && m != "ROUND_ROBIN" {
									log.Errorf("loadBalancer info not match:[%s][%+v]", result.DestinationRule)
									RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("慢启动仅支持:LEAST_REQUEST,ROUND_ROBIN,当前使用的负载均衡算法:%s", m))
									return
								}
							}
						}
					}
				}
			}
		}

		dstRule := &models.DestinationRuleM{}
		err = json.Unmarshal(body, dstRule)
		if err != nil {
			log.Infof("err:[%s]", err)
			handleErrorResponse(w, err)
			return
		}
		dstRule.Spec.Host = fmt.Sprintf("%s.%s.svc.cluster.local", dstRule.Metadata.Name, dstRule.Metadata.Namespace)
		body, err = json.Marshal(dstRule)
		if err != nil {
			log.Infof("err:[%s]", err)
			handleErrorResponse(w, err)
			return
		}

		if networkType == TYPE_OUTLIERDETECTION {
			if dstRule.Spec.TrafficPolicy.OutlierDetection != nil {
				outli := dstRule.Spec.TrafficPolicy.OutlierDetection
				if !isDuration(*outli.BaseEjectionTime) || !isDuration(*outli.Interval) {
					RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("时间格式错误:%s,%s", *outli.BaseEjectionTime, *outli.Interval))
					return
				}
			}
		}
	}

	_, err = business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	if err == nil {
		jsonPatch := string(body)
		updatedConfigDetails, err := business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, jsonPatch)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "UPDATE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object+" Patch: "+jsonPatch)
		RespondWithJSON(w, http.StatusOK, updatedConfigDetails)
	} else {
		fmt.Println("create-xxx:", namespace, objectType, string(body))

		createdConfigDetails, err := business.IstioConfig.CreateIstioConfigDetail(api, namespace, objectType, body)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "CREATE on Namespace: "+namespace+" Type: "+objectType+" Object: "+string(body))
		RespondWithJSON(w, http.StatusOK, createdConfigDetails)
	}
}
func IstioConfigCreate(w http.ResponseWriter, r *http.Request) {
	// Feels kinda replicated for multiple functions..
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]

	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	// Get business layer
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
	}

	newBody, err := fixRequestBody(objectType, body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "request body error: "+err.Error())
	}

	createdConfigDetails, err := business.IstioConfig.CreateIstioConfigDetail(api, namespace, objectType, newBody)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	audit(r, "CREATE on Namespace: "+namespace+" Type: "+objectType+" Object: "+string(newBody))
	RespondWithJSON(w, http.StatusOK, createdConfigDetails)
}

func fixRequestBody(objectType string, body []byte) (res []byte, err error) {

	switch objectType {
	case kubernetes.VirtualServices:
		res = body
		// v := &models.VirtualServiceM{}
		// err = json.Unmarshal(body, v)
		// if err != nil {
		// 	log.Infof("err:[%s]", err)
		// 	return
		// }
		// v.Metadata.Namespace = "sample"

		// for i, _ := range v.Spec.Http {
		// 	for j, _ := range v.Spec.Http[i].Route {
		// 		v.Spec.Http[i].Route[j].Destination.Host = v.Spec.Http[i].Route[j].Destination.Host + ".svc.cluster.local"
		// 	}
		// }
		// for i := range v.Spec.Http {
		// 	if v.Spec.Http[i].Retries == nil {
		// 		v.Spec.Http[i].Retries = &models.HTTPRetry{
		// 			Attempts: 0,
		// 		}
		// 	}
		// }

		// res, err = json.Marshal(v)
		// if err != nil {
		// 	log.Infof("err:[%s]", err)
		// 	return
		// }
	case kubernetes.DestinationRules:
		v := &models.DestinationRuleM{}
		err = json.Unmarshal(body, v)
		if err != nil {
			log.Infof("err:[%s]", err)
			return
		}

		// if v.Spec.TrafficPolicy.ConnectionPool.Http.Http1MaxPendingRequests > 2147483647 {
		// 	v.Spec.TrafficPolicy.ConnectionPool.Http.Http1MaxPendingRequests = 2147483647
		// }
		// if v.Spec.TrafficPolicy.ConnectionPool.Http.Http2MaxRequests > 2147483647 {
		// 	v.Spec.TrafficPolicy.ConnectionPool.Http.Http2MaxRequests = 2147483647
		// }

		// if v.Spec.TrafficPolicy.ConnectionPool.Http.MaxRequestsPerConnection > 536870912 {
		// 	v.Spec.TrafficPolicy.ConnectionPool.Http.MaxRequestsPerConnection = 536870912
		// }

		// if v.Spec.TrafficPolicy.ConnectionPool.Http.MaxRetries > 2147483647 {
		// 	v.Spec.TrafficPolicy.ConnectionPool.Http.MaxRetries = 2147483647
		// }

		// if v.Spec.TrafficPolicy.ConnectionPool.Tcp.MaxConnections > 2147483647 {
		// 	v.Spec.TrafficPolicy.ConnectionPool.Tcp.MaxConnections = 2147483647
		// }

		v.Spec.Host = fmt.Sprintf("%s.%s.svc.cluster.local", v.Metadata.Name, v.Metadata.Namespace)
		res, err = json.Marshal(v)
		if err != nil {
			log.Infof("err:[%s]", err)
			return
		}
	}

	fmt.Println("before:", string(body))
	fmt.Printf("after:%+v\n", string(res))

	return
}

func checkObjectType(objectType string) bool {
	return business.GetIstioAPI(objectType) != ""
}

func audit(r *http.Request, message string) {
	if config.Get().Server.AuditLog {
		user := r.Header.Get("Kiali-User")
		log.Infof("AUDIT User [%s] Msg [%s]", user, message)
	}
}

func IstioConfigPermissions(w http.ResponseWriter, r *http.Request) {
	// query params
	params := r.URL.Query()
	namespaces := params.Get("namespaces") // csl of namespaces

	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	istioConfigPermissions := models.IstioConfigPermissions{}
	if len(namespaces) > 0 {
		ns := strings.Split(namespaces, ",")
		istioConfigPermissions = business.IstioConfig.GetIstioConfigPermissions(ns)
	}
	RespondWithJSON(w, http.StatusOK, istioConfigPermissions)
}
