package cmd

import (
	"fmt"
	"io"
	"net/url"
	"runtime"

	"github.com/porkbeans/hashi/internal/ioutils"
	"github.com/porkbeans/hashi/pkg/parseutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
)

func getList(rawURL string) (parseutils.LinkList, error) {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	resp, err := ioutils.Get(nil, rawURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return parseutils.ParseLinkList(baseURL, root), err
}

func showLinkList(linkList parseutils.LinkList, writer io.Writer) {
	for _, link := range linkList {
		fmt.Fprintf(writer, "%s\n", link.Name)
	}
}

func showProductVersionList(linkList parseutils.ProductVersionList, writer io.Writer) {
	for _, link := range linkList {
		fmt.Fprintf(writer, "%s\n", link.Version)
	}
}

func showProductZipList(linkList parseutils.ProductZipList, writer io.Writer) {
	for _, zipLink := range linkList {
		mark := ""
		if zipLink.Os == runtime.GOOS && zipLink.Arch == runtime.GOARCH {
			mark = "*"
		}

		fmt.Fprintf(writer, "%s %s %s\n", zipLink.Os, zipLink.Arch, mark)
	}
}

var listCmd = &cobra.Command{
	Use:   "list [name] [version]",
	Short: "List HashiCorp tools.",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch len(args) {
		case 0:
			rawURL := urlutils.HashicorpProductList
			linkList, err := getList(rawURL)
			if err != nil {
				return err
			}
			showLinkList(linkList, cmd.OutOrStdout())
		case 1:
			rawURL := urlutils.ProductVersionListURL(args[0])
			linkList, err := getList(rawURL)
			if err != nil {
				return err
			}
			showProductVersionList(linkList.ProductVersionList(), cmd.OutOrStdout())
		case 2:
			rawURL := urlutils.ProductZipListURL(args[0], args[1])
			linkList, err := getList(rawURL)
			if err != nil {
				return err
			}
			showProductZipList(linkList.ProductZipList(), cmd.OutOrStdout())
		}

		return nil
	},
}
