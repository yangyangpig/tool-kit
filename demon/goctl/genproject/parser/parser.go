package parser

import (
	"errors"
	"fmt"
	"gather/toolkitcl/demon/goctl/genproject/rpc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"strings"

	"github.com/jhump/protoreflect/desc"
)

type (
	defaultProtoParser struct {
		FilePaths []string
	}

)

var (
	errIllegalPath = errors.New("the proto path is illegal")

)
// path是指定到proto文件
func NewDefaultProtoParser(paths []string) (*defaultProtoParser, error)  {
	// TODO Check the file path validity
	for _, v := range paths {
		if !strings.Contains(v, ".proto") {
			return nil, errIllegalPath
		}
	}

	return &defaultProtoParser{FilePaths: paths}, nil
}

func (p *defaultProtoParser) Parse() (map[string]*desc.FileDescriptor, error) {
	return p.parse()
}

func (p *defaultProtoParser) parse() (map[string]*desc.FileDescriptor, error) {
	parser := protoparse.Parser{}
	fileDescriptors, err := parser.ParseFiles(p.FilePaths...)
	if err != nil {
		fmt.Printf("parse file from proto error: %v", err)
		return nil, err
	}
	fileDescriptorMap := make(map[string]*desc.FileDescriptor)
	for i, v := range p.FilePaths {
		protoFileName := rpc.GetProtoName(v)
		fmt.Printf("proto file name %s\n", protoFileName)
		fileDescriptorMap[protoFileName] = fileDescriptors[i]
	}
	return fileDescriptorMap, nil
}


func (p *defaultProtoParser) GetFileDescriptorByProtoName(protoName string) (
	res *desc.FileDescriptor, err error) {
	fileDescriptorMap, err := p.parse()
	if err != nil {
		return nil, fmt.Errorf("parse proto error: %v", err)
	}
	if fileDescriptorMap[protoName] != nil {
		return fileDescriptorMap[protoName], nil
	}
	return
}

