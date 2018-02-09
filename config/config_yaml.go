package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
	"time"

	"github.com/DataDog/datadog-trace-agent/backoff"
	writerconfig "github.com/DataDog/datadog-trace-agent/writer/config"

	"github.com/DataDog/datadog-trace-agent/utils"
)

// YamlAgentConfig is a structure used for marshaling the datadog.yaml configuration
// available in Agent versions >= 6

type traceWriter struct {
	MaxSpansPerPayload     int                    `yaml:"max_spans_per_payload"`
	FlushPeriod            int                    `yaml:"flush_period_seconds"`
	UpdateInfoPeriod       int                    `yaml:"update_info_period_seconds"`
	QueueablePayloadSender queueablePayloadSender `yaml:"queueable_payload_sender"`
}

type serviceWriter struct {
	UpdateInfoPeriod       int                    `yaml:"update_info_period_seconds"`
	FlushPeriod            int                    `yaml:"flush_period_seconds"`
	QueueablePayloadSender queueablePayloadSender `yaml:"queueable_payload_sender"`
}

type statsWriter struct {
	UpdateInfoPeriod       int                    `yaml:"update_info_period_seconds"`
	QueueablePayloadSender queueablePayloadSender `yaml:"queueable_payload_sender"`
}

type queueablePayloadSender struct {
	MaxAge            int   `yaml:"queue_max_age_seconds"`
	MaxQueuedBytes    int64 `yaml:"queue_max_bytes"`
	MaxQueuedPayloads int   `yaml:"queue_max_payloads"`
	BackoffDuration   int   `yaml:"exp_backoff_max_duration_seconds"`
	BackoffBase       int   `yaml:"exp_backoff_base_milliseconds"`
	BackoffGrowth     int   `yaml:"exp_backoff_growth_base"`
}

type traceAgent struct {
	Enabled            bool    `yaml:"enabled"`
	Env                string  `yaml:"env"`
	ExtraSampleRate    float64 `yaml:"extra_sample_rate"`
	MaxTracesPerSecond float64 `yaml:"max_traces_per_second"`
	Ignore             string  `yaml:"ignore_resource"`
	ReceiverPort       int     `yaml:"receiver_port"`
	ConnectionLimit    int     `yaml:"connection_limit"`
	NonLocalTraffic    string  `yaml:"trace_non_local_traffic"`

	TraceWriter   traceWriter   `yaml:"trace_writer"`
	ServiceWriter serviceWriter `yaml:"service_writer"`
	StatsWriter   statsWriter   `yaml:"stats_writer"`
}

//YamlAgentConfig is the Primary Object we retrieve from Datadog.yaml
type YamlAgentConfig struct {
	APIKey   string `yaml:"api_key"`
	HostName string `yaml:"hostname"`

	StatsdHost   string `yaml:"bind_host"`
	ReceiverHost string ""

	StatsdPort int    `yaml:"StatsdPort"`
	LogLevel   string `yaml:"log_level"`

	DefaultEnv string `yaml:"env"`

<<<<<<< HEAD
	TraceAgent traceAgent `yaml:"apm_config"`
=======
	TraceAgent struct {
		Enabled            bool    `yaml:"enabled"`
		Env                string  `yaml:"env"`
		ExtraSampleRate    float64 `yaml:"extra_sample_rate"`
		MaxTracesPerSecond float64 `yaml:"max_traces_per_second"`
		Ignore             string  `yaml:"ignore_resource"`
		ReceiverPort       int     `yaml:"receiver_port"`
		ConnectionLimit    int     `yaml:"connection_limit"`
		NonLocalTraffic    string  `yaml:"trace_non_local_traffic"`

		//TODO Merge these into config
		TraceWriter struct {
			MaxSpansPerPayload int   `yaml:"max_spans_per_payload"`
			FlushPeriod        int   `yaml:"flush_period_seconds"`
			UpdateInfoPeriod   int   `yaml:"update_info_period_seconds"`
			MaxAge             int   `yaml:"queue_max_age_seconds"`
			MaxQueuedBytes     int64 `yaml:"queue_max_bytes"`
			MaxQueuedPayloads  int   `yaml:"queue_max_payloads"`
			BackoffDuration    int   `yaml:"exp_backoff_max_duration_seconds"`
			BackoffBase        int   `yaml:"exp_backoff_base_milliseconds"`
			BackoffGrowth      int   `yaml:"exp_backoff_growth_base"`
		} `yaml:"trace_writer"`
		ServiceWriter struct {
			FlushPeriod       int   `yaml:"flush_period_seconds"`
			UpdateInfoPeriod  int   `yaml:"'update_info_period_seconds"`
			MaxAge            int   `yaml:"queue_max_age_seconds"`
			MaxQueuedBytes    int64 `yaml:"queue_max_bytes"`
			MaxQueuedPayloads int   `yaml:"queue_max_payloads"`
			BackoffDuration   int   `yaml:"exp_backoff_max_duration_seconds"`
			BackoffBase       int   `yaml:"exp_backoff_base_milliseconds"`
			BackoffGrowth     int   `yaml:"exp_backoff_growth_base"`
		} `yaml:"service_writer"`
		StatsWriter struct {
			UpdateInfoPeriod  int   `yaml:"update_info_period_seconds"`
			MaxAge            int   `yaml:"queue_max_age_seconds"`
			MaxQueuedBytes    int64 `yaml:"queue_max_bytes"`
			MaxQueuedPayloads int   `yaml:"queue_max_payloads"`
			BackoffDuration   int   `yaml:"exp_backoff_max_duration_seconds"`
			BackoffBase       int   `yaml:"exp_backoff_base_milliseconds"`
			BackoffGrowth     int   `yaml:"exp_backoff_growth_base"`
		} `yaml:"stats_writer"`
	} `yaml:"apm_config"`
>>>>>>> 44da32f2a5fc61f55b9c097bb522a9c2288a58b1
}

