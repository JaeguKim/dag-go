package node

import (
	"strings"
	"sync"

	"github.com/JaeguKim/dag-go/config"
	"github.com/JaeguKim/dag-go/shellexecmd"
)

type Node struct {
	Id string

	Children  []*Node
	Parent    []*Node
	//ParentDag *dag.Dag

	Commands       string
	Status         config.RunningStatus
	ChildrenVertex []chan int
	ParentVertex   []chan int

	Runner func(n *Node, result chan config.PrintStatus)
}

// TODO 고루틴이 특정 시간을 넘어서면 타임아웃도 걸어야 한다. (pipeline.go 에서 파이프라인의 타임 아웃 설정값을 세팅해야 하고 이것을 여기서 적용해야한다.- 직접적용보단 간접적으로 취소를 할 수 있어야 한다.)
// TODO context 향후 적용, 취소할때 취소가 되어야 하기때문에 중요함.
// TODO 채널 close 하는 것.
func SetFunc(n *Node) {
	n.Runner = func(n *Node, result chan config.PrintStatus) {
		//defer close(result)  // TODO 처리 해야 함. 최종적으로 채널에 보내는 모든 작업이 끝나면 RunningStatus chan printStatus close 를 실행해줘야 함.

		r := preFlight(n)
		result <- r

		r = inFlight(n)
		result <- r

		r = postFlight(n)
		result <- r
	}
}

func preFlight(n *Node) config.PrintStatus {
	if n == nil {
		return config.PrintStatus{config.PreFlightFailed, n.Id}
	}
	i := len(n.ParentVertex)
	wg := new(sync.WaitGroup)
	for j := 0; j < i; j++ {
		wg.Add(1)
		go func(c chan int) {
			defer wg.Done()
			<-c
			close(c)

		}(n.ParentVertex[j])
	}
	wg.Wait()
	n.Status = config.PreFlightEnd
	return config.PrintStatus{RStatus: config.PreFlightEnd, NodeId: n.Id}
}

func inFlight(n *Node) config.PrintStatus {

	if n == nil {
		return config.PrintStatus{RStatus: config.InFlightFailed, NodeId: n.Id}
	}

	var bResult = false
	if n.Id == config.StartNode || n.Id == config.EndNode {
		bResult = true
	} else {
		if len(strings.TrimSpace(n.Commands)) == 0 {
			bResult = true
		} else {
			bResult = shellexecmd.Runner(n.Commands)
		}
	}

	if bResult {
		n.Status = config.InFlightEnd
		return config.PrintStatus{RStatus: config.InFlightEnd, NodeId: n.Id}
	} else {
		n.Status = config.InFlightFailed
		return config.PrintStatus{RStatus: config.InFlightFailed, NodeId: n.Id}
	}
}

func postFlight(n *Node) config.PrintStatus {

	if n == nil {
		return config.PrintStatus{RStatus: config.PostFlightFailed, NodeId: n.Id}
	}

	k := len(n.ChildrenVertex)
	for j := 0; j < k; j++ {
		n.ChildrenVertex[j] <- 1
	}
	n.Status = config.PostFlightEnd
	return config.PrintStatus{RStatus: config.PostFlightEnd, NodeId: n.Id}
}

// 제일 먼저 들어간 노드를 가져오고 해당 노드는 삭제한다.
func GetNextNode(n *Node) *Node {

	if n == nil {
		return nil
	}
	if len(n.Children) < 1 {
		return nil
	}

	ch := n.Children[0]
	n.Children = append(n.Children[:0], n.Children[1:]...)

	return ch
}
