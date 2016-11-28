// Code generated by protoc-gen-go.
// source: Manage.proto
// DO NOT EDIT!

package Report

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type ManageProtocol_ManageType int32

const (
	ManageProtocol_REG    ManageProtocol_ManageType = 0
	ManageProtocol_LOGIN  ManageProtocol_ManageType = 1
	ManageProtocol_LOGOUT ManageProtocol_ManageType = 2
	ManageProtocol_CANCEL ManageProtocol_ManageType = 3
)

var ManageProtocol_ManageType_name = map[int32]string{
	0: "REG",
	1: "LOGIN",
	2: "LOGOUT",
	3: "CANCEL",
}
var ManageProtocol_ManageType_value = map[string]int32{
	"REG":    0,
	"LOGIN":  1,
	"LOGOUT": 2,
	"CANCEL": 3,
}

func (x ManageProtocol_ManageType) String() string {
	return proto.EnumName(ManageProtocol_ManageType_name, int32(x))
}
func (ManageProtocol_ManageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0, 0} }

type ManageProtocol_ResultType int32

const (
	ManageProtocol_SUCCESS   ManageProtocol_ResultType = 0
	ManageProtocol_ARREARAGE ManageProtocol_ResultType = 1
	ManageProtocol_INVALID   ManageProtocol_ResultType = 2
)

var ManageProtocol_ResultType_name = map[int32]string{
	0: "SUCCESS",
	1: "ARREARAGE",
	2: "INVALID",
}
var ManageProtocol_ResultType_value = map[string]int32{
	"SUCCESS":   0,
	"ARREARAGE": 1,
	"INVALID":   2,
}

func (x ManageProtocol_ResultType) String() string {
	return proto.EnumName(ManageProtocol_ResultType_name, int32(x))
}
func (ManageProtocol_ResultType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0, 1} }

type ManageProtocol struct {
	TimeReq      string                    `protobuf:"bytes,1,opt,name=time_req,json=timeReq" json:"time_req,omitempty"`
	SerialNumber string                    `protobuf:"bytes,2,opt,name=serial_number,json=serialNumber" json:"serial_number,omitempty"`
	Tid          string                    `protobuf:"bytes,3,opt,name=tid" json:"tid,omitempty"`
	Type         ManageProtocol_ManageType `protobuf:"varint,4,opt,name=type,enum=Report.ManageProtocol_ManageType" json:"type,omitempty"`
	Result       ManageProtocol_ResultType `protobuf:"varint,5,opt,name=result,enum=Report.ManageProtocol_ResultType" json:"result,omitempty"`
	TerminalType string                    `protobuf:"bytes,6,opt,name=terminalType" json:"terminalType,omitempty"`
	ProtocolType string                    `protobuf:"bytes,7,opt,name=protocolType" json:"protocolType,omitempty"`
}

func (m *ManageProtocol) Reset()                    { *m = ManageProtocol{} }
func (m *ManageProtocol) String() string            { return proto.CompactTextString(m) }
func (*ManageProtocol) ProtoMessage()               {}
func (*ManageProtocol) Descriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

func init() {
	proto.RegisterType((*ManageProtocol)(nil), "Report.ManageProtocol")
	proto.RegisterEnum("Report.ManageProtocol_ManageType", ManageProtocol_ManageType_name, ManageProtocol_ManageType_value)
	proto.RegisterEnum("Report.ManageProtocol_ResultType", ManageProtocol_ResultType_name, ManageProtocol_ResultType_value)
}

func init() { proto.RegisterFile("Manage.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 288 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x7c, 0x8e, 0x4f, 0x4f, 0x83, 0x40,
	0x10, 0xc5, 0x6d, 0x69, 0xc1, 0x8e, 0xb4, 0xd9, 0xcc, 0x69, 0xbd, 0x29, 0x5e, 0x3c, 0x71, 0xd0,
	0x34, 0xd1, 0x23, 0x41, 0x42, 0x48, 0x90, 0x9a, 0xa5, 0xf5, 0xda, 0x50, 0xdd, 0x18, 0x12, 0xfe,
	0x75, 0xbb, 0x3d, 0xf8, 0x21, 0xfc, 0xce, 0xee, 0x2e, 0x24, 0xad, 0x17, 0x6f, 0x33, 0xbf, 0xf7,
	0xde, 0xcc, 0x03, 0xf7, 0xb5, 0x68, 0x8a, 0x2f, 0xee, 0x77, 0xa2, 0x95, 0x2d, 0xda, 0x8c, 0x77,
	0xad, 0x90, 0xde, 0x8f, 0x05, 0x8b, 0x5e, 0x78, 0xd3, 0xfc, 0xa3, 0xad, 0xf0, 0x1a, 0x2e, 0x65,
	0x59, 0xf3, 0xad, 0xe0, 0x7b, 0x3a, 0xba, 0x19, 0xdd, 0xcf, 0x98, 0xa3, 0x77, 0xc6, 0xf7, 0x78,
	0x07, 0xf3, 0x03, 0x17, 0x65, 0x51, 0x6d, 0x9b, 0x63, 0xbd, 0xe3, 0x82, 0x8e, 0x8d, 0xee, 0xf6,
	0x30, 0x33, 0x0c, 0x09, 0x58, 0xb2, 0xfc, 0xa4, 0x96, 0x91, 0xf4, 0x88, 0x4b, 0x98, 0xc8, 0xef,
	0x8e, 0xd3, 0x89, 0x42, 0x8b, 0x87, 0x5b, 0xbf, 0xff, 0xed, 0xff, 0xfd, 0x3b, 0xac, 0x6b, 0x65,
	0x64, 0xc6, 0x8e, 0xcf, 0x60, 0x0b, 0x7e, 0x38, 0x56, 0x92, 0x4e, 0xff, 0x0d, 0x32, 0x63, 0x32,
	0xc1, 0x21, 0x80, 0x1e, 0xb8, 0x92, 0x8b, 0xba, 0x6c, 0x8a, 0x4a, 0x73, 0x6a, 0xf7, 0x3d, 0xcf,
	0x99, 0xf6, 0x74, 0xc3, 0x09, 0xe3, 0x71, 0x7a, 0xcf, 0x39, 0xf3, 0x9e, 0x00, 0x4e, 0xb5, 0xd0,
	0x01, 0x8b, 0x45, 0x31, 0xb9, 0xc0, 0x19, 0x4c, 0xd3, 0x55, 0x9c, 0x64, 0x64, 0x84, 0x00, 0xb6,
	0x1a, 0x57, 0x9b, 0x35, 0x19, 0xeb, 0x39, 0x0c, 0xb2, 0x30, 0x4a, 0x89, 0xe5, 0x2d, 0x01, 0x4e,
	0xbd, 0xf0, 0x0a, 0x9c, 0x7c, 0x13, 0x86, 0x51, 0x9e, 0xab, 0xf4, 0x1c, 0x66, 0x01, 0x63, 0x51,
	0xc0, 0x82, 0x38, 0x52, 0x17, 0x94, 0x96, 0x64, 0xef, 0x41, 0x9a, 0xbc, 0x90, 0xf1, 0xce, 0x36,
	0xef, 0x1f, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xc6, 0x01, 0x90, 0x16, 0xae, 0x01, 0x00, 0x00,
}
