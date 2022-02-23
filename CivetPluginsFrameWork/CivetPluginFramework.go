package CivetPluginsFrameWork

import (
	"errors"
	"fmt"
	"plugin"
	"sync"
)

var once sync.Once
var PluginsIns *Plugins

type PlugConfig struct {
	Name     string
	Location string
}
type Plugins struct {
	pluginsMake    map[string]plugin.Plugin
	pluginsBind    map[string]string
	PluginsConfigs []*PlugConfig
}

func GetIns() *Plugins {
	once.Do(func() {
		PluginsIns = &Plugins{
			pluginsMake: make(map[string]plugin.Plugin),
			pluginsBind: make(map[string]string),
		}
	})
	return PluginsIns
}

func (p *Plugins) RegPlug(name string, location string) {
	p.PluginsConfigs = append(p.PluginsConfigs, &PlugConfig{Name: name, Location: location})
	plug, err := plugin.Open(location + name + ".so")
	if err != nil {
		panic(err)
	} else {
		p.pluginsMake[name] = *plug
	}
	fmt.Print(p.pluginsMake)
}
func (p *Plugins) BindPlug(IDName string, PluginName string) {
	if _, ok := p.pluginsMake[PluginName]; ok {
		p.pluginsBind[IDName] = PluginName
	} else {
		fmt.Println("不存在")
	}

}
func (p *Plugins) CallFunc(IDName string, FuncName string, arg1 string) (bool, error) {
	if _, ok := p.pluginsBind[IDName]; ok {
		p2 := p.pluginsMake[p.pluginsBind[IDName]]
		lookup, err := p2.Lookup(FuncName)
		if err != nil {
			panic(err)
		}
		lookup.(func(id string, df string))(FuncName, arg1)
		return true, nil
	} else {
		fmt.Println("产品", IDName, "不存在注册的插件")
		return false, errors.New(IDName + "不存在注册的插件")
	}
}
