package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type input struct {
	Host *string `json:"host"`
}

type output struct {
	Thumbprint string `json:"thumbprint"`
}

var inBytes []byte
var err error
var in *input

func main() {

	in = new(input)

	tf := flag.Bool("terraform", false, "reads data from stdin and writes to stdout/stderr conformant to the external program specification.")

	flag.Parse()

	if *tf {

		inBytes, err = ioutil.ReadAll(os.Stdin)

		if err != nil {
			os.Stderr.WriteString("error reading input from stdin\n")
			os.Exit(1)
		}

		err := json.Unmarshal(inBytes, in)

		if err != nil || in.Host == nil {
			os.Stderr.WriteString("query object must be in the form { host: \"my.host.com\" }\n")
			os.Exit(1)
		}

	} else {

		args := os.Args[1:]

		if len(args) == 0 {
			os.Stderr.WriteString("FQDN is required\n")
			os.Exit(1)
		}

		if len(args) != 1 {
			os.Stderr.WriteString("only one argument permitted\n")
			os.Exit(1)
		}

		in.Host = &args[0]

	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:443", *in.Host), &tls.Config{})

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("failed to connect: %s\n", err))
		os.Exit(1)
	}

	err = conn.Handshake()

	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("failed to handshake: %s\n", err))
		os.Exit(1)
	}

	state := conn.ConnectionState()

	if len(state.PeerCertificates) > 0 {

		thumbprint := sha1.Sum(state.PeerCertificates[len(state.PeerCertificates)-1].Raw)

		formatted := fmt.Sprintf("%X", thumbprint)

		if !*tf {
			_, err = os.Stdout.WriteString(formatted)

			if err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}

		res := &output{
			Thumbprint: formatted,
		}

		out, err := json.Marshal(res)

		if err != nil {
			os.Stderr.WriteString("internal error, unable to read certificate\n")
			os.Exit(1)
		}

		_, err = os.Stdout.Write(out)

		if err != nil {
			os.Exit(1)
		}

		os.Exit(0)

	}

}
