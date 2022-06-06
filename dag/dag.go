package dag

import (
	"context"
	"fmt"
	"strings"

	"github.com/JaeguKim/dag-go/config"
	"github.com/JaeguKim/dag-go/node"
)

type Dag struct {
	Id            string
	Nodes         map[string]*node.Node
	edges         map[string]*Edge
	StartNode     *node.Node
	endNode       *node.Node
	validated     bool
	RunningStatus chan config.PrintStatus
}

type Edge struct {
	parentId string
	childId  string
	ch       chan int
}

func NewDag(dagId string) *Dag {
	dag := new(Dag)
	dag.Nodes = make(map[string]*node.Node)
	dag.edges = make(map[string]*Edge)
	dag.Id = dagId
	dag.validated = false
	dag.StartNode = dag.createNode(config.StartNode)
	dag.endNode = dag.createNode(config.EndNode)
	if dag.StartNode == nil {
		return nil
	}

	dag.StartNode.ParentVertex = append(dag.StartNode.ParentVertex, make(chan int, 1))

	dag.RunningStatus = make(chan config.PrintStatus, 1000)

	return dag
}

func (dag *Dag) CreateEdge(parentId, childId string) config.CreateEdgeErrorType {

	if len(strings.TrimSpace(parentId)) == 0 {
		return config.Fault
	}
	if len(strings.TrimSpace(childId)) == 0 {
		return config.Fault
	}
	edgeKey := fmt.Sprintf("%s-%s", parentId, childId)
	if dag.edges[edgeKey] != nil {
		return config.Exist
	}

	edge := new(Edge)
	edge.parentId = parentId
	edge.childId = childId
	edge.ch = make(chan int, 1)

	dag.edges[edgeKey] = edge

	return config.Create
}

func (dag *Dag) GetEdgeChannel(parentId, childId string) chan int {
	edgeKey := fmt.Sprintf("%s-%s", parentId, childId)
	if dag.edges[edgeKey] == nil {
		return nil
	}
	return dag.edges[edgeKey].ch
}

func (dag *Dag) createNode(id string) *node.Node {

	_, exists := dag.Nodes[id]
	if exists {
		return nil
	}

	newNode := &node.Node{Id: id}
	//newNode.ParentDag = dag
	dag.Nodes[id] = newNode

	return newNode
}

func (dag *Dag) GetNode(s string) *node.Node {
	if len(strings.TrimSpace(s)) == 0 {
		return nil
	}

	size := len(dag.Nodes)
	if size <= 0 {
		return nil
	}

	n := dag.Nodes[s]
	return n
}

func (dag *Dag) AddEdge(from, to string) error {
	edgeKey := fmt.Sprintf("%s-%s", from, to)
	check := dag.edges[edgeKey]
	if check != nil {
		return fmt.Errorf("edge already exists")
	}

	fromNode := dag.Nodes[from]
	if fromNode == nil {
		fromNode = dag.createNode(from)
	}
	toNode := dag.Nodes[to]
	if toNode == nil {
		toNode = dag.createNode(to)
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
		return fmt.Errorf("error")
	}

	return nil
}

func (dag *Dag) GetEdge(s string) *Edge {
	if len(strings.TrimSpace(s)) == 0 {
		return nil
	}

	size := len(dag.edges)
	if size <= 0 {
		return nil
	}

	e := dag.edges[s]
	return e
}

func (dag *Dag) AddEndNode(fromNode *node.Node) error {

	if fromNode == nil {
		return nil
	}
	if dag.endNode == nil {
		return nil
	}

	fromNode.Children = append(fromNode.Children, dag.endNode)
	dag.endNode.Parent = append(dag.endNode.Parent, fromNode)

	dag.CreateEdge(fromNode.Id, dag.endNode.Id)
	v := dag.GetEdgeChannel(fromNode.Id, dag.endNode.Id)

	if v != nil {
		fromNode.ChildrenVertex = append(fromNode.ChildrenVertex, v)
		dag.endNode.ParentVertex = append(dag.endNode.ParentVertex, v)
	} else {
		fmt.Println("error")
	}

	return nil
}

func (dag *Dag) FinishDag() error {
	dag.validated = dag.IsCyclic() == false
	if dag.validated == false {
		return nil
	}

	if len(dag.Nodes) == 0 {
		return fmt.Errorf("no vertex")
	}

	for _, n := range dag.Nodes {
		if len(n.Children) == 0 {
			if n.Id != config.EndNode {
				err := dag.AddEndNode(n)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (dag *Dag) isCycleUtil(key string, visited map[string]bool, recStack map[string]bool) bool {
	if recStack[key] {
		return true
	}
	if visited[key] {
		return false
	}
	visited[key] = true
	recStack[key] = true
	for _, v := range dag.Nodes[key].Children {
		if dag.isCycleUtil(v.Id, visited, recStack) {
			return true
		}
	}
	recStack[key] = false
	return false
}

func (dag *Dag) IsCyclic() bool {
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	for nodeId := range dag.Nodes {
		if dag.isCycleUtil(nodeId, visited, recStack) {
			return true
		}
	}
	return false
}

func (dag *Dag) DagSetFunc() {

	n := len(dag.Nodes)
	if n < 1 {
		return
	}

	for _, v := range dag.Nodes {
		node.SetFunc(v)
	}
}

func (dag *Dag) SetNodeToReadyState() bool {
	n := len(dag.Nodes)
	if n < 1 {
		return false
	}

	for _, v := range dag.Nodes {
		go v.Runner(v, dag.RunningStatus)
	}

	return true
}

func (dag *Dag) start() bool {
	n := len(dag.StartNode.ParentVertex)
	if n != 1 {
		return false
	}

	go func(c chan int) {
		ch := c
		ch <- 1
	}(dag.StartNode.ParentVertex[0])

	return true
}

// TODO context 넣어서 무한 루프 방지하자.
func (dag *Dag) waitTilOver(ctx context.Context) bool {
	// Wait V1
	//for {
	//	flag := true
	//	for _, node := range dag.nodes {
	//		if node.status != PostFlightEnd {
	//			time.Sleep(time.Second)
	//			flag = false
	//			break
	//		}
	//	}
	//	if flag {
	//		return true
	//	}
	//}

	for {
		select {
		case c := <-dag.RunningStatus:
			fmt.Printf("nodeId : %s, status : %d\n", c.NodeId, c.RStatus)
			if c.NodeId == config.EndNode {
				if c.RStatus == config.PreFlightFailed || c.RStatus == config.InFlightFailed || c.RStatus == config.PostFlightFailed {
					return false
				}
				if c.RStatus == config.PostFlightEnd {
					return true
				}
			}
		}
	}

}

func (dag *Dag) Start() {
	dag.start()
	dag.waitTilOver(nil)
}
