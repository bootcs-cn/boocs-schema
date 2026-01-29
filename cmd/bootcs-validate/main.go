package main

import (
	"fmt"
	"os"

	"github.com/bootcs/bootcs-schema/internal/validator"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	verbose bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "bootcs-validate [course-dir]",
		Short: "Validate bootcs course and stage configurations",
		Long: `bootcs-validate validates course.yml and stage.yml files against
the bootcs schema definitions, and checks additional business rules.`,
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		RunE:    runValidate,
	}

	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// 子命令：生成文档（预留）
	genCmd := &cobra.Command{
		Use:   "generate [course-dir]",
		Short: "Generate README documentation from course/stages",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runGenerate,
	}
	rootCmd.AddCommand(genCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runValidate(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	v := validator.New(verbose)
	result := v.ValidateCourse(dir)

	// 打印结果
	for _, msg := range result.Messages {
		fmt.Println(msg)
	}

	if !result.Valid {
		return fmt.Errorf("validation failed with %d errors", result.ErrorCount)
	}

	fmt.Printf("\n✅ All validations passed! (%d stages checked)\n", result.StageCount)
	return nil
}

func runGenerate(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	fmt.Printf("📝 Generating documentation for %s...\n", dir)
	// TODO: 实现文档生成
	fmt.Println("⚠️  Document generation not yet implemented")
	return nil
}
