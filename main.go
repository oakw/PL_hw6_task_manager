package main

import (
	"task_manager/src/app"

	"github.com/londek/reactea"
)

func main() {
	program := reactea.NewProgram(app.New())

	if err := program.Start(); err != nil {
		panic(err)
	}
}
