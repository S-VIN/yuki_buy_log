package stores

import (
	"errors"
	"testing"
	"time"
	"yuki_buy_log/mocks"
	"yuki_buy_log/models"

	"go.uber.org/mock/gomock"
)

func TestNewInviteStore(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("успешное создание store с инвайтами", func(t *testing.T) {
		invites := []models.Invite{
			{
				Id:         1,
				FromUserId: 1,
				ToUserId:   2,
				FromLogin:  "user1",
				ToLogin:    "user2",
				CreatedAt:  time.Now(),
			},
			{
				Id:         2,
				FromUserId: 2,
				ToUserId:   3,
				FromLogin:  "user2",
				ToLogin:    "user3",
				CreatedAt:  time.Now(),
			},
		}

		mockDB.EXPECT().GetAllInvites().Return(invites, nil)

		store, err := NewInviteStore(mockDB)
		if err != nil {
			t.Fatalf("ожидали успешное создание store, получили ошибку: %v", err)
		}

		if store == nil {
			t.Fatal("store не должен быть nil")
		}

		if len(store.data) != 2 {
			t.Errorf("ожидали 2 инвайта в store, получили %d", len(store.data))
		}

		if store.data[0].Id != 1 {
			t.Errorf("ожидали ID 1, получили %d", store.data[0].Id)
		}
	})

	t.Run("создание store при ошибке получения инвайтов", func(t *testing.T) {
		mockDB.EXPECT().GetAllInvites().Return(nil, errors.New("database error"))

		store, err := NewInviteStore(mockDB)
		if err != nil {
			t.Fatalf("не ожидали ошибку при создании store: %v", err)
		}

		if store == nil {
			t.Fatal("store не должен быть nil")
		}

		if len(store.data) != 0 {
			t.Errorf("ожидали пустой store, получили %d инвайтов", len(store.data))
		}
	})
}

func TestGetInviteById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	now := time.Now()
	invites := []models.Invite{
		{
			Id:         1,
			FromUserId: 1,
			ToUserId:   2,
			FromLogin:  "user1",
			ToLogin:    "user2",
			CreatedAt:  now,
		},
		{
			Id:         2,
			FromUserId: 2,
			ToUserId:   3,
			FromLogin:  "user2",
			ToLogin:    "user3",
			CreatedAt:  now,
		},
	}

	mockDB.EXPECT().GetAllInvites().Return(invites, nil)
	store, _ := NewInviteStore(mockDB)

	t.Run("получение существующего инвайта", func(t *testing.T) {
		invite := store.GetInviteById(1)
		if invite == nil {
			t.Fatal("ожидали найти инвайт, получили nil")
		}

		if invite.FromUserId != 1 {
			t.Errorf("ожидали FromUserId 1, получили %d", invite.FromUserId)
		}

		if invite.ToUserId != 2 {
			t.Errorf("ожидали ToUserId 2, получили %d", invite.ToUserId)
		}
	})

	t.Run("получение несуществующего инвайта", func(t *testing.T) {
		invite := store.GetInviteById(999)
		if invite != nil {
			t.Error("ожидали nil для несуществующего инвайта")
		}
	})
}

func TestGetInvitesFromUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	now := time.Now()
	invites := []models.Invite{
		{
			Id:         1,
			FromUserId: 1,
			ToUserId:   2,
			FromLogin:  "user1",
			ToLogin:    "user2",
			CreatedAt:  now,
		},
		{
			Id:         2,
			FromUserId: 1,
			ToUserId:   3,
			FromLogin:  "user1",
			ToLogin:    "user3",
			CreatedAt:  now,
		},
		{
			Id:         3,
			FromUserId: 2,
			ToUserId:   3,
			FromLogin:  "user2",
			ToLogin:    "user3",
			CreatedAt:  now,
		},
	}

	mockDB.EXPECT().GetAllInvites().Return(invites, nil)
	store, _ := NewInviteStore(mockDB)

	t.Run("получение инвайтов от пользователя с инвайтами", func(t *testing.T) {
		result := store.GetInvitesFromUser(1)
		if len(result) != 2 {
			t.Errorf("ожидали 2 инвайта, получили %d", len(result))
		}

		for _, inv := range result {
			if inv.FromUserId != 1 {
				t.Errorf("все инвайты должны быть от пользователя 1, получили от %d", inv.FromUserId)
			}
		}
	})

	t.Run("получение инвайтов от пользователя без инвайтов", func(t *testing.T) {
		result := store.GetInvitesFromUser(999)
		if len(result) != 0 {
			t.Errorf("ожидали пустой список, получили %d инвайтов", len(result))
		}
	})
}

func TestGetInvitesToUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	now := time.Now()
	invites := []models.Invite{
		{
			Id:         1,
			FromUserId: 1,
			ToUserId:   2,
			FromLogin:  "user1",
			ToLogin:    "user2",
			CreatedAt:  now,
		},
		{
			Id:         2,
			FromUserId: 3,
			ToUserId:   2,
			FromLogin:  "user3",
			ToLogin:    "user2",
			CreatedAt:  now,
		},
		{
			Id:         3,
			FromUserId: 1,
			ToUserId:   3,
			FromLogin:  "user1",
			ToLogin:    "user3",
			CreatedAt:  now,
		},
	}

	mockDB.EXPECT().GetAllInvites().Return(invites, nil)
	store, _ := NewInviteStore(mockDB)

	t.Run("получение инвайтов к пользователю с инвайтами", func(t *testing.T) {
		result := store.GetInvitesToUser(2)
		if len(result) != 2 {
			t.Errorf("ожидали 2 инвайта, получили %d", len(result))
		}

		for _, inv := range result {
			if inv.ToUserId != 2 {
				t.Errorf("все инвайты должны быть к пользователю 2, получили к %d", inv.ToUserId)
			}
		}
	})

	t.Run("получение инвайтов к пользователю без инвайтов", func(t *testing.T) {
		result := store.GetInvitesToUser(999)
		if len(result) != 0 {
			t.Errorf("ожидали пустой список, получили %d инвайтов", len(result))
		}
	})
}

func TestGetInvite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	now := time.Now()
	invites := []models.Invite{
		{
			Id:         1,
			FromUserId: 1,
			ToUserId:   2,
			FromLogin:  "user1",
			ToLogin:    "user2",
			CreatedAt:  now,
		},
		{
			Id:         2,
			FromUserId: 2,
			ToUserId:   3,
			FromLogin:  "user2",
			ToLogin:    "user3",
			CreatedAt:  now,
		},
	}

	mockDB.EXPECT().GetAllInvites().Return(invites, nil)
	store, _ := NewInviteStore(mockDB)

	t.Run("получение существующего инвайта между пользователями", func(t *testing.T) {
		invite := store.GetInvite(1, 2)
		if invite == nil {
			t.Fatal("ожидали найти инвайт, получили nil")
		}

		if invite.FromUserId != 1 || invite.ToUserId != 2 {
			t.Errorf("ожидали инвайт от 1 к 2, получили от %d к %d", invite.FromUserId, invite.ToUserId)
		}
	})

	t.Run("получение несуществующего инвайта", func(t *testing.T) {
		invite := store.GetInvite(1, 3)
		if invite != nil {
			t.Error("ожидали nil для несуществующего инвайта")
		}
	})

	t.Run("направление имеет значение", func(t *testing.T) {
		// Инвайт от 1 к 2 существует
		invite1 := store.GetInvite(1, 2)
		if invite1 == nil {
			t.Error("ожидали найти инвайт от 1 к 2")
		}

		// Инвайт от 2 к 1 не существует
		invite2 := store.GetInvite(2, 1)
		if invite2 != nil {
			t.Error("не ожидали найти инвайт от 2 к 1")
		}
	})
}

