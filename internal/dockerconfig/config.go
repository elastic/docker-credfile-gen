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
	"encoding/base64"
	"fmt"

	"github.com/docker/cli/cli/config/types"
)

// File is a partial representation of the full docker configuration file,
// containing a map => Credential which contains credentials for registry
// authentication.
type File struct {
	Auths map[string]Credential `json:"auths,omitempty"`
}

// Credential represents credential entry.
type Credential struct {
	Auth string `json:"auth,omitempty"`
}

// credentialGetter abstracts the config file through this interface to facilitate
// testing.
type credentialGetter interface {
	GetAllCredentials() (map[string]types.AuthConfig, error)
}

// LoadCredentials takes in a docker configuration file and obtains all the
// credentials for it and returns a valid structure to be saved as a docker file
// to be sent to ce-aws environments.
func LoadCredentials(cfg credentialGetter) (File, error) {
	creds, err := cfg.GetAllCredentials()
	if err != nil {
		return File{}, err
	}

	credentials := File{Auths: make(map[string]Credential)}
	for k, cred := range creds {
		if cred.Auth != "" {
			credentials.Auths[k] = Credential{Auth: cred.Auth}
			continue
		}

		if cred.Username == "" && cred.Password == "" {
			continue
		}

		credentials.Auths[k] = Credential{
			Auth: base64.StdEncoding.EncodeToString([]byte(
				fmt.Sprintf("%s:%s", cred.Username, cred.Password),
			)),
		}
	}

	return credentials, nil
}
