package main

import (
	"github.com/src-d/ghsync/cmd/ghsync/subcmd"

	"gopkg.in/src-d/go-cli.v0"
	_ "gopkg.in/src-d/go-queue.v1/amqp"
	_ "gopkg.in/src-d/go-queue.v1/memory"
)

// rewritten during the CI build step
var (
	version = "master"
	build   = "dev"
)

var app = cli.New("ghsync", version, build, "GitHub metadata sync")

func main() {
	app.AddCommand(&subcmd.ShallowCommand{})
	app.AddCommand(&subcmd.DeepCommand{})
	app.AddCommand(&subcmd.MigrateCommand{})

	app.RunMain()
}