func TestAddInvite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("добавление инвайта в пустой store", func(t *testing.T) {
		mockDB.EXPECT().GetAllInvites().Return([]models.Invite{}, nil)
		store, _ := NewInviteStore(mockDB)

		newInvite := models.Invite{
			FromUserId: 1,
			ToUserId:   2,
			FromLogin:  "user1",
			ToLogin:    "user2",
			CreatedAt:  time.Now(),
		}

		err := store.AddInvite(newInvite)
		if err != nil {
			t.Fatalf("не ожидали ошибку при добавлении инвайта: %v", err)
		}

		if len(store.data) != 1 {
			t.Errorf("ожидали 1 инвайт в store, получили %d", len(store.data))
		}

		// Проверяем, что ID был установлен
		if store.data[0].Id != 1 {
			t.Errorf("ожидали ID 1, получили %d", store.data[0].Id)
		}
	})

	t.Run("добавление инвайта в непустой store", func(t *testing.T) {
		existingInvites := []models.Invite{
			{
				Id:         5,
				FromUserId: 1,
				ToUserId:   2,
				FromLogin:  "user1",
				ToLogin:    "user2",
				CreatedAt:  time.Now(),
			},
		}

		mockDB.EXPECT().GetAllInvites().Return(existingInvites, nil)
		store, _ := NewInviteStore(mockDB)

		newInvite := models.Invite{
			FromUserId: 2,
			ToUserId:   3,
			FromLogin:  "user2",
			ToLogin:    "user3",
			CreatedAt:  time.Now(),
		}

		err := store.AddInvite(newInvite)
		if err != nil {
			t.Fatalf("не ожидали ошибку при добавлении инвайта: %v", err)
		}

		if len(store.data) != 2 {
			t.Errorf("ожидали 2 инвайта в store, получили %d", len(store.data))
		}

		// Проверяем, что новый ID больше максимального существующего
		if store.data[1].Id != 6 {
			t.Errorf("ожидали ID 6, получили %d", store.data[1].Id)
		}
	})
}

func TestDeleteInvites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("успешное удаление инвайтов между пользователями", func(t *testing.T) {
		invites := []models.Invite{
			{
				Id:         1,
				FromUserId: 1,
				ToUserId:   2,
				FromLogin:  "user1",
				ToLogin:    "user2",
				CreatedAt:  time.Now(),
			},
			{
				Id:         2,
				FromUserId: 2,
				ToUserId:   1,
				FromLogin:  "user2",
				ToLogin:    "user1",
				CreatedAt:  time.Now(),
			},
			{
				Id:         3,
				FromUserId: 1,
				ToUserId:   3,
				FromLogin:  "user1",
				ToLogin:    "user3",
				CreatedAt:  time.Now(),
			},
		}

		mockDB.EXPECT().GetAllInvites().Return(invites, nil)
		store, _ := NewInviteStore(mockDB)

		mockDB.EXPECT().DeleteInvitesBetweenUsers(models.UserId(1), models.UserId(2)).Return(nil)

		err := store.DeleteInvites(1, 2)
		if err != nil {
			t.Fatalf("не ожидали ошибку при удалении инвайтов: %v", err)
		}

		// Должен остаться только один инвайт (от 1 к 3)
		if len(store.data) != 1 {
			t.Errorf("ожидали 1 инвайт в store, получили %d", len(store.data))
		}

		if store.data[0].Id != 3 {
			t.Errorf("должен остаться инвайт с ID 3, получили %d", store.data[0].Id)
		}
	})

	t.Run("ошибка при удалении инвайтов из БД", func(t *testing.T) {
		invites := []models.Invite{
			{
				Id:         1,
				FromUserId: 1,
				ToUserId:   2,
				FromLogin:  "user1",
				ToLogin:    "user2",
				CreatedAt:  time.Now(),
			},
		}

		mockDB.EXPECT().GetAllInvites().Return(invites, nil)
		store, _ := NewInviteStore(mockDB)

		mockDB.EXPECT().DeleteInvitesBetweenUsers(models.UserId(1), models.UserId(2)).Return(errors.New("database error"))

		err := store.DeleteInvites(1, 2)
		if err == nil {
			t.Error("ожидали ошибку при удалении инвайтов")
		}

		// Инвайты не должны быть удалены из store при ошибке БД
		if len(store.data) != 1 {
			t.Errorf("инвайты не должны быть удалены из store при ошибке БД")
		}
	})
}

