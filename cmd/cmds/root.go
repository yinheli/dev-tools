package cmds

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
	"github.com/yinheli/dev-tools/version"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	configDir string

	rootCmd = &cobra.Command{
		Use:     version.AppName,
		Short:   "dev tools",
		Version: version.Build,
	}
)

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	configDir = path.Join(home, ".dev-tools")

	_, err = os.Stat(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(configDir, 0700)
			if err != nil {
				panic(err)
			}
		}
	}
}

func readUntil(r *strings.Reader, chs ...rune) string {
	var buf strings.Builder
	for {
		c, _, err := r.ReadRune()
		if err == io.EOF {
			break
		}

		for _, ch := range chs {
			if c == ch {
				return buf.String()
			}
		}

		buf.WriteRune(c)
	}

	return buf.String()
}

func cmdLineToArgs(cmdLine string) []string {
	args := make([]string, 0, 10)

	r := strings.NewReader(cmdLine)

	for {
		c, _, err := r.ReadRune()
		if err == io.EOF {
			break
		}
		v := ""
		if c == '"' {
			v = readUntil(r, '"')
		} else if c == '\'' {
			v = readUntil(r, '\'')
		} else if c == ' ' {
			continue
		} else {
			v += string(c)
			v += readUntil(r, ' ', '\'', '"')
		}
		if v != "" {
			args = append(args, v)
		}
	}
	return args
}

func Execute() {
	if len(os.Args) > 1 {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	}

	historyFile := filepath.Join(configDir, ".history")
	history, err := os.Open(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			history, _ = os.Create(historyFile)
		}
	}

	line := liner.NewLiner()
	defer func() {
		if history != nil {
			line.WriteHistory(history)
		}
		line.Close()
	}()

	if history != nil {
		line.ReadHistory(history)
	}

	line.SetCtrlCAborts(true)
	line.SetMultiLineMode(true)

	line.SetCompleter(func(line string) []string {
		s := make([]string, 0, 10)
		for _, it := range rootCmd.Commands() {
			name := it.Name()

			if strings.HasPrefix(name, line) {
				s = append(s, name)
			}
		}
		return s
	})

	for {
		cmdLine, err := line.Prompt("> ")
		switch err {
		case liner.ErrPromptAborted:
			return
		case nil:
		default:
			fmt.Println(err)
		}

		line.AppendHistory(cmdLine)
		args := cmdLineToArgs(cmdLine)
		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
	}
}
