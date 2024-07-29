// For more information about the usage of this plugin, see:
// https://protobuf.dev/reference/go/go-generated.
package main

import (
	"flag"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {

	var (
		flags flag.FlagSet
		//plugins = flags.String("plugins", "", "deprecated option")
		//expose_all = flags.Bool("expose_all", false, "test protoc-gen-debug bool parameter")
		//foo        = flags.String("foo", "", "test protoc-gen-debug string parameter")
	)
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		//if *plugins != "" {
		//	return errors.New("protoc-gen-go: plugins are not supported; use 'protoc --go-grpc_out=...' to generate gRPC\n\n" +
		//		"See " + grpcDocURL + " for more information.")
		//}

		//fmt.Println("expose_all", *expose_all)
		//fmt.Println("foo", *foo)

		q := &Queue{}
		q.Generate(gen)

		return nil
	})
}

type Queue struct{}

func (q *Queue) Generate(gen *protogen.Plugin) error {
	for _, f := range gen.Files {
		if !f.Generate {
			continue
		}
		q.genCode(gen, f)
	}

	gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL) | uint64(pluginpb.CodeGeneratorResponse_FEATURE_SUPPORTS_EDITIONS)
	gen.SupportedEditionsMinimum = descriptorpb.Edition_EDITION_PROTO3
	gen.SupportedEditionsMaximum = descriptorpb.Edition_EDITION_2024

	return nil
}

func (q *Queue) genCode(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	filename := file.GeneratedFilenamePrefix + ".queue.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// source: ", file.Desc.Path())
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()
	g.P("// 测试文件 " + file.GoPackageName)
	g.P("import (")
	g.P("\"encoding/json\"")
	g.P("\"reflect\"")
	g.P()
	g.P("msg \"protobuf-plugin-debug/internal/messager\"")
	g.P(")")

	for _, f := range gen.Files {
		for _, message := range f.Proto.GetMessageType() {
			msgName := *message.Name
			titleName := *message.Name

			g.P("type ", msgName, " struct {")
			g.P("}")
			g.P("func (a *", msgName, ") Encode() (*msg.Package, error) {")
			g.P("\tprType := reflect.TypeOf(a)\n\n\tdata, err := json.Marshal(a)\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\n\treturn &msg.Package{\n\t\tData: data,\n\t\tMeta: prType,\n\t}, nil\n}")
			g.P()
			g.P("func (a *", msgName, ") Decode(data []byte) error {")
			g.P("\treturn json.Unmarshal(data, a)\n}")
			g.P()
			g.P("func (a *", msgName, ") Title() string {")
			g.P("\treturn \"", titleName, "\"")
			g.P("}")
			g.P()
		}
	}

	return g
}
