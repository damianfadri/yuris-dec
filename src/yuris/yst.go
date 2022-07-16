package yuris

import (
	"io/ioutil"
	"log"
	"fmt"

	"github.com/damianfadri/yuris-decompiler/utils"
	"github.com/damianfadri/yuris-decompiler/utils/dsa"
)

type Attribute struct {
	Id					int16
	Type				[]byte
	ValueLength			int
	ValueOffset			int
	Bytes				[]byte
}

type Command struct {
	Id					byte
	NumAttributes		byte
	Offset				byte
}

type Script struct {
	Commands 	[]Command
	Attributes	[]Attribute
}

func ReadYST(path string) Script {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	br := utils.NewBinaryReader(data)

	magic := br.ReadString(4)
	if magic != "YSTB" {
		log.Fatal("Invalid magic in ysc.ybn")
	}

	// YU-RIS version
	_ = br.ReadInt32()

	// Number of instructions
	numInstructions := br.ReadInt32()
	szInstructions := br.ReadInt32()
	if (szInstructions != numInstructions * 4) {
		log.Fatal("Instruction size does not match instruction count.")
	}

	szAttrDescriptors := br.ReadInt32()

	// Attribute values size
	_ = br.ReadInt32()

	// Line numbers size
	_ = br.ReadInt32()

	br.Skip(4)

	offsetInstructions := br.Position
	offsetAttrDescriptors := offsetInstructions + szInstructions
	offsetAttrValues := offsetAttrDescriptors + szAttrDescriptors

	key := uint32(0)
	if (szAttrDescriptors > 0) {
		br.Seek(offsetAttrDescriptors + 8)
		key = uint32(br.ReadInt32())
	}

	// Decrypt script data if possible.
	decrypt(br, key)

	script := Script{}

	br = utils.NewBinaryReader(data)

	attributes := dsa.NewList[Attribute]()
	br.Seek(offsetAttrDescriptors)
	for br.Position < offsetAttrValues {
		attr := Attribute{}
		attr.Id = br.ReadInt16()
		attr.Type = br.ReadBytes(2)
		attr.ValueLength = br.ReadInt32()
		attr.ValueOffset = br.ReadInt32() + offsetAttrValues
		attr.Bytes = br.Bytes[attr.ValueOffset:attr.ValueOffset+attr.ValueLength]

		attributes.Add(attr)
	} 

	commands := dsa.NewList[Command]()
	br.Seek(offsetInstructions)
	for br.Position < offsetAttrDescriptors {
		command := Command{}
		command.Id = br.ReadByte()
		command.NumAttributes = br.ReadByte()
		command.Offset = br.ReadByte()
		br.Skip(1)

		commands.Add(command)
	}

	script.Attributes = attributes.Items
	script.Commands = commands.Items

	return script
}

func (attr *Attribute) Decompile() string {
	stack := dsa.NewStack[string]()
	br := utils.NewBinaryReader(attr.Bytes)

	decompileAttribute(br, stack)

	return stack.Pop()
}

