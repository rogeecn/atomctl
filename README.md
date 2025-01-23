# Atomctl 脚手架工具命令指南

## 命令列表

### 生成命令

#### gen model (别名：m)
- 描述：生成jet模型
- 功能：
  - 从PostgreSQL数据库生成模型
  - 使用database/transform.yaml配置文件进行类型转换
  - 支持忽略特定表和枚举
  - 生成JSON标签
  - 支持自定义字段类型映射
  - 自动生成数据库schema文件

#### gen provider (别名：p)
- 描述：生成provider
- 参数：
  - path：可选，指定生成路径（默认当前目录）
- 功能：
  - 解析指定目录下的.go文件
  - 查找带有@provider注释的结构体
  - 支持 `@provider(grpc|event|job):[except|only] [returnType] [group]` 注释
  - 自动生成provider文件
  - 支持分组生成

#### gen route
- 描述：生成路由
- 参数：
  - path：可选，指定生成路径（默认当前目录）
- 功能：
  - 解析app/http目录下的controller文件
  - 自动生成路由定义
  - 支持分组生成路由
  - 生成完成后自动执行gen provider命令

### 数据库命令

#### migrate (别名：m)
- 描述：数据库迁移
- 参数：
  - action：必选，迁移操作（up|up-by-one|up-to|create|down|down-to|fix|redo|reset|status|version）
  - args：可选，操作参数
- 选项：
  - -c/--config：指定数据库配置文件（默认config.toml）
- 功能：
  - 执行数据库迁移
  - 支持创建迁移文件
  - 支持回滚、重置等操作
  - 查看迁移状态和版本

### 新建命令

#### new project (别名：p)
- 描述：创建新项目
- 参数：
  - moduleName：必选，项目模块名（需符合Go包名规范）
- 选项：
  - --force：强制覆盖已存在项目
- 功能：
  - 根据模板生成项目结构
  - 自动处理隐藏文件（将模板中的-前缀转换为.）
  - 支持模板渲染
  - 生成完成后提示后续步骤

#### new provider
- 描述：创建新的provider
- 参数：
  - providerName：必选，provider名称
- 功能：
  - 在providers目录下创建新的provider
  - 自动生成provider模板文件
  - 支持模板渲染
  - 自动处理命名转换（如驼峰命名）
