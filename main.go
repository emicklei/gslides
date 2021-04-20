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
					ArgsUsage: `export thumbnails`,
				},
			},
		},
	}
	return app
}
