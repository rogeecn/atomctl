package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"git.ipao.vip/rogeecn/atomctl/templates"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// 验证包名是否合法：支持域名、路径分隔符和常见字符
var goPackageRegexp = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_.]*[a-zA-Z0-9](/[a-zA-Z0-9][a-zA-Z0-9-_.]*)*$`)

func isValidGoPackageName(name string) bool {
	return goPackageRegexp.MatchString(name)
}

func CommandNewProject(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "project",
		Aliases: []string{"p"},
		Short:   "new project",
		RunE:    commandNewProjectE,
	}

	root.AddCommand(cmd)
}

func commandNewProjectE(cmd *cobra.Command, args []string) error {
	moduleName := args[0]
	if !isValidGoPackageName(moduleName) {
		return fmt.Errorf("invalid module name: %s, should be a valid go package name", moduleName)
	}

	log.Info("创建项目: ", moduleName)

	var projectInfo struct {
		ModuleName  string
		ProjectName string
	}

	projectInfo.ModuleName = moduleName
	moduleSplitInfo := strings.Split(projectInfo.ModuleName, "/")
	projectInfo.ProjectName = moduleSplitInfo[len(moduleSplitInfo)-1]

	// 检查目录是否存在
	force, _ := cmd.Flags().GetBool("force")
	if _, err := os.Stat(projectInfo.ProjectName); err == nil {
		if !force {
			return fmt.Errorf("project directory %s already exists", projectInfo.ProjectName)
		}
		log.Warnf("强制删除已存在的目录: %s", projectInfo.ProjectName)
		if err := os.RemoveAll(projectInfo.ProjectName); err != nil {
			return fmt.Errorf("failed to remove existing directory: %v", err)
		}
	}

	// 创建项目根目录
	if err := os.MkdirAll(projectInfo.ProjectName, 0o755); err != nil {
		return fmt.Errorf("failed to create project directory: %v", err)
	}

	// 遍历和处理模板文件
	if err := fs.WalkDir(templates.Project, "project", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 计算相对路径，并处理隐藏文件
		relPath, err := filepath.Rel("project", path)
		if err != nil {
			return err
		}

		// 如果是隐藏文件模板，将文件名中的前缀 "-" 替换为 "."
		fileName := filepath.Base(relPath)
		if strings.HasPrefix(fileName, "-") {
			fileName = "." + strings.TrimPrefix(fileName, "-")
			relPath = filepath.Join(filepath.Dir(relPath), fileName)
		}

		targetPath := filepath.Join(projectInfo.ProjectName, relPath)
		if d.IsDir() {
			log.Infof("创建目录: %s", targetPath)
			return os.MkdirAll(targetPath, 0o755)
		}

		// 读取模板内容
		content, err := templates.Project.ReadFile(path)
		if err != nil {
			return err
		}

		// 处理模板文件
		if strings.HasSuffix(path, ".tpl") {
			tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
			if err != nil {
				return err
			}

			// 创建目标文件（去除.tpl后缀）
			targetPath = strings.TrimSuffix(targetPath, ".tpl")
			log.Infof("创建文件: %s", targetPath)
			f, err := os.Create(targetPath)
			if err != nil {
				return errors.Wrapf(err, "创建文件失败 %s", targetPath)
			}
			defer f.Close()

			return tmpl.Execute(f, projectInfo)
		}

		// 处理模板文件
		if strings.HasSuffix(path, ".raw") {
			// 创建目标文件（去除.tpl后缀）
			targetPath = strings.TrimSuffix(targetPath, ".raw")
			log.Infof("创建文件: %s", targetPath)
		}

		// 复制非模板文件
		return os.WriteFile(targetPath, content, 0o644)
	}); err != nil {
		return err
	}

	// 添加成功提示
	log.Info("🎉 项目创建成功!")
	log.Info("后续步骤:")
	log.Infof("  cd %s", projectInfo.ProjectName)
	log.Info("  go mod tidy")

	return nil
}
