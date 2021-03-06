// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package ioctl

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"

	"github.com/iotexproject/iotex-core/ioctl/config"
	"github.com/iotexproject/iotex-core/ioctl/util"
)

type (
	// Client defines the interface of an ioctl client
	Client interface {
		// Start starts the client
		Start(context.Context) error
		// Stop stops the client
		Stop(context.Context) error
		// Config returns the config of the client
		Config() config.Config
		// APIServiceClient returns an API service client
		APIServiceClient(APIServiceConfig) (iotexapi.APIServiceClient, error)
		// SelectTranslation select a translation based on UILanguage
		SelectTranslation(map[config.Language]string) (string, config.Language)
		// AskToConfirm asks user to confirm from terminal, true to continue
		AskToConfirm() bool
		// ReadSecret reads password from terminal
		ReadSecret() (string, error)
	}

	// APIServiceConfig defines a config of APIServiceClient
	APIServiceConfig struct {
		Endpoint string
		Insecure bool
	}

	client struct {
		cfg  config.Config
		conn *grpc.ClientConn
		// TODO: merge into config
		lang config.Language
	}
)

// NewClient creates a new ioctl client
func NewClient() Client {
	return &client{}
}

var confirmMessages = map[config.Language]string{
	config.English: "Do you want to continue? [yes/NO]",
	config.Chinese: "是否继续？【是/否】",
}

func (c *client) Start(context.Context) error {
	return nil
}

func (c *client) Stop(context.Context) error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
		c.conn = nil
	}
	return nil
}

func (c *client) Config() config.Config {
	return c.cfg
}

func (c *client) AskToConfirm() bool {
	msg, lang := c.SelectTranslation(confirmMessages)
	fmt.Println(msg)
	var confirm string
	fmt.Scanf("%s", &confirm)
	switch lang {
	case config.Chinese:
		return strings.EqualFold(confirm, "是")
	default: // config.English
		return strings.EqualFold(confirm, "yes")
	}
}

func (c *client) SelectTranslation(trls map[config.Language]string) (string, config.Language) {
	trl, ok := trls[c.lang]
	if ok {
		return trl, c.lang
	}

	trl, ok = trls[config.English]
	if !ok {
		panic("failed to pick a translation")
	}

	return trl, config.English
}

func (c *client) ReadSecret() (string, error) {
	// TODO: delete util.ReadSecretFromStdin, and move code to here
	return util.ReadSecretFromStdin()
}

func (c *client) APIServiceClient(cfg APIServiceConfig) (iotexapi.APIServiceClient, error) {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return nil, err
		}
	}
	var err error
	if cfg.Insecure {
		c.conn, err = grpc.Dial(cfg.Endpoint, grpc.WithInsecure())
	} else {
		c.conn, err = grpc.Dial(cfg.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	}
	if err != nil {
		return nil, err
	}
	return iotexapi.NewAPIServiceClient(c.conn), nil
}
