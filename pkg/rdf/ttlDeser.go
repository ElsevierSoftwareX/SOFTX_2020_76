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

package rdf

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

// token represents one ttl expression
type token struct {
	typ   int    // token type
	value string // value of token
}

// parser parses a ttl document
type parser struct {
	reader       *bufio.Reader        // reader of turtle document
	runes        []rune               // turtle document as rune slice
	posStatement int                  // starting position of current statement
	prefix       map[string]string    // prefixes
	base         string               // base IRI
	curSubject   string               // current Subject
	curPredicate string               // current Predicate
	triples      []Triple             // list of all extracted triples
	bnCounter    int                  // blank node counter
	blank        map[string]BlankNode // blankNode map
}

// predObjList is a Predicate Object List
type predObjList struct {
	pred Predicate
	obj  []Object
}

// DecodeTTL decodes a ttl input to rdf triples
func DecodeTTL(input io.Reader) (trip []Triple, err error) {
	p := &parser{reader: bufio.NewReader(input), prefix: make(map[string]string),
		blank: make(map[string]BlankNode)}
	err = p.parseRunes()
	if err != nil {
		return
	}
	// remove whitepaces at beginning
	p.posStatement = p.consumeWS(0)
	for {
		err = p.parseStatement()
		if err != nil {
			break
		}
		if p.posStatement >= len(p.runes) {
			break
		}
	}
	trip = p.triples

	return
}

// parseRunes parses all runes from the reader and omits empty lines and comments
func (p *parser) parseRunes() (err error) {
	eof := false
	for {
		var line []byte
		line, err = p.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.New("Error parsing runes: " + err.Error())
				return
			}
			eof = true
		}
		if len(line) == 0 {
			if eof {
				break
			}
			continue
		}
		pos := 0
		// omit new line at end of slice (but only if not last line)
		stop := len(line) - 1
		if eof {
			stop = len(line)
		}
		for pos < stop {
			var r rune
			var s int
			r, s = utf8.DecodeRune(line[pos:])
			if r == utf8.RuneError {
				err = errors.New("Error parsing runes: Rune error")
			}
			if pos == 0 && r == '#' {
				break
			}
			if r == '\t' {
				r = ' '
			} else if r == '\n' {
				r = ' '
				// continue
			} else if r == '\r' {
				r = ' '
			}
			p.runes = append(p.runes, r)
			pos += s
		}
		if eof {
			break
		}
	}
	return
}

// parseStatement decodes one statement beginning from the position stored in parser
func (p *parser) parseStatement() (err error) {
	if len(p.runes) <= p.posStatement {
		return
	}
	length := 0
	switch p.runes[p.posStatement] {
	case '@':
		// @prefix or @base
		length, err = p.parseDirective(p.posStatement + 1)
		length++
	default:
		if p.isEqual(p.posStatement, "BASE ") || p.isEqual(p.posStatement, "PREFIX ") {
			// sparqlBase or sparqlPrefix
			length, err = p.parseDirective(p.posStatement)
		} else {
			length, err = p.parseTriples(p.posStatement)
		}
	}
	if err != nil {
		return
	}
	p.posStatement += length
	return
}

