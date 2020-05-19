# Yasmim [WIP]

[![Build Status via Travis CI](https://travis-ci.org/tsouza/yasmim.svg?branch=master)](https://travis-ci.org/tsouza/yasmim)

`yasmim` is a golang library that helps to implement chain of responsibility design pattern.

## Installation

If you are not using an IDE with go modules support, install it with go:
```
$ go get -u github.com/tsouza/yasmim
```
Import it into your code
```
import "github.com/tsouza/yasmim"
```

## Quick Start

Register a command and run it
```
package main

import (
	"github.com/tsouza/yasmim"
	"github.com/tsouza/yasmim/pkg/command"
	"github.com/tsouza/yasmim/pkg/log"
)

const Name = "my_first_command"

func main() {
  	yasmim.Register(func(define command.Define) {
  		define.Command("my_first_command").
  			Handler(func(_ command.Runtime, _ *log.Logger, _, _ interface{}) error {
                // my_first_command Handler
            })
  	})

    _ = yasmim.New().Run("my_first_command", nil, nil)
}
```

## License

Code and documentation released under [The MIT License (MIT)](LICENSE).