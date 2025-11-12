package utilExcel

type ExcelCell struct {
	columnTitle  string
	columnNumber int
	row          *ExcelRow
}

func (s *ExcelSheet) GetCellByColumnTitle(rowNumber int, columnTitle string) (cell *ExcelCell) {
	row := s.GetRow(rowNumber)
	return row.GetCellByColumnTitle(columnTitle)
}
func (s *ExcelSheet) GetCellByColumnNumber(rowNumber int, columnNumber int) (cell *ExcelCell) {
	row := s.GetRow(rowNumber)
	return row.GetCellByColumnNumber(columnNumber)
}
func (s *ExcelSheet) GetCellByColumnIndex(rowNumber int, columnIndex int) (cell *ExcelCell) {
	row := s.GetRow(rowNumber)
	return row.GetCellByColumnIndex(columnIndex)
}

func (r *ExcelRow) GetCellByColumnTitle(columnTitle string) (cell *ExcelCell) {
	cell = &ExcelCell{
		columnTitle:  columnTitle,
		columnNumber: ColumnTitleToNumber(columnTitle),
		row:          r,
	}
	return
}
func (r *ExcelRow) GetCellByColumnNumber(columnNumber int) (cell *ExcelCell) {
	cell = &ExcelCell{
		columnTitle:  ColumnNumberToTitle(columnNumber),
		columnNumber: columnNumber,
		row:          r,
	}
	return
}
func (r *ExcelRow) GetCellByColumnIndex(columnIndex int) (cell *ExcelCell) {
	cell = &ExcelCell{
		columnTitle:  ColumnIndexToTitle(columnIndex),
		columnNumber: ColumnIndexToNumber(columnIndex),
		row:          r,
	}
	return
}

func (r *ExcelRow) ForeachCell(callback func(cell *ExcelCell) error) (err error) {
	err = r.Sheet().ForeachRow(func(row *ExcelRow) error {
		if row.rowNumber == r.RowNumber() {
			if nil != r.rows {
				columns, err1 := r.rows.Columns()
				if nil != err1 {
					return err1
				}
				for _, columnName := range columns {
					cell := &ExcelCell{
						columnTitle:  columnName,
						columnNumber: ColumnTitleToNumber(columnName),
						row:          r,
					}
					err1 = callback(cell)
					cell = nil
					if nil != err1 {
						return err1
					}
				}
			}
			return ErrorBreak()
		}

		return nil
	})
	err = ErrorIfBreakNil(err)
	return
}
func (r *ExcelRow) CellCount() (v int) {
	_ = r.ForeachCell(func(cell *ExcelCell) error {
		v++
		return nil
	})
	return
}
func (c *ExcelCell) ColumnTitle() string {
	return c.columnTitle
}
func (c *ExcelCell) ColumnNumber() int {
	return c.columnNumber
}
func (c *ExcelCell) Row() *ExcelRow {
	return c.row
}

func (c *ExcelCell) BuildId() string {
	return c.Row().BuildCellId(c.ColumnTitle())
}
func (c *ExcelCell) GetValue() (v string, err error) {
	v, err = c.Row().GetCellValue(c.ColumnTitle())
	return
}
func (c *ExcelCell) GetHyperLink() (hasLink bool, link string, err error) {
	return c.Row().GetCellHyperLink(c.ColumnTitle())
}
func (c *ExcelCell) GetFormula() (v string, err error) {
	return c.Row().GetCellFormula(c.ColumnTitle())
}
func (c *ExcelCell) GetRichText() (v []*ExcelRichText, err error) {
	return c.Row().GetCellRichText(c.ColumnTitle())
}
func (c *ExcelCell) GetType() (v ExcelCellType, err error) {
	return c.Row().GetCellType(c.ColumnTitle())
}
func (c *ExcelCell) GetStyle() (v int, err error) {
	return c.Row().GetCellStyle(c.ColumnTitle())
}

func (c *ExcelCell) SetValue(value interface{}) (err error) {
	return c.Row().SetCellValue(c.ColumnTitle(), value)
}
func (c *ExcelCell) SetHyperLink(link string, linkType string) (err error) {
	return c.Row().SetCellHyperLink(c.ColumnTitle(), link, linkType)
}
func (c *ExcelCell) SetFormula(v string) (err error) {
	return c.Row().SetCellFormula(c.ColumnTitle(), v)
}
func (c *ExcelCell) SetRichText(v []*ExcelRichText) (err error) {
	return c.Row().SetCellRichText(c.ColumnTitle(), v)
}
func (c *ExcelCell) SetStyle(styleId int) (err error) {
	return c.Row().SetCellStyle(c.ColumnTitle(), styleId)
}
func (c *ExcelCell) SetBool(v bool) (err error) {
	return c.Row().SetCellBool(c.ColumnTitle(), v)
}
func (c *ExcelCell) SetFloat(v float64, precision int) (err error) {
	return c.Row().SetCellFloat(c.ColumnTitle(), v, precision)
}
func (c *ExcelCell) SetInt(v int64) (err error) {
	return c.Row().SetCellInt(c.ColumnTitle(), v)
}
func (c *ExcelCell) SetUint(v uint64) (err error) {
	return c.Row().SetCellUint(c.ColumnTitle(), v)
}
func (c *ExcelCell) SetStr(v string) (err error) {
	return c.Row().SetCellStr(c.ColumnTitle(), v)
}

func (c *ExcelCell) SetDefault(v string) (err error) {
	return c.Row().SetCellDefault(c.ColumnTitle(), v)
}

func (c *ExcelCell) GetPictures() (pictures []*ExcelPicture, err error) {
	return c.Row().GetPictures(c.ColumnTitle())
}
func (c *ExcelCell) AddPicture(name string, opts *ExcelGraphicOptions) (err error) {
	return c.Row().AddPicture(c.ColumnTitle(), name, opts)
}
func (c *ExcelCell) AddPictureFromBytes(pic *ExcelPicture) (err error) {
	return c.Row().AddPictureFromBytes(c.ColumnTitle(), pic)
}
func (c *ExcelCell) DeletePicture() (err error) {
	return c.Row().DeletePicture(c.ColumnTitle())
}
