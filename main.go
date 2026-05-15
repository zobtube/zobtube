package main

import (
	"context"
	"embed"
	"fmt"
	"net/mail"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	altsrc "github.com/urfave/cli-altsrc/v3"
	"github.com/urfave/cli-altsrc/v3/yaml"
	"github.com/urfave/cli/v3"

	"github.com/zobtube/zobtube/cli/passwordreset"
	"github.com/zobtube/zobtube/cli/server"
	"github.com/zobtube/zobtube/internal/config"
)

//go:embed web
var webFS embed.FS

// goreleaser build-time variables
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// global logger
var logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

// @title        ZobTube
// @description  ZobTube is a video management system.
// @contact.name ZobTube Issues
// @contact.url  https://github.com/zobtube/zobtube/issues
// @license.name MIT
// @license.url  https://github.com/zobtube/zobtube?tab=MIT-1-ov-file#readme
// @BasePath     /api
// @schemes      http https
func main() {
	var configurationFile string

	cmd := &cli.Command{
		Name:      "zobtube",
		Usage:     "passion of the zob, lube for the tube!",
		Version:   fmt.Sprintf("%s (commit %s), built at %s", version, commit, date),
		Copyright: "(c) 2025 ZobTube",
		Authors: []any{
			mail.Address{Name: "sblablaha", Address: "sblablaha@gmail.com"},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config-file",
				Usage:       "path to configuration file",
				Sources:     cli.EnvVars("ZT_CONFIG_FILE"),
				Value:       "config.yml",
				Destination: &configurationFile,
			},
			&cli.BoolFlag{
				Name:  "gin-debug",
				Usage: "enables gin debugging mode",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_GIN_DEBUG"),
					yaml.YAML("log.gin.debug", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
			},
			&cli.IntFlag{
				Name:  "log-level",
				Usage: "select log verbosity (5: panic / 4: fatal / 3: error / 2: warn / 1: info / 0: debug / -1: trace)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_LOG_LEVEL"),
					yaml.YAML("log.level", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
				Value:       1,
				DefaultText: "1 - info",
				Action: func(ctx context.Context, cmd *cli.Command, v int) error {
					if v < -1 || v > 5 {
						return fmt.Errorf("parameter log-level value %v out of range (must be between -1 and 5)", v)
					}
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "server-bind",
				Usage: "address the http server will bind to",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_SERVER_BIND"),
					yaml.YAML("server.bind", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
				Value: "0.0.0.0:8069",
			},
			&cli.StringFlag{
				Name:  "db-driver",
				Usage: "database driver to use (sqlite or postgresql)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_DB_DRIVER"),
					yaml.YAML("db.driver", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
				Value: "sqlite",
			},
			&cli.StringFlag{
				Name:  "db-connstring",
				Usage: "connection string to the database",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_DB_CONNSTRING"),
					yaml.YAML("db.connstring", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
				Value: "zobtube.sqlite",
			},
			&cli.StringFlag{
				Name:  "media-path",
				Usage: "path to the media folder",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_MEDIA_PATH"),
					yaml.YAML("media.path", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
				Value: "data",
			},
			&cli.StringFlag{
				Name:  "metadata-type",
				Usage: "metadata storage type (filesystem or s3)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_TYPE"),
					yaml.YAML("metadata.type", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
				Value: "filesystem",
			},
			&cli.StringFlag{
				Name:  "metadata-path",
				Usage: "filesystem path for metadata storage (defaults to media-path)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_PATH"),
					yaml.YAML("metadata.path", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
			},
			&cli.StringFlag{
				Name:  "metadata-s3-bucket",
				Usage: "S3 bucket for metadata storage",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_S3_BUCKET"),
					yaml.YAML("metadata.s3.bucket", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
			},
			&cli.StringFlag{
				Name:  "metadata-s3-region",
				Usage: "S3 region for metadata storage",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_S3_REGION"),
					yaml.YAML("metadata.s3.region", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
				Value: "us-east-1",
			},
			&cli.StringFlag{
				Name:  "metadata-s3-prefix",
				Usage: "optional S3 key prefix for metadata storage",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_S3_PREFIX"),
					yaml.YAML("metadata.s3.prefix", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
			},
			&cli.StringFlag{
				Name:  "metadata-s3-endpoint",
				Usage: "optional S3 endpoint for metadata storage (e.g. Minio)",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_S3_ENDPOINT"),
					yaml.YAML("metadata.s3.endpoint", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
			},
			&cli.StringFlag{
				Name:  "metadata-s3-access-key-id",
				Usage: "optional S3 access key for metadata storage",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_S3_ACCESS_KEY_ID"),
					yaml.YAML("metadata.s3.access_key_id", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
			},
			&cli.StringFlag{
				Name:  "metadata-s3-secret-access-key",
				Usage: "optional S3 secret key for metadata storage",
				Sources: cli.NewValueSourceChain(
					cli.EnvVar("ZT_METADATA_S3_SECRET_ACCESS_KEY"),
					yaml.YAML("metadata.s3.secret_access_key", altsrc.NewStringPtrSourcer(&configurationFile)),
				),
			},
		},
		Action: startServer,
		Commands: []*cli.Command{
			{
				Name:   "server",
				Action: startServer,
				Usage:  "start zobtube server, default action if no command passed",
			},
			{
				Name:     "password-reset",
				Category: "user",
				Usage:    "reset password of a user interactively",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "user-id",
						Usage: "user id to reset password. If empty, will list all users with their ids",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					return passwordreset.Run(cmd, &logger)
				},
			},
		},
	}

	err := cmd.Run(context.Background(), os.Args)
	if err != nil {
		logger.Error().Err(err).Send()
	}
}

func metadataParamsFromCmd(cmd *cli.Command) config.MetadataParams {
	return config.MetadataParams{
		Type:              cmd.String("metadata-type"),
		Path:              cmd.String("metadata-path"),
		S3Bucket:          cmd.String("metadata-s3-bucket"),
		S3Region:          cmd.String("metadata-s3-region"),
		S3Prefix:          cmd.String("metadata-s3-prefix"),
		S3Endpoint:        cmd.String("metadata-s3-endpoint"),
		S3AccessKeyID:     cmd.String("metadata-s3-access-key-id"),
		S3SecretAccessKey: cmd.String("metadata-s3-secret-access-key"),
	}
}

func startServer(ctx context.Context, cmd *cli.Command) error {
	return server.Start(&server.Parameters{
		Ctx:     ctx,
		Cmd:     cmd,
		Logger:  &logger,
		Version: version,
		Commit:  commit,
		Date:    date,
		WebFS:    &webFS,
		Metadata: metadataParamsFromCmd(cmd),
	})
}
