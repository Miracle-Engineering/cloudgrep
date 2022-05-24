package cmd

import (
	"fmt"
	"github.com/run-x/cloudgrep/pkg/version"
	"io"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

type versionOptions struct {
	short    bool
	template string
}

var vO versionOptions

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the cloudgrep version",
	Long: `Show the version for cloudgrep.
This will print a representation the version of cloudgrep.
The output will look something like this:
version.BuildInfo{Version:"v0.1.2", GitCommit:"af379c8ce85305912b1e726d1de3d7a052946d52", BuildTime:"2022-05-23T04:25:01Z", GoVersion:"go1.18.2"}
- Version is the semantic version of the release.
- GitCommit is the SHA for the commit that this version was built from.
- BuildTime is the UTC time when the binary was built.
- GoVersion is the version of Go that was used to compile cloudgrep.

When using the --template flag the following properties are available to use in
the template:
- .Version is the semantic version of the release.
- .GitCommit is is the SHA for the commit that this version was built from.
- .BuildTime is the UTC time when the binary was built
- .GoVersion contains the version of Go that cloudgrep was compiled with
For example, --template='Version: {{.Version}}' outputs 'Version: v0.1.2'.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return vO.run(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	f := versionCmd.Flags()
	f.BoolVar(&vO.short, "short", false, "print the version number")
	f.StringVar(&vO.template, "template", "", "template for version string format")
}

func (o *versionOptions) run(out io.Writer) error {
	if o.template != "" {
		tt, err := template.New("_").Parse(o.template)
		if err != nil {
			return err
		}
		return tt.Execute(out, version.Get())
	}
	fmt.Fprintln(out, formatVersion(o.short))
	return nil
}

func formatVersion(short bool) string {
	v := version.Get()
	if short {
		if len(v.GitCommit) >= 7 {
			return fmt.Sprintf("%s+g%s", v.Version, v.GitCommit[:7])
		}
		return v.Version
	}
	return fmt.Sprintf("%#v", v)
}
