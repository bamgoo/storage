package storage

import (
	"sync"

	"github.com/bamgoo/bamgoo"
	. "github.com/bamgoo/base"
	"github.com/bamgoo/util"
)

func init() {
	bamgoo.Mount(module)
}

var module = &Module{
	filecfg: fileConfig{
		Download:  "store/download",
		Thumbnail: "store/thumbnail",
		Preview:   "store/preview",
		Salt:      bamgoo.BAMGOO,
	},
	configs:   make(Configs, 0),
	drivers:   make(map[string]Driver, 0),
	instances: make(map[string]*Instance, 0),
}

type (
	Module struct {
		mutex sync.Mutex

		initialized bool
		connected   bool
		started     bool

		filecfg fileConfig
		configs Configs
		drivers map[string]Driver

		instances map[string]*Instance
		weights   map[string]int
		hashring  *util.HashRing
	}

	fileConfig struct {
		Download  string
		Thumbnail string
		Preview   string
		Salt      string
	}

	Configs map[string]Config
	Config  struct {
		Driver  string
		Weight  int
		Prefix  string
		Proxy   bool
		Remote  bool
		Setting Map
	}
)

func (m *Module) Register(name string, value Any) {
	switch v := value.(type) {
	case Driver:
		m.RegisterDriver(name, v)
	case Config:
		m.RegisterConfig(name, v)
	case Configs:
		m.RegisterConfigs(v)
	}
}

func (m *Module) RegisterDriver(name string, driver Driver) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if name == "" {
		name = bamgoo.DEFAULT
	}
	if driver == nil {
		panic("Invalid storage driver: " + name)
	}
	if bamgoo.Override() {
		m.drivers[name] = driver
	} else {
		if _, ok := m.drivers[name]; !ok {
			m.drivers[name] = driver
		}
	}
}

func (m *Module) RegisterConfig(name string, cfg Config) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if name == "" {
		name = bamgoo.DEFAULT
	}
	if bamgoo.Override() {
		m.configs[name] = cfg
	} else {
		if _, ok := m.configs[name]; !ok {
			m.configs[name] = cfg
		}
	}
}

func (m *Module) RegisterConfigs(configs Configs) {
	for k, v := range configs {
		m.RegisterConfig(k, v)
	}
}

func (m *Module) Config(global Map) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if cfg, ok := global["file"].(Map); ok {
		if v, ok := cfg["download"].(string); ok {
			m.filecfg.Download = v
		}
		if v, ok := cfg["thumbnail"].(string); ok {
			m.filecfg.Thumbnail = v
		}
		if v, ok := cfg["thumb"].(string); ok {
			m.filecfg.Thumbnail = v
		}
		if v, ok := cfg["preview"].(string); ok {
			m.filecfg.Preview = v
		}
		if v, ok := cfg["salt"].(string); ok {
			m.filecfg.Salt = v
		}
	}

	cfgAny, ok := global["storage"]
	if !ok {
		return
	}
	cfg, ok := cfgAny.(Map)
	if !ok || cfg == nil {
		return
	}

	root := Map{}
	for key, val := range cfg {
		if item, ok := val.(Map); ok && key != "setting" {
			m.configure(key, item)
		} else {
			root[key] = val
		}
	}
	if len(root) > 0 {
		m.configure(bamgoo.DEFAULT, root)
	}
}

func (m *Module) configure(name string, cfg Map) {
	out := Config{Driver: bamgoo.DEFAULT, Weight: 1}
	if vv, ok := m.configs[name]; ok {
		out = vv
	}

	if v, ok := cfg["driver"].(string); ok && v != "" {
		out.Driver = v
	}
	if v, ok := cfg["weight"].(int); ok {
		out.Weight = v
	}
	if v, ok := cfg["weight"].(int64); ok {
		out.Weight = int(v)
	}
	if v, ok := cfg["weight"].(float64); ok {
		out.Weight = int(v)
	}
	if v, ok := cfg["prefix"].(string); ok {
		out.Prefix = v
	}
	if v, ok := cfg["proxy"].(bool); ok {
		out.Proxy = v
	}
	if v, ok := cfg["remote"].(bool); ok {
		out.Remote = v
	}
	if v, ok := cfg["setting"].(Map); ok {
		out.Setting = v
	}

	m.configs[name] = out
}

func (m *Module) Setup() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.initialized {
		return
	}

	if len(m.configs) == 0 {
		m.configs[bamgoo.DEFAULT] = Config{Driver: bamgoo.DEFAULT, Weight: 1}
	}
	for k, v := range m.configs {
		if v.Driver == "" {
			v.Driver = bamgoo.DEFAULT
		}
		if v.Weight == 0 {
			v.Weight = 1
		}
		m.configs[k] = v
	}
	m.initialized = true
}

func (m *Module) Open() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.connected {
		return
	}

	weights := make(map[string]int, 0)
	for name, cfg := range m.configs {
		drv := m.drivers[cfg.Driver]
		if drv == nil {
			panic("Invalid storage driver: " + cfg.Driver)
		}
		inst := &Instance{Name: name, Config: cfg, Setting: cfg.Setting}
		conn, err := drv.Connect(inst)
		if err != nil {
			panic("Failed to connect storage: " + err.Error())
		}
		if err := conn.Open(); err != nil {
			panic("Failed to open storage: " + err.Error())
		}
		inst.conn = conn
		m.instances[name] = inst
		if cfg.Weight > 0 {
			weights[name] = cfg.Weight
		}
	}
	m.weights = weights
	m.hashring = util.NewHashRing(weights)
	m.connected = true
}

func (m *Module) Start() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.started {
		return
	}
	m.started = true
}

func (m *Module) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if !m.started {
		return
	}
	m.started = false
}

func (m *Module) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for _, inst := range m.instances {
		if inst.conn != nil {
			_ = inst.conn.Close()
		}
	}
	m.instances = make(map[string]*Instance, 0)
	m.hashring = nil
	m.connected = false
	m.initialized = false
}
