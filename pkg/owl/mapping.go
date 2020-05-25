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

package owl

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// GoModel holds all classes
type GoModel struct {
	Class       map[string]GoClass // map className->Class
	Individual  []GoIndividual     // all individuals
	IRI         string             // ontology iri
	Description string             // description extracted from ontology
	Name        string             // name (part of iri)
	Content     []byte             // ontology in ttl
	Module      string             // Go module name
}

// GoClass holds properties of a class
type GoClass struct {
	IRI          string       // iri of class
	Name         string       // class name
	Parent       []string     // parent classes (also transitive)
	DirectParent []string     // direct parent classes
	Property     []GoProperty // properties of class
	Comment      string       // comment for doc
	Model        *GoModel     // pointer to model
}

// GoProperty holds condiguration of a property
type GoProperty struct {
	IRI          string      // IRI
	Name         string      // name
	Capital      string      // name with capital first letter
	Typ          [2]string   // type in struct (child of BaseTyp, 0: typ; 1: IRI)
	BaseTyp      [2]string   // type in interface
	AllowedTyp   [][2]string // allowed types (children of Typ)
	XSDTyp       string      // XSD type if any
	Individual   []string    // predefined Individuals
	Multi        bool        // multiple values possible?
	Multiplicity string      // multiplicity (slice or single value)
	Comment      string      // comment for doc
	Inverse      string      // inverse property if any
	Symmetric    bool        // symmetric?
}

// GoIndividual individuals
type GoIndividual struct {
	IRI  string // Individual iri
	Name string // name of individual
	Typ  string // type of individual (class)
}

// MapModel extracts the model from an ontology
func MapModel(ont *Ontology, moduleName string) (mod GoModel, err error) {
	fmt.Println("Map ontology to Go model")
	mod.IRI = ont.IRI
	mod.Description = ont.Description[ont.IRI]
	mod.Module = moduleName
	temp := strings.Split(ont.IRI, "/")
	mod.Name = temp[len(temp)-1]
	mod.Class = make(map[string]GoClass)

	for i := range ont.Imports {
		mod.createGoClasses(ont, i)

		for j := range ont.Individual {
			if strings.HasPrefix(ont.Individual[j].Type.Name, i+"#") {
				temp := mod.extractIndividual(ont.Individual[j], ont)
				if temp.IRI != "" {
					mod.Individual = append(mod.Individual, temp)
				}
			}
		}
	}

	return
}

// createGoClasses creates all necessary GoClasses and fills their information if possible
func (mod *GoModel) createGoClasses(ont *Ontology, ontName string) (err error) {
	fmt.Println("\tCreate Go classes for " + ontName)

	for i := range ont.Class {
		var temp GoClass
		temp, err = mod.extractClass(ont.Class[i], ont, ontName)
		if err != nil {
			return
		}
		if temp.Name != "" {
			mod.Class[temp.Name] = temp
		}
	}
	return
}

