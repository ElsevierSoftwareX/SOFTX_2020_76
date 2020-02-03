package template

// ModelHeader template
var ModelHeader = "package ###pkgName###\n\n" +
	"import (\n" +
	"\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/graph\"\n" +
	"\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/rdf\"\n" +
	"\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/owl\"\n" +
	"\t\"io\"\n" +
	"###imports###" +
	")\n\n"

// ModelStruct template
var ModelStruct = "// Model holds all objects\n" +
	"type Model struct {\n" +
	"\tmThing map[string]owl.Thing\n" +
	"###objectMaps###" +
	"###importModels###" +
	"}\n\n"

// StructMap template
var StructMap = "\tm###className### map[string]###className###\n"

// StructImport template
var StructImport = "\tmodel###importName### *im###importName###.Model\n"

// ModelNew template
var ModelNew = "// NewModel creates a new model and initializes class maps\n" +
	"func NewModel() (mod *Model) {\n" +
	"\tmod = &Model{}\n" +
	"\tmod.mThing = make(map[string]owl.Thing)\n" +
	"###makeMaps###" +
	"###newImportModels###" +
	"\treturn\n" +
	"}\n\n"

// NewObjectMap template
var NewObjectMap = "\tmod.m###className### = make(map[string]###className###)\n"

// NewImport template
var NewImport = "\tmod.model###importName### = im###importName###.NewModel()\n"

// ModelExists template
var ModelExists = "// Exist checks for the existance of a resource by the given iri\n" +
	"func (mod *Model) Exist(iri string) (ret bool) {\n" +
	"\t_, ret = mod.mThing[iri]\n" +
	"\treturn\n" +
	"}\n\n"

// ModelNewFromTTL template
var ModelNewFromTTL = "// NewModelFromTTL creates a new model from a ttl io reader\n" +
	"func NewModelFromTTL(input io.Reader) (mod *Model, err error) {\n" +
	"\tdec := rdf.NewTripleDecoder(input, rdf.Turtle)\n" +
	"\ttriples, err := dec.DecodeAll()\n" +
	"\tif err != nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tg, err := graph.New(triples)\n" +
	"\tif err != nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tmod, err = NewModelFromGraph(g)\n" +
	"\treturn\n" +
	"}\n\n"

// ModelNewFromGraph template
var ModelNewFromGraph = "// NewModelFromGraph creates a new model from a owl graph\n" +
	"func NewModelFromGraph(g graph.Graph) (mod *Model, err error) {\n" +
	"\tmod = NewModel()\n" +
	"\tfor i := range g.Nodes {\n" +
	"\t\tfor j := range g.Nodes[i].Predicates {\n" +
	"\t\t\tif g.Nodes[i].Predicates[j].Name == \"http://www.w3.org/1999/02/22-rdf-syntax-ns#type\" {\n" +
	"\t\t\t\tswitch g.Nodes[i].Predicates[j].Object.Name {\n" +
	"###newObjects###" +
	"\t\t\t\tdefault:\n" +
	"\t\t\t\t}\n" +
	"\t\t\t}\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\tfor i := range g.Nodes {\n" +
	"\t\tif res, ok := mod.mThing[g.Nodes[i].Name]; ok {\n" +
	"\t\t\tres.InitFromNode(g.Nodes[i])\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// NewObject template
var NewObject = "\t\t\t\tcase \"###classIRI###\":\n" +
	"\t\t\t\t\tmod.New###capImportName######className###(g.Nodes[i].Name)\n"

// ModelToGraph template
var ModelToGraph = "// ToGraph extracts an owl graph from an existing model\n" +
	"func (mod* Model) ToGraph() (g *graph.Graph) {\n" +
	"\tg = &graph.Graph{}\n" +
	"\tg.Nodes = make(map[string]*graph.Node)\n" +
	"\tfor i := range mod.mThing {\n" +
	"\t\tmod.mThing[i].ToGraph(g)\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// ModelToTTL template
