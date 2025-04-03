package models

import (
	"encoding/json"
	"testing"
)

func TestTheaterLayoutMarshaling(t *testing.T) {
	// Create a test TheaterLayout
	layout := TheaterLayout{
		TheaterID: 1,
		Name:      "Test Theater",
		Rows:      3,
		Columns:   4,
		Layout: [][]SeatStatus{
			{
				{ID: 1, Row: 1, Column: 1, Status: "available"},
				{ID: 2, Row: 1, Column: 2, Status: "booked"},
				{ID: 3, Row: 1, Column: 3, Status: "selected"},
				{ID: 4, Row: 1, Column: 4, Status: "unavailable"},
			},
			{
				{ID: 5, Row: 2, Column: 1, Status: "available"},
				{ID: 6, Row: 2, Column: 2, Status: "available"},
				{ID: 7, Row: 2, Column: 3, Status: "booked"},
				{ID: 8, Row: 2, Column: 4, Status: "booked"},
			},
			{
				{ID: 9, Row: 3, Column: 1, Status: "selected"},
				{ID: 10, Row: 3, Column: 2, Status: "selected"},
				{ID: 11, Row: 3, Column: 3, Status: "unavailable"},
				{ID: 12, Row: 3, Column: 4, Status: "unavailable"},
			},
		},
	}

	// Marshal the layout to JSON
	data, err := json.Marshal(layout)
	if err != nil {
		t.Fatalf("Failed to marshal TheaterLayout: %v", err)
	}

	// Verify the resulting JSON contains the expected layout string
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal test result: %v", err)
	}

	expectedLayout := "ABSX|AABB|SSXX|"
	if result["layout"] != expectedLayout {
		t.Errorf("Expected layout to be %q, got %q", expectedLayout, result["layout"])
	}

	// Test unmarshaling
	var unmarshaledLayout TheaterLayout
	if err := json.Unmarshal(data, &unmarshaledLayout); err != nil {
		t.Fatalf("Failed to unmarshal TheaterLayout: %v", err)
	}

	// Verify the rows and columns are correct
	if unmarshaledLayout.Rows != 3 || unmarshaledLayout.Columns != 4 {
		t.Errorf("Expected dimensions to be 3x4, got %dx%d", unmarshaledLayout.Rows, unmarshaledLayout.Columns)
	}

	// Verify seat statuses are correctly reconstructed
	expectedStatuses := []struct {
		row    int
		col    int
		status string
	}{
		{0, 0, "available"},
		{0, 1, "booked"},
		{0, 2, "selected"},
		{0, 3, "unavailable"},
		{1, 0, "available"},
		{1, 1, "available"},
		{1, 2, "booked"},
		{1, 3, "booked"},
		{2, 0, "selected"},
		{2, 1, "selected"},
		{2, 2, "unavailable"},
		{2, 3, "unavailable"},
	}

	for _, expected := range expectedStatuses {
		actual := unmarshaledLayout.Layout[expected.row][expected.col].Status
		if actual != expected.status {
			t.Errorf("Expected status at [%d][%d] to be %q, got %q",
				expected.row, expected.col, expected.status, actual)
		}
	}
}

func TestTheaterLayoutRoundTrip(t *testing.T) {
	// Create a simple TheaterLayout
	original := TheaterLayout{
		TheaterID: 1,
		Name:      "Test Theater",
		Rows:      2,
		Columns:   2,
		Layout: [][]SeatStatus{
			{
				{ID: 1, Row: 1, Column: 1, Status: "available"},
				{ID: 2, Row: 1, Column: 2, Status: "booked"},
			},
			{
				{ID: 3, Row: 2, Column: 1, Status: "selected"},
				{ID: 4, Row: 2, Column: 2, Status: "unavailable"},
			},
		},
	}

	// Marshal to JSON and back
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal TheaterLayout: %v", err)
	}

	var reconstructed TheaterLayout
	if err := json.Unmarshal(data, &reconstructed); err != nil {
		t.Fatalf("Failed to unmarshal TheaterLayout: %v", err)
	}

	// Check core properties
	if original.TheaterID != reconstructed.TheaterID ||
		original.Name != reconstructed.Name ||
		original.Rows != reconstructed.Rows ||
		original.Columns != reconstructed.Columns {
		t.Errorf("Core properties don't match after round-trip")
	}

	// Check seat positions and statuses
	for i := 0; i < original.Rows; i++ {
		for j := 0; j < original.Columns; j++ {
			origSeat := original.Layout[i][j]
			recSeat := reconstructed.Layout[i][j]

			if origSeat.Status != recSeat.Status ||
				origSeat.Row != recSeat.Row ||
				origSeat.Column != recSeat.Column {
				t.Errorf("Seat at [%d][%d] doesn't match after round-trip: %+v vs %+v",
					i, j, origSeat, recSeat)
			}
		}
	}
}

func TestTheaterLayoutUnmarshalEdgeCases(t *testing.T) {
	testCases := []struct {
		name           string
		jsonLayout     string
		expectedError  bool
		expectedRows   int
		expectedCols   int
		expectedStatus string
	}{
		{
			name:          "Empty Layout",
			jsonLayout:    `{"theater_id": 1, "name": "Empty", "layout": ""}`,
			expectedError: false,
			expectedRows:  0,
			expectedCols:  0,
		},
		{
			name:          "Single Row",
			jsonLayout:    `{"theater_id": 1, "name": "Single", "layout": "AB|"}`,
			expectedError: false,
			expectedRows:  1,
			expectedCols:  2,
		},
		{
			name:           "First Seat Available",
			jsonLayout:     `{"theater_id": 1, "name": "Test", "layout": "A|"}`,
			expectedError:  false,
			expectedRows:   1,
			expectedCols:   1,
			expectedStatus: "available",
		},
		{
			name:           "First Seat Booked",
			jsonLayout:     `{"theater_id": 1, "name": "Test", "layout": "B|"}`,
			expectedError:  false,
			expectedRows:   1,
			expectedCols:   1,
			expectedStatus: "booked",
		},
		{
			name:          "Invalid JSON",
			jsonLayout:    `{"theater_id": 1, "name": "Invalid"`, // Missing closing brace
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var layout TheaterLayout
			err := json.Unmarshal([]byte(tc.jsonLayout), &layout)

			if tc.expectedError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if layout.Rows != tc.expectedRows {
				t.Errorf("Expected %d rows, got %d", tc.expectedRows, layout.Rows)
			}
			if layout.Columns != tc.expectedCols {
				t.Errorf("Expected %d columns, got %d", tc.expectedCols, layout.Columns)
			}

			if tc.expectedStatus != "" && layout.Rows > 0 && layout.Columns > 0 {
				if layout.Layout[0][0].Status != tc.expectedStatus {
					t.Errorf("Expected status %q, got %q", tc.expectedStatus, layout.Layout[0][0].Status)
				}
			}
		})
	}
}
