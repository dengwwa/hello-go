package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// ProjectTables 存储项目信息
type ProjectTables struct {
	Devices bool `json:"devices"`
	Users   bool `json:"users"`
}

// Metadata 元数据
type Metadata struct {
	Name string `json:"name"`
}

// IcebergTable 代表一个表项
type IcebergTable struct {
	ApiVersion string   `json:"apiVersion"`
	Kind       string   `json:"kind"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
}

// Spec 规格
type Spec struct {
	Schema          map[string]interface{} `json:"schema"`
	SortedBy        []string               `json:"sorted-by"`
	TableProperties map[string]interface{} `json:"table-properties"`
}

// PagedList 页面列表
type PagedList struct {
	Kind     string         `json:"kind"`
	Metadata interface{}    `json:"metadata"`
	Items    []IcebergTable `json:"items"`
}

func main() {
	// 读取JSON文件
	filePath := "/Users/zyq/Dev/IdeaProjects/learn/hello-go/work/icberg/data/FeHelper-20251107104806.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	// 解析JSON
	var pagedList PagedList
	err = json.Unmarshal(data, &pagedList)
	if err != nil {
		log.Fatalf("解析JSON失败: %v", err)
	}

	// 存储项目表信息
	projects := make(map[string]ProjectTables)

	// 提取所有表名并分析
	tableCount := 0
	for _, item := range pagedList.Items {
		tableCount++
		tableName := item.Metadata.Name

		// 提取项目名（去掉.devices或.users后缀）
		var projectName string
		var tableType string

		if strings.HasSuffix(tableName, ".devices") {
			projectName = strings.TrimSuffix(tableName, ".devices")
			tableType = "devices"
		} else if strings.HasSuffix(tableName, ".users") {
			projectName = strings.TrimSuffix(tableName, ".users")
			tableType = "users"
		}

		// 更新项目信息
		if projectName != "" {
			if _, exists := projects[projectName]; !exists {
				projects[projectName] = ProjectTables{
					Devices: false,
					Users:   false,
				}
			}

			// 获取当前值并更新
			current := projects[projectName]
			if tableType == "devices" {
				current.Devices = true
			} else if tableType == "users" {
				current.Users = true
			}
			projects[projectName] = current
		}
	}

	// 分析结果
	completeProjects := []string{}
	incompleteProjects := []string{}
	completeCount := 0
	incompleteCount := 0

	// 按项目名排序
	projectNames := make([]string, 0, len(projects))
	for projectName := range projects {
		projectNames = append(projectNames, projectName)
	}
	sort.Strings(projectNames)

	fmt.Println("表完整性分析报告 (Go语言版本)")
	fmt.Println("=" + strings.Repeat("=", 58))
	fmt.Printf("总表数: %d\n", tableCount)
	fmt.Printf("项目数: %d\n", len(projects))
	fmt.Println()

	// 检查每个项目
	missingTables := []string{}

	for _, projectName := range projectNames {
		tables := projects[projectName]
		if tables.Devices && tables.Users {
			completeProjects = append(completeProjects, projectName)
			completeCount++
			fmt.Printf("✓ 完整: %s (devices + users)\n", projectName)
		} else {
			incompleteProjects = append(incompleteProjects, projectName)
			incompleteCount++

			missing := []string{}
			if !tables.Devices {
				missing = append(missing, "devices")
				missingTables = append(missingTables, fmt.Sprintf("%s.devices", projectName))
			}
			if !tables.Users {
				missing = append(missing, "users")
				missingTables = append(missingTables, fmt.Sprintf("%s.users", projectName))
			}

			fmt.Printf("✗ 不完整: %s - 缺少: %s\n", projectName, strings.Join(missing, ", "))
		}
	}

	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("完整项目: %d 个\n", completeCount)
	fmt.Printf("不完整项目: %d 个\n", incompleteCount)
	fmt.Println()

	if len(missingTables) > 0 {
		fmt.Println("缺失的表:")
		fmt.Println("-" + strings.Repeat("-", 47))
		for i, tableName := range missingTables {
			fmt.Printf("%d. %s\n", i+1, tableName)
		}

		// 保存缺失表到文件
		err = saveMissingTables(missingTables)
		if err != nil {
			fmt.Printf("保存缺失表文件失败: %v\n", err)
		} else {
			fmt.Println()
			fmt.Println("缺失表已保存到: /workspace/data/missing_tables_go.txt")
		}
	} else {
		fmt.Println("所有表都完整！")
	}

	// 保存详细报告
	err = saveDetailedReport(projects, projectNames, completeProjects, incompleteProjects)
	if err != nil {
		fmt.Printf("保存详细报告失败: %v\n", err)
	} else {
		fmt.Println("详细报告已保存到: /workspace/data/detailed_report_go.txt")
	}
}

// saveMissingTables 保存缺失表列表
func saveMissingTables(missingTables []string) error {
	// 确保目录存在
	err := os.MkdirAll("/Users/zyq/Dev/IdeaProjects/learn/hello-go/work/data", 0755)
	if err != nil {
		return err
	}

	content := "缺失的表列表 (Go语言分析)\n"
	content += "===========================\n"
	content += fmt.Sprintf("总计: %d 个表\n\n", len(missingTables))

	for i, tableName := range missingTables {
		content += fmt.Sprintf("%d. %s\n", i+1, tableName)
	}

	return os.WriteFile("/Users/zyq/Dev/IdeaProjects/learn/hello-go/work/data/missing_tables_go.txt", []byte(content), 0644)
}

// saveDetailedReport 保存详细报告
func saveDetailedReport(projects map[string]ProjectTables, projectNames, completeProjects, incompleteProjects []string) error {
	// 确保目录存在
	err := os.MkdirAll("/Users/zyq/Dev/IdeaProjects/learn/hello-go/work/icberg/data", 0755)
	if err != nil {
		return err
	}

	content := "Go语言生成的表完整性详细报告\n"
	content += "=======================================================\n\n"

	// 统计信息
	content += "统计信息:\n"
	content += "----------\n"
	content += fmt.Sprintf("项目总数: %d\n", len(projects))
	content += fmt.Sprintf("完整项目: %d\n", len(completeProjects))
	content += fmt.Sprintf("不完整项目: %d\n\n", len(incompleteProjects))

	// 完整项目列表
	content += "完整项目 (有devices和users表):\n"
	content += "-------------------------------\n"
	for i, projectName := range completeProjects {
		content += fmt.Sprintf("%d. %s\n", i+1, projectName)
	}
	content += "\n"

	// 不完整项目列表
	if len(incompleteProjects) > 0 {
		content += "不完整项目:\n"
		content += "----------------\n"
		for _, projectName := range incompleteProjects {
			tables := projects[projectName]
			missing := []string{}
			if !tables.Devices {
				missing = append(missing, "devices")
			}
			if !tables.Users {
				missing = append(missing, "users")
			}
			content += fmt.Sprintf("%s - 缺少: %s\n", projectName, strings.Join(missing, ", "))
		}
		content += "\n"
	}

	// 预期vs实际
	expectedTables := len(projects) * 2
	actualTables := len(completeProjects)*2 + len(incompleteProjects)
	content += "预期vs实际:\n"
	content += "------------\n"
	content += fmt.Sprintf("预期表数: %d (个项目 × 2表/项目)\n", expectedTables)
	content += fmt.Sprintf("实际表数: %d\n", actualTables)
	content += fmt.Sprintf("缺失表数: %d\n", expectedTables-actualTables)

	return os.WriteFile("/Users/zyq/Dev/IdeaProjects/learn/hello-go/work/icberg/data/detailed_report_go.txt", []byte(content), 0644)
}
