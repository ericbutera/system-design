package main

import "os"

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	return nil
}

/*
if etl is batch
how can we make a streaming ingestion pipeline

sources:
qualys
tenable
crowdstrike

*/
