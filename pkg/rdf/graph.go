package rdf

// Node is a node (subject and/or object) in a rdf graph
type Node struct {
	Term        Term
	Edge        []*Edge
	InverseEdge []*Edge
}

// Edge is a edge (predicate) in a rdf graph
type Edge struct {
	Pred    Predicate
	Subject *Node
	Object  *Node
}

// Graph is a rdf grapgh containing nodes and edges
type Graph struct {
	Nodes map[string]*Node
	Edges []*Edge
}

// New creates a graph from an rdf triple slice
func New(triple []Triple) (graph Graph, err error) {
	graph.Nodes = make(map[string]*Node)
	for i := range triple {
		// object
		obj, ok := graph.Nodes[triple[i].Obj.String()]
		if !ok {
			obj = &Node{
				Term: triple[i].Obj,
			}
			graph.Nodes[obj.Term.String()] = obj
		}

		// subject
		subj, ok := graph.Nodes[triple[i].Sub.String()]
		if !ok {
			subj = &Node{
				Term: triple[i].Sub,
			}
			graph.Nodes[subj.Term.String()] = subj
		}

		// predicate
		edge := &Edge{
			Pred:    triple[i].Pred,
			Subject: subj,
			Object:  obj,
		}
		subj.Edge = append(subj.Edge, edge)
		obj.InverseEdge = append(obj.InverseEdge, edge)
		graph.Edges = append(graph.Edges, edge)
	}
	err = nil
	return
}

// ToTriples extracts triples from a graph
func (graph *Graph) ToTriples() (ret []Triple) {
	for i := range graph.Edges {
		trip := Triple{
			Sub:  graph.Edges[i].Subject.Term,
			Pred: graph.Edges[i].Pred,
			Obj:  graph.Edges[i].Object.Term,
		}
		ret = append(ret, trip)
	}
	return
}
