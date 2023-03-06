package tasks

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/domain"
	"github.com/inst-api/poster/internal/pager"
	"github.com/jackc/pgx/v4"
)

func (s *Store) TaskProgress(ctx context.Context, taskID uuid.UUID, pager *pager.Pager) (domain.TaskProgress, error) {
	q := dbmodel.New(s.dbtxf(ctx))

	task, err := q.FindTaskByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.TaskProgress{}, ErrTaskNotFound
		}
	}

	if task.Status < dbmodel.StartedTaskStatus {
		return domain.TaskProgress{}, fmt.Errorf("%w: expected at least %d status, got %d",
			ErrTaskInvalidStatus, dbmodel.StartedTaskStatus, task.Status)
	}

	query := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Select("b.username",
		"b.status",
		"b.file_order",
		"count(tu.*) filter ( where interaction_type = 'post_description' ) as post_description_targets",
		"count(tu.*) filter ( where interaction_type = 'photo_tag' ) as photo_tags_targets",
		"count(distinct medias.pk) as posts_count",
		"count(distinct medias.pk) filter ( where is_edited = true)         as edited_posts_count").
		From("bot_accounts as b").
		LeftJoin("medias on b.id = medias.bot_id").
		LeftJoin("target_users tu on medias.pk = tu.media_fk").
		GroupBy("b.id", "b."+strings.Split(pager.GetSort()[0], " ")[0]).
		Where(sq.Eq{"b.task_id": taskID}).Limit(pager.GetLimit64()).Offset(pager.GetOffset64()).OrderBy(pager.GetSort()...)

	stmt, args, err := query.ToSql()
	if err != nil {
		return domain.TaskProgress{}, fmt.Errorf("failed to build query: %v", err)
	}

	rows, err := s.db.Query(ctx, stmt, args...)
	if err != nil {
		return domain.TaskProgress{}, err
	}

	defer rows.Close()
	bots := make([]domain.BotProgress, 0)
	for rows.Next() {
		var b domain.BotProgress
		err = rows.Scan(
			&b.Username,
			&b.Status,
			&b.FileOrder,
			&b.PostDescriptionTargets,
			&b.PhotoTaggedTargets,
			&b.PostsCount,
			&b.EditedPostsCount,
		)
		if err != nil {
			return domain.TaskProgress{}, fmt.Errorf("failed to scan bot progress: %v", err)
		}
		bots = append(bots, b)
	}

	targetCounters, err := q.GetTaskTargetsCount(ctx, taskID)
	if err != nil {
		return domain.TaskProgress{}, fmt.Errorf("failed to get target counters: %v", err)
	}

	botsCount, err := q.GetTaskBotsCount(ctx, taskID)
	if err != nil {
		return domain.TaskProgress{}, fmt.Errorf("failed to get task bots count: %v", err)
	}

	progress := domain.TaskProgress{
		BotsProgress:   bots,
		TargetCounters: targetCounters,
		Done:           task.Status == dbmodel.DoneTaskStatus || task.Status == dbmodel.AllDoneTaskStatus,
		BotsCount:      int(botsCount),
	}
	return progress, rows.Err()
}
