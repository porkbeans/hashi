package urlutils

import (
	"fmt"
)

const (
	// HashicorpProductList returns a URL containing HashiCorp product list.
	HashicorpProductList = "https://releases.hashicorp.com/"
)

/*
ProductVersionListURL returns a URL containing version list of specified HashiCorp product.

URL Examples

  https://releases.hashicorp.com/consul/
  https://releases.hashicorp.com/nomad/
  https://releases.hashicorp.com/packer/
  https://releases.hashicorp.com/terraform/
  https://releases.hashicorp.com/vagrant/
  https://releases.hashicorp.com/vault/
*/
func ProductVersionListURL(product string) string {
	return fmt.Sprintf("%s%s/", HashicorpProductList, product)
}

/*
ProductBuildListURL returns a URL containing build list of specified HashiCorp product's version.

URL Examples

  https://releases.hashicorp.com/consul/1.4.0/
  https://releases.hashicorp.com/nomad/0.8.7/
  https://releases.hashicorp.com/packer/1.3.3/
  https://releases.hashicorp.com/terraform/0.11.11/
  https://releases.hashicorp.com/vagrant/2.2.3/
  https://releases.hashicorp.com/vault/1.0.1/
*/
func ProductBuildListURL(product string, version string) string {
	return fmt.Sprintf("%s%s/%s/", HashicorpProductList, product, version)
}

/*
ProductBuildChecksumURL returns a checksum URL of specified HashiCorp product's version.

URL Examples

  https://releases.hashicorp.com/consul/1.4.0/consul_1.4.0_SHA256SUMS
  https://releases.hashicorp.com/nomad/0.8.7/nomad_0.8.7_SHA256SUMS
  https://releases.hashicorp.com/packer/1.3.3/packer_1.3.3_SHA256SUMS
  https://releases.hashicorp.com/terraform/0.11.11/terraform_0.11.11_SHA256SUMS
  https://releases.hashicorp.com/vagrant/2.2.3/vagrant_2.2.3_SHA256SUMS
  https://releases.hashicorp.com/vault/1.0.1/vault_1.0.1_SHA256SUMS
*/
func ProductBuildChecksumURL(product string, version string) string {
	return fmt.Sprintf("%s%s/%s/%s_%s_SHA256SUMS", HashicorpProductList, product, version, product, version)
}

/*
ProductBuildURL returns a build file URL of specified HashiCorp product's version.

URL Examples

  https://releases.hashicorp.com/consul/1.4.0/consul_1.4.0_linux_amd64.zip
  https://releases.hashicorp.com/nomad/0.8.7/nomad_0.8.7_linux_amd64.zip
  https://releases.hashicorp.com/packer/1.3.3/packer_1.3.3_linux_amd64.zip
  https://releases.hashicorp.com/terraform/0.11.11/terraform_0.11.11_linux_amd64.zip
  https://releases.hashicorp.com/vagrant/2.2.3/vagrant_2.2.3_linux_amd64.zip
  https://releases.hashicorp.com/vault/1.0.1/vault_1.0.1_linux_amd64.zip
*/
func ProductBuildURL(product string, version string, os string, arch string) string {
	return fmt.Sprintf("%s%s/%s/%s_%s_%s_%s.zip", HashicorpProductList, product, version, product, version, os, arch)
}
