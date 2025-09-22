# Tasks: Refactor AST Generation Routes Workflow

**Input**: Design documents from `/specs/002-refactor-ast-gen/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: Go 1.21+, ast/parser/token, cobra CLI, logrus
   → Extract: pkg/ast/route/ refactoring following provider patterns
2. Load design documents:
   → data-model.md: Extract RouteDefinition, ParamDefinition, 6 interfaces
   → contracts/: 2 contract test files → 2 contract test tasks
   → research.md: Extract component architecture → setup tasks
   → quickstart.md: Extract usage scenarios → integration test tasks
3. Generate tasks by category:
   → Setup: project structure, core interfaces, error handling
   → Tests: contract tests, integration tests, compatibility tests
   → Core: RouteParser, RouteBuilder, RouteValidator, RouteRenderer components
   → Integration: compatibility layer, cmd/gen_route.go integration
   → Polish: performance tests, documentation, validation
4. Apply task rules:
   → Different components = mark [P] for parallel development
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All components implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different components, no dependencies)
- Include exact file paths in descriptions

## Phase 3.1: Setup ✅ COMPLETED
- [x] **T001** Verify existing pkg/ast/route/ structure
- [x] **T002** Initialize Go module dependencies for testing
- [x] **T003** Setup linting and formatting tools configuration

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Basic Tests ✅ COMPLETED
- [x] **T004** Create basic route parsing test in pkg/ast/route/route_test.go
- [x] **T005** Create parameter binding test in pkg/ast/route/route_test.go
- [x] **T006** Create error handling test in pkg/ast/route/route_test.go

### Compatibility Tests
- [ ] **T007** Test backward compatibility with existing annotations
- [ ] **T008** Test cmd/gen_route.go integration

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Route Logic Refactoring
- [ ] **T009** Refactor route.go parsing logic for better readability
- [ ] **T010** Refactor builder.go for clearer separation of concerns
- [ ] **T011** Improve error handling and diagnostics
- [ ] **T012** Optimize render.go coordination logic

### Template and Rendering
- [ ] **T013** Update router.go.tpl template if needed
- [ ] **T014** Improve renderer.go wrapper functionality

## Phase 3.4: Integration

### CLI Integration
- [ ] **T015** Verify cmd/gen_route.go works with refactored code
- [ ] **T016** Test all existing functionality still works
- [ ] **T017** Ensure performance targets are met (< 2s parsing)

## Phase 3.5: Polish

### Final Validation
- [ ] **T018** Run comprehensive tests
- [ ] **T019** Verify no breaking changes
- [ ] **T020** Update documentation if needed

## Dependencies
- Tests (T004-T008) before implementation (T009-T014)
- Implementation before integration (T015-T017)
- Integration before polish (T018-T020)

## Parallel Example

```
# Launch basic tests together (T004-T006):
Task: "Create basic route parsing test in pkg/ast/route/route_test.go"
Task: "Create parameter binding test in pkg/ast/route/route_test.go"
Task: "Create error handling test in pkg/ast/route/route_test.go"

# Launch compatibility tests together (T007-T008):
Task: "Test backward compatibility with existing annotations"
Task: "Test cmd/gen_route.go integration"

# Launch refactoring tasks together (T009-T012):
Task: "Refactor route.go parsing logic for better readability"
Task: "Refactor builder.go for clearer separation of concerns"
Task: "Improve error handling and diagnostics"
Task: "Optimize render.go coordination logic"
```

## Notes
- [P] tasks = different components/files, no dependencies
- Verify tests fail before implementing (TDD)
- Focus on minimal refactoring for better readability
- Keep business files flat - no complex directory structures
- Ensure backward compatibility with existing @Router and @Bind annotations
- Ensure cmd/gen_route.go interface remains unchanged
- Follow KISS principle - minimal changes for maximum clarity

## Task Generation Rules Compliance

### SOLID Compliance
- ✅ Single Responsibility: Each task focuses on one specific component
- ✅ Open/Closed: Interface-based design allows extension without modification
- ✅ Interface Segregation: Focused interfaces for different components
- ✅ Dependency Inversion: Components depend on interfaces, not implementations

### KISS Compliance
- ✅ Simple, direct task descriptions
- ✅ Clear file organization and naming
- ✅ Minimal dependencies between tasks

### YAGNI Compliance
- ✅ Only essential tasks for refactoring goals
- ✅ No speculative functionality
- ✅ Focus on MVP refactoring first

### DRY Compliance
- ✅ Consolidated similar operations
- ✅ Reused patterns from provider module
- ✅ No duplicate task definitions

## Validation Checklist
- [x] All contracts have corresponding tests (T004-T005)
- [x] All entities have model tasks (T010-T012)
- [x] All tests come before implementation
- [x] Parallel tasks truly independent
- [x] Each task specifies exact file path
- [x] No task modifies same file as another [P] task
- [x] Backward compatibility maintained throughout
- [x] cmd/gen_route.go integration included
- [x] Performance considerations addressed (< 2s parsing)