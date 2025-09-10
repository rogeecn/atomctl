package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.ipao.vip/atomctl/v2/pkg/utils/gomod"
	"go.ipao.vip/atomctl/v2/templates"
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
	var (
		moduleName string
		inPlace    bool
	)

	// shared flags
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	baseDir, _ := cmd.Flags().GetString("dir")

	if len(args) == 0 {
		if _, err := os.Stat("go.mod"); err == nil {
			pwd, _ := os.Getwd()
			if err := gomod.Parse(filepath.Join(pwd, "go.mod")); err != nil {
				return fmt.Errorf("parse go.mod failed: %v", err)
			}
			moduleName = gomod.GetModuleName()
			inPlace = true
		} else {
			return fmt.Errorf("module name required or run inside an existing module (go.mod)")
		}
	} else {
		moduleName = args[0]
		if !isValidGoPackageName(moduleName) {
			return fmt.Errorf("invalid module name: %s, should be a valid go package name", moduleName)
		}
	}

	log.Info("创建项目: ", moduleName)

	var projectInfo struct {
		ModuleName  string
		ProjectName string
	}

	projectInfo.ModuleName = moduleName
	moduleSplitInfo := strings.Split(projectInfo.ModuleName, "/")
	projectInfo.ProjectName = moduleSplitInfo[len(moduleSplitInfo)-1]

	force, _ := cmd.Flags().GetBool("force")

	rootDir := "."
	if !inPlace {
		// honor base dir when creating a new project
		rootDir = filepath.Join(baseDir, projectInfo.ProjectName)
		if _, err := os.Stat(rootDir); err == nil {
			if !force {
				return fmt.Errorf("project directory %s already exists", rootDir)
			}
			log.Warnf("强制删除已存在的目录: %s", rootDir)
			if dryRun {
				log.Infof("[dry-run] 将删除目录: %s", rootDir)
			} else if err := os.RemoveAll(rootDir); err != nil {
				return fmt.Errorf("failed to remove existing directory: %v", err)
			}
		}
		if dryRun {
			log.Infof("[dry-run] 将创建目录: %s", rootDir)
		} else if err := os.MkdirAll(rootDir, 0o755); err != nil {
			return fmt.Errorf("failed to create project directory: %v", err)
		}
	}

	if err := fs.WalkDir(templates.Project, "project", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel("project", path)
		if err != nil {
			return err
		}

		fileName := filepath.Base(relPath)
		if strings.HasPrefix(fileName, "-") {
			fileName = "." + strings.TrimPrefix(fileName, "-")
			relPath = filepath.Join(filepath.Dir(relPath), fileName)
		}

		targetPath := filepath.Join(rootDir, relPath)
		if d.IsDir() {
			log.Infof("创建目录: %s", targetPath)
			if dryRun {
				log.Infof("[dry-run] mkdir -p %s", targetPath)
				return nil
			}
			return os.MkdirAll(targetPath, 0o755)
		}

		content, err := templates.Project.ReadFile(path)
		if err != nil {
			return err
		}

		// 渲染判定：优先(.tpl/.raw) > 行内指令 > 模板分隔符
		isTpl := false
		if strings.HasSuffix(path, ".tpl") {
			isTpl = true
			targetPath = strings.TrimSuffix(targetPath, ".tpl")
		} else if strings.HasSuffix(path, ".raw") {
			isTpl = false
			targetPath = strings.TrimSuffix(targetPath, ".raw")
		} else if bytes.Contains(content, []byte("atomctl:mode=tpl")) {
			isTpl = true
		} else if bytes.Contains(content, []byte("atomctl:mode=raw")) {
			isTpl = false
		} else if bytes.Contains(content, []byte("{{")) && bytes.Contains(content, []byte("}}")) {
			isTpl = true
		}

		if inPlace && strings.HasSuffix(path, string(os.PathSeparator)+"go.mod.tpl") {
			if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
				log.Infof("跳过已有文件: %s", filepath.Join(rootDir, "go.mod"))
				return nil
			}
		}

		if !force {
			if _, err := os.Stat(targetPath); err == nil {
				log.Warnf("文件已存在，跳过: %s", targetPath)
				return nil
			}
		}

		if isTpl {
			tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
			if err != nil {
				return err
			}
			log.Infof("[渲染] 文件: %s", targetPath)
			if dryRun {
				log.Infof("[dry-run] render > %s", targetPath)
				return nil
			}
			f, err := os.Create(targetPath)
			if err != nil {
				return errors.Wrapf(err, "创建文件失败 %s", targetPath)
			}
			defer f.Close()
			return tmpl.Execute(f, projectInfo)
		}

		log.Infof("[复制] 文件: %s", targetPath)
		if dryRun {
			log.Infof("[dry-run] write > %s", targetPath)
			return nil
		}
		return os.WriteFile(targetPath, content, 0o644)
	}); err != nil {
		return err
	}

	if inPlace {
		log.Info("🎉 项目初始化成功 (当前目录)!")
	} else {
		log.Info("🎉 项目创建成功!")
		log.Info("后续步骤:")
		log.Infof("  cd %s", projectInfo.ProjectName)
	}
	log.Info("  go mod tidy")

	return nil
}
