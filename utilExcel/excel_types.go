package utilExcel

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
)

const (
	CultureNameUnknown         CultureName = CultureName(excelize.CultureNameUnknown)
	CultureNameEnUSCultureName             = CultureName(excelize.CultureNameEnUS)
	CultureNameZhCN            CultureName = CultureName(excelize.CultureNameZhCN)
)
const (
	LinkTypeExternal = "External"
	LinkTypeLocation = "Location"
)

var errorBreak = fmt.Errorf("utilExcel break")

func ErrorBreak() error {
	return errorBreak
}
func ErrorIsBreak(err error) bool {
	return reflect.DeepEqual(err, errorBreak)
}
func ErrorIfBreakNil(err error) error {
	if ErrorIsBreak(err) {
		err = nil
	}
	return err
}

type CultureName excelize.CultureName

type ExcelFileOptions struct {
	MaxCalcIterations uint
	Password          string
	RawCellValue      bool
	UnzipSizeLimit    int64
	UnzipXMLSizeLimit int64
	ShortDatePattern  string
	LongDatePattern   string
	LongTimePattern   string
	CultureInfo       CultureName
}

func (eo *ExcelFileOptions) ToExcelize() excelize.Options {
	return excelize.Options{
		MaxCalcIterations: eo.MaxCalcIterations,
		Password:          eo.Password,
		RawCellValue:      eo.RawCellValue,
		UnzipSizeLimit:    eo.UnzipSizeLimit,
		UnzipXMLSizeLimit: eo.UnzipXMLSizeLimit,
		ShortDatePattern:  eo.ShortDatePattern,
		LongDatePattern:   eo.LongDatePattern,
		LongTimePattern:   eo.LongTimePattern,
		CultureInfo:       excelize.CultureName(eo.CultureInfo),
	}
}

func optionsToExcelize(opts []ExcelFileOptions) (excelizeOptions []excelize.Options) {
	for _, opt := range opts {
		excelizeOptions = append(excelizeOptions, opt.ToExcelize())
	}
	return
}

type ExcelCellType byte

// Cell value types enumeration.
const (
	CellTypeUnset        = ExcelCellType(excelize.CellTypeUnset)
	CellTypeBool         = ExcelCellType(excelize.CellTypeBool)
	CellTypeDate         = ExcelCellType(excelize.CellTypeDate)
	CellTypeError        = ExcelCellType(excelize.CellTypeError)
	CellTypeFormula      = ExcelCellType(excelize.CellTypeFormula)
	CellTypeInlineString = ExcelCellType(excelize.CellTypeInlineString)
	CellTypeNumber       = ExcelCellType(excelize.CellTypeNumber)
	CellTypeSharedString = ExcelCellType(excelize.CellTypeSharedString)
)

type ExcelGraphicOptions struct {
	AltText             string
	PrintObject         bool
	Locked              bool
	LockAspectRatio     bool
	AutoFit             bool
	AutoFitIgnoreAspect bool
	OffsetX             int
	OffsetY             int
	ScaleX              float64
	ScaleY              float64
	Hyperlink           string
	HyperlinkType       string
	Positioning         string
}

func (ego *ExcelGraphicOptions) ToExcelize() *excelize.GraphicOptions {
	return &excelize.GraphicOptions{
		AltText:             ego.AltText,
		PrintObject:         &ego.PrintObject,
		Locked:              &ego.Locked,
		LockAspectRatio:     ego.LockAspectRatio,
		AutoFit:             ego.AutoFit,
		AutoFitIgnoreAspect: ego.AutoFitIgnoreAspect,
		OffsetX:             ego.OffsetX,
		OffsetY:             ego.OffsetY,
		ScaleX:              ego.ScaleX,
		ScaleY:              ego.ScaleY,
		Hyperlink:           ego.Hyperlink,
		HyperlinkType:       ego.HyperlinkType,
		Positioning:         ego.Positioning,
	}
}
func excelGraphicOptionsFromExcelize(opt *excelize.GraphicOptions) *ExcelGraphicOptions {
	return &ExcelGraphicOptions{
		AltText:             opt.AltText,
		PrintObject:         *opt.PrintObject,
		Locked:              *opt.Locked,
		LockAspectRatio:     opt.LockAspectRatio,
		AutoFit:             opt.AutoFit,
		AutoFitIgnoreAspect: opt.AutoFitIgnoreAspect,
		OffsetX:             opt.OffsetX,
		OffsetY:             opt.OffsetY,
		ScaleX:              opt.ScaleX,
		ScaleY:              opt.ScaleY,
		Hyperlink:           opt.Hyperlink,
		HyperlinkType:       opt.HyperlinkType,
		Positioning:         opt.Positioning,
	}
}

