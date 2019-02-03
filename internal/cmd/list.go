package cmd

import (
	"fmt"
	"io"
	"net/url"
	"runtime"

	"github.com/porkbeans/hashi/internal/httputils"
	"github.com/porkbeans/hashi/internal/ioutils"

	"github.com/porkbeans/hashi/pkg/parseutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

func parseURL(args []string) string {
	switch len(args) {
	case 0:
		return urlutils.HashicorpProductList
	case 1:
		return urlutils.ProductVersionListURL(args[0])
	case 2:
		return urlutils.ProductZipListURL(args[0], args[1])
	default:
		return ""
	}
}

func getList(client httputils.HTTPGetClient, rawURL string) (parseutils.LinkList, error) {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	resp, err := httputils.Get(client, rawURL)
	if err != nil {
		return nil, err
	}
	defer ioutils.Close(resp.Body)

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseutils.ParseLinkList(baseURL, root), err
}

func showLinkList(linkList parseutils.LinkList, writer io.Writer) {
	for _, link := range linkList {
		_, _ = fmt.Fprintf(writer, "%s\n", link.Name)
	}
}

func showProductVersionList(linkList parseutils.ProductVersionList, writer io.Writer) {
	for _, link := range linkList {
		_, _ = fmt.Fprintf(writer, "%s\n", link.Version)
	}
}

func showProductZipList(linkList parseutils.ProductZipList, writer io.Writer) {
	for _, zipLink := range linkList {
		mark := ""
		if zipLink.Os == runtime.GOOS && zipLink.Arch == runtime.GOARCH {
			mark = "*"
		}

		_, _ = fmt.Fprintf(writer, "%s %s %s\n", zipLink.Os, zipLink.Arch, mark)
	}
}

var listCmd = &cobra.Command{
	Use:   "list [name] [version]",
	Short: "List HashiCorp tools.",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		rawURL := parseURL(args)
		linkList, err := getList(nil, rawURL)
		if err != nil {
			return err
		}

		switch len(args) {
		case 0:
			showLinkList(linkList, cmd.OutOrStdout())
		case 1:
			showProductVersionList(linkList.ProductVersionList(), cmd.OutOrStdout())
		case 2:
			showProductZipList(linkList.ProductZipList(), cmd.OutOrStdout())
		}

		return nil
	},
}
