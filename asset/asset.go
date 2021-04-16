package asset

import (
	"fmt"
	"github.com/cloudinary/cloudinary-go/api"
	"github.com/cloudinary/cloudinary-go/config"
	"github.com/cloudinary/cloudinary-go/internal/signature"
	"github.com/cloudinary/cloudinary-go/transformation"
	"hash/crc32"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Asset is the Asset struct.
type Asset struct {
	AssetType      api.AssetType
	DeliveryType   api.DeliveryType
	Transformation transformation.RawTransformation
	Version        int
	PublicID       string
	Suffix         string
	Config         config.Configuration
}

func New(publicID string, conf *config.Configuration) (*Asset, error) {
	if conf == nil {
		var err error
		conf, err = config.New()
		if err != nil {
			return nil, err
		}
	}

	return &Asset{PublicID: publicID, Config: *conf}, nil
}

func Image(publicID string, conf *config.Configuration) (*Asset, error) {
	return New(publicID, conf)
}

func Video(publicID string, conf *config.Configuration) (*Asset, error) {
	v, err := New(publicID, conf)
	if err != nil {
		return nil, err
	}
	v.AssetType = api.Video

	return v, nil
}

func File(publicID string, conf *config.Configuration) (*Asset, error) {
	f, err := New(publicID, conf)
	if err != nil {
		return nil, err
	}
	f.AssetType = api.File

	return f, nil
}

func Media(publicID string, conf *config.Configuration) (*Asset, error) {
	return New(publicID, conf)
}

// String serializes Asset to string.
func (a Asset) String() (string, error) {
	return joinUrl([]interface{}{a.distribution(), a.assetType(), a.signature(), a.Transformation, a.version(), a.PublicID}), nil
}

// version finalizes the version part (v123) of the asset URL.
func (a Asset) version() string {
	var versionRegexp = regexp.MustCompile(`^v\d+`)
	var urlRegexp = regexp.MustCompile(`^https?://`)
	version := a.Version

	if version == 0 &&
		a.Config.URL.ForceVersion &&
		filepath.Dir(a.PublicID) != "." &&
		!urlRegexp.MatchString(a.PublicID) &&
		!versionRegexp.MatchString(a.PublicID) {
		version = 1
	}

	if version != 0 {
		return fmt.Sprintf("v%d", version)
	}

	return ""
}

// distribution builds the hostname for the asset distribution.
//
//   1) Customers in shared distribution (e.g. res.cloudinary.com)
//      If CDNSubDomain is true uses res-[1-5].cloudinary.com for both http and https.
//      Setting secureCDNSubDomain to false disables this for https.
//   2) Customers with private cdn
//      If CDNSubDomain is true uses cloudName-res-[1-5].cloudinary.com for http
//      If secureCDNSubDomain is true uses cloudName-res-[1-5].cloudinary.com for https
//      (please contact support if you require this)
//   3) Customers with cname
//      If CDNSubDomain is true uses a[1-5].cname for http.
//      For https, uses the same naming scheme as 1 for shared distribution and as 2 for private distribution.
//
func (a Asset) distribution() string {
	uc := a.Config.URL
	useSharedHost := !uc.PrivateCDN
	var hostName string
	if uc.Secure {
		hostName = uc.SecureCName
		if hostName == "" {
			if uc.PrivateCDN {
				hostName = buildHostName(a.Config.Cloud.CloudName, "", uc.SubDomain, uc.Domain)
			} else {
				hostName = uc.SharedHost
				useSharedHost = true
			}
		}

		secureCDNSubDomain := uc.SecureCDNSubDomain
		if useSharedHost && !secureCDNSubDomain {
			secureCDNSubDomain = uc.CDNSubDomain
		}

		if secureCDNSubDomain {
			hostName = strings.Replace(hostName, uc.SharedHost, buildHostName("", domainShard(a.PublicID), uc.SubDomain, uc.Domain), 1)
		}
	} else {
		if uc.CName != "" {
			subDomain := ""
			if uc.CDNSubDomain {
				subDomain = "a" + domainShard(a.PublicID)
			}
			hostName = buildHostName("", "", subDomain, uc.CName)
		} else {
			prefix := ""
			if uc.PrivateCDN {
				prefix = a.Config.Cloud.CloudName
			}
			suffix := ""
			if uc.CDNSubDomain {
				suffix = domainShard(a.PublicID)
			}
			hostName = buildHostName(prefix, suffix, uc.SubDomain, uc.Domain)
		}
	}

	distribution := fmt.Sprintf("%s://%s", uc.Protocol(), hostName)

	if useSharedHost {
		distribution += "/" + a.Config.Cloud.CloudName
	}

	return distribution
}

// buildHostName is a helper method for building hostname of form:
// 		subDomainPrefix-subDomain-subDomainSuffix.domain
// For example:
//      cloudName-res-3.cloudinary.com
func buildHostName(subDomainPrefix string, subDomainSuffix string, subDomain string, domain string) string {
	return joinNonEmpty([]interface{}{joinNonEmpty([]interface{}{subDomainPrefix, subDomain, subDomainSuffix}, "-"), domain}, ".")
}

var crc32q = crc32.MakeTable(crc32.IEEE)

// domainShard computes the domain shard from the source.
func domainShard(source string) string {
	return strconv.Itoa(int(crc32.Checksum([]byte(source), crc32q)%5 + 1))
}

var suffixSupportedDeliveryTypes = map[api.AssetType]map[api.DeliveryType]string{
	api.Image: {
		api.Upload: "images",
		api.Private: "private_images",
		api.Authenticated: "authenticated_images",
	},
	api.Video: {
		api.Upload: "videos",
	},
	api.File: {
		api.Upload: "files",
	},
}

func (a Asset)assetType() string {
	if a.Suffix == "" {
		return joinUrl([]interface{}{a.AssetType, a.DeliveryType})
	}

	if a.AssetType
}

// signature returns URL signature.
//
// https://cloudinary.com/documentation/advanced_url_delivery_options#generating_delivery_url_signatures
func (a Asset) signature() string {
	if !a.Config.URL.SignURL {
		return ""
	}

	algo, length := a.getSignatureAlgorithmAndLength()

	toSign := joinUrl([]interface{}{a.Transformation, a.PublicID})

	return signature.SignURL(toSign, a.Config.Cloud.APISecret, algo, length)
}

func (a Asset) getSignatureAlgorithmAndLength() (signature.Algo, signature.Length) {
	if a.Config.URL.GetSignatureLength() == signature.Long {
		return signature.SHA256, signature.Long
	}

	return a.Config.Cloud.GetSignatureAlgorithm(), a.Config.URL.GetSignatureLength()
}

func joinNonEmpty(items []interface{}, sep string) string {
	var parts []string
	for _, i := range items {
		s := fmt.Sprintf("%v", i)
		if strings.TrimSpace(s) != "" {
			parts = append(parts, s)
		}
	}

	return strings.Join(parts, sep)
}

func joinUrl(items []interface{}) string {
	return joinNonEmpty(items, "/")
}
