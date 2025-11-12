package utilExcel

import (
	"fmt"
)

type ExcelSheet struct {
	sheetIndex int
	sheetName  string
	file       *ExcelFile
}

func (f *ExcelFile) SheetCount() (v int) {
	v = f.excelizeFile().SheetCount
	return
}
func (f *ExcelFile) GetSheetByIndex(index int) (sheet *ExcelSheet, err error) {
	if index >= f.SheetCount() {
		err = fmt.Errorf("sheet index %d 不存在", index)
		return
	}

	sheet = &ExcelSheet{
		sheetIndex: index,
		sheetName:  f.excelizeFile().GetSheetName(index),
		file:       f,
	}
	return
}
func (f *ExcelFile) GetSheetByName(name string) (sheet *ExcelSheet, err error) {
	index, err := f.excelizeFile().GetSheetIndex(name)
	if nil != err {
		err = fmt.Errorf("sheet name %s 不存在", name)
		return
	}

	sheet = &ExcelSheet{
		sheetIndex: index,
		sheetName:  name,
		file:       f,
	}
	return
}
func (f *ExcelFile) ForeachSheet(callback func(sheet *ExcelSheet) error) (err error) {
	for _, s := range f.excelizeFile().GetSheetMap() {
		sheet, err1 := f.GetSheetByName(s)
		if nil != err1 {
			err = err1
			return
		}
		err1 = callback(sheet)
		sheet = nil
		if nil != err1 {
			err = err1
			break
		}
	}
	return
}

func (s *ExcelSheet) Index() int {
	return s.sheetIndex
}
func (s *ExcelSheet) Name() string {
	return s.sheetName
}
func (s *ExcelSheet) File() *ExcelFile {
	return s.file
}

func (s *ExcelSheet) GetCellValue(rowNumber int, columnTitle string) (v string, err error) {
	return s.File().excelizeFile().GetCellValue(s.sheetName, FormatCellId(rowNumber, columnTitle))
}
func (s *ExcelSheet) GetCellHyperLink(rowNumber int, columnTitle string) (hasLink bool, link string, err error) {
	return s.File().excelizeFile().GetCellHyperLink(s.sheetName, FormatCellId(rowNumber, columnTitle))
}

func (s *ExcelSheet) GetCellFormula(rowNumber int, columnTitle string) (v string, err error) {
	return s.File().excelizeFile().GetCellFormula(s.sheetName, FormatCellId(rowNumber, columnTitle))
}
func (s *ExcelSheet) GetCellRichText(rowNumber int, columnTitle string) (v []*ExcelRichText, err error) {
	rts, err := s.File().excelizeFile().GetCellRichText(s.sheetName, FormatCellId(rowNumber, columnTitle))
	if nil != err {
		return
	}
	v = excelRichTextsFromExcelize(rts)
	return
}
func (s *ExcelSheet) GetCellType(rowNumber int, columnTitle string) (v ExcelCellType, err error) {
	t, err := s.File().excelizeFile().GetCellType(s.sheetName, FormatCellId(rowNumber, columnTitle))
	if nil != err {
		return
	}
	v = ExcelCellType(t)
	return
}

func (s *ExcelSheet) GetCellStyle(rowNumber int, columnTitle string) (styleId int, err error) {
	return s.File().excelizeFile().GetCellStyle(s.sheetName, FormatCellId(rowNumber, columnTitle))
}

func (s *ExcelSheet) SetCellValue(rowNumber int, columnTitle string, value interface{}) (err error) {
	return s.File().excelizeFile().SetCellValue(s.sheetName, FormatCellId(rowNumber, columnTitle), value)
}
func (s *ExcelSheet) SetCellHyperLink(rowNumber int, columnTitle string, link string, linkType string) (err error) {
	return s.File().excelizeFile().SetCellHyperLink(s.sheetName, FormatCellId(rowNumber, columnTitle), link, linkType)
}
func (s *ExcelSheet) SetCellFormula(rowNumber int, columnTitle string, v string) (err error) {
	return s.File().excelizeFile().SetCellFormula(s.sheetName, FormatCellId(rowNumber, columnTitle), v)
}
func (s *ExcelSheet) SetCellRichText(rowNumber int, columnTitle string, v []*ExcelRichText) (err error) {
	return s.File().excelizeFile().SetCellRichText(s.sheetName, FormatCellId(rowNumber, columnTitle), excelRichTextsToExcelize(v))
}
func (s *ExcelSheet) SetCellStyle(rowNumber int, columnTitle string, styleId int) (err error) {
	return s.File().excelizeFile().SetCellStyle(s.sheetName, FormatCellId(rowNumber, columnTitle), FormatCellId(rowNumber, columnTitle), styleId)
}
func (s *ExcelSheet) SetCellBool(rowNumber int, columnTitle string, v bool) (err error) {
	return s.File().excelizeFile().SetCellBool(s.sheetName, FormatCellId(rowNumber, columnTitle), v)
}
func (s *ExcelSheet) SetCellFloat(rowNumber int, columnTitle string, v float64, precision int) (err error) {
	return s.File().excelizeFile().SetCellFloat(s.sheetName, FormatCellId(rowNumber, columnTitle), v, precision, 32)
}
func (s *ExcelSheet) SetCellInt(rowNumber int, columnTitle string, v int64) (err error) {
	return s.File().excelizeFile().SetCellInt(s.sheetName, FormatCellId(rowNumber, columnTitle), v)
}
func (s *ExcelSheet) SetCellUint(rowNumber int, columnTitle string, v uint64) (err error) {
	return s.File().excelizeFile().SetCellUint(s.sheetName, FormatCellId(rowNumber, columnTitle), v)
}
func (s *ExcelSheet) SetCellStr(rowNumber int, columnTitle string, v string) (err error) {
	return s.File().excelizeFile().SetCellStr(s.sheetName, FormatCellId(rowNumber, columnTitle), v)
}

func (s *ExcelSheet) SetCellDefault(rowNumber int, columnTitle string, v string) (err error) {
	return s.File().excelizeFile().SetCellDefault(s.sheetName, FormatCellId(rowNumber, columnTitle), v)
}

func (s *ExcelSheet) GetPictureCells() (cells []string, err error) {
	return s.File().excelizeFile().GetPictureCells(s.sheetName)
}
func (s *ExcelSheet) GetPictures(rowNumber int, columnTitle string) (pictures []*ExcelPicture, err error) {
	pics, err := s.File().excelizeFile().GetPictures(s.sheetName, FormatCellId(rowNumber, columnTitle))
	if nil != err {
		return
	}
	pictures = excelPicturesFromExcelize(pics)
	return
}
func (s *ExcelSheet) AddPicture(rowNumber int, columnTitle string, name string, opts *ExcelGraphicOptions) (err error) {
	return s.File().excelizeFile().AddPicture(s.sheetName, FormatCellId(rowNumber, columnTitle), name, opts.ToExcelize())
}
func (s *ExcelSheet) AddPictureFromBytes(rowNumber int, columnTitle string, pic *ExcelPicture) (err error) {
	return s.File().excelizeFile().AddPictureFromBytes(s.sheetName, FormatCellId(rowNumber, columnTitle), pic.ToExcelize())
}
func (s *ExcelSheet) DeletePicture(rowNumber int, columnTitle string) (err error) {
	return s.File().excelizeFile().DeletePicture(s.sheetName, FormatCellId(rowNumber, columnTitle))
}
