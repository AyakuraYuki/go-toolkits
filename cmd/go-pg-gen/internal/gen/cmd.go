package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const BinName = "go-pg-gen"

var (
	Cmd *cobra.Command
	Log zerolog.Logger

	host          string // database server host
	port          uint32 // database server port
	username      string // database user name
	password      string // database user password
	database      string // database name to connect to
	table         string // table name to connect to and generate go struct
	dst           string // output code path
	jsonCamelCase bool   // use camel case for json names
	verbose       bool   // verbose output

	dsn      string // Gorm data source name
	logLevel = zerolog.InfoLevel
)

func init() {
	Cmd = &cobra.Command{
		Use:              fmt.Sprintf("%s [table_name_arg]", BinName),
		Short:            "Generate PostgreSQL table into Go struct.",
		Args:             cobra.MaximumNArgs(1),
		PersistentPreRun: persistentPreRun,
		PreRunE:          preRunE,
		RunE:             runE,
	}

	Cmd.Flags().SortFlags = false
	Cmd.PersistentFlags().SortFlags = false

	Cmd.Flags().StringVarP(&host, "host", "H", "127.0.0.1", "database server host")
	Cmd.Flags().Uint32VarP(&port, "port", "p", 5432, "database server port")
	Cmd.Flags().StringVarP(&username, "username", "U", "postgres", "database user name")
	Cmd.Flags().StringVarP(&password, "password", "W", "", "database user password")
	Cmd.Flags().StringVarP(&database, "dbname", "d", "postgres", "database name to connect to")
	Cmd.Flags().StringVarP(&table, "table", "t", "", "table name to connect to and generate go struct")
	Cmd.Flags().StringVarP(&dst, "out", "o", ".", "output code path")
	Cmd.Flags().BoolVar(&jsonCamelCase, "json-camel-case", false, "use camel case for json names")
	Cmd.Flags().BoolVar(&verbose, "verbose", false, "verbose output")
}

func persistentPreRun(_ *cobra.Command, _ []string) {
	if verbose {
		logLevel = zerolog.DebugLevel
	}
}

func preRunE(_ *cobra.Command, args []string) error {
	initLogger()

	if dst != "" {
		dst, _ = filepath.Abs(dst)
	}
	if dst == "" {
		dst, _ = os.Getwd()
	}
	if dst == "" {
		dst, _ = os.UserHomeDir()
	}
	if dst == "" {
		dst = os.TempDir()
	}
	dst = filepath.Join(dst, "gen")
	Log.Info().Msgf("generated code save to %q", dst)

	if password != "" {
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, database)
	} else {
		dsn = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, username, database)
	}

	if table == "" && len(args) > 0 {
		table = strings.TrimSpace(args[0])
	}
	if table == "" {
		return fmt.Errorf("no tablename specified")
	}

	if verbose {
		Log.Debug().
			Str("host", host).
			Uint32("port", port).
			Str("username", username).
			Str("database", database).
			Str("table", table).
			Str("dst", dst).
			Msg("parameters")
	}

	return nil
}

func initLogger() {
	output := zerolog.ConsoleWriter{
		Out:              os.Stdout,
		TimeFormat:       time.DateTime,
		TimeLocation:     time.Now().Location(),
		FormatLevel:      func(i any) string { return strings.ToUpper(fmt.Sprintf("| %-6s |", i)) },
		FormatFieldName:  func(i any) string { return fmt.Sprintf("[%s: ", i) },
		FormatFieldValue: func(i any) string { return fmt.Sprintf("%v]", i) },
	}
	Log = zerolog.New(output).Level(logLevel)
}
