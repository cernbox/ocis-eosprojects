package main

import (
	"os"

	"github.com/cernbox/ocis-eosprojects/pkg/command"
	"github.com/cernbox/ocis-eosprojects/pkg/config"
)

func main() {
	if err := command.Execute(config.New()); err != nil {
		os.Exit(1)
	}
}
