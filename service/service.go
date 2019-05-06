package service

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/netlify/netlify-commons/nconf"
	"github.com/netlify/netlify-commons/tracing"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type RootArgs struct {
	prefix  string
	envFile string
}

func AddRootArgs(cmd *cobra.Command) *RootArgs {
	args := new(RootArgs)
	cmd.PersistentFlags().StringVarP(&args.envFile, "config", "c", "", "An env file to load for config")
	cmd.PersistentFlags().StringVarP(&args.prefix, "prefix", "p", "", "A prefix to use for env vars")
	return args
}

func (args RootArgs) InitService(name, version string, config interface{}) (logrus.FieldLogger, error) {
	if args.envFile != "" {
		if err := godotenv.Load(args.envFile); err != nil {
			return nil, errors.Wrapf(err, "Failed to load configuration from file: %s", args.envFile)
		}
	}

	// Load logging info under the LOG_* environment vars
	logConfig := struct{ Log nconf.LoggingConfig }{
		Log: nconf.LoggingConfig{
			QuoteEmptyFields: true,
		},
	}
	if err := envconfig.Process(args.prefix, &logConfig); err != nil {
		return nil, errors.Wrapf(err, "Failed to load logging configuration from the enviroment")
	}

	log, err := nconf.ConfigureLogging(&logConfig.Log)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to logger")
	}
	log = log.WithField("version", version)
	log.Debug("Configured logging")

	// Load the tracing configuration using the TRACING_* environment vars
	traceConfig := struct{ Tracing tracing.Config }{}
	if err := envconfig.Process(args.prefix, &traceConfig); err != nil {
		return nil, errors.Wrap(err, "Failed to load tracing configurtion from the environment")
	}
	tracing.Configure(&traceConfig.Tracing, name)
	log.Debug("Configured tracing")

	// Load the config for the service
	if config != nil {
		if err := envconfig.Process(args.prefix, config); err != nil {
			return nil, errors.Wrap(err, "Failed to load the service configuration")
		}
		log.Debugf("Configured service %s", name)
	}

	return log, nil
}
