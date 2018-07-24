package xlsx

import (
	"github.com/tealeg/xlsx"
)

type RowSetter func(row *xlsx.Row, data interface{})

func MakeFile(fileName string, setter RowSetter, list []interface{}, firstRow ...string) error {
	file := xlsx.NewFile()

	sheet, err := file.AddSheet("link")
	if err != nil {
		panic(err)
	}

	if len(firstRow) > 0 {
		row := sheet.AddRow()
		for _, t := range firstRow {
			row.AddCell().SetString(t)
		}
	}

	for _, t := range list {
		row := sheet.AddRow()
		setter(row, t)
	}

	err = file.Save(fileName)
	return err
}