// NewYamlIfExists returns a new YamlAgentConfig if the given configPath is exists.
func NewYamlIfExists(configPath string) (*YamlAgentConfig, error) {
	var yamlConf YamlAgentConfig
	if utils.PathExists(configPath) {
		fileContent, err := ioutil.ReadFile(configPath)
		if err = yaml.Unmarshal([]byte(fileContent), &yamlConf); err != nil {
			return nil, fmt.Errorf("parse error: %s", err)
		}
		return &yamlConf, nil
	}
	return nil, nil
}

func mergeYamlConfig(agentConf *AgentConfig, yc *YamlAgentConfig) (*AgentConfig, error) {
	agentConf.APIKey = yc.APIKey
	agentConf.HostName = yc.HostName
	agentConf.Enabled = yc.TraceAgent.Enabled
	agentConf.DefaultEnv = yc.DefaultEnv

	agentConf.ReceiverPort = yc.TraceAgent.ReceiverPort
	agentConf.ExtraSampleRate = yc.TraceAgent.ExtraSampleRate
	agentConf.MaxTPS = yc.TraceAgent.MaxTracesPerSecond

	agentConf.Ignore["resource"] = strings.Split(yc.TraceAgent.Ignore, ",")

	agentConf.ConnectionLimit = yc.TraceAgent.ConnectionLimit

	//Allow user to specify a different ENV for APM Specifically
	if yc.TraceAgent.Env != "" {
		agentConf.DefaultEnv = yc.TraceAgent.Env
	}

	if yc.StatsdHost != "" {
		yc.ReceiverHost = yc.StatsdHost
	}

	//Respect non_local_traffic
	if v := strings.ToLower(yc.TraceAgent.NonLocalTraffic); v == "yes" || v == "true" {
		yc.StatsdHost = "0.0.0.0"
		yc.ReceiverHost = "0.0.0.0"
	}

	agentConf.StatsdHost = yc.StatsdHost
	agentConf.ReceiverHost = yc.ReceiverHost

	agentConf.ServiceWriterConfig = readServiceWriterConfigYaml(yc.TraceAgent.ServiceWriter)
	agentConf.StatsWriterConfig = readStatsWriterConfigYaml(yc.TraceAgent.StatsWriter)
	agentConf.TraceWriterConfig = readTraceWriterConfigYaml(yc.TraceAgent.TraceWriter)
	return agentConf, nil
}

func readServiceWriterConfigYaml(yc serviceWriter) writerconfig.ServiceWriterConfig {
	c := writerconfig.DefaultServiceWriterConfig()

	if yc.FlushPeriod > 0 {
		c.FlushPeriod = utils.GetDuration(yc.FlushPeriod)
	}

	if yc.UpdateInfoPeriod > 0 {
		c.UpdateInfoPeriod = utils.GetDuration(yc.UpdateInfoPeriod)
	}

	c.SenderConfig = readQueueablePayloadSenderConfigYaml(yc.QueueablePayloadSender)
	return c
}

func readStatsWriterConfigYaml(yc statsWriter) writerconfig.StatsWriterConfig {
	c := writerconfig.DefaultStatsWriterConfig()

	if yc.UpdateInfoPeriod > 0 {
		c.UpdateInfoPeriod = utils.GetDuration(yc.UpdateInfoPeriod)
	}

	c.SenderConfig = readQueueablePayloadSenderConfigYaml(yc.QueueablePayloadSender)

	return c
}

func readTraceWriterConfigYaml(yc traceWriter) writerconfig.TraceWriterConfig {
	c := writerconfig.DefaultTraceWriterConfig()

	if yc.MaxSpansPerPayload > 0 {
		c.UpdateInfoPeriod = utils.GetDuration(yc.MaxSpansPerPayload)
	}
	if yc.FlushPeriod > 0 {
		c.FlushPeriod = utils.GetDuration(yc.FlushPeriod)
	}
	if yc.UpdateInfoPeriod > 0 {
		c.UpdateInfoPeriod = utils.GetDuration(yc.UpdateInfoPeriod)
	}

	c.SenderConfig = readQueueablePayloadSenderConfigYaml(yc.QueueablePayloadSender)

	return c
}

func readQueueablePayloadSenderConfigYaml(yc queueablePayloadSender) writerconfig.QueuablePayloadSenderConf {
	c := writerconfig.DefaultQueuablePayloadSenderConf()

	if yc.MaxAge > 0 {
		c.MaxAge = utils.GetDuration(yc.MaxAge)
	}

	if yc.MaxQueuedBytes > 0 {
		c.MaxQueuedBytes = yc.MaxQueuedBytes
	}

	if yc.MaxQueuedPayloads > 0 {
		c.MaxQueuedPayloads = yc.MaxQueuedPayloads
	}

	c.ExponentialBackoff = readExponentialBackoffConfigYaml(yc)

	return c
}

func readExponentialBackoffConfigYaml(yc queueablePayloadSender) backoff.ExponentialConfig {
	c := backoff.DefaultExponentialConfig()

	if yc.BackoffDuration > 0 {
		c.MaxDuration = utils.GetDuration(yc.BackoffDuration)
	}
	if yc.BackoffBase > 0 {
		c.Base = time.Duration(yc.BackoffBase) * time.Millisecond
	}
	if yc.BackoffGrowth > 0 {
		c.GrowthBase = yc.BackoffGrowth
	}

	return c
}
