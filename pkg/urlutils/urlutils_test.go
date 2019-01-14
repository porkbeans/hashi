package urlutils

import (
	"github.com/porkbeans/hashi/pkg/parseutils"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	productListURL, err := url.Parse(HashicorpProductList)
	if err != nil {
		t.Error(err)
	}

	if productListURL == nil {
		t.Errorf("failed to parse URL")
	}
}

func TestProductVersionListUrl(t *testing.T) {
	testCases := parseutils.LinkList{
		{Name: "terraform", URL: "https://releases.hashicorp.com/terraform/"},
		{Name: "vault", URL: "https://releases.hashicorp.com/vault/"},
		{Name: "consul", URL: "https://releases.hashicorp.com/consul/"},
		{Name: "nomad", URL: "https://releases.hashicorp.com/nomad/"},
		{Name: "vagrant", URL: "https://releases.hashicorp.com/vagrant/"},
		{Name: "packer", URL: "https://releases.hashicorp.com/packer/"},
	}

	for _, expected := range testCases {
		actualURL := ProductVersionListURL(expected.Name)

		if actualURL == expected.URL {
			t.Logf("%s -> %s", expected.Name, actualURL)
		} else {
			t.Errorf("%s is not %s", actualURL, expected.URL)
		}
	}
}

func TestProductZipListURL(t *testing.T) {
	testCases := parseutils.ProductVersionList{
		{Name: "terraform", Version: "0.11.11", URL: "https://releases.hashicorp.com/terraform/0.11.11/"},
		{Name: "vault", Version: "1.0.1", URL: "https://releases.hashicorp.com/vault/1.0.1/"},
		{Name: "consul", Version: "1.4.0", URL: "https://releases.hashicorp.com/consul/1.4.0/"},
		{Name: "nomad", Version: "0.8.6", URL: "https://releases.hashicorp.com/nomad/0.8.6/"},
		{Name: "vagrant", Version: "2.2.2", URL: "https://releases.hashicorp.com/vagrant/2.2.2/"},
		{Name: "packer", Version: "1.3.3", URL: "https://releases.hashicorp.com/packer/1.3.3/"},
	}

	for _, expected := range testCases {
		actualURL := ProductZipListURL(expected.Name, expected.Version)

		if actualURL == expected.URL {
			t.Logf("%s %s -> %s", expected.Name, expected.Version, actualURL)
		} else {
			t.Errorf("%s is not %s", actualURL, expected.URL)
		}
	}
}

func TestProductZipChecksumURL(t *testing.T) {
	testCases := parseutils.ProductVersionList{
		{Name: "terraform", Version: "0.11.11", URL: "https://releases.hashicorp.com/terraform/0.11.11/terraform_0.11.11_SHA256SUMS"},
		{Name: "vault", Version: "1.0.1", URL: "https://releases.hashicorp.com/vault/1.0.1/vault_1.0.1_SHA256SUMS"},
		{Name: "consul", Version: "1.4.0", URL: "https://releases.hashicorp.com/consul/1.4.0/consul_1.4.0_SHA256SUMS"},
		{Name: "nomad", Version: "0.8.6", URL: "https://releases.hashicorp.com/nomad/0.8.6/nomad_0.8.6_SHA256SUMS"},
		{Name: "vagrant", Version: "2.2.2", URL: "https://releases.hashicorp.com/vagrant/2.2.2/vagrant_2.2.2_SHA256SUMS"},
		{Name: "packer", Version: "1.3.3", URL: "https://releases.hashicorp.com/packer/1.3.3/packer_1.3.3_SHA256SUMS"},
	}

	for _, expected := range testCases {
		actualURL := ProductZipChecksumURL(expected.Name, expected.Version)

		if actualURL == expected.URL {
			t.Logf("%s %s -> %s", expected.Name, expected.Version, actualURL)
		} else {
			t.Errorf("%s is not %s", actualURL, expected.URL)
		}
	}
}

func TestProductZipURL(t *testing.T) {
	testCases := parseutils.ProductZipList{
		{Name: "terraform", Version: "0.11.11", Os: "linux", Arch: "amd64", URL: "https://releases.hashicorp.com/terraform/0.11.11/terraform_0.11.11_linux_amd64.zip"},
		{Name: "vault", Version: "1.0.1", Os: "linux", Arch: "amd64", URL: "https://releases.hashicorp.com/vault/1.0.1/vault_1.0.1_linux_amd64.zip"},
		{Name: "consul", Version: "1.4.0", Os: "linux", Arch: "amd64", URL: "https://releases.hashicorp.com/consul/1.4.0/consul_1.4.0_linux_amd64.zip"},
		{Name: "nomad", Version: "0.8.6", Os: "linux", Arch: "amd64", URL: "https://releases.hashicorp.com/nomad/0.8.6/nomad_0.8.6_linux_amd64.zip"},
		{Name: "vagrant", Version: "2.2.2", Os: "linux", Arch: "amd64", URL: "https://releases.hashicorp.com/vagrant/2.2.2/vagrant_2.2.2_linux_amd64.zip"},
		{Name: "packer", Version: "1.3.3", Os: "linux", Arch: "amd64", URL: "https://releases.hashicorp.com/packer/1.3.3/packer_1.3.3_linux_amd64.zip"},
	}

	for _, expected := range testCases {
		actualURL := ProductZipURL(expected.Name, expected.Version, expected.Os, expected.Arch)

		if actualURL == expected.URL {
			t.Logf("%s %s -> %s", expected.Name, expected.Version, actualURL)
		} else {
			t.Errorf("%s is not %s", actualURL, expected.URL)
		}
	}
}
