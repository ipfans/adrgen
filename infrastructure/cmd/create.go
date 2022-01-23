package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/asiermarques/adrgen/application"
	"github.com/asiermarques/adrgen/domain"
	"github.com/asiermarques/adrgen/infrastructure"
	"github.com/spf13/cobra"
)

// SupersededADRId represents the ID for a superseded ADR File
//
var SupersededADRId int

// AmendedADRId represents the ID for a amended ADR File
//
var AmendedADRId int

// NewCreateCmd creates the 'create' CLI Command related to the ADR file creation
//
func NewCreateCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "create [the ADR title]",
		Short: "Create a new ADR File in the current directory",
		Long:  `Create a new ADR File in the current directory, you can add meta parameters for decisions tracing`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config, err := GetConfig()
			if err != nil {
				if config.TargetDirectory != "" {
					fmt.Printf("error creating file: %s", err)
				}

				fmt.Printf(
					"config file not found, working in the %s directory\n",
					config.TargetDirectory,
				)
			} else {
				fmt.Printf("config file found, working in the %s directory\n", config.TargetDirectory)
			}

			meta, metaError := cmd.LocalFlags().GetStringSlice("meta")
			if metaError != nil {
				fmt.Printf("an error occurred processing the meta parameter %s\n", metaError)
				return
			}
			for i, value := range meta {
				meta[i] = strings.TrimSpace(value)
			}
			config.MetaParams = append(config.MetaParams, meta...)

			supersedesADRId, supersedesError := cmd.LocalFlags().GetInt("supersedes")
			if supersedesError != nil {
				fmt.Printf(
					"an error occurred processing the supersedes parameter %s\n",
					supersedesError,
				)
				return
			}
			amendsADRId, amendsErr := cmd.LocalFlags().GetInt("amends")
			if amendsErr != nil {
				fmt.Printf("an error occurred processing the amends parameter %s\n", amendsErr)
				return
			}

			currentTime := time.Now()
			date := currentTime.Format("2006-01-02")

			adrWriter := infrastructure.CreateFileADRWriter(config.TargetDirectory)
			templateService := domain.CreateTemplateService(
				infrastructure.CreateCustomTemplateContentFileReader(config),
			)

			filename, creationError := application.CreateADRFile(
				date,
				args[0],
				meta,
				supersedesADRId,
				amendsADRId,
				config,
				infrastructure.CreateADRDirectoryRepository(config.TargetDirectory),
				adrWriter,
				templateService,
				domain.CreateRelationsManager(
					templateService,
					domain.CreateADRStatusManager(config),
				),
			)
			if creationError != nil {
				fmt.Println(creationError)
				return
			}
			fmt.Printf("file created: %s\n\n", filename)
		},
	}

	command.Flags().IntVarP(&SupersededADRId, "supersedes", "s", 0, "")
	command.Flags().IntVarP(&AmendedADRId, "amends", "a", 0, "")
	command.Flags().StringSliceVarP(&MetaFlag, "meta", "m", []string{}, "")

	command.Example = "adrgen create \"Using ADR to record and maintain decisions records\""

	return command
}

func init() {
	rootCmd.AddCommand(NewCreateCmd())
}
