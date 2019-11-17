// Package toolkit implements basic functionality to operate in a GitHub Action context. It supports
// commenting on certain lines of code in a pull request, reading current action's metadata, setting
// an action's outputs etc.
package toolkit

import (
	"fmt"
	"io"
	"os"
	"strings"
)

var out io.Writer = os.Stdout

func println(message string) (n int, err error) {
	return fmt.Fprintln(out, message)
}

// Metadata shows information about current action's environment, runtime & event which triggered the workflow.
type Metadata struct {
	Action     string
	Actor      string
	BaseRef    string
	EventName  string
	EventPath  string
	HeadRef    string
	Ref        string
	Repository string
	RunnerOS   string
	Sha        string
	Workflow   string
	Workspace  string
}

// GetMetadata retrieves the current action run's metadata.
func GetMetadata() *Metadata {
	meta := &Metadata{}
	meta.Action = os.Getenv("GITHUB_ACTION")
	meta.Actor = os.Getenv("GITHUB_ACTOR")
	meta.BaseRef = os.Getenv("GITHUB_BASE_REF")
	meta.EventName = os.Getenv("GITHUB_EVENT_NAME")
	meta.EventPath = os.Getenv("GITHUB_EVENT_PATH")
	meta.HeadRef = os.Getenv("GITHUB_HEAD_REF")
	meta.Ref = os.Getenv("GITHUB_REF")
	meta.Repository = os.Getenv("GITHUB_REPOSITORY")
	meta.RunnerOS = os.Getenv("RUNNER_OS")
	meta.Sha = os.Getenv("GITHUB_SHA")
	meta.Workflow = os.Getenv("GITHUB_WORKFLOW")
	meta.Workspace = os.Getenv("GITHUB_WORKSPACE")

	return meta
}

// Annotation represents a comment on a specific location in a file.
type Annotation struct {
	level   string
	message string
	File    string
	Line    int
	Col     int
}

// String serialises an annotation into Action-compatible console entry.
func (a Annotation) String() string {
	var params = make([]string, 0)

	if len(a.File) != 0 {
		params = append(params, fmt.Sprintf("file=%s", a.File))
	}

	// Lines are 1-indexed so a Line of 0 means uninitialised
	if a.Line != 0 {
		params = append(params, fmt.Sprintf("line=%d", a.Line))
	}

	// Columns are 1-indexed so a Col of 0 means uninitialised
	if a.Col != 0 {
		params = append(params, fmt.Sprintf("col=%d", a.Col))
	}

	output := fmt.Sprintf("::%s", a.level)

	if len(params) != 0 {
		output += " " + strings.Join(params, ",")
	}

	// Escape carriage return and newline characters
	// @see https://github.com/actions/toolkit/blob/master/packages/core/src/command.ts#L71
	a.message = strings.ReplaceAll(a.message, "\r", "%0D")
	a.message = strings.ReplaceAll(a.message, "\n", "%0A")

	return fmt.Sprintf("%s::%s", output, a.message)
}

// NewDebug creates a new debug-level annotation.
// You should set File, Line & Col positions after creation.
func NewDebug(message string) Annotation {
	return Annotation{level: "debug", message: message}
}

// NewWarning creates a new warning-level annotation.
// You should set File, Line & Col positions after creation.
func NewWarning(message string) Annotation {
	return Annotation{level: "warning", message: message}
}

// NewError creates a new error-level annotation.
// You should set File, Line & Col positions after creation.
func NewError(message string) Annotation {
	return Annotation{level: "error", message: message}
}

// Setenv creates or updates an environment variable for any actions running next in a job.
// The action that creates or updates the environment variable does not have access to the new
// value, but all subsequent actions in a job will have access. Environment variables are
// case-sensitive and you can include punctuation.
func Setenv(key string, value string) (n int, err error) {
	os.Setenv(key, value)
	return println(fmt.Sprintf("::set-env name=%s::%s", key, value))
}

// SetOutput sets an action's output parameter.
// Output parameters are defined in an action's metadata file. You will receive an error if you
// attempt to set an output value that was not declared in the action's metadata file.
func SetOutput(name string, value string) (n int, err error) {
	return println(fmt.Sprintf("::set-output name=%s::%s", name, value))
}

// PrependPath prepends a directory to the system PATH variable for all subsequent actions in the
// current job. The currently running action cannot access the new path variable.
func PrependPath(path string) (n int, err error) {
	parts := []string{path, os.Getenv("PATH")}

	if err := os.Setenv("PATH", strings.Join(parts, string(os.PathListSeparator))); err != nil {
		return 0, err
	}

	return println(fmt.Sprintf("::add-path::%s", path))
}

// SetSecret registers a secret which will get masked from logs.
func SetSecret(secret string) (n int, err error) {
	return println(fmt.Sprintf("::add-mask::%s", secret))
}

// GetInput gets the value of an input.  The value is also trimmed.
func GetInput(name string) (string, error) {
	key := "INPUT_" + strings.ReplaceAll(strings.ToUpper(name), " ", "_")
	value := strings.TrimSpace(os.Getenv(key))

	if len(value) == 0 {
		return "", fmt.Errorf("Input %s not supplied or empty string", name)
	}

	return value, nil
}

// Annotate writes an Annotation to the log and to the pull request if file/line/col position is set.
func Annotate(annotation Annotation) (n int, err error) {
	return println(annotation.String())
}

// Error Writes an error-level message to the action output.
func Error(message string) (n int, err error) {
	return Annotate(NewError(message))
}

// Warning writes a warning-level message to the action output.
func Warning(message string) (n int, err error) {
	return Annotate(NewWarning(message))
}

// Debug writes a debug-level message to the action output. Only visible if debugging is enabled.
func Debug(message string) (n int, err error) {
	return Annotate(NewDebug(message))
}

// StartGroup starts an output group. Output will be foldable in this group until the next EndGroup.
func StartGroup(name string) (n int, err error) {
	return println(fmt.Sprintf("::group name=%s", name))
}

// EndGroup ends an output group.
func EndGroup() (n int, err error) {
	return println("::endgroup")
}

// StopCommands stops processing any logging commands.
// This allows you to log anything without accidentally triggering any command.
func StopCommands(endtoken string) (n int, err error) {
	return println(fmt.Sprintf("::stop-commands::%s", endtoken))
}

// ResumeCommands resumes processing logging commands.
func ResumeCommands(endtoken string) (n int, err error) {
	return println(fmt.Sprintf("::%s::", endtoken))
}
