package CivetPluginsFrameWork

import (
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
	plug, err := plugin.Open(name + ".so")
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
func (p *Plugins) CallFunc(IDName string, FuncName string) {
	p2 := p.pluginsMake[p.pluginsBind[IDName]]
	lookup, err := p2.Lookup("Analyst")
	if err != nil {
		panic(err)
	}
	lookup.(func(id string, df string))(FuncName, FuncName)
}
