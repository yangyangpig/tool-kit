package rpc

import (
	"gather/toolkitcl/demon/goctl/util"
	"gather/toolkitcl/demon/goctl/util/ctx"
	"path/filepath"
)

const (
	api = "api"
	server = "server"
	// store = "store"
)



type (
	DirContext interface {
		GetServer() Dir
	}

	// 用来定位到一个目录元数据
	Dir struct {
		Base string // 用于创建目录第一层的路径，绝对路径
		//FileName string // 创建文件名称
		ServerName string // 递归的当前目录
	}
	defaultDirContext struct {
		inner map[string]Dir
		serviceName string
	}
)

func mkdir(ctx *ctx.ProjectContext) (*defaultDirContext, error){
	// 根据proto的package最后的.的包名确定api的项目名称，一般包名和rpc定义的服务名称相同，不知道往后怎么变化，最后
	// 抽象成可配置
	inner := make(map[string]Dir)

	// 目前按照
	apiServerDir := filepath.Join(ctx.ProjectDir, ctx.PackageName)

	inner[api] = Dir{
		Base:apiServerDir,
		ServerName: ctx.PackageName,
	}

	for _, v := range inner {
		err := util.MkDirIfNotExist(v.Base)
		if err != nil {
			return nil, err
		}
	}

	return &defaultDirContext{inner: inner,serviceName:ctx.PackageName}, nil

}