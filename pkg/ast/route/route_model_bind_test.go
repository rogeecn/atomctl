package route

import (
    "os"
    "path/filepath"
    "strings"
    "testing"

    "go.ipao.vip/atomctl/v2/pkg/utils/gomod"
)

// Test that @Bind with model(field) on a path parameter generates PathModel[T](field, key)
func Test_PathModelBind_FromRouteComments(t *testing.T) {
    dir := t.TempDir()
    src := `package v1

import (
    "context"
)

type User struct{}

type Demo struct{}

// @Router /users/:id [get]
// @Bind user path key(id) model(id)
func (d *Demo) Show(ctx context.Context, user *User) (*User, error) {
    return nil, nil
}`

    // minimal go.mod so gomod.GetPackageModuleName works without panic
    gomodPath := filepath.Join(dir, "go.mod")
    goModContent := "module example.com/test\n\ngo 1.23\n"
    if err := os.WriteFile(gomodPath, []byte(goModContent), 0o644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }

    if err := gomod.Parse(gomodPath); err != nil {
        t.Fatalf("gomod.Parse error: %v", err)
    }

    file := filepath.Join(dir, "demo.go")
    if err := os.WriteFile(file, []byte(src), 0o644); err != nil {
        t.Fatalf("write file: %v", err)
    }

    defs := ParseFile(file)
    if len(defs) != 1 {
        t.Fatalf("expected 1 route definition, got %d", len(defs))
    }
    if len(defs[0].Actions) != 1 {
        t.Fatalf("expected 1 action, got %d", len(defs[0].Actions))
    }
    act := defs[0].Actions[0]
    if len(act.Params) != 1 {
        t.Fatalf("expected 1 param, got %d", len(act.Params))
    }
    p := act.Params[0]
    if p.Position != PositionPath {
        t.Fatalf("expected path position, got %s", p.Position)
    }
    if p.Key != "id" {
        t.Fatalf("expected key=id, got %s", p.Key)
    }
    if p.ModelField != "id" {
        t.Fatalf("expected ModelField=id, got %s", p.ModelField)
    }
    if p.Type != "User" { // pointer should be trimmed for non-local
        t.Fatalf("expected Type=User, got %s", p.Type)
    }

    // Build render data and check binder token
    rd, err := buildRenderData(RenderBuildOpts{
        PackageName:    "v1",
        ProjectPackage: "example.com/test",
        Routes:         defs,
    })
    if err != nil {
        t.Fatalf("buildRenderData error: %v", err)
    }
    // Render to text and assert PathModel usage
    out, err := renderTemplate(rd)
    if err != nil {
        t.Fatalf("renderTemplate error: %v", err)
    }
    got := string(out)
    if !strings.Contains(got, "PathModel[User](\"id\", \"id\")") {
        t.Fatalf("expected generated code to contain PathModel[User](\"id\", \"id\"), got:\n%s", got)
    }
}
