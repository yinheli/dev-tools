package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yinheli/dev-tools/pkg/config"
	"github.com/yinheli/dev-tools/pkg/database"
	"github.com/yinheli/dev-tools/pkg/table"
	"github.com/yinheli/go-toolbox/logger/log"
	"go.uber.org/zap"
	"path"
)

func init() {
	var mysqlCfg *config.MySqlConfig

	tableCmd := &cobra.Command{
		Use:   "table",
		Short: "table functions",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			mysqlCfg, _ = config.GetMysqlConfig(path.Join(configDir, configMysqlFile))
			err := database.InitDB(mysqlCfg.CurrentAddr())
			if err != nil {
				log.Warn("init db error", zap.Error(err))
			}
		},
	}

	var (
		tableName string
	)
	goStructCmd := &cobra.Command{
		Use:   "gostruct -t table",
		Short: "table struct to go struct",
		Run: func(cmd *cobra.Command, args []string) {
			if tableName == "" {
				fmt.Println("table name should not empty")
				return
			}
			ret, err := table.ToGo(tableName)
			if err != nil {
				log.Warn("", zap.Error(err))
				return
			}
			fmt.Println()
			fmt.Println(ret)
			fmt.Println()
		},
	}
	goStructCmd.Flags().StringVarP(&tableName, "table", "t", "", "table")

	tableCmd.AddCommand(
		goStructCmd,
	)

	rootCmd.AddCommand(tableCmd)
}
