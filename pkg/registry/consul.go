package registry

import (
	"log"
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
)

// ConsulClient 是 Consul 客户端的封装
type ConsulClient struct {
	client *api.Client
}

// NewConsulClient 创建一个新的 Consul 客户端
func NewConsulClient() (*ConsulClient, error) {
	// 创建 Consul 客户端配置
	config := api.DefaultConfig()

	// 创建 Consul 客户端
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulClient{client: client}, nil
}

// RegisterService 将服务注册到 Consul
func (c *ConsulClient) RegisterService(serviceName string, port int) error {
	// 获取本地 IP 地址
	ip, err := getLocalIP()
	if err != nil {
		return err
	}

	// 创建服务注册信息
	registration := &api.AgentServiceRegistration{
		ID:      serviceName, // 服务唯一标识
		Name:    serviceName, // 服务名称
		Port:    port,        // 服务端口
		Address: ip,          // 服务地址
		Check: &api.AgentServiceCheck{
			HTTP:     "http://" + ip + ":" + strconv.Itoa(port) + "/health", // 健康检查地址
			Interval: "10s",                                               // 健康检查间隔
			Timeout:  "1s",                                                // 健康检查超时时间
		},
	}

	// 注册服务
	err = c.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}

	log.Printf("Service %s registered successfully with IP %s and port %d", serviceName, ip, port)
	return nil
}

// DiscoverService 从 Consul 中发现服务
func (c *ConsulClient) DiscoverService(serviceName string) (string, error) {
	// 查询服务
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", err
	}

	// 如果没有找到服务，返回空
	if len(services) == 0 {
		return "", nil
	}

	// 返回第一个服务的地址
	service := services[0] // 从切片中获取第一个服务
	return service.Service.Address + ":" + strconv.Itoa(service.Service.Port), nil
}

// getLocalIP 获取本地 IP 地址
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}