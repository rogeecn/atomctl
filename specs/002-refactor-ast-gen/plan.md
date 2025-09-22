
# Implementation Plan: Refactor AST Generation Routes Workflow

**Branch**: `002-refactor-ast-gen` | **Date**: 2025-09-22 | **Spec**: [/specs/002-refactor-ast-gen/spec.md](/specs/002-refactor-ast-gen/spec.md)
**Input**: Feature specification from `/specs/002-refactor-ast-gen/spec.md`
**User Requirements**: 1. 重构 @pkg/ast/route/ 的实现流程，使更易读，逻辑更清晰，2.保证 @cmd/gen_route.go 对重构后方法调用的生效，3. 一切功能重构保证测试优先。

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
重构AST生成路由工作流程，提高代码可读性和逻辑清晰度，确保与现有gen_route.go命令的兼容性，并采用测试驱动开发方法。

## Technical Context
**Language/Version**: Go 1.21+
**Primary Dependencies**: go standard library (ast, parser, token), cobra CLI, logrus
**Storage**: File-based route definitions and generated Go code
**Testing**: Go testing with TDD approach (testing/fstest for filesystem tests)
**Target Platform**: CLI tool for Go projects
**Project Type**: Single project with existing pkg/ast/provider refactoring patterns
**Performance Goals**: Fast parsing (< 2s for typical project), minimal memory overhead
**Constraints**: Must maintain backward compatibility with existing @Router and @Bind annotations
**Scale/Scope**: Support medium to large Go projects with extensive route definitions

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### SOLID Principles Compliance
- [x] **Single Responsibility**: Route parsing, generation, and rendering will have separate, focused components
- [x] **Open/Closed**: Design will follow existing provider patterns with extensible interfaces
- [x] **Liskov Substitution**: New route components will implement consistent interfaces
- [x] **Interface Segregation**: Specific interfaces for parsing, generation, and validation
- [x] **Dependency Inversion**: Core functionality will depend on interfaces, not concrete implementations

### KISS Principle Compliance
- [x] Design avoids unnecessary complexity - will follow existing refactored provider patterns
- [x] CLI interface maintains consistency - existing gen_route.go interface preserved
- [x] Code generation logic is simple and direct - clear separation of concerns
- [x] Solutions are intuitive and easy to understand - follows established patterns

### YAGNI Principle Compliance
- [x] Only implementing clearly needed functionality - focus on readability and clarity improvements
- [x] No over-engineering or future-proofing without requirements - minimal changes to achieve goals
- [x] Each feature has explicit user requirements - based on gen_route.go compatibility needs
- [x] No "might be useful" features without justification - scope limited to refactoring

### DRY Principle Compliance
- [x] No code duplication across components - will share patterns with pkg/ast/provider
- [x] Common functionality is abstracted and reused - leverage existing interfaces and utilities
- [x] Template system avoids repetitive implementations - consistent with provider generation
- [x] Shared utilities are properly abstracted - reuse existing AST parsing infrastructure

### Code Quality Standards
- [x] **Testing Discipline**: TDD approach with Red-Green-Refactor cycle - testing first requirement
- [x] **CLI Consistency**: Unified parameter formats and output standards - existing interface maintained
- [x] **Error Handling**: Complete error information and recovery mechanisms - consistent with provider patterns
- [x] **Performance**: Generation speed and memory usage requirements met - <2s parsing goal

### Complexity Tracking
| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [Document any deviations from constitutional principles] | [Justification for complexity] | [Why simpler approach insufficient] |

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure]
```

**Structure Decision**: Option 1 - Single project with pkg/ast/route refactoring following pkg/ast/provider patterns

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude` for your AI assistant
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each contract → contract test implementation task [P]
- Each data model entity → implementation task [P]
- Each interface → component implementation task
- Integration tasks to ensure compatibility with gen_route.go
- Test-driven implementation following TDD principles

**Ordering Strategy**:
- TDD order: Write failing tests first, then implement to make tests pass
- Component dependency order: Core interfaces → Parsers → Builders → Validators → Renderers
- Backward compatibility: Ensure gen_route.go works throughout implementation
- Mark [P] for parallel execution (independent components)

**Estimated Output**: 25-30 numbered, ordered tasks in tasks.md covering:
- Core interface implementations
- Data model and error handling
- Route parsing and validation
- Template rendering and code generation
- Integration and compatibility testing
- Performance and validation testing

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v1.0.0 - See `/memory/constitution.md`*
