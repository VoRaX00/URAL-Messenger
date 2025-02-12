package postgres

//type Repository struct {
//	cfg           Config
//	messengerRepo *messenger.Messenger
//	db            *sqlx.DB
//}
//
//func New(config Config) storage.Storage {
//	return &Repository{
//		cfg: config,
//	}
//}

//func (r *Repository) MustConnect() {
//	err := r.Connect()
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (r *Repository) Connect() error {
//	db, err := sqlx.Open("postgres",
//		fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
//			r.cfg.Host, r.cfg.Port, r.cfg.DBName, r.cfg.User, r.cfg.Password, r.cfg.SSLMode))
//	if err != nil {
//		return err
//	}
//	r.db = db
//
//	r.messengerRepo = messenger.New(r.db)
//	return nil
//}
//
//func (r *Repository) MustClose() {
//	err := r.Close()
//	if err != nil {
//		panic(err)
//	}
//}
//
//func (r *Repository) Close() error {
//	return r.db.Close()
//}
