// Copyright 2023 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package context

import (
	"sort"

	"github.com/redpanda-data/redpanda/src/go/rpk/pkg/config"
	"github.com/redpanda-data/redpanda/src/go/rpk/pkg/out"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func newListCommand(fs afero.Fs, p *config.Params) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List rpk contexts",
		Args:  cobra.ExactArgs(0),
		Run: func(*cobra.Command, []string) {
			cfg, err := p.Load(fs)
			out.MaybeDie(err, "unable to load config: %v", err)

			tw := out.NewTable("name", "description")
			defer tw.Flush()

			y, ok := cfg.ActualRpkYaml()
			if !ok {
				return
			}

			sort.Slice(y.Contexts, func(i, j int) bool {
				return y.Contexts[i].Name < y.Contexts[j].Name
			})

			for _, cx := range y.Contexts {
				name := cx.Name
				if name == y.CurrentContext {
					name += "*"
				}
				tw.Print(name, cx.Description)
			}
		},
	}
}
