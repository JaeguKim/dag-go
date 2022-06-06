package shellexecmd

import (
	"testing"
)

func TestDag_ScriptRunner(t *testing.T) {
	s := "./test01.sh"

	cmd, r := ScriptRunner(s)
	StartThenWait(cmd)
	ch := Reply(r)
	PrintOutput(ch)
}

func TestScriptRunnerString(t *testing.T) {
	script := `
	set -e
	sleep 1
	echo "Hello World"
	sleep 1
	echo "one"
	sleep 1
	echo "two"
	sleep 1
	echo "three"
	sleep 1
	echo "four"
	sleep 1
	echo "Sleep 10s"
	sleep 10
	echo "End"`

	cmd, r := ScriptRunnerString(script)
	StartThenWait(cmd)
	ch := Reply(r)
	PrintOutput(ch)
}

func TestRunner(t *testing.T) {
	script := `
	set -e
	sleep 1
	echo "Hello World"
	sleep 1
	echo "one"
	sleep 1
	echo "two"
	sleep 1
	echo "three"
	sleep 1
	echo "four"
	sleep 1
	echo "Sleep 10s"
	sleep 10
	echo "End"`

	Runner(script)
}