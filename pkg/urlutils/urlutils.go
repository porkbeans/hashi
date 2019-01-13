package urlutils

import (
	"fmt"
)

const (
	HashicorpProductList = "https://releases.hashicorp.com/"
)

func ProductVersionListUrl(product string) string {
	return fmt.Sprintf("%s%s/", HashicorpProductList, product)
}

func ProductBuildListUrl(product string, version string) string {
	return fmt.Sprintf("%s%s/%s/", HashicorpProductList, product, version)
}

func ProductBuildChecksumUrl(product string, version string) string {
	return fmt.Sprintf("%s%s/%s/%s_%s_SHA256SUMS", HashicorpProductList, product, version, product, version)
}

func ProductBuildUrl(product string, version string, os string, arch string) string {
	return fmt.Sprintf("%s%s/%s/%s_%s_%s_%s.zip", HashicorpProductList, product, version, product, version, os, arch)
}