// extractClass extracts a single class
func (mod *GoModel) extractClass(class *Class, ont *Ontology, ontName string) (goClass GoClass,
	err error) {

	// get class name and IRI
	if !strings.HasPrefix(class.Name, ontName+"#") {
		//err = errors.New("wrong iri: " + class.Name)
		return
	}
	goClass.Name = trimName(class.Name, ont)
	goClass.IRI = class.Name

	// get comment
	goClass.Comment = class.Comment

	// get parents
	parents := class.GetAllParents()
	for i := range parents {
		parentName := trimName(parents[i].Name, ont)
		if parentName == "" {
			err = errors.New("Class " + class.Name + ": wrong parent: " + parents[i].Name)
			return
		}
		if parentName != goClass.Name {
			goClass.Parent = append(goClass.Parent, parentName)
		}
	}
	for i := range class.Parent {
		parentName := trimName(class.Parent[i].Name, ont)
		if parentName == "" {
			err = errors.New("Class " + class.Name + ": wrong parent: " + class.Parent[i].Name)
			return
		}
		if parentName != goClass.Name {
			goClass.DirectParent = append(goClass.DirectParent, parentName)
		}
	}

	// get properties
	restInv := class.GetRestrictionsInverse()
	rest := class.GetRestrictions()
	for i := range restInv {
		var property GoProperty
		property.Name, property.Capital, property.IRI = getRestrictionNameAndIRI(restInv[i], ont)
		if property.Name == "" {
			err = errors.New("Class " + class.Name + " unknown property " +
				restInv[i].Property.Name)
			return
		}
		if _, ok := ont.Property[restInv[i].Property.Name]; ok {
			property.Comment = ont.Property[restInv[i].Property.Name].Comment
		}
		var exist bool
		property.BaseTyp, exist = getRestrictionType(restInv[i], ont)
		if !exist {
			err = errors.New("Class " + class.Name + " Restriction " + restInv[i].Property.Name +
				" unknown base type " + property.BaseTyp[0])
			return
		}

		property.Multi, property.Multiplicity = getRestrictionMultiplicity(restInv[i])
		if restInv[i].Property.Inverse != nil {
			property.Inverse = trimName(restInv[i].Property.Inverse.Name, ont)
		}

		for j := range rest {
			if rest[j].Property.Name == restInv[i].Property.Name {
				property.Typ, exist = getRestrictionType(rest[j], ont)
				if !exist || property.Typ[0] == "string" {
					// fmt.Println("Warning: Class " + class.Name + " Restriction " + rest[j].Property.Name + " unknown type " + property.Typ)
					property.Typ = property.BaseTyp
				}

				//if property.Typ == "string" || property.Typ == "int" || property.Typ == "float64" ||
				if property.Typ[0] == "time.Time" || property.Typ[0] == "time.Duration" {
					if len(rest[j].Value) > 0 {
						property.XSDTyp = rest[j].Value[0]
					} else if len(restInv[i].Value) > 0 {
						property.XSDTyp = restInv[i].Value[0]
					} else {
						err = errors.New("Class " + class.Name + " Restriction " +
							restInv[i].Property.Name + " no xsd type")
						return
					}
				}
				allowedType := getRestrictionAllowedTypes(rest[j], ont)
				if len(allowedType) == 0 {
					property.AllowedTyp = append(property.AllowedTyp, property.Typ)
				} else {
					for k := range allowedType {
						property.AllowedTyp = append(property.AllowedTyp, allowedType[k])
					}
				}
				property.Individual, err = getIndividuals(rest[j], ont)
			}
		}
		b, _ := GetBaseClass([]string{property.BaseTyp[1], property.Typ[1]}, ont.Class)
		if property.BaseTyp[0] != property.Typ[0] && b == nil {
			if property.Typ[0] == "string" || property.Typ[0] == "time.Time" ||
				property.Typ[0] == "time.Duration" || property.Typ[0] == "int" ||
				property.Typ[0] == "float64" || property.Typ[0] == "bool" {

			} else {
				// WARNING: This removes properties that are present in specification (relevant for saref4ener:PowerProfile:consists of; saref4ener:AlternativesGroup does not inherit from saref:Profile)
				continue
			}
		}
		goClass.Property = append(goClass.Property, property)
	}

	return
}

// getRestrictionNameAndIRI returns the name and iri of a restriction
func getRestrictionNameAndIRI(rest *Restriction, ont *Ontology) (name string, capital string,
	iri string) {
	name = trimName(rest.Property.Name, ont)
	a := []rune(name)
	a[0] = unicode.ToLower(a[0])
	name = string(a)
	iri = rest.Property.Name
	capital = strings.Title(name)
	return
}

// getRestrictionMultiplicity returns an indicator if restrictions allows multipe values and a
// string for slice or array creation
func getRestrictionMultiplicity(rest *Restriction) (multi bool, multiplicity string) {
	if rest.ValueConstraint == "http://www.w3.org/2002/07/owl#allValuesFrom" ||
		rest.ValueConstraint == "http://www.w3.org/2002/07/owl#someValuesFrom" ||
		rest.CardinalityConstraint == "http://www.w3.org/2002/07/owl#minQualifiedCardinality" ||
		rest.CardinalityConstraint == "http://www.w3.org/2002/07/owl#minCardinality" {
		multi = true
		multiplicity = "[]"
	} else if rest.CardinalityConstraint == "http://www.w3.org/2002/07/owl#maxCardinality" ||
		rest.CardinalityConstraint == "http://www.w3.org/2002/07/owl#qualifiedCardinality" ||
		rest.CardinalityConstraint == "http://www.w3.org/2002/07/owl#maxQualifiedCardinality" ||
		rest.CardinalityConstraint == "http://www.w3.org/2002/07/owl#cardinality" {
		if rest.Multiplicity > 1 {
			multi = true
			multiplicity = "[" + strconv.Itoa(rest.Multiplicity) + "]"
		} else {
			multi = false
			multiplicity = ""
		}
	}
	return
}

// getRestrictionAllowedTypes returns the allowed types of a restriction
func getRestrictionAllowedTypes(rest *Restriction, ont *Ontology) (values [][2]string) {
	if rest.ValueConstraint == "http://www.w3.org/2002/07/owl#hasValue" {
		base, err := GetBaseType(rest.Value, ont.Individual, ont.Class)
		if err != nil {
			return
		}
		if base != nil {
			allowedType := trimName(base.Name, ont)
			values = append(values, [2]string{allowedType, base.Name})
		} else {
			fmt.Println("unknown basetype for " + rest.Property.Name)
		}
	} else {
		for i := range rest.Value {
			if _, ok := ont.Class[rest.Value[i]]; ok {
				allowedType := trimName(rest.Value[i], ont)
				values = append(values, [2]string{allowedType, rest.Value[i]})
			} else if strings.HasPrefix(rest.Value[i], "http://www.w3.org/2001/XMLSchema") {
				allowedType, err := mapLiteralType(rest.Value[i])
				if err == nil {
					values = append(values, [2]string{allowedType, ""})
				}
			}
		}
	}
	return
}

