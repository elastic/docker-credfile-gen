# Docker credential file generator

Generates a `.docker-config.json` file with the registry authentication credentials that are accessed through the Docker Go SDK, the output of this file is meant to be used for the development environments so the generated file containing the docker registry credentials can be copied directly to the remote VM without necessarily relying on what the current docker configuration file contains.

## Usage

```
usage: docker-credfile-gen [-output output]

generates a json formatted with all the the docker registry credentials found in the credential store

available configuration options:

  -output string
    	resulting docker credential file path (default ".docker-config.json")
```

## Building

If you've got Go set up, you can go ahead and build the binary using `make build`, which will download all the dependencies and compile a binary in `bin/docker-credfile-gen`.

## Building in Docker

If you've do not have Go set up, but have Docker, you can also build the binary in Docker container with `make docker-build`, which will also compile a binary in `bin/docker-credfile-gen`.
