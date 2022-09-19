package reader

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/pkg/logger"
)

func ParseUsersList(ctx context.Context, reader io.Reader) ([]domain.BotAccount, []error) {
	csvReader := csv.NewReader(reader)

	csvReader.Comma = '|'
	csvReader.FieldsPerRecord = 4

	var botAccounts []domain.BotAccount
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

		var botAccount domain.BotAccount

		err = botAccount.Parse(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to parse bot account on line %d: %v", line, err))
			continue
		}

		botAccounts = append(botAccounts, botAccount)
	}

	logger.Debugf(ctx, "read %d lines\n", line)

	return botAccounts, errs
}
