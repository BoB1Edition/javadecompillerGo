package decompiler

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"os"
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
		log.Printf("i: %d", i)
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
		case CONSTANT_Class,
			CONSTANT_String,
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
	for i := uint16(0); i < class.fields_count; i++ {
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
	writer := bufio.NewWriter(f)
	defer writer.Flush()
	flags := this.class.access_flags
	for flags != 0 {
		switch flags {
		case flags | ACC_PUBLIC:
			flags = flags ^ ACC_PUBLIC
			fmt.Println("public")
		case flags | ACC_FINAL:
			flags = flags ^ ACC_FINAL
			fmt.Println("final")
		case flags | ACC_SUPER:
			flags = flags ^ ACC_SUPER
			fmt.Println("super")
		case flags | ACC_INTERFACE:
			flags = flags ^ ACC_INTERFACE
			fmt.Println("interface")
		case flags | ACC_ABSTRACT:
			flags = flags ^ ACC_ABSTRACT
			fmt.Println("abstract")
		case flags | ACC_SYNTHETIC:
			flags = flags ^ ACC_SYNTHETIC
			fmt.Println("synthetic")
		case flags | ACC_ANNOTATION:
			flags = flags ^ ACC_ANNOTATION
			fmt.Println("annotation")
		case flags | ACC_ENUM:
			flags = flags ^ ACC_ENUM
			fmt.Println("enum")
		case flags | ACC_MODULE:
			flags = flags ^ ACC_MODULE
			fmt.Println("module")
		default:
			break
		}
	}
	return err
}
