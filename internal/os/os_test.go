package cmd

import (
	"testing"

	root "github.com/equinix/metal-cli/internal/cli"
	"github.com/equinix/metal-cli/internal/os"
	outputPkg "github.com/equinix/metal-cli/internal/outputs"
	"github.com/spf13/cobra"
)

func TestCli_RegisterCommands(t *testing.T) {
	rootClient := root.NewClient(consumerToken, apiURL, Version)
	rootCmd := rootClient.NewCommand()
	rootCmd.DisableSuggestions = false
	type fields struct {
		MainCmd  *cobra.Command
		Outputer outputPkg.Outputer
		cli *Cli
	}
	type args struct {
		client *root.Client
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		cmdFunc func(*testing.T, *cobra.Command)
	}{
		{
			name: "test",
			fields: fields{
				MainCmd:  &cobra.Command{},
				Outputer: outputPkg.Outputer(&outputPkg.Standard{}),
			},
			args: args{
				client: &root.Client{},
			},
			cmdFunc: func(t *testing.T, c *cobra.Command) {
				t.Helper()
				root := c.Root()
				root.SetArgs([]string{subCommand})
				if err := root.Execute(); err != nil {
					t.Logf("%+v", root.Args)
					t.Error("expected an error")
				}
			},
		},
		{
			name: "os",
			fields: fields{
				MainCmd: os.NewClient(rootClient, Outputer).NewCommand(),
				Outputer: outputPkg.Outputer(&outputPkg.Standard{}),
			},
			args: args{
				client: rootClient,
			},
			cmdFunc: func(t *testing.T, c *cobra.Command) {
				t.Helper()
				root := c.Root()
				root.SetArgs([]string{subCommand})
				if err := root.Execute(); err != nil {
					t.Logf("%+v", root.Args)
					t.Error("expected an error")
				}
			},
		},
	},
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &Cli{
				MainCmd:  tt.fields.MainCmd,
				Outputer: tt.fields.Outputer,
			}
			cli.RegisterCommands(tt.args.client)
		})
	}
}
