package typeset

import (
	"fmt"
	"strings"
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

func formatANSI(codes []ansiCMD) (string, error) {
    // Start ANSI sequence
    ansi := "\033["

    // Append codes if they exist
    if len(codes) > 0 {
        ansi += fmt.Sprintf("%d", codes[0])
        for i := 1; i < len(codes); i++ {
            ansi += fmt.Sprintf(";%d", codes[i])
        }
    } else {
        return "", fmt.Errorf("no ANSI code to apply")
    }

    // Finish ANSI code and return
    ansi += "m"
    // ansi := fmt.Sprintf("\033[%vm", codes)
    
    return ansi, nil
}

func getCode(format string) (ansiCMD, error) {
    var code ansiCMD
    // Switch depending on input format
    switch strings.ToLower(strings.Replace(format, " ", "", -1)) {
    // Behaviour codes
    case "reset", "default":
        code = Normal
    case "bold":
        code = Bold
    case "faint":
        code = Faint
    case "italics":
        code = Italics
    case "underline":
        code = Underline
    case "blink":
        code = Blink
    case "inverse":
        code = Inverse
    case "hidden":
        code = Hidden
    case "crossout":
        code = CrossOut
    case "doubleunderline":
        code = DoubleUnderline
    case "notitalics":
        code = NotItalics
    case "notunderline":
        code = NotUnderline
    case "notblink":
        code = NotBlink
    case "notinverse":
        code = NotInverse
    case "nothidden":
        code = NotHidden
    case "notcrossout":
        code = NotCrossOut
    // Foreground color codes
    case "black", "foreblack":
        code = ForeBlack
    case "red", "forered":
        code = ForeRed
    case "green", "foregreen":
        code = ForeGreen
    case "yellow", "foreyellow":
        code = ForeYellow
    case "blue", "foreblue":
        code = ForeBlue
    case "magenta", "foremagenta":
        code = ForeMagenta
    case "cyan", "forecyan":
        code = ForeCyan
    case "white", "forewhite":
        code = ForeWhite
    case "foredefault":
        code = ForeDefault
    // Background color codes
    case "backblack":
        code = BackBlack
    case "backred":
        code = BackRed
    case "backgreen":
        code = BackGreen
    case "backyellow":
        code = BackYellow
    case "backblue":
        code = BackBlue
    case "backmagenta":
        code = BackMagenta
    case "backcyan":
        code = BackCyan
    case "backwhite":
        code = BackWhite
    case "backdefault":
        code = BackDefault
    // Default error
    default:
        return Normal, fmt.Errorf("error parsing format %v.\n" +
            "You can specify 'normal', 'bold', 'faint', 'italics', " +
            "'underline', 'blink', 'crossout', as well as those keys preceded by 'not'.\n" +
            "You can also specify the colors 'black', 'red', 'green', 'yellow', 'blue', " +
            "'magenta', 'cyan', and 'white', but you should also specify 'fore' or 'back' " +
            "before the color to set the foreground or background. " +
            "The default is the foreground", format)
    }
    return code, nil
}

func FormatString(line string, format []string) (string, error) {
    // Return unformatted string if no format is specified
    if len(format) == 0 {
        return line, nil
    }

    // Go through formats to get ANSI codes
    var codes []ansiCMD
    for i := 0; i < len(format); i++ {
        code, err := getCode(format[i])
        if err != nil {
            return "", fmt.Errorf("error getting ANSI CODE: %v", err)
        }
        codes = append(codes, code)
    }

    // Format string prefix
    prefix, err := formatANSI(codes)
    if err != nil {
        return "", fmt.Errorf("error combining ANSI codes: %v", err)
    }

    // Format reset string
    suffix, err := formatANSI([]ansiCMD{Normal})
    if err != nil {
        return "", fmt.Errorf("error creating suffix: %v", err)
    }

    return prefix + line + suffix, nil
}





