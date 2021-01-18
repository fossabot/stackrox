package initbundles

import (
	"context"
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	v1 "github.com/stackrox/rox/generated/api/v1"
	pkgCommon "github.com/stackrox/rox/pkg/roxctl/common"
	"github.com/stackrox/rox/pkg/utils"
	"github.com/stackrox/rox/roxctl/common"
	"github.com/stackrox/rox/roxctl/common/util"
)

func listInitBundles() error {
	ctx, cancel := context.WithTimeout(pkgCommon.Context(), contextTimeout)
	defer cancel()

	conn, err := common.GetGRPCConnection()
	if err != nil {
		return err
	}
	defer utils.IgnoreError(conn.Close)
	svc := v1.NewClusterInitServiceClient(conn)

	tabWriter := tabwriter.NewWriter(os.Stdout, 4, 8, 2, '\t', 0)

	rsp, err := svc.GetInitBundles(ctx, &v1.Empty{})
	if err != nil {
		return errors.Wrap(err, "getting all init bundles")
	}

	bundles := rsp.GetItems()
	sort.Slice(bundles, func(i, j int) bool { return bundles[i].GetName() < bundles[j].GetName() })

	fmt.Fprintln(tabWriter, " \tName\tCreated at\tExpires at\tID")
	fmt.Fprintln(tabWriter, " \t====\t==========\t==========\t==")

	for idx, meta := range bundles {
		name := meta.GetName()
		if name == "" {
			name = "(empty)"
		}
		fmt.Fprintf(tabWriter, "%4d\t%s\t%s\t%v\t%v\n",
			idx,
			name,
			meta.GetCreatedAt(),
			meta.GetExpiresAt(),
			meta.GetId(),
		)
	}
	return errors.Wrap(tabWriter.Flush(), "flushing tabular output")
}

// listCommand implements the command for listing init bundles.
func listCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "list",
		RunE: util.RunENoArgs(func(c *cobra.Command) error {
			return listInitBundles()
		}),
	}
	return c
}
