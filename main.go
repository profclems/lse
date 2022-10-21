package main

import (
	"github.com/pee2pee/lse/ls"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	lsf := ls.LS{
		StdOut: os.Stdout,
		Stderr: os.Stderr,
	}

	cmd := &cobra.Command{
		Use:   "lse [dir]",
		Short: "A cross platform drop-in replacement for ls",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			lsf.Dir = "."
			if len(args) == 1 {
				lsf.Dir = args[0]
			}

			return lsf.ListDir()
		},
	}

	cmd.Flags().BoolVarP(&lsf.A, "all", "a", false, "show all files including hidden files")
	cmd.Flags().BoolVarP(&lsf.L, "tabular", "l", false, "show detailed directory structure in tabular form")
	cmd.Flags().BoolVarP(&lsf.R, "recursive", "R", false, "show all subdirectories encountered")
	cmd.Flags().BoolVarP(&lsf.D, "directory", "d", false, "show directory structure")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
