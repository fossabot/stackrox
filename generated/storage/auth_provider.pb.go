// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: storage/auth_provider.proto

package storage

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
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

// Next Tag: 9
type AuthProvider struct {
	Id         string            `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" sql:"pk"`
	Name       string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty" sql:"unique"`
	Type       string            `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	UiEndpoint string            `protobuf:"bytes,4,opt,name=ui_endpoint,json=uiEndpoint,proto3" json:"ui_endpoint,omitempty"`
	Enabled    bool              `protobuf:"varint,5,opt,name=enabled,proto3" json:"enabled,omitempty"`
	Config     map[string]string `protobuf:"bytes,6,rep,name=config,proto3" json:"config,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// The login URL will be provided by the backend, and may not be specified in a request.
	LoginUrl  string `protobuf:"bytes,7,opt,name=login_url,json=loginUrl,proto3" json:"login_url,omitempty"`
	Validated bool   `protobuf:"varint,8,opt,name=validated,proto3" json:"validated,omitempty"` // Deprecated: Do not use.
	// UI endpoints which to allow in addition to `ui_endpoint`. I.e., if a login request
	// is coming from any of these, the auth request will use these for the callback URL,
	// not ui_endpoint.
	ExtraUiEndpoints     []string `protobuf:"bytes,9,rep,name=extra_ui_endpoints,json=extraUiEndpoints,proto3" json:"extra_ui_endpoints,omitempty"`
	Active               bool     `protobuf:"varint,10,opt,name=active,proto3" json:"active,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AuthProvider) Reset()         { *m = AuthProvider{} }
func (m *AuthProvider) String() string { return proto.CompactTextString(m) }
func (*AuthProvider) ProtoMessage()    {}
func (*AuthProvider) Descriptor() ([]byte, []int) {
	return fileDescriptor_4ed6b69aa5a381c8, []int{0}
}
func (m *AuthProvider) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AuthProvider) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AuthProvider.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AuthProvider) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AuthProvider.Merge(m, src)
}
func (m *AuthProvider) XXX_Size() int {
	return m.Size()
}
func (m *AuthProvider) XXX_DiscardUnknown() {
	xxx_messageInfo_AuthProvider.DiscardUnknown(m)
}

var xxx_messageInfo_AuthProvider proto.InternalMessageInfo

func (m *AuthProvider) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *AuthProvider) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AuthProvider) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *AuthProvider) GetUiEndpoint() string {
	if m != nil {
		return m.UiEndpoint
	}
	return ""
}

func (m *AuthProvider) GetEnabled() bool {
	if m != nil {
		return m.Enabled
	}
	return false
}

func (m *AuthProvider) GetConfig() map[string]string {
	if m != nil {
		return m.Config
	}
	return nil
}

func (m *AuthProvider) GetLoginUrl() string {
	if m != nil {
		return m.LoginUrl
	}
	return ""
}

// Deprecated: Do not use.
func (m *AuthProvider) GetValidated() bool {
	if m != nil {
		return m.Validated
	}
	return false
}

func (m *AuthProvider) GetExtraUiEndpoints() []string {
	if m != nil {
		return m.ExtraUiEndpoints
	}
	return nil
}

func (m *AuthProvider) GetActive() bool {
	if m != nil {
		return m.Active
	}
	return false
}

func (m *AuthProvider) MessageClone() proto.Message {
	return m.Clone()
}
func (m *AuthProvider) Clone() *AuthProvider {
	if m == nil {
		return nil
	}
	cloned := new(AuthProvider)
	*cloned = *m

	if m.Config != nil {
		cloned.Config = make(map[string]string, len(m.Config))
		for k, v := range m.Config {
			cloned.Config[k] = v
		}
	}
	if m.ExtraUiEndpoints != nil {
		cloned.ExtraUiEndpoints = make([]string, len(m.ExtraUiEndpoints))
		copy(cloned.ExtraUiEndpoints, m.ExtraUiEndpoints)
	}
	return cloned
}

func init() {
	proto.RegisterType((*AuthProvider)(nil), "storage.AuthProvider")
	proto.RegisterMapType((map[string]string)(nil), "storage.AuthProvider.ConfigEntry")
}

func init() { proto.RegisterFile("storage/auth_provider.proto", fileDescriptor_4ed6b69aa5a381c8) }

var fileDescriptor_4ed6b69aa5a381c8 = []byte{
	// 387 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x86, 0x59, 0x3b, 0x75, 0xec, 0x49, 0x0e, 0xd1, 0xaa, 0x42, 0x4b, 0x8b, 0x1c, 0x63, 0x71,
	0xf0, 0x01, 0xb9, 0x12, 0x70, 0xa0, 0xbd, 0x11, 0xd4, 0x3b, 0x5a, 0xa9, 0x17, 0x2e, 0xd6, 0x36,
	0x5e, 0xdc, 0x55, 0xcc, 0xae, 0xbb, 0xde, 0xb5, 0x9a, 0x37, 0xe1, 0xc0, 0x03, 0x71, 0xe4, 0x09,
	0x22, 0x14, 0xde, 0x20, 0x4f, 0x80, 0xbc, 0x76, 0x20, 0xb7, 0x99, 0xff, 0x9b, 0xd1, 0xcc, 0xfc,
	0x03, 0x97, 0xad, 0x51, 0x9a, 0x55, 0xfc, 0x8a, 0x59, 0xf3, 0x50, 0x34, 0x5a, 0x75, 0xa2, 0xe4,
	0x3a, 0x6f, 0xb4, 0x32, 0x0a, 0x4f, 0x47, 0x78, 0x71, 0x5e, 0xa9, 0x4a, 0x39, 0xed, 0xaa, 0x8f,
	0x06, 0x9c, 0xfe, 0xf0, 0x61, 0xfe, 0xd1, 0x9a, 0x87, 0xcf, 0x63, 0x17, 0x7e, 0x09, 0x9e, 0x28,
	0x09, 0x4a, 0x50, 0x16, 0xad, 0xe6, 0x87, 0xdd, 0x32, 0x6c, 0x1f, 0xeb, 0x9b, 0xb4, 0xd9, 0xa4,
	0xd4, 0x13, 0x25, 0x7e, 0x0d, 0x13, 0xc9, 0xbe, 0x71, 0xe2, 0x39, 0xbe, 0x38, 0xec, 0x96, 0x73,
	0xc7, 0xad, 0x14, 0x8f, 0x96, 0xa7, 0xd4, 0x51, 0x8c, 0x61, 0x62, 0xb6, 0x0d, 0x27, 0x7e, 0x5f,
	0x45, 0x5d, 0x8c, 0x97, 0x30, 0xb3, 0xa2, 0xe0, 0xb2, 0x6c, 0x94, 0x90, 0x86, 0x4c, 0x1c, 0x02,
	0x2b, 0x6e, 0x47, 0x05, 0x13, 0x98, 0x72, 0xc9, 0xee, 0x6b, 0x5e, 0x92, 0xb3, 0x04, 0x65, 0x21,
	0x3d, 0xa6, 0xf8, 0x1a, 0x82, 0xb5, 0x92, 0x5f, 0x45, 0x45, 0x82, 0xc4, 0xcf, 0x66, 0x6f, 0x5f,
	0xe5, 0xe3, 0x4d, 0xf9, 0xe9, 0xe6, 0xf9, 0x27, 0x57, 0x73, 0x2b, 0x8d, 0xde, 0xd2, 0xb1, 0x01,
	0x5f, 0x42, 0x54, 0xab, 0x4a, 0xc8, 0xc2, 0xea, 0x9a, 0x4c, 0xdd, 0xcc, 0xd0, 0x09, 0x77, 0xba,
	0xc6, 0x09, 0x44, 0x1d, 0xab, 0x45, 0xc9, 0x0c, 0x2f, 0x49, 0xd8, 0xcf, 0x5c, 0x79, 0x04, 0xd1,
	0xff, 0x22, 0x7e, 0x03, 0x98, 0x3f, 0x19, 0xcd, 0x8a, 0x93, 0xd5, 0x5b, 0x12, 0x25, 0x7e, 0x16,
	0xd1, 0x85, 0x23, 0x77, 0xff, 0x0e, 0x68, 0xf1, 0x73, 0x08, 0xd8, 0xda, 0x88, 0x8e, 0x13, 0x70,
	0x07, 0x8c, 0xd9, 0xc5, 0x35, 0xcc, 0x4e, 0x76, 0xc3, 0x0b, 0xf0, 0x37, 0x7c, 0x3b, 0x58, 0x4c,
	0xfb, 0x10, 0x9f, 0xc3, 0x59, 0xc7, 0x6a, 0x3b, 0xda, 0x4a, 0x87, 0xe4, 0xc6, 0xfb, 0x80, 0x56,
	0xef, 0x7f, 0xee, 0x63, 0xf4, 0x6b, 0x1f, 0xa3, 0xdf, 0xfb, 0x18, 0x7d, 0xff, 0x13, 0x3f, 0x83,
	0x17, 0x42, 0xe5, 0xad, 0x61, 0xeb, 0x8d, 0x56, 0x4f, 0xc3, 0x0f, 0x8f, 0x6e, 0x7c, 0x39, 0xbe,
	0xfa, 0x3e, 0x70, 0xfa, 0xbb, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff, 0x2c, 0x6b, 0xa9, 0xe3, 0x19,
	0x02, 0x00, 0x00,
}

func (m *AuthProvider) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AuthProvider) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AuthProvider) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Active {
		i--
		if m.Active {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x50
	}
	if len(m.ExtraUiEndpoints) > 0 {
		for iNdEx := len(m.ExtraUiEndpoints) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.ExtraUiEndpoints[iNdEx])
			copy(dAtA[i:], m.ExtraUiEndpoints[iNdEx])
			i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.ExtraUiEndpoints[iNdEx])))
			i--
			dAtA[i] = 0x4a
		}
	}
	if m.Validated {
		i--
		if m.Validated {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x40
	}
	if len(m.LoginUrl) > 0 {
		i -= len(m.LoginUrl)
		copy(dAtA[i:], m.LoginUrl)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.LoginUrl)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.Config) > 0 {
		for k := range m.Config {
			v := m.Config[k]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintAuthProvider(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(k)
			copy(dAtA[i:], k)
			i = encodeVarintAuthProvider(dAtA, i, uint64(len(k)))
			i--
			dAtA[i] = 0xa
			i = encodeVarintAuthProvider(dAtA, i, uint64(baseI-i))
			i--
			dAtA[i] = 0x32
		}
	}
	if m.Enabled {
		i--
		if m.Enabled {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if len(m.UiEndpoint) > 0 {
		i -= len(m.UiEndpoint)
		copy(dAtA[i:], m.UiEndpoint)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.UiEndpoint)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Type) > 0 {
		i -= len(m.Type)
		copy(dAtA[i:], m.Type)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.Type)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintAuthProvider(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAuthProvider(dAtA []byte, offset int, v uint64) int {
	offset -= sovAuthProvider(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *AuthProvider) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	l = len(m.Type)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	l = len(m.UiEndpoint)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	if m.Enabled {
		n += 2
	}
	if len(m.Config) > 0 {
		for k, v := range m.Config {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovAuthProvider(uint64(len(k))) + 1 + len(v) + sovAuthProvider(uint64(len(v)))
			n += mapEntrySize + 1 + sovAuthProvider(uint64(mapEntrySize))
		}
	}
	l = len(m.LoginUrl)
	if l > 0 {
		n += 1 + l + sovAuthProvider(uint64(l))
	}
	if m.Validated {
		n += 2
	}
	if len(m.ExtraUiEndpoints) > 0 {
		for _, s := range m.ExtraUiEndpoints {
			l = len(s)
			n += 1 + l + sovAuthProvider(uint64(l))
		}
	}
	if m.Active {
		n += 2
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovAuthProvider(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAuthProvider(x uint64) (n int) {
	return sovAuthProvider(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *AuthProvider) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAuthProvider
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
			return fmt.Errorf("proto: AuthProvider: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AuthProvider: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UiEndpoint", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.UiEndpoint = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Enabled", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
			m.Enabled = bool(v != 0)
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Config == nil {
				m.Config = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
				var wire uint64
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowAuthProvider
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAuthProvider
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthAuthProvider
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthAuthProvider
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowAuthProvider
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthAuthProvider
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthAuthProvider
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipAuthProvider(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthAuthProvider
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Config[mapkey] = mapvalue
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LoginUrl", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LoginUrl = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Validated", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
			m.Validated = bool(v != 0)
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtraUiEndpoints", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
				return ErrInvalidLengthAuthProvider
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAuthProvider
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtraUiEndpoints = append(m.ExtraUiEndpoints, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 10:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Active", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAuthProvider
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
			m.Active = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipAuthProvider(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAuthProvider
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
func skipAuthProvider(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAuthProvider
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
					return 0, ErrIntOverflowAuthProvider
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
					return 0, ErrIntOverflowAuthProvider
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
				return 0, ErrInvalidLengthAuthProvider
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAuthProvider
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAuthProvider
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAuthProvider        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAuthProvider          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAuthProvider = fmt.Errorf("proto: unexpected end of group")
)
