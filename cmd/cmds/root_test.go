package cmds

import (
	"fmt"
	"testing"
)

func TestCmdlineToArgs(t *testing.T) {
	// line := `"131" "121"  'af' fafa`
	// line := `fafa faf "afa" "'s'" fa`
	line := `aaa`
	fmt.Println(line)
	args := cmdLineToArgs(line)
	for _, c := range args {
		fmt.Println(c)
	}
}

func Test02(t *testing.T) {
	rootCmd.SetArgs([]string{"config", "--help"})
	rootCmd.Execute()
}