// parseDirective decodes one directive (prefix, base, sparql prefix or sparql base) beginning from
// the specified position
func (p *parser) parseDirective(pos int) (length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of directive")
		return
	}
	switch p.runes[pos] {
	case 'p':
		var prefix, iri string
		var tempLength int
		if p.isEqual(pos, "prefix") {
			length = 6
			length += p.consumeWS(pos + length)
			prefix, tempLength, err = p.parsePrefix(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			length += p.consumeWS(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.prefix[prefix] = iri
		} else {
			err = errors.New("wrong prefix")
			return
		}
	case 'b':
		var iri string
		var tempLength int
		if p.isEqual(pos, "base") {
			length = 4
			length += p.consumeWS(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.base = iri
		} else {
			err = errors.New("wrong base")
			return
		}
	case 'P':
		var prefix, iri string
		var tempLength int
		if p.isEqual(pos, "PREFIX") {
			length = 6
			length += p.consumeWS(pos + length)
			prefix, tempLength, err = p.parsePrefix(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			length += p.consumeWS(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.prefix[prefix] = iri
		} else {
			err = errors.New("wrong PREFIX")
			return
		}
	case 'B':
		var iri string
		var tempLength int
		if p.isEqual(pos, "BASE") {
			length = 4
			length += p.consumeWS(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.base = iri
		} else {
			err = errors.New("wrong BASE")
			return
		}
	default:
		err = errors.New("Invalid directive " + string(p.runes[pos]))
		return
	}
	// consumer dot
	length += p.consumeWS(pos + length)
	if p.isEqual(pos+length, ".") {
		length++
	} else {
		err = errors.New("missing dot at end of directive")
		return
	}
	length += p.consumeWS(pos + length)
	return
}

// parsePrefix parses one prefix beginning from the specified position
func (p *parser) parsePrefix(pos int) (prefix string, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of prefix")
	}
	prefix, length, err = p.parseUntil(pos, ':')
	length++
	return
}

// parseTriples parses all tripels (subject predicateObjectList) in a statement and adds them to
// p.triples
func (p *parser) parseTriples(pos int) (length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of triples")
		return
	}
	var trip Triple

	// subject
	trip.Sub, length, err = p.parseSubject(pos)
	if err != nil {
		return
	}
	length += p.consumeWS(pos + length)

	// predicateObjectList
	var tempLength int
	var poList []predObjList
	poList, tempLength, err = p.parsePredicateObjectList(pos + length)
	if err != nil {
		return
	}
	length += tempLength
	length += p.consumeWS(pos + length)

	// add triples
	for i := range poList {
		trip.Pred = poList[i].pred
		for j := range poList[i].obj {
			trip.Obj = poList[i].obj[j]
			p.triples = append(p.triples, trip)
		}
	}

	// consume dot
	length += p.consumeWS(pos + length)
	if p.isEqual(pos+length, ".") {
		length++
	} else {
		err = errors.New("missing dot at end of triples with subject " + trip.Sub.String())
		return
	}
	length += p.consumeWS(pos + length)

	return
}

// parseSubject parses the subject (iri | BlankNode | collection) of a triple
func (p *parser) parseSubject(pos int) (subj Subject, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("reached eof before end of subject")
		return
	}
	if p.runes[pos] == '_' && p.runes[pos+1] == ':' {
		subj, length, err = p.parseBlankNode(pos)
	} else if p.runes[pos] == '(' {
		subj, length, err = p.parseCollection(pos)
	} else if p.runes[pos] == '[' {
		subj, length, err = p.parseBlankNodePropertyList(pos)
	} else {
		subj, length, err = p.parseIRI(pos)
	}
	return
}

// parsePredicateObjectList parses a predicateObjectList (verb objectList (';' (verb objectList)?)*)
func (p *parser) parsePredicateObjectList(pos int) (list []predObjList, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of predicate object list")
		return
	}
	length = 0

	for {
		var poListTemp predObjList
		var tempLength int
		// verb
		poListTemp.pred, tempLength, err = p.parsePredicate(pos + length)
		if err != nil {
			return
		}
		length += tempLength
		length += p.consumeWS(pos + length)

		//objectList
		poListTemp.obj, tempLength, err = p.parseObjectList(pos + length)
		if err != nil {
			return
		}
		length += tempLength
		length += p.consumeWS(pos + length)
		list = append(list, poListTemp)

		if !p.isEqual(pos+length, ";") {
			break
		}
		length++
		length += p.consumeWS(pos + length)
		if p.isEqual(pos+length, ".") || p.isEqual(pos+length, "]") {
			break
		}
	}
	return
}

// parsePredicate parses the next predicate
func (p *parser) parsePredicate(pos int) (pred Predicate, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of predicate")
		return
	}
	var temp string
	temp, _, err = p.parseUntil(pos+length, ' ')
	if err != nil {
		return
	}
	if temp == "a" {
		pred = IRI{name: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"}
		length = 1
	} else {
		pred, length, err = p.parseIRI(pos)
	}
	return
}

// parseObjectList parses an objectList (object (',' object)*)
func (p *parser) parseObjectList(pos int) (obj []Object, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of object list")
		return
	}
	for {
		var temp Object
		var tempLength int
		temp, tempLength, err = p.parseObject(pos + length)
		if err != nil {
			return
		}
		obj = append(obj, temp)
		length += tempLength
		length += p.consumeWS(pos + length)
		if !p.isEqual(pos+length, ",") {
			break
		}
		length++
		length += p.consumeWS(pos + length)
	}
	return
}

