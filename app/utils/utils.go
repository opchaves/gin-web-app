package utils

import (
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/opchaves/gin-web-app/app/model/apperrors"
)

func ConvertToUUID(id string) (*pgtype.UUID, error) {
	uuid := pgtype.UUID{}
	err := uuid.Scan(id)
	if err != nil {
		err = apperrors.NewBadRequest(apperrors.InvalidId)
		return nil, err
	}

	return &uuid, nil
}

func UUIDtoString(uuid pgtype.UUID) string {
	uBytes, _ := uuid.MarshalJSON()
	id := string(uBytes[:])

	return strings.Replace(id, "\"", "", -1)
}
