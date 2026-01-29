package main

import (
	"fmt"
	"os"

	"github.com/bootcs/bootcs-schema/internal/generator"
	"github.com/bootcs/bootcs-schema/internal/validator"
	"github.com/spf13/cobra"
)

var (
	version    = "dev"
	verbose    bool
	outputFile string
	dryRun     bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "schema-tool [course-dir]",
		Short: "Bootcs course schema validation and documentation tool",
		Long: `schema-tool validates course.yml and stage.yml files against
the bootcs schema definitions, generates README documentation, and more.`,
		Version: version,
		Args:    cobra.MaximumNArgs(1),
		RunE:    runValidate,
	}

	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// 子命令：生成文档
	genCmd := &cobra.Command{
		Use:   "generate [course-dir]",
		Short: "Generate README documentation from course/stages",
		Long: `Generate a README.md file with course information and stages table.

Examples:
  schema-tool generate .
  schema-tool generate /path/to/course --output README.md
  schema-tool generate . --dry-run`,
		Args: cobra.MaximumNArgs(1),
		RunE: runGenerate,
	}
	genCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path (default: course-dir/README.md)")
	genCmd.Flags().BoolVar(&dryRun, "dry-run", false, "print to stdout instead of writing file")
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

	gen := generator.New(dir)

	if dryRun {
		content, err := gen.GenerateREADME()
		if err != nil {
			return err
		}
		fmt.Print(content)
		return nil
	}

	if err := gen.WriteREADME(outputFile); err != nil {
		return err
	}

	output := outputFile
	if output == "" {
		output = dir + "/README.md"
	}
	fmt.Printf("✅ Generated %s\n", output)
	return nil
}
