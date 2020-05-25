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

// ClassHeader template
var ClassHeader = "package ###pkgName###\n\n" +
	"import (\n" +
	"###imports###" +
	")\n\n"

// Import template
var Import = "\t\"###import###\"\n"

// ClassInterface template
var ClassInterface = "// ###className### ###comment###\n" +
	"type ###className### interface {\n" +
	"###interfaceMethods###" +
	"###interfaceInheritance###" +
	"}\n\n"

// InterfaceInterface template
var InterfaceInterface = "\t###propName######multi######propBaseTypeNoImp###\n"

// InterfaceInheritance template
var InterfaceInheritance = "\tIs###parentName###() bool // indicates base class\n"

// ClassStruct template
var ClassStruct = "// s###className### ###comment###\n" +
	"type s###className### struct {\n" +
	"###structProperties###" +
	"}\n\n"

// StructProperty template
var StructProperty = "\t###propLongName###\n"

//ClassNew template
var ClassNew = "// New###className### creates a new ###className###\n" +
	"func (mod *Model) New###className###(iri string) (ret ###className###, err error) {\n" +
	"\tif mod.Exist(iri) {\n" +
	"\t\terr = errors.New(\"Resource already exists\")\n" +
	"\t\treturn\n" +
	"\t}\n" +
	"\tres := &s###className###{}\n" +
	"\tpc := propCommon{iri: iri, typ: \"###className###\", model: mod}\n" +
	"\tres.propCommon = pc\n" +
	"\tmod.add###className###(res)\n" +
	"\tres.makeMaps()\n" +
	"\tret = res\n" +
	"\treturn\n" +
	"}\n\n"

// ClassMakeMaps template
var ClassMakeMaps = "// makeMaps creates initiliazes maps\n" +
	"func (res *s###className###) makeMaps() {\n" +
	"###newMakeMaps###" +
	"###newInitProps###" +
	"\treturn\n" +
	"}\n\n"

// NewMakeMap template
var NewMakeMap = "\tres.###propName### = make(map[string]###propType###)\n"

// NewInitPropClassMultiple template
var NewInitPropClassMultiple = "\tif v, ok := res.model.m###propType###[\"###value###\"]; ok {\n" +
	"\t\tres.Add###propCapital###(v)\n" +
	"\t}\n"

// NewInitPropClassSingle template
var NewInitPropClassSingle = "\tif v, ok := res.model.m###propType###[\"###value###\"]; ok {\n" +
	"\t\tres.Set###propCapital###(v)\n" +
	"\t}\n"

// NewInitPropLiteralMultiple template
var NewInitPropLiteralMultiple = "\tAdd###propCapital###(###value###)\n"

// NewInitPropLiteralSingle template
var NewInitPropLiteralSingle = "\tSet###propCapital###(###value###)\n"

// ClassAdd template
var ClassAdd = "// add###className### adds ###className### to model\n" +
	"func (mod *Model) add###className###(res ###className###) {\n" +
	"###newAddToMaps###" +
	"\treturn\n" +
	"}\n\n"

// AddToMap template
var AddToMap = "\tmod.m###parentName###[res.IRI()] = res\n"

// ClassGet template
var ClassGet = "// ###className### returns all resources with given prefix\n" +
	"func (mod *Model) ###className###(prefix string) (res []###className###) {\n" +
	"\tfor i := range mod.m###className### {\n" +
	"\t\tif strings.HasPrefix(i, prefix) {\n" +
	"\t\t\tres = append(res, mod.m###className###[i])\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// ClassRemove template
var ClassRemove = "// RemoveObject deletes all its references in this object\n" +
	"func (res *s###className###) RemoveObject(obj owl.Thing, prop string) {\n" +
	"\tswitch prop {\n" +
	"###removeProps###" +
	"\tdefault:\n" +
	"###parentRemove###" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// RemovePropMultiple template
var RemovePropMultiple = "\tcase \"###propIRI###\":\n" +
	"\t\tif v, ok := obj.(###propBaseType###); ok {\n" +
	"\t\t\tres.Del###propCapital###(v)\n" +
	"\t\t}\n"

