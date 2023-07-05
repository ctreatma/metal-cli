package os

import (
	"testing"

	root "github.com/equinix/metal-cli/internal/cli"
	outputPkg "github.com/equinix/metal-cli/internal/outputs"
	"github.com/spf13/cobra"
)

func TestCli_RegisterCommands(t *testing.T) {
	rootClient := root.NewClient("", "https://api.equinix.com/metal/v1/", "test")
	rootCmd := rootClient.NewCommand()
	rootCmd.DisableSuggestions = false
	type fields struct {
		MainCmd  *cobra.Command
		Outputer outputPkg.Outputer
	}
	type args struct {
		client *root.Client
	}
	outputter := outputPkg.Outputer(&outputPkg.Standard{})
	tests := []struct {
		name    string
		fields  fields
		args    args
		cmdFunc func(*testing.T, *cobra.Command)
	}{
		{
			name: "test",
			fields: fields{
				MainCmd:  &cobra.Command{},
				Outputer: outputPkg.Outputer(outputter),
			},
			args: args{
				client: &root.Client{},
			},
			cmdFunc: func(t *testing.T, c *cobra.Command) {
				t.Helper()
				root := c.Root()
				if err := root.Execute(); err != nil {
					t.Logf("%+v", root.Args)
					t.Error("expected an error")
				}
			},
		},
		{
			name: "os",
			fields: fields{
				MainCmd:  NewClient(rootClient, outputter).NewCommand(),
				Outputer: outputPkg.Outputer(outputter),
			},
			args: args{
				client: rootClient,
			},
			cmdFunc: func(t *testing.T, c *cobra.Command) {
				t.Helper()
				root := c.Root()
				if err := root.Execute(); err != nil {
					t.Logf("%+v", root.Args)
					t.Error("expected an error")
				} else {
					t.Error("no error happened...but what did happen?")
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.cmdFunc(t, tt.fields.MainCmd)
		})
	}
}
