package proxy

type Host struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Service struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
	Hosts  []Host `yaml:"hosts"`
}

type Proxy struct {
	Listen     Host      `yaml:"listen"`
	Services   []Service `yaml:"services"`
	ServiceMap map[string]*Service
}

type Config struct {
	Proxy Proxy `yaml:"proxy"`
}
