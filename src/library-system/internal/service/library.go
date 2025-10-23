package service

import (
	"context"
	"errors"
	"lab2-rsoi/library-system/internal/dto"
	//"lab2-rsoi/library-system/internal/models"
	"lab2-rsoi/library-system/internal/repo"

	"github.com/google/uuid"
)

var (
	CountOfBooksIsZero = errors.New("count of books is zero")
)

type LibraryServiceIface interface {
	ListLibraries(ctx context.Context, city string, page, size int) (*dto.LibraryPaginationResponse, error)
	ListBooks(ctx context.Context, libraryUID uuid.UUID, showAll bool, page, size int) (*dto.BookPaginationResponse, error)
	GetBookByUID(ctx context.Context, uid uuid.UUID) (*dto.BookResponse, error)
	GetLibraryByUID(ctx context.Context, uid uuid.UUID) (*dto.LibraryResponse, error)
	UpdateBookCount(ctx context.Context, bookUID, libraryUID uuid.UUID, inc int) error
	UpdateBookCondition(ctx context.Context, bookUID uuid.UUID, condition string) error
}

type LibraryService struct {
	repo repo.LibraryRepository
}

func NewLibraryService(r repo.LibraryRepository) LibraryServiceIface {
	return &LibraryService{repo: r}
}

func (s *LibraryService) ListLibraries(ctx context.Context, city string, page, size int) (*dto.LibraryPaginationResponse, error) {
	libraries, err := s.repo.FetchLibrariesByCity(ctx, city, page, size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountLibrariesByCity(ctx, city)
	if err != nil {
		return nil, err
	}

	items := make([]dto.LibraryResponse, len(libraries))
	for i, l := range libraries {
		items[i] = dto.LibraryResponse{
			LibraryUID: l.LibraryUID,
			Name:       l.Name,
			City:       l.City,
			Address:    l.Address,
		}
	}

	resp := &dto.LibraryPaginationResponse{
		Page:          page,
		PageSize:      len(items),
		TotalElements: total,
		Items:         items,
	}

	return resp, nil
}

func (s *LibraryService) ListBooks(ctx context.Context, libraryUID uuid.UUID, showAll bool, page, size int) (*dto.BookPaginationResponse, error) {
	books, err := s.repo.FetchBooksByLibrary(ctx, libraryUID, showAll, page, size)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountBooksByLibrary(ctx, libraryUID, showAll)
	if err != nil {
		return nil, err
	}

	items := make([]dto.BookResponse, len(books))
	for i, b := range books {
		items[i] = dto.BookResponse{
			BookUID:        b.BookUID,
			Name:           b.Name,
			Author:         b.Author,
			Genre:          b.Genre,
			Condition:      b.Condition,
			AvailableCount: b.AvailableCount,
		}
	}

	resp := &dto.BookPaginationResponse{
		Page:          page,
		PageSize:      len(items),
		TotalElements: total,
		Items:         items,
	}

	return resp, nil
}

func (s *LibraryService) GetLibraryByUID(ctx context.Context, uid uuid.UUID) (*dto.LibraryResponse, error) {
	lib, err := s.repo.GetLibraryByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := &dto.LibraryResponse{
		ID:         lib.ID,
		LibraryUID: lib.LibraryUID,
		Name:       lib.Name,
		City:       lib.City,
		Address:    lib.Address,
	}

	return resp, nil
}

func (s *LibraryService) GetBookByUID(ctx context.Context, uid uuid.UUID) (*dto.BookResponse, error) {
	book, err := s.repo.GetBookByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	resp := &dto.BookResponse{
		ID:             book.ID,
		BookUID:        book.BookUID,
		Name:           book.Name,
		Author:         book.Author,
		Genre:          book.Genre,
		Condition:      book.Condition,
		AvailableCount: book.AvailableCount,
	}

	return resp, nil
}

func (s *LibraryService) UpdateBookCondition(ctx context.Context, bookUID uuid.UUID, condition string) error {
	return s.repo.UpdateCondition(ctx, bookUID, condition)
}

func (s *LibraryService) UpdateBookCount(ctx context.Context, bookUID, libraryUID uuid.UUID, inc int) error {
	book, err := s.repo.GetBookByUID(ctx, bookUID)
	if err != nil {
		return err
	}

	newCount := book.AvailableCount + inc
	if newCount < 0 {
		newCount = 0
		return CountOfBooksIsZero
	}

	return s.repo.UpdateCount(ctx, bookUID, libraryUID, newCount)
}
