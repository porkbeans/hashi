package parseutils

import (
	"encoding/hex"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/url"
	"os"
	"testing"
)

func getNodeFromFile(path string) (*html.Node, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return html.Parse(file)
}

func TestParseLinkList(t *testing.T) {
	node, err := getNodeFromFile("testdata/product_list.html")
	if err != nil {
		t.Error(err)
	}

	expectedList := LinkList{
		{Name: "consul", Url: urlutils.HashicorpProductList + "consul/"},
		{Name: "nomad", Url: urlutils.HashicorpProductList + "nomad/"},
		{Name: "packer", Url: urlutils.HashicorpProductList + "packer/"},
		{Name: "terraform", Url: urlutils.HashicorpProductList + "terraform/"},
		{Name: "vagrant", Url: urlutils.HashicorpProductList + "vagrant/"},
		{Name: "vault", Url: urlutils.HashicorpProductList + "vault/"},
	}

	baseUrl, _ := url.Parse(urlutils.HashicorpProductList)
	actualList := ParseLinkList(baseUrl, node)

	if len(expectedList) != len(actualList) {
		t.Errorf("failed to parse list")
	}

	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != actualList[i] {
			t.Errorf("%+v is not equal to %+v", expectedList[i], actualList[i])
		}
		t.Logf("%s -> %s\n", actualList[i].Name, actualList[i].Url)
	}
}

func TestLinkList_ProductVersionList(t *testing.T) {
	node, err := getNodeFromFile("testdata/product_version_list.html")
	if err != nil {
		t.Error(err)
	}

	expectedList := ProductVersionList{
		{Name: "consul", Version: "1.4.0", Url: urlutils.HashicorpProductList + "consul/1.4.0/"},
		{Name: "consul", Version: "1.4.0-rc1", Url: urlutils.HashicorpProductList + "consul/1.4.0-rc1/"},
		{Name: "consul", Version: "1.3.1", Url: urlutils.HashicorpProductList + "consul/1.3.1/"},
	}

	baseUrl, _ := url.Parse(urlutils.HashicorpProductList)
	linkList := ParseLinkList(baseUrl, node)
	actualList := linkList.ProductVersionList()

	if len(expectedList) != len(actualList) {
		t.Errorf("failed to parse list")
	}

	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != actualList[i] {
			t.Errorf("%+v is not equal to %+v", expectedList[i], actualList[i])
		}
		t.Logf("%s %s -> %s\n", actualList[i].Name, actualList[i].Version, actualList[i].Url)
	}
}

func TestLinkList_ProductBuildList(t *testing.T) {
	node, err := getNodeFromFile("testdata/product_build_list.html")
	if err != nil {
		t.Error(err)
	}

	baseUrl, _ := url.Parse(urlutils.HashicorpProductList)

	expectedList := ProductBuildList{
		{Name: "consul", Version: "1.4.0", Os: "darwin", Arch: "386", Url: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_darwin_386.zip"},
		{Name: "consul", Version: "1.4.0", Os: "darwin", Arch: "amd64", Url: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_darwin_amd64.zip"},
		{Name: "consul", Version: "1.4.0", Os: "freebsd", Arch: "386", Url: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_freebsd_386.zip"},
		{Name: "consul", Version: "1.4.0", Os: "freebsd", Arch: "amd64", Url: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_freebsd_amd64.zip"},
	}

	linkList := ParseLinkList(baseUrl, node)
	actualList := linkList.ProductBuildList()

	if len(expectedList) != len(actualList) {
		t.Errorf("failed to parse list")
	}

	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != actualList[i] {
			t.Errorf("%+v is not equal to %+v", expectedList[i], actualList[i])
		}
		t.Logf("%s %s %s %s -> %s\n", actualList[i].Name, actualList[i].Version, actualList[i].Os, actualList[i].Arch, actualList[i].Url)
	}
}

func decodeHex(hexString string) [32]byte {
	b := [32]byte{}
	h, _ := hex.DecodeString(hexString)
	copy(b[:], h[0:32])
	return b
}

func TestParseChecksumList(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/checksums")
	if err != nil {
		t.Error(err)
	}

	expectedList := ChecksumList{
		{"consul", "1.4.0", "darwin", "386", decodeHex("bf1e3f225c7af45d10efe1541a0a647cf534566f57a34d8929b8f4524cd20189")},
		{"consul", "1.4.0", "darwin", "amd64", decodeHex("8a7118cf29c697ddd072eaf40080d61aea16f606e50ef4e6784ca121c5fa1c1e")},
		{"consul", "1.4.0", "freebsd", "386", decodeHex("86b1a3fb550bbf4699e8a590f120bbe00e6274b9e39e3406faafcc0f58e5ba9f")},
		{"consul", "1.4.0", "freebsd", "amd64", decodeHex("dfc629df4bb697ffd1fcee2b2936270f29ee7e2dc80aeb4c9cd178f68c89683f")},
	}

	actualList := ParseChecksumList(string(content))

	if len(expectedList) != len(actualList) {
		t.Errorf("failed to parse list")
	}

	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != actualList[i] {
			t.Errorf("%+v is not equal to %+v", expectedList[i], actualList[i])
		}
		t.Logf("%s %s %s %s -> %s\n", actualList[i].Name, actualList[i].Version, actualList[i].Os, actualList[i].Arch, hex.EncodeToString(actualList[i].Checksum[:]))
	}
}
