// Code generated by protoc-gen-go.
// source: item.proto
// DO NOT EDIT!

package msg

import proto "code.google.com/p/goprotobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type MQItemList struct {
	XXX_unrecognized []byte `json:"-"`
}

func (m *MQItemList) Reset()         { *m = MQItemList{} }
func (m *MQItemList) String() string { return proto.CompactTextString(m) }
func (*MQItemList) ProtoMessage()    {}

type MItem struct {
	Id               *uint32 `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Level            *uint32 `protobuf:"varint,2,req,name=level" json:"level,omitempty"`
	Data             *string `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *MItem) Reset()         { *m = MItem{} }
func (m *MItem) String() string { return proto.CompactTextString(m) }
func (*MItem) ProtoMessage()    {}

func (m *MItem) GetId() uint32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *MItem) GetLevel() uint32 {
	if m != nil && m.Level != nil {
		return *m.Level
	}
	return 0
}

func (m *MItem) GetData() string {
	if m != nil && m.Data != nil {
		return *m.Data
	}
	return ""
}

type MRItemList struct {
	ItemList         []*MItem `protobuf:"bytes,1,rep,name=item_list" json:"item_list,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *MRItemList) Reset()         { *m = MRItemList{} }
func (m *MRItemList) String() string { return proto.CompactTextString(m) }
func (*MRItemList) ProtoMessage()    {}

func (m *MRItemList) GetItemList() []*MItem {
	if m != nil {
		return m.ItemList
	}
	return nil
}

func init() {
}