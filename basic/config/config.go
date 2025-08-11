package config

type Config struct {
	Mysql  Mysql      `mapstructure:"mysql" json:"mysql"`
	Server Server     `mapstructure:"server" json:"server"`
	NaCos  NaCos      `mapstructure:"nacos" json:"nacos"`
	Redis  Redis      `mapstructure:"redis" json:"redis"`
	MQTT   MQTTConfig `mapstructure:"mqtt" json:"mqtt"`
}

// MQTTConfig MQTT配置
type MQTTConfig struct {
	Broker   string `mapstructure:"broker"`
	Port     int    `mapstructure:"port"`
	ClientID string `mapstructure:"client_id"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

// ElasticsearchConfig Elasticsearch配置
type ElasticsearchConfig struct {
	Addrs []string `mapstructure:"addrs"`
}

type Mysql struct {
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	DBName   string `mapstructure:"dbname" json:"dbname"`
}

type Server struct {
	Port int `mapstructure:"port" json:"port"`
}

type Redis struct {
	Address  string `mapstructure:"address" json:"address"`
	Password string `mapstructure:"password" json:"password"`
	DB       int    `mapstructure:"db" json:"db"`
}
type NaCos struct {
	NamespaceId string `mapstructure:"namespace_id" json:"namespace_id"`
	IpAddr      string `mapstructure:"ip_addr" json:"ip_addr"`
	Port        uint64 `mapstructure:"port" json:"port"`
	DataId      string `mapstructure:"data_id" json:"data_id"`
	Group       string
}
