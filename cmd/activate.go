/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"os/exec"
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

func generateRandomString(length int) string {
	bytes := make([]byte, length/2+1)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(bytes)[:length]
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

func promptForInstructions(instructions []InstructionFile) map[string]bool {
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
	return selectedMap
}

func buildClaudeMdContent(instructions []InstructionFile, selected map[string]bool) string {
	var content string
	for _, inst := range instructions {
		if !selected[inst.Name] {
			continue
		}
		fileContent, _ := Read(filepath.Join(".agentscope", inst.FileName))
		content += extractContent(fileContent) + "\n"
	}
	return content
}

func createPortalDirectory(projectRoot string) (string, func()) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	projectName := filepath.Base(projectRoot)
	sessionID := generateRandomString(16)
	portalDir := filepath.Join(homeDir, ".agentscope", projectName, sessionID)

	if err := os.MkdirAll(portalDir, 0755); err != nil {
		log.Fatalf("Failed to create portal directory: %v", err)
	}

	cleanup := func() {
		if err := os.RemoveAll(portalDir); err != nil {
			log.Printf("Warning: failed to cleanup portal directory: %v", err)
		}
	}

	return portalDir, cleanup
}

func createProjectSymlink(portalDir, projectRoot string) {
	symlinkName := generateRandomString(8)
	symlinkPath := filepath.Join(portalDir, symlinkName)

	relPath, err := filepath.Rel(portalDir, projectRoot)
	if err != nil {
		log.Fatalf("Failed to create relative path: %v", err)
	}

	if err := os.Symlink(relPath, symlinkPath); err != nil {
		log.Fatalf("Failed to create symlink: %v", err)
	}
}

func writeClaudeMd(portalDir, content string) {
	claudeMdPath := filepath.Join(portalDir, "CLAUDE.md")
	if err := os.WriteFile(claudeMdPath, []byte(content), 0644); err != nil {
		log.Fatalf("Failed to write CLAUDE.md: %v", err)
	}
}

func runClaude(portalDir string, args []string) {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		log.Fatal("claude CLI not found. Please install it first.")
	}

	cmd := exec.Command(claudePath, args...)
	cmd.Dir = portalDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		log.Fatalf("Failed to run claude: %v", err)
	}
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

var activateCmd = &cobra.Command{
	Use:   "activate [claude-args...]",
	Short: "Activate an agent session with Claude",
	Long: `Creates an isolated portal directory for running Claude with project-specific context.
All arguments after 'activate' are passed directly to the claude CLI.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectRoot, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}

		instructions, err := getInstructionFiles()
		if err != nil {
			log.Fatal(err)
		}

		selected := promptForInstructions(instructions)
		content := buildClaudeMdContent(instructions, selected)

		portalDir, cleanup := createPortalDirectory(projectRoot)
		defer cleanup()

		createProjectSymlink(portalDir, projectRoot)
		writeClaudeMd(portalDir, content)
		runClaude(portalDir, args)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
