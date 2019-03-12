package doc

//IDocumented is anything that is described in documentation
type IDocumented interface {
	Doc() IDoc
}

//IDoc ...
type IDoc interface {
	//Title(text string)
	//Par(text string)
	//Table(caption string) ITable
}

//ITable ...
type ITable interface {
	Row() IRow
}

//IRow ...
type IRow interface {
	Col() IDoc
}

//New ...
func New() IDoc {
	return doc{}
}

type doc struct {
}

type text struct {
	value string
}
