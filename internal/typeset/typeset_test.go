package typeset_test

import (
	"fmt"
	"testing"

	"github.com/CraigYanitski/mescli/internal/typeset"
)

func TestBlueText(t *testing.T) {
    type testCase struct {
        message   string
        format    []string
        expected  string
    }

    tests := []testCase{
        {"some", []string{"bold"}, 
            "\033[1msome\033[0m"},
        {"text that", []string{"faint"}, 
            "\033[2mtext that\033[0m"},
        {"I want to", []string{"italics"}, 
            "\033[3mI want to\033[0m"},
        {"test typesetting", []string{"underline"}, 
            "\033[4mtest typesetting\033[0m"},
        {"in different styles", []string{"crossout"}, 
            "\033[9min different styles\033[0m"},
        {"and combinations", []string{"inverse"}, 
            "\033[7mand combinations\033[0m"},
        {"of colours and ", []string{"blue"}, 
            "\033[34mof colours and \033[0m"},
        {" using \n\n escape \n\n sequences. ", []string{"faint", "yellow"}, 
            "\033[2;33m using \n\n escape \n\n sequences. \033[0m"},
    }

    failCount := 0
    passCount := 0
	
    fmt.Println("\n\nTesting typesetting")

    for _, test := range tests {
        fmt.Println("----------------------------------------")
        fmt.Printf("Typesetting %q with format %q...\n", test.message, test.format)
        
        formattedMessage, err := typeset.FormatString(test.message, test.format)
        if err != nil {
            t.Errorf("error encountered during Typesetting: %v", err)
            failCount++
            continue
        }
        
        if formattedMessage != test.expected {
            failCount++
            t.Errorf(`
Inputs:    message: %q, format %q
Expected:  formattedMessage: %s
Actual:    formattedMessage: %s
`, test.message, test.format, test.expected, formattedMessage)
        } else {
            passCount++
            fmt.Printf(`
Inputs:    message: %q, format %q
Expected:  formattedMessage: %s
Actual:    formattedMessage: %s
`, test.message, test.format, test.expected, formattedMessage)
        }
    }
    
	fmt.Println("========================================")
	fmt.Printf("%d passed, %d failed\n\n\n", passCount, failCount)
}
