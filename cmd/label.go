package cmd

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/garaemon/paperpile/internal/api"
	"github.com/garaemon/paperpile/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	labelCmd.AddCommand(labelListCmd)
	rootCmd.AddCommand(labelCmd)
}

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage labels on library items",
}

var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available labels",
	Args:  cobra.NoArgs,
	RunE:  runLabelList,
}

func runLabelList(cmd *cobra.Command, args []string) error {
	client := api.NewClient(config.GetSession())
	return execLabelList(client, os.Stdout)
}

func execLabelList(fetcher LabelFetcher, out io.Writer) error {
	labels, err := fetcher.FetchLabels()
	if err != nil {
		return fmt.Errorf("failed to fetch labels: %w", err)
	}

	if len(labels) == 0 {
		fmt.Fprintln(out, "(no labels)")
		return nil
	}

	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tCOUNT")
	for _, label := range labels {
		fmt.Fprintf(w, "%s\t%s\t%d\n", label.ID, label.Name, label.Count)
	}
	w.Flush()
	return nil
}
