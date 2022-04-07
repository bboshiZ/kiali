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
)

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
	fmt.Printf("xxxxxx-xxxx-%+v,%+v\n", istioConfigDetails, err)

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

		object = fmt.Sprintf("filter-local-ratelimit-%s", object)
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
					fmt.Printf("xxxxxx-%+v,%+v\n", patch.Patch.Value.TypedConfig.Value.TokenBucket, err)
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

			// if configP, ok := result.EnvoyFilter.Spec.ConfigPatches.([]ConfigPatches); ok {
			// 	fmt.Printf("xxxxxx-%+v\n", configP[0].Patch.Value.TypedConfig.Value.TokenBucket)
			// 	// configP.Patch.Value.TypedConfig.Value.TokenBucket.MaxTokens
			// }
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

// func IstioVirtualServiceDelete(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	namespace := params["namespace"]
// 	objectType := "virtualservices"
// 	object := params["object"]

// 	api := business.GetIstioAPI(objectType)
// 	if api == "" {
// 		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
// 		return
// 	}

// 	// Get business layer
// 	business, err := getBusiness(r)
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
// 		return
// 	}
// 	err = business.IstioConfig.DeleteIstioConfigDetail(api, namespace, objectType, object)
// 	if err != nil {
// 		handleErrorResponse(w, err)
// 		return
// 	} else {
// 		audit(r, "DELETE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object)
// 		RespondWithCode(w, http.StatusOK)
// 	}

// }
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

// func IstioVirtualServiceUpdate(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	namespace := params["namespace"]
// 	objectType := "virtualservices"
// 	object := params["object"]
// 	api := business.GetIstioAPI(objectType)
// 	if api == "" {
// 		RespondWithError(w, http.StatusBadRequest, "Object type not managed: "+objectType)
// 		return
// 	}

// 	// Get business layer
// 	business, err := getBusiness(r)
// 	if err != nil {
// 		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
// 		return
// 	}

// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		RespondWithError(w, http.StatusBadRequest, "Update request with bad update patch: "+err.Error())
// 	}
// 	jsonPatch := string(body)
// 	updatedConfigDetails, err := business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, jsonPatch)

// 	if err != nil {
// 		handleErrorResponse(w, err)
// 		return
// 	}

// 	audit(r, "UPDATE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object+" Patch: "+jsonPatch)
// 	RespondWithJSON(w, http.StatusOK, updatedConfigDetails)
// }

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
	}

	newBody, err := fixRequestBody(objectType, body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "request body error: "+err.Error())
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

