package dag

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDot_opts(t *testing.T) {
	var v testDotVertex
	var g GraphBase
	g.Add(&v)

	opts := &DotOpts{MaxDepth: 42}
	actual := g.Dot(opts)
	if len(actual) == 0 {
		t.Fatal("should not be empty")
	}

	if !v.DotNodeCalled {
		t.Fatal("should call DotNode")
	}
	if !reflect.DeepEqual(v.DotNodeOpts, opts) {
		t.Fatalf("bad; %#v", v.DotNodeOpts)
	}
}

type testDotVertex struct {
	DotNodeCalled bool
	DotNodeTitle  string
	DotNodeOpts   *DotOpts
	DotNodeReturn *DotNode
}

func (v *testDotVertex) DotNode(title string, opts *DotOpts) *DotNode {
	v.DotNodeCalled = true
	v.DotNodeTitle = title
	v.DotNodeOpts = opts
	return v.DotNodeReturn
}

func TestGraphDot_MultiGraph(t *testing.T) {
	graph := createConnectedMultiSubgraph()

	marshaledGraph := graph.Marshal(&MarshalOpts{
		ConnectSubgraphHeads: true,
		ConnectSubgraphTails: true,
	})

	jsonGraph, err := json.MarshalIndent(marshaledGraph, "", "  ")
	assert.NoError(t, err)

	fmt.Println(string(jsonGraph))

	dot := marshaledGraph.Dot(&DotOpts{})

	expectedDot := `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] itemOne" -> "[root] itemTwo"
		"[root] itemTwo" -> "[subgraphOne] itemThree"
		"[subgraphOne] itemFour" -> "[root] itemFive"
	}
	subgraph "subgraphOne" {
		"[subgraphOne] itemThree" -> "[subgraphOne] itemFour"
	}
}
`

	fmt.Println(string(dot))

	assert.Equal(t, expectedDot, string(dot))
}
