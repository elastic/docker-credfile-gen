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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/cli/cli/config"

	"github.com/elastic/docker-credfile-gen/internal/dockerconfig"
)

func main() {
	var credFile string
	var raw bool
	flagSet := flag.NewFlagSet("docker-credential-file-gen", flag.ExitOnError)
	flagSet.StringVar(&credFile, "output", ".docker-config.json", "resulting docker credential file path")
	flagSet.BoolVar(&raw, "raw", false, "skip logging and write the output directly to file")
	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "usage: %s [-output output]\n\n", flagSet.Name())
		fmt.Fprintf(flagSet.Output(), "generates a json formatted with all the the docker registry credentials found in the credential store\n\n")
		fmt.Fprintf(flagSet.Output(), "available configuration options:\n\n")
		flagSet.PrintDefaults()
	}

	loggerFlags := log.Lshortfile | log.Lmsgprefix
	errLogger := log.New(os.Stderr, "ERROR ", loggerFlags)

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		errLogger.Println(err)
		os.Exit(1)
	}

	out := io.Writer(os.Stdout)
	if raw {
		out = io.Discard
	}
	logger := log.New(out, "INFO ", loggerFlags)
	logger.Println("loading docker configuration file from system")
	cfgFile := config.LoadDefaultConfigFile(os.Stderr)

	credentials, err := dockerconfig.LoadCredentials(cfgFile)
	if err != nil {
		errLogger.Println(err)
		os.Exit(2)
	}

	logger.Printf("generating docker configuration file: %s\n", credFile)
	if err := writeFile(credFile, credentials); err != nil {
		errLogger.Println(err)
		os.Exit(3)
	}
}

func writeFile(name string, contents dockerconfig.File) error {
	var w io.Writer = os.Stdout
	if name != "-" {
		f, err := os.Create(name)
		if err != nil {
			return err
		}
		defer f.Close()
		w = f
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(contents)
}
