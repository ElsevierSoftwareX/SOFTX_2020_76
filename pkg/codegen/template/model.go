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

// Package template defines the templates used for code generation in string format.
package template

// ModelHeader template
var ModelHeader = "package ###pkgName###\n\n" +
	"import (\n" +
	// "\t\"git-ce.rwth-aachen.de/acs/private/research/ensure/owl/owl.git/pkg/graph\"\n" +
	"\t\"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/rdf\"\n" +
	"\t\"git.rwth-aachen.de/acs/public/ontology/owl/owl2go/pkg/owl\"\n" +
	"\t\"io\"\n" +
	")\n\n"

// ModelStruct template
var ModelStruct = "// Model holds all objects\n" +
	"type Model struct {\n" +
	"\tmThing map[string]owl.Thing\n" +
	"###objectMaps###" +
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
	"\ttriples, err := rdf.DecodeTTL(input)\n" +
	"\tif err != nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tg, err := rdf.NewGraph(triples)\n" +
	"\tif err != nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tmod, err = NewModelFromGraph(g)\n" +
	"\treturn\n" +
	"}\n\n"

// ModelNewFromJSONLD template
var ModelNewFromJSONLD = "// NewModelFromJSONLD creates a new model from a jsonld io reader\n" +
	"func NewModelFromJSONLD(input io.Reader) (mod *Model, err error) {\n" +
	"\ttriples, err := rdf.DecodeJSONLD(input)\n" +
	"\tif err != nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tg, err := rdf.NewGraph(triples)\n" +
	"\tif err != nil {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tmod, err = NewModelFromGraph(g)\n" +
	"\treturn\n" +
	"}\n\n"

// ModelNewFromGraph template
var ModelNewFromGraph = "// NewModelFromGraph creates a new model from a owl graph\n" +
	"func NewModelFromGraph(g rdf.Graph) (mod *Model, err error) {\n" +
	"\tmod = NewModel()\n" +
	"\tfor i := range g.Nodes {\n" +
	"\t\tfor j := range g.Nodes[i].Edge {\n" +
	"\t\t\tif g.Nodes[i].Edge[j].Pred.String() == \"http://www.w3.org/1999/02/22-rdf-syntax-ns#type\" {\n" +
	"\t\t\t\tswitch g.Nodes[i].Edge[j].Object.Term.String() {\n" +
	"###newObjects###" +
	"\t\t\t\tdefault:\n" +
	"\t\t\t\t}\n" +
	"\t\t\t}\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\tfor i := range g.Nodes {\n" +
	"\t\tif res, ok := mod.mThing[g.Nodes[i].Term.String()]; ok {\n" +
	"\t\t\tres.InitFromNode(g.Nodes[i])\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// NewObject template
var NewObject = "\t\t\t\tcase \"###classIRI###\":\n" +
	"\t\t\t\t\tmod.New###capImportName######className###(g.Nodes[i].Term.String())\n"

// ModelToGraph template
var ModelToGraph = "// ToGraph extracts an owl graph from an existing model\n" +
	"func (mod* Model) ToGraph() (g *rdf.Graph) {\n" +
	"\tg = &rdf.Graph{}\n" +
	"\tg.Nodes = make(map[string]*rdf.Node)\n" +
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
	"\terr = rdf.EncodeTTL(newTriples, output)\n" +
	// "\toutput.Close()\n" +
	"\treturn\n" +
	"}\n\n"

// ModelToJSONLD template
var ModelToJSONLD = "// ToJSONLD writes a jsonld file from an existing model\n" +
	"func (mod* Model) ToJSONLD(output io.Writer) (err error) {\n" +
	"\tg := mod.ToGraph()\n" +
	"\tnewTriples := g.ToTriples()\n" +
	"\terr = rdf.EncodeJSONLD(newTriples, output)\n" +
	// "\toutput.Close()\n" +
	"\treturn\n" +
	"}\n\n"

// ModelDeleteObject template
var ModelDeleteObject = "// DeleteObject deletes an object from the model along with its references\n" +
	"func (mod *Model) DeleteObject(obj owl.Thing) (err error) {\n" +
	"\tif !mod.Exist(obj.IRI()) {\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tg := mod.ToGraph()\n" +
	"\tn := g.Nodes[obj.IRI()]\n" +
	"\tfor i := range n.InverseEdge {\n" +
	"\t\tif subj, ok := mod.mThing[n.InverseEdge[i].Subject.Term.String()]; ok {\n" +
	"\t\t\tsubj.RemoveObject(obj, n.InverseEdge[i].Pred.String())\n" +
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
