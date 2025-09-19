<!--
# Sync Impact Report
**Version Change**: Template (unversioned) → 1.0.0
**Bump Rationale**: Initial constitution creation with SOLID, KISS, YAGNI, DRY principles

## Modified Principles
- Template placeholders → Concrete SOLID, KISS, YAGNI, DRY principles
- Generic governance → Specific CLI tool governance for atomctl

## Added Sections
- SOLID Principles (5 principles with detailed explanations)
- KISS Principle (simplicity focus)
- YAGNI Principle (avoiding over-engineering)
- DRY Principle (eliminating duplication)
- Code Quality Standards (testing, CLI consistency, error handling)
- Performance Requirements (generation speed, UX consistency)
- Development Workflow (code generation process, template management)
- Governance (amendment process, versioning policy, compliance review)

## Templates Updated
- ✅ `.specify/templates/plan-template.md` - Updated Constitution Check section
- ✅ `.specify/templates/tasks-template.md` - Updated Task Generation Rules
- ✅ `.specify/templates/spec-template.md` - No changes needed (already principle-agnostic)
- ✅ `.specify/templates/agent-file-template.md` - No changes needed (not found)

## Follow-up TODOs
- Review all generated code for compliance with new constitution
- Update any existing command documentation to reference new principles
- Consider adding principle-specific validation to CI/CD pipeline
-->

# atomctl Constitution

## Core Principles

### I. SOLID Principles
**Single Responsibility Principle**: 每个组件、类和函数必须有单一明确的职责。命令、生成器和工具应当各司其职，避免功能混杂。

**Open/Closed Principle**: 软件实体应该对扩展开放，对修改关闭。通过插件化的模板系统和可配置的生成器来实现扩展性，而不是修改核心代码。

**Liskov Substitution Principle**: 子类型必须能够无缝替换其基类型。所有生成器接口和命令处理器都应遵循一致的契约。

**Interface Segregation Principle**: 接口应该特定且专注，避免"胖接口"。为不同类型的生成操作定义专门的接口。

**Dependency Inversion Principle**: 依赖抽象而非具体实现。核心功能应依赖接口，而不是特定的模板引擎或数据库驱动。

### II. KISS (Keep It Simple, Stupid)
追求终极的简洁性和直观性，避免不必要的复杂性。命令行接口必须保持一致性，参数命名和输出格式应当直观易懂。代码生成逻辑应当简单直接，避免过度工程化。

### III. YAGNI (You Aren't Gonna Need It)
只实现当前明确需要的功能，避免过度设计和为未来功能预留。每个新功能都必须有明确的用户需求支撑，避免基于猜测的"可能有用"的功能。

### IV. DRY (Don't Repeat Yourself)
识别并消除代码或逻辑中的重复模式，提高可重用性。模板系统、代码生成器和公共工具函数应当避免重复实现相同功能。

## Code Quality Standards

### Testing Discipline
- **TDD Mandatory**: 测试驱动开发必须严格执行，先写测试 → 用户确认 → 测试失败 → 然后实现
- **Red-Green-Refactor**: 严格遵守红-绿-重构循环
- **Integration Testing**: 需要集成测试的重点领域：新生成器合同测试、合同变更、服务间通信、共享模式

### CLI Interface Consistency
- **统一参数格式**: 所有命令遵循相同的参数命名约定（`--force`, `--dry-run`, `--dir`）
- **标准输出格式**: 支持 JSON 和人类可读格式，错误信息输出到 stderr
- **文本 I/O 协议**: stdin/args → stdout，错误 → stderr

### Error Handling Completeness
- **完整的错误信息**: 提供清晰的错误描述和建议解决方案
- **错误恢复机制**: 在可能的情况下提供恢复选项
- **错误日志记录**: 结构化日志记录，便于调试和监控

## Performance Requirements

### Generation Performance
- **代码生成速度**: 大型项目的代码生成必须在合理时间内完成（< 5秒）
- **内存使用**: 避免在生成过程中出现不必要的内存峰值
- **并行处理**: 在可能的情况下支持并行生成独立组件

### User Experience Consistency
- **响应性**: 所有命令操作必须提供即时反馈
- **进度指示**: 长时间运行的操作必须显示进度
- **一致性**: 所有命令的行为和输出格式保持一致

## Development Workflow

### Code Generation Process
1. **分析阶段**: 扫描项目结构，理解现有代码模式
2. **生成阶段**: 基于模板和配置生成代码
3. **验证阶段**: 确保生成代码的正确性和一致性
4. **格式化阶段**: 自动应用代码格式化标准

### Template Management
- **模板版本控制**: 所有模板必须版本化，确保可重现性
- **模板测试**: 每个模板必须有对应的测试用例
- **模板文档**: 复杂模板必须有详细的使用说明

## Governance

### Constitution Supremacy
本章程优先于所有其他实践规范。所有代码审查和 PR 都必须验证合规性。

### Amendment Process
- **修改提案**: 任何章程修改必须有明确的提案和理由
- **审查批准**: 修改需要经过核心团队审查和批准
- **迁移计划**: 重大变更必须有详细的迁移计划
- **文档更新**: 修改必须同步更新所有相关文档

### Versioning Policy
- **MAJOR**: 向后不兼容的治理原则删除或重新定义
- **MINOR**: 新原则或部分添加或实质性扩展指导
- **PATCH**: 澄清、措辞、拼写修正、非语义优化

### Compliance Review
- **定期审查**: 每季度审查章程的执行情况
- **违规处理**: 对违反章程的行为必须有明确的处理流程
- **持续改进**: 基于实际使用经验持续优化章程内容

**Version**: 1.0.0 | **Ratified**: 2025-09-19 | **Last Amended**: 2025-09-19