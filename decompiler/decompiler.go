package decompiler

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type decompiler struct {
	filename string
	class    ClassFile
}

const MAGIC uint32 = 0xCAFEBABE

func New(filename string) *decompiler {
	d := &decompiler{
		filename: filename,
	}
	return d
}

func (this *decompiler) getAtribules(reader *bufio.Reader) (AttributeInfo, error) {
	atr := AttributeInfo{}
	buff := make([]byte, 2)
	_, err := reader.Read(buff)
	if err != nil {
		return atr, err
	}
	atr.attribute_name_index = binary.BigEndian.Uint16(buff)
	buff = make([]byte, 4)
	_, err = reader.Read(buff)
	if err != nil {
		return atr, err
	}
	atr.attribute_length = binary.BigEndian.Uint32(buff)
	atr.info = make([]byte, atr.attribute_length)
	buff = make([]byte, atr.attribute_length)
	_, err = reader.Read(buff)
	if err != nil {
		return atr, err
	}
	atr.info = buff
	return atr, err
}

func (this *decompiler) ParseFile() error {
	f, err := os.Open(this.filename)
	if err != nil {
		return err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	buff := make([]byte, 4)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	if MAGIC != binary.BigEndian.Uint32(buff) {
		return fmt.Errorf("magic signature wrong. Magic signature file: %X", binary.BigEndian.Uint32(buff))
	}
	class := ClassFile{
		magic: binary.BigEndian.Uint32(buff),
	}
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.minor_version = binary.BigEndian.Uint16(buff)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.major_version = binary.BigEndian.Uint16(buff)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.constant_pool_count = binary.BigEndian.Uint16(buff)
	class.constant_pool = make([]CpInfo, class.constant_pool_count-1)
	for i := uint16(0); i < class.constant_pool_count-1; i++ {
		tag, err := reader.ReadByte()
		if err != nil {
			return err
		}
		class.constant_pool[i].tag = ConstantTag(tag)
		switch class.constant_pool[i].tag {
		case CONSTANT_Utf8:
			lbuff := make([]byte, 2)
			_, err = reader.Read(lbuff)
			if err != nil {
				return err
			}
			len := binary.BigEndian.Uint16(lbuff)
			sbuff := make([]byte, len)
			_, err = reader.Read(sbuff)
			if err != nil {
				return err
			}
			class.constant_pool[i].info = make([]byte, 2+len)
			class.constant_pool[i].info = append(lbuff, sbuff...)
			continue
		case CONSTANT_MethodHandle:
			class.constant_pool[i].info = make([]byte, 3)
			buff = make([]byte, 3)
		case CONSTANT_Long,
			CONSTANT_Double:
			class.constant_pool[i].info = make([]byte, 8)
			buff = make([]byte, 8)
		case CONSTANT_Class:
			class.constant_pool[i].info = make([]byte, 2)
			buff = make([]byte, 2)
		case CONSTANT_String,
			CONSTANT_MethodType,
			CONSTANT_Module,
			CONSTANT_Package:
			class.constant_pool[i].info = make([]byte, 2)
			buff = make([]byte, 2)
		case CONSTANT_Methodref,
			CONSTANT_Fieldref,
			CONSTANT_InterfaceMethodref,
			CONSTANT_Integer,
			CONSTANT_Float,
			CONSTANT_NameAndType,
			CONSTANT_InvokeDynamic,
			CONSTANT_Dynamic:
			class.constant_pool[i].info = make([]byte, 4)
			buff = make([]byte, 4)
		default:
			fmt.Printf("TODO: case for %v\n", class.constant_pool[i].tag)
			continue
		}
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		class.constant_pool[i].info = buff
	}
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.access_flags = AccessFlags(binary.BigEndian.Uint16(buff))
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.this_class = binary.BigEndian.Uint16(buff)
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.super_class = binary.BigEndian.Uint16(buff)
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.interfaces_count = binary.BigEndian.Uint16(buff)
	class.interfaces = make([]uint16, class.interfaces_count)
	buff = make([]byte, 2)
	for i := uint16(0); i < class.interfaces_count; i++ {
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		class.interfaces[i] = binary.BigEndian.Uint16(buff)
	}
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.fields_count = binary.BigEndian.Uint16(buff)
	class.fields = make([]FieldInfo, class.fields_count)
	for i := uint16(0); i < class.fields_count; i++ {
		fi := FieldInfo{}
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		fi.access_flags = AccessFlags(binary.BigEndian.Uint16(buff))
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		fi.name_index = binary.BigEndian.Uint16(buff)
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		fi.descriptor_index = binary.BigEndian.Uint16(buff)
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		fi.attributes_count = binary.BigEndian.Uint16(buff)
		fi.attributes = make([]AttributeInfo, fi.attributes_count)
		for i := uint16(0); i < class.attributes_count; i++ {
			fi.attributes[i], err = this.getAtribules(reader)
			if err != nil {
				return err
			}
		}
		class.fields[i] = fi
	}
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.methods_count = binary.BigEndian.Uint16(buff)
	class.methods = make([]MethodInfo, class.methods_count)
	for i := uint16(0); i < class.methods_count; i++ {
		mt := MethodInfo{}
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		mt.access_flags = AccessFlags(binary.BigEndian.Uint16(buff))
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		mt.name_index = binary.BigEndian.Uint16(buff)
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		mt.descriptor_index = binary.BigEndian.Uint16(buff)
		buff = make([]byte, 2)
		_, err = reader.Read(buff)
		if err != nil {
			return err
		}
		mt.attributes_count = binary.BigEndian.Uint16(buff)
		mt.attributes = make([]AttributeInfo, mt.attributes_count)
		for i := uint16(0); i < mt.attributes_count; i++ {
			mt.attributes[i], err = this.getAtribules(reader)
			if err != nil {
				return err
			}
		}
		class.methods[i] = mt
	}
	buff = make([]byte, 2)
	_, err = reader.Read(buff)
	if err != nil {
		return err
	}
	class.attributes_count = binary.BigEndian.Uint16(buff)
	class.attributes = make([]AttributeInfo, class.attributes_count)
	for i := uint16(0); i < class.attributes_count; i++ {
		class.attributes[i], err = this.getAtribules(reader)
		if err != nil {
			return err
		}
	}
	this.class = class
	return nil
}

func (this *decompiler) WriteFile(ofile string) error {
	var err error
	f, err := os.Create(ofile)
	if err != nil {
		return err
	}
	defer f.Close()
	text := ""
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	str, err := this.accessFlagsToString()
	if err != nil {
		return err
	}
	text += str
	thclass := CONSTANT_Class_info{
		this.class.constant_pool[this.class.this_class-1],
	}
	i := thclass.NameIndex() - 1
	nameclass := CONSTANT_Utf8_info{
		CpInfo: this.class.constant_pool[i],
	}
	fmt.Print(nameclass.Values())
	s, err := nameclass.Values()
	if err != nil {
		return err
	}
	imports := this.addImport(s)
	classnames := strings.Split(s, "/")
	text += "/* " + s + " */" + classnames[len(classnames)-1]
	if this.class.super_class > 0 {
		text += " extends "
		super := CONSTANT_Class_info{
			this.class.constant_pool[this.class.super_class-1],
		}
		i = super.NameIndex() - 1
		nameclass = CONSTANT_Utf8_info{
			CpInfo: this.class.constant_pool[i],
		}
		s, err := nameclass.Values()
		if err != nil {
			return err
		}
		imports = this.addImport(s) + imports
		classnames := strings.Split(s, "/")
		text += "/* " + s + " */" + classnames[len(classnames)-1]
	}
	if this.class.interfaces_count > 0 {
		text += "implements"
		for i, inter := range this.class.interfaces {
			sinter := CONSTANT_Class_info{
				this.class.constant_pool[inter-1],
			}
			nameclass = CONSTANT_Utf8_info{}
			nameclass.CpInfo = this.class.constant_pool[sinter.NameIndex()-1]
			interfacename, err := nameclass.Values()
			if err != nil {
				return err
			}
			text += " " + interfacename
			if i < len(this.class.interfaces) {
				text += ","
			}
		}
	}
	text += " {\n\t"

	if this.class.fields_count > 0 {
		for _, field := range this.class.fields {
			text += field.GetCode(this.class.constant_pool) + "\n\t"
			imports += field.imports
			if err != nil {
				return err
			}
			//text += s + "\n"

		}
	}
	if this.class.methods_count > 0 {
		for _, method := range this.class.methods {
			text += method.GetCode(this.class.constant_pool) + "\n\t"
			imports += method.imports
		}
	}

	text += "\n}"
	text = imports + "\n\n" + text
	name := classnames[len(classnames)-1]
	text = strings.ReplaceAll(text, "<init>", name) //Пока так
	_, err = writer.WriteString(text)
	if err != nil {
		return err
	}
	return err
}

/*
func (this *decompiler) fieldToString(fi FieldInfo) (string, error) {
	str := ""
	access, err := this.AccessFlagsToString(fi.access_flags)
	if err != nil {
		return str, err
	}
	str += access
	desriptor := CONSTANT_Utf8_info{this.class.constant_pool[fi.descriptor_index-1]}
	desriptor.Values()

	log.Print(desriptor.Values())
	return str, err
}
*/

func (this *decompiler) accessFlagsToString() (string, error) {
	str := ""
	flags := this.class.access_flags

MAINLOOP:
	for {
		switch flags {
		case flags | ACC_PUBLIC:
			flags = flags ^ ACC_PUBLIC
			str = fmt.Sprint(str, "public ")
		case flags | ACC_FINAL:
			flags = flags ^ ACC_FINAL
			str = fmt.Sprint(str, "final ")
		case flags | ACC_SUPER:
			flags = flags ^ ACC_SUPER
			str = fmt.Sprint(str, "/* super */ ")
		case flags | ACC_ABSTRACT:
			flags = flags ^ ACC_ABSTRACT
			str = fmt.Sprint(str, "abstract ")
		case flags | ACC_INTERFACE:
			flags = flags ^ ACC_INTERFACE
			str = fmt.Sprint(str, "interface ")
		case flags | ACC_SYNTHETIC:
			flags = flags ^ ACC_SYNTHETIC
			str = fmt.Sprint(str, "/* synthetic */ ")
		case flags | ACC_ANNOTATION:
			flags = flags ^ ACC_ANNOTATION
			str = fmt.Sprint(str, "@interface ")
		case flags | ACC_ENUM:
			flags = flags ^ ACC_ENUM
			str = fmt.Sprint(str, "enum ")
		case flags | ACC_MODULE:
			flags = flags ^ ACC_MODULE
			str = fmt.Sprint(str, "module ")
		default:
			str = fmt.Sprint(str, "class ")
			break MAINLOOP
		}
	}
	return str, nil
}

func (this *decompiler) addImport(classname string) string {
	imports := strings.ReplaceAll(classname, "/", ".")
	return "import " + imports + ";\n"
}
