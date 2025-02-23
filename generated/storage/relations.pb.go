// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/relations.proto

package storage

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	types "github.com/gogo/protobuf/types"
	proto "github.com/golang/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ImageComponentEdge struct {
	// id is base 64 encoded Image:Component ids.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" sql:"pk,id"`
	/// Layer that contains this component
	//
	// Types that are valid to be assigned to HasLayerIndex:
	//	*ImageComponentEdge_LayerIndex
	HasLayerIndex        isImageComponentEdge_HasLayerIndex `protobuf_oneof:"has_layer_index"`
	Location             string                             `protobuf:"bytes,3,opt,name=location,proto3" json:"location,omitempty" search:"Component Location,store,hidden"`
	ImageId              string                             `protobuf:"bytes,4,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty" sql:"pk,fk(Image:id)"`
	ImageComponentId     string                             `protobuf:"bytes,5,opt,name=image_component_id,json=imageComponentId,proto3" json:"image_component_id,omitempty" sql:"pk,fk(ImageComponent:id),no-fk-constraint"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_unrecognized     []byte                             `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *ImageComponentEdge) Reset()         { *m = ImageComponentEdge{} }
func (m *ImageComponentEdge) String() string { return proto.CompactTextString(m) }
func (*ImageComponentEdge) ProtoMessage()    {}
func (*ImageComponentEdge) Descriptor() ([]byte, []int) {
	return fileDescriptor_62f882e266fcf764, []int{0}
}
func (m *ImageComponentEdge) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ImageComponentEdge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ImageComponentEdge.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ImageComponentEdge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageComponentEdge.Merge(m, src)
}
func (m *ImageComponentEdge) XXX_Size() int {
	return m.Size()
}
func (m *ImageComponentEdge) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageComponentEdge.DiscardUnknown(m)
}

var xxx_messageInfo_ImageComponentEdge proto.InternalMessageInfo

type isImageComponentEdge_HasLayerIndex interface {
	isImageComponentEdge_HasLayerIndex()
	MarshalTo([]byte) (int, error)
	Size() int
	Clone() isImageComponentEdge_HasLayerIndex
}

type ImageComponentEdge_LayerIndex struct {
	LayerIndex int32 `protobuf:"varint,2,opt,name=layer_index,json=layerIndex,proto3,oneof" json:"layer_index,omitempty"`
}

func (*ImageComponentEdge_LayerIndex) isImageComponentEdge_HasLayerIndex() {}
func (m *ImageComponentEdge_LayerIndex) Clone() isImageComponentEdge_HasLayerIndex {
	if m == nil {
		return nil
	}
	cloned := new(ImageComponentEdge_LayerIndex)
	*cloned = *m

	return cloned
}

func (m *ImageComponentEdge) GetHasLayerIndex() isImageComponentEdge_HasLayerIndex {
	if m != nil {
		return m.HasLayerIndex
	}
	return nil
}

func (m *ImageComponentEdge) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ImageComponentEdge) GetLayerIndex() int32 {
	if x, ok := m.GetHasLayerIndex().(*ImageComponentEdge_LayerIndex); ok {
		return x.LayerIndex
	}
	return 0
}

func (m *ImageComponentEdge) GetLocation() string {
	if m != nil {
		return m.Location
	}
	return ""
}

func (m *ImageComponentEdge) GetImageId() string {
	if m != nil {
		return m.ImageId
	}
	return ""
}

func (m *ImageComponentEdge) GetImageComponentId() string {
	if m != nil {
		return m.ImageComponentId
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ImageComponentEdge) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ImageComponentEdge_LayerIndex)(nil),
	}
}

func (m *ImageComponentEdge) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ImageComponentEdge) Clone() *ImageComponentEdge {
	if m == nil {
		return nil
	}
	cloned := new(ImageComponentEdge)
	*cloned = *m

	if m.HasLayerIndex != nil {
		cloned.HasLayerIndex = m.HasLayerIndex.Clone()
	}
	return cloned
}

