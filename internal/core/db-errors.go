package core

type DbErrType string

const (
	Resource_Not_Found       DbErrType = "resource not found in database"
	Duplicate_Value_Resource DbErrType = "duplicate resource"
	Null_Value_Resource      DbErrType = "null value in non-value field"
	Internal_Db_Server       DbErrType = "problem with the db server"
)

type DB_ERROR struct {
	Type DbErrType
}

func (dbErr DB_ERROR) Error() string {
	return string(dbErr.Type)
}
