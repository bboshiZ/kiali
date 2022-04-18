package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
)

func MeshClusterList(w http.ResponseWriter, r *http.Request) {
	result := map[string][]models.ClusterM{}

	hulkUrl := "http://scmp-hulk.sgt:80/hulk/clusters"
	resp, err := http.Get(hulkUrl)
	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	type Cluster struct {
		Result []struct {
			Id      int    `json:"id"`
			Name    string `json:"name"`
			Address string `json:"address"`
		} `json:"result"`
	}

	var cluster Cluster
	err = json.Unmarshal(body, &cluster)
	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	var inMesh bool
	for _, c := range cluster.Result {
		inMesh = false
		for cName, _ := range business.ClusterMap {
			if c.Name == "shareit-cce-test" {
				break
			}
			if c.Name == cName {
				inMesh = true
				result["inMesh"] = append(result["inMesh"], models.ClusterM{Id: c.Id, Name: c.Name, Address: c.Address})
			}
		}
		if !inMesh {
			result["outMesh"] = append(result["outMesh"], models.ClusterM{Id: c.Id, Name: c.Name, Address: c.Address})
		}
	}

	RespondWithJSON(w, http.StatusOK, result)

}

func ClusterList(w http.ResponseWriter, r *http.Request) {
	// resp := RespList{
	// 	Data: []interface{}{},
	// }

	var data []models.ClusterM
	for c, _ := range business.ClusterMap {
		d := models.ClusterM{
			Name: c,
		}
		data = append(data, d)
	}
	// resp.CurrentPage = 1
	// resp.PageCount = 1
	// resp.PageSize = 20
	// resp.TotalCount = len(data)

	// resp.Data = data

	RespondWithJSON(w, http.StatusOK, data)
}

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
