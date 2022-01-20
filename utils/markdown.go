package utils

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Table struct {
	header []string
	data   [][]string
	name   string
}

func NewTable(header []string, data [][]string, name string) Table {
	return Table{
		header: header,
		data:   data,
		name:   name,
	}
}

func (T *Table) String() string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader(T.header)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.AppendBulk(T.data)
	table.Render()

	return tableString.String()
}

func (T *Table) WriteToFile(fd *os.File) (err error) {
	_, err = fd.WriteString(T.String())
	_, err = fd.WriteString("\n")
	return
}
