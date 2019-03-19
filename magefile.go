// +build mage

package main

import (
	"bufio"
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	appName = "dev-tools"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// tidy code
func Fmt() error {
	packages := strings.Split("cmd pkg version", " ")
	return sh.Run("gofmt", append([]string{"-s", "-l", "-w", "mage.go", "magefile.go"}, packages...)...)
}

// for local machine build
func Build() error {
	return buildTarget("", runtime.GOOS, runtime.GOARCH, nil)
}

// build linux
func Linux() error {
	return buildTarget("", "linux", "amd64", nil)
}

// build windows
func Windows() error {
	return buildTarget("", "windows", "amd64", nil)
}

// build to target (cross build)
func buildTarget(name, OS, arch string, envs map[string]string) error {
	tag := tag()
	if name == "" {
		name = fmt.Sprintf("%s-%s-%s-%s", appName, tag, OS, arch)
	} else {
		name = fmt.Sprintf("%s-%s-%s", appName, tag, name)
	}
	dir := fmt.Sprintf("dist/%s", name)
	target := fmt.Sprintf("%s/%s", dir, appName)

	args := make([]string, 0, 10)
	args = append(args, "build", "-o", target)
	args = append(args, "-ldflags", flags(), "cmd/main.go")

	fmt.Println("build", name)
	env := make(map[string]string)
	env["GOOS"] = OS
	env["GOARCH"] = arch
	env["CGO_ENABLED"] = "0"

	if envs != nil {
		for k, v := range envs {
			env[k] = v
		}
	}

	if err := sh.RunWith(env, mg.GoCmd(), args...); err != nil {
		return err
	}

	sh.Run("tar", "-czf", fmt.Sprintf("%s.tar.gz", dir), "-C", "dist", name)
	return nil
}

func flags() string {
	timestamp := time.Now().Format(time.RFC3339)
	hash := hash()
	tag := tag()
	mod := mod()
	var buf strings.Builder
	buf.WriteString("-s -w ")
	buf.WriteString(fmt.Sprintf(`-X "%s/version.Build=%s-%s" `, mod, tag, hash))
	buf.WriteString(fmt.Sprintf(`-X "%s/version.BuildTime=%s" `, mod, timestamp))
	buf.WriteString(`-extldflags "-static" `)
	return buf.String()
}

// tag returns the git tag for the current branch or "" if none.
func tag() string {
	s, _ := sh.Output("bash", "-c", "git describe --abbrev=0 --tags 2> /dev/null")
	if s == "" {
		return "dev"
	}
	return s
}

// hash returns the git hash for the current repo or "" if none.
func hash() string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return hash
}

func mod() string {
	f, err := os.Open("go.mod")
	if err == nil {
		reader := bufio.NewReader(f)
		line, _, _ := reader.ReadLine()
		return strings.Replace(string(line), "module ", "", 1)
	}
	return ""
}

// cleanup all build files
func Clean() {
	sh.Rm("dist")
}
