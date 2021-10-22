package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/util"
)

func ServiceInject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	service := params["service"]
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Workloads initialization error: "+err.Error())
		return
	}

	namespace := params["namespace"]
	workloadType := "Deployment"

	serviceDetails, err := business.Svc.GetService(namespace, service, defaultHealthRateInterval, util.Clock.Now())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	jsonPatch := string(`{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject":"true"}}}}}`)
	for _, deploy := range serviceDetails.Workloads {
		_, err1 := business.Workload.UpdateWorkload(namespace, deploy.Name, workloadType, true, jsonPatch)
		if err1 != nil {
			err = errors.New(err.Error() + err1.Error())
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
	service := params["service"]
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Workloads initialization error: "+err.Error())
		return
	}

	namespace := params["namespace"]
	workloadType := "Deployment"

	serviceDetails, err := business.Svc.GetService(namespace, service, defaultHealthRateInterval, util.Clock.Now())
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	jsonPatch := string(`{"spec":{"template":{"metadata":{"annotations":{"sidecar.istio.io/inject":"false"}}}}}`)
	for _, deploy := range serviceDetails.Workloads {
		_, err1 := business.Workload.UpdateWorkload(namespace, deploy.Name, workloadType, true, jsonPatch)
		if err1 != nil {
			err = errors.New(err.Error() + err1.Error())
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
		// TotalCount:  10,
		// PageCount:   1,
		// CurrentPage: 1,
		// PageSize:    10,
		Data: []interface{}{},
	}

	params := mux.Vars(r)
	cluster := r.URL.Query().Get("cluster")
	if _, ok := business.ClusterMap[cluster]; !ok {
		RespondWithJSON(w, http.StatusOK, resp)
		return
	}
	// Get business layer
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	namespace := params["namespace"]

	// Fetch and build services
	serviceList, err := business.Svc.GetServiceList(cluster, namespace, true)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	resp.CurrentPage = 1
	resp.PageCount = 1
	resp.PageSize = 20
	resp.TotalCount = len(serviceList.Services)

	resp.Data = serviceList

	RespondWithJSON(w, http.StatusOK, resp)
}

// ServiceDetails is the API handler to fetch full details of an specific service
func ServiceDetails(w http.ResponseWriter, r *http.Request) {
	// Get business layer
	business, err := getBusiness(r)
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

	serviceDetails, err := business.Svc.GetService(namespace, service, rateInterval, queryTime)
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
	business, err := getBusiness(r)
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
