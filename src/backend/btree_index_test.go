package backend

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch_emptyIndex_returnsZeroValueAndFalse(t *testing.T) {
	// Arrange
	index := NewBtreeIndex[string]()

	// Act
	value, found := index.Search(1)

	// Assert
	assert.False(t, found)
	assert.Zero(t, value)
}

func TestInsert_newID_returnsZeroPreviousValueAndFalse(t *testing.T) {
	// Arrange
	index := NewBtreeIndex[string]()

	// Act
	previous, replaced := index.Insert(1, "Tim")

	// Assert
	assert.False(t, replaced)
	assert.Zero(t, previous)
}

func TestSearch_afterInsert_returnsValueAndTrue(t *testing.T) {
	// Arrange
	index := NewBtreeIndex[string]()
	index.Insert(1, "Tim")

	// Act
	value, found := index.Search(1)

	// Assert
	assert.True(t, found)
	assert.Equal(t, "Tim", value)
}

func TestSearch_missingIDBetweenStoredIDs_returnsZeroValueAndFalse(t *testing.T) {
	// Arrange
	index := NewBtreeIndex[string]()
	index.Insert(1, "Tim")
	index.Insert(3, "Sarah")

	// Act
	value, found := index.Search(2)

	// Assert
	assert.False(t, found)
	assert.Zero(t, value)
}

func TestInsert_existingID_replacesValueAndReturnsPrevious(t *testing.T) {
	// Arrange
	index := NewBtreeIndex[string]()
	index.Insert(1, "Tim")

	// Act
	previous, replaced := index.Insert(1, "Sarah")

	// Assert
	assert.True(t, replaced)
	assert.Equal(t, "Tim", previous)

	value, found := index.Search(1)
	require.True(t, found)
	assert.Equal(t, "Sarah", value)
}

func TestSearch_multipleIDsInsertedOutOfOrder_returnsValueForEachID(t *testing.T) {
	// Arrange
	index := NewBtreeIndex[string]()
	inserted := []struct {
		id    int
		value string
	}{
		{id: 5, value: "five"},
		{id: 1, value: "one"},
		{id: 4, value: "four"},
		{id: 2, value: "two"},
		{id: 3, value: "three"},
	}
	for _, item := range inserted {
		index.Insert(item.id, item.value)
	}

	// Act & Assert
	for _, item := range inserted {
		value, found := index.Search(item.id)
		require.True(t, found, "id %d not found", item.id)
		assert.Equal(t, item.value, value, "id %d", item.id)
	}
}

func TestSearch_edgeIDs_returnsStoredValues(t *testing.T) {
	testCases := []struct {
		name  string
		id    int
		value string
	}{
		{name: "zero", id: 0, value: "zero"},
		{name: "negative", id: -7, value: "negative"},
		{name: "max int", id: math.MaxInt, value: "max"},
		{name: "min int", id: math.MinInt, value: "min"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			index := NewBtreeIndex[string]()
			index.Insert(tc.id, tc.value)

			// Act
			value, found := index.Search(tc.id)

			// Assert
			require.True(t, found)
			assert.Equal(t, tc.value, value)
		})
	}
}

func TestSearch_manyIDsInsertedDescending_returnsValueForEachID(t *testing.T) {
	// Arrange
	const count = 100
	index := NewBtreeIndex[string]()
	for id := count - 1; id >= 0; id-- {
		index.Insert(id, fmt.Sprintf("user-%d", id))
	}

	// Act & Assert
	for id := 0; id < count; id++ {
		value, found := index.Search(id)
		require.True(t, found, "id %d not found", id)
		assert.Equal(t, fmt.Sprintf("user-%d", id), value, "id %d", id)
	}
}

func TestBTreeIndex_nonComparableValueType_storesAndRetrievesValue(t *testing.T) {
	// Arrange
	index := NewBtreeIndex[[]string]()

	// Act
	previous, replaced := index.Insert(1, []string{"Tim", "Smith"})

	// Assert
	assert.False(t, replaced)
	assert.Nil(t, previous)

	value, found := index.Search(1)
	require.True(t, found)
	assert.Equal(t, []string{"Tim", "Smith"}, value)
}
