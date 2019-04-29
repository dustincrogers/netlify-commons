package messaging

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/nats-io/go-nats"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func silentIfNil(log logrus.FieldLogger) *logrus.Entry {
	if log == nil {
		l := logrus.New()
		l.Out = ioutil.Discard
		log = logrus.NewEntry(l)
	}
	return log.WithField("component", "nats")
}

func (config *NatsConfig) ConnectToNats(log logrus.FieldLogger, opts ...nats.Option) (*nats.Conn, error) {
	log = silentIfNil(log)

	if err := config.LoadServerNames(); err != nil {
		return nil, errors.Wrap(err, "Failed to discover new servers")
	}

	log.WithFields(config.Fields()).Info("Going to connect to nats servers")
	if len(opts) == 0 {
		opts = []nats.Option{
			ErrorHandler(log),
			nats.MaxReconnects(-1),
		}
	}

	if config.TLS != nil {
		tlsConfig, err := config.TLS.TLSConfig()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to configure TLS")
		}
		if tlsConfig != nil {
			opts = append(opts, nats.Secure(tlsConfig))
			log.Info("Configured TLS connection")
		}
	}

	return nats.Connect(config.ServerString(), opts...)
}

func (config *NatsConfig) ConnectToNatsStreaming(log logrus.FieldLogger, opts ...nats.Option) (*nats.Conn, stan.Conn, error) {
	log = silentIfNil(log)
	if config.ClusterID == "" {
		return nil, nil, errors.New("Must provide a cluster ID to connect to streaming nats")
	}

	if config.ClientID == "" {
		config.ClientID = fmt.Sprintf("generated-%d", time.Now().Nanosecond())
		log.WithField("client_id", config.ClientID).Info("No client ID specified, generating a random one")
	}

	nc, err := config.ConnectToNats(log, opts...)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to connect to nats")
	}

	log.WithFields(config.Fields()).Infof("Connecting to nats streaming cluster %s", config.ClusterID)
	sc, err := stan.Connect(config.ClusterID, config.ClientID, stan.NatsConn(nc))
	if err != nil {
		defer nc.Close()
		return nil, nil, err
	}
	return nc, sc, err
}

func ErrorHandler(log logrus.FieldLogger) nats.Option {
	errLogger := log.WithField("component", "error-logger")
	handler := func(conn *nats.Conn, sub *nats.Subscription, natsErr error) {
		err := natsErr

		l := errLogger.WithFields(logrus.Fields{
			"subject":     sub.Subject,
			"group":       sub.Queue,
			"conn_status": conn.Status(),
		})

		if err == nats.ErrSlowConsumer {
			pendingMsgs, _, perr := sub.Pending()
			if perr != nil {
				err = perr
			} else {
				l = l.WithField("pending_messages", pendingMsgs)
			}
		}

		l.WithError(err).Error("Error while consuming from " + sub.Subject)
	}
	return nats.ErrorHandler(handler)
}
