package dag

import (
	"fmt"
	"testing"
	"time"

	"github.com/JaeguKim/dag-go/config"
	"github.com/JaeguKim/dag-go/node"
)

func TestCreateEdge(t *testing.T) {
	dag := NewDag("testRootDagId")
	res := dag.CreateEdge("1","2")
	if res == config.Exist {
		t.Errorf("error")
	}
	res = dag.CreateEdge("1","2")
	if res != config.Exist {
		t.Errorf("error")
	}
}

func TestGetEdgeChannel(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1","2")
	v := dag.GetEdgeChannel("1","2")
	if v == nil {
		t.Errorf("getVertext returned wrong value: got nil, expected not nil")
	}
	v = dag.GetEdgeChannel("1","3")
	if v != nil {
		t.Errorf("getEdgeChannel returned wrong value: got not nil, expected nil")
	}
}

func TestCreateNode(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.createNode("1")
	if dag.createNode("1") != nil || dag.GetNode("1") == nil{
		t.Errorf("error")
	}
}

func TestAddEdge(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	if dag.GetEdge("1-2") == nil {
		t.Errorf("error")
	}
}

func TestDag_AddEndNode(t *testing.T) {
	dag := NewDag("testRootDagId")
	node := dag.createNode("testNode")
	dag.AddEndNode(node)
	edgeKey := fmt.Sprintf("%s-%s",node.Id,dag.endNode.Id)
	if dag.GetEdge(edgeKey) == nil {
		t.Errorf("error")
	}
}

func TestDag_IsCyclic1(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.AddEdge("0","1")
	dag.AddEdge("0","2")
	dag.AddEdge("1","2")
	dag.AddEdge("2","0")
	dag.AddEdge("2","3")
	dag.AddEdge("3","3")
	if dag.IsCyclic() == false {
		t.Errorf("error")
	}
}

func TestDag_IsCyclic2(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.AddEdge("0","1")
	dag.AddEdge("0","2")
	dag.AddEdge("2","3")
	if dag.IsCyclic() {
		t.Errorf("error")
	}
}

// circle 이면 true, 아니면 false
func TestDag_IsCyclic3(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("2", "3")
	dag.AddEdge("3", "1")

	result := dag.IsCyclic()

	if result == false {
		t.Errorf("Error")
	}
}

func TestDag_IsCyclic4(t *testing.T) {
	//var result bool = false

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", dag.StartNode.Id)

	result := dag.IsCyclic()
	if result == false {
		t.Errorf("Error")
	}
}

func TestDag_IsCyclic5(t *testing.T) {

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge(dag.StartNode.Id, "2")
	result := dag.IsCyclic()
	if result == true {
		t.Errorf("Error")
	}
}

func TestDag_IsCyclic6(t *testing.T) {

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("1", "3")
	dag.AddEdge("1", "4")
	result := dag.IsCyclic()
	if result == true {
		t.Errorf("Error")
	}
}

func TestDag_GetNextNode(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "2")
	dag.AddEdge("2", "3")
	dag.AddEdge("2", "4")
	dag.AddEdge("3", "5")
	// 첫번째 자식노드를 리턴함.
	first := dag.Nodes[dag.StartNode.Id]
	// n 은 "2"
	n := node.GetNextNode(first)
	s := dag.Nodes[dag.StartNode.Id]
	fmt.Println("첫번째 자식노드(startNode) 수:", len(s.Children), "자식노드의 이름:", n.Id)

	// nn 은 3 또는 4
	nn := node.GetNextNode(n)
	ss := dag.Nodes[n.Id]
	fmt.Println(n.Id, "자식노드 수:", len(ss.Children), "자식노드의 이름:", nn.Id)
	// nnn 은 3 또는 4
	// 여기까지는 nil 이 안나오기 때문에 nil check 는 하지 않음.
	nnn := node.GetNextNode(n)
	fmt.Println(n.Id, "자식노드 수:", len(ss.Children), "자식노드의 이름:", nnn.Id)

	// nnnn 은 nil 이어야 한다.
	n4 := node.GetNextNode(n)
	if n4 == nil {
		fmt.Println("nil 정상")
	}

}

// TODO 다양한 파이프라인 테스트
// 파이프라인 테스트가 끝나고 메서드 및 함수에 대한 정리를 진행하자.
func TestDag_PipelineTest(t *testing.T) {
	//var result bool = false

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("2", "3")
	dag.AddEdge("3", "4")

	//visited := dag.visitReset()

	// 처음 시작은 start_node, start_node 로 해줘야 함.
	result := dag.IsCyclic()

	dag.DagSetFunc()
	dag.SetNodeToReadyState()
	dag.start()

	time.Sleep(10 * time.Second)
	fmt.Println("모든 고루틴이 종료될때까지 그냥 기다림.")

	fmt.Printf("%t", result)
}

func TestDag_PipelineTest1(t *testing.T) {
	//var result bool = false

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("1", "3")
	dag.AddEdge("1", "4")

	//visited := dag.visitReset()

	// 처음 시작은 start_node, start_node 로 해줘야 함.
	result := dag.IsCyclic()

	dag.DagSetFunc()
	dag.SetNodeToReadyState()
	dag.start()

	time.Sleep(10 * time.Second)
	fmt.Println("모든 고루틴이 종료될때까지 그냥 기다림.")

	fmt.Printf("%t", result)
}

