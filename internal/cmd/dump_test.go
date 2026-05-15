package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
)

func TestDump(t *testing.T) {
	suite.Run(t, &DumpTestSuite{})
}

type DumpTestSuite struct {
	suite.Suite
}

func (suite *DumpTestSuite) TestDumpWritesYAML() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "secret"}},
	}
	stor := newFakeStorage(ns)
	outputPath := filepath.Join(suite.T().TempDir(), "dump.yaml")

	err := executeDump(stor, outputPath)

	suite.Require().NoError(err)

	data, readErr := os.ReadFile(outputPath)
	suite.Require().NoError(readErr)
	suite.Contains(string(data), "myns")
	suite.Contains(string(data), "github")
}

func (suite *DumpTestSuite) TestDumpInvalidOutputPath() {
	stor := newFakeStorage()

	err := executeDump(stor, "/nonexistent/path/dump.yaml")

	suite.Require().Error(err)
}