// RemovePropSingle template
var RemovePropSingle = "\tcase \"###propIRI###\":\n" +
	"\t\tif v, ok := obj.(###propBaseType###); ok {\n" +
	"\t\t\tif res.###propName###.IRI() == v.IRI() {\n" +
	"\t\t\t\tres.Set###propCapital###(nil)\n" +
	"\t\t\t}\n" +
	"\t\t}\n"

// RemoveNonInverse template
var RemoveNonInverse = "\tcase \"###propIRI###\":\n" +
	"\t\tres.###propLongName.removeObject(obj)\n"

// ClassIRI template
var ClassIRI = "// IRI is the resource iri\n" +
	"func (res *s###className###) IRI() (out string) {\n" +
	"\tout = res.iri\n" +
	"\treturn\n" +
	"}\n\n"

// ClassInheritance template
var ClassInheritance = "// is###parentName### indicates base class\n" +
	"func (res *s###className###) Is###parentName###() bool {\n" +
	"\treturn true\n" +
	"}\n\n"

// ClassInverseMultipleMultiple template
var ClassInverseMultipleMultiple = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *s###className###) Set###propCapital###(in []###propBaseType###) (err error) {\n" +
	"\ttemp := res.###propName###\n" +
	"\tres.###propName### = make(map[string]###propType###)\n" +
	"\tfor i := range temp {\n" +
	"\t\ttemp[i].Del###propInverse###(res)\n" +
	"\t}\n" +
	"\terr = res.Add###propCapital###(in...)\n" +
	"\treturn\n" +
	"}\n\n" +
	"// Add###propCapital### adds ###comment###\n" +
	"func (res *s###className###) Add###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	"###inverseAddMultipleMultipleAllowed###" +
	"\t\terr = errors.New(\"Wrong ###propType### type. Allowed types are ###allowedTypes###\")\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n" +
	"// Del###propCapital### deletes ###comment###\n" +
	"func (res *s###className###) Del###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	"###inverseDelMultipleMultipleAllowed###" +
	"\t\terr = errors.New(\"Wrong ###propType### type. Allowed types are ###allowedTypes###\")\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// InverseAddMultipleMultiple template
var InverseAddMultipleMultiple = "\t\tif v, ok := in[i].(###propAllowedType###); ok {\n" +
	"\t\t\tif _, ok := res.###propName###[v.IRI()]; !ok {\n" +
	"\t\t\t\tres.###propName###[v.IRI()] = v\n" +
	"\t\t\t\tv.Add###propInverse###(res)\n" +
	"\t\t\t}\n" +
	"\t\t\tcontinue\n" +
	"\t\t}\n"

// InverseDelMultipleMultiple template
var InverseDelMultipleMultiple = "\t\tif v, ok := in[i].(###propAllowedType###); ok {\n" +
	"\t\t\tif _, ok := res.###propName###[v.IRI()]; !ok {\n" +
	"\t\t\t\tdelete(res.###propName###, v.IRI())\n" +
	"\t\t\t\tv.Del###propInverse###(res)\n" +
	"\t\t\t}\n" +
	"\t\t\tcontinue\n" +
	"\t\t}\n"

// ClassInverseMultipleSingle template
var ClassInverseMultipleSingle = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *s###className###) Set###propCapital###(in []###propBaseType###) (err error) {\n" +
	"\ttemp := res.###propName###\n" +
	"\tres.###propName### = make(map[string]###propType###)\n" +
	"\tfor i := range temp {\n" +
	"\t\ttemp[i].Del###propInverse###(res)\n" +
	"\t}\n" +
	"\terr = res.Add###propCapital###(in...)\n" +
	"\treturn\n" +
	"}\n\n" +
	"// Add###propCapital### adds ###comment###\n" +
	"func (res *s###className###) Add###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	"\t\tif _, ok := res.###propName###[in[i].IRI()]; !ok {\n" +
	"\t\t\tres.###propName###[in[i].IRI()] = in[i]\n" +
	"\t\t\tin[i].Add###propInverse###(res)\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n" +
	"// Del###propCapital### deletes ###comment###\n" +
	"func (res *s###className###) Del###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	"\t\tif _, ok := res.###propName###[in[i].IRI()]; ok {\n" +
	"\t\t\tdelete(res.###propName###, in[i].IRI())\n" +
	"\t\t\tin[i].Del###propInverse###(res)\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// ClassInverseSingleMultiple template
