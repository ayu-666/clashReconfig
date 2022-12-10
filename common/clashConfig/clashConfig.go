package clashConfig

import (
	"gopkg.in/yaml.v3"
)

type Main struct {
	Secret         string                   `yaml:"secret"`
	AllowLan       bool                     `yaml:"allow-lan"`
	MixedPort      int                      `yaml:"mixed-port"`
	ProxyGroups    []*ProxyGroup            `yaml:"proxy-groups"`
	Proxies        []Proxy                  `yaml:"proxies"`
	ProxyProviders map[string]ProxyProvider `yaml:"proxy-providers"`
	RuleProviders  map[string]RuleProvider  `yaml:"rule-providers"`
	Rules          []Rule                   `yaml:"rules"`
}

type ProxyGroup struct {
	Interval  int      `yaml:"interval"`
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Proxies   []string `yaml:"proxies"`
	Url       string   `yaml:"url"`
	Use       []string `yaml:"use"`
	Lazy      bool     `yaml:"lazy"`
	Tolerance int      `yaml:"tolerance"`
}
type Proxy map[string]any
type HealthCheckOptions struct {
	Enable   bool   `yaml:"enable"`
	Url      string `yaml:"url"`
	Interval int    `yaml:"interval"`
}

type ProxyProvider struct {
	Type        string             `yaml:"type"`
	Url         string             `yaml:"url"`
	Interval    int                `yaml:"interval"`
	Path        string             `yaml:"path"`
	HealthCheck HealthCheckOptions `yaml:"health-check"`
}

type RuleProvider struct {
	Type     string `yaml:"type"`
	Url      string `yaml:"url"`
	Interval int    `yaml:"interval"`
	Path     string `yaml:"path"`
	Behavior string `yaml:"behavior"`
}
type Rule string

func AddProxyNamePrefixSuffix(_proxies []Proxy, prefix string, suffix string) []Proxy {
	_res := []Proxy{}
	for _, _proxy := range _proxies {
		_proxy["name"] = prefix + _proxy["name"].(string) + suffix
		_res = append(_res, _proxy)
	}
	return _res
}

func Yaml2Struct(content []byte) (Main, error) {
	mainData := Main{}
	err := yaml.Unmarshal(content, &mainData)
	return mainData, err
}
