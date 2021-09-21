// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ethermint/nameservice/v1/genesis.proto

package types

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// GenesisState defines the nameservice module's genesis state.
type GenesisState struct {
	// params defines all the nameservice of the module.
	Params Params `protobuf:"bytes,1,opt,name=params,proto3" json:"params"`
	// records
	Records []Record `protobuf:"bytes,2,rep,name=records,proto3" json:"records" json:"records" yaml:"records"`
	// authorities
	Authorities []AuthorityEntry `protobuf:"bytes,3,rep,name=authorities,proto3" json:"authorities" json:"authorities" yaml:"authorities"`
	// names
	Names []NameEntry `protobuf:"bytes,4,rep,name=names,proto3" json:"names" json:"names" yaml:"names"`
}

func (m *GenesisState) Reset()         { *m = GenesisState{} }
func (m *GenesisState) String() string { return proto.CompactTextString(m) }
func (*GenesisState) ProtoMessage()    {}
func (*GenesisState) Descriptor() ([]byte, []int) {
	return fileDescriptor_f532cf62390df3f0, []int{0}
}
func (m *GenesisState) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GenesisState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GenesisState.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GenesisState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenesisState.Merge(m, src)
}
func (m *GenesisState) XXX_Size() int {
	return m.Size()
}
func (m *GenesisState) XXX_DiscardUnknown() {
	xxx_messageInfo_GenesisState.DiscardUnknown(m)
}

var xxx_messageInfo_GenesisState proto.InternalMessageInfo

func (m *GenesisState) GetParams() Params {
	if m != nil {
		return m.Params
	}
	return Params{}
}

func (m *GenesisState) GetRecords() []Record {
	if m != nil {
		return m.Records
	}
	return nil
}

func (m *GenesisState) GetAuthorities() []AuthorityEntry {
	if m != nil {
		return m.Authorities
	}
	return nil
}

func (m *GenesisState) GetNames() []NameEntry {
	if m != nil {
		return m.Names
	}
	return nil
}

func init() {
	proto.RegisterType((*GenesisState)(nil), "ethermint.nameservice.v1.GenesisState")
}

func init() {
	proto.RegisterFile("ethermint/nameservice/v1/genesis.proto", fileDescriptor_f532cf62390df3f0)
}