func TestDeleteOldInvites(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)

	t.Run("успешное удаление старых инвайтов", func(t *testing.T) {
		now := time.Now()
		oldTime := now.Add(-24 * time.Hour)
		cutoffTime := now.Add(-12 * time.Hour)

		invites := []models.Invite{
			{
				Id:         1,
				FromUserId: 1,
				ToUserId:   2,
				FromLogin:  "user1",
				ToLogin:    "user2",
				CreatedAt:  oldTime, // Старый инвайт
			},
			{
				Id:         2,
				FromUserId: 2,
				ToUserId:   3,
				FromLogin:  "user2",
				ToLogin:    "user3",
				CreatedAt:  now, // Новый инвайт
			},
		}

		mockDB.EXPECT().GetAllInvites().Return(invites, nil)
		store, _ := NewInviteStore(mockDB)

		mockDB.EXPECT().DeleteOldInvites(cutoffTime).Return(int64(1), nil)

		rowsAffected, err := store.DeleteOldInvites(cutoffTime)
		if err != nil {
			t.Fatalf("не ожидали ошибку при удалении старых инвайтов: %v", err)
		}

		if rowsAffected != 1 {
			t.Errorf("ожидали 1 удаленную строку, получили %d", rowsAffected)
		}

		// Должен остаться только новый инвайт
		if len(store.data) != 1 {
			t.Errorf("ожидали 1 инвайт в store, получили %d", len(store.data))
		}

		if store.data[0].Id != 2 {
			t.Errorf("должен остаться инвайт с ID 2, получили %d", store.data[0].Id)
		}
	})

	t.Run("удаление инвайтов с временем равным cutoff", func(t *testing.T) {
		now := time.Now()
		cutoffTime := now

		invites := []models.Invite{
			{
				Id:         1,
				FromUserId: 1,
				ToUserId:   2,
				FromLogin:  "user1",
				ToLogin:    "user2",
				CreatedAt:  cutoffTime, // Время равно cutoff
			},
		}

		mockDB.EXPECT().GetAllInvites().Return(invites, nil)
		store, _ := NewInviteStore(mockDB)

		mockDB.EXPECT().DeleteOldInvites(cutoffTime).Return(int64(0), nil)

		_, err := store.DeleteOldInvites(cutoffTime)
		if err != nil {
			t.Fatalf("не ожидали ошибку: %v", err)
		}

		// Инвайт с временем равным cutoff должен остаться
		if len(store.data) != 1 {
			t.Errorf("инвайт с временем равным cutoff должен остаться")
		}
	})

	t.Run("ошибка при удалении старых инвайтов из БД", func(t *testing.T) {
		now := time.Now()
		cutoffTime := now.Add(-12 * time.Hour)

		invites := []models.Invite{
			{
				Id:         1,
				FromUserId: 1,
				ToUserId:   2,
				FromLogin:  "user1",
				ToLogin:    "user2",
				CreatedAt:  now.Add(-24 * time.Hour),
			},
		}

		mockDB.EXPECT().GetAllInvites().Return(invites, nil)
		store, _ := NewInviteStore(mockDB)

		mockDB.EXPECT().DeleteOldInvites(cutoffTime).Return(int64(0), errors.New("database error"))

		_, err := store.DeleteOldInvites(cutoffTime)
		if err == nil {
			t.Error("ожидали ошибку при удалении старых инвайтов")
		}

		// Инвайты не должны быть удалены из store при ошибке БД
		if len(store.data) != 1 {
			t.Errorf("инвайты не должны быть удалены из store при ошибке БД")
		}
	})
}

func TestInviteStoreConcurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockIDataBaseManager(ctrl)
	invites := []models.Invite{
		{
			Id:         1,
			FromUserId: 1,
			ToUserId:   2,
			FromLogin:  "user1",
			ToLogin:    "user2",
			CreatedAt:  time.Now(),
		},
	}
	mockDB.EXPECT().GetAllInvites().Return(invites, nil)
	store, _ := NewInviteStore(mockDB)

	t.Run("конкурентное чтение безопасно", func(t *testing.T) {
		done := make(chan bool)

		for i := 0; i < 10; i++ {
			go func() {
				invite := store.GetInviteById(1)
				if invite == nil {
					t.Error("инвайт должен существовать")
				}
				done <- true
			}()
		}

		for i := 0; i < 10; i++ {
			<-done
		}
	})
}
