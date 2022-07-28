package decompiler

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
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
	imports    string
	valuestype string
	m          sync.Mutex
}

func (this *CONSTANT_Utf8_info) isValid() bool {
	return this.tag == CONSTANT_Utf8
}

func (this *CONSTANT_Utf8_info) Values() (string, error) {
	if !this.isValid() {
		return "", fmt.Errorf("this class not valid")
	}
	l, err := this.Len()
	if err != nil {
		return "", err
	}
	tmp := this.info[2 : l+2]
	return string(tmp), nil
}

func (this *CONSTANT_Utf8_info) Len() (uint, error) {
	if !this.isValid() {
		return 0, fmt.Errorf("this class not valid")
	}
	tmp := this.info[:2]
	return uint(binary.BigEndian.Uint16(tmp)), nil
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
	ACC_PUBLIC       AccessFlags = 0x0001
	ACC_PRIVATE      AccessFlags = 0x0002
	ACC_PROTECTED    AccessFlags = 0x0004
	ACC_STATIC       AccessFlags = 0x0008
	ACC_FINAL        AccessFlags = 0x0010
	ACC_SYNCHRONIZED AccessFlags = 0x0020
	ACC_SUPER        AccessFlags = 0x0020
	ACC_BRIDGE       AccessFlags = 0x0040
	ACC_VOLATILE     AccessFlags = 0x0040
	ACC_VARARGS      AccessFlags = 0x0080
	ACC_TRANSIENT    AccessFlags = 0x0080
	ACC_NATIVE       AccessFlags = 0x0100
	ACC_INTERFACE    AccessFlags = 0x0200
	ACC_ABSTRACT     AccessFlags = 0x0400
	ACC_STRICT       AccessFlags = 0x0800
	ACC_SYNTHETIC    AccessFlags = 0x1000
	ACC_ANNOTATION   AccessFlags = 0x2000
	ACC_ENUM         AccessFlags = 0x4000
	ACC_MODULE       AccessFlags = 0x8000
)

type FieldInfo struct {
	access_flags     AccessFlags
	name_index       uint16
	descriptor_index uint16
	attributes_count uint16
	attributes       []AttributeInfo
	imports          string
}

func (this *CONSTANT_Utf8_info) IsSimple() (bool, error) {
	v, e := this.Values()
	if e != nil {
		return false, e
	}
	switch v {
	case "I", "Z":
		this.m.Lock()
		defer this.m.Unlock()
		this.valuestype, _ = typeFromSignature(v)
		return true, nil
	default:
		this.parseType()
		return false, nil
	}
}

func typeFromSignature(signature string) (string, string) {
	if len(signature) == 0 {
		return "", ""
	}
	switch signature[0] {
	case 'I':
		return "int ", ""
	case 'Z':
		return "boolean ", ""
	case 'V':
		return "boolean ", ""
	case 'L':
		CLASSL := regexp.MustCompile(`(?m)L(.*\/)*(?P<classname>.*);`)
		if CLASSL.MatchString(signature) {
			m := CLASSL.FindStringSubmatch(signature)
			//log.Print(m)
			for i, name := range CLASSL.SubexpNames() {
				if "classname" == name {
					t := m[i] + " "
					i := "import " + strings.ReplaceAll(signature[1:], "/", ".") + "// " + signature + "\n"
					return t, i
				}
			}
		}
	default:
		log.Panicf("signature: %s", signature)
		//return false, nil
	}
	return "", ""
}

func (this *CONSTANT_Utf8_info) parseType() {
	this.m.Lock()
	defer this.m.Unlock()
	CLASSL := regexp.MustCompile(`(?m)L(.*\/)*(?P<classname>.*);`)
	v, _ := this.Values()
	if CLASSL.MatchString(v) {
		m := CLASSL.FindStringSubmatch(v)
		//log.Print(m)
		for i, name := range CLASSL.SubexpNames() {
			if "classname" == name {
				this.valuestype = m[i]
				this.imports = "import " + strings.ReplaceAll(v[1:], "/", ".") + "// " + v + "\n"
			}
		}
	}
}

func (this *CONSTANT_Utf8_info) GetImports() string {
	this.m.Lock()
	defer this.m.Unlock()
	return this.imports
}

func (this *CONSTANT_Utf8_info) GetType() string {
	this.m.Lock()
	defer this.m.Unlock()
	return this.valuestype
}

func (this *FieldInfo) GetCode(cp []CpInfo) string {
	text, err := this.accessToString()
	if err != nil {
		log.Panic(err)
		return ""
	}
	c := cp[this.name_index-1]
	CUtf8 := CONSTANT_Utf8_info{}
	CUtf8.CpInfo = c
	name, err := CUtf8.Values()
	if err != nil {
		log.Panic(err)
		return ""
	}
	d := CONSTANT_Utf8_info{}
	d.CpInfo = cp[this.descriptor_index-1]
	b, err := d.IsSimple()
	if err != nil {
		log.Panic(err)
		return ""
	}
	if b {
		ft := d.GetType()
		text += ft + " " + name
	} else {
		this.imports = d.GetImports()
		ft := d.GetType()
		text += ft + " " + name
	}
	if this.attributes_count > 0 {
		log.Printf("TODO: attributes")
	}
	return text + ";"
}

