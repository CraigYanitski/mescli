package typeset

import (
    "fmt"
)

type AnsiCMD int

const (
    Normal           AnsiCMD = 0
    Bold             AnsiCMD = 1
    Faint            AnsiCMD = 2
    Italics          AnsiCMD = 3
    Underline        AnsiCMD = 4
    Blink            AnsiCMD = 5
    Inverse          AnsiCMD = 7
    Hidden           AnsiCMD = 8
    CrossOut         AnsiCMD = 9
    DoubleUnderline  AnsiCMD = 21
    NotItalics       AnsiCMD = 23
    NotUnderline     AnsiCMD = 24
    NotBlink         AnsiCMD = 25
    NotInverse       AnsiCMD = 27
    NotHidden        AnsiCMD = 28
    NotCrossOut      AnsiCMD = 29
    ForeBlack        AnsiCMD = 30
    ForeRed          AnsiCMD = 31
    ForeGreen        AnsiCMD = 32
    ForeYellow       AnsiCMD = 33
    ForeBlue         AnsiCMD = 34
    ForeMagenta      AnsiCMD = 35
    ForeCyan         AnsiCMD = 36
    ForeWhite        AnsiCMD = 37
    ForeDefault      AnsiCMD = 39
    BackBlack        AnsiCMD = 40
    BackRed          AnsiCMD = 41
    BackGreen        AnsiCMD = 42
    BackYellow       AnsiCMD = 43
    BackBlue         AnsiCMD = 44
    BackMagenta      AnsiCMD = 45
    BackCyan         AnsiCMD = 46
    BackWhite        AnsiCMD = 47
    BackDefault      AnsiCMD = 49
)

func FormatANSI(codes []AnsiCMD) (string, error) {
    ansi := "\033["
    if len(codes) > 0 {
        ansi += fmt.Sprintf("%d", codes[0])
        for i := 1; i < len(codes); i++ {
            ansi += fmt.Sprintf(";%d", codes[i])
        }
    } else {
        return "", fmt.Errorf("No ANSI code to apply...")
    }
    ansi += "m"
    // ansi := fmt.Sprintf("\033[%vm", codes)
    return ansi, nil
}





