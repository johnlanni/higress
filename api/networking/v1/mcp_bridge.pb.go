// Copyright (c) 2022 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: networking/v1/mcp_bridge.proto

// $schema: higress.io.v1.McpBridge
// $title: Mcp Bridge
// $description: Configuration affecting service discovery from multi registries
// $mode: none

package v1

import (
	duration "github.com/golang/protobuf/ptypes/duration"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// <!-- crd generation tags
// +cue-gen:McpBridge:groupName:networking.higress.io
// +cue-gen:McpBridge:version:v1
// +cue-gen:McpBridge:storageVersion
// +cue-gen:McpBridge:annotations:helm.sh/resource-policy=keep
// +cue-gen:McpBridge:subresource:status
// +cue-gen:McpBridge:scope:Namespaced
// +cue-gen:McpBridge:resource:categories=higress-io,plural=mcpbridges
// +cue-gen:McpBridge:preserveUnknownFields:false
// -->
//
// <!-- go code generation tags
// +kubetype-gen
// +kubetype-gen:groupVersion=networking.higress.io/v1
// +genclient
// +k8s:deepcopy-gen=true
// -->
type McpBridge struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Registries []*RegistryConfig `protobuf:"bytes,1,rep,name=registries,proto3" json:"registries,omitempty"`
}

func (x *McpBridge) Reset() {
	*x = McpBridge{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_v1_mcp_bridge_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *McpBridge) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*McpBridge) ProtoMessage() {}

