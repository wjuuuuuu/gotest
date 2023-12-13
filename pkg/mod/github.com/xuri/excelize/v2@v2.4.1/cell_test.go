package excelize

import (
	"fmt"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	_ "image/jpeg"

	"github.com/stretchr/testify/assert"
)

func TestConcurrency(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "Book1.xlsx"))
	assert.NoError(t, err)
	wg := new(sync.WaitGroup)
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(val int, t *testing.T) {
			// Concurrency set cell value
			assert.NoError(t, f.SetCellValue("Sheet1", fmt.Sprintf("A%d", val), val))
			assert.NoError(t, f.SetCellValue("Sheet1", fmt.Sprintf("B%d", val), strconv.Itoa(val)))
			// Concurrency get cell value
			_, err := f.GetCellValue("Sheet1", fmt.Sprintf("A%d", val))
			assert.NoError(t, err)
			// Concurrency set rows
			assert.NoError(t, f.SetSheetRow("Sheet1", "B6", &[]interface{}{" Hello",
				[]byte("World"), 42, int8(1<<8/2 - 1), int16(1<<16/2 - 1), int32(1<<32/2 - 1),
				int64(1<<32/2 - 1), float32(42.65418), float64(-42.65418), float32(42), float64(42),
				uint(1<<32 - 1), uint8(1<<8 - 1), uint16(1<<16 - 1), uint32(1<<32 - 1),
				uint64(1<<32 - 1), true, complex64(5 + 10i)}))
			// Concurrency create style
			style, err := f.NewStyle(`{"font":{"color":"#1265BE","underline":"single"}}`)
			assert.NoError(t, err)
			// Concurrency set cell style
			assert.NoError(t, f.SetCellStyle("Sheet1", "A3", "A3", style))
			// Concurrency add picture
			assert.NoError(t, f.AddPicture("Sheet1", "F21", filepath.Join("test", "images", "excel.jpg"),
				`{"x_offset": 10, "y_offset": 10, "hyperlink": "https://github.com/xuri/excelize", "hyperlink_type": "External", "positioning": "oneCell"}`))
			// Concurrency get cell picture
			name, raw, err := f.GetPicture("Sheet1", "A1")
			assert.Equal(t, "", name)
			assert.Nil(t, raw)
			assert.NoError(t, err)
			// Concurrency iterate rows
			rows, err := f.Rows("Sheet1")
			assert.NoError(t, err)
			for rows.Next() {
				_, err := rows.Columns()
				assert.NoError(t, err)
			}
			// Concurrency iterate columns
			cols, err := f.Cols("Sheet1")
			assert.NoError(t, err)
			for rows.Next() {
				_, err := cols.Rows()
				assert.NoError(t, err)
			}

			wg.Done()
		}(i, t)
	}
	wg.Wait()
	val, err := f.GetCellValue("Sheet1", "A1")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "1", val)
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestConcurrency.xlsx")))
}

