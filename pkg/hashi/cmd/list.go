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

func showProductBuildList(l parseutils.ProductBuildList, cmd *cobra.Command) {
	for _, buildLink := range l {
		mark := ""
		if buildLink.Os == runtime.GOOS && buildLink.Arch == runtime.GOARCH {
			mark = "*"
		}

		cmd.Printf("%s %s %s\n", buildLink.Os, buildLink.Arch, mark)

	}
}

var listCmd = &cobra.Command{
	Use:   "list [name] [version]",
	Short: "List hashicorp tools.",
	Args:  cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var rawUrl string
		switch len(args) {
		case 0:
			rawUrl = urlutils.HashicorpProductList
		case 1:
			rawUrl = urlutils.ProductVersionListUrl(args[0])
		case 2:
			rawUrl = urlutils.ProductBuildListUrl(args[0], args[1])
		}

		baseUrl, err := url.Parse(rawUrl)
		if err != nil {
			cmd.SetOutput(os.Stderr)
			cmd.Printf("Error: %s\n", err)
			return
		}

		root, err := getNode(baseUrl.String())
		if err != nil {
			cmd.SetOutput(os.Stderr)
			cmd.Printf("Error: %s\n", err)
			return
		}

		linkList := parseutils.ParseLinkList(baseUrl, root)

		switch len(args) {
		case 0:
			showLinkList(linkList, cmd)
		case 1:
			showProductVersionList(linkList.ProductVersionList(), cmd)
		case 2:
			showProductBuildList(linkList.ProductBuildList(), cmd)
		}
	},
}
