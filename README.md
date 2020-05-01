# Thumbprint

A simple Terraform compatible CLI tool to retrieve a Root CA thumbprint, to use with [IAM OIDC Identity Providers](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_create_oidc_verify-thumbprint.html), for example. 

## Installation

Simple `go build main.go`, or `go get github.com/matt-potter/thumbprint`

Ensure your `$PATH` is set up to include your `$GOPATH/bin` directory 

## Options

```
➜  ~ thumbprint --help 
Usage of thumbprint:
  -terraform
        reads data from stdin and writes to stdout/stderr conformant to the external program specification.
➜  ~ 
```

## Use as a CLI

```
➜  ~ thumbprint google.com 
DFE2070C79E7FF36A925FFA327FFE3DEECF8F9C2
```

## Use in terraform 

Pass in the `-terraform` flag. The tool will now read input from `STDIN` as per the Terraform spec. The tool expects a json object to `STDIN` in the form of {"host": " FQDN value "}, which is achieved with a Terraform data external block as below. The result is accessible in the "thumbprint" field of the result map.

```
data "external" "thumbprint" {
  program = ["thumbprint", "-terraform"]

  query = {
    host = ** FQDN HERE **
  }
}

resource "aws_iam_openid_connect_provider" "k8s" {
  url = ...
  client_id_list = [
    [...]
  ]
  thumbprint_list = [data.external.thumbprint.result["thumbprint"]]
}
```