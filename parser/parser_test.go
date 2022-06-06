package parser

import (
	"fmt"
	"testing"

	"github.com/JaeguKim/dag-go/dag"
)

func TestProcessXML(t *testing.T) {
	d := serveXml()
	decoder := newDecoder(d)
	_, xmlNodes, _ := processXML(decoder)

	for _, node := range xmlNodes {
		fmt.Println("Node Id: ", node.Id)
		for _, t := range node.To {
			fmt.Println("To", t)
		}
		for _, f := range node.From {
			fmt.Println("From", f)
		}
		fmt.Println("Command ", node.Command)

	}
}

func TestProcessXml(t *testing.T) {
	xmls := getTestXmlList()
	num := len(xmls)
	if num == 0 {
		fmt.Println("값이 없음")
		return
	}

	for _, xml := range xmls {
		d := []byte(xml)

		decoder := newDecoder(d)
		_, xmlNodes, _ := processXML(decoder)

		for _, node := range xmlNodes {
			fmt.Println("Node Id: ", node.Id)
			for _, t := range node.To {
				fmt.Println("To", t)
			}
			for _, f := range node.From {
				fmt.Println("From", f)
			}
			fmt.Println("Command ", node.Command)

		}
	}

}

// 여러 형태의 xml 테스트 필요, 깨진 xml 로도 테스트 진행 필요.
// TODO nodes id 를 고유하게 만들어줘서, 이걸 dag id 로 넣어주자.
func serveXml() []byte {
	// id, to, from, command
	// id 는 attribute, to, from, command 는 tag
	// to, from 은 복수 가능.
	// from 이 없으면 시작노드, 파싱된 후에 start_node, end_node 추가 됨.
	xml := `
	<nodes>
		<node id = "1">
			<to>2</to>
			<command> echo "hello world 1"</command>
		</node>
		<node id = "2" >
			<from>1</from>
			<to>3</to>
			<to>8</to>
			<command>echo "hello world 2"</command>
		</node>
		<node id ="3" >
			<from>2</from>
			<to>4</to>
			<to>5</to>
			<command>echo "hello world 3"</command>
		</node>
		<node id ="4" >
			<from>3</from>
			<to>6</to>
			<command>echo "hello world 4"</command>
		</node>
		<node id ="5" >
			<from>3</from>
			<to>6</to>
			<command>echo "hello world 5"</command>
		</node>
		<node id ="6" >
			<from>4</from>
			<from>5</from>
			<to>7</to>
			<command>echo "hello world 6"</command>
		</node>
		<node id ="7" >
			<from>6</from>
			<to>9</to>
			<to>10</to>
			<command>echo "hello world 7"</command>
		</node>
		<node id ="8" >
			<from>2</from>
			<command>echo "hello world 8"</command>
		</node>
		<node id ="9" >
			<from>7</from>
			<to>11</to>
			<command>echo "hello world 9"</command>
		</node>
		<node id ="10" >
			<from>7</from>
			<to>11</to>
			<command>echo "hello world 10"</command>
		</node>
		<node id ="11" >
			<from>9</from>
			<from>10</from>
			<command>echo "hello world 11"</command>
		</node>
	</nodes>`

	return []byte(xml)
}

func getTestXmlList() []string {
	var xs []string

	// nodes 가 없는 경우
	failedXml3 := `
		<node id = "1" >
			<to>2</to>
			<command> echo "hello world 1"</command>
		</node>
		<node id = "2" >
			<from>1</from>
			<to>3</to>
			<command>echo "hello world 2"</command>
		</node>
		<node id ="3" >
			<from>2</from>
			<command>echo "hello world 3"</command>
		</node>`

	//xs = append(xs, failedXml1)
	//xs = append(xs, failedXml2)
	xs = append(xs, failedXml3)
	//xs = append(xs, failedXml4)

	return xs
}

func xmlss() {
	// id 가 없는 node
	failedXml1 := `
	<nodes>
		<node>
			<to>2</to>
			<command> echo "hello world 1"</command>
		</node>
		<node id = "2" >
			<from>1</from>
			<to>3</to>
			<command>echo "hello world 2"</command>
		</node>
		<node id ="3" >
			<from>2</from>
			<command>echo "hello world 3"</command>
		</node>
	</nodes>`

	// node id 가 중복된 경우
	failedXml2 := `
	<nodes>
		<node id = "2" >
			<to>2</to>
			<command> echo "hello world 1"</command>
		</node>
		<node id = "2" >
			<from>1</from>
			<to>3</to>
			<command>echo "hello world 2"</command>
		</node>
		<node id ="3" >
			<from>2</from>
			<command>echo "hello world 3"</command>
		</node>
	</nodes>`

	// nodes 가 없는 경우
	failedXml3 := `
		<node id = "1" >
			<to>2</to>
			<command> echo "hello world 1"</command>
		</node>
		<node id = "2" >
			<from>1</from>
			<to>3</to>
			<command>echo "hello world 2"</command>
		</node>
		<node id ="3" >
			<from>2</from>
			<command>echo "hello world 3"</command>
		</node>`

	// circle TODO 다수 테스트 진행해야함.
	failedXml4 := `
	<nodes>
		<node id = "1">
			<to>2</to>
			<command> echo "hello world 1"</command>
		</node>
		<node id = "2" >
			<from>1</from>
			<to>3</to>
			<to>1</to>
			<command>echo "hello world 2"</command>
		</node>
		<node id ="3" >
			<from>2</from>
			<command>echo "hello world 3"</command>
		</node>
	</nodes>`

	fmt.Println(failedXml1)
	fmt.Println(failedXml2)
	fmt.Println(failedXml3)
	fmt.Println(failedXml4)
}

// parser 테스트
// TODO 버그 있음. read |0: file already closed
// generateDag 에서 입력 파라미터 포인터로 할지 생각하자.
func TestXmlParser(t *testing.T) {
	d := serveXml()
	decoder := newDecoder(d)
	_, xmlNodes, xmlNodeMap := processXML(decoder)
	generateDag(xmlNodes, xmlNodeMap)
}

func TestCreateNodeFromXmlNode(t *testing.T) {
	newDag := dag.NewDag("testRootDagId")
	createNodeFromXmlNode(newDag,&XmlNode{Id: "testXmlNodeId", Command: "echo test"})
	if newDag.GetNode("testXmlNodeId") == nil {
		t.Errorf("error")
	}
}

func TestAddEdgeFromXmlNode(t *testing.T) {
	newDag := dag.NewDag("testRootDagId")
	addEdgeFromXmlNode(newDag,&XmlNode{Id: "1"},&XmlNode{Id: "2"})
	if newDag.GetEdge("1-2") == nil {
		t.Errorf("error")
	}
}

func TestAddEdgeFromStartNodeToXmlNode(t *testing.T) {
	newDag := dag.NewDag("testRootDagId")
	addEdgeFromStartNodeToXmlNode(newDag,&XmlNode{Id: "1"})
	edgeKey := fmt.Sprintf("%s-1",newDag.StartNode.Id)
	if newDag.GetEdge(edgeKey) == nil {
		t.Errorf("error")
	}
}