package usecase

import (
	"bytes"
	"fmt"
	"path/filepath"
	"quiz-app/internal/domain/entity"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFService interface {
	GenerateTestResultPDF(data *entity.ClassExport) ([]byte, error)
}

type PDFServiceImpl struct{}

func NewPDFService() PDFService {
	return &PDFServiceImpl{}
}

func (s *PDFServiceImpl) GenerateTestResultPDF(
	data *entity.ClassExport,
) ([]byte, error) {

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 10)

	// ===== Load UTF-8 fonts =====
	pdf.AddUTF8Font("Roboto", "", filepath.Join("assets", "fonts", "Roboto-Regular.ttf"))
	pdf.AddUTF8Font("Roboto", "B", filepath.Join("assets", "fonts", "Roboto-Bold.ttf"))

	if pdf.Error() != nil {
		return nil, pdf.Error()
	}

	pdf.AddPage()

	// ===== Title =====
	pdf.SetFont("Roboto", "B", 16)
	pdf.Cell(0, 10, "KẾT QUẢ BÀI KIỂM TRA")
	pdf.Ln(12)

	// ===== Class & Test info =====
	pdf.SetFont("Roboto", "", 12)
	pdf.Cell(0, 8, "Lớp: "+data.ClassName)
	pdf.Ln(6)

	pdf.Cell(0, 8, "Bài test: "+data.Test.TestName)
	pdf.Ln(6)

	pdf.Cell(0, 8, fmt.Sprintf(
		"Thời gian: %d phút | Tổng điểm: %.1f",
		data.Test.DurationMinutes,
		data.Test.TestScore,
	))
	pdf.Ln(6)

	pdf.Cell(0, 8, fmt.Sprintf(
		"Mở: %s  -  Đóng: %s",
		formatTime(data.Test.StartTime),
		formatTime(data.Test.EndTime),
	))
	pdf.Ln(10)

	// ===== Table header =====
	pdf.SetFont("Roboto", "B", 11)
	pdf.CellFormat(55, 8, "Họ tên", "1", 0, "C", false, 0, "")
	pdf.CellFormat(65, 8, "Email", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 8, "Điểm", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Trạng thái", "1", 1, "C", false, 0, "")

	// ===== Table body =====
	pdf.SetFont("Roboto", "", 11)

	for _, u := range data.Test.Users {

		fullName := fmt.Sprintf("%s %s", u.FirstName, u.LastName)

		status := "Chưa nộp"
		if u.Submitted {
			status = "Đã nộp"
		}

		pdf.CellFormat(55, 8, fullName, "1", 0, "", false, 0, "")
		pdf.CellFormat(65, 8, u.Email, "1", 0, "", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%.1f", u.Score), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, status, "1", 1, "C", false, 0, "")
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.Format("02/01/2006 15:04")
}