func (x *McpBridge) ProtoReflect() protoreflect.Message {
	mi := &file_networking_v1_mcp_bridge_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use McpBridge.ProtoReflect.Descriptor instead.
func (*McpBridge) Descriptor() ([]byte, []int) {
	return file_networking_v1_mcp_bridge_proto_rawDescGZIP(), []int{0}
}

func (x *McpBridge) GetRegistries() []*RegistryConfig {
	if x != nil {
		return x.Registries
	}
	return nil
}

type RegistryConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type                 string             `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Name                 string             `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Domain               string             `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Port                 uint32             `protobuf:"varint,4,opt,name=port,proto3" json:"port,omitempty"`
	NacosAddressServer   string             `protobuf:"bytes,5,opt,name=nacosAddressServer,proto3" json:"nacosAddressServer,omitempty"`
	NacosAccessKey       string             `protobuf:"bytes,6,opt,name=nacosAccessKey,proto3" json:"nacosAccessKey,omitempty"`
	NacosSecretKey       string             `protobuf:"bytes,7,opt,name=nacosSecretKey,proto3" json:"nacosSecretKey,omitempty"`
	NacosNamespaceId     string             `protobuf:"bytes,8,opt,name=nacosNamespaceId,proto3" json:"nacosNamespaceId,omitempty"`
	NacosNamespace       string             `protobuf:"bytes,9,opt,name=nacosNamespace,proto3" json:"nacosNamespace,omitempty"`
	NacosGroups          []string           `protobuf:"bytes,10,rep,name=nacosGroups,proto3" json:"nacosGroups,omitempty"`
	NacosRefreshInterval *duration.Duration `protobuf:"bytes,11,opt,name=nacosRefreshInterval,proto3" json:"nacosRefreshInterval,omitempty"`
	ConsulNamespace      string             `protobuf:"bytes,12,opt,name=consulNamespace,proto3" json:"consulNamespace,omitempty"`
	ZkServicesPath       []string           `protobuf:"bytes,13,rep,name=zkServicesPath,proto3" json:"zkServicesPath,omitempty"`
}

func (x *RegistryConfig) Reset() {
	*x = RegistryConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_networking_v1_mcp_bridge_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegistryConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegistryConfig) ProtoMessage() {}

func (x *RegistryConfig) ProtoReflect() protoreflect.Message {
	mi := &file_networking_v1_mcp_bridge_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegistryConfig.ProtoReflect.Descriptor instead.
func (*RegistryConfig) Descriptor() ([]byte, []int) {
	return file_networking_v1_mcp_bridge_proto_rawDescGZIP(), []int{1}
}

func (x *RegistryConfig) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *RegistryConfig) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *RegistryConfig) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *RegistryConfig) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *RegistryConfig) GetNacosAddressServer() string {
	if x != nil {
		return x.NacosAddressServer
	}
	return ""
}

func (x *RegistryConfig) GetNacosAccessKey() string {
	if x != nil {
		return x.NacosAccessKey
	}
	return ""
}

func (x *RegistryConfig) GetNacosSecretKey() string {
	if x != nil {
		return x.NacosSecretKey
	}
	return ""
}

func (x *RegistryConfig) GetNacosNamespaceId() string {
	if x != nil {
		return x.NacosNamespaceId
	}
	return ""
}

func (x *RegistryConfig) GetNacosNamespace() string {
	if x != nil {
		return x.NacosNamespace
	}
	return ""
}

func (x *RegistryConfig) GetNacosGroups() []string {
	if x != nil {
		return x.NacosGroups
	}
	return nil
}

func (x *RegistryConfig) GetNacosRefreshInterval() *duration.Duration {
	if x != nil {
		return x.NacosRefreshInterval
	}
	return nil
}

func (x *RegistryConfig) GetConsulNamespace() string {
	if x != nil {
		return x.ConsulNamespace
	}
	return ""
}

func (x *RegistryConfig) GetZkServicesPath() []string {
	if x != nil {
		return x.ZkServicesPath
	}
	return nil
}

var File_networking_v1_mcp_bridge_proto protoreflect.FileDescriptor

var file_networking_v1_mcp_bridge_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x2f,
	0x6d, 0x63, 0x70, 0x5f, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0d, 0x68, 0x69, 0x67, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x76, 0x31, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x4a, 0x0a, 0x09, 0x4d, 0x63, 0x70, 0x42, 0x72, 0x69, 0x64, 0x67, 0x65, 0x12, 0x3d, 0x0a,
	0x0a, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1d, 0x2e, 0x68, 0x69, 0x67, 0x72, 0x65, 0x73, 0x73, 0x2e, 0x69, 0x6f, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x52, 0x0a, 0x72, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x69, 0x65, 0x73, 0x22, 0x8d, 0x04, 0x0a,
	0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x72, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x18, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xe2,
	0x41, 0x01, 0x02, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0xe2,
	0x41, 0x01, 0x02, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x18, 0x0a, 0x04, 0x70,
	0x6f, 0x72, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x42, 0x04, 0xe2, 0x41, 0x01, 0x02, 0x52,
	0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x2e, 0x0a, 0x12, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x41, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x12, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x12, 0x26, 0x0a, 0x0e, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x41, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e,
	0x61, 0x63, 0x6f, 0x73, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x4b, 0x65, 0x79, 0x12, 0x26, 0x0a,
	0x0e, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x53, 0x65, 0x63, 0x72,
	0x65, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x2a, 0x0a, 0x10, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x4e, 0x61,
	0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x10, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x49,
	0x64, 0x12, 0x26, 0x0a, 0x0e, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6e, 0x61, 0x63, 0x6f, 0x73,
	0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x6e, 0x61, 0x63,
	0x6f, 0x73, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b,
	0x6e, 0x61, 0x63, 0x6f, 0x73, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x73, 0x12, 0x4d, 0x0a, 0x14, 0x6e,
	0x61, 0x63, 0x6f, 0x73, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x76, 0x61, 0x6c, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x52, 0x14, 0x6e, 0x61, 0x63, 0x6f, 0x73, 0x52, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x6f,
	0x6e, 0x73, 0x75, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x70, 0x61, 0x63, 0x65, 0x12, 0x26, 0x0a, 0x0e, 0x7a, 0x6b, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x73, 0x50, 0x61, 0x74, 0x68, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x7a, 0x6b,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x50, 0x61, 0x74, 0x68, 0x42, 0x2e, 0x5a, 0x2c,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x69, 0x62, 0x61,
	0x62, 0x61, 0x2f, 0x68, 0x69, 0x67, 0x72, 0x65, 0x73, 0x73, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_networking_v1_mcp_bridge_proto_rawDescOnce sync.Once
	file_networking_v1_mcp_bridge_proto_rawDescData = file_networking_v1_mcp_bridge_proto_rawDesc
)

func file_networking_v1_mcp_bridge_proto_rawDescGZIP() []byte {
	file_networking_v1_mcp_bridge_proto_rawDescOnce.Do(func() {
		file_networking_v1_mcp_bridge_proto_rawDescData = protoimpl.X.CompressGZIP(file_networking_v1_mcp_bridge_proto_rawDescData)
	})
	return file_networking_v1_mcp_bridge_proto_rawDescData
}

var file_networking_v1_mcp_bridge_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_networking_v1_mcp_bridge_proto_goTypes = []interface{}{
	(*McpBridge)(nil),         // 0: higress.io.v1.McpBridge
	(*RegistryConfig)(nil),    // 1: higress.io.v1.RegistryConfig
	(*duration.Duration)(nil), // 2: google.protobuf.Duration
}
var file_networking_v1_mcp_bridge_proto_depIdxs = []int32{
	1, // 0: higress.io.v1.McpBridge.registries:type_name -> higress.io.v1.RegistryConfig
	2, // 1: higress.io.v1.RegistryConfig.nacosRefreshInterval:type_name -> google.protobuf.Duration
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_networking_v1_mcp_bridge_proto_init() }
func file_networking_v1_mcp_bridge_proto_init() {
	if File_networking_v1_mcp_bridge_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_networking_v1_mcp_bridge_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*McpBridge); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_networking_v1_mcp_bridge_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegistryConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_networking_v1_mcp_bridge_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_networking_v1_mcp_bridge_proto_goTypes,
		DependencyIndexes: file_networking_v1_mcp_bridge_proto_depIdxs,
		MessageInfos:      file_networking_v1_mcp_bridge_proto_msgTypes,
	}.Build()
	File_networking_v1_mcp_bridge_proto = out.File
	file_networking_v1_mcp_bridge_proto_rawDesc = nil
	file_networking_v1_mcp_bridge_proto_goTypes = nil
	file_networking_v1_mcp_bridge_proto_depIdxs = nil
}
