package utilExcel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type ExcelFile struct {
	ef *excelize.File
}

func OpenFile(filename string, opts ...ExcelFileOptions) (excelFile *ExcelFile, err error) {

	f, err := excelize.OpenFile(filename, optionsToExcelize(opts)...)
	if nil != err {
		return
	}
	excelFile = &ExcelFile{ef: f}
	return
}
func OpenReader(r io.Reader, opts ...ExcelFileOptions) (excelFile *ExcelFile, err error) {

	f, err := excelize.OpenReader(r, optionsToExcelize(opts)...)
	if nil != err {
		return
	}
	excelFile = &ExcelFile{ef: f}
	return
}

func NewFile(opts ...ExcelFileOptions) (excelFile *ExcelFile) {
	f := excelize.NewFile(optionsToExcelize(opts)...)
	excelFile = &ExcelFile{ef: f}
	return
}

func (f *ExcelFile) excelizeFile() *excelize.File {
	return f.ef
}
func (f *ExcelFile) NewStyle(style *ExcelStyle) (styleId int, err error) {
	if nil == style {
		return
	}
	styleId, err = f.ef.NewStyle(style.ToExcelize())
	return
}

func (f *ExcelFile) GetStyle(styleId int) (style *ExcelStyle, err error) {
	v, err := f.ef.GetStyle(styleId)
	if nil != err {
		return
	}
	style = excelStyleFromExcelize(v)
	return
}
func (f *ExcelFile) SaveAs(name string, opts ...ExcelFileOptions) (err error) {
	err = f.ef.SaveAs(name, optionsToExcelize(opts)...)
	return
}
