package decompiler

import (
	"bytes"
	"encoding/binary"
)

type ClassFile struct {
	magic               uint32
	minor_version       uint16
	major_version       uint16
	constant_pool_count uint16
	constant_pool       []CpInfo
	access_flags        AccessFlags
	this_class          uint16
	super_class         uint16
	interfaces_count    uint16
	interfaces          []uint16
	fields_count        uint16
	fields              []FieldInfo
	methods_count       uint16
	methods             []MethodInfo
	attributes_count    uint16
	attributes          []AttributeInfo
}

type ConstantTag byte

const (
	CONSTANT_Class              ConstantTag = 7  //--
	CONSTANT_Fieldref           ConstantTag = 9  //--
	CONSTANT_Methodref          ConstantTag = 10 //--
	CONSTANT_InterfaceMethodref ConstantTag = 11 //--
	CONSTANT_String             ConstantTag = 8  //--
	CONSTANT_Integer            ConstantTag = 3  //--
	CONSTANT_Float              ConstantTag = 4  //--
	CONSTANT_Long               ConstantTag = 5  //--
	CONSTANT_Double             ConstantTag = 6  //--
	CONSTANT_NameAndType        ConstantTag = 12 //--
	CONSTANT_Utf8               ConstantTag = 1  //--
	CONSTANT_MethodHandle       ConstantTag = 15 //--
	CONSTANT_MethodType         ConstantTag = 16 //--
	CONSTANT_Dynamic            ConstantTag = 17 //--
	CONSTANT_InvokeDynamic      ConstantTag = 18 //--
	CONSTANT_Module             ConstantTag = 19
	CONSTANT_Package            ConstantTag = 20
)

type CONSTANT_Module_info struct {
	CpInfo
}

func (this *CONSTANT_Module_info) NameIndex() uint16 {
	return binary.BigEndian.Uint16(this.info)
}

type CONSTANT_Package_info struct {
	CpInfo
}

func (this *CONSTANT_Package_info) NameIndex() uint16 {
	return binary.BigEndian.Uint16(this.info)
}

type CONSTANT_Dynamic_info struct {
	CpInfo
}

func (this *CONSTANT_Dynamic_info) BootstrapMethodAttrIndex() uint16 {
	tmp := this.info[:2]
	return binary.BigEndian.Uint16(tmp)
}

func (this *CONSTANT_Dynamic_info) NameAndTypeIndex() uint16 {
	tmp := this.info[2:]
	return binary.BigEndian.Uint16(tmp)
}

type CONSTANT_InvokeDynamic_info struct {
	CpInfo
}

func (this *CONSTANT_InvokeDynamic_info) BootstrapMethodAttrIndex() uint16 {
	tmp := this.info[:2]
	return binary.BigEndian.Uint16(tmp)
}

func (this *CONSTANT_InvokeDynamic_info) NameAndTypeIndex() uint16 {
	tmp := this.info[2:]
	return binary.BigEndian.Uint16(tmp)
}

type CONSTANT_MethodType_info struct {
	CpInfo
}

func (this *CONSTANT_MethodType_info) DescriptorIndex() uint16 {
	return binary.BigEndian.Uint16(this.info)
}

type ReferenceKind byte

const (
	REF_getField ReferenceKind = iota + 1
	REF_getStatic
	REF_putField
	REF_putStatic
	REF_invokeVirtual
	REF_invokeStatic
	REF_invokeSpecial
	REF_newInvokeSpecial
	REF_invokeInterface
)

type CONSTANT_MethodHandle_info struct {
	CpInfo
}

func (this *CONSTANT_MethodHandle_info) ReferenceKind() ReferenceKind {
	tmp := this.info[0]
	return ReferenceKind(tmp)
}

func (this *CONSTANT_MethodHandle_info) ReferenceIndex() uint16 {
	tmp := this.info[1:]
	return binary.BigEndian.Uint16(tmp)
}

type CONSTANT_Long_info struct {
	CpInfo
}

func (this *CONSTANT_Long_info) HighBytes() uint32 {
	tmp := this.info[:4]
	return binary.BigEndian.Uint32(tmp)
}

func (this *CONSTANT_Long_info) LowBytes() uint32 {
	tmp := this.info[4:]
	return binary.BigEndian.Uint32(tmp)
}

