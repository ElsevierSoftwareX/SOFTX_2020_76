package template

// HelperHeader template
var HelperHeader = "package helper\n\n" +
	"import (\n" +
	"\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/rdf\"\n" +
	"\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/owl\"\n" +
	"\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/graph\"\n" +
	"\t\"fmt\"\n" +
	"\t\"strconv\"\n" +
	"\t\"strings\"\n" +
	"\t\"time\"\n" +
	")\n\n"

// HelperAddToGraph template
var HelperAddToGraph = "// AddObjectToGraph adds the specified object to the graph\n" +
	"func AddObjectToGraph(g *graph.Graph, typeIRI string, res owl.Thing) (node *graph.Node) {\n" +
	"\t var ok bool\n" +
	"\tif node, ok = g.Nodes[res.IRI()]; !ok {\n" +
	"\t\tnode = &graph.Node{\n" +
	"\t\t\tName: res.IRI(),\n" +
	"\t\t\tType: rdf.TermIRI,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[res.IRI()] = node\n" +
	"\t}\n" +
	"\t var typ *graph.Node\n" +
	"\tif typ, ok = g.Nodes[typeIRI]; !ok {\n" +
	"\t\ttyp = &graph.Node{\n" +
	"\t\t\tName: typeIRI,\n" +
	"\t\t\tType: rdf.TermIRI,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[typeIRI] = typ\n" +
	"\t}\n" +
	"\tpred := &graph.Edge{\n" +
	"\t\tName: \"http://www.w3.org/1999/02/22-rdf-syntax-ns#type\",\n" +
	"\t\tSubject: node,\n" +
	"\t\tObject: typ,\n" +
	"\t}\n" +
	"\tg.Edges = append(g.Edges, pred)\n" +
	"\tnode.Predicates = append(node.Predicates, pred)\n" +
	"\ttyp.InversePredicates = append(typ.InversePredicates, pred)\n" +
	"\treturn\n" +
	"}\n\n"

// HelperAddClassPropertyToGraph template
var HelperAddClassPropertyToGraph = "// AddClassPropertyToGraph adds the specified property to the graph\n" +
	"func AddClassPropertyToGraph(g *graph.Graph, propIRI string, subjNode *graph.Node, obj owl.Thing) {\n" +
	"\tif obj == nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tobjNode, ok := g.Nodes[obj.IRI()]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &graph.Node{\n" +
	"\t\t\tName: obj.IRI(),\n" +
	"\t\t\tType: rdf.TermIRI,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[obj.IRI()] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddStringPropertyToGraph template
var HelperAddStringPropertyToGraph = "// AddStringPropertyToGraph adds the specified property to the graph\n" +
	"func AddStringPropertyToGraph(g *graph.Graph, propIRI string, subjNode *graph.Node, obj string) {\n" +
	"\tif obj == \"\" {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tlit, _ := rdf.NewLiteral(obj)\n" +
	"\tobjNode, ok := g.Nodes[obj]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &graph.Node{\n" +
	"\t\t\tName: obj,\n" +
	"\t\t\tType: rdf.TermLiteral,\n" +
	"\t\t\tLiteral: &lit,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[obj] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddIntPropertyToGraph template
var HelperAddIntPropertyToGraph = "// AddIntPropertyToGraph adds the specified property to the graph\n" +
	"func AddIntPropertyToGraph(g *graph.Graph, propIRI string, subjNode *graph.Node, obj int) {\n" +
	"\tlit, _ := rdf.NewLiteral(obj)\n" +
	"\tobjNode, ok := g.Nodes[fmt.Sprintf(\"%d\", obj)]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &graph.Node{\n" +
	"\t\t\tName:fmt.Sprintf(\"%d\", obj),\n" +
	"\t\t\tType: rdf.TermLiteral,\n" +
	"\t\t\tLiteral: &lit,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[fmt.Sprintf(\"%d\", obj)] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddFloatPropertyToGraph template
var HelperAddFloatPropertyToGraph = "// AddFloatPropertyToGraph adds the specified property to the graph\n" +
	"func AddFloatPropertyToGraph(g *graph.Graph, propIRI string, subjNode *graph.Node, obj float64) {\n" +
	"\tlit, _ := rdf.NewLiteral(obj)\n" +
	"\tobjNode, ok := g.Nodes[fmt.Sprintf(\"%f\", obj)]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &graph.Node{\n" +
	"\t\t\tName:fmt.Sprintf(\"%f\", obj),\n" +
	"\t\t\tType: rdf.TermLiteral,\n" +
	"\t\t\tLiteral: &lit,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[fmt.Sprintf(\"%f\", obj)] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"\treturn\n" +
	"}\n\n"