type ExcelPictureInsertType int

// Insert picture types.
const (
	PictureInsertTypePlaceOverCells = ExcelPictureInsertType(excelize.PictureInsertTypePlaceOverCells)
	PictureInsertTypePlaceInCell    = ExcelPictureInsertType(excelize.PictureInsertTypePlaceInCell)
	PictureInsertTypeIMAGE          = ExcelPictureInsertType(excelize.PictureInsertTypeIMAGE)
	PictureInsertTypeDISPIMG        = ExcelPictureInsertType(excelize.PictureInsertTypeDISPIMG)
)

type ExcelPicture struct {
	Extension  string
	File       []byte
	Format     *ExcelGraphicOptions
	InsertType ExcelPictureInsertType
}

func (ep *ExcelPicture) ToExcelize() *excelize.Picture {
	v := &excelize.Picture{
		Extension:  ep.Extension,
		File:       ep.File,
		InsertType: excelize.PictureInsertType(ep.InsertType),
	}
	if nil != ep.Format {
		v.Format = ep.Format.ToExcelize()
	}
	return v
}

func excelPictureFromExcelize(pic *excelize.Picture) *ExcelPicture {
	return &ExcelPicture{
		Extension:  pic.Extension,
		File:       pic.File,
		Format:     excelGraphicOptionsFromExcelize(pic.Format),
		InsertType: ExcelPictureInsertType(pic.InsertType),
	}
}
func excelPicturesFromExcelize(pics []excelize.Picture) (excelPics []*ExcelPicture) {
	for _, pic := range pics {
		excelPics = append(excelPics, excelPictureFromExcelize(&pic))
	}
	return
}

type ExcelFont struct {
	Bold         bool
	Italic       bool
	Underline    bool
	Family       string
	Size         float64
	Strike       bool
	Color        string
	ColorIndexed int
	ColorTheme   int
	ColorTint    float64
	VertAlign    string
}

func (ef *ExcelFont) ToExcelize() *excelize.Font {
	v := &excelize.Font{
		Bold:         ef.Bold,
		Italic:       ef.Italic,
		Underline:    "none",
		Family:       ef.Family,
		Size:         ef.Size,
		Strike:       ef.Strike,
		Color:        ef.Color,
		ColorIndexed: ef.ColorIndexed,
		ColorTheme:   nil,
		ColorTint:    ef.ColorTint,
		VertAlign:    ef.VertAlign,
	}
	if ef.Underline {
		v.Underline = "single"
	}
	if ef.ColorTheme > 0 {
		v.ColorTheme = &ef.ColorTheme
	}
	return v
}

func excelFontFromExcelize(font *excelize.Font) *ExcelFont {
	v := &ExcelFont{
		Bold:         font.Bold,
		Italic:       font.Italic,
		Underline:    false,
		Family:       font.Family,
		Size:         font.Size,
		Strike:       font.Strike,
		Color:        font.Color,
		ColorIndexed: font.ColorIndexed,
		ColorTheme:   0,
		ColorTint:    font.ColorTint,
		VertAlign:    font.VertAlign,
	}

	if "none" != font.Underline {
		v.Underline = true
	}
	if nil != font.ColorTheme {
		v.ColorTheme = *font.ColorTheme
	}

	return v
}

type ExcelRichText struct {
	Font *ExcelFont
	Text string
}

func (ert *ExcelRichText) ToExcelize() *excelize.RichTextRun {
	v := &excelize.RichTextRun{
		Text: ert.Text,
	}
	if nil != ert.Font {
		v.Font = ert.Font.ToExcelize()
	}

	return v
}
func excelRichTextFromExcelize(rt *excelize.RichTextRun) *ExcelRichText {
	return &ExcelRichText{
		Font: excelFontFromExcelize(rt.Font),
		Text: rt.Text,
	}
}
func excelRichTextsFromExcelize(rts []excelize.RichTextRun) (texts []*ExcelRichText) {
	for _, rt := range rts {
		texts = append(texts, excelRichTextFromExcelize(&rt))
	}
	return
}
func excelRichTextsToExcelize(rts []*ExcelRichText) (texts []excelize.RichTextRun) {
	for _, rt := range rts {
		texts = append(texts, *rt.ToExcelize())
	}
	return
}

type ExcelBorder struct {
	Type  string
	Color string
	Style int
}

func (eb *ExcelBorder) ToExcelize() *excelize.Border {
	return &excelize.Border{
		Type:  eb.Type,
		Color: eb.Color,
		Style: eb.Style,
	}
}
func excelBorderFromExcelize(eb *excelize.Border) *ExcelBorder {
	return &ExcelBorder{
		Type:  eb.Type,
		Color: eb.Color,
		Style: eb.Style,
	}
}

