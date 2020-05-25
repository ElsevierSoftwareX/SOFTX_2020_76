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

// PropertyHeader template
var PropertyHeader = "package ###pkgName###\n\n" +
	"###propImports###"

// PropertyStructCommon template
var PropertyStructCommon = "type propCommon struct {\n" +
	"\tiri string // resource iri\n" +
	"\ttyp string // type of resource\n" +
	"\tmodel *Model // pointer to model\n" +
	"}\n\n"

// PropertyStructMultipleClass template
var PropertyStructMultipleClass = "// ###propLongName### is ###comment###\n" +
	"type ###propLongName### struct {\n" +
	"\t###propName### map[string]###propType###\n" +
	"}\n\n"

// PropertyStructMultipleLiteral template
var PropertyStructMultipleLiteral = "// ###propLongName### is ###comment###\n" +
	"type ###propLongName### struct {\n" +
	"\t###propName### []###propType###\n" +
	"}\n\n"

// PropertyStructSingle template
var PropertyStructSingle = "// ###propLongName### is ###comment###\n" +
	"type ###propLongName### struct {\n" +
	"\t###propName### ###propType###\n" +
	"}\n\n"

// PropertyInterfaceSingle template
var PropertyInterfaceSingle = "// ###propName###Single###propBaseTypeNoImp### is interface for property ###propCapital### with single type ###propBaseType###\n" +
	"type ###propName###Single###propBaseTypeNoImp### interface {\n" +
	"\t###propCapital###() ###propBaseType### // ###comment###\n" +
	"\tSet###propCapital###(###propBaseType###) error // set ###comment###\n" +
	"}\n\n"

// PropertyInterfaceMultiple template
var PropertyInterfaceMultiple = "// ###propName###Multiple###propBaseTypeNoImp### is interface for property ###propCapital### with multiple type ###propBaseType###\n" +
	"type ###propName###Multiple###propBaseTypeNoImp### interface {\n" +
	"\t###propCapital###() []###propBaseType### // ###comment###\n" +
	"\tSet###propCapital###([]###propBaseType###) error // set ###comment###\n" +
	"\tAdd###propCapital###(...###propBaseType###) error // add ###comment###\n" +
	"\tDel###propCapital###(...###propBaseType###) error // delete ###comment###\n" +
	"}\n\n"

// PropertyIRI template
var PropertyIRI = "// IRI is the resource iri\n" +
	"func (res *propCommon) IRI() (out string) {\n" +
	"\tout = res.iri\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyInitClass template
var PropertyInitClass = "// init initializes the property\n" +
	"func (res *###propLongName###) init(model *Model, in string) {\n" +
	"###PropInit###" +
	"\treturn\n" +
	"}\n\n"

// PropertyInitLiteral template
var PropertyInitLiteral = "// init initializes the property\n" +
	"func (res *###propLongName###) init(in string) {\n" +
	"###PropInit###" +
	"\treturn\n" +
	"}\n\n"

// PropInitClassBaseThing template
var PropInitClassBaseThing = "\tif obj, ok := model.mThing[in]; ok {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitClassDefault template
var PropInitClassDefault = "\tif obj, ok := model.m###propBaseType###[in]; ok {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitClassImport template
var PropInitClassImport = "\tif temp := model.###capImportName######propBaseType###(in); len(temp) > 0 {\n" +
	"\t\tfor j := range temp {\n" +
	"\t\t\tif temp[j].IRI() == in {\n" +
	"\t\t\t\tres.###Multiplicity######propCapital###(temp[j])\n" +
	"\t\t\t}\n" +
	"\t\t}\n" +
	"\t}\n"

// PropInitTime template
var PropInitTime = "\tif obj, err := time.Parse(\"15:04:05Z07:00\", in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitDate template
var PropInitDate = "\tif obj, err := time.Parse(\"2006-01-02Z07:00\", in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitDateTime template
var PropInitDateTime = "\tif obj, err := time.Parse(time.RFC3339, in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitDateTimeStamp template
var PropInitDateTimeStamp = "\tif obj, err := time.Parse(time.RFC3339, in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitGDay template
var PropInitGDay = "\tif obj, err := time.Parse(\"---02\", in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitGMonth template
var PropInitGMonth = "\tif obj, err := time.Parse(\"--01\", in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitGYear template
var PropInitGYear = "\tif obj, err := time.Parse(\"2006\", in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitGYearMonth template
var PropInitGYearMonth = "\tif obj, err := time.Parse(\"2006-01\", in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitDuration template
var PropInitDuration = "\tif obj, err := owl.ParseXsdDuration(in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitFloat template
var PropInitFloat = "\tif obj, err := strconv.ParseFloat(in, 32); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(float64(obj))\n" +
	"\t}\n"