var ClassInverseSingleMultiple = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *s###className###) Set###propCapital###(in ###propBaseType###) (err error) {\n" +
	"\tif in != nil {\n" +
	"\t\tif res.###propName### != nil {\n" +
	"\t\t\tif res.###propName###.IRI() != in.IRI() {\n" +
	"\t\t\t\ttemp := res.###propName###\n" +
	"###inverseSetSingleMultipleOne###" +
	"\t\t\t}\n" +
	"\t\t} else {\n" +
	"###inverseSetSingleMultipleTwo###" +
	"\t\t}\n" +
	"\t} else {\n" +
	"\t\tif res.###propName### != nil {\n" +
	"\t\t\ttemp := res.###propName###\n" +
	"###inverseSetSingleMultipleThree###" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// InverseSetSingleMultipleOne template
var InverseSetSingleMultipleOne = "\t\t\t\tif _, ok := in.(###propAllowedType###); ok {\n" +
	"\t\t\t\t\tres.###propName### = in\n" +
	"\t\t\t\t\ttemp.Set###propInverse###(nil)\n" +
	"\t\t\t\t\tin.Set###propInverse###(res)\n" +
	"\t\t\t\t\treturn\n" +
	"\t\t\t\t}\n"

// InverseSetSingleMultipleTwo template
var InverseSetSingleMultipleTwo = "\t\t\tif _, ok := in.(###propAllowedType###); ok {\n" +
	"\t\t\t\tres.###propName### = in\n" +
	"\t\t\t\tin.Set###propInverse###(res)\n" +
	"\t\t\t\treturn\n" +
	"\t\t\t}\n"

// InverseSetSingleMultipleThree template
var InverseSetSingleMultipleThree = "\t\t\tif _, ok := in.(###propAllowedType###); ok {\n" +
	"\t\t\t\tres.###propName### = in\n" +
	"\t\t\t\ttemp.Set###propInverse###(nil)\n" +
	"\t\t\t\treturn\n" +
	"\t\t\t}\n"

// ClassInverseSingleSingle template
var ClassInverseSingleSingle = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *s###className###) Set###propCapital###(in ###propBaseType###) (err error) {\n" +
	"\tif in != nil {\n" +
	"\t\tif res.###propName### != nil {\n" +
	"\t\t\tif res.###propName###.IRI() != in.IRI() {\n" +
	"\t\t\t\ttemp := res.###propName###\n" +
	"\t\t\t\tres.###propName### = in\n" +
	"\t\t\t\ttemp.Set###propInverse###(nil)\n" +
	"\t\t\t\tin.Set###propInverse###(res)\n" +
	"\t\t\t}\n" +
	"\t\t} else {\n" +
	"\t\t\tres.###propName### = in\n" +
	"\t\t\tin.Set###propInverse###(res)\n" +
	"\t\t}\n" +
	"\t} else {\n" +
	"\t\tif res.###propName### != nil {\n" +
	"\t\t\ttemp := res.###propName###\n" +
	"\t\t\tres.###propName### = in\n" +
	"\t\t\ttemp.Set###propInverse###(nil)\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// ClassInit template
var ClassInit = "// InitFromNode initializes the resource from a graph node\n" +
	"func (res *s###className###) InitFromNode(node *rdf.Node) (err error) {\n" +
	"\tfor i := range node.Edge {\n" +
	"\t\tres.propsInit(node.Edge[i])\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// PropsInit template
var PropsInit = "// propsInit initializes the property from a graph node\n" +
	"func (res *s###className###) propsInit(pred *rdf.Edge) (err error) {\n" +
	"\tswitch pred.Pred.String() {\n" +
	"###initSwitchProps###" +
	"\tdefault:\n" +
	"###parentPropInit###" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// InitSwitchProp template
var InitSwitchProp = "\tcase \"###propIRI###\":\n" +
	"###PropInit###"

// PropInitClassNonInverse template
var PropInitClassNonInverse = "\t\tres.###propLongName###.init(res.model, pred.Object.Term.String())\n"

// PropInitLiteralNonInverse template
var PropInitLiteralNonInverse = "\t\tres.###propLongName###.init(pred.Object.Term.String())\n"

