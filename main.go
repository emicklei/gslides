package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var version = time.Now().String()
var Verbose = false

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
	cli.VersionFlag = &cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Flags = []cli.Flag{&cli.BoolFlag{Name: "v"}}
	app.Commands = []*cli.Command{
		{
			Name:  "export",
			Usage: "Retrieving information related to user accounts",
			Subcommands: []*cli.Command{
				{
					Name:  "thumbnails",
					Usage: "Export a PNG file per slide",
					Action: func(c *cli.Context) error {
						Verbose = c.Bool("v")
						return cmdExportThumbnails(c)
					},
					ArgsUsage: `export thumbnails <presentation-id>`,
				},
				{
					Name:  "notes",
					Usage: "Export a TXT file with notes per slide",
					Action: func(c *cli.Context) error {
						Verbose = c.Bool("v")
						return cmdExportNotes(c)
					},
					ArgsUsage: `export notes <presentation-id>`,
				},
				{
					Name:  "pdf",
					Usage: "Export a presentation to a PDF file",
					Action: func(c *cli.Context) error {
						Verbose = c.Bool("v")
						return cmdExportPDF(c)
					},
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:        "o",
							Usage:       "output file name",
							DefaultText: "gslides-output.pdf",
						},
					},
					ArgsUsage: `export pdf <presentation-id>`,
				},
			},
		},
		{
			Name:  "append",
			Usage: "Append a slide from one prestentation to another",
			Action: func(c *cli.Context) error {
				Verbose = c.Bool("v")
				return cmdAppendSlide(c)
			},
			ArgsUsage: `append <presentation-id> <other-presentation-id> <slide-index>`,
		},
		{
			Name:  "inspect",
			Usage: "Inspect a presentation or slide",
			Action: func(c *cli.Context) error {
				Verbose = c.Bool("v")
				return cmdInspect(c)
			},
			ArgsUsage: `inspect <presentation-id> <slide-index>`,
		},
		{
			Name:  "list",
			Usage: "List all presentations from your Google Drive folder",
			Action: func(c *cli.Context) error {
				Verbose = c.Bool("v")
				return cmdList(c)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "owner",
					Usage: "filter by owner account",
				},
				&cli.StringFlag{
					Name:  "folder",
					Usage: "files in other folder than root",
				},
			},
			ArgsUsage: `list <folder>`,
		},
	}
	return app
}
