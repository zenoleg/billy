package rating

import (
	"sync"

	"emperror.dev/errors"
)

type (
	MemeStorage interface {
		Get(id string) (Meme, error)
		Save(memes ...Meme) error
	}

	InMemoryMemeStorage struct {
		memes map[string]Meme
		mx    sync.RWMutex
	}
)

func NewInMemoryMemeStorage(memes map[string]Meme) MemeStorage {
	return &InMemoryMemeStorage{memes: memes, mx: sync.RWMutex{}}
}

func (i *InMemoryMemeStorage) Get(id string) (Meme, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()

	meme, ok := i.memes[id]
	if !ok {
		return Meme{}, errors.New("Meme not found")
	}

	return meme, nil
}

func (i *InMemoryMemeStorage) Save(memes ...Meme) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	for _, meme := range memes {
		i.memes[meme.id] = meme
	}

	return nil
}