// parseObject parses one object (iri | BlankNode | collection | blankNodePropertyList | literal)
func (p *parser) parseObject(pos int) (obj Object, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of object")
		return
	}
	if p.runes[pos] == '_' && p.runes[pos+1] == ':' {
		obj, length, err = p.parseBlankNode(pos)
	} else if p.runes[pos] == '(' {
		obj, length, err = p.parseCollection(pos)
	} else if p.runes[pos] == '[' {
		obj, length, err = p.parseBlankNodePropertyList(pos)
	} else if p.runes[pos] == '+' || p.runes[pos] == '-' || p.runes[pos] == '0' ||
		p.runes[pos] == '1' || p.runes[pos] == '2' || p.runes[pos] == '3' || p.runes[pos] == '4' ||
		p.runes[pos] == '5' || p.runes[pos] == '6' || p.runes[pos] == '7' || p.runes[pos] == '8' ||
		p.runes[pos] == '9' || p.runes[pos] == '"' || p.runes[pos] == '\'' ||
		p.isEqual(pos, "true") || p.isEqual(pos, "false") || p.runes[pos] == '.' {
		obj, length, err = p.parseLiteral(pos)
	} else {
		obj, length, err = p.parseIRI(pos)
	}
	return
}

// parseIRI parses the next iri (IRIRef | prefixedName)
func (p *parser) parseIRI(pos int) (iri IRI, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of iri")
		return
	}
	var i string
	if p.runes[pos] == '<' {
		i, length, err = p.parseIRIRef(pos)
	} else {
		i, length, err = p.parsePrefixedName(pos)
	}
	if err != nil {
		return
	}
	iri = IRI{name: i}
	return
}

// parseIRI parses IRIRef (<iri>)
func (p *parser) parseIRIRef(pos int) (iri string, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of iri")
		return
	}
	if p.runes[pos] != '<' {
		err = errors.New("No IRI; missing <")
		return
	}
	iri, length, err = p.parseUntil(pos+1, '>')
	length += 2
	return
}

// parsePrefixedName parses prefixed name (prefix:name)
func (p *parser) parsePrefixedName(pos int) (iri string, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of iri")
		return
	}
	var prefix string
	prefix, length, err = p.parsePrefix(pos)
	if err != nil {
		return
	}
	ok := false
	if iri, ok = p.prefix[prefix]; !ok {
		err = errors.New("no such prefix " + prefix)
	}
	var name string
	var tempLength int
	name, tempLength, err = p.parseUntil(pos+length, ' ')
	if strings.ContainsRune(name, ',') {
		name, tempLength, err = p.parseUntil(pos+length, ',')
	} else if strings.ContainsRune(name, ';') {
		name, tempLength, err = p.parseUntil(pos+length, ';')
	} else if strings.ContainsRune(name, '.') {
		name, tempLength, err = p.parseUntil(pos+length, '.')
	}
	if err != nil {
		return
	}
	iri = iri + name
	length += tempLength
	return
}

// parseLiteral parses a literal (RDFLiteral | NumericLiteral | BooleanLiteral)
func (p *parser) parseLiteral(pos int) (lit Literal, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of literal")
		return
	}
	if p.runes[pos] == '"' || p.runes[pos] == '\'' {
		lit, length, err = p.parseRDFLiteral(pos)
	} else if p.isEqual(pos, "true") || p.isEqual(pos, "false") {
		lit, length, err = p.parseBooleanLiteral(pos)
	} else {
		lit, length, err = p.parseNumericLiteral(pos)
	}
	return
}

// parseRDFLiteral parses a rdf literal (String (LANGTAG | '^^' iri))
func (p *parser) parseRDFLiteral(pos int) (lit Literal, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("reached eof before end of literal")
		return
	}
	if p.runes[pos] == '"' {
		if p.runes[pos+1] == '"' {
			lit.str, length, err = p.parseUntil(pos+3, '"')
			for {
				if p.runes[pos+3+length+1] != '"' {
					var tempLength int
					var tempStr string
					tempStr, tempLength, err = p.parseUntil(pos+3+length+1, '"')
					lit.str += "\"" + tempStr
					length += 1 + tempLength
					continue
				}
				break
			}
			length += 6
		} else {
			lit.str, length, err = p.parseUntil(pos+1, '"')
			length += 2
		}
	} else if p.runes[pos] != '\'' {
		if p.runes[pos+1] == '\'' {
			lit.str, length, err = p.parseUntil(pos+3, '\'')
			for {
				if p.runes[pos+3+length+1] != '\'' {
					var tempLength int
					var tempStr string
					tempStr, tempLength, err = p.parseUntil(pos+3+length+1, '\'')
					lit.str += "'" + tempStr
					length += 1 + tempLength
					continue
				}
				break
			}
			length += 6
		} else {
			lit.str, length, err = p.parseUntil(pos+1, '\'')
			length += 2
		}
	} else {
		err = errors.New("no rdf literal; missing quotes")
	}
	if err != nil {
		return
	}
	if len(p.runes) <= pos+length {
		return
	}
	// langtag
	if p.runes[pos+length] == '@' {
		length++
		var tempLength int
		lit.langTag, tempLength, err = p.parseUntil(pos+length, ' ')
		if strings.ContainsRune(lit.langTag, ',') {
			lit.langTag, tempLength, err = p.parseUntil(pos+length, ',')
		} else if strings.ContainsRune(lit.langTag, ';') {
			lit.langTag, tempLength, err = p.parseUntil(pos+length, ';')
		} else if strings.ContainsRune(lit.langTag, '.') {
			lit.langTag, tempLength, err = p.parseUntil(pos+length, '.')
		}
		if err != nil {
			return
		}
		length += tempLength
	}
	// type iri
	if p.runes[pos+length] == '^' {
		length += 2
		var tempLength int
		var iriTemp IRI
		iriTemp, tempLength, err = p.parseIRI(pos + length)
		if err != nil {
			return
		}
		lit.typeIRI = iriTemp.name
		length += tempLength
	}

	return
}

