package data

import (
	"database/sql"

	"smecalculus/rolevod/lib/id"
)

func NullStringFromID(id id.ADT) sql.NullString {
	if id.IsEmpty() {
		return sql.NullString{}
	}
	return sql.NullString{String: id.String(), Valid: true}
}

func NullStringFromID2(id *id.ADT) sql.NullString {
	if id == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: id.String(), Valid: true}
}

func NullStringToID(dto sql.NullString) (id.ADT, error) {
	if dto.Valid {
		return id.StringToID(dto.String)
	}
	return id.Empty(), nil
}
