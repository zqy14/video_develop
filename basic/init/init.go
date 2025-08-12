package init

import (
	"devicemanage/devicerpc/basic/config"
	"devicemanage/devicerpc/types"
	"devicemanage/devicerpc/utils"
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
	"log"
)

func init() {
	InitConfig()
	InitDB()
   utils.InitEs()

}

func InitDB() {
	config.BD = utils.InitMysqlS()
	config.Red = utils.InitRedis()
}

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("./dev.yaml")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Read config file failed, err: %v", err)
	}
	if err := v.Unmarshal(&config.GloBalConfig); err != nil {
		log.Fatalf("Unmarshal config file failed, err: %v", err)
	}
	//create clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         config.GloBalConfig.NaCos.NamespaceId, //we can create multiple clients with different namespaceId to support multiple namespace.When namespace is public, fill in the blank string here.
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      config.GloBalConfig.NaCos.IpAddr,
			ContextPath: "/nacos",
			Port:        config.GloBalConfig.NaCos.Port,
			Scheme:      "http",
		},
	}
	// Create config client for dynamic configuration
	configClient, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	content, _ := configClient.GetConfig(vo.ConfigParam{
		DataId: config.GloBalConfig.NaCos.DataId,
		Group:  config.GloBalConfig.NaCos.Group})

	var T types.T
	err := json.Unmarshal([]byte(content), &T)
	if err != nil {
		log.Fatalf("Unmarshal config file failed, err: %v", err)
	}

	err = configClient.ListenConfig(vo.ConfigParam{
		DataId: config.GloBalConfig.NaCos.DataId,
		Group:  config.GloBalConfig.NaCos.Group,
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})

}
