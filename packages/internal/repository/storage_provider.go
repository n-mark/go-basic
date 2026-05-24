package repository

type StorageProvider interface {
	Save(e Entity)
	Load(r *Repo)
	Close()
}