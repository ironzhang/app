package pluginapp

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ironzhang/pluginapp/jsonconfig"
	"github.com/ironzhang/x-pearls/log"
)

var exit = os.Exit

var G Application

type Plugin interface {
	Name() string
	Init() error
	Fini() error
}

type Flagger interface {
	SetFlags(fs *FlagSet)
}

type Runner interface {
	Run(ctx context.Context) error
}

type Options struct {
	Version       bool
	ConfigFile    string
	ConfigExample string
}

type Application struct {
	CommandLine *flag.FlagSet
	Options     Options
	VersionInfo func() string
	LoadConfig  func(filename string, configs map[string]interface{}) error
	WriteConfig func(filename string, configs map[string]interface{}) error

	configs map[string]interface{}
	plugins []Plugin
}

func (app *Application) Register(p Plugin, c interface{}) {
	for _, v := range app.plugins {
		if v.Name() == p.Name() {
			panic(fmt.Sprintf("%q plugin is registered", p.Name()))
		}
	}
	app.plugins = append(app.plugins, p)

	if app.configs == nil {
		app.configs = make(map[string]interface{})
	}
	if c != nil {
		app.configs[p.Name()] = c
	}
}

func (app *Application) Main(args []string) {
	var err error

	// 解析命令行参数
	if err = app.parseCommandLine(args); err != nil {
		fmt.Fprintf(os.Stderr, "parse command line: %v\n", err)
		exit(3)
	}

	// 执行命令
	if err = app.doCommand(); err != nil {
		fmt.Fprintf(os.Stderr, "do command: %v\n", err)
		exit(3)
	}

	// 加载配置
	if err = app.loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		exit(3)
	}

	// 输出程序基本信息
	app.printInfo()

	// 运行插件
	if err = app.runPlugins(); err != nil {
		fmt.Fprintf(os.Stderr, "run plugins: %v\n", err)
		exit(3)
	}
}

func (app *Application) setupCommandLine() {
	if app.VersionInfo != nil {
		app.CommandLine.BoolVar(&app.Options.Version, "version", app.Options.Version, "输出版本信息")
	}
	if len(app.configs) > 0 {
		app.CommandLine.StringVar(&app.Options.ConfigFile, "config-file", app.Options.ConfigFile, "指定配置文件")
		app.CommandLine.StringVar(&app.Options.ConfigExample, "config-example", app.Options.ConfigExample, "生成配置示例")
	}

	for _, p := range app.plugins {
		if f, ok := p.(Flagger); ok {
			f.SetFlags(&FlagSet{
				prefix:  p.Name() + ".",
				flagset: app.CommandLine,
			})
		}
	}
}

func (app *Application) parseCommandLine(args []string) error {
	if app.CommandLine == nil {
		app.CommandLine = flag.CommandLine
	}
	app.setupCommandLine()
	return app.CommandLine.Parse(args[1:])
}

func (app *Application) doCommand() (err error) {
	var quit bool

	if app.Options.Version {
		if app.VersionInfo == nil {
			return errors.New("Application.VersionInfo is nil")
		}
		fmt.Fprintf(os.Stdout, "%s", app.VersionInfo())
		quit = true
	} else if app.Options.ConfigExample != "" {
		if app.WriteConfig == nil {
			app.WriteConfig = jsonconfig.Write
		}
		if err = app.WriteConfig(app.Options.ConfigExample, app.configs); err != nil {
			return fmt.Errorf("generate config example: %v", err)
		}
		fmt.Fprintf(os.Stdout, "generate config example %s success\n", app.Options.ConfigExample)
		quit = true
	}

	if quit {
		exit(0)
	}

	return nil
}

func (app *Application) loadConfig() (err error) {
	if app.Options.ConfigFile == "" {
		return nil
	}
	if len(app.configs) <= 0 {
		return nil
	}
	if app.LoadConfig == nil {
		app.LoadConfig = jsonconfig.Load
	}
	return app.LoadConfig(app.Options.ConfigFile, app.configs)
}

func (app *Application) configInfo() string {
	data, _ := json.MarshalIndent(app.configs, "", "\t")
	return string(data)
}

func (app *Application) printInfo() {
	if app.VersionInfo != nil {
		log.Infof("version info:\n%s", app.VersionInfo())
	}
	if len(app.configs) > 0 {
		log.Infof("config info:\n%s", app.configInfo())
	}
}

func (app *Application) runPlugins() (err error) {
	// init plugins
	for _, p := range app.plugins {
		if err = p.Init(); err != nil {
			log.Errorw("init", "plugin", p.Name(), "error", err)
			return fmt.Errorf("init %s plugin: %v", p.Name(), err)
		}
		log.Debugw("init success", "plugin", p.Name())
	}

	// quit signal
	quit := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		cancel()
		time.Sleep(10 * time.Second)
		quit <- errors.New("wait 10s, force exit")
	}()

	// run plugins
	var wg sync.WaitGroup
	for _, p := range app.plugins {
		if r, ok := p.(Runner); ok {
			wg.Add(1)
			go func(name string, runner Runner) {
				if err := runner.Run(ctx); err != nil {
					quit <- fmt.Errorf("run %s plugin: %v", name, err)
				}
				wg.Done()
			}(p.Name(), r)
		}
	}
	go func() {
		wg.Wait()
		quit <- nil
	}()

	// fini plugins
	err = <-quit
	for i := len(app.plugins) - 1; i >= 0; i-- {
		p := app.plugins[i]
		if err = p.Fini(); err != nil {
			log.Errorw("fini", "plugin", p.Name(), "error", err)
			continue
		}
		log.Debugw("fini success", "plugin", p.Name())
	}
	return err
}
