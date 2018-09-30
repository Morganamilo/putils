package main // import "github.com/Morganamilo/putils"

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var tmp *os.File

func editor() []string {
	if len(os.Args) > 1 {
		return os.Args[1:]
	}

	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			return nil
		}
	}

	return strings.Fields(editor)
}

func execEditor(editor string, args []string) {
	tty, err := os.Open("/dev/tty")
	handleErr(err, "failed to open tty", 1)

	cmd := exec.Command(editor, args...)
	cmd.Stdin = tty
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	handleErr(err, "failed to open editor", 1)
}

func printErr(str string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", str, err)
}

func handleErr(err error, str string, ret int) {
	if err != nil {
		printErr(str, err)
		cleanupTmp()
		os.Exit(ret)
	}
}

func cleanupTmp() int {
	ret := 0
	if tmp != nil {
		err := tmp.Close()
		if err != nil {
			printErr("failed to close tmp file", err)
			ret = 1
		}
		err = os.Remove(tmp.Name())
		if err != nil {
			printErr("failed to remove tmp file", err)
			ret = 1
		}
		tmp = nil
	}

	return ret
}

func main() {
	var err error

	args := editor()
	if args == nil {
		handleErr(fmt.Errorf("%s", "please set $VISUAL or $EDITOR"), "no suitable editor", 2)
	}

	tmp, err = ioutil.TempFile(os.TempDir(), "pedit.")
	handleErr(err, "failed to get tmp file", 1)

	_, err = io.Copy(tmp, os.Stdin)
	handleErr(err, "failed to copy stdin to tmp file", 1)

	editor := args[0]
	args = append(args, tmp.Name())
	execEditor(editor, args[1:])

	err = tmp.Sync()
	handleErr(err, "failed to sync tmp file", 1)

	_, err = tmp.Seek(0, io.SeekStart)
	handleErr(err, "failed to seek tmp file", 1)

	_, err = io.Copy(os.Stdout, tmp)
	handleErr(err, "failes to copy tmp file to stdout", 1)

	os.Exit(cleanupTmp())
}
