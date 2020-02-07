package pkg

import (
	"os"
	"testing"

	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/codegen"
	"git-ce.rwth-aachen.de/acs/private/research/agent/owl2go.git/pkg/owl"
)

func TestOwl2go(t *testing.T) {
	var err error
	var file *os.File
	file, err = os.Open("../test/test.ttl")
	var on owl.Ontology
	on, err = owl.ExtractOntology(file)
	if err != nil {
		t.Error(err)
	}
	file.Close()

	var mod []owl.GoModel
	mod, err = owl.MapModel(&on, "git-ce.rwth-aachen.de/acs/private/research/agent/saref.git")
	if err != nil {
		t.Error(err)
	}

	err = codegen.GenerateGoCode(mod, "../../saref")
	if err != nil {
		t.Error(err)
	}
}
