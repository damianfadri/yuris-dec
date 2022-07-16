package main

import (
	"fmt"
	"os"
	"bufio"
	"log"
	
	"github.com/damianfadri/yuris-decompiler/utils/dsa"
	"github.com/damianfadri/yuris-decompiler/yuris"
)

func main() {
	script := yuris.ReadYST("G:\\Erewhon\\02_extracted\\ysbin\\yst00042.ybn")	
	labels := yuris.ReadYSL("G:\\Erewhon\\02_extracted\\ysbin\\ysl.ybn")
	compiler := yuris.ReadYSCom("G:\\Code\\github-yuris-decompiler\\yu-ris_0481_002\\システム\\system\\YSCom\\YSCom.ycd")

	// Get labels for the current script
	scriptLabels := dsa.NewList[yuris.Label]()
	for i := 0; i < len(labels); i++ {
		label := labels[i]
		if (label.ScriptIndex == 42) {
			scriptLabels.Add(label)
		}
	}

	iterCommands := dsa.NewIterator[yuris.Command](script.Commands)
	iterAttributes := dsa.NewIterator[yuris.Attribute](script.Attributes)
	iterLabels := dsa.NewIterator[yuris.Label](scriptLabels.Items)
	
	stack := dsa.NewStack[yuris.Line]()

	label := iterLabels.Next()
	command := iterCommands.Next()
	commandCount := 0
	for iterCommands.HasNext() {
		item := yuris.Line{}

		// TODO: Handle labels without return
		if label != nil && commandCount == label.Offset {
			args := dsa.NewList[string]()
			args.Add(label.Name)

			names := dsa.NewList[string]()
			names.Add("LabelName")

			item.Command = "LABEL"
			item.Arguments = args.Items
			item.Names = names.Items

			stack.Push(item)

			label = iterLabels.Next()
		} else {
			commandName := compiler.Commands[command.Id]
			item.Command = commandName
	
			names := dsa.NewList[string]()
			args := dsa.NewList[string]()
	
			// Set current line value.
			switch commandName {
			case "IF":
				fallthrough
			case "ELSE":
				fallthrough
			case "LOOP":
				if command.NumAttributes == 0 {
					break
				}
				attribute := iterAttributes.Next()
				conditionAttr := compiler.Attributes[command.Id][byte(attribute.Id)]
				conditionValue := attribute.Decompile()
	
				names.Add(conditionAttr)
				args.Add(conditionValue)
	
				// Skip the rest of the attributes
				for i := 1; i < int(command.NumAttributes); i++ {
					iterAttributes.Next()
				}
			case "LET":	
				varNameAttr := iterAttributes.Next()
				varName := varNameAttr.Decompile()
	
				varValueAttr := iterAttributes.Next()
				varValue := varValueAttr.Decompile()
	
				varOperation := "="
				switch varNameAttr.Type[1] {
				case byte(1):
					varOperation = "+="
				case byte(2):
					varOperation = "-="
				}
	
				names.Add("Operand1")
				args.Add(varName)

				names.Add("Operation")
				args.Add(varOperation)

				names.Add("Operand2")
				args.Add(varValue)
			default:
				for i := 0; i < int(command.NumAttributes); i++ {
					attribute := iterAttributes.Next()
					attrName := compiler.Attributes[command.Id][byte(attribute.Id)]
					attrValue := attribute.Decompile()
	
					names.Add(attrName)
					args.Add(attrValue)
				}
			}
	
			item.Names = names.Items
			item.Arguments = args.Items
			
			switch commandName {
			case "RETURN":
				curr := item
				children := dsa.NewList[yuris.Line]()
				for stack.Count() > 0 && (curr.Command != "LABEL" || curr.Visited) {
					children.Add(curr)
					curr = stack.Pop()
				}
	
				children.Reverse()
	
				start := curr
				start.Visited = true
				start.Children = children.Items
				stack.Push(start)
			case "RETURNCODE":
				curr := item
				children := dsa.NewList[yuris.Line]()
				for stack.Count() > 0 && (curr.Command != "WORD" || curr.Visited) {
					children.Add(curr)
					curr = stack.Pop()
				}
	
				children.Reverse()
	
				start := curr
				start.Visited = true
				start.Children = children.Items
				stack.Push(start)
			case "IFBLEND":
				fallthrough
			case "IFEND":
				curr := stack.Pop()
				children := dsa.NewList[yuris.Line]()
				for stack.Count() > 0 && (curr.Command != "IF" || curr.Visited) && (curr.Command != "ELSE" || curr.Visited) {
					children.Add(curr)
					curr = stack.Pop()
				}
	
				children.Reverse()

				start := curr
				start.Visited = true
				start.Children = children.Items
				stack.Push(start)
	
				if item.Command == "IFEND" {
					stack.Push(item)
				}
			case "LOOPEND":
				curr := stack.Pop()
				children := dsa.NewList[yuris.Line]()
				for stack.Count() > 0 && (curr.Command != "LOOP" || curr.Visited) {
					children.Add(curr)
					curr = stack.Pop()
				}
	
				start := curr
				start.Visited = true
				start.Children = children.Items
				stack.Push(start)
	
				if item.Command == "LOOPEND" {
					stack.Push(item)
				}
			default:
				stack.Push(item)
			}

			command = iterCommands.Next()
			commandCount += 1
		}
	}

	lines := dsa.NewList[yuris.Line]()

	for stack.Count() > 0 {
		item := stack.Pop()
		lines.Add(item)
	}

	lines.Reverse()

	file, err := os.Create("G:\\Code\\github-yuris-decompiler\\lines.txt")
	if err != nil {
		log.Fatal("")
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for i := 0; i < lines.Count(); i++ {
		fmt.Fprintln(w, lines.Items[i].ToString(0))
	}

	w.Flush()
}