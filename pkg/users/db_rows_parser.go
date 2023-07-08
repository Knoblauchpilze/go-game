package users

import (
	"github.com/KnoblauchPilze/go-game/pkg/db"
	"github.com/google/uuid"
)

type userRowParser struct {
	user User
}

func (p *userRowParser) Parse(row db.Scannable) error {
	return row.Scan(&p.user.Id, &p.user.Mail, &p.user.Name, &p.user.Password, &p.user.CreatedAt)
}

type userIdsParser struct {
	ids []uuid.UUID
}

func (p *userIdsParser) Parse(row db.Scannable) error {
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return err
	}

	p.ids = append(p.ids, id)
	return nil
}
