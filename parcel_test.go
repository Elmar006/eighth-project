// parcel_test.go
package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

// getTestParcel
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "/mnt/d/8 sprint/final_8sprint/eighth-project/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err, "Ошибка при добавлении посылки")
	require.NotZero(t, id, "ID посылки не должен быть равен нулю")

	parcel.Number = id

	// get
	packParcel, err := store.Get(id)
	require.NoError(t, err, "Ошибка при получении посылки")

	assert.Equal(t, parcel.Number, packParcel.Number, "Ошибка Number")
	assert.Equal(t, parcel.Client, packParcel.Client, "Клиенты не совпадают")
	assert.Equal(t, parcel.Status, packParcel.Status, "Статусы не совпадают")
	assert.Equal(t, parcel.Address, packParcel.Address, "Адреса не совпадают")
	assert.Equal(t, parcel.CreatedAt, packParcel.CreatedAt, "CreatedAt не совпадают")

	// delete
	err = store.Delete(id)
	require.NoError(t, err, "ERROR DELETE")

	// проверка, посылка удалена
	_, err = store.Get(id)
	require.Error(t, err)
	require.ErrorIs(t, sql.ErrNoRows, err)
}

// TestSetAddress
func TestSetAddress(t *testing.T) {
	db, err := sql.Open("sqlite", "/mnt/d/8 sprint/final_8sprint/eighth-project/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	parcel.Number = id

	// set address
	newAddress := "new test address"
	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	// check
	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newAddress, updatedParcel.Address)
}

// TestSetStatus
func TestSetStatus(t *testing.T) {
	db, err := sql.Open("sqlite", "/mnt/d/8 sprint/final_8sprint/eighth-project/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	// set status
	newStatus := ParcelStatusSent
	err = store.SetStatus(id, newStatus)
	require.NoError(t, err)

	// check
	updatedParcel, err := store.Get(id)
	require.NoError(t, err)
	assert.Equal(t, newStatus, updatedParcel.Status)
}

// TestGetByClient
func TestGetByClient(t *testing.T) {
	db, err := sql.Open("sqlite", "/mnt/d/8 sprint/final_8sprint/eighth-project/tracker.db")
	require.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, 3, "У клиента должно быть 3 посылки")

	// check
	for _, parcel := range storedParcels {
		expected, exists := parcelMap[parcel.Number]
		//Решил тут тоже assert поменять на require так, как если вернется не та посылка, то дальнейшяя проверка поидее бесмысленна
		require.True(t, exists, "Посылка не найдена в ожидаемых")

		assert.Equal(t, expected, parcel)
	}
}
