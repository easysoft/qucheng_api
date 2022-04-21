package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var versionTpl = `cne-api version:
 Version:           %v
 Go version:        %v
 Git commit:        %v
 Built:             %v
 OS/Arch:           %v
 Experimental:      true
`

var (
	Version       string
	BuildDate     string
	GitCommitHash string
	Mode          string
)

const (
	defaultVersion       = "0.0.0"
	defaultGitCommitHash = "a1b2c3d4"
	defaultBuildDate     = "Mon Aug  3 15:06:50 2020"
)

func NewCmdVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "show build version",
		Run:   showVersion,
	}
	return cmd
}

func showVersion(cmd *cobra.Command, args []string) {
	if Version == "" {
		Version = defaultVersion
	}
	if BuildDate == "" {
		BuildDate = defaultBuildDate
	}
	if GitCommitHash == "" {
		GitCommitHash = defaultGitCommitHash
	}
	osarch := fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH)
	fmt.Printf(versionTpl, Version, runtime.Version(), GitCommitHash, BuildDate, osarch)
}
