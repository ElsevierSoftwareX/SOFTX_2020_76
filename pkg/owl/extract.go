package owl

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/rdf"
)

// ExtractOntology extracts all classes, properties, individuals and imports
func ExtractOntology(link string) (on Ontology, err error) {
	iri := ""
	description := ""
	on.Class = make(map[string]*Class)
	on.Property = make(map[string]*Property)
	on.Individual = make(map[string]*Individual)
	on.Imports = make(map[string][]string)
	on.Description = make(map[string]string)
	on.Content = make(map[string][]byte)

	var resp *http.Response
	resp, err = requestOntology(link)
	if err != nil {
		return
	}
	var g rdf.Graph
	g, iri, description, _, err = parseOntology(resp.Body)
	if err != nil {
		return
	}
	resp.Body.Close()

	on.graph = &g
	on.Description[iri] = description
	on.Imports[iri] = []string{}

	err = on.parseImports(on.graph)
	if err != nil {
		return
	}

	on.Class, err = extractClasses(on.graph)
	if err != nil {
		return
	}

	on.Property, err = extractProperties(on.graph)
	if err != nil {
		return
	}

	on.Individual, err = extractIndividuals(on.graph, on.Class)
	if err != nil {
		return
	}

	err = on.postProcessProperties()
	if err != nil {
		return
	}
	err = on.postProcessClasses()
	if err != nil {
		return
	}
	err = on.addPropertyDomain()
	if err != nil {
		return
	}

	return
}

// parseOntology parses the specified ontology
func parseOntology(input io.Reader) (g rdf.Graph, iri string, description string, content []byte,
	err error) {
	fmt.Println("Read TTL input")
	g, err = readTTL(input)
	if err != nil {
		return
	}

	// get ontology iri
	for i := range g.Edges {
		if g.Edges[i].Pred.String() == "http://www.w3.org/1999/02/22-rdf-syntax-ns#type" &&
			g.Edges[i].Object.Term.String() == "http://www.w3.org/2002/07/owl#Ontology" {
			iri = g.Edges[i].Subject.Term.String()
		} else if g.Edges[i].Pred.String() == "http://purl.org/dc/terms/description" {
			description = g.Edges[i].Object.Term.String()
		}
	}
	return
}

// readTTL reads a ttl file and returns a graph
func readTTL(input io.Reader) (g rdf.Graph, err error) {
	var triples []rdf.Triple
	triples, err = rdf.DecodeTTL(input)
	if err != nil {
		return
	}

	g, err = rdf.NewGraph(triples)
	return
}

// parseImports parses all imports and adds imports to ontologies
func (on *Ontology) parseImports(gIn *rdf.Graph) (err error) {
	var gTemp rdf.Graph
	gTemp.Nodes = make(map[string]*rdf.Node)
	hasImport := false
	for i := range gIn.Edges {
		if gIn.Edges[i].Pred.String() == "http://www.w3.org/2002/07/owl#imports" {
			hasImport = true
			iri := gIn.Edges[i].Subject.Term.String()
			impIRI := gIn.Edges[i].Object.Term.String()

			on.Imports[iri] = append(on.Imports[iri], impIRI)

			var resp *http.Response
			resp, err = requestOntology(gIn.Edges[i].Object.Term.String())
			if err != nil {
				return
			}
			fmt.Println("Parse Imported Ontology " + gIn.Edges[i].Object.Term.String())
			var g rdf.Graph
			var desc string
			g, impIRI, desc, _, err = parseOntology(resp.Body)
			if err != nil {
				return
			}
			resp.Body.Close()

			on.Description[impIRI] = desc
			on.Imports[impIRI] = []string{}

			gTemp.Merge(&g)
		}
	}
	if hasImport {
		on.parseImports(&gTemp)
		on.graph.Merge(&gTemp)
	}
	return
}

// getComment returns a comment if it exists (rdf:comment)
func getComment(node *rdf.Node) (ret string) {
	// find comment
	for j := range node.Edge {
		if node.Edge[j].Pred.String() == "http://www.w3.org/2000/01/rdf-schema#comment" {
			regex := regexp.MustCompile(`\r?\n`)
			ret = regex.ReplaceAllString(node.Edge[j].Object.Term.String(), " ")
			regex = regexp.MustCompile(`\n`)
			ret = regex.ReplaceAllString(ret, " ")
			// ret = strings.Replace(node.Predicates[j].Object.Name, "\n", " ", -1)

			break
		}
	}
	if ret == "" {
		ret = "no comment"
	}
	return
}

// getUnionValues returns all values of a union (rdfs:first and rdfs:rest)
func getUnionValues(node *rdf.Node) (ret []*rdf.Node) {
	for i := range node.Edge {
		if node.Edge[i].Pred.String() == "http://www.w3.org/1999/02/22-rdf-syntax-ns#first" {
			ret = append(ret, node.Edge[i].Object)
		} else if node.Edge[i].Pred.String() == "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest" &&
			node.Edge[i].Object.Term.String() != "http://www.w3.org/1999/02/22-rdf-syntax-ns#nil" {
			ret = append(ret, getUnionValues(node.Edge[i].Object)...)
		}
	}
	return
}

func requestOntology(path string) (resp *http.Response, err error) {
	client := &http.Client{
		Timeout: time.Second * 2,
	}
	resp, err = client.Get(path)
	return
}
