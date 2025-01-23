package atom

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/dig"
	"{{.ModuleName}}/pkg/atom/config"
	"{{.ModuleName}}/pkg/atom/container"
)

var (
	GroupInitialName           = "initials"
	GroupRoutesName            = "routes"
	GroupGrpcServerServiceName = "grpc_server_services"
	GroupCommandName           = "command_services"
	GroupQueueName             = "queue_handler"
	GroupCronJobName           = "cron_jobs"

	GroupInitial    = dig.Group(GroupInitialName)
	GroupRoutes     = dig.Group(GroupRoutesName)
	GroupGrpcServer = dig.Group(GroupGrpcServerServiceName)
	GroupCommand    = dig.Group(GroupCommandName)
	GroupQueue      = dig.Group(GroupQueueName)
	GroupCronJob    = dig.Group(GroupCronJobName)
)

func Serve(opts ...Option) error {
	rootCmd := &cobra.Command{Use: "app"}
	for _, opt := range opts {
		opt(rootCmd)
	}

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return err
	})

	rootCmd.PersistentFlags().StringP("config", "c", "config.toml", "config file")

	return rootCmd.Execute()
}

func LoadProviders(configFile string, providers container.Providers) error {
	configure, err := config.Load(configFile)
	if err != nil {
		return errors.Wrapf(err, "load config file: %s", configFile)
	}

	if err := providers.Provide(configure); err != nil {
		return err
	}
	return nil
}

type Option func(*cobra.Command)

var (
	AppName    string
	AppVersion string
)

func Providers(providers container.Providers) Option {
	return func(cmd *cobra.Command) {
		cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			return LoadProviders(cmd.Flag("config").Value.String(), providers)
		}
	}
}

func Command(opt ...Option) Option {
	return func(parentCmd *cobra.Command) {
		cmd := &cobra.Command{}
		for _, o := range opt {
			o(cmd)
		}
		parentCmd.AddCommand(cmd)
	}
}

func Arguments(f func(cmd *cobra.Command)) Option {
	return f
}

func Version(ver string) Option {
	return func(cmd *cobra.Command) {
		cmd.Version = ver
		AppVersion = ver
	}
}

func Name(name string) Option {
	return func(cmd *cobra.Command) {
		cmd.Use = name
		AppName = name
	}
}

func Short(short string) Option {
	return func(cmd *cobra.Command) {
		cmd.Short = short
	}
}

func Long(long string) Option {
	return func(cmd *cobra.Command) {
		cmd.Long = long
	}
}

func Example(example string) Option {
	return func(cmd *cobra.Command) {
		cmd.Example = example
	}
}

func Run(run func(cmd *cobra.Command, args []string)) Option {
	return func(cmd *cobra.Command) {
		cmd.Run = run
	}
}

func RunE(run func(cmd *cobra.Command, args []string) error) Option {
	return func(cmd *cobra.Command) {
		cmd.RunE = run
	}
}

func PostRun(run func(cmd *cobra.Command, args []string)) Option {
	return func(cmd *cobra.Command) {
		cmd.PostRun = run
	}
}

func PostRunE(run func(cmd *cobra.Command, args []string) error) Option {
	return func(cmd *cobra.Command) {
		cmd.PostRunE = run
	}
}

func PreRun(run func(cmd *cobra.Command, args []string)) Option {
	return func(cmd *cobra.Command) {
		cmd.PreRun = run
	}
}

func PreRunE(run func(cmd *cobra.Command, args []string) error) Option {
	return func(cmd *cobra.Command) {
		cmd.PreRunE = run
	}
}

func Config(file string) Option {
	return func(cmd *cobra.Command) {
		_ = cmd.PersistentFlags().Set("config", file)
	}
}
