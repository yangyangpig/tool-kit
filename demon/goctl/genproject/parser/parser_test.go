package parser

import (
	"fmt"
	"testing"
)

func TestDefaultProtoParser_Parse(t *testing.T) {
	relFilePaths := []string{
		"protocol/album.proto",
		"protocol/avatar.proto",
	}
	defaultParse,err := NewDefaultProtoParser(relFilePaths)
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