func TestCheckCellInArea(t *testing.T) {
	f := NewFile()
	expectedTrueCellInAreaList := [][2]string{
		{"c2", "A1:AAZ32"},
		{"B9", "A1:B9"},
		{"C2", "C2:C2"},
	}

	for _, expectedTrueCellInArea := range expectedTrueCellInAreaList {
		cell := expectedTrueCellInArea[0]
		area := expectedTrueCellInArea[1]
		ok, err := f.checkCellInArea(cell, area)
		assert.NoError(t, err)
		assert.Truef(t, ok,
			"Expected cell %v to be in area %v, got false\n", cell, area)
	}

	expectedFalseCellInAreaList := [][2]string{
		{"c2", "A4:AAZ32"},
		{"C4", "D6:A1"}, // weird case, but you never know
		{"AEF42", "BZ40:AEF41"},
	}

	for _, expectedFalseCellInArea := range expectedFalseCellInAreaList {
		cell := expectedFalseCellInArea[0]
		area := expectedFalseCellInArea[1]
		ok, err := f.checkCellInArea(cell, area)
		assert.NoError(t, err)
		assert.Falsef(t, ok,
			"Expected cell %v not to be inside of area %v, but got true\n", cell, area)
	}

	ok, err := f.checkCellInArea("A1", "A:B")
	assert.EqualError(t, err, `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	assert.False(t, ok)

	ok, err = f.checkCellInArea("AA0", "Z0:AB1")
	assert.EqualError(t, err, `cannot convert cell "AA0" to coordinates: invalid cell name "AA0"`)
	assert.False(t, ok)
}

func TestSetCellFloat(t *testing.T) {
	sheet := "Sheet1"
	t.Run("with no decimal", func(t *testing.T) {
		f := NewFile()
		assert.NoError(t, f.SetCellFloat(sheet, "A1", 123.0, -1, 64))
		assert.NoError(t, f.SetCellFloat(sheet, "A2", 123.0, 1, 64))
		val, err := f.GetCellValue(sheet, "A1")
		assert.NoError(t, err)
		assert.Equal(t, "123", val, "A1 should be 123")
		val, err = f.GetCellValue(sheet, "A2")
		assert.NoError(t, err)
		assert.Equal(t, "123.0", val, "A2 should be 123.0")
	})

	t.Run("with a decimal and precision limit", func(t *testing.T) {
		f := NewFile()
		assert.NoError(t, f.SetCellFloat(sheet, "A1", 123.42, 1, 64))
		val, err := f.GetCellValue(sheet, "A1")
		assert.NoError(t, err)
		assert.Equal(t, "123.4", val, "A1 should be 123.4")
	})

	t.Run("with a decimal and no limit", func(t *testing.T) {
		f := NewFile()
		assert.NoError(t, f.SetCellFloat(sheet, "A1", 123.42, -1, 64))
		val, err := f.GetCellValue(sheet, "A1")
		assert.NoError(t, err)
		assert.Equal(t, "123.42", val, "A1 should be 123.42")
	})
	f := NewFile()
	assert.EqualError(t, f.SetCellFloat(sheet, "A", 123.42, -1, 64), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestSetCellValue(t *testing.T) {
	f := NewFile()
	assert.EqualError(t, f.SetCellValue("Sheet1", "A", time.Now().UTC()), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	assert.EqualError(t, f.SetCellValue("Sheet1", "A", time.Duration(1e13)), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestSetCellValues(t *testing.T) {
	f := NewFile()
	err := f.SetCellValue("Sheet1", "A1", time.Date(2010, time.December, 31, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)

	v, err := f.GetCellValue("Sheet1", "A1")
	assert.NoError(t, err)
	assert.Equal(t, v, "12/31/10 00:00")

	// test date value lower than min date supported by Excel
	err = f.SetCellValue("Sheet1", "A1", time.Date(1600, time.December, 31, 0, 0, 0, 0, time.UTC))
	assert.NoError(t, err)

	v, err = f.GetCellValue("Sheet1", "A1")
	assert.NoError(t, err)
	assert.Equal(t, v, "1600-12-31T00:00:00Z")
}

func TestSetCellBool(t *testing.T) {
	f := NewFile()
	assert.EqualError(t, f.SetCellBool("Sheet1", "A", true), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}

func TestGetCellValue(t *testing.T) {
	// Test get cell value without r attribute of the row.
	f := NewFile()
	sheetData := `<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>%s</sheetData></worksheet>`
	f.Sheet.Delete("xl/worksheets/sheet1.xml")
	f.Pkg.Store("xl/worksheets/sheet1.xml", []byte(fmt.Sprintf(sheetData, `<row r="3"><c t="str"><v>A3</v></c></row><row><c t="str"><v>A4</v></c><c t="str"><v>B4</v></c></row><row r="7"><c t="str"><v>A7</v></c><c t="str"><v>B7</v></c></row><row><c t="str"><v>A8</v></c><c t="str"><v>B8</v></c></row>`)))
	f.checked = nil
	cells := []string{"A3", "A4", "B4", "A7", "B7"}
	rows, err := f.GetRows("Sheet1")
	assert.Equal(t, [][]string{nil, nil, {"A3"}, {"A4", "B4"}, nil, nil, {"A7", "B7"}, {"A8", "B8"}}, rows)
	assert.NoError(t, err)
	for _, cell := range cells {
		value, err := f.GetCellValue("Sheet1", cell)
		assert.Equal(t, cell, value)
		assert.NoError(t, err)
	}
	cols, err := f.GetCols("Sheet1")
	assert.Equal(t, [][]string{{"", "", "A3", "A4", "", "", "A7", "A8"}, {"", "", "", "B4", "", "", "B7", "B8"}}, cols)
	assert.NoError(t, err)
	f.Sheet.Delete("xl/worksheets/sheet1.xml")
	f.Pkg.Store("xl/worksheets/sheet1.xml", []byte(fmt.Sprintf(sheetData, `<row r="2"><c r="A2" t="str"><v>A2</v></c></row><row r="2"><c r="B2" t="str"><v>B2</v></c></row>`)))
	f.checked = nil
	cell, err := f.GetCellValue("Sheet1", "A2")
	assert.Equal(t, "A2", cell)
	assert.NoError(t, err)
	f.Sheet.Delete("xl/worksheets/sheet1.xml")
	f.Pkg.Store("xl/worksheets/sheet1.xml", []byte(fmt.Sprintf(sheetData, `<row r="2"><c r="A2" t="str"><v>A2</v></c></row><row r="2"><c r="B2" t="str"><v>B2</v></c></row>`)))
	f.checked = nil
	rows, err = f.GetRows("Sheet1")
	assert.Equal(t, [][]string{nil, {"A2", "B2"}}, rows)
	assert.NoError(t, err)
	f.Sheet.Delete("xl/worksheets/sheet1.xml")
	f.Pkg.Store("xl/worksheets/sheet1.xml", []byte(fmt.Sprintf(sheetData, `<row r="1"><c r="A1" t="str"><v>A1</v></c></row><row r="1"><c r="B1" t="str"><v>B1</v></c></row>`)))
	f.checked = nil
	rows, err = f.GetRows("Sheet1")
	assert.Equal(t, [][]string{{"A1", "B1"}}, rows)
	assert.NoError(t, err)
}

func TestGetCellFormula(t *testing.T) {
	// Test get cell formula on not exist worksheet.
	f := NewFile()
	_, err := f.GetCellFormula("SheetN", "A1")
	assert.EqualError(t, err, "sheet SheetN is not exist")

	// Test get cell formula on no formula cell.
	assert.NoError(t, f.SetCellValue("Sheet1", "A1", true))
	_, err = f.GetCellFormula("Sheet1", "A1")
	assert.NoError(t, err)
}

func ExampleFile_SetCellFloat() {
	f := NewFile()
	var x = 3.14159265
	if err := f.SetCellFloat("Sheet1", "A1", x, 2, 64); err != nil {
		fmt.Println(err)
	}
	val, _ := f.GetCellValue("Sheet1", "A1")
	fmt.Println(val)
	// Output: 3.14
}

func BenchmarkSetCellValue(b *testing.B) {
	values := []string{"First", "Second", "Third", "Fourth", "Fifth", "Sixth"}
	cols := []string{"A", "B", "C", "D", "E", "F"}
	f := NewFile()
	b.ResetTimer()
	for i := 1; i <= b.N; i++ {
		for j := 0; j < len(values); j++ {
			if err := f.SetCellValue("Sheet1", cols[j]+strconv.Itoa(i), values[j]); err != nil {
				b.Error(err)
			}
		}
	}
}

func TestOverflowNumericCell(t *testing.T) {
	f, err := OpenFile(filepath.Join("test", "OverflowNumericCell.xlsx"))
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	val, err := f.GetCellValue("Sheet1", "A1")
	assert.NoError(t, err)
	// GOARCH=amd64 - all ok; GOARCH=386 - actual: "-2147483648"
	assert.Equal(t, "8595602512225", val, "A1 should be 8595602512225")
}
func TestGetCellRichText(t *testing.T) {
	f := NewFile()

	runsSource := []RichTextRun{
		{
			Text: "a\n",
		},
		{
			Text: "b",
			Font: &Font{
				Underline: "single",
				Color:     "ff0000",
				Bold:      true,
				Italic:    true,
				Family:    "Times New Roman",
				Size:      100,
				Strike:    true,
			},
		},
	}
	assert.NoError(t, f.SetCellRichText("Sheet1", "A1", runsSource))

	runs, err := f.GetCellRichText("Sheet1", "A1")
	assert.NoError(t, err)

	assert.Equal(t, runsSource[0].Text, runs[0].Text)
	assert.Nil(t, runs[0].Font)
	assert.NotNil(t, runs[1].Font)

	runsSource[1].Font.Color = strings.ToUpper(runsSource[1].Font.Color)
	assert.True(t, reflect.DeepEqual(runsSource[1].Font, runs[1].Font), "should get the same font")

	// Test get cell rich text when string item index overflow
	ws, ok := f.Sheet.Load("xl/worksheets/sheet1.xml")
	assert.True(t, ok)
	ws.(*xlsxWorksheet).SheetData.Row[0].C[0].V = "2"
	runs, err = f.GetCellRichText("Sheet1", "A1")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(runs))
	// Test get cell rich text when string item index is negative
	ws, ok = f.Sheet.Load("xl/worksheets/sheet1.xml")
	assert.True(t, ok)
	ws.(*xlsxWorksheet).SheetData.Row[0].C[0].V = "-1"
	runs, err = f.GetCellRichText("Sheet1", "A1")
	assert.NoError(t, err)
	assert.Equal(t, 0, len(runs))
	// Test get cell rich text on invalid string item index
	ws, ok = f.Sheet.Load("xl/worksheets/sheet1.xml")
	assert.True(t, ok)
	ws.(*xlsxWorksheet).SheetData.Row[0].C[0].V = "x"
	_, err = f.GetCellRichText("Sheet1", "A1")
	assert.EqualError(t, err, "strconv.Atoi: parsing \"x\": invalid syntax")
	// Test set cell rich text on not exists worksheet
	_, err = f.GetCellRichText("SheetN", "A1")
	assert.EqualError(t, err, "sheet SheetN is not exist")
	// Test set cell rich text with illegal cell coordinates
	_, err = f.GetCellRichText("Sheet1", "A")
	assert.EqualError(t, err, `cannot convert cell "A" to coordinates: invalid cell name "A"`)
}
func TestSetCellRichText(t *testing.T) {
	f := NewFile()
	assert.NoError(t, f.SetRowHeight("Sheet1", 1, 35))
	assert.NoError(t, f.SetColWidth("Sheet1", "A", "A", 44))
	richTextRun := []RichTextRun{
		{
			Text: "bold",
			Font: &Font{
				Bold:   true,
				Color:  "2354e8",
				Family: "Times New Roman",
			},
		},
		{
			Text: " and ",
			Font: &Font{
				Family: "Times New Roman",
			},
		},
		{
			Text: "italic ",
			Font: &Font{
				Bold:   true,
				Color:  "e83723",
				Italic: true,
				Family: "Times New Roman",
			},
		},
		{
			Text: "text with color and font-family,",
			Font: &Font{
				Bold:   true,
				Color:  "2354e8",
				Family: "Times New Roman",
			},
		},
		{
			Text: "\r\nlarge text with ",
			Font: &Font{
				Size:  14,
				Color: "ad23e8",
			},
		},
		{
			Text: "strike",
			Font: &Font{
				Color:  "e89923",
				Strike: true,
			},
		},
		{
			Text: " and ",
			Font: &Font{
				Size:  14,
				Color: "ad23e8",
			},
		},
		{
			Text: "underline.",
			Font: &Font{
				Color:     "23e833",
				Underline: "single",
			},
		},
	}
	assert.NoError(t, f.SetCellRichText("Sheet1", "A1", richTextRun))
	assert.NoError(t, f.SetCellRichText("Sheet1", "A2", richTextRun))
	style, err := f.NewStyle(&Style{
		Alignment: &Alignment{
			WrapText: true,
		},
	})
	assert.NoError(t, err)
	assert.NoError(t, f.SetCellStyle("Sheet1", "A1", "A1", style))
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestSetCellRichText.xlsx")))
	// Test set cell rich text on not exists worksheet
	assert.EqualError(t, f.SetCellRichText("SheetN", "A1", richTextRun), "sheet SheetN is not exist")
	// Test set cell rich text with illegal cell coordinates
	assert.EqualError(t, f.SetCellRichText("Sheet1", "A", richTextRun), `cannot convert cell "A" to coordinates: invalid cell name "A"`)
	richTextRun = []RichTextRun{{Text: strings.Repeat("s", TotalCellChars+1)}}
	// Test set cell rich text with characters over the maximum limit
	assert.EqualError(t, f.SetCellRichText("Sheet1", "A1", richTextRun), ErrCellCharsLength.Error())
}

func TestFormattedValue2(t *testing.T) {
	f := NewFile()
	v := f.formattedValue(0, "43528")
	assert.Equal(t, "43528", v)

	v = f.formattedValue(15, "43528")
	assert.Equal(t, "43528", v)

	v = f.formattedValue(1, "43528")
	assert.Equal(t, "43528", v)
	customNumFmt := "[$-409]MM/DD/YYYY"
	_, err := f.NewStyle(&Style{
		CustomNumFmt: &customNumFmt,
	})
	assert.NoError(t, err)
	v = f.formattedValue(1, "43528")
	assert.Equal(t, "03/04/2019", v)

	// formatted value with no built-in number format ID
	numFmtID := 5
	f.Styles.CellXfs.Xf = append(f.Styles.CellXfs.Xf, xlsxXf{
		NumFmtID: &numFmtID,
	})
	v = f.formattedValue(2, "43528")
	assert.Equal(t, "43528", v)

	// formatted value with invalid number format ID
	f.Styles.CellXfs.Xf = append(f.Styles.CellXfs.Xf, xlsxXf{
		NumFmtID: nil,
	})
	_ = f.formattedValue(3, "43528")

	// formatted value with empty number format
	f.Styles.NumFmts = nil
	f.Styles.CellXfs.Xf = append(f.Styles.CellXfs.Xf, xlsxXf{
		NumFmtID: &numFmtID,
	})
	v = f.formattedValue(1, "43528")
	assert.Equal(t, "43528", v)
}
