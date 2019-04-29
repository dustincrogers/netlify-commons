package messaging

import (
	"fmt"
	"strings"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"
	"github.com/sirupsen/logrus"

	"github.com/netlify/netlify-commons/discovery"
	"github.com/netlify/netlify-commons/nconf"
)

type NatsConfig struct {
	TLS           *nconf.TLSConfig `mapstructure:"tls_conf"`
	DiscoveryName string           `split_words:"true" mapstructure:"discovery_name"`
	Servers       []string         `mapstructure:"servers"`

	// for streaming
	ClusterID string `mapstructure:"cluster_id" envconfig:"cluster_id"`
	ClientID  string `mapstructure:"client_id" envconfig:"client_id"`
	StartPos  string `mapstructure:"start_pos" split_words:"true"`

	Subject string `mapstructure:"subject"`
	Group   string `mapstructure:"group"`
}

func (c *NatsConfig) LoadServerNames() error {
	if c.DiscoveryName == "" {
		return nil
	}

	natsURLs := []string{}
	endpoints, err := discovery.DiscoverEndpoints(c.DiscoveryName)
	if err != nil {
		return err
	}

	for _, endpoint := range endpoints {
		natsURLs = append(natsURLs, fmt.Sprintf("nats://%s:%d", endpoint.Target, endpoint.Port))
	}

	c.Servers = natsURLs
	return nil
}

// ServerString will build the proper string for nats connect
func (config *NatsConfig) ServerString() string {
	return strings.Join(config.Servers, ",")
}

func (config *NatsConfig) Fields() logrus.Fields {
	f := logrus.Fields{
		"servers": strings.Join(config.Servers, ","),
		"group":   config.Group,
		"subject": config.Subject,
	}

	if config.TLS != nil {
		f["ca_files"] = strings.Join(config.TLS.CAFiles, ",")
		f["key_file"] = config.TLS.KeyFile
		f["cert_file"] = config.TLS.CertFile
	}

	if config.ClusterID != "" {
		f["client_id"] = config.ClientID
		f["cluster_id"] = config.ClusterID
	}

	return f
}

func (config *NatsConfig) StartPoint() (stan.SubscriptionOption, error) {
	switch v := strings.ToLower(config.StartPos); v {
	case "all":
		return stan.DeliverAllAvailable(), nil
	case "last":
		return stan.StartWithLastReceived(), nil
	case "new":
		return stan.StartAt(pb.StartPosition_NewOnly), nil
	case "", "first":
		return stan.StartAt(pb.StartPosition_First), nil
	}
	return nil, fmt.Errorf("Unknown start position '%s', possible values are all, last, new, first and ''", config.StartPos)
}
