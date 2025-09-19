
# Implementation Plan: 优化 pkg/ast/provider 目录的代码组织逻辑与功能实现

**Branch**: `001-pkg-ast-provider` | **Date**: 2025-09-19 | **Spec**: `/projects/atomctl/specs/001-pkg-ast-provider/spec.md`
**Input**: Feature specification from `/projects/atomctl/specs/001-pkg-ast-provider/spec.md`

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
主要需求是优化 pkg/ast/provider 目录的代码组织逻辑与功能实现，补充完善测试用例。当前代码包含两个主要文件：provider.go（337行复杂的解析逻辑）和 render.go（65行渲染逻辑）。需要重构代码以提高可维护性、可测试性，并添加完整的测试覆盖。

## Technical Context
**Language/Version**: Go 1.24.0
**Primary Dependencies**: go/ast, go/parser, go/token, samber/lo, sirupsen/logrus, golang.org/x/tools/imports
**Storage**: N/A (file processing)
**Testing**: standard Go testing package with testify
**Target Platform**: Linux/macOS/Windows (CLI tool)
**Project Type**: Single project (CLI tool)
**Performance Goals**: <5s for large project parsing, <100MB memory usage
**Constraints**: Must maintain backward compatibility with existing @provider annotations
**Scale/Scope**: ~400 lines of existing code to refactor, target 90% test coverage

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### SOLID Principles Compliance
- [ ] **Single Responsibility**: Each component has single, clear responsibility
- [ ] **Open/Closed**: Design allows extension without modification
- [ ] **Liskov Substitution**: Subtypes can replace base types seamlessly
- [ ] **Interface Segregation**: Interfaces are specific and focused
- [ ] **Dependency Inversion**: Depend on abstractions, not concrete implementations

### KISS Principle Compliance
- [ ] Design avoids unnecessary complexity
- [ ] CLI interface maintains consistency
- [ ] Code generation logic is simple and direct
- [ ] Solutions are intuitive and easy to understand

### YAGNI Principle Compliance
- [ ] Only implementing clearly needed functionality
- [ ] No over-engineering or future-proofing without requirements
- [ ] Each feature has explicit user需求支撑
- [ ] No "might be useful" features without justification

### DRY Principle Compliance
- [ ] No code duplication across components
- [ ] Common functionality is abstracted and reused
- [ ] Template system avoids repetitive implementations
- [ ] Shared utilities are properly abstracted

### Code Quality Standards
- [ ] **Testing Discipline**: TDD approach with Red-Green-Refactor cycle
- [ ] **CLI Consistency**: Unified parameter formats and output standards
- [ ] **Error Handling**: Complete error information and recovery mechanisms
- [ ] **Performance**: Generation speed and memory usage requirements met

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

**Structure Decision**: [DEFAULT to Option 1 unless Technical Context indicates web/mobile app]

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
- Each contract → contract test task [P]
- Each entity → model creation task [P]
- Each user story → integration test task
- Implementation tasks to make tests pass

**Specific Task Categories for This Feature**:
- **Code Organization**: Refactor large Parse function into smaller, focused functions
- **Testing**: Create comprehensive test suite for all components
- **Interface Design**: Define clear interfaces for parser, validator, and renderer
- **Error Handling**: Implement robust error handling and recovery
- **Performance**: Ensure performance requirements are met
- **Documentation**: Add comprehensive documentation and examples

**Ordering Strategy**:
- TDD order: Tests before implementation
- Dependency order: Interfaces → Implementations → Integrations
- Mark [P] for parallel execution (independent files)

**Estimated Output**: 30-35 numbered, ordered tasks in tasks.md

**Key Implementation Notes**:
- Maintain backward compatibility with existing Parse() function
- Follow SOLID principles throughout refactoring
- Achieve 90% test coverage target
- Keep performance within requirements (<5s parsing, <100MB memory)

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
