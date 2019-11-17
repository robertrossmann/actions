package toolkit

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_GetMetadata(t *testing.T) {
	t.Run("type", func(t *testing.T) {
		meta := GetMetadata()

		assert.IsType(t, meta, &Metadata{})
	})
}

func Test_NewDebug(t *testing.T) {
	want := "::debug::hello world"
	got := NewDebug("hello world").String()

	assert.Equal(t, want, got)
}

func Test_NewWarning(t *testing.T) {
	want := "::warning::hello world"
	got := NewWarning("hello world").String()

	assert.Equal(t, want, got)
}

func Test_NewError(t *testing.T) {
	want := "::error::hello world"
	got := NewError("hello world").String()

	assert.Equal(t, want, got)
}

func Test_AnnotationFields(t *testing.T) {
	t.Parallel()

	t.Run("File", func(t *testing.T) {
		want := "::debug file=/path/to/file.js::hello world"
		a := NewDebug("hello world")
		a.File = "/path/to/file.js"
		got := a.String()

		assert.Equal(t, want, got)
	})

	t.Run("Line", func(t *testing.T) {
		want := "::debug line=5::hello world"
		a := NewDebug("hello world")
		a.Line = 5
		got := a.String()

		assert.Equal(t, want, got)
	})

	t.Run("Col", func(t *testing.T) {
		want := "::debug col=5::hello world"
		a := NewDebug("hello world")
		a.Col = 5
		got := a.String()

		assert.Equal(t, want, got)
	})

	t.Run("All", func(t *testing.T) {
		want := "::debug file=/test/file.js,line=5,col=4::hello world"
		a := NewDebug("hello world")
		a.File = "/test/file.js"
		a.Line = 5
		a.Col = 4
		got := a.String()

		assert.Equal(t, want, got)
	})
}

func Test_Setenv(t *testing.T) {
	assert.Empty(t, os.Getenv("TEST_ENV_VAR"))
	defer os.Unsetenv("TEST_ENV_VAR")

	want := "::set-env name=TEST_ENV_VAR::testvalue\n"
	got := capture(func() {
		Setenv("TEST_ENV_VAR", "testvalue")
	})

	assert.Equal(t, want, got)
	assert.Equal(t, "testvalue", os.Getenv("TEST_ENV_VAR"))
}

func Test_SetOutput(t *testing.T) {
	want := "::set-output name=testkey::testvalue\n"
	got := capture(func() {
		SetOutput("testkey", "testvalue")
	})

	assert.Equal(t, want, got)
}

func Test_PrependPath(t *testing.T) {
	path := os.Getenv("PATH")
	defer os.Setenv("PATH", path)

	want := "::add-path::/usr/dummy/bin\n"
	got := capture(func() {
		PrependPath("/usr/dummy/bin")
	})

	assert.Contains(t, os.Getenv("PATH"), "/usr/dummy/bin")
	assert.Equal(t, want, got)
}

func Test_SetSecret(t *testing.T) {
	want := "::add-mask::supersecret\n"
	got := capture(func() {
		SetSecret("supersecret")
	})

	assert.Equal(t, want, got)
}

func Test_GetInput(t *testing.T) {
	t.Run("All caps, no spaces", func(t *testing.T) {
		os.Setenv("INPUT_TESTINPUT", "testval")
		defer os.Unsetenv("INPUT_TESTINPUT")

		want := "testval"
		got, _ := GetInput("TESTINPUT")

		assert.Equal(t, want, got)
	})

	t.Run("All caps, with spaces", func(t *testing.T) {
		os.Setenv("INPUT_TEST_INPUT", "testval")
		defer os.Unsetenv("INPUT_TEST_INPUT")

		want := "testval"
		got, _ := GetInput("TEST INPUT")

		assert.Equal(t, want, got)
	})

	t.Run("Mixed caps, no spaces", func(t *testing.T) {
		os.Setenv("INPUT_TESTINPUT", "testval")
		defer os.Unsetenv("INPUT_TESTINPUT")

		want := "testval"
		got, _ := GetInput("TestInput")

		assert.Equal(t, want, got)
	})

	t.Run("No caps, no spaces", func(t *testing.T) {
		os.Setenv("INPUT_TESTINPUT", "testval")
		defer os.Unsetenv("INPUT_TESTINPUT")

		want := "testval"
		got, _ := GetInput("testinput")

		assert.Equal(t, want, got)
	})

	t.Run("Non-existent input", func(t *testing.T) {
		want := ""
		got, err := GetInput("TESTINPUT")

		assert.Equal(t, want, got)
		assert.EqualError(t, err, "Input TESTINPUT not supplied or empty string")
	})

	t.Run("Leading/trailing whitespace in value", func(t *testing.T) {
		os.Setenv("INPUT_TESTINPUT", "  testval\n  ")
		defer os.Unsetenv("INPUT_TESTINPUT")

		want := "testval"
		got, _ := GetInput("testinput")

		assert.Equal(t, want, got)
	})
}

func Test_Error(t *testing.T) {
	want := "::error::hello world\n"
	got := capture(func() {
		Error("hello world")
	})

	assert.Equal(t, want, got)
}

func Test_Warning(t *testing.T) {
	want := "::warning::hello world\n"
	got := capture(func() {
		Warning("hello world")
	})

	assert.Equal(t, want, got)
}

func Test_Debug(t *testing.T) {
	want := "::debug::hello world\n"
	got := capture(func() {
		Debug("hello world")
	})

	assert.Equal(t, want, got)
}
func Test_StartGroup(t *testing.T) {
	want := "::group name=hello world\n"
	got := capture(func() {
		StartGroup("hello world")
	})

	assert.Equal(t, want, got)
}

func Test_EndGroup(t *testing.T) {
	want := "::endgroup\n"
	got := capture(func() {
		EndGroup()
	})

	assert.Equal(t, want, got)
}

func Test_StopCommands(t *testing.T) {
	want := "::stop-commands::hello world\n"
	got := capture(func() {
		StopCommands("hello world")
	})

	assert.Equal(t, want, got)
}

func Test_ResumeCommands(t *testing.T) {
	want := "::hello world::\n"
	got := capture(func() {
		ResumeCommands("hello world")
	})

	assert.Equal(t, want, got)
}

// capture stubs the package's output to stdout and instead stores the output in a buffer.
func capture(f func()) string {
	original := out
	buffer := &bytes.Buffer{}
	out = buffer
	f()
	out = original

	return buffer.String()
}
