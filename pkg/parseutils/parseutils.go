package parseutils

import (
	"encoding/hex"
	"golang.org/x/net/html"
	"net/url"
	"regexp"
	"strings"
)

type LinkEntry struct {
	Name string
	Url  string
}

type LinkList []LinkEntry

func ParseLinkList(baseUrl *url.URL, node *html.Node) LinkList {
	links := LinkList{}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode && c.Parent.Data == "li" && c.Data == "a" {
			text := c.FirstChild.Data

			for _, attr := range c.Attr {
				if attr.Key == "href" {
					relativeUrl, err := url.Parse(attr.Val)
					if err != nil {
						continue
					}

					links = append(
						links, LinkEntry{
							Name: text,
							Url:  baseUrl.ResolveReference(relativeUrl).String(),
						})
				}
			}
		} else if c.FirstChild != nil {
			links = append(links, ParseLinkList(baseUrl, c)...)
		}
	}

	return links
}

type ProductVersionEntry struct {
	Name    string
	Version string
	Url     string
}

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

func (l LinkList) ProductVersionList() ProductVersionList {
	pattern := regexp.MustCompile(`^/(?P<product>.+)/(?P<version>.+)/$`)
	versionLinks := ProductVersionList{}

	for _, link := range l {
		versionUrl, err := url.Parse(link.Url)
		if err != nil {
			continue
		}

		if matches, matched := mapSubmatchNames(pattern, versionUrl.Path); matched {
			versionLinks = append(versionLinks, ProductVersionEntry{
				Name:    matches["product"],
				Version: matches["version"],
				Url:     link.Url,
			})
		}
	}

	return versionLinks
}

type ProductBuildEntry struct {
	Name    string
	Version string
	Os      string
	Arch    string
	Url     string
}

type ProductBuildList []ProductBuildEntry

func (l LinkList) ProductBuildList() ProductBuildList {
	pattern := regexp.MustCompile(`^/(?P<product>.+)/(?P<version>.+)/.+_(?P<os>.+)_(?P<arch>.+)\.zip$`)
	buildLinks := ProductBuildList{}

	for _, link := range l {
		buildUrl, err := url.Parse(link.Url)
		if err != nil {
			continue
		}

		if matches, matched := mapSubmatchNames(pattern, buildUrl.Path); matched {
			buildLinks = append(buildLinks, ProductBuildEntry{
				Name:    matches["product"],
				Version: matches["version"],
				Os:      matches["os"],
				Arch:    matches["arch"],
				Url:     link.Url,
			})
		}
	}

	return buildLinks
}

type ChecksumEntry struct {
	Name     string
	Version  string
	Os       string
	Arch     string
	Checksum [32]byte
}

type ChecksumList []ChecksumEntry

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
