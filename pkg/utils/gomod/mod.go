package gomod

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
)

var goMod *GoMod

type GoMod struct {
	file    *modfile.File
	modules map[string]ModuleInfo
}

type ModuleInfo struct {
	Name    string
	Version string
	Path    string
}

// ParseGoMod 解析当前目录下的go.mod文件
func Parse(modPath string) error {
	// 查找当前目录下的go.mod文件
	// 读取文件内容
	content, err := os.ReadFile(modPath)
	if err != nil {
		return err
	}

	// 使用官方包解析go.mod
	f, err := modfile.Parse(modPath, content, nil)
	if err != nil {
		return err
	}
	goMod = &GoMod{file: f, modules: make(map[string]ModuleInfo)}

	for _, require := range f.Require {
		if !require.Indirect {
			continue
		}

		name, err := getPackageName(require.Mod.Path, require.Mod.Version)
		if err != nil {
			continue
		}

		goMod.modules[require.Mod.Path] = ModuleInfo{
			Name:    name,
			Version: require.Mod.Version,
			Path:    require.Mod.Path,
		}
	}

	return nil
}

// GetModuleName 获取模块名
func GetModuleName() string {
	return goMod.file.Module.Mod.Path
}

// GetModuleVersion 获取模块版本
func GetModuleVersion() string {
	return goMod.file.Module.Mod.Version
}

func GetPackageModuleName(pkg string) string {
	if module, ok := goMod.modules[pkg]; ok {
		return module.Name
	}

	return filepath.Base(pkg)
}

// GetPackageModuleName 获取包的真实包名
func getPackageName(pkg, version string) (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}

	pkgPath := fmt.Sprintf("%s@%s", pkg, version)
	// 构建包的本地路径
	pkgLocalPath := filepath.Join(gopath, "pkg", "mod", pkgPath)

	// 获取目录下任意一个非_test.go文件，读取他的package name
	files, err := filepath.Glob(filepath.Join(pkgLocalPath, "*.go"))
	if err != nil {
		return "", err
	}

	packagePattern := regexp.MustCompile(`package\s+(\w+)`)
	if len(files) > 0 {
		for _, file := range files {
			if strings.HasSuffix(file, "_test.go") {
				continue
			}
			// 读取文件内容

			content, err := os.ReadFile(file)
			if err != nil {
				return "", err
			}

			packageName := packagePattern.FindStringSubmatch(string(content))
			if len(packageName) == 2 {
				return packageName[1], nil
			}
		}
	}
	// 读取go.mod 文件内容
	modFile := filepath.Join(pkgLocalPath, "go.mod")
	content, err := os.ReadFile(modFile)
	if err != nil {
		return "", err
	}

	f, err := modfile.Parse(modFile, content, nil)
	if err != nil {
		return "", err
	}

	path := f.Module.Mod.Path

	// 获取包名
	path, name := filepath.Split(path)
	versionPattern := regexp.MustCompile(`^v\d+$`)
	if versionPattern.MatchString(name) {
		_, name = filepath.Split(strings.TrimSuffix(path, "/"))
	}

	if strings.Contains(name, "-") {
		name = strings.ReplaceAll(name, "-", "")
	}

	return name, nil
}
