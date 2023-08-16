package metrics

import (
	"fmt"
	"net/http"

	"github.com/idle-ape/mixgo/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

// DefaultPath default path for prometheus http handler
var DefaultPath = "/metrics"

// DefaultEndpoint default endpoint for prometheus http handler
var DefaultEndpoint = ":9090"

var metrics = &Metrics{}

func init() {
	plugin.Register(metrics)
}

// Config config for metrics
type Config struct {
	EndPoint   string          `yaml:"endpoint"` // metrics server endpoint, default :9090
	Path       string          `yaml:"path"`     // path for prometheus to pull data, default /metrics
	Namespace  string          `yaml:"namespace"`
	Subsystem  string          `yaml:"subsystem"`
	Counters   []CounterConf   `yaml:"counters"`
	Histograms []HistogramConf `yaml:"histograms"`
}

type CounterConf struct {
	Opts   prometheus.CounterOpts `yaml:"opts"`
	Labels []string               `yaml:"labels"`
}

type HistogramConf struct {
	Opts   prometheus.HistogramOpts `yaml:"opts"`
	Labels []string                 `yaml:"labels"`
}

type Metrics struct {
	cfg        *Config
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
}

func (m *Metrics) Name() string {
	return "metrics"
}

func (m *Metrics) Setup(config yaml.Node) error {
	m.cfg = &Config{}
	if err := config.Decode(m.cfg); err != nil {
		return fmt.Errorf("%s setup err: %v", m.Name(), err)
	}

	m.counters = make(map[string]*prometheus.CounterVec)
	m.histograms = make(map[string]*prometheus.HistogramVec)

	// register counter
	m.registerCounter()

	// register histogram
	m.registerHistogram()

	go func() {
		endpoint, path := DefaultEndpoint, DefaultPath
		if m.cfg.EndPoint != "" {
			endpoint = m.cfg.EndPoint
		}
		if m.cfg.Path != "" {
			path = m.cfg.Path
		}
		http.Handle(path, promhttp.Handler())
		if err := http.ListenAndServe(endpoint, nil); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (m *Metrics) registerCounter() {
	for _, counter := range m.cfg.Counters {
		if counter.Opts.Namespace == "" {
			counter.Opts.Namespace = m.cfg.Namespace
		}
		if counter.Opts.Subsystem == "" {
			counter.Opts.Subsystem = m.cfg.Subsystem
		}
		m.counters[counter.Opts.Name] = prometheus.NewCounterVec(counter.Opts, counter.Labels)
		prometheus.MustRegister(m.counters[counter.Opts.Name])
	}
}

func (m *Metrics) registerHistogram() {
	for _, histogram := range m.cfg.Histograms {
		if histogram.Opts.Namespace == "" {
			histogram.Opts.Namespace = m.cfg.Namespace
		}
		if histogram.Opts.Subsystem == "" {
			histogram.Opts.Subsystem = m.cfg.Subsystem
		}
		m.histograms[histogram.Opts.Name] = prometheus.NewHistogramVec(histogram.Opts, histogram.Labels)
		prometheus.MustRegister(m.histograms[histogram.Opts.Name])
	}
}

// Counter get counter by name
func Counter(name string) *prometheus.CounterVec {
	return metrics.counters[name]
}

// Histogram get histogram by name
func Histogram(name string) *prometheus.HistogramVec {
	return metrics.histograms[name]
}
