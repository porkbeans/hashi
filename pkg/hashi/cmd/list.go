package cmd

import (
	"fmt"
	"github.com/porkbeans/hashi/pkg/parseutils"
	"github.com/porkbeans/hashi/pkg/urlutils"
	"github.com/spf13/cobra"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
	"runtime"
)

func getNode(url string) (*html.Node, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		break
	case 403:
		return nil, fmt.Errorf("not found")
	default:
		return nil, fmt.Errorf("failed to get %s (status: %d)", url, res.StatusCode)
	}

	return html.Parse(res.Body)
}

func showLinkList(l parseutils.LinkList, cmd *cobra.Command) {
	for _, link := range l {
		cmd.Printf("%s\n", link.Name)
	}
}

func showProductVersionList(l parseutils.ProductVersionList, cmd *cobra.Command) {
	for _, link := range l {
		cmd.Printf("%s\n", link.Version)
	}
}

func showProductZipList(l parseutils.ProductZipList, cmd *cobra.Command) {
	for _, zipLink := range l {
		mark := ""
		if zipLink.Os == runtime.GOOS && zipLink.Arch == runtime.GOARCH {
			mark = "*"
		}

		cmd.Printf("%s %s %s\n", zipLink.Os, zipLink.Arch, mark)

	}
}

var listCmd = &cobra.Command{
	Use:   "list [name] [version]",
	Short: "List HashiCorp tools.",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var rawURL string
		switch len(args) {
		case 0:
			rawURL = urlutils.HashicorpProductList
		case 1:
			rawURL = urlutils.ProductVersionListURL(args[0])
		case 2:
			rawURL = urlutils.ProductZipListURL(args[0], args[1])
		}

		baseURL, err := url.Parse(rawURL)
		if err != nil {
			cmd.SetOutput(os.Stderr)
			cmd.Printf("Error: %s\n", err)
			return
		}

		root, err := getNode(baseURL.String())
		if err != nil {
			cmd.SetOutput(os.Stderr)
			cmd.Printf("Error: %s\n", err)
			return
		}

		linkList := parseutils.ParseLinkList(baseURL, root)

		switch len(args) {
		case 0:
			showLinkList(linkList, cmd)
		case 1:
			showProductVersionList(linkList.ProductVersionList(), cmd)
		case 2:
			showProductZipList(linkList.ProductZipList(), cmd)
		}
	},
}