func TestDag_PipelineTest2(t *testing.T) {
	//var result bool = false

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("1", "3")
	dag.AddEdge("2", "4")

	//visited := dag.visitReset()

	// 처음 시작은 start_node, start_node 로 해줘야 함.
	result := dag.IsCyclic()

	dag.DagSetFunc()
	dag.SetNodeToReadyState()
	dag.start()

	time.Sleep(10 * time.Second)
	fmt.Println("모든 고루틴이 종료될때까지 그냥 기다림.")

	fmt.Printf("%t", result)
}

func TestDag_PipelineTest3(t *testing.T) {
	//var result bool = false

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("1", "3")
	dag.AddEdge("1", "6")
	dag.AddEdge("2", "4")
	dag.AddEdge("3", "4")
	dag.AddEdge("4", "8")
	dag.AddEdge("4", "5")
	dag.AddEdge("4", "7")
	dag.AddEdge("5", "9")
	dag.AddEdge("7", "9")
	dag.AddEdge("9", "10")

	//visited := dag.visitReset()

	// 처음 시작은 start_node, start_node 로 해줘야 함.
	result := dag.IsCyclic()

	dag.DagSetFunc()
	dag.SetNodeToReadyState()
	dag.start()

	time.Sleep(10 * time.Second)
	fmt.Println("모든 고루틴이 종료될때까지 그냥 기다림.")

	fmt.Printf("%t", result)
}

// TestDag_PipelineTest3 와 동일하지만 AddEdge 의 순서를 바꿨음.
func TestDag_PipelineTest4(t *testing.T) {
	//var result bool = false

	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("4", "5")
	dag.AddEdge("4", "7")
	dag.AddEdge("5", "9")
	dag.AddEdge("7", "9")
	dag.AddEdge("9", "10")
	dag.AddEdge("1", "2")
	dag.AddEdge("1", "3")
	dag.AddEdge("2", "4")
	dag.AddEdge("3", "4")
	dag.AddEdge("4", "8")
	dag.AddEdge("1", "6")

	//visited := dag.visitReset()

	// 처음 시작은 start_node, start_node 로 해줘야 함.
	result := dag.IsCyclic()

	dag.DagSetFunc()
	dag.SetNodeToReadyState()
	dag.start()

	time.Sleep(10 * time.Second)
	fmt.Println("모든 고루틴이 종료될때까지 그냥 기다림.")

	fmt.Printf("%t", result)
}

// 각 노드의 channel 이 동일한지 또는 리셋되었는지 확인한다.
func TestChannelTest(t *testing.T) {
	dag := NewDag("testRootDagId")
	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")

	n1 := dag.Nodes[dag.StartNode.Id]
	n2 := dag.Nodes["1"]
	n3 := dag.Nodes["2"]
	// n1
	num1CHc := len(n1.ChildrenVertex)
	if num1CHc != 1 {
		fmt.Println("실패1")
	}
	// n2
	num2CHp := len(n2.ParentVertex)
	if num2CHp != 1 {
		fmt.Println("실패2 p")
	}

	num2CHc := len(n2.ChildrenVertex)
	if num2CHc != 1 {
		fmt.Println("실패2 c")
	}
	// n3

	num3CHp := len(n3.ParentVertex)
	if num3CHp != 1 {
		fmt.Println("실패3 p")
	}

	num3CHc := len(n3.ChildrenVertex)
	if num3CHc != 0 {
		fmt.Println("실패3 c")
	}

	// channel 이 같은지 검사
	t00 := n1.ChildrenVertex[0]
	t01 := n2.ParentVertex[0]

	if t00 != t01 {
		fmt.Println("실패 t00, t01")
	}

	t10 := n2.ChildrenVertex[0]
	t11 := n3.ParentVertex[0]

	if t10 != t11 {
		fmt.Println("실패 t10, t11")
	}
	// buffer channel 이라 그냥 이렇게 테스트 함.
	t00 <- 1
	v0 := <-t01

	fmt.Println(v0)

	t10 <- 1
	v1 := <-t11
	fmt.Println(v1)

}

func TestRunningStatus(t *testing.T) {
	dag := NewDag("testRootDagId")
	//go dag.printRunningStatus()

	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("1", "3")
	dag.AddEdge("1", "4")

	//visited := dag.visitReset()

	// 처음 시작은 start_node, start_node 로 해줘야 함.
	result := dag.IsCyclic()

	dag.DagSetFunc()
	dag.SetNodeToReadyState()

	//dag.SetFunc()
	//dag.getReady()
	dag.start()

	//dag.printRunningStatus()
	time.Sleep(10 * time.Second)
	fmt.Println("모든 고루틴이 종료될때까지 그냥 기다림.")

	fmt.Printf("%t", result)
	fmt.Println("false 이면 정상")
}

func TestFinishDag(t *testing.T) {
	dag := NewDag("testRootDagId")

	dag.AddEdge(dag.StartNode.Id, "1")
	dag.AddEdge("1", "2")
	dag.AddEdge("1", "3")
	dag.AddEdge("1", "4")

	dag.FinishDag()
	dag.DagSetFunc()
	dag.SetNodeToReadyState()

	dag.start()

	b := dag.waitTilOver(nil)
	fmt.Println("모든 고루틴이 종료될때까지 그냥 기다림.")
	if b == false {
		t.Errorf("error")
	}
}



// TODO 비정상 graph 도 테스트 진행해야함.
