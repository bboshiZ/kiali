package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/config"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
)

func IstioConfigList(w http.ResponseWriter, r *http.Request) {
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

func IstioVirtualServiceDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := "virtualservices"
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

func IstioVirtualServiceUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := "virtualservices"
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
	jsonPatch := string(body)
	updatedConfigDetails, err := business.IstioConfig.UpdateIstioConfigDetail(api, namespace, objectType, object, jsonPatch)

	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	audit(r, "UPDATE on Namespace: "+namespace+" Type: "+objectType+" Name: "+object+" Patch: "+jsonPatch)
	RespondWithJSON(w, http.StatusOK, updatedConfigDetails)
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

func IstioVirtualServiceCreate(w http.ResponseWriter, r *http.Request) {
	// Feels kinda replicated for multiple functions..
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := "virtualservices"

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

func IstioConfigCreate(w http.ResponseWriter, r *http.Request) {
	// Feels kinda replicated for multiple functions..
	params := mux.Vars(r)
	namespace := params["namespace"]
	objectType := params["object_type"]

	// cluster := r.URL.Query().Get("cluster")
	// if _, ok := business.ClusterMap[cluster]; !ok {
	// 	RespondWithJSON(w, http.StatusOK, "")
	// 	return
	// }

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
		v := &models.VirtualServiceM{}
		err = json.Unmarshal(body, v)
		if err != nil {
			log.Infof("err:[%s]", err)
			return
		}
		// v.Metadata.Namespace = "sample"

		// for i, _ := range v.Spec.Http {
		// 	for j, _ := range v.Spec.Http[i].Route {
		// 		v.Spec.Http[i].Route[j].Destination.Host = v.Spec.Http[i].Route[j].Destination.Host + ".svc.cluster.local"
		// 	}
		// }
		for i := range v.Spec.Http {
			if v.Spec.Http[i].Retries == nil {
				v.Spec.Http[i].Retries = &models.HTTPRetry{
					Attempts: 0,
				}
			}
		}

		res, err = json.Marshal(v)
		if err != nil {
			log.Infof("err:[%s]", err)
			return
		}
	case kubernetes.DestinationRules:
		v := &models.DestinationRuleM{}
		err = json.Unmarshal(body, v)
		if err != nil {
			log.Infof("err:[%s]", err)
			return
		}
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
