package yuris

import (
	"github.com/damianfadri/yuris-decompiler/utils"
)

type Line struct {
	Command		string
	Arguments	[]string
	Names		[]string
	Children	[]Line
	Visited		bool
}

func getIndent(count int) string {
	sb := utils.NewStringBuilder()
	for i := 0; i < count; i++ {
		sb.Append(" ")
	}

	return sb.ToString()
}

func (item *Line) ToStringSingle(indent int) string {
	sb := utils.NewStringBuilder()

	sb.Append(getIndent(indent))

	switch item.Command {
	case "LET":
		sb.Append(item.Arguments[0])
		sb.Append(" ")
		sb.Append(item.Arguments[1])
		sb.Append(" ")
		sb.Append(item.Arguments[2])
	case "LABEL":
		sb.Append("#=")
		sb.Append(item.Arguments[0])
	case "IF":
		sb.Append("IF")
		sb.Append("[")
		sb.Append(item.Arguments[0])
		sb.Append("]")
	case "ELSE":
		sb.Append("ELSE")
		sb.Append("[")
		if (len(item.Arguments) > 0) {
			sb.Append(item.Arguments[0])
		}
		sb.Append("]")
	case "LOOP":
		sb.Append("LOOP")
		sb.Append("[")
		if (item.Arguments[0] != "255") {
			sb.Append(item.Names[0])
			sb.Append(" = ")
			sb.Append(item.Arguments[0])
		}
		sb.Append("]")
	case "IFEND":
		fallthrough
	case "LOOPEND":
		fallthrough
	case "LOOPBREAK":
		fallthrough
	case "LOOPCONTINUE":
		fallthrough
	case "END":
		sb.Append(item.Command)
		sb.Append("[]")
	case "S_INT":
		fallthrough
	case "S_STR":
		fallthrough
	case "INT":
		fallthrough
	case "STR":
		sb.Append(item.Command)
		sb.Append("[")
		sb.Append(item.Arguments[0])
		if (item.Arguments[1] != "0") {
			sb.Append(" = ")
			sb.Append(item.Arguments[1])
		}
		sb.Append("]")
	default:
		sb.Append(item.Command)
		sb.Append("[")
		for i := 0; i < len(item.Arguments); i++ {
			sb.Append(item.Names[i])
			sb.Append("=")
			sb.Append(item.Arguments[i])

			if (i + 1 < len(item.Arguments)) {
				sb.Append(" ")
			}
		}
		sb.Append("]")
	}

	sb.Append("\n")
	return sb.ToString()
}

func (item *Line) ToString(indent int) string {
	if len(item.Children) == 0 {
		return item.ToStringSingle(indent)
	}

	sb := utils.NewStringBuilder()
	ind := getIndent(indent)

	sb.Append(item.ToStringSingle(indent))
	sb.Append(ind)
	sb.Append("{")
	sb.Append("\n")
	for i := 0; i < len(item.Children); i++ {
		child := item.Children[i]
		sb.Append(child.ToString(indent + 2))
	}
	sb.Append(ind)
	sb.Append("}")
	sb.Append("\n")

	return sb.ToString()
}