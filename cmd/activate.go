/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

type InstructionFile struct {
	FileName         string
	Name             string
	Position         int
	EnabledByDefault bool
}

func extractContent(fileContent string) string {
	lines := strings.Split(fileContent, "\n")
	dashCount := 0
	startIdx := 0

	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			dashCount++
			if dashCount == 2 {
				startIdx = i + 1
				break
			}
		}
	}

	return strings.Join(lines[startIdx:], "\n")
}

func getInstructionFiles() ([]InstructionFile, error) {
	f, err := os.Open("./.agentscope")
	if err != nil {
		return nil, err
	}
	fileNames, _ := f.Readdirnames(0)

	var instructions []InstructionFile
	for _, fileName := range fileNames {
		content, _ := Read(filepath.Join(".agentscope", fileName))
		opts := Parse(content)
		pos, _ := strconv.Atoi(opts["position"])
		instructions = append(instructions, InstructionFile{
			FileName:         fileName,
			Name:             opts["name"],
			Position:         pos,
			EnabledByDefault: opts["enabledByDefault"] == "True" || opts["enabledByDefault"] == "true",
		})
	}

	sort.Slice(instructions, func(i, j int) bool {
		iPos := instructions[i].Position
		jPos := instructions[j].Position

		if iPos == 0 {
			iPos = 999999
		}
		if jPos == 0 {
			jPos = 999999
		}
		return iPos < jPos
	})

	return instructions, nil
}

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		instructions, err := getInstructionFiles()
		if err != nil {
			log.Fatal(err)
		}

		var options []string
		var defaults []string
		for _, inst := range instructions {
			options = append(options, inst.Name)
			if inst.EnabledByDefault {
				defaults = append(defaults, inst.Name)
			}
		}

		var selected []string
		prompt := &survey.MultiSelect{
			Message: "Select instructions to activate:",
			Options: options,
			Default: defaults,
		}
		if err := survey.AskOne(prompt, &selected); err != nil {
			log.Fatal(err)
		}

		selectedMap := make(map[string]bool)
		for _, s := range selected {
			selectedMap[s] = true
		}

		var content string
		for _, inst := range instructions {
			if !selectedMap[inst.Name] {
				continue
			}
			fileContent, _ := Read(filepath.Join(".agentscope", inst.FileName))
			content += extractContent(fileContent) + "\n"
		}

		if err := os.WriteFile("./CLAUDE.md", []byte(content), 0644); err != nil {
			log.Fatal(err)
		}

		fmt.Println("CLAUDE.md updated successfully")
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
