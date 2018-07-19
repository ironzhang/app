package main

import (
	"fmt"
	"os"

	"github.com/ironzhang/pluginapp"
)

type Config struct {
	Environment string
	Node        string
	Addr        string
}

type Plugin struct {
	Config Config
}

func (p *Plugin) Name() string {
	return "main"
}

func (p *Plugin) Init() error {
	fmt.Println(p.Config.Environment, p.Config.Node, p.Config.Addr)
	return nil
}

func (p *Plugin) Fini() error {
	return nil
}

var G = &Plugin{
	Config: Config{
		Environment: "production",
		Node:        "node1",
		Addr:        "localhost:8000",
	},
}

func init() {
	pluginapp.G.Register(G, &G.Config)
}

func main() {
	pluginapp.G.Main(os.Args)
}
