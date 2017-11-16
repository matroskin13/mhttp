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

func do(host string, method string, httpType string, body []byte) (*mhttp.Response, error) {
	req, err := mhttp.NewRequest(host, method, httpType)

	if err != nil {
		return nil, err
	}

	resp, err := req.Do(body)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

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

		res, err := do(host, method, mhttp.GetTypeByAlias(mType), params)

		if err != nil {
			return cli.NewExitError(err, 0)
		}

		s, err := prettyjson.Format(res.BodyRaw)

		if err != nil {
			return cli.NewExitError(err, 0)
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

	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type, t",
					Value: mhttp.GetTypeByAlias(mhttp.TypeJSON),
				},

				cli.BoolFlag{
					Name: "interactive, i",
				},
			},
			Action: requestAction(http.MethodGet),
		},

		{
			Name:    "post",
			Aliases: []string{"p"},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type, t",
					Value: mhttp.GetTypeByAlias(mhttp.TypeJSON),
				},

				cli.BoolFlag{
					Name: "interactive, i",
				},
			},
			Action: requestAction(http.MethodPost),
		},
	}

	app.Run(os.Args)
}