func (this *FieldInfo) accessToString() (string, error) {
	str := ""
	flag := this.access_flags
MAINLOOP:
	for {
		switch flag {
		case flag | ACC_PUBLIC:
			flag = flag ^ ACC_PUBLIC
			str += "public "
		case flag | ACC_PRIVATE:
			flag = flag ^ ACC_PRIVATE
			str += "private "
		case flag | ACC_PROTECTED:
			flag = flag ^ ACC_PROTECTED
			str += "protected "
		case flag | ACC_STATIC:
			flag = flag ^ ACC_STATIC
			str += "static "
		case flag | ACC_FINAL:
			flag = flag ^ ACC_FINAL
			str += "final "
		case flag | ACC_VOLATILE:
			flag = flag ^ ACC_VOLATILE
			str += "volatile "
		case flag | ACC_TRANSIENT:
			flag = flag ^ ACC_TRANSIENT
			str += "transient "
		case flag | ACC_SYNTHETIC:
			flag = flag ^ ACC_SYNTHETIC
			str += "/* synthetic */ "
		case flag | ACC_ENUM:
			flag = flag ^ ACC_ENUM
			str += "enum "
		case 0:
			break MAINLOOP
		default:
			return str, fmt.Errorf("Field access flag unknown: %v", flag)
		}
	}
	return str, nil
}

type AttributeInfo struct {
	attribute_name_index uint16
	attribute_length     uint32
	info                 []byte
}

func (this *AttributeInfo) ToCodeAttribute() CodeAttribute {
	ca := CodeAttribute{
		attribute_name_index: this.attribute_name_index,
		attribute_length:     this.attribute_length,
	}
	ca.max_stack = binary.BigEndian.Uint16(this.info[:2])
	ca.max_locals = binary.BigEndian.Uint16(this.info[2:4])
	ca.code_length = binary.BigEndian.Uint32(this.info[4:8])
	ca.code = make([]byte, ca.code_length)
	last := 8 + ca.code_length
	copy(ca.code, this.info[8:last])
	ca.exception_table_length = binary.BigEndian.Uint16(this.info[last : last+2])
	last += 2
	ca.exception_table = make([]ExceptionTable, ca.exception_table_length)
	for i := range ca.exception_table {
		et := ExceptionTable{}
		et.start_pc = binary.BigEndian.Uint16(this.info[last : last+2])
		last += 2
		et.end_pc = binary.BigEndian.Uint16(this.info[last : last+2])
		last += 2
		et.handler_pc = binary.BigEndian.Uint16(this.info[last : last+2])
		last += 2
		et.catch_type = binary.BigEndian.Uint16(this.info[last : last+2])
		last += 2
		ca.exception_table[i] = et
	}
	ca.attributes_count = binary.BigEndian.Uint16(this.info[last : last+2])
	last += 2
	ca.attributes = make([]AttributeInfo, ca.attributes_count)
	for i := range ca.attributes {
		a := AttributeInfo{
			attribute_name_index: binary.BigEndian.Uint16(this.info[last : last+2]),
		}
		last += 2
		a.attribute_length = binary.BigEndian.Uint32(this.info[last : last+4])
		last += 4
		a.info = this.info[last : last+a.attribute_length]
		last += a.attribute_length
		ca.attributes[i] = a
	}
	return ca
}

type ExceptionTable struct {
	start_pc   uint16
	end_pc     uint16
	handler_pc uint16
	catch_type uint16
}
type CodeAttribute struct {
	attribute_name_index   uint16
	attribute_length       uint32
	max_stack              uint16
	max_locals             uint16
	code_length            uint32
	code                   []byte
	exception_table_length uint16
	exception_table        []ExceptionTable
	attributes_count       uint16
	attributes             []AttributeInfo
}

type MethodInfo struct {
	access_flags     AccessFlags
	name_index       uint16
	descriptor_index uint16
	attributes_count uint16
	attributes       []AttributeInfo
	imports          string
}

