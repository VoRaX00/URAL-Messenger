package redis

import "github.com/go-redis/redis"

type Repository struct {
	cfg Config
	db  *redis.Client
}

func New(cfg Config) *Repository {
	return &Repository{
		cfg: cfg,
	}
}

func (r *Repository) MustConnect() {
	err := r.Connect()
	if err != nil {
		panic(err)
	}
}

func (r *Repository) Connect() error {
	db := redis.NewClient(&redis.Options{
		Addr:         r.cfg.Addr,
		Password:     r.cfg.Password,
		DB:           r.cfg.DB,
		MaxRetries:   r.cfg.MaxRetries,
		DialTimeout:  r.cfg.DialTimeout,
		ReadTimeout:  r.cfg.Timeout,
		WriteTimeout: r.cfg.Timeout,
	})

	if err := db.Ping().Err(); err != nil {
		return err
	}

	r.db = db
	return nil
}

func (r *Repository) MustClose() {
	err := r.db.Close()
	if err != nil {
		panic(err)
	}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
