package command

import (
	"fmt"
	"strings"

	"github.com/mohae/cli"
	"github.com/mohae/contour"
	"github.com/mohae/kraul/app"
)

// IndexCommand is a Command implementation that says hello world
type IndexCommand struct {
	UI cli.Ui
}

// Help prints the help text for the run sub-command.
func (c *IndexCommand) Help() string {
	helpText := `
Usage: kraul index [flags] <urls...>

index indexs the website urls specified in the kraul's
config file. Its behavior is controlled both by the config 
file and by any rules set in custom processors. 

Ouput is in the form of a csv file.

Optionally, 1 or more URLs can be passed and they will be
indexed instead. If the scheme is omitted, it is assumed to
be http.

    $ kraul index

	$ kraul index www.example.com http://hoopybits.com
	
kraul flags:

    --lower=(true, false)    true lowercases the output.

    --logging=(true, false)  enable/disable log output
    -l                       alias to --logging
`
	return strings.TrimSpace(helpText)
}

// Run runs the index command; optionally, the config file can be added.
func (c *IndexCommand) Run(args []string) int {
	// set up the command flags
	contour.SetFlagSetUsage(func() {
		c.UI.Output(c.Help())
	})

	// Filter the flags from the args and update the config with them.
	// The args remaining after being filtered are returned.
	filteredArgs, err := contour.FilterArgs(args)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	err = app.SetLogging()
	if err != nil {
		c.UI.Error(fmt.Sprintf("setup and configuration of application logging failed: %s", err))
		return 1
	}

	// Run the command in the package; 0 is the start url, or should be.
	message, err := app.Index(filteredArgs[0])
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(message)
	return 0
}

// Synopsis provides a precis of the index command.
func (c *IndexCommand) Synopsis() string {
	ret := `Indexs urls specified in a config file using its defined rules.
`
	return ret
}