type ComponentCVEEdge struct {
	// base 64 encoded Component:CVE ids.
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" sql:"pk,id"`
	IsFixable bool   `protobuf:"varint,2,opt,name=is_fixable,json=isFixable,proto3" json:"is_fixable,omitempty" search:"Fixable,store"`
	// Whether there is a version the CVE is fixed in the component.
	//
	// Types that are valid to be assigned to HasFixedBy:
	//	*ComponentCVEEdge_FixedBy
	HasFixedBy           isComponentCVEEdge_HasFixedBy `protobuf_oneof:"has_fixed_by"`
	ImageComponentId     string                        `protobuf:"bytes,4,opt,name=image_component_id,json=imageComponentId,proto3" json:"image_component_id,omitempty" sql:"pk,fk(ImageComponent:id)"`
	ImageCveId           string                        `protobuf:"bytes,5,opt,name=image_cve_id,json=imageCveId,proto3" json:"image_cve_id,omitempty" sql:"pk,fk(CVE:id),no-fk-constraint"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *ComponentCVEEdge) Reset()         { *m = ComponentCVEEdge{} }
func (m *ComponentCVEEdge) String() string { return proto.CompactTextString(m) }
func (*ComponentCVEEdge) ProtoMessage()    {}
func (*ComponentCVEEdge) Descriptor() ([]byte, []int) {
	return fileDescriptor_62f882e266fcf764, []int{1}
}
func (m *ComponentCVEEdge) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ComponentCVEEdge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ComponentCVEEdge.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ComponentCVEEdge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ComponentCVEEdge.Merge(m, src)
}
func (m *ComponentCVEEdge) XXX_Size() int {
	return m.Size()
}
func (m *ComponentCVEEdge) XXX_DiscardUnknown() {
	xxx_messageInfo_ComponentCVEEdge.DiscardUnknown(m)
}

var xxx_messageInfo_ComponentCVEEdge proto.InternalMessageInfo

type isComponentCVEEdge_HasFixedBy interface {
	isComponentCVEEdge_HasFixedBy()
	MarshalTo([]byte) (int, error)
	Size() int
	Clone() isComponentCVEEdge_HasFixedBy
}

type ComponentCVEEdge_FixedBy struct {
	FixedBy string `protobuf:"bytes,3,opt,name=fixed_by,json=fixedBy,proto3,oneof" json:"fixed_by,omitempty" search:"Fixed By,store,hidden"`
}

func (*ComponentCVEEdge_FixedBy) isComponentCVEEdge_HasFixedBy() {}
func (m *ComponentCVEEdge_FixedBy) Clone() isComponentCVEEdge_HasFixedBy {
	if m == nil {
		return nil
	}
	cloned := new(ComponentCVEEdge_FixedBy)
	*cloned = *m

	return cloned
}

func (m *ComponentCVEEdge) GetHasFixedBy() isComponentCVEEdge_HasFixedBy {
	if m != nil {
		return m.HasFixedBy
	}
	return nil
}

func (m *ComponentCVEEdge) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ComponentCVEEdge) GetIsFixable() bool {
	if m != nil {
		return m.IsFixable
	}
	return false
}

func (m *ComponentCVEEdge) GetFixedBy() string {
	if x, ok := m.GetHasFixedBy().(*ComponentCVEEdge_FixedBy); ok {
		return x.FixedBy
	}
	return ""
}

func (m *ComponentCVEEdge) GetImageComponentId() string {
	if m != nil {
		return m.ImageComponentId
	}
	return ""
}

func (m *ComponentCVEEdge) GetImageCveId() string {
	if m != nil {
		return m.ImageCveId
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*ComponentCVEEdge) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*ComponentCVEEdge_FixedBy)(nil),
	}
}

func (m *ComponentCVEEdge) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ComponentCVEEdge) Clone() *ComponentCVEEdge {
	if m == nil {
		return nil
	}
	cloned := new(ComponentCVEEdge)
	*cloned = *m

	if m.HasFixedBy != nil {
		cloned.HasFixedBy = m.HasFixedBy.Clone()
	}
	return cloned
}

type ImageCVEEdge struct {
	// base 64 encoded Image:CVE ids.
	Id                   string             `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" sql:"pk,id"`
	FirstImageOccurrence *types.Timestamp   `protobuf:"bytes,2,opt,name=first_image_occurrence,json=firstImageOccurrence,proto3" json:"first_image_occurrence,omitempty" search:"First Image Occurrence Timestamp,hidden"`
	State                VulnerabilityState `protobuf:"varint,3,opt,name=state,proto3,enum=storage.VulnerabilityState" json:"state,omitempty" search:"Vulnerability State"`
	ImageId              string             `protobuf:"bytes,4,opt,name=image_id,json=imageId,proto3" json:"image_id,omitempty" sql:"pk,fk(Image:id)"`
	ImageCveId           string             `protobuf:"bytes,5,opt,name=image_cve_id,json=imageCveId,proto3" json:"image_cve_id,omitempty" sql:"pk,fk(CVE:id),no-fk-constraint"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *ImageCVEEdge) Reset()         { *m = ImageCVEEdge{} }
func (m *ImageCVEEdge) String() string { return proto.CompactTextString(m) }
func (*ImageCVEEdge) ProtoMessage()    {}
func (*ImageCVEEdge) Descriptor() ([]byte, []int) {
	return fileDescriptor_62f882e266fcf764, []int{2}
}
func (m *ImageCVEEdge) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ImageCVEEdge) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ImageCVEEdge.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ImageCVEEdge) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ImageCVEEdge.Merge(m, src)
}
func (m *ImageCVEEdge) XXX_Size() int {
	return m.Size()
}
func (m *ImageCVEEdge) XXX_DiscardUnknown() {
	xxx_messageInfo_ImageCVEEdge.DiscardUnknown(m)
}

var xxx_messageInfo_ImageCVEEdge proto.InternalMessageInfo

func (m *ImageCVEEdge) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ImageCVEEdge) GetFirstImageOccurrence() *types.Timestamp {
	if m != nil {
		return m.FirstImageOccurrence
	}
	return nil
}

func (m *ImageCVEEdge) GetState() VulnerabilityState {
	if m != nil {
		return m.State
	}
	return VulnerabilityState_OBSERVED
}

func (m *ImageCVEEdge) GetImageId() string {
	if m != nil {
		return m.ImageId
	}
	return ""
}

func (m *ImageCVEEdge) GetImageCveId() string {
	if m != nil {
		return m.ImageCveId
	}
	return ""
}

func (m *ImageCVEEdge) MessageClone() proto.Message {
	return m.Clone()
}
func (m *ImageCVEEdge) Clone() *ImageCVEEdge {
	if m == nil {
		return nil
	}
	cloned := new(ImageCVEEdge)
	*cloned = *m

	cloned.FirstImageOccurrence = m.FirstImageOccurrence.Clone()
	return cloned
}

func init() {
	proto.RegisterType((*ImageComponentEdge)(nil), "storage.ImageComponentEdge")
	proto.RegisterType((*ComponentCVEEdge)(nil), "storage.ComponentCVEEdge")
	proto.RegisterType((*ImageCVEEdge)(nil), "storage.ImageCVEEdge")
}

func init() { proto.RegisterFile("storage/relations.proto", fileDescriptor_62f882e266fcf764) }

var fileDescriptor_62f882e266fcf764 = []byte{
	// 613 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x94, 0xdf, 0x6e, 0xd3, 0x30,
	0x14, 0xc6, 0x97, 0x6e, 0x63, 0x9d, 0x37, 0xed, 0x8f, 0x35, 0x46, 0x57, 0x50, 0x13, 0x2c, 0x2e,
	0x8a, 0xd4, 0xa5, 0x68, 0xdb, 0x0d, 0xbb, 0x41, 0xca, 0x34, 0xb4, 0x22, 0x24, 0x50, 0x86, 0x76,
	0xc1, 0x4d, 0xe5, 0xc6, 0x4e, 0x66, 0x35, 0x8d, 0x8b, 0xed, 0x4d, 0xed, 0x0b, 0x70, 0xcb, 0x2d,
	0x8f, 0xc4, 0x15, 0xe2, 0x05, 0x88, 0xd0, 0x78, 0x83, 0x3c, 0x01, 0x8a, 0x9d, 0x64, 0x1d, 0x1a,
	0x68, 0x42, 0xdc, 0x59, 0xc7, 0xdf, 0xf9, 0xce, 0x39, 0x3f, 0xe7, 0x04, 0x3c, 0x90, 0x8a, 0x0b,
	0x1c, 0xd1, 0xae, 0xa0, 0x31, 0x56, 0x8c, 0x27, 0xd2, 0x1d, 0x0b, 0xae, 0x38, 0x5c, 0x2a, 0x2e,
	0x9a, 0x76, 0xc4, 0x79, 0x14, 0xd3, 0xae, 0x0e, 0x0f, 0x2e, 0xc2, 0xae, 0x62, 0x23, 0x2a, 0x15,
	0x1e, 0x8d, 0x8d, 0xb2, 0xb9, 0x59, 0x5a, 0x04, 0x97, 0xb4, 0x08, 0x6d, 0x45, 0x3c, 0xe2, 0xfa,
	0xd8, 0xcd, 0x4f, 0x26, 0x8a, 0xbe, 0xd6, 0x00, 0xec, 0x8d, 0x70, 0x44, 0x8f, 0xf8, 0x68, 0xcc,
	0x13, 0x9a, 0xa8, 0x63, 0x12, 0x51, 0x68, 0x83, 0x1a, 0x23, 0x0d, 0xcb, 0xb1, 0xda, 0xcb, 0xde,
	0x7a, 0x96, 0xda, 0x2b, 0xf2, 0x43, 0x7c, 0x88, 0xc6, 0xc3, 0x0e, 0x23, 0xc8, 0xaf, 0x31, 0x02,
	0x1f, 0x83, 0x95, 0x18, 0x4f, 0xa9, 0xe8, 0xb3, 0x84, 0xd0, 0x49, 0xa3, 0xe6, 0x58, 0xed, 0xc5,
	0x93, 0x39, 0x1f, 0xe8, 0x60, 0x2f, 0x8f, 0xc1, 0x13, 0x50, 0x8f, 0x79, 0xa0, 0x07, 0x68, 0xcc,
	0x6b, 0xa7, 0x4e, 0x96, 0xda, 0x6d, 0x49, 0xb1, 0x08, 0xce, 0x0f, 0x51, 0x55, 0xd0, 0x79, 0x5d,
	0xa8, 0x3a, 0x79, 0xd3, 0xb4, 0x73, 0xce, 0x08, 0xa1, 0x09, 0xf2, 0xab, 0x6c, 0x78, 0x00, 0xea,
	0x2c, 0xef, 0xb1, 0xcf, 0x48, 0x63, 0x41, 0x3b, 0xed, 0x64, 0xa9, 0x7d, 0xbf, 0xec, 0x29, 0x1c,
	0xb6, 0xf5, 0x08, 0x87, 0x8c, 0x3c, 0x45, 0xfe, 0x92, 0x96, 0xf6, 0x08, 0xc4, 0x00, 0x9a, 0xac,
	0xa0, 0xac, 0x94, 0xe7, 0x2f, 0xea, 0xfc, 0xfd, 0x2c, 0xb5, 0xbb, 0xbf, 0xe7, 0x57, 0x1d, 0xe5,
	0x46, 0x9d, 0x84, 0xef, 0x86, 0xc3, 0xdd, 0x80, 0x27, 0x52, 0x09, 0xcc, 0x12, 0x85, 0xfc, 0x0d,
	0x76, 0x43, 0xd5, 0x23, 0xde, 0x26, 0x58, 0x3f, 0xc7, 0xb2, 0x3f, 0x43, 0x02, 0x7d, 0xaf, 0x81,
	0x8d, 0x4a, 0x72, 0x74, 0x76, 0x7c, 0x37, 0x9c, 0xcf, 0x01, 0x60, 0xb2, 0x1f, 0xb2, 0x09, 0x1e,
	0xc4, 0x54, 0xd3, 0xac, 0x7b, 0xcd, 0x2c, 0xb5, 0xb7, 0x4b, 0x5a, 0x2f, 0xcd, 0x95, 0x41, 0x84,
	0xfc, 0x65, 0x26, 0x8b, 0x08, 0x7c, 0x01, 0xea, 0x21, 0x9b, 0x50, 0xd2, 0x1f, 0x4c, 0x0b, 0xcc,
	0x28, 0x4b, 0xed, 0xd6, 0x4c, 0x22, 0x25, 0x8e, 0x37, 0xbd, 0x09, 0xf7, 0x64, 0xce, 0x5f, 0xd2,
	0x59, 0xde, 0x14, 0xbe, 0xbd, 0x95, 0xd3, 0xc2, 0x8c, 0xd5, 0xdf, 0x38, 0xdd, 0x82, 0x05, 0xbe,
	0x02, 0xab, 0x85, 0xe3, 0x25, 0xbd, 0x66, 0xde, 0xce, 0x52, 0xfb, 0xc9, 0x8c, 0xd7, 0xd1, 0xd9,
	0xf1, 0x1f, 0x40, 0x03, 0xe3, 0x78, 0x49, 0x7b, 0xc4, 0x5b, 0x03, 0xab, 0x39, 0xe2, 0x72, 0x44,
	0xf4, 0x69, 0x1e, 0xac, 0x9a, 0x2e, 0xee, 0xca, 0xf6, 0xa3, 0x05, 0xb6, 0x43, 0x26, 0xa4, 0xea,
	0x9b, 0xa6, 0x78, 0x10, 0x5c, 0x08, 0x41, 0x93, 0xc0, 0x80, 0x5e, 0xd9, 0x6b, 0xba, 0x66, 0x9d,
	0xdc, 0x72, 0x9d, 0xdc, 0x77, 0xe5, 0x3a, 0x79, 0x07, 0x59, 0x6a, 0x3f, 0xbb, 0x66, 0x29, 0xa4,
	0x72, 0x74, 0x75, 0xe7, 0x4d, 0xe5, 0xe2, 0x54, 0xea, 0xea, 0xd3, 0xdd, 0xd2, 0xf5, 0xb4, 0xf0,
	0x5a, 0x07, 0x4f, 0xc1, 0xa2, 0x54, 0x58, 0x51, 0xfd, 0x4c, 0x6b, 0x7b, 0x0f, 0xdd, 0x62, 0x49,
	0xdd, 0xb3, 0x8b, 0x38, 0xa1, 0x02, 0x0f, 0x58, 0xcc, 0xd4, 0xf4, 0x34, 0x97, 0x78, 0x4e, 0x96,
	0xda, 0x8f, 0xca, 0xba, 0x37, 0xee, 0x1d, 0x2d, 0x40, 0xbe, 0xf1, 0xfa, 0xc7, 0xdd, 0xf8, 0x9f,
	0x2f, 0x74, 0xf0, 0xe5, 0xaa, 0x65, 0x7d, 0xbb, 0x6a, 0x59, 0x3f, 0xae, 0x5a, 0xd6, 0xe7, 0x9f,
	0xad, 0x39, 0xb0, 0xc3, 0xb8, 0x2b, 0x15, 0x0e, 0x86, 0x82, 0x4f, 0x0c, 0xd4, 0x72, 0xd4, 0xf7,
	0xe5, 0x2f, 0x6c, 0x70, 0x4f, 0xc7, 0xf7, 0x7f, 0x05, 0x00, 0x00, 0xff, 0xff, 0x9a, 0xdd, 0x90,
	0xba, 0xed, 0x04, 0x00, 0x00,
}

func (m *ImageComponentEdge) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ImageComponentEdge) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ImageComponentEdge) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ImageComponentId) > 0 {
		i -= len(m.ImageComponentId)
		copy(dAtA[i:], m.ImageComponentId)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.ImageComponentId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.ImageId) > 0 {
		i -= len(m.ImageId)
		copy(dAtA[i:], m.ImageId)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.ImageId)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Location) > 0 {
		i -= len(m.Location)
		copy(dAtA[i:], m.Location)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.Location)))
		i--
		dAtA[i] = 0x1a
	}
	if m.HasLayerIndex != nil {
		{
			size := m.HasLayerIndex.Size()
			i -= size
			if _, err := m.HasLayerIndex.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ImageComponentEdge_LayerIndex) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ImageComponentEdge_LayerIndex) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i = encodeVarintRelations(dAtA, i, uint64(m.LayerIndex))
	i--
	dAtA[i] = 0x10
	return len(dAtA) - i, nil
}
func (m *ComponentCVEEdge) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ComponentCVEEdge) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ComponentCVEEdge) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ImageCveId) > 0 {
		i -= len(m.ImageCveId)
		copy(dAtA[i:], m.ImageCveId)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.ImageCveId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.ImageComponentId) > 0 {
		i -= len(m.ImageComponentId)
		copy(dAtA[i:], m.ImageComponentId)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.ImageComponentId)))
		i--
		dAtA[i] = 0x22
	}
	if m.HasFixedBy != nil {
		{
			size := m.HasFixedBy.Size()
			i -= size
			if _, err := m.HasFixedBy.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	if m.IsFixable {
		i--
		if m.IsFixable {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x10
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ComponentCVEEdge_FixedBy) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ComponentCVEEdge_FixedBy) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	i -= len(m.FixedBy)
	copy(dAtA[i:], m.FixedBy)
	i = encodeVarintRelations(dAtA, i, uint64(len(m.FixedBy)))
	i--
	dAtA[i] = 0x1a
	return len(dAtA) - i, nil
}
func (m *ImageCVEEdge) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ImageCVEEdge) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ImageCVEEdge) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.ImageCveId) > 0 {
		i -= len(m.ImageCveId)
		copy(dAtA[i:], m.ImageCveId)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.ImageCveId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.ImageId) > 0 {
		i -= len(m.ImageId)
		copy(dAtA[i:], m.ImageId)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.ImageId)))
		i--
		dAtA[i] = 0x22
	}
	if m.State != 0 {
		i = encodeVarintRelations(dAtA, i, uint64(m.State))
		i--
		dAtA[i] = 0x18
	}
	if m.FirstImageOccurrence != nil {
		{
			size, err := m.FirstImageOccurrence.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintRelations(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintRelations(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintRelations(dAtA []byte, offset int, v uint64) int {
	offset -= sovRelations(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ImageComponentEdge) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	if m.HasLayerIndex != nil {
		n += m.HasLayerIndex.Size()
	}
	l = len(m.Location)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	l = len(m.ImageId)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	l = len(m.ImageComponentId)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ImageComponentEdge_LayerIndex) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovRelations(uint64(m.LayerIndex))
	return n
}
func (m *ComponentCVEEdge) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	if m.IsFixable {
		n += 2
	}
	if m.HasFixedBy != nil {
		n += m.HasFixedBy.Size()
	}
	l = len(m.ImageComponentId)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	l = len(m.ImageCveId)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *ComponentCVEEdge_FixedBy) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.FixedBy)
	n += 1 + l + sovRelations(uint64(l))
	return n
}
func (m *ImageCVEEdge) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	if m.FirstImageOccurrence != nil {
		l = m.FirstImageOccurrence.Size()
		n += 1 + l + sovRelations(uint64(l))
	}
	if m.State != 0 {
		n += 1 + sovRelations(uint64(m.State))
	}
	l = len(m.ImageId)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	l = len(m.ImageCveId)
	if l > 0 {
		n += 1 + l + sovRelations(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovRelations(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozRelations(x uint64) (n int) {
	return sovRelations(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *ImageComponentEdge) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRelations
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ImageComponentEdge: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ImageComponentEdge: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LayerIndex", wireType)
			}
			var v int32
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.HasLayerIndex = &ImageComponentEdge_LayerIndex{v}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Location", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Location = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageComponentId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageComponentId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRelations(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRelations
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ComponentCVEEdge) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRelations
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ComponentCVEEdge: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ComponentCVEEdge: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field IsFixable", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.IsFixable = bool(v != 0)
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FixedBy", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HasFixedBy = &ComponentCVEEdge_FixedBy{string(dAtA[iNdEx:postIndex])}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageComponentId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageComponentId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageCveId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageCveId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRelations(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRelations
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ImageCVEEdge) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowRelations
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ImageCVEEdge: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ImageCVEEdge: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FirstImageOccurrence", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.FirstImageOccurrence == nil {
				m.FirstImageOccurrence = &types.Timestamp{}
			}
			if err := m.FirstImageOccurrence.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= VulnerabilityState(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ImageCveId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthRelations
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthRelations
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ImageCveId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipRelations(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthRelations
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipRelations(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowRelations
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowRelations
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthRelations
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupRelations
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthRelations
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthRelations        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowRelations          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupRelations = fmt.Errorf("proto: unexpected end of group")
)
