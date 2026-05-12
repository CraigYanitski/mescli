package main

import (
    "os"

    "github.com/CraigYanitski/mescli/cmd"
)

func main() {
    err := cmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

