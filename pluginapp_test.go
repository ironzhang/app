package pluginapp_test

import (
	"testing"

	"github.com/ironzhang/pluginapp"
	"github.com/ironzhang/pluginapp/plugins/backend"
	"github.com/ironzhang/x-pearls/log"
)

func TestMain(m *testing.M) {
	log.SetLevel("debug")
	m.Run()
}

func TestPlugins(t *testing.T) {
	var _ = backend.G

	args := []string{"test"}
	pluginapp.G.Main(args)
}
