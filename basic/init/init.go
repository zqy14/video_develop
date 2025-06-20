package init

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
	"public_comment/basic/config"
	"public_comment/basic/utils"
	"public_comment/consul"
)

func init() {
	InitNaCos()
	InitConsul()
	InitConfig()
	InitDB()
	InitZap()
}

func InitZap() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"user.log"}
	logger, err := zapConfig.Build()
	if err != nil {
		log.Println("初始化日志失败", err)
		return
	}
	zap.ReplaceGlobals(logger)
	zap.S().Infof("zap初始化日志成功")
}


func InitNaCos() {
	v := viper.New()
	v.SetConfigFile("./dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		return
	}
	NaCosConfig := config.NaCos{}
	if err := v.Unmarshal(&NaCosConfig); err != nil {
		return
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         NaCosConfig.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// At least one ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      NaCosConfig.IpAddr,
			ContextPath: "/nacos",
			Port:        NaCosConfig.Port,
			Scheme:      "http",
		},
	}

	// Create config client for dynamic configuration
	configClient, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	content, _ := configClient.GetConfig(vo.ConfigParam{
		DataId: NaCosConfig.DataId,
		Group:  NaCosConfig.Group})

	json.Unmarshal([]byte(content), &config.GlobalNaCos)

}
func InitConsul() {
	config, err := consul.ConsulConfig("14.103.243.153", 8500)
	if err != nil {
		return
	}
	register := consul.ConsulRegister{
		Id:      "os",
		Name:    "os",
		Tags:    []string{"os"},
		Port:    8503,
		Address: "localhost",
	}
	err = config.ConsulRegister(register)
	if err != nil {
		return  
	} else {
		log.Println("consul init success")
	}
}

func InitDB() {
	config.DB = utils.GlobalMysql()
	config.Red = utils.ExampleClient()
}

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("D:\\GoWork\\src\\videos\\user-rpc\\dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		return
	}
	if err := v.Unmarshal(&config.Global); err != nil {
		return
	}
	log.Println("viper init success")
}
