package ctx

import (
	"fmt"
	"github.com/jhump/protoreflect/desc/protoparse"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var projectName = "toolkit"

type (
	// 主要用于透传项目公共常用的数据，包括proto的路径，项目路径等
	ProjectContext struct {
		ProjectDir string // 输入的路径

		// 服务名，和proto定义的服务名对应
		Name string

		// 指定的proto文件路径
		ProtoPath string

		// proto协议指定的包名
		PackageName string
	}
)
// 所有的解析proto的逻辑，先确定访问在prepare这里，通过ctx带到所有往下面的方法中去
func Prepare(protoDirPath string) (*ProjectContext, error) {
	// TODO 这里可以扩展加入检测goMod相关逻辑
	pCtx := &ProjectContext{
		ProtoPath: protoDirPath,
	}
	absPath, err := filepath.Abs(protoDirPath)
	if err != nil {
		log.Fatalf("get proto file abs path error: %v", err)
		return nil, err
	}
	fmt.Printf("proto abs path : %s\n", absPath)
	// 根据proto文件目录，获取该目录下所有文件路径
	fileInfo, err := ioutil.ReadDir(absPath)
	if err != nil {
		log.Fatalf("get proto file list error: %v", err)
		return nil, err
	}
	realFilePaths := make([]string, 0, len(fileInfo))

	for _, v := range fileInfo {
		if !strings.EqualFold(path.Ext(v.Name()), ".proto") {
			continue
		}
		realFilePaths = append(realFilePaths, absPath + "/" + v.Name())
	}
	os.Chdir(absPath)
	fmt.Printf("proto real file path %v\n", realFilePaths)
	parser := protoparse.Parser{}
	fileDescriptors, err := parser.ParseFiles(realFilePaths...)
	if err != nil {
		fmt.Printf("parse file from proto error: %v", err)
		return nil, err
	}

	// 获取服务名
	for _, v := range fileDescriptors {
		if v.AsFileDescriptorProto().GetService() != nil {
			for _, s := range v.AsFileDescriptorProto().GetService() {
				if s.GetName() != "" {
					fmt.Printf("prc server name : %s\n", s.GetName())
					pCtx.Name = s.GetName()
				}
			}
		}
	}

	// 获取项目目录路径
	pCtx.ProjectDir = getProjectDir(protoDirPath)
	return pCtx, nil
}

func getProjectDir(protoDirPath string) (projectDir string) {
	if strings.Contains(protoDirPath, projectName) {
		resp := strings.SplitN(protoDirPath, projectName, 2)
		if len(resp) >1 {
			projectDir = resp[0]
		}
	}
	return
}
