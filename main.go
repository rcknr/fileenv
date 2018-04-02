package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

//usage prints fileenv usage
func usage() {
	fmt.Printf("Usage: %s [flags...] /path/to/program [program arguments...]\n", os.Args[0])
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func main() {
	// set up flags
	flag.Usage = usage
	debug := flag.Bool("debug", false, "print debug info messages")
	fail := flag.Bool("fail", false, "immediately exit on warning")
	flag.Parse()

	// check for program path
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	// loop over environment variables
	for _, line := range os.Environ() {

		// get key
		splits := strings.Split(line, "=")
		if len(splits) == 0 {
			continue
		}

		key := splits[0]

		// skip if key doesn't end in _FILE
		if !strings.HasSuffix(strings.ToUpper(key), "_FILE") {
			continue
		}

		path := os.Getenv(key)

		if *debug {
			fmt.Printf("INFO: %s: opening: %s\n", key, path)
		}

		// open and read file
		file, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %s: unable to open %s: %v\n", key, path, err)
			if *fail {
				os.Exit(1)
			}
			continue
		}

		buf := new(bytes.Buffer)

		_, err = io.Copy(buf, file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %s: unable to read %s: %v\n", key, path, err)
			if *fail {
				os.Exit(1)
			}
			continue
		}

		newKey := key[0 : len(key)-5]
		val := strings.TrimSpace(buf.String())

		// set environment variable
		if *debug {
			fmt.Printf("INFO: setting %s=%#v\n", newKey, val)
		}

		if err := os.Setenv(newKey, val); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: %s: unable to set environment variable: %v\n", newKey, err)
			if *fail {
				os.Exit(1)
			}
		}
	}

	// execute program
	path, err := exec.LookPath(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: \"%s\" not found: %v\n", flag.Arg(0), err)
	}

	cmd := exec.Command(path, flag.Args()[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: program terminated with error: %s\n", err)
		// return child exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(ws.ExitStatus())
			}
		}
	}
}