var ModelToTTL = "// ToTTL writes a ttl file from an existing model\n" +
	"func (mod* Model) ToTTL(output io.Writer) (err error) {\n" +
	"\tg := mod.ToGraph()\n" +
	"\tnewTriples := g.ToTriples()\n" +
	"\tout := rdf.NewTripleEncoder(output, rdf.Turtle)\n" +
	"\terr = out.EncodeAll(newTriples)\n" +
	"\tout.Close()\n" +
	"\treturn\n" +
	"}\n\n"

// ModelDeleteObject template
var ModelDeleteObject = "// DeleteObject deletes an object from the model along with its references\n" +
	"func (mod *Model) DeleteObject(obj owl.Thing) (err error) {\n" +
	"\tif !mod.Exist(obj.IRI()) {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"###deleteFromImports###" +
	"\tg := mod.ToGraph()\n" +
	"\tn := g.Nodes[obj.IRI()]\n" +
	"\tfor i := range n.InversePredicates {\n" +
	"\t\tif subj, ok := mod.mThing[n.InversePredicates[i].Subject.Name]; ok {\n" +
	"\t\t\tsubj.RemoveObject(obj, n.InversePredicates[i].Name)\n" +
	"\t\t}\n" +
	"\t}\n" +
	"###deleteFromMaps###" +
	"\tdelete(mod.mThing, obj.IRI())\n" +
	"\treturn\n" +
	"}\n\n"

// DeleteFromImport template
var DeleteFromImport = "\terr = mod.model###importName###.DeleteObject(obj)\n"

// DeleteFromMap template
var DeleteFromMap = "\tdelete(mod.m###className###, obj.IRI())\n"

// ModelString template
var ModelString = "// String prints the model\n" +
	"func (mod *Model) String() (ret string) {\n" +
	"\tret = \"\"\n" +
	"\tfor i := range mod.mThing {\n" +
	"\t\tret += mod.mThing[i].String()\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// ModelToDot template
var ModelToDot = "// ToDot writes a dot file from an existing model\n" +
	"func (mod* Model) ToDot(output io.Writer) (err error) {\n" +
	"\treplace := make(map[string]string)\n" +
	"\tshape := make(map[string]string)\n" +
	"###importReplace###" +
	"###importShape###" +
	"\treplace[\"http://www.wurvoc.org/vocabularies/om-1.8/\"] = \"om:\"\n" +
	"\treplace[\"http://www.w3.org/1999/02/22-rdf-syntax-ns#\"] = \"rdf:\"\n" +
	"\treplace[\"http://www.w3.org/2000/01/rdf-schema#\"] = \"rdfs:\"\n" +
	"\tg := mod.ToGraph()\n" +
	"\tg.ToGraphvizDot(output, replace, shape)\n" +
	"\treturn\n" +
	"}\n\n"

// ImportReplace template
var ImportReplace = "\treplace[\"###importIRI###\"] = \"###importName###:\"\n"

// ImportShape template
var ImportShape = "\tshape[\"###importIRI###\"] = \"box\"\n"

// ImportsHeader template
var ImportsHeader = "package ###pkgName###\n\n" +
	"import (\n" +
	"\t\"errors\"\n" +
	"###imports###" +
	")\n\n"

// ImportsNewGetMethods template
var ImportsNewGetMethods = "// New###capImportName######className### creates a new ###importName### ###className###\n" +
	"func (mod *Model) New###capImportName######className###(iri string) (ret im###importName###.###className###, err error) {\n" +
	"\tif mod.Exist(iri) {\n" +
	"\t\terr = errors.New(\"Resource already exists\")\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tret, err = mod.model###importModelName###.New###importCapImportName######className###(iri)\n" +
	"\tif err != nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tmod.mThing[ret.IRI()] = ret\n" +
	"\treturn\n" +
	"}\n\n" +
	"// ###capImportName######className### returns all resources with given prefix\n" +
	"func (mod *Model) ###capImportName######className###(prefix string) (res []im###importName###.###className###) {\n" +
	"\tres = mod.model###importModelName###.###importCapImportName######className###(prefix)\n" +
	"\treturn\n" +
	"}\n\n"
