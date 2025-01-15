package utilExcel

import "github.com/xuri/excelize/v2"

type ExcelRow struct {
	rowNumber int
	sheet     *ExcelSheet
	rows      *excelize.Rows
}

func (s *ExcelSheet) GetRow(rowNumber int) (row *ExcelRow) {
	row = &ExcelRow{rowNumber: rowNumber, sheet: s, rows: nil}
	return
}
func (s *ExcelSheet) getRows() (rows *excelize.Rows, err error) {
	return s.File().excelizeFile().Rows(s.sheetName)
}
func (s *ExcelSheet) ForeachRow(callback func(row *ExcelRow) error) (err error) {
	rows, err := s.getRows()
	if err != nil {
		return
	}
	defer func() {
		_ = rows.Close()
	}()

	cur := 0
	for rows.Next() {
		cur++
		row := &ExcelRow{rowNumber: cur, sheet: s, rows: rows}
		err = callback(row)
		row = nil
		if nil != err {
			break
		}
	}
	err = ErrorIfBreakNil(err)

	return
}

func (s *ExcelSheet) RowCount() (v int) {
	_ = s.ForeachRow(func(row *ExcelRow) error {
		v++
		return nil
	})
	return
}

func (r *ExcelRow) RowNumber() int {
	return r.rowNumber
}
func (r *ExcelRow) Sheet() *ExcelSheet {
	return r.sheet
}

func (r *ExcelRow) BuildCellId(columnTitle string) string {
	return FormatCellId(r.RowNumber(), columnTitle)
}

func (r *ExcelRow) GetCellValue(columnTitle string) (v string, err error) {
	return r.Sheet().GetCellValue(r.RowNumber(), columnTitle)
}
func (r *ExcelRow) GetCellHyperLink(columnTitle string) (hasLink bool, link string, err error) {
	return r.Sheet().GetCellHyperLink(r.RowNumber(), columnTitle)
}
func (r *ExcelRow) GetCellFormula(columnTitle string) (v string, err error) {
	return r.Sheet().GetCellFormula(r.RowNumber(), columnTitle)
}
func (r *ExcelRow) GetCellRichText(columnTitle string) (v []*ExcelRichText, err error) {
	return r.Sheet().GetCellRichText(r.RowNumber(), columnTitle)
}
func (r *ExcelRow) GetCellType(columnTitle string) (v ExcelCellType, err error) {
	return r.Sheet().GetCellType(r.RowNumber(), columnTitle)
}
func (r *ExcelRow) GetCellStyle(columnTitle string) (v int, err error) {
	return r.Sheet().GetCellStyle(r.RowNumber(), columnTitle)
}

func (r *ExcelRow) SetCellValue(columnTitle string, value interface{}) (err error) {
	return r.Sheet().SetCellValue(r.RowNumber(), columnTitle, value)
}
func (r *ExcelRow) SetCellHyperLink(columnTitle string, link string, linkType string) (err error) {
	return r.Sheet().SetCellHyperLink(r.RowNumber(), columnTitle, link, linkType)
}
func (r *ExcelRow) SetCellFormula(columnTitle string, v string) (err error) {
	return r.Sheet().SetCellFormula(r.RowNumber(), columnTitle, v)
}
func (r *ExcelRow) SetCellRichText(columnTitle string, v []*ExcelRichText) (err error) {
	return r.Sheet().SetCellRichText(r.RowNumber(), columnTitle, v)
}
func (r *ExcelRow) SetCellStyle(columnTitle string, styleId int) (err error) {
	return r.Sheet().SetCellStyle(r.RowNumber(), columnTitle, styleId)
}
func (r *ExcelRow) SetCellBool(columnTitle string, v bool) (err error) {
	return r.Sheet().SetCellBool(r.RowNumber(), columnTitle, v)
}
func (r *ExcelRow) SetCellFloat(columnTitle string, v float64, precision int) (err error) {
	return r.Sheet().SetCellFloat(r.RowNumber(), columnTitle, v, precision)
}
func (r *ExcelRow) SetCellInt(columnTitle string, v int) (err error) {
	return r.Sheet().SetCellInt(r.RowNumber(), columnTitle, v)
}
func (r *ExcelRow) SetCellUint(columnTitle string, v uint64) (err error) {
	return r.Sheet().SetCellUint(r.RowNumber(), columnTitle, v)
}
func (r *ExcelRow) SetCellStr(columnTitle string, v string) (err error) {
	return r.Sheet().SetCellStr(r.RowNumber(), columnTitle, v)
}

func (r *ExcelRow) SetCellDefault(columnTitle string, v string) (err error) {
	return r.Sheet().SetCellDefault(r.RowNumber(), columnTitle, v)
}

func (r *ExcelRow) GetPictures(columnTitle string) (pictures []*ExcelPicture, err error) {
	return r.Sheet().GetPictures(r.RowNumber(), columnTitle)
}
func (r *ExcelRow) AddPicture(columnTitle string, name string, opts *ExcelGraphicOptions) (err error) {
	return r.Sheet().AddPicture(r.RowNumber(), columnTitle, name, opts)
}
func (r *ExcelRow) AddPictureFromBytes(columnTitle string, pic *ExcelPicture) (err error) {
	return r.Sheet().AddPictureFromBytes(r.RowNumber(), columnTitle, pic)
}
func (r *ExcelRow) DeletePicture(columnTitle string) (err error) {
	return r.Sheet().DeletePicture(r.RowNumber(), columnTitle)
}