var fileDescriptor_f532cf62390df3f0 = []byte{
	// 334 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0x31, 0x4b, 0xf3, 0x40,
	0x18, 0x80, 0x93, 0xaf, 0xfd, 0x2a, 0xa4, 0x4e, 0xc1, 0x21, 0x16, 0xbc, 0xc6, 0x4a, 0xa5, 0x08,
	0xde, 0x59, 0xdd, 0x1c, 0x04, 0x03, 0x22, 0x38, 0x88, 0xc4, 0x4d, 0x5c, 0xae, 0xf5, 0x25, 0x39,
	0xf1, 0x72, 0xe5, 0xee, 0x5a, 0xcc, 0x5f, 0x70, 0xf2, 0x67, 0x75, 0xec, 0xe8, 0x54, 0xa4, 0xfd,
	0x07, 0xfe, 0x02, 0xe9, 0x5d, 0x6a, 0xe3, 0x10, 0xb7, 0x7b, 0x8f, 0xe7, 0x7d, 0x9e, 0xe1, 0xf5,
	0x0e, 0x41, 0xa7, 0x20, 0x39, 0xcb, 0x34, 0xc9, 0x28, 0x07, 0x05, 0x72, 0xc2, 0x86, 0x40, 0x26,
	0x7d, 0x92, 0x40, 0x06, 0x8a, 0x29, 0x3c, 0x92, 0x42, 0x0b, 0x3f, 0xf8, 0xe1, 0x70, 0x89, 0xc3,
	0x93, 0x7e, 0x6b, 0x27, 0x11, 0x89, 0x30, 0x10, 0x59, 0xbd, 0x2c, 0xdf, 0x3a, 0xaa, 0xf4, 0x96,
	0xd7, 0x0d, 0xdb, 0x79, 0xab, 0x79, 0xdb, 0xd7, 0xb6, 0x76, 0xaf, 0xa9, 0x06, 0xff, 0xc2, 0x6b,
	0x8c, 0xa8, 0xa4, 0x5c, 0x05, 0x6e, 0xe8, 0xf6, 0x9a, 0xa7, 0x21, 0xae, 0xaa, 0xe3, 0x3b, 0xc3,
	0x45, 0xf5, 0xe9, 0xbc, 0xed, 0xc4, 0xc5, 0x96, 0x4f, 0xbd, 0x2d, 0x09, 0x43, 0x21, 0x9f, 0x54,
	0xf0, 0x2f, 0xac, 0xfd, 0x2d, 0x88, 0x0d, 0x18, 0x75, 0x57, 0x82, 0xaf, 0x79, 0x7b, 0xef, 0x59,
	0x89, 0xec, 0xbc, 0x53, 0xac, 0x77, 0xc2, 0x9c, 0xf2, 0x97, 0xcd, 0x18, 0xaf, 0xbd, 0x7e, 0xee,
	0x35, 0xe9, 0x58, 0xa7, 0x42, 0x32, 0xcd, 0x40, 0x05, 0x35, 0x93, 0xe9, 0x55, 0x67, 0x2e, 0x0b,
	0x38, 0xbf, 0xca, 0xb4, 0xcc, 0xa3, 0xe3, 0x22, 0xd7, 0xb5, 0xb9, 0x92, 0x6a, 0x9d, 0x2c, 0x7f,
	0xc5, 0xe5, 0x96, 0xff, 0xe8, 0xfd, 0x37, 0xf2, 0xa0, 0x6e, 0xa2, 0x07, 0xd5, 0xd1, 0x5b, 0xca,
	0xc1, 0xf6, 0xf6, 0x8b, 0xde, 0xae, 0xed, 0x19, 0x6e, 0x5d, 0xb2, 0x43, 0x6c, 0xa5, 0xd1, 0xcd,
	0x74, 0x81, 0xdc, 0xd9, 0x02, 0xb9, 0x9f, 0x0b, 0xe4, 0xbe, 0x2f, 0x91, 0x33, 0x5b, 0x22, 0xe7,
	0x63, 0x89, 0x9c, 0x87, 0x93, 0x84, 0xe9, 0x74, 0x3c, 0xc0, 0x43, 0xc1, 0x89, 0x4e, 0xa9, 0x54,
	0x4c, 0x91, 0xcd, 0x95, 0x5f, 0x7f, 0xdd, 0x59, 0xe7, 0x23, 0x50, 0x83, 0x86, 0xb9, 0xef, 0xd9,
	0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2e, 0x1b, 0xb0, 0x86, 0x65, 0x02, 0x00, 0x00,
}

func (m *GenesisState) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GenesisState) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GenesisState) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Names) > 0 {
		for iNdEx := len(m.Names) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Names[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Authorities) > 0 {
		for iNdEx := len(m.Authorities) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Authorities[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Records) > 0 {
		for iNdEx := len(m.Records) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Records[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenesis(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.Params.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenesis(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenesis(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenesis(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GenesisState) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Params.Size()
	n += 1 + l + sovGenesis(uint64(l))
	if len(m.Records) > 0 {
		for _, e := range m.Records {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Authorities) > 0 {
		for _, e := range m.Authorities {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	if len(m.Names) > 0 {
		for _, e := range m.Names {
			l = e.Size()
			n += 1 + l + sovGenesis(uint64(l))
		}
	}
	return n
}

func sovGenesis(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenesis(x uint64) (n int) {
	return sovGenesis(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GenesisState) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenesis
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
			return fmt.Errorf("proto: GenesisState: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GenesisState: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Params", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Params.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Records", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Records = append(m.Records, Record{})
			if err := m.Records[len(m.Records)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authorities", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Authorities = append(m.Authorities, AuthorityEntry{})
			if err := m.Authorities[len(m.Authorities)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Names", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenesis
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
				return ErrInvalidLengthGenesis
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenesis
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Names = append(m.Names, NameEntry{})
			if err := m.Names[len(m.Names)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenesis(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenesis
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenesis(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
					return 0, ErrIntOverflowGenesis
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
				return 0, ErrInvalidLengthGenesis
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenesis
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenesis
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenesis        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenesis          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenesis = fmt.Errorf("proto: unexpected end of group")
)
