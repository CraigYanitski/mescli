package typeset

import (
    "fmt"
)

type ansiCMD int

const (
    Normal           ansiCMD = 0
    Bold             ansiCMD = 1
    Faint            ansiCMD = 2
    Italics          ansiCMD = 3
    Underline        ansiCMD = 4
    Blink            ansiCMD = 5
    Inverse          ansiCMD = 7
    Hidden           ansiCMD = 8
    CrossOut         ansiCMD = 9
    DoubleUnderline  ansiCMD = 21
    NotItalics       ansiCMD = 23
    NotUnderline     ansiCMD = 24
    NotBlink         ansiCMD = 25
    NotInverse       ansiCMD = 27
    NotHidden        ansiCMD = 28
    NotCrossOut      ansiCMD = 29
    ForeBlack        ansiCMD = 30
    ForeRed          ansiCMD = 31
    ForeGreen        ansiCMD = 32
    ForeYellow       ansiCMD = 33
    ForeBlue         ansiCMD = 34
    ForeMagenta      ansiCMD = 35
    ForeCyan         ansiCMD = 36
    ForeWhite        ansiCMD = 37
    ForeDefault      ansiCMD = 39
    BackBlack        ansiCMD = 40
    BackRed          ansiCMD = 41
    BackGreen        ansiCMD = 42
    BackYellow       ansiCMD = 43
    BackBlue         ansiCMD = 44
    BackMagenta      ansiCMD = 45
    BackCyan         ansiCMD = 46
    BackWhite        ansiCMD = 47
    BackDefault      ansiCMD = 49
)

func formatANSI(codes []ansiCMD) string, error {
    ansi := "\e["
    if len(codes) > 0 {
        ansi += fmt.Sprintf("%d", code[0])
        for i := 1; i < len(codes); i++ {
            ansi += fmt.Sprintf(";%d", code[i]
        }
    } else {
        return "", fmt.Errorf("No ANSI code to apply...")
    ansi += "m"
    return ansi, nil
}





