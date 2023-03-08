package apollogf

import (
	"context"
	"fmt"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/agcache"
	"github.com/apolloconfig/agollo/v4/constant"
	apolloconfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/extension"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

// default namespace for g.Cfg()
var defaultNamespace string

type ApolloAdapter struct {
	Client agollo.Client
	Cache  agcache.CacheInterface
	Config *apolloconfig.AppConfig
}

type ApolloArg struct {
	Namespaces []string
	AppId      string
	IP         string
	Cluster    string
}

func SetDefaultNamespace(n string) {
	defaultNamespace = n
}

func SetAdapters(arg ApolloArg) {
	for i, namespace := range arg.Namespaces {
		cnf := &apolloconfig.AppConfig{
			AppID:          arg.AppId,
			Cluster:        arg.Cluster,
			IP:             arg.IP,
			NamespaceName:  namespace,
			IsBackupConfig: true,
		}
		adapter := NewAdapter(cnf)
		if defaultNamespace == "" {
			if i == 0 {
				g.Cfg().SetAdapter(adapter)
			}
		} else {
			if namespace == defaultNamespace {
				g.Cfg().SetAdapter(adapter)
			}
		}
		g.Cfg(namespace).SetAdapter(adapter)
	}
}

func NewAdapter(c *apolloconfig.AppConfig) *ApolloAdapter {
	extension.AddFormatParser(constant.YAML, &Parser{})
	extension.AddFormatParser(constant.YML, &Parser{})

	client, err := agollo.StartWithConfig(func() (*apolloconfig.AppConfig, error) {
		return c, nil
	})

	ctx := context.Background()
	if err != nil {
		g.Log().Panic(ctx, err)
	}
	cache := client.GetConfigCache(c.NamespaceName)
	return &ApolloAdapter{
		Client: client,
		Config: c,
		Cache:  cache,
	}
}

// Available checks and returns the backend configuration service is available.
// The optional parameter `resource` specifies certain configuration resource.
//
// Note that this function does not return error as it just does simply check for
// backend configuration service.
func (a *ApolloAdapter) Available(ctx context.Context, resource string) (ok bool) {
	return true
}

// Get retrieves and returns value by specified `pattern` in current resource.
// Pattern like:
// "x.y.z" for map item.
// "x.0.y" for slice item.
func (a *ApolloAdapter) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	m, err := a.Data(ctx)
	if err != nil {
		return nil, err
	}
	v, err := gjson.LoadJson(m)
	if err != nil {
		return nil, err
	}
	return v.Get(pattern).Val(), nil
}

// Data retrieves and returns all configuration data in current resource as map.
// Note that this function may lead lots of memory usage if configuration data is too large,
// you can implement this function if necessary.
func (a *ApolloAdapter) Data(ctx context.Context) (data map[string]interface{}, err error) {
	m := make(map[string]interface{})
	a.Cache.Range(func(key, value interface{}) bool {
		m[fmt.Sprintf("%s", key)] = value
		return true
	})
	return m, nil
}

type Parser struct {
}

// Parse 内存内容=>yml文件转换器
func (d *Parser) Parse(configContent interface{}) (map[string]interface{}, error) {
	configJson, err := gjson.LoadContentType("yaml", configContent, true)
	if err != nil {
		return nil, err
	}
	if configJson != nil {
		return configJson.Var().Map(), nil
	}
	return nil, nil
}
