package console

// GoZeroMakeRpcCommand 从proto文件生成完整的go-zero RPC服务命令
type GoZeroMakeRpcCommand struct {
	generator *GoZeroGenerator
}

// NewGoZeroMakeRpcCommand 创建新的从proto文件生成go-zero RPC服务命令
func NewGoZeroMakeRpcCommand(generator *GoZeroGenerator) *GoZeroMakeRpcCommand {
	return &GoZeroMakeRpcCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *GoZeroMakeRpcCommand) GetName() string {
	return "make:rpc"
}

// GetDescription 获取命令描述
func (cmd *GoZeroMakeRpcCommand) GetDescription() string {
	return "Generate complete go-zero RPC service from proto file"
}

// GetSignature 获取命令签名
func (cmd *GoZeroMakeRpcCommand) GetSignature() string {
	return "make:rpc <proto_file> [--output=]"
}

// GetArguments 获取命令参数
func (cmd *GoZeroMakeRpcCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "proto_file",
			Description: "The proto file path",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *GoZeroMakeRpcCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "output",
			ShortName:   "o",
			Description: "The output directory",
			Required:    false,
			Default:     ".",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *GoZeroMakeRpcCommand) Execute(input Input) error {
	protoFile := input.GetArgument("proto_file").(string)
	outputDir := input.GetOption("output").(string)

	return cmd.generator.GenerateFromProto(protoFile, outputDir)
}

// GoZeroMakeApiCommand 从.api文件生成完整的go-zero API服务命令
type GoZeroMakeApiCommand struct {
	generator *GoZeroGenerator
}

// NewGoZeroMakeApiCommand 创建新的从.api文件生成go-zero API服务命令
func NewGoZeroMakeApiCommand(generator *GoZeroGenerator) *GoZeroMakeApiCommand {
	return &GoZeroMakeApiCommand{
		generator: generator,
	}
}

// GetName 获取命令名称
func (cmd *GoZeroMakeApiCommand) GetName() string {
	return "make:api"
}

// GetDescription 获取命令描述
func (cmd *GoZeroMakeApiCommand) GetDescription() string {
	return "Generate complete go-zero API service from .api file"
}

// GetSignature 获取命令签名
func (cmd *GoZeroMakeApiCommand) GetSignature() string {
	return "make:api <api_file> [--output=]"
}

// GetArguments 获取命令参数
func (cmd *GoZeroMakeApiCommand) GetArguments() []Argument {
	return []Argument{
		{
			Name:        "api_file",
			Description: "The .api file path",
			Required:    true,
		},
	}
}

// GetOptions 获取命令选项
func (cmd *GoZeroMakeApiCommand) GetOptions() []Option {
	return []Option{
		{
			Name:        "output",
			ShortName:   "o",
			Description: "The output directory",
			Required:    false,
			Default:     ".",
			Type:        "string",
		},
	}
}

// Execute 执行命令
func (cmd *GoZeroMakeApiCommand) Execute(input Input) error {
	apiFile := input.GetArgument("api_file").(string)
	outputDir := input.GetOption("output").(string)

	return cmd.generator.GenerateFromApi(apiFile, outputDir)
}
