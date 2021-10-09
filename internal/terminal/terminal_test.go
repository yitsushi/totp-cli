package terminal_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func prepareIO(input []byte) (io.Reader, *bytes.Buffer, *bytes.Buffer) {
	return bytes.NewReader(input), &bytes.Buffer{}, &bytes.Buffer{}
}

func TestTerminal_Read(t *testing.T) {
	type testCase struct {
		Prompt string
		In     string
		Out    string
		Value  string
	}

	testCases := []testCase{
		{Prompt: "pro:", In: "test in\n", Out: "pro: ", Value: "test in"},
		{Prompt: "pro:", In: "test\nin\n", Out: "pro: ", Value: "test"},
		{Prompt: "", In: "   test   input   \n", Out: "", Value: "test   input"},
	}

	for _, tc := range testCases {
		input, output, errorOut := prepareIO([]byte(tc.In))

		term := terminal.New(input, output, errorOut)
		value, err := term.Read(tc.Prompt)

		assert.Nil(t, err)
		assert.Equal(t, tc.Out, output.String())
		assert.Equal(t, tc.Value, value)
	}
}

func TestTerminal_Confirm(t *testing.T) {
	type testCase struct {
		Prompt string
		In     string
		Out    string
		Value  bool
	}

	testCases := []testCase{
		{Prompt: "should be yes:", In: "y\n", Out: "should be yes: ", Value: true},
		{Prompt: "yes:", In: "yes\n", Out: "yes: ", Value: true},
		{Prompt: "sure:", In: "sure\n", Out: "sure: ", Value: true},
		{Prompt: "sure   :", In: "sure   \n", Out: "sure   : ", Value: true},
		{Prompt: "", In: "anything else\n", Out: "", Value: false},
		{Prompt: "", In: "\n", Out: "", Value: false},
		{Prompt: "", In: "", Out: "", Value: false},
	}

	for _, tc := range testCases {
		input, output, errorOut := prepareIO([]byte(tc.In))

		term := terminal.New(input, output, errorOut)
		value := term.Confirm(tc.Prompt)

		assert.Equal(t, tc.Out, output.String())
		assert.Equal(t, tc.Value, value)
	}
}

func TestTerminal_Hidden(t *testing.T) {
	type testCase struct {
		Prompt string
		In     string
		Err    string
		Out    string
		Value  string
	}

	testCases := []testCase{
		{
			Prompt: "prompt:",
			In:     "asd\n",
			Err:    "prompt: ***\n",
			Out:    "",
			Value:  "asd",
		},
	}

	for _, tc := range testCases {
		input, output, errorOut := prepareIO([]byte(tc.In))

		term := terminal.New(input, output, errorOut)
		value, err := term.Hidden(tc.Prompt)

		assert.Nil(t, err)
		assert.Equal(t, tc.Err, errorOut.String())
		assert.Equal(t, tc.Out, output.String())
		assert.Equal(t, tc.Value, value)
	}
}
