/*
Copyright 2020 Institute for Automation of Complex Power Systems,
E.ON Energy Research Center, RWTH Aachen University

This project is licensed under either of
- Apache License, Version 2.0
- MIT License
at your option.

Apache License, Version 2.0:

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

MIT License:

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package template

// HelperHeader template
var HelperHeader = "package helper\n\n" +
	"import (\n" +
	"\t\"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/rdf\"\n" +
	"\t\"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/owl\"\n" +
	// "\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/graph\"\n" +
	"\tvalidator \"gopkg.in/asaskevich/govalidator.v9\"\n" +
	"\t\"fmt\"\n" +
	"\t\"strconv\"\n" +
	"\t\"strings\"\n" +
	"\t\"time\"\n" +
	")\n\n"

// HelperAddToGraph template
var HelperAddToGraph = "// AddObjectToGraph adds the specified object to the graph\n" +
	"func AddObjectToGraph(g *rdf.Graph, typeIRI string, res owl.Thing) (node *rdf.Node) {\n" +
	"\t var ok bool\n" +
	"\tif node, ok = g.Nodes[res.IRI()]; !ok {\n" +
	"\t\tif isIRI(res.IRI()) {\n" +
	"\t\t\tnode = &rdf.Node{Term: rdf.NewIRI(res.IRI())}\n" +
	"\t\t} else {\n" +
	"\t\t\tnode = &rdf.Node{Term: rdf.NewBlankNode(res.IRI())}\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[res.IRI()] = node\n" +
	"\t}\n" +
	"\tvar typ *rdf.Node\n" +
	"\tif typ, ok = g.Nodes[typeIRI]; !ok {\n" +
	"\t\ttyp = &rdf.Node{Term: rdf.NewIRI(typeIRI)}\n" +
	"\t\tg.Nodes[typeIRI] = typ\n" +
	"\t}\n" +
	"\tpred := &rdf.Edge{\n" +
	"\t\tPred: rdf.NewIRI(\"http://www.w3.org/1999/02/22-rdf-syntax-ns#type\"),\n" +
	"\t\tSubject: node,\n" +
	"\t\tObject: typ,\n" +
	"\t}\n" +
	"\tg.Edges = append(g.Edges, pred)\n" +
	"\tnode.Edge = append(node.Edge, pred)\n" +
	"\ttyp.InverseEdge = append(typ.InverseEdge, pred)\n" +
	"\treturn\n" +
	"}\n\n"

// HelperAddClassPropertyToGraph template
var HelperAddClassPropertyToGraph = "// AddClassPropertyToGraph adds the specified property to the graph\n" +
	"func AddClassPropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj owl.Thing) {\n" +
	"\tif obj == nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tobjNode, ok := g.Nodes[obj.IRI()]\n" +
	"\tif !ok {\n" +
	"\t\tif isIRI(obj.IRI()) {\n" +
	"\t\t\tobjNode = &rdf.Node{Term: rdf.NewIRI(obj.IRI())}\n" +
	"\t\t} else {\n" +
	"\t\t\tobjNode = &rdf.Node{Term: rdf.NewBlankNode(obj.IRI())}\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[obj.IRI()] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddStringPropertyToGraph template
var HelperAddStringPropertyToGraph = "// AddStringPropertyToGraph adds the specified property to the graph\n" +
	"func AddStringPropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj string) {\n" +
	"\tif obj == \"\" {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tlit, _ := rdf.NewLiteral(obj, \"\")\n" +
	"\tobjNode, ok := g.Nodes[obj]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &rdf.Node{Term: lit}\n" +
	"\t\tg.Nodes[obj] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddIntPropertyToGraph template
var HelperAddIntPropertyToGraph = "// AddIntPropertyToGraph adds the specified property to the graph\n" +
	"func AddIntPropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj int) {\n" +
	"\tlit, _ := rdf.NewLiteral(obj, \"\")\n" +
	"\tobjNode, ok := g.Nodes[fmt.Sprintf(\"%d\", obj)]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &rdf.Node{Term: lit}\n" +
	"\t\tg.Nodes[fmt.Sprintf(\"%d\", obj)] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddFloatPropertyToGraph template
var HelperAddFloatPropertyToGraph = "// AddFloatPropertyToGraph adds the specified property to the graph\n" +
	"func AddFloatPropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj float64) {\n" +
	"\tlit, _ := rdf.NewLiteral(obj, \"\")\n" +
	"\tobjNode, ok := g.Nodes[fmt.Sprintf(\"%f\", obj)]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &rdf.Node{Term: lit}\n" +
	"\t\tg.Nodes[fmt.Sprintf(\"%f\", obj)] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddBoolPropertyToGraph template
var HelperAddBoolPropertyToGraph = "// AddBoolPropertyToGraph adds the specified property to the graph\n" +
	"func AddBoolPropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj bool) {\n" +
	"\tlit, _ := rdf.NewLiteral(obj, \"\")\n" +
	"\tobjNode, ok := g.Nodes[fmt.Sprintf(\"%v\", obj)]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &rdf.Node{Term: lit}\n" +
	"\t\tg.Nodes[fmt.Sprintf(\"%v\", obj)] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"}\n\n"

// HelperAddInterfacePropertyToGraph template
var HelperAddInterfacePropertyToGraph = "// AddInterfacePropertyToGraph adds the specified property to the graph\n" +
	"func AddInterfacePropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj interface{}) {\n" +
	"\tlit, _ := rdf.NewLiteral(obj, \"\")\n" +
	"\tobjNode, ok := g.Nodes[fmt.Sprintf(\"%v\", obj)]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &rdf.Node{Term: lit}\n" +
	"\t\tg.Nodes[fmt.Sprintf(\"%v\", obj)] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"}\n\n"

// HelperAddTimePropertyToGraph template
var HelperAddTimePropertyToGraph = "// Add###timeType###PropertyToGraph adds the specified property to the graph\n" +
	"func Add###timeType###PropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj time.Time) {\n" +
	"\tvar temp time.Time\n" +
	"\tif obj == temp {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"###timeLiteral###" +
	"\tobjNode, ok := g.Nodes[lit.String()]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &rdf.Node{Term: lit}\n" +
	"\t\tg.Nodes[lit.String()] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"}\n\n"

// LiteralDateTime template
var LiteralDateTime = "\tlit, _ := rdf.NewLiteral(obj, rdf.XsdDateTime)\n"

// LiteralDate template
var LiteralDate = "\tlit, _ := rdf.NewLiteral(obj, rdf.XsdDate)\n"

// LiteralDateTimeStamp template
var LiteralDateTimeStamp = "\tlit, _ := rdf.NewLiteral(obj, rdf.XsdDateTimeStamp)\n"

// LiteralGYear template
var LiteralGYear = "\tlit, _ := rdf.NewLiteral(obj, rdf.XsdYear)\n"

// LiteralGDay template
var LiteralGDay = "\tlit, _ := rdf.NewLiteral(obj, rdf.XsdDay)\n"

// LiteralGYearMonth template
var LiteralGYearMonth = "\tlit, _ := rdf.NewLiteral(obj, rdf.XsdYearMonth)\n"

// LiteralGMonth template
var LiteralGMonth = "\tlit, _ := rdf.NewLiteral(obj, rdf.XsdMonth)\n"

// HelperAddDurationPropertyToGraph template
var HelperAddDurationPropertyToGraph = "// AddDurationPropertyToGraph adds the specified property to the graph\n" +
	"func AddDurationPropertyToGraph(g *rdf.Graph, propIRI string, subjNode *rdf.Node, obj time.Duration) {\n" +
	"\tvar temp time.Duration\n" +
	"\tif obj == temp {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tlit, _ := rdf.NewLiteral(obj, \"\")\n" +
	"\tvar objNode *rdf.Node\n" +
	"\tif objNode, ok := g.Nodes[lit.String()]; !ok {\n" +
	"\t\tobjNode = &rdf.Node{Term: lit}\n" +
	"\t\tg.Nodes[lit.String()] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"}\n\n"

// HelperAddPropertyToGraph template
var HelperAddPropertyToGraph = "\tpred := &rdf.Edge{\n" +
	"\t\tPred: rdf.NewIRI(propIRI),\n" +
	"\t\tObject: objNode,\n" +
	"\t\tSubject: subjNode,\n" +
	"\t}\n" +
	"\tsubjNode.Edge = append(subjNode.Edge, pred)\n" +
	"\tobjNode.InverseEdge = append(objNode.InverseEdge, pred)\n" +
	"\tg.Edges = append(g.Edges, pred)\n"

// HelperParseXsdDuration template
var HelperParseXsdDuration = "// ParseXsdDuration parses xsdDuration\n" +
	"func ParseXsdDuration(in string) (out time.Duration, err error) {\n" +
	"\tp := strings.TrimPrefix(in, \"P\")\n" +
	"\tstr := \"\"\n" +
	"\th := strings.Split(p, \"H\")\n" +
	"\tif len(h) > 1 {\n" +
	"\t\tif n, err := strconv.Atoi(h[0]); err == nil && n > 0 {\n" +
	"\t\t\tstr += h[0] + \"h\"\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\tm := strings.Split(h[len(h)-1], \"M\")\n" +
	"\tif len(m) > 1 {\n" +
	"\t\tif n, err := strconv.Atoi(m[0]); err == nil && n > 0 {\n" +
	"\t\t\tstr += m[0] + \"m\"\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\ts := strings.Split(m[len(m)-1], \"S\")\n" +
	"\tif len(s) > 1 {\n" +
	"\t\tif n, err := strconv.ParseFloat(s[0], 32); err == nil && n > 0 {\n" +
	"\t\t\tstr += s[0] + \"s\"\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\tout, err = time.ParseDuration(str)\n" +
	"\treturn\n" +
	"}\n\n"

// HelperIsIRI template
var HelperIsIRI = "// isIRI checks if string is valid iri\n" +
	"func isIRI(iri string) (ok bool) {\n" +
	"\tok = false\n" +
	"\tif validator.IsURL(iri) {\n" +
	"\t\tok = true\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"
