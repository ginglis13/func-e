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

package debug

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/tetratelabs/getenvoy/pkg/binary"
	"github.com/tetratelabs/getenvoy/pkg/binary/envoy"
	"github.com/tetratelabs/getenvoy/pkg/manifest"
)

const envoyReference = "standard:1.11.0"

func fetchEnvoy(t *testing.T) {
	key, _ := manifest.NewKey(envoyReference)
	r, _ := envoy.NewRuntime()
	if !r.AlreadyDownloaded(key) {
		location, err := manifest.Locate(key, manifest.DefaultURL)
		if err != nil {
			t.Fatalf("unable to retrieve manifest from %v: %v", manifest.DefaultURL, err)
		}
		if err := r.Fetch(key, location); err != nil {
			t.Fatalf("unable to retrieve binary from %v: %v", location, err)
		}
	}
}

func startWaitKillGetEnvoy(r binary.Runner, key *manifest.Key) {
	go r.Run(key, []string{"-c", filepath.Join("testdata", "bootstrap.yaml")})
	r.Wait(binary.StatusReady)
	r.SendSignal(syscall.SIGINT)
	r.Wait(binary.StatusTerminated)
}

// This test relies on a local Envoy binary, if not present it will fetch one from GetEnvoy
// This is more of an integration test than a unit test, but either way is necessary.
func Test_retrieveAdminAPIData(t *testing.T) {
	fetchEnvoy(t)
	t.Run("creates all non-empty files", func(t *testing.T) {
		key, _ := manifest.NewKey(envoyReference)
		r, _ := envoy.NewRuntime(EnableEnvoyAdminDataCollection)
		defer os.RemoveAll(r.DebugStore())
		startWaitKillGetEnvoy(r, key)

		for _, filename := range adminAPIPaths {
			path := filepath.Join(r.DebugStore(), filename)
			f, err := os.Stat(path)
			if err != nil {
				t.Errorf("error stating %v: %v", path, err)
			}
			if f.Size() < 1 {
				t.Errorf("file %v was empty", path)
			}
		}
	})

}
