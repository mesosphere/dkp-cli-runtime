package plugin

import (
	"encoding/json"
	"io"

	"github.com/spf13/cobra"
)

const DiscoveryCommandName = "_plugin_commands"

// Spec describes a CLI plugin in a serializable format.
type Spec struct {
	Commands CommandSpec `json:"commands"`
}

// NewDiscoveryCommand creates a hidden command returning a description of the passed command hierarchy as JSON.
func NewDiscoveryCommand(output io.Writer, rootCmd *cobra.Command) *cobra.Command {
	discoveryCmd := &cobra.Command{
		Use:    DiscoveryCommandName,
		Args:   cobra.NoArgs,
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			spec := Spec{
				Commands: SpecFromCommand(rootCmd),
			}
			encoder := json.NewEncoder(output)
			return encoder.Encode(spec)
		},
	}
	return discoveryCmd
}
