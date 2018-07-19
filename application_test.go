package pluginapp

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

type TConfig struct {
	Environment string
	Node        string
	Addr        string
}

type TPlugin struct {
}

func (p *TPlugin) Name() string {
	return "TPlugin"
}

func (p *TPlugin) Init() error {
	return nil
}

func (p *TPlugin) Fini() error {
	return nil
}

func TestApplicationVersion(t *testing.T) {
	app := Application{
		CommandLine: flag.NewFlagSet("test", flag.ExitOnError),
		VersionInfo: func() string { return fmt.Sprintf("version: v0.0.1\n") },
	}
	args := []string{"test", "-version"}
	app.Main(args)
}

func TestApplicationConfigExample(t *testing.T) {
	app := Application{
		CommandLine: flag.NewFlagSet("test", flag.ExitOnError),
		VersionInfo: func() string { return fmt.Sprintf("version: v0.0.2\n") },
	}
	app.Register(&TPlugin{}, &TConfig{Environment: "production", Node: "node0", Addr: "localhost:8000"})

	args := []string{"test", "-config-example", "test.conf"}
	app.Main(args)
	os.Remove("test.conf")
}

func TestMain(m *testing.M) {
	//config.Default = tomlcfg.TOML
	exit = func(code int) { fmt.Printf("exit(%d)\n", code) }

	m.Run()
}