func decompileAttribute(br *utils.BinaryReader, stack *dsa.Stack[string]) {
	if len(br.Bytes) == br.Position {
		return;
	}

	opcode := br.ReadByte()
	argLength := br.ReadInt16()

	switch opcode {
	case 0x21:	// not equal
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s != %s", first, second)
		stack.Push(result)
	case 0x25:	// modulo
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s %% %s", first, second)
		stack.Push(result)
	case 0x26:	// logical and
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s && %s", first, second)
		stack.Push(result)
	case 0x29:	// end var index
		br.Skip(1)
		values := dsa.NewList[string]()
		for value := stack.Pop(); value != "("; value = stack.Pop() {
			values.Add(value)
		}

		varName := stack.Pop()
		
		b := utils.NewStringBuilder()
		for i := 0; i < values.Count(); i++ {
			b.Append(values.Items[i])
		}

		result := fmt.Sprintf("%s(%s)", varName, b.ToString())
		stack.Push(result)
	case 0x2a:	// multiply
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s * %s", first, second)
		stack.Push(result)
	case 0x2b:	// add
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s + %s", first, second)
		stack.Push(result)
	case 0x2c:	// array separator
		stack.Push(", ")
	case 0x2d:	// subtract
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s - %s", first, second)
		stack.Push(result)
	case 0x2f:	// divide
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s / %s", first, second)
		stack.Push(result)
	case 0x3c:	// less than
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s < %s", first, second)
		stack.Push(result)
	case 0x3d:	// equal
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s == %s", first, second)
		stack.Push(result)
	case 0x3e:	// greater than
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s > %s", first, second)
		stack.Push(result)
	case 0x41:	// binary and
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s & %s", first, second)
		stack.Push(result)
	case 0x42:	// int8
		number := br.ReadByte()

		result := fmt.Sprintf("%d", number)
		stack.Push(result)
	case 0x46:	// double
		number := br.ReadDouble()

		result := fmt.Sprintf("%f", number)
		stack.Push(result)
	case 0x48:	// variable
		prefix := br.ReadString(1)
		varId := br.ReadInt16()

		result := fmt.Sprintf("%svar%x", prefix, varId)
		stack.Push(result)
	case 0x49:	// int32
		number := br.ReadInt32()

		result := fmt.Sprintf("%d", number)
		stack.Push(result)
	case 0x4c:	// int64
		number := br.ReadInt64()

		result := fmt.Sprintf("%d", number)
		stack.Push(result)
	case 0x4d:	// string
		result := br.ReadString(int(argLength))
		stack.Push(result)
	case 0x4f:	// binary or
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s | %s", first, second)
		stack.Push(result)
	case 0x52:	// change sign
		item := stack.Pop()
		result := fmt.Sprintf("-%s", item)
		stack.Push(result)
	case 0x53:	// less than or equal
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s <= %s", first, second)
		stack.Push(result)
	case 0x56:	// start var index
		prefix := br.ReadString(1)
		varId := br.ReadInt16()

		result := fmt.Sprintf("%svar%x", prefix, varId)
		stack.Push(result)
		stack.Push("(")
	case 0x57:	// int16
		number := br.ReadInt16()

		result := fmt.Sprintf("%d", number)
		stack.Push(result)
	case 0x5a:	// greater than or equal
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s >= %s", first, second)
		stack.Push(result)
	case 0x5e:	// binary xor
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s ^ %s", first, second)
		stack.Push(result)
	case 0x69:	// to number
		item := stack.Pop()

		result := fmt.Sprintf("@(%s)", item)
		stack.Push(result)
	case 0x73:	// to string
		item := stack.Pop()

		result := fmt.Sprintf("$(%s)", item)
		stack.Push(result)
	case 0x76:	// array var
		prefix := br.ReadString(1)
		varId := br.ReadInt16()

		result := fmt.Sprintf("%svar%x()", prefix, varId)
		stack.Push(result)
	case 0x7c:	// logical or
		second := stack.Pop()
		first := stack.Pop()

		result := fmt.Sprintf("%s || %s", first, second)
		stack.Push(result)
	}

	decompileAttribute(br, stack)
}

func decrypt(br *utils.BinaryReader, key uint32) {
	if (key == 0) {
		return;
	}

	repeatedKey := []byte{
		byte(key),
		byte(key >> 8),
		byte(key >> 16),
		byte(key >> 24),
		byte(key),
		byte(key >> 8),
		byte(key >> 16),
		byte(key >> 24),
	}

	offsetData := 0x20
	for offsetSize := 0xC; offsetSize < 0x1C; offsetSize += 4 {
		br.Seek(offsetSize)
		size := br.ReadInt32()

		xor(br, offsetData, size, repeatedKey)
		offsetData += size
	}
}

func xor(br *utils.BinaryReader, offset int, length int, key []byte) {
	if (offset < 0 || length < 0 || offset + length > len(br.Bytes)) {
		log.Fatal("Script data cannot be decrypted with given key.")
	}

	offsetKey := 0
	for length > 0 {
		lenChunk := min(length, len(key) - offsetKey)
		if lenChunk >= 8 {
			lenChunk = 8
		} else if lenChunk >= 4 {
			lenChunk = 4
		} else {
			lenChunk = 1
		}

		for i := 0; i < lenChunk; i++ {
			br.Bytes[offset + i] ^= key[offsetKey + i]
		}

		offset += lenChunk
		length -= lenChunk
		offsetKey += lenChunk
		if offsetKey == len(key) {
			offsetKey = 0
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}