package parser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/JaeguKim/dag-go/config"
	"github.com/JaeguKim/dag-go/dag"
	"github.com/JaeguKim/dag-go/node"
)

type (
	XmlNode struct {
		Id   string
		From []string
		To      []string
		Command string
	}

	XmlNodes []*XmlNode
)

func InitWithXML(d []byte) (bool, *dag.Dag) {
	decoder := newDecoder(d)
	e, xmlNodes, xmlNodeMap := processXML(decoder)
	if e != nil {
		fmt.Println("processXML returned 0")
		return false, nil
	}
	b, newDag := generateDag(xmlNodes, xmlNodeMap)
	if b == false {
		fmt.Println("generateDag returned false")
		return false, nil
	}
	return b, newDag
}

func newDecoder(b []byte) *xml.Decoder {
	d := xml.NewDecoder(bytes.NewReader(b))
	return d
}

func processXML(parser *xml.Decoder) (error, XmlNodes, map[string]*XmlNode) {
	var (
		n  *XmlNode = nil
		ns XmlNodes = nil
		// TODO bool 로 하면 안됨 int 로 바꿔야 함. 초기 값은 0, false = 1, true = 2 로
		xStart    = false
		nStart    = false
		cmdStart  = false
		fromStart = false
		toStart   = false
	)
	if parser == nil {
		return fmt.Errorf("parser is null"), nil, nil
	}

	xmlNodeMap := make(map[string]*XmlNode)

	for {
		token, err := parser.Token()

		if err == io.EOF {
			break // TODO break 구문 수정 처리 필요.
		}

		// TODO 중복되는 것 같지만 일단 그냥 넣어둠.
		if token == nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			xmlTag := t.Name.Local
			if xmlTag == config.Nodes {
				xStart = true
			}
			if xmlTag == config.Node {
				nStart = true

				n = new(XmlNode)
				// 현재는 node 의 경우 속성이 1 이도록 강제함. 추후 수정할 필요가 있으면 수정.
				if len(t.Attr) != 1 {
					return fmt.Errorf("length of node attribute should be one"), nil, nil
				}
				if t.Attr[0].Name.Local != config.Id {
					return fmt.Errorf("node attribute should contain id"), nil, nil
				}

				if xmlNodeMap[t.Attr[0].Value] != nil {

					return fmt.Errorf("duplicated node id exists"), nil, nil
				}
				xmlNodeMap[t.Attr[0].Value] = n
				n.Id = t.Attr[0].Value
			}

			if xmlTag == config.Command {
				cmdStart = true
			}

			if xmlTag == config.From {
				fromStart = true
			}

			if xmlTag == config.To {
				toStart = true
			}

		case xml.EndElement:
			xmlTag := t.Name.Local
			if xmlTag == config.Nodes {
				if xStart {
					xStart = false
				} else {

					return fmt.Errorf("xml end element error"), nil, nil
				}
			}
			// TODO 중복 구문들 function 으로 만들자.
			if xmlTag == config.Node {
				if xStart && nStart {
					if n == nil {

						return fmt.Errorf("n should not be null"), nil, nil
					}
					ns = append(ns, n)
					nStart = false
					n = nil
				}
			}

			if xmlTag == config.From {
				if xStart && nStart {
					if n == nil {

						return fmt.Errorf("n should not be null"), nil, nil
					}
					fromStart = false
				}
			}

			if xmlTag == config.To {
				if xStart && nStart {
					if n == nil {

						return fmt.Errorf("n should not be null"), nil, nil
					}
					if toStart == false {

						return fmt.Errorf("toStart should be true"), nil, nil
					}
					toStart = false
				}
			}

			if xmlTag == config.Command {
				if xStart && nStart {
					if n == nil {

						return fmt.Errorf("n should not be null"), nil, nil
					}
					if cmdStart == false {

						return fmt.Errorf("cmdStart should be true"), nil, nil
					}
					cmdStart = false
				}
			}

		case xml.CharData:
			if nStart && n != nil {
				if fromStart {
					n.From = append(n.From, string(t))
				}
				if cmdStart {
					n.Command = string(t)
				}
				if toStart {
					n.To = append(n.To, string(t))
				}
			}
		}
	}

	return nil, ns, xmlNodeMap
}