// parseBooleanLiteral parses a boolean literal
func (p *parser) parseBooleanLiteral(pos int) (lit Literal, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of boolean literal")
		return
	}
	if p.isEqual(pos, "true") {
		lit = Literal{str: "true", value: true}
		length = 4
	} else if p.isEqual(pos, "false") {
		lit = Literal{str: "false", value: false}
		length = 5
	} else {
		err = errors.New("no boolean literal")
	}
	return
}

// parseNumericLiteral parses a numeric literal
func (p *parser) parseNumericLiteral(pos int) (lit Literal, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("reached eof before end of numerical literl")
		return
	}
	// get length of literal
	var value float64
	tempLength := 0
	lit.str, length, err = p.parseUntil(pos, ' ')
	if err != nil {
		return
	}
	// could be object list
	if p.runes[pos+length-1] == ',' {
		lit.str, length, err = p.parseUntil(pos, ',')
		if err != nil {
			return
		}
	}

	// + or -
	if p.runes[pos] == '-' {
		value = -1
		tempLength++
	} else if p.runes[pos] == '+' {
		value = 1
		tempLength++
	} else {
		value = 1
	}

	// look for number before dot or exp
	numLen := 0
	var num []rune
	for i := pos + tempLength; i < pos+length; i++ {
		if p.runes[i] == '.' || p.runes[i] == 'e' || p.runes[i] == 'E' {
			break
		}
		num = append(num, p.runes[i])
		numLen++
	}
	if numLen > 0 {
		var temp int
		temp, err = strconv.Atoi(string(num))
		if err != nil {
			return
		}
		value = value * float64(temp)
	} else {
		value = 0
	}
	tempLength += numLen

	// is integer?
	if tempLength == length {
		lit.value = int(value)
		lit.typeIRI = "http://www.w3.org/2001/XMLSchema#integer"
		return
	}

	// look for dot
	if p.runes[pos+tempLength] == '.' {
		tempLength++
		numLen = 0
		num = nil
		for i := pos + tempLength; i < pos+length; i++ {
			if p.runes[i] == '.' || p.runes[i] == 'e' || p.runes[i] == 'E' {
				break
			}
			num = append(num, p.runes[i])
			numLen++
		}
		if numLen > 0 {
			var temp int
			temp, err = strconv.Atoi(string(num))
			if err != nil {
				return
			}
			// divisor
			div := 1
			for i := 0; i < numLen; i++ {
				div = div * 10
			}
			value = value + (float64(temp) / float64(div))
		}
		tempLength += numLen
	}
	if tempLength == length {
		lit.value = value
		lit.typeIRI = "http://www.w3.org/2001/XMLSchema#decimal"
		return
	}

	// look for exp
	if p.runes[pos+tempLength] == 'E' || p.runes[pos+tempLength] == 'e' {
		tempLength++
		exp := 1.0
		// + or -?
		if p.runes[pos+tempLength] == '-' {
			exp = -exp
			tempLength++
		} else if p.runes[pos+tempLength] == '+' {
			tempLength++
		}
		numLen = 0
		num = nil
		for i := pos + tempLength; i < pos+length; i++ {
			num = append(num, p.runes[i])
			numLen++
		}
		if numLen > 0 {
			var temp int
			temp, err = strconv.Atoi(string(num))
			if err != nil {
				return
			}
			// multiplicator
			mul := 1
			for i := 0; i < temp; i++ {
				mul = mul * 10
			}
			value = value * float64(mul)
		}
		tempLength += numLen
	}

	lit.value = value
	lit.typeIRI = "http://www.w3.org/2001/XMLSchema#double"
	return
}