// getRestrictionType returns restriction type
func getRestrictionType(rest *Restriction, ont *Ontology) (ret [2]string, typeExist bool) {
	typeExist = false
	if rest.ValueConstraint == "http://www.w3.org/2002/07/owl#hasValue" {
		base, err := GetBaseType(rest.Value, ont.Individual, ont.Class)
		if err != nil {
			err = errors.New("Restriction " + rest.Property.Name + " " + fmt.Sprint(err))
			return
		}
		if base != nil {
			ret[0] = trimName(base.Name, ont)
		} else {
			ret[0] = "owl.Thing"
		}
		typeExist = true
	} else {
		isClass := false
		isLiteral := false
		for i := range rest.Value {
			if temp := trimName(rest.Value[i], ont); temp != "" {
				isClass = true
			} else if strings.HasPrefix(rest.Value[i], "http://www.w3.org/2001/XMLSchema") {
				isLiteral = true
			}
		}
		if isClass && !isLiteral {
			base, err := GetBaseClass(rest.Value, ont.Class)
			if err != nil {
				err = errors.New("Restriction " + rest.Property.Name + " " + fmt.Sprint(err))
				return
			}
			if base != nil {
				ret[0] = trimName(base.Name, ont)
				ret[1] = base.Name
			} else {
				ret[0] = "owl.Thing"
			}
			typeExist = true
		} else if isLiteral && !isClass {
			var err error
			ret[0], err = mapLiteralType(rest.Value[0])
			if err == nil {
				typeExist = true
			}
		} else {
			ret[0] = "interface{}"
			typeExist = true
		}
	}

	return
}

// getIndividuals returns all individuals of a restrictions
func getIndividuals(rest *Restriction, ont *Ontology) (inds []string, err error) {
	if rest.ValueConstraint == "http://www.w3.org/2002/07/owl#hasValue" {
		for i := range rest.Value {
			inds = append(inds, rest.Value[i])
		}
	}
	return
}

// extractIndividual
func (mod *GoModel) extractIndividual(individual *Individual,
	ont *Ontology) (goIndividual GoIndividual) {
	if temp := trimName(individual.Name, ont); temp != "" {
		goIndividual.Name = temp
	} else if strings.HasPrefix(individual.Name, "http://www.wurvoc.org/vocabularies/om-1.8") {
		goIndividual.Name = strings.TrimPrefix(individual.Name,
			"http://www.wurvoc.org/vocabularies/om-1.8/")
	}
	if goIndividual.Name == "" {
		return
	}
	temp := strings.Split(goIndividual.Name, "_")
	for i := range temp {
		temp[i] = strings.Title(temp[i])
	}
	goIndividual.Name = strings.Join(temp, "")
	goIndividual.IRI = individual.Name
	goIndividual.Typ = trimName(individual.Type.Name, ont)
	return
}

// mapLiteralType maps the literal type to a go datatype
func mapLiteralType(literal string) (goType string, err error) {
	switch literal {
	case "http://www.w3.org/2001/XMLSchema#time":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#dateTime":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#date":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#duration":
		goType = "time.Duration"
	case "http://www.w3.org/2001/XMLSchema#dateTimeStamp":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#gYear":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#gDay":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#gYearMonth":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#gMonth":
		goType = "time.Time"
	case "http://www.w3.org/2001/XMLSchema#string":
		goType = "string"
	case "http://www.w3.org/2001/XMLSchema#float":
		goType = "float64"
	case "http://www.w3.org/2001/XMLSchema#decimal":
		goType = "float64"
	case "http://www.w3.org/2001/XMLSchema#integer":
		goType = "int"
	case "http://www.w3.org/2001/XMLSchema#nonNegativeInteger":
		goType = "int"
	case "http://www.w3.org/2001/XMLSchema#unsignedInt":
		goType = "int"
	case "http://www.w3.org/2001/XMLSchema#boolean":
		goType = "bool"
	default:
		err = errors.New("Unknown literal")
	}
	return
}

// trimName trims long names (iris) to short ones
func trimName(name string, ont *Ontology) (out string) {
	for i := range ont.Imports {
		if strings.HasPrefix(name, i+"#") {
			out = strings.Title(strings.TrimPrefix(name, i+"#"))
			temp := strings.Split(i, "/")
			pref := strings.Title(temp[len(temp)-1])
			if !strings.HasPrefix(out, pref) {
				out = pref + out
			}
			return
		}
	}
	out = ""
	return
}
