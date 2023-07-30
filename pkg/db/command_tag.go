package db

import (
	"strconv"
	"strings"

	"github.com/KnoblauchPilze/go-game/pkg/errors"
	"github.com/jackc/pgx"
)

var deleteCommandTag = "DELETE"
var insertCommandTag = "INSERT"

func extractAffectedRowsFromCommandTag(tag pgx.CommandTag) (int, error) {
	pieces := strings.Split(string(tag), " ")
	if len(pieces) != 2 && len(pieces) != 3 {
		return 0, errors.NewCode(errors.ErrInvalidSqlCommandTag)
	}

	switch pieces[0] {
	case deleteCommandTag:
		if len(pieces) != 2 {
			return 0, errors.NewCode(errors.ErrInvalidSqlCommandTag)
		}
		return extractRows(pieces[1])
	case insertCommandTag:
		if len(pieces) != 3 {
			return 0, errors.NewCode(errors.ErrInvalidSqlCommandTag)
		}
		return extractRows(pieces[2])
	default:
		return 0, errors.NewCode(errors.ErrUnknownSqlCommandTag)
	}
}

func extractRows(rows string) (int, error) {
	affected, err := strconv.Atoi(rows)
	if err != nil {
		return 0, errors.WrapCode(err, errors.ErrInvalidSqlCommandTag)
	}

	return affected, nil
}
