package template

// Individual template
var Individual = "package ###pkgName###\n\n" +
	"// CreateIndividuals adds all individuals to a model\n" +
	"func (mod *Model) CreateIndividuals() {\n" +
	"###createIndividuals###" +
	"\treturn\n" +
	"}\n\n"

// CreateIndividual template
var CreateIndividual = "\tmod.New###individualType###(\"###individualIRI###\")\n"
