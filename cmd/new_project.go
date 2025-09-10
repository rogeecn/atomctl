package cmd

import (
    "fmt"
    "io/fs"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "text/template"

    "go.ipao.vip/atomctl/templates"
    "go.ipao.vip/atomctl/pkg/utils/gomod"
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
    var (
        moduleName string
        inPlace    bool
    )

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
        rootDir = projectInfo.ProjectName
        if _, err := os.Stat(rootDir); err == nil {
            if !force {
                return fmt.Errorf("project directory %s already exists", rootDir)
            }
            log.Warnf("强制删除已存在的目录: %s", rootDir)
            if err := os.RemoveAll(rootDir); err != nil {
                return fmt.Errorf("failed to remove existing directory: %v", err)
            }
        }
        if err := os.MkdirAll(rootDir, 0o755); err != nil {
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
            return os.MkdirAll(targetPath, 0o755)
        }

        content, err := templates.Project.ReadFile(path)
        if err != nil {
            return err
        }

        if strings.HasSuffix(path, ".tpl") {
            if inPlace && strings.HasSuffix(path, string(os.PathSeparator)+"go.mod.tpl") {
                log.Infof("跳过已有文件: %s", filepath.Join(rootDir, "go.mod"))
                return nil
            }

            tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
            if err != nil {
                return err
            }

            targetPath = strings.TrimSuffix(targetPath, ".tpl")
            if !force {
                if _, err := os.Stat(targetPath); err == nil {
                    log.Warnf("文件已存在，跳过: %s", targetPath)
                    return nil
                }
            }
            log.Infof("创建文件: %s", targetPath)
            f, err := os.Create(targetPath)
            if err != nil {
                return errors.Wrapf(err, "创建文件失败 %s", targetPath)
            }
            defer f.Close()
            return tmpl.Execute(f, projectInfo)
        }

        if strings.HasSuffix(path, ".raw") {
            targetPath = strings.TrimSuffix(targetPath, ".raw")
        }

        if !force {
            if _, err := os.Stat(targetPath); err == nil {
                log.Warnf("文件已存在，跳过: %s", targetPath)
                return nil
            }
        }
        log.Infof("创建文件: %s", targetPath)
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
