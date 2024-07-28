package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Run integration tests",
	Long:  `Run all integration tests for the momo API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running integration tests...")
		testsDir, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "tests"))
		execCmd := exec.Command("bash", filepath.Join(testsDir, "run_integration_tests.sh"))
		execCmd.Dir = testsDir
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil {
			fmt.Printf("Error running integration tests: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
