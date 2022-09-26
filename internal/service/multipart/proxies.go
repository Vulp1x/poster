package multipart

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	tasksservice "github.com/inst-api/poster/gen/tasks_service"
	"github.com/inst-api/poster/pkg/logger"
	"go.uber.org/multierr"
)

func readProxiesList(ctx context.Context, reader io.Reader) ([]*tasksservice.ProxyRecord, error) {
	csvReader := csv.NewReader(reader)

	csvReader.Comma = ':'
	csvReader.FieldsPerRecord = 4

	var botAccounts []*tasksservice.ProxyRecord
	var errs []error

	var line int
	for {
		record, err := csvReader.Read()
		line, _ = csvReader.FieldPos(0)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break // readed all file
			}

			errs = append(errs, fmt.Errorf("failed to read from file at line %d: %v", line, err))
			continue
		}

		botAccounts = append(botAccounts, &tasksservice.ProxyRecord{
			Record:     record,
			LineNumber: line,
		})
	}

	logger.Debugf(ctx, "read %d lines\n", line)

	return botAccounts, multierr.Combine(errs...)
}
