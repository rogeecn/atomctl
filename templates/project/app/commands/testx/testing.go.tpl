package testx

import (
	"os"
	"testing"

	"go.ipao.vip/atom"
	"go.ipao.vip/atom/container"

	"github.com/rogeecn/fabfile"
	. "github.com/smartystreets/goconvey/convey"
)

func Default(providers ...container.ProviderContainer) container.Providers {
	return append(container.Providers{}, providers...)
}

func Serve(providers container.Providers, t *testing.T, invoke any) {
	Convey("tests boot up", t, func() {
		file := fabfile.MustFind("config.toml")

		localEnv := os.Getenv("ENV_LOCAL")
		if localEnv != "" {
			file = fabfile.MustFind("config." + localEnv + ".toml")
		}

		So(atom.LoadProviders(file, providers), ShouldBeNil)
		So(container.Container.Invoke(invoke), ShouldBeNil)
	})
}
