// parcel_test.go
package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

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

	require.Equal(t, parcel.Number, packParcel.Number, "ОШИБКА NUMBER")
	require.Equal(t, parcel.Client, packParcel.Client, "КЛИЕНТЫ НЕ СОВПАДАЮТ")
	require.Equal(t, parcel.Status, packParcel.Status, "СТАТУСЫ НЕ СОВПАДАЮТ")
	require.Equal(t, parcel.Address, packParcel.Address, "АДРЕСА НЕ СОВПАДАЮТ")
	require.Equal(t, parcel.CreatedAt, packParcel.CreatedAt, "CreatedAt НЕ СОВПАДАЮТ")

	// delete
	err = store.Delete(id)
	require.NoError(t, err, "ERROR DELETE")

	// проверка, посылка удалена
	_, err = store.Get(id)
	require.Error(t, err, "ОЖИДАЕТСЯ ОШИБКА ПРИ ПОВТОРНОМ ПОЛУЧЕНИИ ПОСЫЛКИ")
	require.Equal(t, sql.ErrNoRows, err, "ОЖИДАЕМАЯ ОШИБКА: No Rows!")
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
	require.Equal(t, newAddress, updatedParcel.Address)
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
	require.Equal(t, newStatus, updatedParcel.Status)
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
	require.Equal(t, 3, len(storedParcels), "Должно быть 3 посылки")

	// check
	for _, parcel := range storedParcels {
		expected, exists := parcelMap[parcel.Number]
		require.True(t, exists, "Посылка не найдена в ожидаемых")

		require.Equal(t, expected.Client, parcel.Client)
		require.Equal(t, expected.Status, parcel.Status)
		require.Equal(t, expected.Address, parcel.Address)
		require.Equal(t, expected.CreatedAt, parcel.CreatedAt)
	}
}
