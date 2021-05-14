package rpc

import (
	"fmt"
	"gather/toolkitcl/demon/goctl/genproject/parser"
	"testing"
)

func TestMkdir_mkdir(t *testing.T)  {
	// TODO import的路径和proto文件路径怎么样区分开
	relFilePaths := []string{
		"../parser/protocol/album.proto",
		"../parser/protocol/avatar.proto",
	}
	defaultParse,err := parser.NewDefaultProtoParser(relFilePaths)
	if err != nil {
		t.Errorf("new default proto parser error: %v", err)
		return
	}
	resp, err := defaultParse.Parse()
	if err != nil {
		t.Errorf("parse proto error: %v", err)
		return
	}
	fmt.Printf("parse file descriptor: %v\n", resp)
	avatarFileDescriptor, err := defaultParse.GetFileDescriptorByProtoName("avatar")
	if err != nil {
		t.Errorf("get point proto descriptor error: %v", err)
		return
	}
	fmt.Printf("avatar descriptor: %v", avatarFileDescriptor)

}
