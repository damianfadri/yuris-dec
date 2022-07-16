package yuris

import (
	"io/ioutil"
	"log"

	"github.com/damianfadri/yuris-decompiler/utils/dsa"
	"github.com/damianfadri/yuris-decompiler/utils"
)

type Label struct {
	Name				string
	Id					int
	Offset				int
	ScriptIndex			int16
}

func ReadYSL(path string) []Label {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	br := utils.NewBinaryReader(data)
	magic := br.ReadString(4)
	if magic != "YSLB" {
		log.Fatal("Invalid magic in ysl.ybn")
	}

	// YU-RIS version
	_ = br.ReadInt32()
	numLabels := br.ReadInt32()

	arrLabelRangeStartIndices := make([]int, 0x100)
	for j := 0; j < len(arrLabelRangeStartIndices); j++ {
		arrLabelRangeStartIndices[j] = br.ReadInt32()
	}

	labels := dsa.NewList[Label]()
	for j := 0; j < numLabels; j++ {
		label := Label{}
		
		lenName := br.ReadByte()
		label.Name = br.ReadString(int(lenName))
		label.Id = br.ReadInt32()
		label.Offset = br.ReadInt32()
		label.ScriptIndex = br.ReadInt16()
		br.Skip(2)

		labels.Add(label)
	}

	labels.Sort(func(i, j int) bool { 
		return labels.Items[i].Offset < labels.Items[j].Offset
	})

	return labels.Items
}