// PropClassBaseThing template
var PropClassBaseThing = "\t\tif obj, ok := res.model.mThing[pred.Object.Term.String()]; ok {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropClassDefault template
var PropClassDefault = "\t\tif obj, ok := res.model.m###propBaseType###[pred.Object.Term.String()]; ok {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropClassImport template
var PropClassImport = "\t\tif temp := res.model.###capImportName######propBaseType###(pred.Object.Term.String()); len(temp) > 0 {\n" +
	"\t\t\tfor j := range temp {\n" +
	"\t\t\t\tif temp[j].IRI() == pred.Object.Term.String() {\n" +
	"\t\t\t\t\tres.###Multiplicity######propCapital###(temp[j])\n" +
	"\t\t\t\t}\n" +
	"\t\t\t}\n" +
	"\t\t}\n"

// PropTime template
var PropTime = "\t\tif obj, err := time.Parse(\"15:04:05Z07:00\", pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropDate template
var PropDate = "\t\tif obj, err := time.Parse(\"2006-01-02Z07:00\", pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropDateTime template
var PropDateTime = "\t\tif obj, err := time.Parse(time.RFC3339, pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropDateTimeStamp template
var PropDateTimeStamp = "\t\tif obj, err := time.Parse(time.RFC3339, pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropGDay template
var PropGDay = "\t\tif obj, err := time.Parse(\"---02\", pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropGMonth template
var PropGMonth = "\t\tif obj, err := time.Parse(\"--01\", pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropGYear template
var PropGYear = "\t\tif obj, err := time.Parse(\"2006\", pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropGYearMonth template
var PropGYearMonth = "\t\tif obj, err := time.Parse(\"2006-01\", pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropDuration template
var PropDuration = "\t\tif obj, err := owl.ParseXsdDuration(pred.Object.Term.String()); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropFloat template
var PropFloat = "\t\tif obj, err := strconv.ParseFloat(pred.Object.Term.String(), 32); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(float64(obj))\n" +
	"\t\t}\n"

// PropInt template
var PropInt = "\t\tif obj, err := strconv.Atoi(); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropBool template
var PropBool = "\t\tif obj, err := strconv.ParseBool(in); err == nil {\n" +
	"\t\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t\t}\n"

// PropString template
var PropString = "\t\tres.###Multiplicity######propCapital###()\n"

// ClassToGraph template
var ClassToGraph = "// ToGraph creates a new owl graph node and adds it to the graph\n" +
	"func (res *s###className###) ToGraph(g *rdf.Graph) {\n" +
	"\tnode := owl.AddObjectToGraph(g, \"###classIRI###\", res)\n" +
	"\tres.propsToGraph(node, g)\n" +
	"\treturn\n" +
	"}\n\n"

// PropsToGraph template
var PropsToGraph = "// propsToGraph adds all properties to the graph\n" +
	"func (res *s###className###) propsToGraph(node *rdf.Node, g *rdf.Graph) {\n" +
	"###toGraphProps###" +
	"\treturn\n" +
	"}\n\n"

// ToGraphProp template
var ToGraphProp = "\tres.###propLongName###.toGraph(node, g)\n"

// ClassToGraphNoProp template
var ClassToGraphNoProp = "// ToGraph creates a new owl graph node and adds it to the graph\n" +
	"func (res *s###className###) ToGraph(g *rdf.Graph) {\n" +
	"\towl.AddObjectToGraph(g, \"###classIRI###\", res)\n" +
	"\treturn\n" +
	"}\n\n"

// ClassString template
var ClassString = "// String prints the object into a string\n" +
	"func (res *s###className###) String() (ret string) {\n" +
	"\tret = res.iri + \" \" + \"###className###\\n\"\n" +
	"\tret += res.propsString()\n" +
	"\treturn\n" +
	"}\n\n"

// PropsString template
var PropsString = "// propsString prints the object into a string\n" +
	"func (res *s###className###) propsString() (ret string) {\n" +
	"\tret = \"\"\n" +
	"###stringProps###" +
	"\treturn\n" +
	"}\n\n"

// StringProp template
var StringProp = "\tret += res.###propLongName###.String()\n"
