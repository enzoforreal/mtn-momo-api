package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Go dependencies",
	Long:  `Update all Go dependencies for the momo API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Updating Go dependencies...")
		projectDir, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0])))
		execCmd := exec.Command("go", "get", "-u", "./...")
		execCmd.Dir = projectDir
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil {
			fmt.Printf("Error updating dependencies: %v\n", err)
		}

		execCmd = exec.Command("go", "mod", "tidy")
		execCmd.Dir = projectDir
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil {
			fmt.Printf("Error tidying dependencies: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