func (this *CONSTANT_Long_info) Values() int64 {
	var ret int64
	binary.Read(bytes.NewBuffer(this.info), binary.BigEndian, &ret)
	return ret
}

type CONSTANT_Double_info struct {
	CpInfo
}

func (this *CONSTANT_Double_info) HighBytes() uint32 {
	tmp := this.info[:4]
	return binary.BigEndian.Uint32(tmp)
}

func (this *CONSTANT_Double_info) LowBytes() uint32 {
	tmp := this.info[4:]
	return binary.BigEndian.Uint32(tmp)
}

func (this *CONSTANT_Double_info) Values() float64 {
	var ret float64
	binary.Read(bytes.NewBuffer(this.info), binary.BigEndian, &ret)
	return ret
}

type CONSTANT_NameAndType_info struct {
	CpInfo
}

func (this *CONSTANT_NameAndType_info) NameIndex() uint16 {
	tmp := this.info[:2]
	return binary.BigEndian.Uint16(tmp)
}

func (this *CONSTANT_NameAndType_info) DescriptorIndex() uint16 {
	tmp := this.info[2:]
	return binary.BigEndian.Uint16(tmp)
}

type CONSTANT_Utf8_info struct {
	CpInfo
}

func (this *CONSTANT_Utf8_info) Values() string {
	tmp := this.info[2:this.Len()]
	return string(tmp)
}

func (this *CONSTANT_Utf8_info) Len() uint {
	tmp := this.info[:2]
	return uint(binary.BigEndian.Uint16(tmp))
}

type CpInfo struct {
	tag  ConstantTag
	info []byte
}

type CONSTANT_Class_info struct {
	CpInfo
}

func (this *CONSTANT_Class_info) NameIndex() uint16 {
	return binary.BigEndian.Uint16(this.info)
}

type CONSTANT_Methodref_info struct {
	CpInfo
}

func (this *CONSTANT_Methodref_info) ClassIndex() uint16 {
	tmp := this.info[:2]
	return binary.BigEndian.Uint16(tmp)
}

func (this *CONSTANT_Methodref_info) NameAndTypeIndex() uint16 {
	tmp := this.info[2:]
	return binary.BigEndian.Uint16(tmp)
}

type CONSTANT_Fieldref_info CONSTANT_Methodref_info

type CONSTANT_InterfaceMethodref_info CONSTANT_Methodref_info

type CONSTANT_String_info struct {
	CpInfo
}

func (this *CONSTANT_String_info) StringIndex() uint16 {
	return binary.BigEndian.Uint16(this.info)
}

type CONSTANT_Integer_info struct {
	CpInfo
}

func (this *CONSTANT_Integer_info) Values() int {
	var ret int32
	binary.Read(bytes.NewBuffer(this.info), binary.BigEndian, &ret)
	return int(ret)
}

type CONSTANT_Float_info struct {
	CpInfo
}

func (this *CONSTANT_Float_info) Values() float32 {
	var ret float32
	binary.Read(bytes.NewBuffer(this.info), binary.BigEndian, &ret)
	return ret
}

type AccessFlags uint16

const (
	ACC_PUBLIC     AccessFlags = 0x0001
	ACC_PRIVATE    AccessFlags = 0x0002
	ACC_PROTECTED  AccessFlags = 0x0004
	ACC_STATIC     AccessFlags = 0x0008
	ACC_FINAL      AccessFlags = 0x0010
	ACC_SUPER      AccessFlags = 0x0020
	ACC_VOLATILE   AccessFlags = 0x0040
	ACC_TRANSIENT  AccessFlags = 0x0080
	ACC_INTERFACE  AccessFlags = 0x0200
	ACC_ABSTRACT   AccessFlags = 0x0400
	ACC_SYNTHETIC  AccessFlags = 0x1000
	ACC_ANNOTATION AccessFlags = 0x2000
	ACC_ENUM       AccessFlags = 0x4000
	ACC_MODULE     AccessFlags = 0x8000
)

type FieldInfo struct {
	access_flags     AccessFlags
	name_index       uint16
	descriptor_index uint16
	attributes_count uint16
	attributes       []AttributeInfo
}

type AttributeInfo struct {
	attribute_name_index uint16
	attribute_length     uint32
	info                 []byte
}

type MethodInfo struct {
	access_flags     AccessFlags
	name_index       uint16
	descriptor_index uint16
	attributes_count uint16
	attributes       []AttributeInfo
}
