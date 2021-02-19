package parser

import (
	"bytes"
	"fmt"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/desc/protoprint"
	_ "github.com/jhump/protoreflect/dynamic"
)

func main()  {
	relFilePaths := []string{
		"protocol/album.proto",
		"protocol/avatar.proto",
	}

	parser := protoparse.Parser{}

	descs, err := parser.ParseFiles(relFilePaths...)
	if err != nil {
		fmt.Printf("parse file from proto error: %v", err)
		return
	}

	printer := &protoprint.Printer{}

	var buf bytes.Buffer
	err = printer.PrintProtoFile(descs[0], &buf)
	if err != nil {
		fmt.Printf("print proto file error: %v", err)
		return
	}
	// 打印出所有的proto内容
	//fmt.Printf("descsStr=%s\n", buf.String())

	// descs name=./protocol/user.proto
	fmt.Printf("descs name=%s\n", descs[0].GetName())


	//fmt.Printf("descs string=%s\n", descs[0].String())

	// 获取AsFileDescriptorProto

	fmt.Printf("descriptor proto %v", descs[0].AsFileDescriptorProto())

	// 获取所有message

	// traverseMessage(descs[0].AsFileDescriptorProto())

}

func traverseMessage(des *dpb.FileDescriptorProto)  {
	for _, v := range des.GetMessageType() {
		fmt.Printf("message Name : %s Field: %v\n", v.GetName(), v.GetField())
	}
}

func GetMessageBin() []byte {
	//req := &User.GetAlbumPhotosReq{
	//	Uid: 12321,
	//}
	//bin, err := proto.Marshal(req)
	//if err != nil {
	//	fmt.Printf("bin=%v, err=%v", bin, err)
	//}
	return []byte{}
}




