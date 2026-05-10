package usecase

import (
	"context"
	"errors"
	"fmt"

	"quiz-app/internal/domain/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubmissionUseCase struct {
	submissionRepo repository.SubmissionRepository
	pdfService     PDFService
}

func NewSubmissionUseCase(
	repo repository.SubmissionRepository,
	pdfService PDFService,
) *SubmissionUseCase {
	return &SubmissionUseCase{
		submissionRepo: repo,
		pdfService:     pdfService,
	}
}
func (uc *SubmissionUseCase) ExportPDF(
	ctx context.Context,
	classID primitive.ObjectID,
	testID primitive.ObjectID,
) ([]byte, string, error) {

	// 1. Lấy dữ liệu từ Mongo (aggregation)
	data, err := uc.submissionRepo.ExportData(ctx, classID, testID)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	if data == nil {
		return nil, "", errors.New("no data to export")
	}

	// 2. Generate PDF
	pdfBytes, err := uc.pdfService.GenerateTestResultPDF(data)
	if err != nil {
		return nil, "", err
	}

	// 3. File name
	fileName := fmt.Sprintf(
		"ket-qua-%s-%s.pdf",
		data.ClassName,
		data.Test.TestName,
	)

	return pdfBytes, fileName, nil
}
