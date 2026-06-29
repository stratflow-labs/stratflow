package db

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenderMappingMatchesProtoValues(t *testing.T) {
	require.Equal(t, sql.NullString{}, userGenderIntToDB(0))
	require.Equal(t, sql.NullString{String: "male", Valid: true}, userGenderIntToDB(1))
	require.Equal(t, sql.NullString{String: "female", Valid: true}, userGenderIntToDB(2))

	require.Equal(t, 0, genderDBToInt(sql.NullString{}))
	require.Equal(t, 1, genderDBToInt(sql.NullString{String: "male", Valid: true}))
	require.Equal(t, 2, genderDBToInt(sql.NullString{String: "female", Valid: true}))
}
