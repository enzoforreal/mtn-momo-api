package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Long:  `Start the momo API server.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting the server...")
		exampleDir, _ := filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "example"))
		execCmd := exec.Command("go", "run", "main.go")
		execCmd.Dir = exampleDir
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr
		if err := execCmd.Run(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