// PropInitInt template
var PropInitInt = "\tif obj, err := strconv.Atoi(in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitBool template
var PropInitBool = "\tif obj, err := strconv.ParseBool(in); err == nil {\n" +
	"\t\tres.###Multiplicity######propCapital###(obj)\n" +
	"\t}\n"

// PropInitString template
var PropInitString = "\tres.###Multiplicity######propCapital###(in)\n"

// PropInitInterface template
var PropInitInterface = "\tres.###Multiplicity######propCapital###(in)\n"

// MultiplicityMultiple template
var MultiplicityMultiple = "Add"

// MultiplicitySingle template
var MultiplicitySingle = "Set"

// PropertyGetSingle template
var PropertyGetSingle = "// ###propCapital### ###comment###\n" +
	"func (res *###propLongName###) ###propCapital###() (out ###propBaseType###) {\n" +
	"\tout = res.###propName###\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyGetMultipleLiteral template
var PropertyGetMultipleLiteral = "// ###propCapital### ###comment###\n" +
	"func (res *###propLongName###) ###propCapital###() (out []###propBaseType###) {\n" +
	"\tout = res.###propName###\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyGetMultipleClass template
var PropertyGetMultipleClass = "// ###propCapital### ###comment###\n" +
	"func (res *###propLongName###) ###propCapital###() (out []###propBaseType###) {\n" +
	"\tout = make([]###propBaseType###, len(res.###propName###))\n" +
	"\tindex := 0\n" +
	"\tfor i := range res.###propName### {\n" +
	"\t\tout[index] = res.###propName###[i]\n" +
	"\t\tindex++\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// PropertySetSingleLiteral template
var PropertySetSingleLiteral = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *###propLongName###) Set###propCapital###(in ###propBaseType###) (err error) {\n" +
	"\tres.###propName### = in\n" +
	"\treturn\n" +
	"}\n\n"

// PropertySetSingleLiteralMultiple template
var PropertySetSingleLiteralMultiple = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *###propLongName###) Set###propCapital###(in ###propBaseType###) (err error) {\n" +
	"###setSingleLiteralMultiple###" +
	"\terr = errors.New(\"Wrong ###propType### type. Allowed types are ###propAllowedTypes###\")\n" +
	"\treturn\n" +
	"}\n\n"

// SetSingleLiteralMultiple template
var SetSingleLiteralMultiple = "\tif v, ok := in.(###propAllowedType###); ok {\n" +
	"\t\tres.###propName### = v\n" +
	"\t\treturn\n" +
	"\t}\n"

// PropertySetSingleClassSingle template
var PropertySetSingleClassSingle = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *###propLongName###) Set###propCapital###(in ###propBaseType###) (err error) {\n" +
	"\tres.###propName### = in\n" +
	"\treturn\n" +
	"}\n\n"

// PropertySetSingleClassMultiple template
var PropertySetSingleClassMultiple = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *###propLongName###) Set###propCapital###(in ###propBaseType###) (err error) {\n" +
	"###setSingleClassMultiple###" +
	"\terr = errors.New(\"Wrong ###propType### type. Allowed types are ###propAllowedTypes###\")\n" +
	"\treturn\n" +
	"}\n\n"

// SetSingleClassMultiple template
var SetSingleClassMultiple = "\tif v, ok := in.(###propAllowedType###); ok {\n" +
	"\t\tres.###propName### = v\n" +
	"\t\treturn\n" +
	"\t}\n"

// PropertySetMultipleLiteral template
var PropertySetMultipleLiteral = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *###propLongName###) Set###propCapital###(in []###propBaseType###) (err error) {\n" +
	"\tres.###propName### = in\n" +
	"\treturn\n" +
	"}\n\n"

