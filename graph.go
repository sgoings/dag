package dag

type HasSubgraph interface {
	Subgraph() Graph
}

type Graph interface {
	String() string

	Connect(Edge)
	Edges() []Edge
	HasEdge(Edge) bool
	EdgesFrom(Vertex) []Edge
	EdgesTo(Vertex) []Edge
	DownEdges(Vertex) Set
	UpEdges(Vertex) Set
	RemoveEdge(Edge)

	Add(Vertex) Vertex
	Vertices() []Vertex
	HasVertex(Vertex) bool
	Remove(Vertex) Vertex
	Replace(original Vertex, replacement Vertex) bool

	Marshal(*MarshalOpts) *marshalGraph
	Dot(*DotOpts) []byte
}
