package dag

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDot_empty(t *testing.T) {
	var g GraphBase
	g.Add(1)
	g.Add(2)
	g.Add(3)

	actual := strings.TrimSpace(string(g.Dot(nil)))
	expected := strings.TrimSpace(testGraphDotEmptyStr)
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestGraphDot_basic(t *testing.T) {
	var g GraphBase
	g.Add(1)
	g.Add(2)
	g.Add(3)
	g.Connect(BasicEdge(1, 3))

	actual := strings.TrimSpace(string(g.Dot(nil)))
	expected := strings.TrimSpace(testGraphDotBasicStr)
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

func TestGraphDot_quoted(t *testing.T) {
	var g GraphBase
	quoted := `name["with-quotes"]`
	other := `other`
	g.Add(quoted)
	g.Add(other)
	g.Connect(BasicEdge(quoted, other))

	actual := strings.TrimSpace(string(g.Dot(nil)))
	expected := strings.TrimSpace(testGraphDotQuotedStr)
	if actual != expected {
		t.Fatalf("\ngot:   %q\nwanted %q\n", actual, expected)
	}
}

func TestGraphDot_attrs(t *testing.T) {
	var g GraphBase
	g.Add(&testGraphNodeDotter{
		Result: &DotNode{
			Name:  "foo",
			Attrs: map[string]string{"foo": "bar"},
		},
	})

	actual := strings.TrimSpace(string(g.Dot(nil)))
	expected := strings.TrimSpace(testGraphDotAttrsStr)
	if actual != expected {
		t.Fatalf("bad: %s", actual)
	}
}

type testGraphNodeDotter struct{ Result *DotNode }

func (n *testGraphNodeDotter) Name() string                      { return n.Result.Name }
func (n *testGraphNodeDotter) DotNode(string, *DotOpts) *DotNode { return n.Result }

const testGraphDotQuotedStr = `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] name[\"with-quotes\"]" -> "[root] other"
	}
}`

const testGraphDotBasicStr = `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] 1" -> "[root] 3"
	}
}
`

const testGraphDotEmptyStr = `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
	}
}`

const testGraphDotAttrsStr = `digraph {
	compound = "true"
	newrank = "true"
	subgraph "root" {
		"[root] foo" [foo = "bar"]
	}
}`

func TestGraph_MultiGraph(t *testing.T) {
	graph := createConnectedMultiSubgraph()

	marshaledGraph := graph.Marshal(&MarshalOpts{})

	jsonGraph, err := json.MarshalIndent(marshaledGraph, "", "  ")
	assert.NoError(t, err)

	fmt.Println(string(jsonGraph))

	assert.Equal(t, 4, len(marshaledGraph.Vertices))

	assert.Equal(t, "itemFive", marshaledGraph.Vertices[0].Name)
	assert.Equal(t, "itemOne", marshaledGraph.Vertices[1].Name)
	assert.Equal(t, "itemTwo", marshaledGraph.Vertices[2].Name)
	assert.Equal(t, "subgraphOne", marshaledGraph.Vertices[3].Name)

	assert.Equal(t, 1, len(marshaledGraph.Subgraphs))
	assert.Equal(t, "subgraphOne", marshaledGraph.Subgraphs[0].Name)

	assert.Equal(t, 2, len(marshaledGraph.Subgraphs[0].Vertices))
	assert.Equal(t, "itemFour", marshaledGraph.Subgraphs[0].Vertices[0].Name)
	assert.Equal(t, "itemThree", marshaledGraph.Subgraphs[0].Vertices[1].Name)

	assert.Equal(t, &marshalEdge{
		Name:   "itemOne|itemTwo",
		Source: marshaledGraph.Vertices[1].ID,
		Target: marshaledGraph.Vertices[2].ID,
		Attrs:  map[string]string{},
	}, marshaledGraph.Edges[0])

	assert.Equal(t, &marshalEdge{
		Name:   "itemTwo|subgraphOne",
		Source: marshaledGraph.Vertices[2].ID,
		Target: marshaledGraph.Vertices[3].ID,
		Attrs:  map[string]string{},
	}, marshaledGraph.Edges[1])

	assert.Equal(t, &marshalEdge{
		Name:   "itemThree|itemFour",
		Source: marshaledGraph.Subgraphs[0].Vertices[1].ID,
		Target: marshaledGraph.Subgraphs[0].Vertices[0].ID,
		Attrs:  map[string]string{},
	}, marshaledGraph.Subgraphs[0].Edges[0])
}
