# Tasks: 优化 pkg/ast/provider 目录的代码组织逻辑与功能实现

**Input**: Design documents from `/projects/atomctl/specs/001-pkg-ast-provider/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [x] T001 Create test directory structure and test data files
- [x] T002 Set up testing dependencies (testify, gomock) and test configuration
- [x] T003 [P] Configure linting and benchmark tools for code quality

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Contract Tests (from contracts/)
- [x] T004 [P] Contract test Parser interface in tests/contract/test_parser_api.go
- [x] T005 [P] Contract test Validator interface in tests/contract/test_validator_api.go
- [x] T006 [P] Contract test Renderer interface in tests/contract/test_renderer_api.go
- [x] T007 [P] Contract test Validation Rules in tests/contract/test_validation_rules.go

### Integration Tests (from user stories)
- [x] T008 [P] Integration test single file parsing workflow in tests/integration/test_file_parsing.go
- [x] T009 [P] Integration test complex project with multiple providers in tests/integration/test_project_parsing.go
- [x] T010 [P] Integration test error handling and recovery in tests/integration/test_error_handling.go
- [x] T011 [P] Integration test performance requirements in tests/integration/test_performance.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### Data Model Implementation (from data-model.md)
- [x] T012 [P] Provider data structure in pkg/ast/provider/types.go
- [x] T013 [P] ProviderMode and InjectionMode enums in pkg/ast/provider/modes.go
- [x] T014 [P] InjectParam and SourceLocation structs in pkg/ast/provider/types.go
- [x] T015 [P] ParserConfig and ParserContext in pkg/ast/provider/config.go

### Interface Implementation (from contracts/)
- [x] T016 [P] Parser interface implementation in pkg/ast/provider/parser_interface.go
- [x] T017 [P] Validator interface implementation in pkg/ast/provider/validator.go
- [x] T018 [P] Renderer interface implementation in pkg/ast/provider/renderer.go
- [x] T019 [P] Error types and error handling in pkg/ast/provider/errors.go

### Core Logic Refactoring
- [ ] T020 Extract comment parsing logic from Parse function in pkg/ast/provider/comment_parser.go
- [ ] T021 Extract AST traversal logic in pkg/ast/provider/ast_walker.go
- [ ] T022 Extract import resolution logic in pkg/ast/provider/import_resolver.go
- [ ] T023 Extract provider building logic in pkg/ast/provider/builder.go
- [ ] T024 Refactor main Parse function to use new components in pkg/ast/provider/parser.go

### Validation Rules Implementation (from contracts/)
- [x] T025 [P] Implement struct name validation rule in pkg/ast/provider/validator.go
- [x] T026 [P] Implement return type validation rule in pkg/ast/provider/validator.go
- [x] T027 [P] Implement provider mode validation rule in pkg/ast/provider/validator.go
- [x] T028 [P] Implement injection params validation rule in pkg/ast/provider/validator.go
- [x] T029 [P] Implement validation report generation in pkg/ast/provider/report_generator.go

## Phase 3.4: Integration

### Parser Integration
- [ ] T030 Integrate new parser components in pkg/ast/provider/parser.go
- [ ] T031 Implement backward compatibility layer for existing Parse function
- [ ] T032 Add configuration and context management to parser
- [ ] T033 Implement caching mechanism for performance optimization

### Validator Integration
- [ ] T034 [P] Integrate validation rules in validator implementation
- [ ] T035 Implement validation report generation
- [ ] T036 Add custom validation rule registration system

### Renderer Integration
- [ ] T037 [P] Update renderer to use new data structures
- [ ] T038 Implement template function registration system
- [ ] T039 Add custom template support for extensibility

### Error Handling Integration
- [ ] T040 [P] Implement comprehensive error handling across all components
- [ ] T041 Add error recovery mechanisms
- [ ] T042 Implement structured logging for debugging

## Phase 3.5: Polish

### Unit Tests
- [ ] T043 [P] Unit tests for comment parsing in tests/unit/test_comment_parser.go
- [ ] T044 [P] Unit tests for AST walker in tests/unit/test_ast_walker.go
- [ ] T045 [P] Unit tests for import resolver in tests/unit/test_import_resolver.go
- [ ] T046 [P] Unit tests for provider builder in tests/unit/test_builder.go
- [ ] T047 [P] Unit tests for validation rules in tests/unit/test_validation_rules.go

### Performance Tests
- [ ] T048 Performance benchmark for single file parsing (<100ms)
- [ ] T049 Performance benchmark for large project parsing (<5s)
- [ ] T050 Memory usage validation (<100MB normal usage)
- [ ] T051 Stress test with concurrent file parsing

### Documentation and Examples
- [ ] T052 [P] Update package documentation and godoc comments
- [ ] T053 [P] Create usage examples and migration guide
- [ ] T054 [P] Document new interfaces and extension points
- [ ] T055 [P] Update README and quickstart guide

### Final Integration
- [ ] T056 Remove code duplication and consolidate utilities
- [ ] T057 Final performance optimization and validation
- [ ] T058 Integration test for complete backward compatibility
- [ ] T059 Run full test suite and validate 90% coverage
- [ ] T060 Final code review and cleanup

## Dependencies
- Tests (T004-T011) before implementation (T012-T042)
- Data models (T012-T015) before interfaces (T016-T019)
- Core logic refactoring (T020-T024) before integration (T030-T042)
- Integration (T030-T042) before polish (T043-T060)
- Unit tests (T043-T047) before performance tests (T048-T051)
- Documentation (T052-T055) before final integration (T056-T060)

## Parallel Execution Groups

### Group 1: Contract Tests (Can run in parallel)
```
Task: "Contract test Parser interface in tests/contract/test_parser_api.go"
Task: "Contract test Validator interface in tests/contract/test_validator_api.go"
Task: "Contract test Renderer interface in tests/contract/test_renderer_api.go"
Task: "Contract test Validation Rules in tests/contract/test_validation_rules.go"
```

### Group 2: Integration Tests (Can run in parallel)
```
Task: "Integration test single file parsing workflow in tests/integration/test_file_parsing.go"
Task: "Integration test complex project with multiple providers in tests/integration/test_project_parsing.go"
Task: "Integration test error handling and recovery in tests/integration/test_error_handling.go"
Task: "Integration test performance requirements in tests/integration/test_performance.go"
```

### Group 3: Data Model Implementation (Can run in parallel)
```
Task: "Provider data structure in pkg/ast/provider/types.go"
Task: "ProviderMode and InjectionMode enums in pkg/ast/provider/modes.go"
Task: "InjectParam and SourceLocation structs in pkg/ast/provider/types.go"
Task: "ParserConfig and ParserContext in pkg/ast/provider/config.go"
```

### Group 4: Interface Implementation (Can run in parallel)
```
Task: "Parser interface implementation in pkg/ast/provider/parser.go"
Task: "Validator interface implementation in pkg/ast/provider/validator.go"
Task: "Renderer interface implementation in pkg/ast/provider/renderer.go"
Task: "Error types and error handling in pkg/ast/provider/errors.go"
```

### Group 5: Validation Rules (Can run in parallel)
```
Task: "Struct name validation rule in pkg/ast/provider/validation/struct_name.go"
Task: "Provider mode validation rules in pkg/ast/provider/validation/mode_rules.go"
Task: "Inject params validation rule in pkg/ast/provider/validation/inject_params.go"
Task: "Package alias validation rule in pkg/ast/provider/validation/package_alias.go"
Task: "Mode-specific validation rules in pkg/ast/provider/validation/mode_specific.go"
```

### Group 6: Unit Tests (Can run in parallel)
```
Task: "Unit tests for comment parsing in tests/unit/test_comment_parser.go"
Task: "Unit tests for AST walker in tests/unit/test_ast_walker.go"
Task: "Unit tests for import resolver in tests/unit/test_import_resolver.go"
Task: "Unit tests for provider builder in tests/unit/test_builder.go"
Task: "Unit tests for validation rules in tests/unit/test_validation_rules.go"
```

## Critical Implementation Notes

### Backward Compatibility
- Maintain existing `Parse(filename string) []Provider` function signature
- Keep existing `Render(filename string, conf []Provider)` function signature
- Ensure all existing @provider annotation formats continue to work

### SOLID Principles
- **Single Responsibility**: Each component has one clear purpose
- **Open/Closed**: Design for extension through interfaces
- **Liskov Substitution**: All implementations can be substituted
- **Interface Segregation**: Keep interfaces focused and specific
- **Dependency Inversion**: Depend on abstractions, not concrete types

### Performance Requirements
- Single file parsing: <100ms
- Large project parsing: <5s
- Memory usage: <100MB (normal operation)
- Test coverage: ≥90%

### Quality Gates
- All tests must pass before merging
- Code must follow Go formatting standards
- No linting warnings or errors
- Documentation must be complete and accurate

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing (TDD approach)
- Commit after each significant task completion
- Each task must be specific and actionable
- Maintain backward compatibility throughout refactoring

## Task Generation Rules

### SOLID Compliance
- **Single Responsibility**: Each task focuses on one specific component
- **Open/Closed**: Design tasks to allow extension without modification
- **Interface Segregation**: Create focused interfaces for different task types

### KISS Compliance
- Keep task descriptions simple and direct
- Avoid over-complicating task dependencies
- Use intuitive file naming and structure

### YAGNI Compliance
- Only create tasks for clearly needed functionality
- Avoid speculative tasks without direct requirements
- Focus on MVP implementation first

### DRY Compliance
- Abstract common patterns into reusable task templates
- Avoid duplicate task definitions
- Consolidate similar operations where possible

### From Design Documents
- **Contracts**: Each contract file → contract test task [P]
- **Data Model**: Each entity → model creation task [P]
- **User Stories**: Each story → integration test [P]

### Ordering
- Setup → Tests → Models → Services → Endpoints → Polish
- Dependencies block parallel execution

### Validation Checklist
- [ ] All contracts have corresponding tests
- [ ] All entities have model tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task