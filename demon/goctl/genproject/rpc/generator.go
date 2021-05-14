package rpc

type GeneratorApi interface {
	Prepare() error
	GenMain() error
	GenServer() error
	GenPush() error
	GenUtil() error
	GenError() error
	GenRpcClient() error
	GenPb() error
}

type GeneratorServer interface {
	Prepare() error
	GenMain() error
	GenServer() error
	GenConfig() error
	GenStore() error
}