func generateDag(xmlNodes XmlNodes, xmlNodeMap map[string]*XmlNode) (bool, *dag.Dag) {

	if xmlNodes == nil {
		return false, nil
	}

	if len(xmlNodes) <= 0 {
		return false, nil
	}

	newDag := dag.NewDag("testDagRootId")
	for _, xmlNode := range xmlNodes {
		rn := len(xmlNode.From)
		if rn == 0 {
			err := addEdgeFromStartNodeToXmlNode(newDag, xmlNode)
			if err != nil {
				return false, nil
			}
		}
		for _, nodeId := range xmlNode.To {
			err := addEdgeFromXmlNode(newDag, xmlNode, xmlNodeMap[nodeId])
			if err != nil {
				return false, nil
			}
		}
	}

	err := newDag.FinishDag()
	if err != nil {
		return false, nil
	}
	newDag.DagSetFunc()
	newDag.SetNodeToReadyState()
	return true, newDag
}

func addEdgeFromStartNodeToXmlNode(dag *dag.Dag, to *XmlNode) error {
	if to == nil {
		return fmt.Errorf("xmlNode is nil")
	}

	fromNode := dag.StartNode

	toNode := createNodeFromXmlNode(dag,to)
	if toNode == nil {
		return fmt.Errorf("node creation error")
	}
	//toNode := dag.Nodes[to.Id]
	//if toNode == nil {
	//	toNode = createNodeFromXmlNode(dag,to)
	//}

	if fromNode == toNode {
		return fmt.Errorf("from-node and to-node are same")
	}

	// 원은 허용하지 않는다.
	// TODO 향후 어떻게 고칠지 생각해 보자.
	// 일단은 주석처리 한다.
	/*if strings.Contains(toNode.Id, "start_node") {
		return fmt.Errorf("circle is not allowed.")
	}*/

	fromNode.Children = append(fromNode.Children, toNode)
	toNode.Parent = append(toNode.Parent, fromNode)

	//fromNode.outdegree++
	//toNode.indegree++

	//check := dag.createEdge(fromNode.Id, toNode.Id)
	dag.CreateEdge(fromNode.Id, toNode.Id)
	// 생성 시키면 0, 이 존재하면 1, 에러면 2
	// TODO 향후 에러코드 만들면 수정해야함.
	/*if check == 0 {
		fmt.Println("만들어줌.")
	}
	if check == 1 {
		fmt.Println("존재함")
	}

	if check == 2 {
		fmt.Println("에러")
	}*/

	v := dag.GetEdgeChannel(fromNode.Id, toNode.Id)

	if v != nil {
		fromNode.ChildrenVertex = append(fromNode.ChildrenVertex, v)
		toNode.ParentVertex = append(toNode.ParentVertex, v)
	} else {
		return fmt.Errorf("dag has duplicated edge")
	}

	return nil
}

func addEdgeFromXmlNode(dag *dag.Dag, from, to *XmlNode) error {

	if from == nil {
		return fmt.Errorf("nil")
	}

	if to == nil {
		return fmt.Errorf("nil")
	}

	fromNode := dag.Nodes[from.Id]
	if fromNode == nil {
		fromNode = createNodeFromXmlNode(dag,from)
	}
	toNode := dag.Nodes[to.Id]
	if toNode == nil {
		toNode = createNodeFromXmlNode(dag,to)
	}

	if fromNode == toNode {
		return fmt.Errorf("from-node and to-node are same")
	}

	fromNode.Children = append(fromNode.Children, toNode)
	toNode.Parent = append(toNode.Parent, fromNode)

	dag.CreateEdge(fromNode.Id, toNode.Id)

	v := dag.GetEdgeChannel(fromNode.Id, toNode.Id)

	if v != nil {
		fromNode.ChildrenVertex = append(fromNode.ChildrenVertex, v)
		toNode.ParentVertex = append(toNode.ParentVertex, v)
	} else {
		fmt.Println("error")
	}

	return nil
}

func createNodeFromXmlNode(dag *dag.Dag, xnode *XmlNode) *node.Node {

	if xnode == nil {
		return nil
	}

	_, exists := dag.Nodes[xnode.Id]
	if exists {
		return nil
	}

	newNode := &node.Node{Id: xnode.Id}
	newNode.Commands = xnode.Command
	//newNode.ParentDag = dag
	dag.Nodes[xnode.Id] = newNode

	return newNode
}