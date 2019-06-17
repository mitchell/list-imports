package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/list-imports/exploration"
)

func Execute() {
	var includeVendor bool

	rootCmd := &cobra.Command{
		Use:   "list-imports [dir]",
		Short: "List the imports of a go project",
		Long: `List the imports of a go project or specified directory as JSON. Optionally shows 
imports of vendor folder, essentially showing transitive dependencies. 
Specifying a dir is optional, and will default to the working directory.`,
		Args: cobra.MaximumNArgs(1),
		Run:  makeRootRun(&includeVendor),
	}

	rootCmd.Flags().BoolVarP(&includeVendor, "include-vendor", "i", false, "include vendor dir in listing of imports")

	check(rootCmd.Execute())
}

func makeRootRun(includeVendor *bool) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		root := "."

		if len(args) == 1 {
			root = args[0]
		}

		imports, err := exploration.FindImports(root, *includeVendor)
		check(err)

		importsJSON, err := json.MarshalIndent(imports, "", "  ")
		check(err)

		fmt.Println(string(importsJSON))
	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
