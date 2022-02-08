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

package dockerconfig

import (
	"errors"
	"testing"

	"github.com/docker/cli/cli/config/types"
	"github.com/stretchr/testify/assert"
)

type mockConfigGetter struct {
	res map[string]types.AuthConfig
	err error
}

func (m mockConfigGetter) GetAllCredentials() (map[string]types.AuthConfig, error) {
	return m.res, m.err
}

func TestLoadCredentials(t *testing.T) {
	type args struct {
		cfg credentialGetter
	}
	tests := []struct {
		name string
		args args
		want File
		err  string
	}{
		{
			name: "fails to obtain all credentials",
			args: args{cfg: mockConfigGetter{err: errors.New("some error")}},
			err:  "some error",
		},
		{
			name: "obtains a bunch of credentials",
			args: args{cfg: mockConfigGetter{
				res: map[string]types.AuthConfig{
					"docker.elastic.co": {
						Username: "my-user",
						Password: "my-password",
					},
					"different-docker.elastic.co": {
						Auth: "bXktdXNlcjpteS1wYXNzCg==",
					},
					"ignore-entry": {},
				},
			}},
			want: File{Auths: map[string]Credential{
				"different-docker.elastic.co": {
					Auth: "bXktdXNlcjpteS1wYXNzCg==",
				},
				"docker.elastic.co": {
					Auth: "bXktdXNlcjpteS1wYXNzd29yZA==",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadCredentials(tt.args.cfg)
			if err != nil {
				assert.EqualError(t, err, tt.err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
