package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	readPlainCmd = &cobra.Command{
		Use:   "readplain [text] ...",
		Short: "Let Anna read plain text input",
		Long:  "Let Anna read plain text input",
		Run:   readPlainRun,
	}
)

func readPlainRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		os.Exit(1)
	}

	ctx := context.Background()

	ID, err := textInterface.ReadPlainWithPlain(ctx, strings.Join(args, " "))
	if len(args) == 0 {
		fmt.Printf("%#v\n", maskAny(err))
		os.Exit(1)
	}

	data, err := textInterface.ReadPlainWithID(ctx, ID)
	if len(args) == 0 {
		fmt.Printf("%#v\n", maskAny(err))
		os.Exit(1)
	}

	fmt.Printf("%s\n", data)
}