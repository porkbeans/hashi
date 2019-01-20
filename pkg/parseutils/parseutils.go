package parseutils

import (
	"encoding/hex"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// LinkEntry represents a link with label.
type LinkEntry struct {
	Name string
	URL  string
}

// LinkList represents a list of LinkEntry.
type LinkList []LinkEntry

// ParseLinkList parses html to a list of links.
func ParseLinkList(baseURL *url.URL, node *html.Node) LinkList {
	linkList := LinkList{}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Parent.Data == "li" && c.Data == "a" {
			text := c.FirstChild.Data

			for _, attr := range c.Attr {
				if attr.Key == "href" {
					relativeURL, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}

					linkList = append(
						linkList, LinkEntry{
							Name: text,
							URL:  baseURL.ResolveReference(relativeURL).String(),
						})
				}
			}
		} else if c.FirstChild != nil {
			linkList = append(linkList, ParseLinkList(baseURL, c)...)
		}
	}

	return linkList
}

// ProductVersionEntry represents a link of specific version of a HashiCorp product.
type ProductVersionEntry struct {
	Name    string
	Version string
	URL     string
}

// ProductVersionList represents a list of ProductVersionEntry.
type ProductVersionList []ProductVersionEntry

func mapSubmatchNames(pattern *regexp.Regexp, str string) (map[string]string, bool) {
	matches := map[string]string{}
	if !pattern.MatchString(str) {
		return matches, false
	}

	indices := map[string]int{}
	for i, name := range pattern.SubexpNames() {
		if len(name) > 0 {
			indices[name] = i
		}
	}

	rawMatches := pattern.FindStringSubmatch(str)
	for name, i := range indices {
		matches[name] = rawMatches[i]
	}

	return matches, true
}

// ProductVersionList parses LinkList to ProductVersionList
func (l LinkList) ProductVersionList() ProductVersionList {
	pattern := regexp.MustCompile(`^/(?P<product>.+)/(?P<version>.+)/$`)
	versionList := ProductVersionList{}

	for _, link := range l {
		versionURL, err := url.Parse(link.URL)
		if err != nil {
			continue
		}

		if matches, matched := mapSubmatchNames(pattern, versionURL.Path); matched {
			versionList = append(versionList, ProductVersionEntry{
				Name:    matches["product"],
				Version: matches["version"],
				URL:     link.URL,
			})
		}
	}

	return versionList
}

// ProductZipEntry represents a link to a HashiCorp product's zip.
type ProductZipEntry struct {
	Name    string
	Version string
	Os      string
	Arch    string
	URL     string
}

// ProductZipList represents a list of ProductZipEntry.
type ProductZipList []ProductZipEntry

// ProductZipList parses a LinkList to a ProductZipList.
func (l LinkList) ProductZipList() ProductZipList {
	pattern := regexp.MustCompile(`^/(?P<product>.+)/(?P<version>.+)/.+_(?P<os>.+)_(?P<arch>.+)\.zip$`)
	zipList := ProductZipList{}

	for _, link := range l {
		zipURL, err := url.Parse(link.URL)
		if err != nil {
			continue
		}

		if matches, matched := mapSubmatchNames(pattern, zipURL.Path); matched {
			zipList = append(zipList, ProductZipEntry{
				Name:    matches["product"],
				Version: matches["version"],
				Os:      matches["os"],
				Arch:    matches["arch"],
				URL:     link.URL,
			})
		}
	}

	return zipList
}

// ChecksumEntry represents a checksum for a HashiCorp product's zip file.
type ChecksumEntry struct {
	Name     string
	Version  string
	Os       string
	Arch     string
	Checksum [32]byte
}

// ChecksumList represents a list of ChecksumEntry.
type ChecksumList []ChecksumEntry

// ParseChecksumList parses a checksum file to a ChecksumList.
func ParseChecksumList(rawChecksumList string) ChecksumList {
	entryPattern := regexp.MustCompile(`^(?P<checksum>[0-9a-f]{64})\s+(?P<product>.+)_(?P<version>.+)_(?P<os>.+)_(?P<arch>.+)\.zip$`)
	checksumList := ChecksumList{}

	for _, rawEntry := range strings.Split(rawChecksumList, "\n") {
		matches, matched := mapSubmatchNames(entryPattern, rawEntry)
		if !matched {
			continue
		}

		checksum := [32]byte{}
		checksumBuffer, _ := hex.DecodeString(matches["checksum"])
		copy(checksum[:], checksumBuffer[0:32])

		checksumList = append(checksumList, ChecksumEntry{
			Name:     matches["product"],
			Version:  matches["version"],
			Os:       matches["os"],
			Arch:     matches["arch"],
			Checksum: checksum,
		})
	}

	return checksumList
}
