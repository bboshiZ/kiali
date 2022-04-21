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

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/util"
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
	business, err := getBusiness(r)
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
			// RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
			// return
			continue
		}
		nodeList, err := business.Svc.GetAllNode()

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
	// if _, found := query["validate"]; found {
	// 	includeValidations = true
	// }

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

func IstioConfigDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]
	object := params["object"]
	clusterMap := business.ClusterMap

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
	business, err := getBusiness(r)
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
	// fmt.Printf("xxxxxx-xxxx-%+v,%+v\n", istioConfigDetails, err)

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

		object = fmt.Sprintf("filter-mirror-%s", object)
		result, err := business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, object)
		if err == nil && result.EnvoyFilter != nil {
			if configP, ok := result.EnvoyFilter.Spec.ConfigPatches.([]interface{}); ok {
				date, err := json.Marshal(configP[0])
				if err == nil {
					var patch MirrorConfigPatches
					err = json.Unmarshal(date, &patch)
					if err == nil {
						if httpList, ok := istioConfigDetails.VirtualService.Spec.Http.([]interface{}); ok {
							// fmt.Printf("xxxxxx-httpList-%+v\n", httpList)
							if len(httpList) > 0 {
								vsHttp := httpList[0].(map[string]interface{})
								// vsHttp["mirror"] = []interface{}{}
								hCluster, err := GetHulkClusters()
								if err != nil {
									log.Errorf("GetHulkClusters err:[%s]", err)
									RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
									return
								}
								mirrorConfig := []interface{}{}
								cid := 0
								for _, mirror := range patch.Patch.Value.Route.RequestMirrorPolicies {
									// fmt.Printf("xxxxxx-mirror-%+v\n", mirror)
									k8sClusterType := "outMesh"
									s := strings.Split(mirror.ServiceInfo, "|")
									// label := result.EnvoyFilter.Metadata.Labels
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
										// RespondWithError(w, http.StatusInternalServerError, "cid or targetPort error")
										// return
									}
									mirrorConfig = append(mirrorConfig, map[string]interface{}{
										"clusterType":      k8sClusterType,
										"cluster":          s[3],
										"cid":              cid,
										"namespace":        s[2],
										"service":          s[1],
										"targetPort":       tPort,
										"mirrorPercentage": float64(mirror.RuntimeFraction.DefaultValue.Numerator) / 10000,
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

		object = fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, object)
		result, err := business.IstioConfig.GetIstioConfigDetails(namespace, kubernetes.EnvoyFilters, object)
		// fmt.Printf("xxxxxx-aaa-%+v,%+v\n", result, err)
		if err == nil && result.EnvoyFilter != nil {
			// fmt.Printf("xxxxxx-bbb-%+v\n", result.EnvoyFilter.Spec.ConfigPatches)
			if configP, ok := result.EnvoyFilter.Spec.ConfigPatches.([]interface{}); ok {
				date, err := json.Marshal(configP[0])
				if err == nil {
					var patch ConfigPatches
					err = json.Unmarshal(date, &patch)
					bucket := patch.Patch.Value.TypedConfig.Value.TokenBucket
					// fmt.Printf("xxxxxx-%+v,%+v\n", patch.Patch.Value.TypedConfig.Value.TokenBucket, err)
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
	business, err := getBusiness(r)
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
	business, err := getBusiness(r)
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
	TYPE_LOADBALANCE           = "load_balance"
	TYPE_CONNECTIONPOOL        = "connection_pool"
	TYPE_OUTLIERDETECTION      = "outlier-detection"
	TYPE_RETRY                 = "retry"
	TYPE_RATELIMIT             = "ratelimit"
	TYPE_MIRROR                = "mirror"
	TYPE_LOCALITYLB            = "locality-lbsetting"
	TYPE_FAULTINJECT           = "fault-inject"
	TYPE_SLOWSTART             = "slow-start"
	TYPE_SUBNET                = "subset"
	reteLimitEnvoyFilterPrefix = "filter-local-ratelimit-"

	serviceEntryTpl = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "ServiceEntry",
	"metadata": {
		"name": "{{.Name}}",
		"namespace": "{{.Namespace}}",
		"labels": {
			"mirror": "{{.MirrorName}}"
		}
	},
	"spec": {
		"hosts": [
		"{{.Host}}"
		],
		"ports": [
		{
			"number": {{.Port}},
			"name": "http",
			"protocol": "HTTP"
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

	mirrorTpl = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "filter-mirror-{{.Name}}",
		"namespace": "{{.Namespace}}"
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
			"app": "{{.Name}}"
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

func IstioNetworkConfigDelete(w http.ResponseWriter, r *http.Request) {
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

	business, err := getBusiness(r)
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
		filterName := fmt.Sprintf("filter-mirror-%s", object)
		err := business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.EnvoyFilters, filterName)
		if err != nil {
			log.Infof("err:[%s]", err)
			// RespondWithError(w, http.StatusInternalServerError, "delete  ratelimit error: "+err.Error())
			// return
		}
		label := fmt.Sprintf("mirror=mirror-%s-%s", object, namespace)
		objects, err := business.IstioConfig.GetIstioObject(api, namespace, kubernetes.ServiceEntries, label)
		if err == nil {
			for _, ob := range objects {
				err = business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.ServiceEntries, ob.GetObjectMeta().Name)
				if err != nil {
					log.Infof("err:[%s]", err)
				}
			}
		}
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

func IstioNetworkConfigUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]
	networkType := params["network_type"]
	object := params["object"]
	clusterMap := business.ClusterMap
	cidStr := r.URL.Query().Get("cid")
	sourceCid, _ := strconv.Atoi(cidStr)
	api := business.GetIstioAPI(objectType)
	if api == "" {
		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
		return
	}

	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Create request could not be read: "+err.Error())
	}

	if networkType == TYPE_MIRROR {

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
			}
		}

		targerService, err := business.Svc.GetService(cName, dstRule.Metadata.Namespace, dstRule.Metadata.Name, defaultHealthRateInterval, util.Clock.Now())
		if err != nil {
			log.Infof("err:[%s]", err)
			handleErrorResponse(w, err)
			return
		}
		mConfig := Mirror{
			Name:      dstRule.Metadata.Name,
			Namespace: dstRule.Metadata.Namespace,
			Host:      fmt.Sprintf("%s.%s.svc.cluster.local:%d", dstRule.Metadata.Name, dstRule.Metadata.Namespace, targerService.Service.Ports[0].Port),
		}

		label := fmt.Sprintf("mirror=mirror-%s-%s", object, namespace)
		objects, err := business.IstioConfig.GetIstioObject(api, namespace, kubernetes.ServiceEntries, label)
		// fmt.Println("xxxxx-", objects, err)
		if err == nil {
			for _, ob := range objects {
				err = business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.ServiceEntries, ob.GetObjectMeta().Name)
				if err != nil {
					log.Infof("err:[%s]", err)
				}
			}
		}

		for _, m := range dstRule.Spec.Http[0].Mirror {

			inMesh := false
			if m.TargetPort <= 0 {
				log.Errorf("m.TargetPort <= 0:[%s]", m.TargetPort)
				RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("TargetPort err:%d", m.TargetPort))
				return
			}

			for cName := range clusterMap {
				if m.Cluster == "shareit-cce-test" {
					break
				}
				if m.Cluster == cName {
					inMesh = true
					mirrorService, err := business.Svc.GetService(m.Cluster, m.Namespace, m.Service, defaultHealthRateInterval, util.Clock.Now())
					if err != nil {
						log.Infof("err:[%s]", err)
						handleErrorResponse(w, err)
						return
					}
					mConfig.Array = append(mConfig.Array, Service{
						MirrorCluster: fmt.Sprintf("outbound|%d||%s.%s.svc.cluster.local", mirrorService.Service.Ports[0].Port, m.Service, m.Namespace),
						Numerator:     strconv.FormatInt(int64(m.MirrorPercentage*10000), 10),
						ServiceInfo:   fmt.Sprintf("%d|%s|%s|%s", m.TargetPort, m.Service, m.Namespace, m.Cluster),
					})
					inMesh = true
					break
					// result["inMesh"] = append(result["inMesh"], models.ClusterM{Id: c.Id, Name: c.Name, Address: c.Address})
				}
			}

			if inMesh {
				continue
			}

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
						Name       string   `json:"name"`
						Namespace  string   `json:"namespace"`
						Address    []string `json:"address"`
						Host       string   `json:"host"`
						Port       string   `json:"port"`
						MirrorName string   `json:"mirror_name"`
						// MirrorCluster string `json:"mirror_cluster"`
						// Numerator     string `json:"numerator"`
					}

					snName := fmt.Sprintf("mirror-%s-to-%s", object, m.Service)
					sn := ServiceEn{
						Name:       snName,
						Namespace:  namespace,
						Host:       fmt.Sprintf("mirror-%s-to-%s.ushareit", object, m.Service),
						Port:       strconv.Itoa(m.TargetPort),
						Address:    podIps,
						MirrorName: fmt.Sprintf("mirror-%s-%s", object, namespace),
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
					// object = fmt.Sprintf("mirror-entry-%s", m.Service)

					log.Infof("body-111:[%s]", string(body))
					business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.ServiceEntries, snName)
					// if err != nil {
					// 	log.Infof("err:[%s]", err)
					// 	handleErrorResponse(w, err)
					// 	return
					// }
					_, err = business.IstioConfig.CreateIstioConfigDetail(api, namespace, kubernetes.ServiceEntries, body)
					if err != nil {
						handleErrorResponse(w, err)
						return
					}
					// appFound = true
					// break

					// appUrl := fmt.Sprintf("http://scmp-hulk.sgt:80/hulk/api/v2/apps/%s/all?limit=10000&page=1", m.Namespace)
					// req, err := http.NewRequest("GET", appUrl, nil)
					// if err != nil {
					// 	log.Infof("err:[%s]", err)
					// 	handleErrorResponse(w, err)
					// 	return
					// }
					// req.Header.Set("cid", strconv.Itoa(m.Cid))
					// client := &http.Client{}
					// resp, err := client.Do(req)
					// if err != nil {
					// 	log.Errorf("err:[%s]", err)
					// 	handleErrorResponse(w, err)
					// 	return
					// }
					// defer resp.Body.Close()
					// appBody, _ := ioutil.ReadAll(resp.Body)
					// log.Infof("mirror-xxx:[%+v][%+v]", string(appBody), err)

					// type AppInfo struct {
					// 	Code   int `json:"code"`
					// 	Result []struct {
					// 		Id     int               `json:"id"`
					// 		Name   string            `json:"name"`
					// 		Type   int               `json:"type"`
					// 		Labels map[string]string `json:"labels"`
					// 	} `json:"result"`
					// }

					// var app AppInfo
					// err = json.Unmarshal(appBody, &app)
					// if err != nil {
					// 	log.Error(err)
					// 	RespondWithError(w, http.StatusInternalServerError, err.Error())
					// 	return
					// }
					// if app.Code != 0 {
					// 	RespondWithError(w, http.StatusInternalServerError, "Get apps error:"+string(appBody))
					// 	return
					// }

					// var appType, podUrl string
					// var podIps []string
					// var appFound bool
					// for _, a := range app.Result {
					// 	if a.Name == m.Service {
					// 		if a.Type == 0 {
					// 			appType = "deployment"
					// 			podUrl = fmt.Sprintf("http://scmp-hulk.sgt:80/hulk/api/v2/apps/%s/%s/%s/pods?cid=%d", m.Namespace, appType, a.Name, m.Cid)

					// 		} else if a.Type == 2 {
					// 			appType = "cloneset"
					// 			podUrl = fmt.Sprintf("http://scmp-hulk.sgt:80/hulk/api/v2/apps/%s/%s/%s/pods?cid=%d", m.Namespace, appType, a.Name, m.Cid)

					// 		}

					// 		resp, err := http.Get(podUrl)
					// 		if err != nil {
					// 			log.Error(err)
					// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
					// 			return
					// 		}

					// 		defer resp.Body.Close()
					// 		podBody, _ := ioutil.ReadAll(resp.Body)

					// 		type PodInfo struct {
					// 			Result []struct {
					// 				Id    int    `json:"id"`
					// 				PodIp string `json:"pod_ip"`
					// 			} `json:"result"`
					// 		}

					// 		var pod PodInfo
					// 		err = json.Unmarshal(podBody, &pod)
					// 		if err != nil {
					// 			log.Error(err)
					// 			RespondWithError(w, http.StatusInternalServerError, err.Error())
					// 			return
					// 		}
					// 		for _, p := range pod.Result {
					// 			podIps = append(podIps, p.PodIp)
					// 		}

					// 		type ServiceEn struct {
					// 			Name       string   `json:"name"`
					// 			Namespace  string   `json:"namespace"`
					// 			Address    []string `json:"address"`
					// 			Host       string   `json:"host"`
					// 			Port       string   `json:"port"`
					// 			MirrorName string   `json:"mirror_name"`
					// 		}

					// 		if m.TargetPort <= 0 {
					// 			// m.TargetPort = 80
					// 			log.Errorf("m.TargetPort <= 0:[%s]", m.TargetPort)
					// 			msg := fmt.Sprintf("TargetPort err:%d", m.TargetPort)
					// 			RespondWithError(w, http.StatusBadRequest, msg)
					// 			return
					// 		}

					// 		snName := fmt.Sprintf("mirror-%s-to-%s", object, m.Service)
					// 		sn := ServiceEn{
					// 			Name:       snName,
					// 			Namespace:  namespace,
					// 			Host:       fmt.Sprintf("mirror-%s-to-%s.ushareit", object, m.Service),
					// 			Port:       strconv.Itoa(m.TargetPort),
					// 			Address:    podIps,
					// 			MirrorName: fmt.Sprintf("mirror-%s-%s", object, namespace),
					// 		}

					// 		t, err := template.New("test").Funcs(tplFuncMap).Parse(serviceEntryTpl)
					// 		if err != nil {
					// 			log.Infof("template:[%s]", err)
					// 			handleErrorResponse(w, err)
					// 			return
					// 		}
					// 		var buf bytes.Buffer
					// 		err = t.Execute(&buf, sn)
					// 		if err != nil {
					// 			log.Errorf("err:[%s]", err)
					// 			handleErrorResponse(w, err)
					// 			return
					// 		}

					// 		body = buf.Bytes()
					// 		log.Infof("body-111:[%s]", string(body))
					// 		business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.ServiceEntries, snName)
					// 		_, err = business.IstioConfig.CreateIstioConfigDetail(api, namespace, kubernetes.ServiceEntries, body)
					// 		if err != nil {
					// 			handleErrorResponse(w, err)
					// 			return
					// 		}
					// 		appFound = true
					// 		break
					// 	}
					// }

					// if !appFound {
					// 	RespondWithError(w, http.StatusBadRequest, "App not found")
					// 	return

					// }

					//Request URL: https://scmp.ushareit.me/hulk/api/v2/apps/sprs/deployment/store-house/pods
					// podUrl := fmt.Sprintf("http://scmp-hulk.sgt:80/hulk/api/v2/apps/%s/%s/%s/pods", m.Namespace, appType, m.Service)

					mConfig.Array = append(mConfig.Array, Service{
						MirrorCluster: fmt.Sprintf("outbound|%d||mirror-%s-to-%s.ushareit", m.TargetPort, object, m.Service),
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
		object = fmt.Sprintf("filter-mirror-%s", dstRule.Metadata.Name)

	}
	if networkType == TYPE_RATELIMIT {
		objectType = kubernetes.EnvoyFilters
		dstRule := &models.DestinationRuleM{}
		err = json.Unmarshal(body, dstRule)
		if err != nil {
			log.Infof("err:[%s]", err)
			handleErrorResponse(w, err)
			return
		}

		type limit struct {
			Name          string `json:name`
			Namespace     string `json:namespace`
			MaxTokens     int    `json:max_tokens`
			TokensPerFill int    `json:tokens_per_fill`
			FillInterval  string `json:fill_interval`
		}
		trafic := dstRule.Spec.TrafficPolicy.(map[string]interface{})
		li := trafic["rateLimit"].(map[string]interface{})
		object = fmt.Sprintf("%s%s", reteLimitEnvoyFilterPrefix, dstRule.Metadata.Name)
		p := limit{
			Name:      object,
			Namespace: dstRule.Metadata.Namespace,
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

	}

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

	if objectType == kubernetes.DestinationRules {
		v := &models.DestinationRuleM{}
		err = json.Unmarshal(body, v)
		if err != nil {
			log.Infof("err:[%s]", err)
			handleErrorResponse(w, err)
			return
		}
		v.Spec.Host = fmt.Sprintf("%s.%s.svc.cluster.local", v.Metadata.Name, v.Metadata.Namespace)
		body, err = json.Marshal(v)
		if err != nil {
			log.Infof("err:[%s]", err)
			handleErrorResponse(w, err)
			return
		}
	}

	_, err = business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	if err == nil {
		// fmt.Printf("IstioNetworkConfig-updatemirror-xxx:%+v\n", string(body))

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
