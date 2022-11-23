package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/dbtx"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/store"
	"github.com/inst-api/poster/pkg/logger"
	"github.com/jackc/pgx/v4"
)

func (s *Store) AssignProxies(ctx context.Context, taskID uuid.UUID) (int, error) {
	tx, err := s.txf(ctx)
	if err != nil {
		return 0, store.ErrTransactionFail
	}

	defer dbtx.RollbackUnlessCommitted(ctx, tx)

	q := dbmodel.New(tx)

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrTaskNotFound
		}

		return 0, fmt.Errorf("failed to find task with id '%s': %v", taskID, err)
	}

	if task.Status != dbmodel.DataUploadedTaskStatus {
		return 0, fmt.Errorf("%w: expected %d got %d", ErrTaskInvalidStatus, dbmodel.DataUploadedTaskStatus, task.Status)
	}

	botAccounts, err := q.FindBotsForTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("failed to find bot accounts for task: %v", err)
	}

	residentialProxies, err := q.FindResidentialProxiesForTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("failed to find proxiesIds for task: %v", err)
	}

	cheapProxies, err := q.FindCheapProxiesForTask(ctx, taskID)
	if err != nil {
		return 0, fmt.Errorf("failed to find cheap proxies for task: %v", err)
	}

	// after deleting botAccounts and residentialProxies would have same length
	botAccounts, residentialProxies, cheapProxies, err = s.deleteUnnecessaryRows(ctx, tx, botAccounts, residentialProxies, cheapProxies)
	if err != nil {
		return 0, err
	}

	botIds := domain.Ids(botAccounts)
	err = q.AssignProxiesToBotsForTask(ctx, dbmodel.AssignProxiesToBotsForTaskParams{
		TaskID:             taskID,
		ResidentialProxies: domain.Strings(residentialProxies),
		CheapProxies:       domain.Strings(cheapProxies),
		Ids:                botIds,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to set bot accounts cheap proxies: %v", err)
	}

	err = q.AssignBotsToProxiesForTask(ctx, dbmodel.AssignBotsToProxiesForTaskParams{
		TaskID: taskID,
		Ids:    domain.Ids(residentialProxies),
		BotIds: botIds,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to assign bot to residential proxies: %v", err)
	}

	err = q.AssignBotsToProxiesForTask(ctx, dbmodel.AssignBotsToProxiesForTaskParams{
		TaskID: taskID,
		Ids:    domain.Ids(cheapProxies),
		BotIds: botIds,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to assign bot to cheap proxies: %v", err)
	}

	err = q.UpdateTaskStatus(ctx, dbmodel.UpdateTaskStatusParams{Status: dbmodel.ReadyTaskStatus, ID: taskID})
	if err != nil {
		return 0, fmt.Errorf("failed to update task status: %v", err)
	}

	return len(botIds), tx.Commit(ctx)
}

func (s *Store) deleteUnnecessaryRows(
	ctx context.Context,
	tx dbmodel.Tx,
	accounts []dbmodel.BotAccount,
	residentialProxies, cheapProxies []dbmodel.Proxy,
) ([]dbmodel.BotAccount, []dbmodel.Proxy, []dbmodel.Proxy, error) {
	q := dbmodel.New(tx)

	accountsLen, proxiesLen := len(accounts), min(len(residentialProxies), len(cheapProxies))
	logger.Infof(ctx, "got %d accounts and %d residentialProxies", accountsLen, proxiesLen)

	var remainRows = min(accountsLen, proxiesLen)

	switch {
	case accountsLen < proxiesLen:
		var err error
		// надо удалить лишние прокси из задачи
		var deletedResidentialProxiesCount, deletedCheapProxiesCount, rowsToDelete int64

		if len(residentialProxies) > accountsLen {
			residentialProxiesToDelete := len(residentialProxies) - accountsLen
			rowsToDelete += int64(residentialProxiesToDelete)
			deletedResidentialProxiesCount, err = q.DeleteProxiesForTask(ctx, proxiesLastIds(residentialProxies, residentialProxiesToDelete))
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to delete residentialProxies: %v", err)
			}
		}

		if len(cheapProxies) > accountsLen {
			cheapProxiesToDelete := len(cheapProxies) - accountsLen
			rowsToDelete += int64(cheapProxiesToDelete)
			deletedCheapProxiesCount, err = q.DeleteProxiesForTask(ctx, proxiesLastIds(cheapProxies, cheapProxiesToDelete))
			if err != nil {
				return nil, nil, nil, fmt.Errorf("failed to delete residentialProxies: %v", err)
			}
		}

		if deletedResidentialProxiesCount+deletedCheapProxiesCount != rowsToDelete {
			return nil, nil, nil, fmt.Errorf("wanted to delete %d residentialProxies, deleted %d", rowsToDelete, deletedResidentialProxiesCount)
		}

	case accountsLen == proxiesLen:
		return accounts, residentialProxies, cheapProxies, nil

	case accountsLen > proxiesLen:
		// надо удалить лишних ботов из задачи

		rowsToDelete := accountsLen - proxiesLen
		deletedRowsCount, err := q.DeleteBotAccountsForTask(ctx, accountsLastIds(accounts, rowsToDelete))
		if err != nil {
			return nil, nil, nil, fmt.Errorf("failed to delete bot accounts: %v", err)
		}

		if int(deletedRowsCount) != rowsToDelete {
			return nil, nil, nil, fmt.Errorf("wanted to delete %d bot accounts, deleted %d", rowsToDelete, deletedRowsCount)
		}
	}

	return accounts[:remainRows], residentialProxies[:remainRows], cheapProxies[:remainRows], nil
}

// accountsLastIds возвращает список из rowsToDelete последних айдишников
func accountsLastIds(arr []dbmodel.BotAccount, rowsToDelete int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, rowsToDelete)
	for _, account := range arr[len(arr)-rowsToDelete:] {
		ids = append(ids, account.ID)
	}

	return ids
}

// proxiesLastIds возвращает список из rowsToDelete последних айдишников
func proxiesLastIds(arr []dbmodel.Proxy, rowsToDelete int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, rowsToDelete)
	for _, account := range arr[len(arr)-rowsToDelete:] {
		ids = append(ids, account.ID)
	}

	return ids
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
