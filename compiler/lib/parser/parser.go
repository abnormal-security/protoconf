package parser

import (
	"os"

	_ "github.com/abnormal-security/protoconf/pb/protoconf/v1"
	_ "github.com/bufbuild/protovalidate-go"
	_ "github.com/bufbuild/protovalidate-go/legacy"

	"github.com/abnormal-security/protoconf/utils"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// Parser provides a wrapper around jhump/protoreflect/protoparse that will keep a cache of dpd.FileDescriptor
type Parser struct {
	LocalResolver   *protoregistry.Types
	FilesResolver   *protoregistry.Files
	FileDescriptors map[string]*desc.FileDescriptor
}

func NewParserWithDescriptorRegistry(registry *utils.DescriptorRegistry) *Parser {
	files := registry.GetFilesResolver()
	return &Parser{
		FilesResolver:   files,
		LocalResolver:   registry.GetTypesResolver(files),
		FileDescriptors: registry.FileRegistry,
	}
}

func (p *Parser) ParseFilesX(filenames ...string) (results []*desc.FileDescriptor, err error) {
	for _, filename := range filenames {
		if fd, ok := p.FileDescriptors[filename]; ok {
			results = append(results, fd)
			continue
		}
		fd, err := p.FilesResolver.FindFileByPath(filename)
		if err != nil {
			return nil, err
		}
		d, err := desc.WrapFile(fd)
		if err != nil {
			f := protodesc.ToFileDescriptorProto(fd)
			deps := []*desc.FileDescriptor{}
			for _, dep := range f.Dependency {
				dd, err := desc.LoadFileDescriptor(dep)
				if err != nil {
					return nil, err
				}
				deps = append(deps, dd)
			}
			d, err = desc.CreateFileDescriptor(f, deps...)
			if err != nil {
				return nil, err
			}

			results = append(results, d)
		} else {
			results = append(results, d)
		}

	}
	return results, nil
}

func (p *Parser) ReadConfig(filename string, msg proto.Message) error {
	configReader, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return protojson.UnmarshalOptions{Resolver: p.LocalResolver}.Unmarshal(configReader, msg)
}
