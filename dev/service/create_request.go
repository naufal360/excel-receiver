package service

import (
	"encoding/csv"
	"excel-receiver/config"
	"excel-receiver/constant"
	"excel-receiver/entity"
	"excel-receiver/http/api/ierr"
	"excel-receiver/provider"
	"excel-receiver/repository"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
)

type sendRequest struct {
	log         provider.ILogger
	artemisRepo repository.QueueInterface
	requestRepo repository.RequestInterface
}

func NewSendRequestService(log provider.ILogger, artemisRepo repository.QueueInterface, requestRepo repository.RequestInterface) *sendRequest {
	return &sendRequest{
		log:         log,
		artemisRepo: artemisRepo,
		requestRepo: requestRepo,
	}
}

func (s *sendRequest) SendRequest(ctx *gin.Context, file *multipart.FileHeader, extensionFile string) (status string, err error) {
	var (
		uploadDir      = config.Configuration.Server.UploadDir
		queueList      []entity.Queue
		requestPayload entity.Request
		reqID          = ctx.Request.Context().Value(constant.RequestIDKey{}).(string)
	)

	openFile, err := file.Open()
	if err != nil {
		return "failed", err
	}

	switch extensionFile {
	case ".xlsx":
		queueList, status, err = s.readXLSX(openFile, file.Filename, reqID)
	case ".csv":
		queueList, status, err = s.readCSV(openFile, file.Filename, reqID)
	default:
		return "failed", ierr.NewF(constant.InvalidRequest, "")
	}
	if err != nil {
		return status, err
	}

	for _, queue := range queueList {
		if err = s.artemisRepo.ProduceQueue(ctx.Request.Context(), &queue); err != nil {
			s.log.ErrorWithFields(provider.AppLog, map[string]interface{}{
				constant.ReqIDLog: reqID,
				"DATA":            queue,
				"ERROR":           err,
			}, "error when produce queue")
			return "failed", err
		}
	}

	filepath := fmt.Sprintf("/%s%s", reqID, extensionFile)
	dst := fmt.Sprintf(".%s%s", uploadDir, filepath)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		return "failed", err
	}

	requestPayload = entity.Request{
		RequestID: reqID,
		Status:    "received",
		Filepath:  filepath,
		CreatedAt: time.Now(),
	}
	status, err = s.requestRepo.CreateRequest(ctx.Request.Context(), &requestPayload)

	return
}

func (s *sendRequest) readXLSX(file multipart.File, filename, reqID string) (queueList []entity.Queue, status string, err error) {
	var (
		sheetName = config.Configuration.Server.SheetName
	)

	readFile, err := excelize.OpenReader(file)
	if err != nil {
		return queueList, "failed", err
	}

	rows := readFile.GetRows(sheetName)
	totalRows := len(rows)
	dataRows := totalRows - 1 // skip count rows name col

	if dataRows < 1 || dataRows > 10 {
		s.log.ErrorWithFields(provider.AppLog, map[string]interface{}{
			constant.ReqIDLog: reqID,
			"ERROR":           "invalid row data length",
		}, "error when produce queue")
		return queueList, "failed", ierr.NewF(constant.InvalidRowData, "")
	}

	queueList, status, err = s.checkMandatory(reqID, rows)

	return
}

func (s *sendRequest) readCSV(file multipart.File, filename, reqID string) (queueList []entity.Queue, status string, err error) {
	var (
		rows                [][]string
		dataRows, totalData int
	)

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return queueList, "failed", err
		}

		rows = append(rows, record)
		totalData += 1
	}

	dataRows = totalData - 1 // skip count rows name col

	if dataRows < 1 || dataRows > 10 {
		s.log.ErrorWithFields(provider.AppLog, map[string]interface{}{
			constant.ReqIDLog: reqID,
			"ERROR":           "invalid row data length",
		}, "error when produce queue")
		return queueList, "failed", ierr.NewF(constant.InvalidRowData, "")
	}

	queueList, status, err = s.checkMandatory(reqID, rows)

	return
}

func (s *sendRequest) checkMandatory(reqID string, rows [][]string) (queueList []entity.Queue, status string, err error) {
	mandatoryColumns := map[string]bool{
		"uniqid":      false,
		"title":       false,
		"description": false,
		"condition":   false,
		"price":       false,
	}

	for idx, row := range rows {
		if idx == 0 {
			for _, col := range row {
				if _, ok := mandatoryColumns[col]; ok {
					mandatoryColumns[col] = true
				}
			}

			// Check if all mandatory columns are present
			for col, present := range mandatoryColumns {
				if !present {
					return queueList, "failed", ierr.NewF(constant.EmptyRowMandatory, col)
				}
			}
		} else {
			var price float64
			var weight float64

			queue := entity.Queue{
				RequestID: reqID,
			}

			for i, col := range row {
				switch colName := rows[0][i]; colName {
				case "uniqid":
					queue.UniqID = col
				case "description":
					queue.Description = col
				case "condition":
					queue.Condition = col
				case "price":
					if col != "" {
						price, _ = strconv.ParseFloat(col, 64)
						queue.Price = price
					}
				case "color":
					queue.Color = col
				case "size":
					queue.Size = col
				case "age_group":
					queue.AgeGroup = col
				case "material":
					queue.Material = col
				case "weight_kg":
					if col != "" {
						weight, _ = strconv.ParseFloat(col, 64)
						queue.WeightKG = weight
					}
				}
			}

			queueList = append(queueList, queue)
		}
	}
	return queueList, status, err
}
