package storage

type Storage interface {
	MustConnect()
	Connect() error
	MustClose()
	Close() error
}
