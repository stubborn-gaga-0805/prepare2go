package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/stubborn-gaga-0805/prepare2go/configs/conf"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/helpers"
	"github.com/stubborn-gaga-0805/prepare2go/pkg/mysql"
	"os"
	"os/exec"
	"reflect"
)

const (
	defaultOutputPath  = "./internal/repo/orm"
	defaultPackageName = "orm"
	defaultDbConn      = "db"
)

type genModelCmd struct {
	*baseCmd
}

type genModelFlags struct {
	flagTables      string
	flagOutputPath  string
	flagPackageName string
	flagDBConn      string
}

var (
	flagTables      = flag{"table", "t", "", `Specify the generated table name (multiple tables are separated by ",")`}
	flagOutputPath  = flag{"output", "o", defaultOutputPath, `The path to execute the generated file, default "./internal/repo/orm"...`}
	flagPackageName = flag{"pkg", "p", defaultPackageName, `The package name of the generated model file, the default is "orm", which needs to correspond to the folder of the generated path...`}
	flagDBConn      = flag{"conn", "c", defaultDbConn, `The database connection configuration in the configuration file, the default "db"...`}

	dnsTpl = `%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local`
)

func newGenModelCmd() *genModelCmd {
	gc := &genModelCmd{newBaseCmd()}
	gc.cmd = &cobra.Command{
		Use:     "gen-model",
		Aliases: []string{},
		Short:   "Generate 'model' files for 'gorm'",
		Long:    `💡 Generate 'model' files for 'gorm', eg: aurora gen-model, enter interactive mode`,
		Run: func(cmd *cobra.Command, args []string) {
			gc.initGenModelRuntime(cmd)
			gc.initConfig()
			gc.run()
		},
	}
	addGenModelRuntimeFlag(gc.cmd, true)

	return gc
}

func (c *genModelCmd) runV2() {
	mysql.NewGenTool(helpers.GetContextWithRequestId(), conf.GetConfig().Data.Db, c.genModelFlags.flagTables, c.genModelFlags.flagOutputPath, c.genModelFlags.flagPackageName)
	return
}

func (c *genModelCmd) run() {
	var (
		cfg = conf.GetConfig().Data.Db
		cmd = "gentool"
		dns = fmt.Sprintf(dnsTpl, cfg.Username, cfg.Password, cfg.Addr, cfg.Database)
	)
	if c.genModelFlags.flagDBConn != defaultDbConn {
		dns = c.getConn()
		if len(dns) == 0 {
			fmt.Println("Failed to generate model file")
			return
		}
	}

	// 组装gen-tools命令
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	command := exec.Command(cmd,
		"-dsn", dns,
		"-db", "mysql",
		"-tables", c.genModelFlags.flagTables,
		"-modelPkgName", c.genModelFlags.flagPackageName,
		"-outPath", c.genModelFlags.flagOutputPath,
		"-outFile", "gen.go",
		"-onlyModel",
		"-fieldWithIndexTag",
		"-fieldWithTypeTag",
		"-fieldNullable",
	)
	command.Dir = wd
	command.Env = os.Environ()

	stdout, _ := command.StdoutPipe()
	stderr, _ := command.StderrPipe()
	if err := command.Start(); err != nil {
		fmt.Println(err)
	}
	// 读取命令的标准输出和错误输出
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	scanner = bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := command.Wait(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\nFinished...!")
	return
}

// 通过命令注入运行环境参数
func addGenModelRuntimeFlag(cmd *cobra.Command, persistent bool) {
	getFlags(cmd, persistent).StringP(flagTables.name, flagTables.shortName, flagTables.defaultValue.(string), flagTables.usage)
	getFlags(cmd, persistent).StringP(flagOutputPath.name, flagOutputPath.shortName, flagOutputPath.defaultValue.(string), flagOutputPath.usage)
	getFlags(cmd, persistent).StringP(flagPackageName.name, flagPackageName.shortName, flagPackageName.defaultValue.(string), flagPackageName.usage)
	getFlags(cmd, persistent).StringP(flagDBConn.name, flagDBConn.shortName, flagDBConn.defaultValue.(string), flagDBConn.usage)
}

func getTables(cmd *cobra.Command) string {
	return cmd.Flag(flagTables.name).Value.String()
}

func getOutputPath(cmd *cobra.Command) string {
	return cmd.Flag(flagOutputPath.name).Value.String()
}

func getPackageName(cmd *cobra.Command) string {
	return cmd.Flag(flagPackageName.name).Value.String()
}

func getDbConn(cmd *cobra.Command) string {
	return cmd.Flag(flagDBConn.name).Value.String()
}

func (c *genModelCmd) getConn() string {
	t := reflect.TypeOf(conf.GetConfig().Data)
	v := reflect.ValueOf(conf.GetConfig().Data)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("json") == c.genModelFlags.flagDBConn && field.Type == reflect.TypeOf(mysql.DB{}) {
			cfg := v.Field(i).Interface().(mysql.DB)
			return fmt.Sprintf(dnsTpl, cfg.Username, cfg.Password, cfg.Addr, cfg.Database)
		}
	}
	fmt.Println(fmt.Sprintf("Database connection configuration error! Connection [%s] does not exist!!！", c.genModelFlags.flagDBConn))
	return ""
}
