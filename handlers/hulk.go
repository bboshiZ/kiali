package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kiali/kiali/log"
)

type HulkCluster struct {
	Code   int `json:"code"`
	Result []struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"result"`
}

func GetHulkClusters() (cluster HulkCluster, err error) {
	hulkUrl := "http://scmp-hulk.sgt:80/hulk/clusters?limit=1000"
	resp, err := http.Get(hulkUrl)
	if err != nil {
		log.Error(err)
		return
	}

	defer resp.Body.Close()
	clusterBody, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(clusterBody, &cluster)
	if err != nil {
		log.Error(err)
		return
	}

	if cluster.Code != 0 {
		return cluster, errors.New(string(clusterBody))
	}

	// for i := range cluster.Result {
	// 	cluster.Result[i].Name = strings.ToLower(cluster.Result[i].Name)
	// }
	return
}

type HulkEndpoints struct {
	Code   int `json:"code"`
	Result struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Subsets []struct {
			ports []struct {
				Name     string `json:"name"`
				Port     string `json:"port"`
				Protocol string `json:"protocol"`
			} `json:"ports"`
			Addresses []struct {
				Ip string `json:"ip"`
			} `json:"addresses"`
		} `json:"subsets"`
	} `json:"result"`
}

func GetHulkClusterEndpointsIps(name string, ns string, cid int) (ips []string, err error) {
	hulkUrl := fmt.Sprintf("http://scmp-hulk.sgt:80/hulk/endpoints/%s/name/%s?cid=%d", ns, name, cid)
	// hulkUrl := "http://scmp-hulk.sgt:80/hulk/clusters"
	resp, err := http.Get(hulkUrl)
	if err != nil {
		log.Error(err)
		return
	}

	var ep HulkEndpoints

	defer resp.Body.Close()
	epBody, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(epBody, &ep)
	if err != nil {
		log.Error(err)
		err = errors.New(string(epBody))
		return
	}
	if ep.Code != 0 {
		return nil, errors.New(string(epBody))
	}
	for _, e := range ep.Result.Subsets {
		for _, a := range e.Addresses {
			ips = append(ips, a.Ip)

		}
	}

	return
}
