// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/docker-credfile-gen/internal/dockerconfig"
)

func Test_writeFile(t *testing.T) {
	dockerConfigFile, err := ioutil.TempFile("", "docker-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		dockerConfigFile.Close()
		os.Remove(dockerConfigFile.Name())
	}()

	invalidPath, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	fileContents := `{
  "auths": {
    "different-docker.elastic.co": {
      "auth": "bXktdXNlcjpteS1wYXNzCg=="
    },
    "docker.elastic.co": {
      "auth": "bXktdXNlcjpteS1wYXNzd29yZA=="
    }
  }
}
`

	type args struct {
		name     string
		contents dockerconfig.File
	}
	tests := []struct {
		name string
		args args
		err  string
		want string
	}{
		{
			name: "writes the file contents to a file",
			args: args{
				name: dockerConfigFile.Name(),
				contents: dockerconfig.File{
					Auths: map[string]dockerconfig.Credential{
						"different-docker.elastic.co": {
							Auth: "bXktdXNlcjpteS1wYXNzCg==",
						},
						"docker.elastic.co": {
							Auth: "bXktdXNlcjpteS1wYXNzd29yZA==",
						},
					},
				},
			},
			want: fileContents,
		},
		{
			name: "fails to create the file contents",
			args: args{
				name:     invalidPath,
				contents: dockerconfig.File{},
			},
			err: fmt.Sprintf("open %s: is a directory", invalidPath),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeFile(tt.args.name, tt.args.contents); err != nil {
				assert.EqualError(t, err, tt.err)
				return
			}

			contents, err := ioutil.ReadFile(tt.args.name)
			if err != nil {
				t.Error(err)
				return
			}
			assert.Equal(t, tt.want, string(contents))
		})
	}
}
