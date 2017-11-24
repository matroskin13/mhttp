package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/hokaccha/go-prettyjson"
	"github.com/urfave/cli"

	"encoding/json"
	"mhttp"
)

var hostIsNotDefined = errors.New("host is not defined")

type SavedRequest struct {
	URI         string            `json:"uri"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
	Params      []string          `json:"params"`
	Type        string            `json:"type"`
	RequestBody []byte            `json:"-"`
	Flags       map[string]bool   `json:"-"`
}

func doAction(savedRequest *SavedRequest) error {
	req, err := mhttp.NewRequest(savedRequest.URI, savedRequest.Method, savedRequest.Headers)

	if err != nil {
		return cli.NewExitError(err, 0)
	}

	res, err := req.Do(savedRequest.RequestBody)

	if err != nil {
		return cli.NewExitError(err, 0)
	}

	fmt.Println(string(res.BodyRaw))

	s, err := prettyjson.Format(res.BodyRaw)

	if err != nil {
		return cli.NewExitError(err, 0)
	}

	if flag, _ := savedRequest.Flags["request-info"]; flag {
		requestInfo, _ := req.GetPrettyRequest()

		fmt.Println(requestInfo)

		return nil
	}

	fmt.Println(string(s))

	if flag, _ := savedRequest.Flags["interactive"]; flag {
		err := mhttp.InitInteractive(res.BodyRaw)

		if err != nil {
			return cli.NewExitError(err, 0)
		}
	}

	return nil
}

func prependRequest(c *cli.Context, savedRequest *SavedRequest) error {
	var requestBody []byte

	host := c.Args().First()

	savedRequest.Params = append(savedRequest.Params, c.Args().Tail()...)

	if savedRequest.Method != http.MethodGet {
		parser, err := mhttp.ParseParams(savedRequest.Params)

		if err != nil {
			return cli.NewExitError(err, 0)
		}

		requestBody, err = parser.ToJSON()

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

	headers["Content-Type"] = mhttp.GetTypeByAlias(savedRequest.Type)

	for key, h := range headers {
		savedRequest.Headers[key] = h
	}

	savedRequest.RequestBody = requestBody
	savedRequest.Flags = map[string]bool{
		"request-info": c.Bool("request-info"),
		"interactive":  c.Bool("interactive"),
	}

	return nil
}

func requestAction(method string) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		host := c.Args().First()
		mType := c.String("type")

		savedRequest := SavedRequest{
			Method:  method,
			Headers: make(map[string]string),
			Params:  c.Args().Tail(),
			URI:     host,
			Type:    mType,
		}

		prependRequest(c, &savedRequest)

		if name := c.String("save"); name != "" {
			config, err := mhttp.GetOrCreateConfig()

			if err != nil {
				return cli.NewExitError(err, 0)
			}

			err = config.AddJSONVar("__mhttp", "save_"+name, savedRequest)

			if err != nil {
				return cli.NewExitError(err, 0)
			}

			err = config.Save()

			if err != nil {
				return cli.NewExitError(err, 0)
			}

			return nil
		}

		err := doAction(&savedRequest)

		if err != nil {
			return cli.NewExitError(err, 0)
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

		cli.StringFlag{
			Name: "save, s",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Flags:   requestFlags,
			Action:  requestAction(http.MethodGet),
		},

		{
			Name:    "post",
			Aliases: []string{"p"},
			Flags:   requestFlags,
			Action:  requestAction(http.MethodPost),
		},

		{
			Name:    "use",
			Aliases: []string{"u"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "interactive, i",
				},

				cli.StringSliceFlag{
					Name: "headers, H",
				},
			},
			Action: func(c *cli.Context) error {
				config, err := mhttp.GetOrCreateConfig()

				if err != nil {
					return cli.NewExitError(err, 0)
				}

				value, err := config.GetVar("__mhttp", "save_"+c.Args().First())

				if err != nil {
					return cli.NewExitError(err, 0)
				}

				savedRequest := SavedRequest{}

				err = json.Unmarshal([]byte(value), &savedRequest)

				if err != nil {
					return cli.NewExitError(err, 0)
				}

				err = prependRequest(c, &savedRequest)

				if err != nil {
					return err
				}

				err = doAction(&savedRequest)

				if err != nil {
					return err
				}

				return nil
			},
		},

		{
			Name: "config",
			Action: func(c *cli.Context) error {
				config, err := mhttp.GetOrCreateConfig()

				if err != nil {
					return cli.NewExitError(err, 0)
				}

				fmt.Println(config)

				return nil
			},
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
