package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/hokaccha/go-prettyjson"
	"github.com/urfave/cli"

	"mhttp"
)

var hostIsNotDefined = errors.New("host is not defined")

func requestAction(method string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		var params []byte

		host := c.Args().First()
		mType := c.String("type")

		if method != http.MethodGet {
			parser, err := mhttp.ParseParams(c.Args().Tail())

			if err != nil {
				return cli.NewExitError(err, 0)
			}

			params, err = parser.ToJSON()

			if err != nil {
				return cli.NewExitError(err, 0)
			}
		}

		if host == "" {
			return cli.NewExitError(hostIsNotDefined, 0)
		}

		headers, err := mhttp.PrependHeaders(c.StringSlice("headers"))

		if err != nil {
			return cli.NewExitError(err, 0)
		}

		headers["Content-Type"] = mhttp.GetTypeByAlias(mType)

		req, err := mhttp.NewRequest(host, method, headers)

		if err != nil {
			return cli.NewExitError(err, 0)
		}

		res, err := req.Do(params)

		if err != nil {
			return cli.NewExitError(err, 0)
		}

		s, err := prettyjson.Format(res.BodyRaw)

		if err != nil {
			return cli.NewExitError(err, 0)
		}

		if c.Bool("request-info") {
			requestInfo, _ := req.GetPrettyRequest()

			fmt.Println(requestInfo)

			return nil
		}

		fmt.Println(string(s))

		if c.Bool("interactive") {
			err := mhttp.InitInteractive(res.BodyRaw)

			if err != nil {
				return cli.NewExitError(err, 0)
			}
		}

		return nil
	}
}

func main() {
	app := cli.NewApp()

	app.Name = "Monster HTTP"
	app.Usage = ""
	app.Description = "Command line HTTP client"
	app.Version = "0.0.1"

	requestFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "type, t",
			Value: mhttp.TypeJSON,
		},

		cli.BoolFlag{
			Name: "interactive, i",
		},

		cli.BoolFlag{
			Name: "request-info, ri",
		},

		cli.StringSliceFlag{
			Name: "headers, H",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Flags: requestFlags,
			Action: requestAction(http.MethodGet),
		},

		{
			Name:    "post",
			Aliases: []string{"p"},
			Flags: requestFlags,
			Action: requestAction(http.MethodPost),
		},

		{
			Name:    "var",
			Aliases: []string{"v"},
			Subcommands: cli.Commands{
				{
					Name: "set",
					Action: func(c *cli.Context) error {
						config, err := mhttp.GetOrCreateConfig()

						if err != nil {
							return cli.NewExitError(err, 0)
						}

						key := c.Args().First()
						value := c.Args().Get(1)

						config.AddVar("", key, value)
						config.Save()

						return nil
					},
				},

				{
					Name: "get",
					Action: func(c *cli.Context) error {
						config, err := mhttp.GetOrCreateConfig()

						if err != nil {
							return cli.NewExitError(err, 0)
						}

						key := c.Args().First()

						value, err := config.GetVar("", key)

						if err != nil {
							return cli.NewExitError(err, 0)
						}

						fmt.Println(value)

						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
