package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yinheli/dev-tools/pkg/config"
	"os"
	"path"
)

var (
	configMysqlFile = "mysql.yaml"

	configCmd = &cobra.Command{
		Use:     "config [mysql]",
		Short:   "config",
		Aliases: []string{"c"},
	}

	configMysqlCmd = &cobra.Command{
		Use: "mysql [list|set|remove]",
	}
)

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configMysqlCmd)

	configMysqlCmd.AddCommand(
		&cobra.Command{
			Use:     "list",
			Aliases: []string{"l"},
			Run: func(cmd *cobra.Command, args []string) {
				if cfg, ok := getMysqlConfig(); ok {
					cfg.Print()
				}
			},
		},

		&cobra.Command{
			Use:     "set",
			Aliases: []string{"s"},
			Args:    cobra.MinimumNArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				cfg, _ := config.GetMysqlConfig(path.Join(configDir, configMysqlFile))
				cfg.Set(args[0], args[1])
			},
		},

		&cobra.Command{
			Use:     "remove",
			Aliases: []string{"d", "r"},
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				if cfg, ok := getMysqlConfig(); ok {
					cfg.Remove(args[0])
				}
			},
		},

		&cobra.Command{
			Use:     "use",
			Aliases: []string{"u"},
			Args:    cobra.MinimumNArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				if cfg, ok := getMysqlConfig(); ok {
					cfg.Use(args[0])
				}
			},
		},
	)
}

func getMysqlConfig() (*config.MySqlConfig, bool) {
	cfg, err := config.GetMysqlConfig(path.Join(configDir, configMysqlFile))
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no config found")
			return nil, false
		}

		fmt.Println("read config file faild", err)
		return nil, false
	}

	return cfg, true
}
