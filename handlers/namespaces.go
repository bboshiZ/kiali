package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
)

func NamespaceList(w http.ResponseWriter, r *http.Request) {
	resp := RespList{
		// TotalCount:  10,
		// PageCount:   1,
		// CurrentPage: 1,
		// PageSize:    10,
		Data: []interface{}{},
	}

	cluster := r.URL.Query().Get("cluster")
	if _, ok := business.ClusterMap[cluster]; !ok {
		RespondWithJSON(w, http.StatusOK, resp)
		return
	}

	business, err := getBusiness(r)
	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	namespaces, err := business.Namespace.GetRemoteNamespaces(cluster)
	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.CurrentPage = 1
	resp.PageCount = 1
	resp.PageSize = 20
	resp.TotalCount = len(namespaces)

	resp.Data = namespaces

	RespondWithJSON(w, http.StatusOK, namespaces)
}

// NamespaceValidationSummary is the API handler to fetch validations summary to be displayed.
// It is related to all the Istio Objects within the namespace
func NamespaceValidationSummary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]

	business, err := getBusiness(r)
	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var validationSummary models.IstioValidationSummary

	istioConfigValidationResults, errValidations := business.Validations.GetValidations(namespace, "")
	if errValidations != nil {
		log.Error(errValidations)
		RespondWithError(w, http.StatusInternalServerError, errValidations.Error())
	} else {
		validationSummary = istioConfigValidationResults.SummarizeValidation(namespace)
	}

	RespondWithJSON(w, http.StatusOK, validationSummary)
}

// NamespaceUpdate is the API to perform a patch on a Namespace configuration
func NamespaceUpdate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Namespace initialization error: "+err.Error())
		return
	}
	namespace := params["namespace"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Update request with bad update patch: "+err.Error())
	}
	jsonPatch := string(body)

	ns, err := business.Namespace.UpdateNamespace(namespace, jsonPatch)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}
	audit(r, "UPDATE on Namespace: "+namespace+" Patch: "+jsonPatch)
	RespondWithJSON(w, http.StatusOK, ns)
}
