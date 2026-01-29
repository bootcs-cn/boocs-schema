package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bootcs/bootcs-schema/schemas"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

// Result represents validation result
type Result struct {
	Valid      bool
	ErrorCount int
	StageCount int
	Messages   []string
}

// Validator validates course and stage configurations
type Validator struct {
	verbose      bool
	courseSchema *jsonschema.Schema
	stageSchema  *jsonschema.Schema
}

// New creates a new Validator
func New(verbose bool) *Validator {
	v := &Validator{verbose: verbose}
	v.loadSchemas()
	return v
}

func (v *Validator) loadSchemas() {
	compiler := jsonschema.NewCompiler()

	// 加载 course schema
	courseData, _ := schemas.FS.ReadFile("course.schema.json")
	if err := compiler.AddResource("course.schema.json", strings.NewReader(string(courseData))); err == nil {
		v.courseSchema, _ = compiler.Compile("course.schema.json")
	}

	// 加载 stage schema
	compiler2 := jsonschema.NewCompiler()
	stageData, _ := schemas.FS.ReadFile("stage.schema.json")
	if err := compiler2.AddResource("stage.schema.json", strings.NewReader(string(stageData))); err == nil {
		v.stageSchema, _ = compiler2.Compile("stage.schema.json")
	}
}

// ValidateCourse validates a course directory
func (v *Validator) ValidateCourse(dir string) *Result {
	result := &Result{Valid: true}

	// 1. 验证 course.yml
	coursePath := filepath.Join(dir, "course.yml")
	if _, err := os.Stat(coursePath); os.IsNotExist(err) {
		result.addError("❌ course.yml not found")
		return result
	}

	courseData, err := v.loadYAML(coursePath)
	if err != nil {
		result.addError("❌ course.yml: %v", err)
		return result
	}

	if err := v.validateSchema(v.courseSchema, courseData, "course.yml"); err != nil {
		result.addError("❌ course.yml schema: %v", err)
	} else {
		result.addInfo("✅ course.yml: schema valid")
	}

	// 2. 获取 stage_order
	stageOrder, _ := v.getStageOrder(courseData)

	// 3. 验证 stages 目录
	stagesDir := filepath.Join(dir, "stages")
	if _, err := os.Stat(stagesDir); os.IsNotExist(err) {
		result.addError("❌ stages/ directory not found")
		return result
	}

	// 4. 验证每个 stage
	entries, _ := os.ReadDir(stagesDir)
	actualStages := make(map[string]bool)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		stageName := entry.Name()
		actualStages[stageName] = true

		stageDir := filepath.Join(stagesDir, stageName)
		v.validateStage(stageDir, stageName, result)
		result.StageCount++
	}

	// 5. 验证 stage_order 一致性
	v.validateStageOrder(stageOrder, actualStages, result)

	return result
}

func (v *Validator) validateStage(stageDir, stageName string, result *Result) {
	prefix := fmt.Sprintf("stages/%s", stageName)

	// 检查必需文件
	requiredFiles := []string{"stage.yml", "README.md", "LEARNING.md"}
	for _, file := range requiredFiles {
		path := filepath.Join(stageDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			result.addError("❌ %s/%s: missing", prefix, file)
		}
	}

	// 验证 stage.yml schema
	stageYmlPath := filepath.Join(stageDir, "stage.yml")
	if stageData, err := v.loadYAML(stageYmlPath); err == nil {
		if err := v.validateSchema(v.stageSchema, stageData, prefix+"/stage.yml"); err != nil {
			result.addError("❌ %s/stage.yml schema: %v", prefix, err)
		} else {
			// 验证 slug 与目录名匹配
			if slug, ok := stageData["slug"].(string); ok && slug != stageName {
				result.addError("❌ %s/stage.yml: slug '%s' does not match directory name '%s'", prefix, slug, stageName)
			}
			result.addInfo("✅ %s/stage.yml: valid", prefix)
		}
	}

	// 验证 LEARNING.md 行数 (60-100)
	learningPath := filepath.Join(stageDir, "LEARNING.md")
	if lines, err := countLines(learningPath); err == nil {
		if lines < 60 {
			result.addWarning("⚠️  %s/LEARNING.md: %d lines (recommended: 60-100)", prefix, lines)
		} else if lines > 100 {
			result.addWarning("⚠️  %s/LEARNING.md: %d lines (recommended: 60-100)", prefix, lines)
		} else if v.verbose {
			result.addInfo("✅ %s/LEARNING.md: %d lines", prefix, lines)
		}
	}

	// 验证 README.md 行数 (30-60)
	readmePath := filepath.Join(stageDir, "README.md")
	if lines, err := countLines(readmePath); err == nil {
		if lines < 30 {
			result.addWarning("⚠️  %s/README.md: %d lines (recommended: 30-60)", prefix, lines)
		} else if lines > 60 {
			result.addWarning("⚠️  %s/README.md: %d lines (recommended: 30-60)", prefix, lines)
		} else if v.verbose {
			result.addInfo("✅ %s/README.md: %d lines", prefix, lines)
		}
	}
}

func (v *Validator) validateStageOrder(stageOrder []string, actualStages map[string]bool, result *Result) {
	// 检查 stage_order 中的 stage 是否都存在
	for _, stage := range stageOrder {
		if !actualStages[stage] {
			result.addError("❌ stage_order: '%s' declared but directory not found", stage)
		}
	}

	// 检查是否有目录未在 stage_order 中
	declaredStages := make(map[string]bool)
	for _, s := range stageOrder {
		declaredStages[s] = true
	}
	for stage := range actualStages {
		if !declaredStages[stage] {
			result.addWarning("⚠️  stages/%s: directory exists but not in stage_order", stage)
		}
	}
}

func (v *Validator) loadYAML(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (v *Validator) validateSchema(schema *jsonschema.Schema, data map[string]interface{}, name string) error {
	if schema == nil {
		return fmt.Errorf("schema not loaded")
	}

	// 转换为 JSON 再验证（jsonschema 库需要）
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var jsonObj interface{}
	if err := json.Unmarshal(jsonData, &jsonObj); err != nil {
		return err
	}

	return schema.Validate(jsonObj)
}

func (v *Validator) getStageOrder(courseData map[string]interface{}) ([]string, bool) {
	order, ok := courseData["stage_order"].([]interface{})
	if !ok {
		return nil, false
	}

	result := make([]string, len(order))
	for i, s := range order {
		result[i], _ = s.(string)
	}
	return result, true
}

func countLines(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return len(strings.Split(string(data), "\n")), nil
}

func (r *Result) addError(format string, args ...interface{}) {
	r.Valid = false
	r.ErrorCount++
	r.Messages = append(r.Messages, fmt.Sprintf(format, args...))
}

func (r *Result) addWarning(format string, args ...interface{}) {
	r.Messages = append(r.Messages, fmt.Sprintf(format, args...))
}

func (r *Result) addInfo(format string, args ...interface{}) {
	r.Messages = append(r.Messages, fmt.Sprintf(format, args...))
}
