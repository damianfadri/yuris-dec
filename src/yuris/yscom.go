package yuris

import (
	"io/ioutil"
	"log"

	"github.com/damianfadri/yuris-decompiler/utils"
)

type CompilerDefinition struct {
	Commands		map[byte]string
	Attributes		[]map[byte]string
}

func ReadYSCom(path string) CompilerDefinition {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	br := utils.NewBinaryReader(data)
	magic := br.ReadString(4)
	if magic != "YSCD" {
		log.Fatal("Invalid magic in ysc.ybn")
	}

	// YU-RIS version
	_ = br.ReadInt32()
	numCommands := br.ReadInt32()

	br.Skip(4)

	yscom := CompilerDefinition{}
	yscom.Commands = make(map[byte]string)

	for commandId := byte(0); commandId < byte(numCommands); commandId++ {
		commandName := br.ReadStringUntilNull()
		yscom.Commands[commandId] = commandName

		numAttrs := byte(0)
		numAttrs = br.ReadByte()
		commandAttrs := make(map[byte]string)
		for attrId := byte(0); attrId < numAttrs; attrId++ {
			attrName := br.ReadStringUntilNull()
			commandAttrs[attrId] = attrName

			br.Skip(4)
		}

		yscom.Attributes = append(yscom.Attributes, commandAttrs)
	}

	return yscom
}