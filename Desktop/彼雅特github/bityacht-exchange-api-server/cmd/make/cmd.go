package make

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"bityacht-exchange-api-server/internal/pkg/logger"

	"github.com/spf13/cobra"
	"gorm.io/gorm/schema"
)

// ref: https://securego.io/docs/rules/g304.html
var errUnsafeFilename = errors.New("unsafe filename")

func init() {
	Cmd.AddCommand(makeMigrationCmd)
	Cmd.AddCommand(makeSeederCmd)
}

var Cmd = &cobra.Command{
	Use:   "make",
	Short: "Make can help you to generate code by template.",
}

const migrationPrefixCreate = "create_"

var migrationPrefixs = []string{"update", "drop"}

type migrationData struct {
	TimeStamp int64
	TableName string

	// For model.go
	Package string
}

func openOutputFile(rootPath string, filename string, mkdirAll bool) (string, *os.File, error) {
	filename = filepath.Clean(filename)
	if strings.ContainsAny(filename, "/\\") {
		return "", nil, errUnsafeFilename
	}

	fullFilePath := filepath.Clean(filepath.Join(rootPath, filename))
	if !strings.Contains(fullFilePath, rootPath) {
		return "", nil, errUnsafeFilename
	}

	if mkdirAll {
		if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
			return "", nil, err
		}
	}

	outputFile, err := os.Create(fullFilePath)
	if err != nil {
		return "", nil, err
	}

	return filename, outputFile, nil
}

var makeMigrationCmd = &cobra.Command{
	Use:   "migration {filename}",
	Short: "Create migration file into migrations.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpl, err := template.ParseFiles("cmd/make/templates/migration.go.tmpl")
		if err != nil {
			return err
		}
		data := migrationData{TimeStamp: time.Now().Unix(), TableName: strings.ToLower(args[0])}

		filename, outputFile, err := openOutputFile("cmd/migrate/migrations/", fmt.Sprintf("%d_%s.go", data.TimeStamp, data.TableName), false)
		if err != nil {
			return err
		}

		var createModel bool
		if strings.HasPrefix(data.TableName, migrationPrefixCreate) {
			data.TableName = strings.TrimPrefix(data.TableName, migrationPrefixCreate)
			createModel = true
		} else {
			for _, prefix := range migrationPrefixs {
				prefix += "_"

				if strings.HasPrefix(data.TableName, prefix) {
					data.TableName = strings.TrimPrefix(data.TableName, prefix)
					break
				}
			}
		}

		data.TableName = strings.TrimSuffix(data.TableName, "_table")

		if err := tmpl.Execute(outputFile, data); err != nil {
			return err
		} else if err := outputFile.Sync(); err != nil {
			return err
		} else if err = outputFile.Close(); err != nil {
			return err
		}

		logger.Logger.Info().Msg("migration file '" + filename + "' has been created!")

		if err := func() error {
			if !createModel {
				return nil
			}

			tmpl, err := template.ParseFiles("cmd/make/templates/model.go.tmpl")
			if err != nil {
				return err
			}

			data.Package = strings.ReplaceAll(data.TableName, "_", "")

			_, outputFile, err := openOutputFile("internal/database/sql/"+data.Package, "model.go", true)
			if err != nil {
				return err
			}

			if err := tmpl.Execute(outputFile, data); err != nil {
				return err
			} else if err := outputFile.Sync(); err != nil {
				return err
			} else if err = outputFile.Close(); err != nil {
				return err
			}

			logger.Logger.Info().Err(err).Msg("'" + data.Package + "/model.go' has been created!")
			return nil
		}(); err != nil {
			logger.Logger.Info().Err(err).Msg("create model.go failed!")
		}

		return nil
	},
}

type seederData struct {
	Name      string
	TableName string
}

var makeSeederCmd = &cobra.Command{
	Use:   "seeder",
	Short: "Create seeder file into migrations.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tmpl, err := template.ParseFiles("cmd/make/templates/seeder.go.tmpl")
		if err != nil {
			return err
		}
		data := seederData{Name: args[0], TableName: schema.NamingStrategy{}.TableName(args[0])}

		filename, outputFile, err := openOutputFile("cmd/seed/seeders/", fmt.Sprintf("%s.go", strings.ToLower(args[0])), false)
		if err != nil {
			return err
		}

		if err := tmpl.Execute(outputFile, data); err != nil {
			return err
		} else if err := outputFile.Sync(); err != nil {
			return err
		} else if err = outputFile.Close(); err != nil {
			return err
		}

		logger.Logger.Info().Msg("seeder file '" + filename + "' has been created! If the order of this seeder is important, make sure to set the order in cmd/seed/seeders/init.go")
		return nil
	},
}