// PropertySetMultipleClass template
var PropertySetMultipleClass = "// Set###propCapital### is setter of ###comment###\n" +
	"func (res *###propLongName###) Set###propCapital###(in []###propBaseType###) (err error) {\n" +
	"\tres.###propName### = make(map[string]###propType###)\n" +
	"\terr = res.Add###propCapital###(in...)\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyAddLiteral template
var PropertyAddLiteral = "// Add###propCapital### adds ###comment###\n" +
	"func (res *###propLongName###) Add###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tres.###propName### = append(res.###propName###, in...)\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyAddClassSingle template
var PropertyAddClassSingle = "// Add###propCapital### adds ###comment###\n" +
	"func (res *###propLongName###) Add###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	// "\t\tif _, ok := res.###propName###[in[i].IRI()]; !ok {\n" +
	"\t\tres.###propName###[in[i].IRI()] = in[i]\n" +
	// "\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyAddClassMultiple template
var PropertyAddClassMultiple = "// Add###propCapital### adds ###comment###\n" +
	"func (res *###propLongName###) Add###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	"###addClassMultiple###" +
	"\t\terr = errors.New(\"Wrong ###propType### type. Allowed types are ###propAllowedTypes###\")\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// AddClassMultiple template
var AddClassMultiple = "\t\tif v, ok := in[i].(###propAllowedType###); ok {\n" +
	// "\t\t\tif _, ok := res.###propName###[v.IRI()]; !ok {\n" +
	"\t\t\tres.###propName###[v.IRI()] = v\n" +
	// "\t\t\t}\n" +
	"\t\tcontinue\n" +
	"\t\t}\n"

// PropertyDelLiteral template
var PropertyDelLiteral = "// Del###propCapital### deletes ###comment###\n" +
	"func (res *###propLongName###) Del###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	"\t\tfor j := range res.###propName### {\n" +
	"\t\t\tif in[i] == res.###propName###[j] {\n" +
	"\t\t\t\tres.###propName### = append(res.###propName###[:j], res.###propName###[j:]...)\n" +
	"\t\t\t\tbreak\n" +
	"\t\t\t}\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyDelClassSingle template
var PropertyDelClassSingle = "// Del###propCapital### deletes ###comment###\n" +
	"func (res *###propLongName###) Del###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	// "\t\tif _, ok := res.###propName###[in[i].IRI()]; ok {\n" +
	"\t\tdelete(res.###propName###, in[i].IRI())\n" +
	// "\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyDelClassMultiple template
var PropertyDelClassMultiple = "// Del###propCapital### deletes ###comment###\n" +
	"func (res *###propLongName###) Del###propCapital###(in ...###propBaseType###) (err error) {\n" +
	"\tfor i := range in {\n" +
	"###delClassMultiple###" +
	"\t\terr = errors.New(\"Wrong ###propType### type. Allowed types are ###propAllowedTypes###\")\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// DelClassMultiple template
var DelClassMultiple = "\t\tif v, ok := in[i].(###propAllowedType###); ok {\n" +
	// "\t\t\tif _, ok := res.###propName###[v.IRI()]; ok {\n" +
	"\t\t\tdelete(res.###propName###, v.IRI())\n" +
	// "\t\t\t}\n" +
	"\t\tcontinue\n" +
	"\t\t}\n"

// PropertyGraphSingle template
var PropertyGraphSingle = "// toGraph adds all predicates corresponding to the property to an owl graph\n" +
	"func (res *###propLongName###) toGraph(node *rdf.Node, g *rdf.Graph) {\n" +
	"###graphProp###" +
	"\treturn\n" +
	"}\n\n"

// PropertyGraphMultiple template
var PropertyGraphMultiple = "// toGraph adds all predicates corresponding to the property to an owl graph\n" +
	"func (res *###propLongName###) toGraph(node *rdf.Node, g *rdf.Graph) {\n" +
	"\tfor i := range res.###propName### {\n" +
	"###graphProp###" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyMultipleRemove template
var PropertyMultipleRemove = "// removeObject removes object from property\n" +
	"func (res *###propLongName###) removeObject(obj owl.Thing) {\n" +
	"\tif v, ok := obj.(###propBaseType###); ok {\n" +
	"\t\tres.Del###propCapital###(v)\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// PropertySingleRemove template