const (
	TYPE_ROUTE            = "route"
	TYPE_LOADBALANCE      = "load_balance"
	TYPE_CONNECTIONPOOL   = "connection_pool"
	TYPE_OUTLIERDETECTION = "outlier-detection"
	TYPE_RETRY            = "retry"
	TYPE_RATELIMIT        = "ratelimit"
	TYPE_MIRROR           = "mirror"
	TYPE_LOCALITYLB       = "locality-lbsetting"
	TYPE_FAULTINJECT      = "fault_inject"
	TYPE_SUBNET           = "subset"

	rateLimtTpl = `
{
	"apiVersion": "networking.istio.io/v1alpha3",
	"kind": "EnvoyFilter",
	"metadata": {
		"name": "filter-local-ratelimit-{{.Name}}",
		"namespace": "{{.Namespace}}"
	},
	"spec": {
		"workloadSelector": {
		"labels": {
			"app": "nginx"
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
	rateLimtTpl_aa = `
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: filter-local-ratelimit-{{.Name}}
  namespace: {{.Namespace}}
spec:
  workloadSelector:
    labels:
      app: nginx
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: SIDECAR_INBOUND
        listener:
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.local_ratelimit
          typed_config:
            "@type": type.googleapis.com/udpa.type.v1.TypedStruct
            type_url: type.googleapis.com/envoy.extensions.filters.http.local_ratelimit.v3.LocalRateLimit
            value:
              stat_prefix: http_local_rate_limiter
              token_bucket:
                max_tokens: {{.Max_tokens}}
                tokens_per_fill: {{.Tokens_per_fill}}
                fill_interval: {{.Fill_interval}}
              filter_enabled:
                runtime_key: local_rate_limit_enabled
                default_value:
                  numerator: 100
                  denominator: HUNDRED
              filter_enforced:
                runtime_key: local_rate_limit_enforced
                default_value:
                  numerator: 100
                  denominator: HUNDRED
              response_headers_to_add:
                - append: false
                  header:
                    key: x-local-rate-limit
                    value: 'true'
	
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

		object = fmt.Sprintf("filter-local-ratelimit-%s", object)

		err := business.IstioConfig.DeleteIstioConfigDetail(api, namespace, kubernetes.EnvoyFilters, object)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "delete  ratelimit error: "+err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, "")

		return
	}

	fmt.Println("IstioNetworkConfig-xxx:", namespace, objectType, object)
	result, err := business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	// result.DestinationRule
	// fmt.Println("IstioNetworkConfig-result-xxx:", result.DestinationRule, err)
	if networkType == TYPE_MIRROR {
		vsHttpList := result.VirtualService.Spec.Http.([]interface{})
		var newVsHttpList []interface{}
		for i, _ := range vsHttpList {
			vs := vsHttpList[i].(map[string]interface{})
			vs["mirror"] = nil
			vs["mirrorPercentage"] = nil
			newVsHttpList = append(newVsHttpList, vs)
		}
		result.VirtualService.Spec.Http = newVsHttpList
		jsonPatch, err := json.Marshal(result.VirtualService)

		fmt.Printf("IstioNetworkConfig-result-xxx:%+v,%+v\n", string(jsonPatch), err)
		_, err = business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, string(jsonPatch))
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "IstioNetworkConfigDelete error: "+err.Error())
			return
		}

		RespondWithJSON(w, http.StatusOK, "")
		return
	}
	fmt.Printf("IstioNetworkConfig-result-xxx:%+v,%s\n", result.DestinationRule.Spec, err)
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

	}

	// tp["outlierDetection"] = nil

	result.DestinationRule.Spec.TrafficPolicy = tp
	// fmt.Printf("IstioNetworkConfig-result-xxx:%+v\n", result.DestinationRule.Spec)

	jsonPatch, err := json.Marshal(result.DestinationRule)
	fmt.Printf("IstioNetworkConfig-result-xxx:%+v,%+v\n", string(jsonPatch), err)

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

		p := limit{
			Name:      dstRule.Metadata.Name,
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

		t := template.New("test")
		t = template.Must(t.Parse(rateLimtTpl))

		var buf bytes.Buffer
		err := t.Execute(&buf, p)
		if err != nil {
			log.Infof("err-111:[%s]", err)
			handleErrorResponse(w, err)
			return
		}

		body = buf.Bytes()
		object = fmt.Sprintf("filter-local-ratelimit-%s", dstRule.Metadata.Name)

		log.Infof("body-111:[%s]", string(body))

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

	// business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	fmt.Println("IstioNetworkConfig-xxx:", namespace, objectType, object)
	// result, err := business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)

	result, err := business.IstioConfig.GetIstioConfigDetails(namespace, objectType, object)
	// fmt.Println("IstioNetworkConfig-result-xxx:", result, err)
	fmt.Printf("IstioNetworkConfig-result-xxx:%+v,%+v\n", result, err)

	if err == nil {
		if networkType == "mirror" {
			fmt.Printf("IstioNetworkConfig-mirror-xxx:%+v,%+v\n", result.VirtualService.Spec.Http, err)

			vsHttpList := result.VirtualService.Spec.Http.([]interface{})

			if len(vsHttpList) > 0 {
				vsHttp := vsHttpList[0].(map[string]interface{})
				v := &models.VirtualService{}
				jsonErr := json.Unmarshal(body, v)
				if err != nil {
					log.Infof("err:[%s]", jsonErr)
					handleErrorResponse(w, jsonErr)
				}
				newVsHttpList := v.Spec.Http.([]interface{})
				if len(newVsHttpList) > 0 {
					newVsHttp := newVsHttpList[0].(map[string]interface{})
					vsHttp["mirror"] = newVsHttp["mirror"]
					vsHttp["mirrorPercentage"] = newVsHttp["mirrorPercentage"]
					result.VirtualService.Spec.Http = []map[string]interface{}{vsHttp}
					fmt.Printf("IstioNetworkConfig-VirtualService-xxx:%+v\n", result.VirtualService.Spec.Http)

					body, _ = json.Marshal(result.VirtualService)
				}

			}

		}

		fmt.Printf("IstioNetworkConfig-updatemirror-xxx:%+v\n", string(body))

		jsonPatch := string(body)
		updatedConfigDetails, err := business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, jsonPatch)
		if err != nil {
			handleErrorResponse(w, err)
			return
		}
		audit(r, "UPDATE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object+" Patch: "+jsonPatch)
		RespondWithJSON(w, http.StatusOK, updatedConfigDetails)
	} else {
		if networkType == "mirror" {
			vs := &models.VirtualService{}
			jsonErr := json.Unmarshal(body, vs)
			if jsonErr != nil {
				log.Infof("err:[%s]", jsonErr)
				handleErrorResponse(w, jsonErr)
			}

			vs.Spec.Hosts = []string{fmt.Sprintf("%s.%s.svc.cluster.local", vs.Metadata.Name, vs.Metadata.Namespace)}
			newVsHttpList := vs.Spec.Http.([]interface{})
			if len(newVsHttpList) > 0 {
				newVsHttp := newVsHttpList[0].(map[string]interface{})
				d := models.HTTPRouteDestination{
					Destination: &models.Destination{
						Host: fmt.Sprintf("%s.%s.svc.cluster.local", vs.Metadata.Name, vs.Metadata.Namespace),
					},
				}
				newVsHttp["route"] = []models.HTTPRouteDestination{d}
				vs.Spec.Http = []interface{}{newVsHttp}
				body, _ = json.Marshal(vs)
				fmt.Printf("create-mirror-xxx:%+v\n", vs)
			}

		}

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
