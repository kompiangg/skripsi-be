package sqlx

type Tx interface {
	Commit() error
	Rollback() error
}
