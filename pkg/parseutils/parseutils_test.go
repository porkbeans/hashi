package parseutils

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"
	"net/url"
	"testing"
	"text/template"

	"github.com/porkbeans/hashi/internal/testutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"golang.org/x/net/html"
)

func TestParseLinkList1(t *testing.T) {
	expectedList := LinkList{
		{Name: "consul", URL: urlutils.HashicorpProductList + "consul/"},
		{Name: "nomad", URL: urlutils.HashicorpProductList + "nomad/"},
		{Name: "packer", URL: urlutils.HashicorpProductList + "packer/"},
		{Name: "terraform", URL: urlutils.HashicorpProductList + "terraform/"},
		{Name: "vagrant", URL: urlutils.HashicorpProductList + "vagrant/"},
		{Name: "vault", URL: urlutils.HashicorpProductList + "vault/"},
	}

	tmpl := template.Must(template.ParseFiles("testdata/product_list.html.tmpl"))
	buf := &bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(buf, "product_list.html.tmpl", expectedList); err != nil {
		t.Error(err)
		return
	}

	node, err := html.Parse(buf)
	if err != nil {
		t.Error(err)
		return
	}

	baseURL, _ := url.Parse(urlutils.HashicorpProductList)
	actualList := ParseLinkList(baseURL, node)

	if len(expectedList) != len(actualList) {
		t.Errorf("failed to parse list")
		return
	}

	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != actualList[i] {
			t.Errorf("%+v is not equal to %+v", expectedList[i], actualList[i])
		}
		t.Logf("%s -> %s\n", actualList[i].Name, actualList[i].URL)
	}
}

func TestParseLinkList2(t *testing.T) {
	linkList := LinkList{
		{Name: "consul", URL: testutils.GenerateInvalidURL()},
		{Name: "nomad", URL: testutils.GenerateInvalidURL()},
		{Name: "packer", URL: testutils.GenerateInvalidURL()},
		{Name: "terraform", URL: testutils.GenerateInvalidURL()},
		{Name: "vagrant", URL: testutils.GenerateInvalidURL()},
		{Name: "vault", URL: testutils.GenerateInvalidURL()},
	}

	tmpl := template.Must(template.ParseFiles("testdata/product_list.html.tmpl"))
	buf := &bytes.Buffer{}

	if err := tmpl.ExecuteTemplate(buf, "product_list.html.tmpl", linkList); err != nil {
		t.Error(err)
		return
	}

	node, err := html.Parse(buf)
	if err != nil {
		t.Error(err)
		return
	}

	baseURL, _ := url.Parse(urlutils.HashicorpProductList)
	actualList := ParseLinkList(baseURL, node)

	if len(actualList) != 0 {
		t.Errorf("failed to parse list")
	}
}

func TestLinkList_ProductVersionList1(t *testing.T) {
	linkList := LinkList{
		{Name: "consul_1.4.0", URL: urlutils.HashicorpProductList + "consul/1.4.0/"},
		{Name: "consul_1.4.0-rc1", URL: urlutils.HashicorpProductList + "consul/1.4.0-rc1/"},
		{Name: "consul_1.3.1", URL: urlutils.HashicorpProductList + "consul/1.3.1/"},
	}
	actualList := linkList.ProductVersionList()
	expectedList := ProductVersionList{
		{Name: "consul", Version: "1.4.0", URL: urlutils.HashicorpProductList + "consul/1.4.0/"},
		{Name: "consul", Version: "1.4.0-rc1", URL: urlutils.HashicorpProductList + "consul/1.4.0-rc1/"},
		{Name: "consul", Version: "1.3.1", URL: urlutils.HashicorpProductList + "consul/1.3.1/"},
	}

	if len(expectedList) != len(actualList) {
		t.Errorf("failed to parse list")
		return
	}

	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != actualList[i] {
			t.Errorf("%+v is not equal to %+v", expectedList[i], actualList[i])
		}
		t.Logf("%s %s -> %s\n", actualList[i].Name, actualList[i].Version, actualList[i].URL)
	}
}

func TestLinkList_ProductVersionList2(t *testing.T) {
	linkList := LinkList{
		{Name: "consul_1.4.0", URL: testutils.GenerateInvalidURL()},
		{Name: "consul_1.4.0-rc1", URL: testutils.GenerateInvalidURL()},
		{Name: "consul_1.3.1", URL: testutils.GenerateInvalidURL()},
	}
	actualList := linkList.ProductVersionList()

	if len(actualList) != 0 {
		t.Errorf("failed to parse list")
	}
}

func TestLinkList_ProductZipList(t *testing.T) {
	linkList := LinkList{
		{Name: "consul_1.4.0_darwin_386.zip", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_darwin_386.zip"},
		{Name: "consul_1.4.0_darwin_amd64.zip", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_darwin_amd64.zip"},
		{Name: "consul_1.4.0_freebsd_386.zip", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_freebsd_386.zip"},
		{Name: "consul_1.4.0_freebsd_amd64.zip", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_freebsd_amd64.zip"},
	}
	actualList := linkList.ProductZipList()
	expectedList := ProductZipList{
		{Name: "consul", Version: "1.4.0", Os: "darwin", Arch: "386", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_darwin_386.zip"},
		{Name: "consul", Version: "1.4.0", Os: "darwin", Arch: "amd64", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_darwin_amd64.zip"},
		{Name: "consul", Version: "1.4.0", Os: "freebsd", Arch: "386", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_freebsd_386.zip"},
		{Name: "consul", Version: "1.4.0", Os: "freebsd", Arch: "amd64", URL: urlutils.HashicorpProductList + "consul/1.4.0/consul_1.4.0_freebsd_amd64.zip"},
	}

	if len(expectedList) != len(actualList) {
		t.Errorf("failed to parse list")
		return
	}

	for i := 0; i < len(expectedList); i++ {
		if expectedList[i] != actualList[i] {
			t.Errorf("%+v is not equal to %+v", expectedList[i], actualList[i])
		}
		t.Logf("%s %s %s %s -> %s\n", actualList[i].Name, actualList[i].Version, actualList[i].Os, actualList[i].Arch, actualList[i].URL)
	}
}

func TestLinkList_ProductZipList2(t *testing.T) {
	linkList := LinkList{
		{Name: "consul_1.4.0_darwin_386.zip", URL: testutils.GenerateInvalidURL()},
		{Name: "consul_1.4.0_darwin_amd64.zip", URL: testutils.GenerateInvalidURL()},
		{Name: "consul_1.4.0_freebsd_386.zip", URL: testutils.GenerateInvalidURL()},
		{Name: "consul_1.4.0_freebsd_amd64.zip", URL: testutils.GenerateInvalidURL()},
	}
	actualList := linkList.ProductZipList()

	if len(actualList) != 0 {
		t.Errorf("failed to parse list")
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
		return
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