// HelperAddBoolPropertyToGraph template
var HelperAddBoolPropertyToGraph = "// AddBoolPropertyToGraph adds the specified property to the graph\n" +
	"func AddBoolPropertyToGraph(g *graph.Graph, propIRI string, subjNode *graph.Node, obj bool) {\n" +
	"\tlit, _ := rdf.NewLiteral(obj)\n" +
	"\tobjNode, ok := g.Nodes[fmt.Sprintf(\"%v\", obj)]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &graph.Node{\n" +
	"\t\t\tName:fmt.Sprintf(\"%v\", obj),\n" +
	"\t\t\tType: rdf.TermLiteral,\n" +
	"\t\t\tLiteral: &lit,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[fmt.Sprintf(\"%v\", obj)] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"}\n\n"

// HelperAddTimePropertyToGraph template
var HelperAddTimePropertyToGraph = "// Add###timeType###PropertyToGraph adds the specified property to the graph\n" +
	"func Add###timeType###PropertyToGraph(g *graph.Graph, propIRI string, subjNode *graph.Node, obj time.Time) {\n" +
	"\tvar temp time.Time\n" +
	"\tif obj == temp {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"###timeLiteral###" +
	"\tobjNode, ok := g.Nodes[lit.String()]\n" +
	"\tif !ok {\n" +
	"\t\tobjNode = &graph.Node{\n" +
	"\t\t\tName:lit.String(),\n" +
	"\t\t\tType: rdf.TermLiteral,\n" +
	"\t\t\tLiteral: &lit,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[lit.String()] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"}\n\n"

// LiteralDateTime template
var LiteralDateTime = "\tlit := rdf.NewxsdDateTimeLiteral(obj)\n"

// LiteralDate template
var LiteralDate = "\tlit := rdf.NewxsdDateLiteral(obj)\n"

// LiteralDateTimeStamp template
var LiteralDateTimeStamp = "\tlit := rdf.NewxsdDateTimeStampLiteral(obj)\n"

// LiteralGYear template
var LiteralGYear = "\tlit := rdf.NewxsdYearLiteral(obj)\n"

// LiteralGDay template
var LiteralGDay = "\tlit := rdf.NewxsdDayLiteral(obj)\n"

// LiteralGYearMonth template
var LiteralGYearMonth = "\tlit := rdf.NewxsdYearMonthLiteral(obj)\n"

// LiteralGMonth template
var LiteralGMonth = "\tlit := rdf.NewxsdMonthLiteral(obj)\n"

// HelperAddDurationPropertyToGraph template
var HelperAddDurationPropertyToGraph = "// AddDurationPropertyToGraph adds the specified property to the graph\n" +
	"func AddDurationPropertyToGraph(g *graph.Graph, propIRI string, subjNode *graph.Node, obj time.Duration) {\n" +
	"\tvar temp time.Duration\n" +
	"\tif obj == temp {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tlit := rdf.NewxsdDurationLiteral(obj)\n" +
	"\tvar objNode *graph.Node\n" +
	"\tif objNode, ok := g.Nodes[lit.String()]; !ok {\n" +
	"\t\tobjNode = &graph.Node{\n" +
	"\t\t\tName:lit.String(),\n" +
	"\t\t\tType: rdf.TermLiteral,\n" +
	"\t\t\tLiteral: &lit,\n" +
	"\t\t}\n" +
	"\t\tg.Nodes[lit.String()] = objNode\n" +
	"\t}\n" +
	HelperAddPropertyToGraph +
	"}\n\n"

// HelperAddPropertyToGraph template
var HelperAddPropertyToGraph = "\tpred := &graph.Edge{\n" +
	"\t\tName: propIRI,\n" +
	"\t\tObject: objNode,\n" +
	"\t\tSubject: subjNode,\n" +
	"\t}\n" +
	"\tsubjNode.Predicates = append(subjNode.Predicates, pred)\n" +
	"\tobjNode.InversePredicates = append(objNode.InversePredicates, pred)\n" +
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
