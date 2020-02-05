package rdf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// IRIs of different literal datatype
var (
	// XsdString string
	XsdString = "http://www.w3.org/2001/XMLSchema#string"
	// XsdBoolean bool
	XsdBoolean = "http://www.w3.org/2001/XMLSchema#boolean"
	// XsdDecimal float
	XsdDecimal = "http://www.w3.org/2001/XMLSchema#decimal"
	// XsdInteger int
	XsdInteger = "http://www.w3.org/2001/XMLSchema#integer"
	// XsdDouble float
	XsdDouble = "http://www.w3.org/2001/XMLSchema#double"

	// XsdDate time
	XsdDate = "http://www.w3.org/2001/XMLSchema#date"
	// XsdDateTime time
	XsdDateTime = "http://www.w3.org/2001/XMLSchema#dateTime"
	// XsdDateTimeStamp time
	XsdDateTimeStamp = "http://www.w3.org/2001/XMLSchema#dateTimeStamp"
	// XsdYear time
	XsdYear = "http://www.w3.org/2001/XMLSchema#gYear"
	// XsdMonth time
	XsdMonth = "http://www.w3.org/2001/XMLSchema#gMonth"
	// XsdDay time
	XsdDay = "http://www.w3.org/2001/XMLSchema#gDay"
	// XsdYearMonth time
	XsdYearMonth = "http://www.w3.org/2001/XMLSchema#gYearMonth"
	// XsdDuration time
	XsdDuration = "http://www.w3.org/2001/XMLSchema#Duration"

	// XsdByte byte
	XsdByte = "http://www.w3.org/2001/XMLSchema#byte"
)

// Triple is one rdf triple consisting of Subject, Predicate and Object
type Triple struct {
	Sub  Subject
	Pred Predicate
	Obj  Object
}

// TermType is the type of a RDF term (IRI, Literal, BlankNode)
type TermType int

// possible RDF term types
const (
	TermBlankNode TermType = iota
	TermLiteral
	TermIRI
)

// Term represents an RDF term. possible types: Blank node, Literal and IRI.
type Term interface {
	Type() (typ TermType)
	String() (str string)
	SerializeTTL(prefix map[string]string) (str string)
}

// Subject represents the subject of a rdf triple
type Subject interface {
	Term
}

// Predicate represents the predicate of a rdf triple
type Predicate interface {
	Term
}

// Object represents the object of a rdf triple
type Object interface {
	Term
}

// IRI is a possible RDF term
type IRI struct {
	name string
}

// Type denotes the term type
func (iri IRI) Type() (typ TermType) {
	typ = TermIRI
	return
}

// String prints the IRI
func (iri IRI) String() (str string) {
	str = iri.name
	return
}

// NewIRI returns an IRI object with specified name
func NewIRI(name string) (iri IRI) {
	iri = IRI{name: name}
	return
}

// Literal is a possible RDF term
type Literal struct {
	str     string
	typeIRI string
	langTag string
	value   interface{}
}

// Type denotes the term type
func (lit Literal) Type() (typ TermType) {
	typ = TermLiteral
	return
}

// String prints the literal string
func (lit Literal) String() (str string) {
	str = lit.str
	return
}

