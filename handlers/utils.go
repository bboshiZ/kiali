package handlers

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/kiali/kiali/business"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/prometheus"
)

type promClientSupplier func() (*prometheus.Client, error)

var defaultPromClientSupplier = prometheus.NewClient

func checkNamespaceAccess(nsServ business.NamespaceService, namespace string) (*models.Namespace, error) {
	if nsInfo, err := nsServ.GetNamespace(namespace); err != nil {
		return nil, err
	} else {
		return nsInfo, nil
	}
}

func createMetricsServiceForNamespace(w http.ResponseWriter, r *http.Request, promSupplier promClientSupplier, namespace string) (*business.MetricsService, *models.Namespace) {
	metrics, infoMap := createMetricsServiceForNamespaces(w, r, promSupplier, []string{namespace})
	if result, ok := infoMap[namespace]; ok {
		if result.err != nil {
			RespondWithError(w, http.StatusForbidden, "Cannot access namespace data: "+result.err.Error())
			return nil, nil
		}
		return metrics, result.info
	}
	return nil, nil
}

type nsInfoError struct {
	info *models.Namespace
	err  error
}

func createMetricsServiceForNamespaces(w http.ResponseWriter, r *http.Request, promSupplier promClientSupplier, namespaces []string) (*business.MetricsService, map[string]nsInfoError) {
	layer, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return nil, nil
	}
	prom, err := promSupplier()
	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusServiceUnavailable, "Prometheus client error: "+err.Error())
		return nil, nil
	}

	nsInfos := make(map[string]nsInfoError)
	for _, ns := range namespaces {
		info, err := checkNamespaceAccess(layer.Namespace, ns)
		nsInfos[ns] = nsInfoError{info: info, err: err}
	}
	metrics := business.NewMetricsService(prom)
	return metrics, nsInfos
}

// getAuthInfo retrieves the token from the request's context
func getAuthInfo(r *http.Request) (*api.AuthInfo, error) {
	authInfoContext := r.Context().Value("authInfo")
	if authInfoContext != nil {
		if authInfo, ok := authInfoContext.(*api.AuthInfo); ok {
			return authInfo, nil
		} else {
			return nil, errors.New("authInfo is not of type *api.AuthInfo")
		}
	} else {
		return nil, errors.New("authInfo missing from the request context")
	}
}

// getBusiness returns the business layer specific to the users's request
func getBusiness(r *http.Request) (*business.Layer, error) {
	authInfo, err := getAuthInfo(r)
	if err != nil {
		return nil, err
	}

	return business.Get(authInfo)
}

func getBusinessByCluster(r *http.Request) (*business.Layer, error) {
	idStr := r.Header.Get("cid")
	cid, _ := strconv.Atoi(idStr)
	if cid == 0 {
		return nil, errors.New("cid not set")
	}
	hCluster, err := GetHulkClusters()
	if err != nil {
		return nil, err
	}
	for _, hc := range hCluster.Result {
		if cid == hc.Id {
			// svcClusterName = hc.Name
			cName := strings.ToLower(hc.Name)
			if c, ok := business.IstioPrimary[cName]; ok {
				return business.GetByCluster(business.IstioPrimary[c])
			}

		}
	}

	return nil, errors.New("cluster not found")
}

func SlicePage(page, pageSize, nums int) (sliceStart, sliceEnd, pageCount int) {
	if page < 0 {
		page = 1
	}

	if pageSize < 0 {
		pageSize = 20
	}

	if pageSize > nums {
		return 0, nums, 1
	}

	// 总页数
	pageCount = int(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0, 1
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize

	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd, pageCount
}

const (
	DNS1123LabelMaxLength = 63 // Public for testing only.
	dns1123LabelFmt       = "[a-zA-Z0-9](?:[-a-zA-Z0-9]*[a-zA-Z0-9])?"
	// a wild-card prefix is an '*', a normal DNS1123 label with a leading '*' or '*-', or a normal DNS1123 label
	wildcardPrefix = `(\*|(\*|\*-)?` + dns1123LabelFmt + `)`

	// Using kubernetes requirement, a valid key must be a non-empty string consist
	// of alphanumeric characters, '-', '_' or '.', and must start and end with an
	// alphanumeric character (e.g. 'MyValue',  or 'my_value',  or '12345'
	qualifiedNameFmt = "(?:[A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]"

	// In Kubernetes, label names can start with a DNS name followed by a '/':
	dnsNamePrefixFmt       = dns1123LabelFmt + `(?:\.` + dns1123LabelFmt + `)*/`
	dnsNamePrefixMaxLength = 253
)

var (
	tagRegexp            = regexp.MustCompile("^(" + dnsNamePrefixFmt + ")?(" + qualifiedNameFmt + ")$") // label value can be an empty string
	labelValueRegexp     = regexp.MustCompile("^" + "(" + qualifiedNameFmt + ")?" + "$")
	dns1123LabelRegexp   = regexp.MustCompile("^" + dns1123LabelFmt + "$")
	wildcardPrefixRegexp = regexp.MustCompile("^" + wildcardPrefix + "$")
)

// encapsulates DNS 1123 checks common to both wildcarded hosts and FQDNs
func checkDNS1123Preconditions(name string) error {
	if len(name) > 255 {
		return fmt.Errorf("domain name %q too long (max 255)", name)
	}
	if len(name) == 0 {
		return fmt.Errorf("empty domain name not allowed")
	}
	return nil
}

func IsWildcardDNS1123Label(value string) bool {
	return len(value) <= DNS1123LabelMaxLength && wildcardPrefixRegexp.MatchString(value)
}

// IsDNS1123Label tests for a string that conforms to the definition of a label in
// DNS (RFC 1123).
func IsDNS1123Label(value string) bool {
	return len(value) <= DNS1123LabelMaxLength && dns1123LabelRegexp.MatchString(value)
}

func validateDNS1123Labels(domain string) error {
	parts := strings.Split(domain, ".")
	topLevelDomain := parts[len(parts)-1]
	if _, err := strconv.Atoi(topLevelDomain); err == nil {
		return fmt.Errorf("domain name %q invalid (top level domain %q cannot be all-numeric)", domain, topLevelDomain)
	}
	for i, label := range parts {
		// Allow the last part to be empty, for unambiguous names like `istio.io.`
		if i == len(parts)-1 && label == "" {
			return nil
		}
		if !IsDNS1123Label(label) {
			return fmt.Errorf("domain name %q invalid (label %q invalid)", domain, label)
		}
	}
	return nil
}

// ValidateWildcardDomain checks that a domain is a valid FQDN, but also allows wildcard prefixes.
func ValidateWildcardDomain(domain string) error {
	if err := checkDNS1123Preconditions(domain); err != nil {
		return err
	}
	// We only allow wildcards in the first label; split off the first label (parts[0]) from the rest of the host (parts[1])
	parts := strings.SplitN(domain, ".", 2)
	if !IsWildcardDNS1123Label(parts[0]) {
		return fmt.Errorf("domain name %q invalid (label %q invalid)", domain, parts[0])
	} else if len(parts) > 1 {
		return validateDNS1123Labels(parts[1])
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

// 用rand.Int63()替换rand.Intn()
func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
