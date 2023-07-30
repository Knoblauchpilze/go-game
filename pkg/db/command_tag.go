package db

import (
	"strconv"
	"strings"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/jackc/pgx"
)

func extractAffectedRowsFromCommandTag(tag pgx.CommandTag) (int, error) {
	pieces := strings.Split(string(tag), " ")
	if len(pieces) != 3 {
		return 0, errors.NewCode(errors.ErrInvalidSqlCommandTag)
	}

	affected, err := strconv.Atoi(pieces[2])
	if err != nil {
		return 0, errors.WrapCode(err, errors.ErrInvalidSqlCommandTag)
	}

	return affected, nil
}
