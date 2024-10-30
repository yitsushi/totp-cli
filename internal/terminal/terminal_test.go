package terminal_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func prepareIO(input []byte) (io.Reader, *bytes.Buffer, *bytes.Buffer) {
	return bytes.NewReader(input), &bytes.Buffer{}, &bytes.Buffer{}
}

func TestTerminal(t *testing.T) {
	suite.Run(t, &TerminalTestSuite{})
}

type TerminalTestSuite struct {
	suite.Suite
}

func (suite *TerminalTestSuite) TestRead() {
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

		suite.Require().NoError(err)
		suite.Equal(tc.Out, output.String())
		suite.Equal(tc.Value, value)
	}
}

func (suite *TerminalTestSuite) TestConfirm() {
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

		suite.Equal(tc.Out, output.String())
		suite.Equal(tc.Value, value)
	}
}

func (suite *TerminalTestSuite) TestHidden() {
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

		suite.Require().NoError(err)
		suite.Equal(tc.Err, errorOut.String())
		suite.Equal(tc.Out, output.String())
		suite.Equal(tc.Value, value)
	}
}