var PropertySingleRemove = "// removeObject removes object from property\n" +
	"func (res *###propLongName###) removeObject(obj owl.Thing) {\n" +
	"\tif v, ok := obj.(###propBaseType###); ok {\n" +
	"\t\tif res.###propName###.IRI() == v.IRI() {\n" +
	"\t\t\tres.Set###propCapital###(nil)\n" +
	"\t\t}\n" +
	"\t}\n" +
	"\treturn\n" +
	"}\n\n"

// GraphPropString template
var GraphPropString = "###indent###\towl.AddStringPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropFloat template
var GraphPropFloat = "###indent###\towl.AddFloatPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropInt template
var GraphPropInt = "###indent###\towl.AddIntPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropBool template
var GraphPropBool = "###indent###\towl.AddBoolPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropInterface template
var GraphPropInterface = "###indent###\towl.AddInterfacePropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSTime template
var GraphPropSTime = "###indent###\towl.AddTimePropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSDateTime template
var GraphPropSDateTime = "###indent###\towl.AddDateTimePropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSDate template
var GraphPropSDate = "###indent###\towl.AddDatePropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSDuration template
var GraphPropSDuration = "###indent###\towl.AddDurationPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSDateTimeStamp template
var GraphPropSDateTimeStamp = "###indent###\towl.AddDateTimeStampPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSGYear template
var GraphPropSGYear = "###indent###\towl.AddGYearPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSGDay template
var GraphPropSGDay = "###indent###\towl.AddGDayPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSGYearMonth template
var GraphPropSGYearMonth = "###indent###\towl.AddGYearMonthPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropSGMonth template
var GraphPropSGMonth = "###indent###\towl.AddGMonthPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// GraphPropClass template
var GraphPropClass = "###indent###\towl.AddClassPropertyToGraph(g, \"###propIRI###\", node, res.###propName######array###)\n"

// PropertyStringSingle template
var PropertyStringSingle = "// String prints the object into a string\n" +
	"func (res *###propLongName###) String() (ret string) {\n" +
	"\tret += \"\\t###propName###: [\"\n" +
	"###stringProp###" +
	"\tret += \"]\\n\"\n" +
	"\treturn\n" +
	"}\n\n"

// PropertyStringMultiple template
var PropertyStringMultiple = "// String prints the object into a string\n" +
	"func (res *###propLongName###) String() (ret string) {\n" +
	"\tret += \"\\t###propName###: [\"\n" +
	"\tfor i := range res.###propName### {\n" +
	"###stringProp###" +
	"\t}\n" +
	"\tret += \"]\\n\"\n" +
	"\treturn\n" +
	"}\n\n"

// StringPropString template
var StringPropString = "###indent###\tret += res.###propName######array### + \", \"\n"

// StringPropFloat template
var StringPropFloat = "###indent###\tret += fmt.Sprintf(\"%f\", res.###propName######array###) + \", \"\n"

// StringPropInt template
var StringPropInt = "###indent###\tret += fmt.Sprintf(\"%d\", res.###propName######array###) + \", \"\n"

// StringPropBool template
var StringPropBool = "###indent###\tret += fmt.Sprintf(\"%v\", res.###propName######array###) + \", \"\n"

// StringPropInterface template
var StringPropInterface = "###indent###\tret += fmt.Sprintf(\"%v\", res.###propName######array###) + \", \"\n"

// StringPropTime template
var StringPropTime = "###indent###\tret += res.###propName######array###.String() + \", \"\n"

// StringPropClassSingle template
var StringPropClassSingle = "\tif res.###propName### != nil {\n" +
	"\t\tret += res.###propName###.IRI() + \", \"\n" +
	"\t}\n"

// StringPropClassMultiple template
var StringPropClassMultiple = "\t\tret += res.###propName###[i].IRI() + \", \"\n"

// ArraySingle template
var ArraySingle = ""

// ArrayMultiple template
var ArrayMultiple = "[i]"

// IndentSingle template
var IndentSingle = ""

// IndentMultiple template
var IndentMultiple = "\t"