type ExcelFill struct {
	Type    string
	Pattern int
	Color   []string
	Shading int
}

func (ef *ExcelFill) ToExcelize() *excelize.Fill {
	return &excelize.Fill{
		Type:    ef.Type,
		Pattern: ef.Pattern,
		Color:   ef.Color,
		Shading: ef.Shading,
	}
}
func excelFillFromExcelize(ef *excelize.Fill) *ExcelFill {
	return &ExcelFill{
		Type:    ef.Type,
		Pattern: ef.Pattern,
		Color:   ef.Color,
		Shading: ef.Shading,
	}
}

type ExcelProtection struct {
	Hidden bool
	Locked bool
}

func (ep *ExcelProtection) ToExcelize() *excelize.Protection {
	return &excelize.Protection{
		Hidden: ep.Hidden,
		Locked: ep.Locked,
	}
}
func ExcelProtectionFromExcelize(ep *excelize.Protection) *ExcelProtection {
	return &ExcelProtection{
		Hidden: ep.Hidden,
		Locked: ep.Locked,
	}
}

type ExcelAlignment struct {
	Horizontal      string
	Indent          int
	JustifyLastLine bool
	ReadingOrder    uint64
	RelativeIndent  int
	ShrinkToFit     bool
	TextRotation    int
	Vertical        string
	WrapText        bool
}

func (ea *ExcelAlignment) ToExcelize() *excelize.Alignment {
	return &excelize.Alignment{
		Horizontal:      ea.Horizontal,
		Indent:          ea.Indent,
		JustifyLastLine: ea.JustifyLastLine,
		ReadingOrder:    ea.ReadingOrder,
		RelativeIndent:  ea.RelativeIndent,
		ShrinkToFit:     ea.ShrinkToFit,
		TextRotation:    ea.TextRotation,
		Vertical:        ea.Vertical,
		WrapText:        ea.WrapText,
	}
}
func excelAlignmentFromExcelize(ea *excelize.Alignment) *ExcelAlignment {
	return &ExcelAlignment{
		Horizontal:      ea.Horizontal,
		Indent:          ea.Indent,
		JustifyLastLine: ea.JustifyLastLine,
		ReadingOrder:    ea.ReadingOrder,
		RelativeIndent:  ea.RelativeIndent,
		ShrinkToFit:     ea.ShrinkToFit,
		TextRotation:    ea.TextRotation,
		Vertical:        ea.Vertical,
		WrapText:        ea.WrapText,
	}
}

type ExcelStyle struct {
	Border        []*ExcelBorder
	Fill          *ExcelFill
	Font          *ExcelFont
	Alignment     *ExcelAlignment
	Protection    *ExcelProtection
	NumFmt        int
	DecimalPlaces int
	CustomNumFmt  string
	NegRed        bool
}

func (es *ExcelStyle) ToExcelize() *excelize.Style {
	style := &excelize.Style{
		NumFmt: es.NumFmt,
		NegRed: es.NegRed,
	}
	if nil != es.Border {
		for _, b := range es.Border {
			style.Border = append(style.Border, *(b.ToExcelize()))
		}
	}
	if nil != es.Fill {
		style.Fill = *(es.Fill.ToExcelize())
	}

	if nil != es.Font {
		style.Font = es.Font.ToExcelize()
	}

	if nil != es.Alignment {
		style.Alignment = es.Alignment.ToExcelize()
	}
	if nil != es.Protection {
		style.Protection = es.Protection.ToExcelize()
	}

	if es.DecimalPlaces > 0 {
		style.DecimalPlaces = &es.DecimalPlaces
	}

	if "" != es.CustomNumFmt {
		style.CustomNumFmt = &es.CustomNumFmt
	}

	return style
}
func excelStyleFromExcelize(es *excelize.Style) *ExcelStyle {
	style := &ExcelStyle{
		Fill:          excelFillFromExcelize(&es.Fill),
		NumFmt:        es.NumFmt,
		DecimalPlaces: *es.DecimalPlaces,
		CustomNumFmt:  *es.CustomNumFmt,
		NegRed:        es.NegRed,
	}
	if nil != es.Border {
		for _, b := range es.Border {
			style.Border = append(style.Border, excelBorderFromExcelize(&b))
		}
	}

	if nil != es.Font {
		style.Font = excelFontFromExcelize(es.Font)
	}

	if nil != es.Alignment {
		style.Alignment = excelAlignmentFromExcelize(es.Alignment)
	}
	if nil != es.Protection {
		style.Protection = ExcelProtectionFromExcelize(es.Protection)
	}

	return style
}
