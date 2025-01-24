package main

import (
    "fmt"
    "github.com/CraigYanitski/mescli/typeset"
)

func main() {
    // Initialise ANSI codes for test
    var codes = []typeset.AnsiCMD{3, 31}
    var reset = []typeset.AnsiCMD{0}

    // Get ANSI strings
    prefix, err := typeset.FormatANSI(codes)
    if err != nil {
        panic(err)
    }
    suffix, err := typeset.FormatANSI(reset)
    if err != nil {
        panic(err)
    }

    // Test printing to terminal
    fmt.Printf("%v%v%v\n", prefix, "A test string!", suffix)
}
