package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

var version = time.Now().String()

func main() {
	if err := newApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Name = "gslides"
	app.Usage = `Google Slides command line tool
	see https://github.com/emicklei/gslides for documentation.
`
	// override -v
	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Flags = []cli.Flag{&cli.BoolFlag{Name: "v"}}
	app.Commands = []cli.Command{
		{
			Name:  "export",
			Usage: "Retrieving information related to user accounts",
			Subcommands: []cli.Command{
				{
					Name:  "thumbnails",
					Usage: "Export a PNG file per slide",
					Action: func(c *cli.Context) error {
						return cmdExportThumbnails(c)
					},
					ArgsUsage: `export thumbnails <presentation-id>`,
				},
				{
					Name:  "notes",
					Usage: "Export a TXT file with notes per slide",
					Action: func(c *cli.Context) error {
						return cmdExportNotes(c)
					},
					ArgsUsage: `export notes <presentation-id>`,
				},
			},
		},
		{
			Name:  "append",
			Usage: "Append a slide from one prestentation to another",
			Action: func(c *cli.Context) error {
				return cmdAppendSlide(c)
			},
			ArgsUsage: `append <presentation-id> <other-presentation-id> <slide-index>`,
		},
	}
	return app
}
