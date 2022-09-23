package common

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"io/ioutil"
)

func LoadProtobuf(filename string) (error, []protoreflect.FileDescriptor) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err, nil
	}

	fds := &descriptorpb.FileDescriptorSet{}
	err = proto.Unmarshal(b, fds)
	if err != nil {
		return err, nil
	}

	ff, err := protodesc.NewFiles(fds)
	if err != nil {
		return err, nil
	}

	var ret []protoreflect.FileDescriptor
	ff.RangeFiles(func(descriptor protoreflect.FileDescriptor) bool {
		ret = append(ret, descriptor)
		return true
	})

	return nil, ret
}

func LoadProtobufMethods(filename string) (error, []protoreflect.MethodDescriptor) {
	err, descriptors := LoadProtobuf(filename)
	if err != nil {
		return err, nil
	}

	var ret []protoreflect.MethodDescriptor
	for _, descriptor := range descriptors {
		for i := 0; i < descriptor.Services().Len(); i++ {
			sd := descriptor.Services().Get(i)
			for j := 0; j < sd.Methods().Len(); j++ {
				m := sd.Methods().Get(j)
				ret = append(ret, m)
			}
		}
	}

	return nil, ret
}

func MessageToFullJson(mi protoreflect.MessageDescriptor) (error, string) {
	message := dynamicpb.NewMessage(mi)
	fullFill(message, mi)
	b, err := protojson.MarshalOptions{Multiline: true, Indent: "  ", EmitUnpopulated: true}.Marshal(message)
	if err != nil {
		return err, ""
	}
	return nil, string(b)
}

func fullFill(message *dynamicpb.Message, mi protoreflect.MessageDescriptor) {
	for i := 0; i < mi.Fields().Len(); i++ {
		fd := mi.Fields().Get(i)
		if !fd.IsMap() && !fd.IsList() && fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
			submi := fd.Message()
			submessage := dynamicpb.NewMessage(submi)
			fullFill(submessage, submi)
			message.Set(fd, protoreflect.ValueOfMessage(submessage))
		}
	}
}
