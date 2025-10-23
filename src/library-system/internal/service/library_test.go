package service_test

import (
	"context"
	"lab2-rsoi/library-system/internal/models"
	"lab2-rsoi/library-system/internal/repo"
	"lab2-rsoi/library-system/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLibraryRepo struct {
	mock.Mock
}

func (m *MockLibraryRepo) FetchLibrariesByCity(ctx context.Context, city string, page, size int) ([]models.Library, error) {
	args := m.Called(ctx, city, page, size)
	return args.Get(0).([]models.Library), args.Error(1)
}

func (m *MockLibraryRepo) CountLibrariesByCity(ctx context.Context, city string) (int, error) {
	args := m.Called(ctx, city)
	return args.Int(0), args.Error(1)
}

func (m *MockLibraryRepo) FetchBooksByLibrary(ctx context.Context, libraryUID uuid.UUID, showAll bool, page, size int) ([]repo.BookWithCount, error) {
	args := m.Called(ctx, libraryUID, showAll, page, size)
	return args.Get(0).([]repo.BookWithCount), args.Error(1)
}

func (m *MockLibraryRepo) CountBooksByLibrary(ctx context.Context, libraryUID uuid.UUID, showAll bool) (int, error) {
	args := m.Called(ctx, libraryUID, showAll)
	return args.Int(0), args.Error(1)
}

func (m *MockLibraryRepo) GetLibraryByUID(ctx context.Context, uid uuid.UUID) (*models.Library, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*models.Library), args.Error(1)
}

func (m *MockLibraryRepo) GetBookByUID(ctx context.Context, uid uuid.UUID) (*repo.BookWithCount, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(*repo.BookWithCount), args.Error(1)
}

func (m *MockLibraryRepo) UpdateCondition(ctx context.Context, bookUID uuid.UUID, condition string) error {
	args := m.Called(ctx, bookUID, condition)
	return args.Error(0)
}

func (m *MockLibraryRepo) UpdateCount(ctx context.Context, bookUID, libraryUID uuid.UUID, newCount int) error {
	args := m.Called(ctx, bookUID, libraryUID, newCount)
	return args.Error(0)
}

func TestListBooks(t *testing.T) {
	mockRepo := new(MockLibraryRepo)
	svc := service.NewLibraryService(mockRepo)

	libUID := uuid.New()
	author := "Автор"
	genre := "Жанр"
	books := []repo.BookWithCount{
		{Book: models.Book{BookUID: uuid.New(), Name: "Книга 1", Author: &author, Genre: &genre, Condition: "NEW"}, AvailableCount: 5},
	}

	mockRepo.On("FetchBooksByLibrary", mock.Anything, libUID, true, 1, 10).Return(books, nil)
	mockRepo.On("CountBooksByLibrary", mock.Anything, libUID, true).Return(1, nil)

	resp, err := svc.ListBooks(context.Background(), libUID, true, 1, 10)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.TotalElements)
	assert.Equal(t, "Книга 1", resp.Items[0].Name)
	mockRepo.AssertExpectations(t)
}
func TestListLibraries(t *testing.T) {
	mockRepo := new(MockLibraryRepo)
	svc := service.NewLibraryService(mockRepo)

	city := "Москва"
	page, size := 1, 10
	libraries := []models.Library{
		{LibraryUID: uuid.New(), Name: "Библиотека 1", City: city, Address: "Адрес 1"},
	}

	mockRepo.On("FetchLibrariesByCity", mock.Anything, city, page, size).Return(libraries, nil)
	mockRepo.On("CountLibrariesByCity", mock.Anything, city).Return(1, nil)

	resp, err := svc.ListLibraries(context.Background(), city, page, size)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.TotalElements)
	assert.Equal(t, city, resp.Items[0].City)
	mockRepo.AssertExpectations(t)
}

func TestGetLibraryByUID(t *testing.T) {
	mockRepo := new(MockLibraryRepo)
	svc := service.NewLibraryService(mockRepo)

	libUID := uuid.New()
	lib := &models.Library{LibraryUID: libUID, Name: "Библиотека", City: "Москва", Address: "Адрес"}
	mockRepo.On("GetLibraryByUID", mock.Anything, libUID).Return(lib, nil)

	resp, err := svc.GetLibraryByUID(context.Background(), libUID)
	assert.NoError(t, err)
	assert.Equal(t, libUID, resp.LibraryUID)
	assert.Equal(t, "Библиотека", resp.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetBookByUID(t *testing.T) {
	mockRepo := new(MockLibraryRepo)
	svc := service.NewLibraryService(mockRepo)

	bookUID := uuid.New()
	author := "Автор"
	genre := "Жанр"

	book := &repo.BookWithCount{Book: models.Book{BookUID: bookUID, Name: "Книга", Author: &author, Genre: &genre, Condition: "EXCELLENT"}, AvailableCount: 3}
	mockRepo.On("GetBookByUID", mock.Anything, bookUID).Return(book, nil)

	resp, err := svc.GetBookByUID(context.Background(), bookUID)
	assert.NoError(t, err)
	assert.Equal(t, bookUID, resp.BookUID)
	assert.Equal(t, "Книга", resp.Name)
	mockRepo.AssertExpectations(t)
}
