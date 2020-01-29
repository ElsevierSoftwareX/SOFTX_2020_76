package rdf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode/utf8"
)

// token represents one ttl expression
type token struct {
	typ   int    // token type
	value string // value of token
}

// parser parses a ttl document
type parser struct {
	reader       *bufio.Reader     // reader of turtle document
	runes        []rune            // turtle document as rune slice
	posStatement int               // starting position of current statement
	prefix       map[string]string // prefixes
	base         string            // base IRI
	curSubject   string            // current Subject
	curPredicate string            // current Predicate
	triples      []Triple          // list of all extracted triples
	bnCounter    int               // blank node counter
}

type predObjList struct {
	pred Predicate
	obj  []Object
}

// DecodeTTL decodes a ttl input to rdf triples
func DecodeTTL(input io.Reader) (trip []Triple, err error) {
	p := &parser{reader: bufio.NewReader(input), prefix: make(map[string]string)}
	err = p.parseRunes()
	if err != nil {
		return
	}
	for {
		err = p.parseStatement()
		if err != nil {
			break
		}
		if p.posStatement >= len(p.runes) {
			break
		}
	}
	fmt.Println(p.prefix)
	fmt.Println(p.base)
	fmt.Println(p.triples)

	return
}

// parseRunes parses all runes from the reader and omits empty lines and comments
func (p *parser) parseRunes() (err error) {
	for {
		var line []byte
		line, err = p.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		if len(line) == 0 {
			continue
		}
		pos := 0
		// omit new line at end of slice
		for pos < len(line)-1 {
			var r rune
			var s int
			r, s = utf8.DecodeRune(line[pos:])
			if r == utf8.RuneError {
				err = errors.New("Rune error")
			}
			if pos == 0 && r == '#' {
				break
			}
			p.runes = append(p.runes, r)
			pos += s
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

// parseDirective decodes one directive beginning from the specified position
func (p *parser) parseDirective(pos int) (length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Invalid directive " + strconv.Itoa(pos))
		return
	}
	switch p.runes[pos] {
	case 'p':
		var prefix, iri string
		var tempLength int
		if p.isEqual(pos, "prefix") {
			length = 6
			length += p.consumeWP(pos + length)
			prefix, tempLength, err = p.parsePrefix(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			length += p.consumeWP(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.prefix[prefix] = iri
		} else {
			err = errors.New("Invalid directive " + string(p.runes[pos]) + strconv.Itoa(pos))
			return
		}
	case 'b':
		var iri string
		var tempLength int
		if p.isEqual(pos, "base") {
			length = 4
			length += p.consumeWP(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.base = iri
		} else {
			err = errors.New("Invalid directive " + string(p.runes[pos]) + strconv.Itoa(pos))
			return
		}
	case 'P':
		var prefix, iri string
		var tempLength int
		if p.isEqual(pos, "PREFIX") {
			length = 6
			length += p.consumeWP(pos + length)
			prefix, tempLength, err = p.parsePrefix(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			length += p.consumeWP(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.prefix[prefix] = iri
		} else {
			err = errors.New("Invalid directive " + string(p.runes[pos]) + strconv.Itoa(pos))
			return
		}
	case 'B':
		var iri string
		var tempLength int
		if p.isEqual(pos, "BASE") {
			length = 4
			length += p.consumeWP(pos + length)
			iri, tempLength, err = p.parseIRIRef(pos + length)
			if err != nil {
				return
			}
			length += tempLength
			p.base = iri
		} else {
			err = errors.New("Invalid directive " + string(p.runes[pos]) + strconv.Itoa(pos))
			return
		}
	default:
		err = errors.New("Invalid directive " + string(p.runes[pos]) + strconv.Itoa(pos))
		return
	}
	// consumer dot
	length += p.consumeWP(pos + length)
	if p.isEqual(pos+length, ".") {
		length++
	} else {
		err = errors.New("No dot")
		return
	}
	length += p.consumeWP(pos + length)
	return
}

// parsePrefix parses one prefix beginning from the specified position
func (p *parser) parsePrefix(pos int) (prefix string, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Prefix error " + strconv.Itoa(pos))
	}
	prefix, length, err = p.parseUntil(pos, ':')
	length++
	return
}

// parseTriples parses all tripls in a statement
func (p *parser) parseTriples(pos int) (length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Invalid triples " + strconv.Itoa(pos))
		return
	}
	var trip Triple
	trip.Sub, length, err = p.parseSubject(pos)
	if err != nil {
		return
	}
	length += p.consumeWP(pos + length)

	var tempLength int
	var poList []predObjList
	poList, tempLength, err = p.parsePredicateObjectList(pos + length)
	if err != nil {
		return
	}
	length += tempLength
	length += p.consumeWP(pos + length)

	for i := range poList {
		trip.Pred = poList[i].pred
		for j := range poList[i].obj {
			trip.Obj = poList[i].obj[j]
			p.triples = append(p.triples, trip)
		}
	}

	// consume dot
	length += p.consumeWP(pos + length)
	if p.isEqual(pos+length, ".") {
		length++
	} else {
		err = errors.New("No dot")
		return
	}
	length += p.consumeWP(pos + length)

	return
}

// parseSubject parses the subject of a triple
func (p *parser) parseSubject(pos int) (subj Subject, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("Invalid subject " + strconv.Itoa(pos))
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

// parsePredicateObjectList parses a predicate object list
func (p *parser) parsePredicateObjectList(pos int) (list []predObjList, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Invalid predicate object list " + strconv.Itoa(pos))
		return
	}
	length = 0

	for {
		var poListTemp predObjList
		var tempLength int
		poListTemp.pred, tempLength, err = p.parsePredicate(pos + length)
		if err != nil {
			return
		}
		length += tempLength
		length += p.consumeWP(pos + length)

		poListTemp.obj, tempLength, err = p.parseObjectList(pos + length)
		if err != nil {
			return
		}
		length += tempLength
		length += p.consumeWP(pos + length)
		list = append(list, poListTemp)

		if !p.isEqual(pos+length, ";") {
			break
		}
		length++
		length += p.consumeWP(pos + length)
	}
	return
}

// parsePredicate parses the next predicate
func (p *parser) parsePredicate(pos int) (pred Predicate, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Invalid predicate " + strconv.Itoa(pos))
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

// parseObjectList parses an object list
func (p *parser) parseObjectList(pos int) (obj []Object, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Invalid object list " + strconv.Itoa(pos))
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
		length += p.consumeWP(pos + length)
		if !p.isEqual(pos+length, ",") {
			break
		}
		length++
		length += p.consumeWP(pos + length)
	}
	return
}

// parseObject parses one object
func (p *parser) parseObject(pos int) (obj Object, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Invalid object " + strconv.Itoa(pos))
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

// parseIRI parses the next iri
func (p *parser) parseIRI(pos int) (iri IRI, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("IRI error " + strconv.Itoa(pos))
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

// parseIRI parses IRIRef <iri>
func (p *parser) parseIRIRef(pos int) (iri string, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("IRI error " + strconv.Itoa(pos))
		return
	}
	if p.runes[pos] != '<' {
		err = errors.New("No IRI: " + string(p.runes[pos]) + strconv.Itoa(pos))
		return
	}
	iri, length, err = p.parseUntil(pos+1, '>')
	length += 2
	return
}

// parsePrefixedName parses prefixed name prefix:name
func (p *parser) parsePrefixedName(pos int) (iri string, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("IRI error " + strconv.Itoa(pos))
		return
	}
	var prefix string
	prefix, length, err = p.parsePrefix(pos)
	if err != nil {
		return
	}
	ok := false
	if iri, ok = p.prefix[prefix]; !ok {
		err = errors.New("No such prefix " + prefix)
	}
	var name string
	var tempLength int
	name, tempLength, err = p.parseUntil(pos+length, ' ')
	if err != nil {
		return
	}
	iri = iri + name
	length += tempLength
	return
}

// parseLiteral parses a literal
func (p *parser) parseLiteral(pos int) (lit Literal, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Literal error " + strconv.Itoa(pos))
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

// parseRDFLiteral parses a rdf literal
func (p *parser) parseRDFLiteral(pos int) (lit Literal, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("Literal error " + strconv.Itoa(pos))
		return
	}
	if p.runes[pos] == '"' {
		if p.runes[pos+1] == '"' {
			lit.str, length, err = p.parseUntil(pos+3, '"')
			length += 6
		} else {
			lit.str, length, err = p.parseUntil(pos+1, '"')
			length += 2
		}
	} else if p.runes[pos] != '\'' {
		if p.runes[pos+1] == '\'' {
			lit.str, length, err = p.parseUntil(pos+3, '\'')
			length += 6
		} else {
			lit.str, length, err = p.parseUntil(pos+1, '\'')
			length += 2
		}
	} else {
		err = errors.New("No rdf literal " + strconv.Itoa(pos))
	}
	if err != nil {
		return
	}
	if len(p.runes) <= pos+length {
		return
	}
	if p.runes[pos+length] == '@' {
		length++
		var tempLength int
		lit.langTag, tempLength, err = p.parseUntil(pos+length, ' ')
		if err != nil {
			return
		}
		length += tempLength
	} else if p.runes[pos+length] == '^' {
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
		err = errors.New("Literal error " + strconv.Itoa(pos))
		return
	}
	if p.isEqual(pos, "true") {
		lit = Literal{str: "true", value: true}
		length = 4
	} else if p.isEqual(pos, "false") {
		lit = Literal{str: "false", value: false}
		length = 5
	} else {
		err = errors.New("No boolean literal " + strconv.Itoa(pos))
	}
	return
}

// parseNumericLiteral parses a numeric literal
func (p *parser) parseNumericLiteral(pos int) (lit Literal, length int, err error) {
	if len(p.runes) <= pos {
		err = errors.New("Literal error " + strconv.Itoa(pos))
		return
	}
	// get length of literal
	var value float64
	tempLength := 0
	lit.str, length, err = p.parseUntil(pos, ' ')

	if err != nil {
		return
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
		err = errors.New("BlankNode error " + strconv.Itoa(pos))
		return
	}
	if p.runes[pos] != '_' || p.runes[pos+1] != ':' {
		err = errors.New("No BlankNode " + strconv.Itoa(pos))
		return
	}
	blank.name, length, err = p.parseUntil(pos, ' ')
	if err != nil {
		return
	}
	length += 2
	return
}

// parseCollection parses a collection
func (p *parser) parseCollection(pos int) (blank BlankNode, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("Collection error " + strconv.Itoa(pos))
		return
	}
	blank = p.blankNode()
	if p.runes[pos] != '(' {
		err = errors.New("No Collection " + strconv.Itoa(pos))
		return
	}
	length = 1
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

		tempLength = p.consumeWP(pos + length)
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

// parseBlankNodePropertyList parses a blankNodePropertyList
func (p *parser) parseBlankNodePropertyList(pos int) (blank BlankNode, length int, err error) {
	if len(p.runes) <= pos+1 {
		err = errors.New("BlankNodePropertyList error " + strconv.Itoa(pos))
		return
	}
	if p.runes[pos] != '[' {
		err = errors.New("No BlankNodePropertyList " + strconv.Itoa(pos))
		return
	}
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
	length += p.consumeWP(pos + length)

	for i := range poList {
		trip.Pred = poList[i].pred
		for j := range poList[i].obj {
			trip.Obj = poList[i].obj[j]
			p.triples = append(p.triples, trip)
		}
	}
	if p.runes[pos+length] != ']' {
		err = errors.New("BlankNodePropertyList error " + strconv.Itoa(pos))
		return
	}
	length++

	return
}

// blankNode creates a new blank node and increments the counter
func (p *parser) blankNode() (blank BlankNode) {
	blank = BlankNode{name: "_:bn" + strconv.Itoa(p.bnCounter)}
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

// consumeWP returns number of consecutive white spaces
func (p *parser) consumeWP(pos int) (num int) {
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
			err = errors.New("No delimiter")
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
