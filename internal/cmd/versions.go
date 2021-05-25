// Copyright 2019 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"runtime"

	"github.com/urfave/cli/v2"

	"github.com/tetratelabs/getenvoy/internal/binary/envoy"
	"github.com/tetratelabs/getenvoy/internal/globals"
	"github.com/tetratelabs/getenvoy/internal/manifest"
)

// NewVersionsCmd returns command that lists available Envoy versions for the current platform.
func NewVersionsCmd(o *globals.GlobalOpts) *cli.Command {
	return &cli.Command{
		Name:  "versions",
		Usage: "List available Envoy versions",
		Action: func(c *cli.Context) error {
			m, err := manifest.GetManifest(o.ManifestURL)
			if err != nil {
				return err
			}
			envoy.PrintVersions(m, runtime.GOOS, c.App.Writer)
			return nil
		},
	}
}