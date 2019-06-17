package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mitchell/list-imports/explore"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use: "list-imports",
		Run: rootRun,
	}

	check(rootCmd.Execute())
}

func rootRun(cmd *cobra.Command, args []string) {
	imports, err := explore.FindImports(".")
	check(err)

	importjson, err := json.MarshalIndent(imports, "", "  ")
	check(err)

	fmt.Println(string(importjson))
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