func (this *MethodInfo) GetCode(cps []CpInfo) string {
	text, err := this.accessToString()
	if err != nil {
		log.Panic(err)
	}
	c := cps[this.name_index-1]
	CUtf8 := CONSTANT_Utf8_info{}
	CUtf8.CpInfo = c
	name, err := CUtf8.Values()
	if err != nil {
		log.Panic(err)
		return ""
	}

	//if name != "<init>" {
	d := CONSTANT_Utf8_info{}
	d.CpInfo = cps[this.descriptor_index-1]
	dstring, _ := d.Values()
	rettype := regexp.MustCompile(`(?m).*\)(.*)`)
	if !rettype.MatchString(dstring) {
		log.Panic("rettype not match")
	}
	srettype := rettype.FindStringSubmatch(dstring)
	t, imports := typeFromSignature(srettype[1])
	if len(imports) > 0 {
		this.imports += imports
	}
	reparam := regexp.MustCompile(`(?m)\((I|Z|L.*\;|.?)*\).*`)
	if !reparam.MatchString(dstring) {
		log.Panic("rettype not match")
	}
	sparam := reparam.FindStringSubmatch(dstring)
	//}
	if name == "<init>" {
		text += name
	} else {
		text += t + name
	}
	text += "("
	for i, param := range sparam[1:] {
		tp, imports := typeFromSignature(param)
		if len(imports) > 0 {
			this.imports += imports
		}
		if len(tp) > 0 {
			str := fmt.Sprintf("%sparam%d,", tp, i)
			text += str
		}

	}
	if text[len(text)-1] == ',' {
		text = text[:len(text)-1]
	}
	text += ") {\n\t"
	for _, attr := range this.attributes {
		attrname := CONSTANT_Utf8_info{
			CpInfo: cps[attr.attribute_name_index-1],
		}
		val, err := attrname.Values()
		if err != nil {
			log.Panic(err)
			return ""
		}
		switch val {
		case "Code":
			ca := attr.ToCodeAttribute()
			text += opcodeTostring(ca.code)
			log.Printf("name: %s Code: %#v", name, ca.code)
		}
	}
	text += "}\n"
	return text
}

func opcodeTostring(opcode []byte) string {
	str := ""
	for i := 0; i < len(opcode); i++ {
		str += "\t"
		switch opcode[i] {
		case 0x0:
			str += "// nop\n"
		case 0x1:
			str += "// aconst_null\n"
		case 0x2:
			str += "// iconst_m1\n"
		case 0x3:
			str += "// iconst_0\n"
		case 0x4:
			str += "// iconst_1\n"
		case 0x5:
			str += "// iconst_2\n"
		case 0x6:
			str += "// iconst_3\n"
		case 0x7:
			str += "// iconst_4\n"
		case 0x8:
			str += "// iconst_5\n"
		case 0x9:
			str += "// lconst_0\n"
		case 0xa:
			str += "// lconst_1\n"
		case 0xb:
			str += "// fconst_0\n"
		case 0xc:
			str += "// fconst_1\n"
		case 0xd:
			str += "// fconst_2\n"
		case 0xe:
			str += "// dconst_0\n"
		case 0xf:
			str += "// dconst_1\n"
		case 0x10:
			str += fmt.Sprintf("// bipush %#x\n", opcode[i+1])
			i += 1
		case 0x20:
			str += "// lload_2\n"
		case 0x2a:
			str += "// aload_0\n"
		case 0x9f:
			str += fmt.Sprintf("// if_icmpeq %#x %#x\n", opcode[i+1], opcode[i+2])
			i += 2
		case 0xa7:
			str += fmt.Sprintf("// goto %#x %#x\n", opcode[i+1], opcode[i+2])
			i += 2
		case 0xac:
			str += "// ireturn\n"
		case 0xb4:
			str += fmt.Sprintf("// getfield %#x %#x\n", opcode[i+1], opcode[i+2])
			i += 2
		default:
			log.Printf("TODO opcode %#x %d", opcode[i], opcode[i])
		}
	}
	return str
}

func (this *MethodInfo) accessToString() (string, error) {
	str := ""
	flag := this.access_flags
MAINLOOP:
	for {
		switch flag {
		case flag | ACC_PUBLIC:
			flag = flag ^ ACC_PUBLIC
			str += "public "
		case flag | ACC_PRIVATE:
			flag = flag ^ ACC_PRIVATE
			str += "private "
		case flag | ACC_PROTECTED:
			flag = flag ^ ACC_PROTECTED
			str += "protected "
		case flag | ACC_STATIC:
			flag = flag ^ ACC_STATIC
			str += "static "
		case flag | ACC_FINAL:
			flag = flag ^ ACC_FINAL
			str += "final "
		case flag | ACC_SYNCHRONIZED:
			flag = flag ^ ACC_SYNCHRONIZED
			str += "synchronized "
		case flag | ACC_BRIDGE:
			flag = flag ^ ACC_BRIDGE
			str += "/* bridge method */ "
		case flag | ACC_VARARGS:
			flag = flag ^ ACC_VARARGS
			str += "/* ACC_VARARGS ??? */ "
		case flag | ACC_NATIVE:
			flag = flag ^ ACC_NATIVE
			str += "native "
		case flag | ACC_ABSTRACT:
			flag = flag ^ ACC_ABSTRACT
			str += "abstract "
		case flag | ACC_STRICT:
			flag = flag ^ ACC_STRICT
			str += "strictfp "
		case flag | ACC_SYNTHETIC:
			flag = flag ^ ACC_SYNTHETIC
			str += "/* synthetic */ "
		case 0:
			break MAINLOOP
		default:
			return str, fmt.Errorf("Method access flag unknown: %v", flag)
		}
	}
	return str, nil
}
