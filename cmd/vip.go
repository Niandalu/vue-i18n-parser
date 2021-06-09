package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"

	"github.com/niandalu/vue-i18n-parser/internal/collector"
	"github.com/niandalu/vue-i18n-parser/internal/feeder"
	"github.com/niandalu/vue-i18n-parser/internal/reader"
	"github.com/urfave/cli"
)

func main() {
	projectRoot, rootMissing := os.Getwd()

	if rootMissing != nil {
		log.Fatal(rootMissing)
	}

	app := prepareApp(projectRoot)
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func prepareApp(projectRoot string) *cli.App {
	app := cli.NewApp()
	app.Name = "Vue I18n Parser"
	app.Usage = "Automate vue app translation process"

	app.Commands = []cli.Command{
		{
			Name:    "collect",
			Aliases: []string{"c"},
			Usage:   "Collect translation KVs under the path",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "diff, d",
					Usage: "Only collect modified files",
				},

				cli.StringFlag{
					Name:  "file, f",
					Value: "./to.be.translated.csv",
					Usage: "Collect all the translations KVs into `FILE",
				},

				cli.StringFlag{
					Name:  "ignore",
					Value: "(node_modules|tourist.yml)",
					Usage: "Files to be ignored, only accept regexp right now",
				},

				cli.StringFlag{
					Name:  "languages, l",
					Value: "cn,en",
					Usage: "Languages your app supported, separated by comma. And the first value is considered to be mandatory language",
				},
			},
			Action: func(c *cli.Context) error {
				diffOnly := c.Bool("diff")
				var files []reader.TranslationFile

				for _, relativeDir := range strings.Split(c.Args().First(), ",") {
					files = append(files, reader.Run(relativeDir, c.String("ignore"))...)
				}
				languages := strings.Split(c.String("languages"), ",")
				sheet := collector.Run(files, languages, diffOnly)

				target := c.String("file")
				f, e := os.Create(target)
				if e != nil {
					log.Fatalf("Could not write %v, %v", target, e)
				}
				defer f.Close()
				f.Write([]byte("\uFEFF")) // BOM

				w := csv.NewWriter(f)
				for _, row := range sheet {
					w.Write(row)
				}
				w.Flush()

				return nil
			},
		},
		{
			Name:    "feed",
			Aliases: []string{"f"},
			Usage:   "Read the translated scripts and write records back to your app",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "patch, p",
					Usage: "Use patch(incremnetal) mode",
				},
				cli.IntFlag{
					Name:  "indent, i",
					Value: 2,
					Usage: "yaml indention",
				},
			},
			Action: func(c *cli.Context) error {
				patchMode := c.Bool("patch")
				indention := c.Int("indent")

				feeder.Run(projectRoot, c.Args().First(),
					feeder.Options{
						PatchMode: patchMode,
						Indent:    indention,
					},
				)
				return nil
			},
		},
	}

	return app
}