// parseBlankNode parses a blank node
func (p *parser) parseBlankNode(pos int) (blank BlankNode, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("reached eof before end of blank node")
		return
	}
	if p.runes[pos] != '_' || p.runes[pos+1] != ':' {
		err = errors.New("no blank node")
		return
	}
	var blankName string
	blankName, length, err = p.parseUntil(pos+2, ' ')
	if err != nil {
		return
	}
	length += 2
	var ok bool
	blank, ok = p.blank[blankName]
	if !ok {
		blank = p.blankNode()
		p.blank[blankName] = blank
	}
	return
}

// parseCollection parses a collection ('(' object* ')')
func (p *parser) parseCollection(pos int) (blank BlankNode, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("reached eof before end of collection")
		return
	}
	blank = p.blankNode()
	if p.runes[pos] != '(' {
		err = errors.New("no collection; missing (")
		return
	}
	length = p.consumeWS(pos + 1)
	length++
	var trip Triple
	trip.Sub = blank
	for {
		var tempLength int
		trip.Pred = IRI{name: "http://www.w3.org/1999/02/22-rdf-syntax-ns#first"}
		trip.Obj, tempLength, err = p.parseObject(pos + length)
		if err != nil {
			return
		}
		length += tempLength
		p.triples = append(p.triples, trip)

		tempLength = p.consumeWS(pos + length)
		length += tempLength
		if p.isEqual(pos+length, ")") {
			break
		}
		trip.Pred = IRI{name: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"}
		bNext := p.blankNode()
		trip.Obj = bNext
		p.triples = append(p.triples, trip)
		trip.Sub = bNext
	}
	trip.Pred = IRI{name: "http://www.w3.org/1999/02/22-rdf-syntax-ns#rest"}
	trip.Obj = IRI{name: "http://www.w3.org/1999/02/22-rdf-syntax-ns#nil"}
	p.triples = append(p.triples, trip)
	length++
	return
}

// parseBlankNodePropertyList parses a blankNodePropertyList ('[' predicateObjectList ']')
func (p *parser) parseBlankNodePropertyList(pos int) (blank BlankNode, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("reached eof before end of blank node property list")
		return
	}
	if p.runes[pos] != '[' {
		err = errors.New("no BlankNodePropertyList; missing [")
		return
	}
	length = p.consumeWS(pos + 1)
	length++
	blank = p.blankNode()
	var trip Triple
	trip.Sub = blank

	var tempLength int
	var poList []predObjList
	poList, tempLength, err = p.parsePredicateObjectList(pos + length)
	if err != nil {
		return
	}
	length += tempLength
	length += p.consumeWS(pos + length)

	for i := range poList {
		trip.Pred = poList[i].pred
		for j := range poList[i].obj {
			trip.Obj = poList[i].obj[j]
			p.triples = append(p.triples, trip)
		}
	}
	if p.runes[pos+length] != ']' {
		err = errors.New("no BlankNodePropertyList; missing ]")
		return
	}
	length++

	return
}

// blankNode creates a new blank node and increments the counter
func (p *parser) blankNode() (blank BlankNode) {
	blank = BlankNode{name: "bn" + strconv.Itoa(p.bnCounter)}
	p.bnCounter++
	return
}

// isEqual checks if runes at position equal specified string
func (p *parser) isEqual(pos int, comp string) (ok bool) {
	ok = false
	compRune, err := toRunes([]byte(comp))
	if err != nil {
		return
	}
	if len(p.runes) < pos+len(compRune) {
		return
	}
	for i := range compRune {
		if compRune[i] != p.runes[pos+i] {
			return
		}
	}
	ok = true
	return
}

// toRunes transforms a byte slice to a rune slice
func toRunes(in []byte) (out []rune, err error) {
	pos := 0
	for pos < len(in) {
		r, s := utf8.DecodeRune(in[pos:])
		if r == utf8.RuneError {
			err = errors.New("Rune error")
		}
		out = append(out, r)
		pos += s
	}
	return
}

// consumeWS returns number of consecutive white spaces
func (p *parser) consumeWS(pos int) (num int) {
	num = 0
	for {
		if len(p.runes) <= pos {
			break
		}
		if p.runes[pos] == ' ' {
			num++
			pos++
		} else {
			break
		}
	}
	return
}

// parseUntil returns a string from current position to next occurance of specified rune
func (p *parser) parseUntil(pos int, delim rune) (res string, length int, err error) {
	length = 0
	var r []rune
	for {
		if len(p.runes) <= pos+length {
			err = errors.New("reached eof before delimiter")
			return
		}
		if p.runes[pos+length] == delim && p.runes[pos+length-1] != '\\' {
			break
		} else {
			r = append(r, p.runes[pos+length])
			length++
		}
	}
	res = string(r)
	return
}
