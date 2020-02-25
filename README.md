# OWL2Go

OWL2Go is a code generator that takes a OWL ontology and generates a Go package which implementes a data model compliant to the input ontology and serialization/deserialization from/to Turtle/JSON-LD documents.

## Mapping

OWL construct | OWL expression (owl:..., rdf(s):...) | Go construct
--- | ---| ---
Class | Class A | ```interface A```; ```struct sA``` implements ```A```
 | | A subClassOf B | ```interface A``` extends ```interface B```
 | | A intersectionOf B and C | ```interface A``` extends ```interfaces B``` and ```C```
 | | A unionOf B and C | ```interface B``` and ```C``` extend ```interface A```
Restriction (of class A) on Property B | allValuesFrom C,D; someValuesFrom; hasValue; xCardinality | ```C```,```D```: allowed types; ```E``` := base type of ```C```, ```D```; ```F``` := common type of ```C```, ```D```; ```A = {...; B() []E; SetB([]E)error; AddB(...E)error; DelB(...E)}```; ```sA = {...; b map[string]F}```
Property A | domain B | add property ```A``` to class ```B``` as for restrictions
 | | range B | add ```B``` to allowed types of property
 | | inverseOf B | add function call to ```AddB()``` and ```DelB()``` in ```AddA()``` and ```DelA()``` functions and vice versa
 | | SymmetricProperty | add function call to ```AddA()``` and ```DelA()``` in ```AddA()``` and ```DelA()``` functions
Individual | type | create individual when creating model

## How to use OWL2Go

For an example on how to generate a package for the SAREF ontology, please refer to ```cmd/main.go```.

## How to use the generated package

## Copyright

2020, Institute for Automation of Complex Power Systems, EONERC

## License

This project is licensed under either of

- Apache License, Version 2.0 ([LICENSE-Apache](LICENSE-Apache) or [http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0))
- MIT license ([LICENSE-MIT](LICENSE-MIT) or [https://opensource.org/licenses/MIT](https://opensource.org/licenses/MIT))

at your option.

## Contact

[![EONERC ACS Logo](docs/eonerc_logo.png)](http://www.acs.eonerc.rwth-aachen.de)

- [Stefan Dähling](mailto:sdaehling@eonerc.rwth-aachen.de)

[Institute for Automation of Complex Power Systems (ACS)](http://www.acs.eonerc.rwth-aachen.de)  
[EON Energy Research Center (EONERC)](http://www.eonerc.rwth-aachen.de)  
[RWTH Aachen University, Germany](http://www.rwth-aachen.de)