// NewLiteral returns a literal
func NewLiteral(val interface{}, typ string) (lit Literal, err error) {
	switch t := val.(type) {
	case int, int32, int64:
		lit = Literal{str: fmt.Sprintf("%v", t), typeIRI: XsdInteger, value: t}
	case bool:
		lit = Literal{str: fmt.Sprintf("%v", t), typeIRI: XsdBoolean, value: t}
	case float32, float64:
		lit = Literal{str: fmt.Sprintf("%v", t), typeIRI: XsdDouble, value: t}
	case string:
		lit = Literal{str: t, typeIRI: XsdString, value: t}
	case time.Time:
		switch typ {
		case XsdDateTime:
			lit = Literal{str: t.Format(time.RFC3339), typeIRI: XsdDateTime, value: t}
		case XsdDateTimeStamp:
			lit = Literal{str: t.Format(time.RFC3339), typeIRI: XsdDateTimeStamp, value: t}
		case XsdDate:
			lit = Literal{str: t.Format("2006-01-02Z07:00"), typeIRI: XsdDate, value: t}
		case XsdDay:
			lit = Literal{str: t.Format("---02"), typeIRI: XsdDay, value: t}
		case XsdMonth:
			lit = Literal{str: t.Format("--01"), typeIRI: XsdMonth, value: t}
		case XsdYear:
			lit = Literal{str: t.Format("2006"), typeIRI: XsdYear, value: t}
		case XsdYearMonth:
			lit = Literal{str: t.Format("2006-01"), typeIRI: XsdYearMonth, value: t}
		default:
			lit = Literal{str: t.Format(time.RFC3339), typeIRI: XsdDateTime, value: t}
		}
	case time.Duration:
		sec := fmt.Sprintf("%v", float32(t.Seconds())-float32(int(t.Seconds()))+float32(int(t.Seconds())%60)) + "S"
		min := fmt.Sprintf("%v", int(t.Minutes())%60) + "M"
		hour := fmt.Sprintf("%v", int(t.Hours())) + "H"
		//day := fmt.Sprintf("%v", int(in.Hours())/24) + "D"
		lit = Literal{str: "P" + hour + min + sec, typeIRI: XsdDuration, value: t}
	case []byte:
		lit = Literal{str: string(t), typeIRI: XsdByte, value: t}
	default:
		err = fmt.Errorf("invalid rdf literal type %v", t)
	}
	return
}

// ToTime converts a xsd literal to time.Time if possible
func (lit Literal) ToTime() (t time.Time, err error) {
	switch lit.typeIRI {
	case XsdDateTime:
		t, err = time.Parse(time.RFC3339, lit.str)
	case XsdDateTimeStamp:
		t, err = time.Parse(time.RFC3339, lit.str)
	case XsdDate:
		t, err = time.Parse("2006-01-02Z07:00", lit.str)
	case XsdDay:
		t, err = time.Parse("---02", lit.str)
	case XsdMonth:
		t, err = time.Parse("--01", lit.str)
	case XsdYear:
		t, err = time.Parse("2006", lit.str)
	case XsdYearMonth:
		t, err = time.Parse("2006-01", lit.str)
	default:
		err = errors.New("Cannot convert xsd datatype " + lit.typeIRI + " to time")
	}
	if err != nil {
		return
	}
	lit.value = t
	return
}

// ToDuration converts a literal to time.Duration if possible
func (lit Literal) ToDuration() (d time.Duration, err error) {
	switch lit.typeIRI {
	case XsdDuration:
		temp := strings.Split(lit.str, "P")
		if len(temp) < 2 {
			err = errors.New("Cannot convert xsd datatype " + lit.typeIRI + " to duration")
			return
		}
		str := ""
		h := strings.Split(temp[1], "H")
		if len(h) > 1 {
			_, err := strconv.Atoi(h[0])
			if err != nil {
				str += h[0] + "h"
			}
		}
		m := strings.Split(h[len(h)-1], "M")
		if len(m) > 1 {
			_, err := strconv.Atoi(m[0])
			if err != nil {
				str += m[0] + "m"
			}
		}
		s := strings.Split(m[len(m)-1], "S")
		if len(s) > 1 {
			_, err := strconv.ParseFloat(s[0], 32)
			if err != nil {
				str += s[0] + "s"
			}
		}
		d, err = time.ParseDuration(str)
	default:
		err = errors.New("Cannot convert xsd datatype " + lit.typeIRI + " to duration")
	}
	if err != nil {
		return
	}
	lit.value = d
	return
}

// BlankNode is a possible RDF term
type BlankNode struct {
	name string
}

// Type denotes the term type
func (blank BlankNode) Type() (typ TermType) {
	typ = TermBlankNode
	return
}

// String prints the name
func (blank BlankNode) String() (str string) {
	str = blank.name
	return
}
