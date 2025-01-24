package db

// Интерфейс Базы данных имплементирован структурой Postgres (for now)
type BoardRepository interface {
	CreateBoard()
	GetBoards()
}

type ThreadRepository interface {
	GetAllThreads()
	CreateThread()
}

type PostRepository interface {
	GetAllPosts()
	CreatePost()
}

type UserRepository interface {
	GetUser()
	CreateUser()
}
