package mescli

import (
    "fmt"
    "github.com/CraigYanitski/mescli/typeset"
)

func main() {
    var codes = []int{1, 31}
    var reset = []int{0}
    fmt.Printf("%v %v %v", typeset.formatANSI(codes), "A test string!", typeset.formatANSI(reset))
}
