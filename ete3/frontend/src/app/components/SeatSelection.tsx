import { useState, useEffect } from 'react';
import { Seat } from '../types';
import useBookingStore from '../store/bookingStore';

interface SeatSelectionProps {
  seats: Seat[];
  bookedSeatIds?: number[];
}

const SeatSelection = ({ seats, bookedSeatIds = [] }: SeatSelectionProps) => {
  const { selectedSeats, toggleSeatSelection } = useBookingStore();
  
  // Create a map of all available seats from the API
  const seatsByPosition = new Map();
  seats.forEach(seat => {
    const key = `${seat.row_number}-${seat.seat_number}`;
    seatsByPosition.set(key, seat);
  });

  // Create grid of rows (A-G) and columns (1-20)
  const rows = [1, 2, 3, 4, 5, 6, 7]; // A through G
  const columns = Array.from({ length: 20 }, (_, i) => i + 1); // 1-20
  
  const isSeatSelected = (seatId: number) => {
    return selectedSeats.some(seat => seat.id === seatId);
  };
  
  const isSeatBooked = (seatId: number) => {
    return bookedSeatIds.includes(seatId);
  };
  
  const isSeatAvailable = (row: number, column: number) => {
    const key = `${row}-${column}`;
    return seatsByPosition.has(key);
  };
  
  const handleSeatClick = (seat: Seat) => {
    if (isSeatBooked(seat.id)) return;
    toggleSeatSelection(seat);
  };
  
  return (
    <div className="my-8 bg-black text-white p-6 rounded-lg">
      <h2 className="text-2xl font-bold mb-6">Select Your Seats</h2>
      
      <div className="flex justify-center mb-6">
        <div className="flex items-center space-x-8">
          <div className="flex items-center">
            <div className="w-6 h-6 bg-gray-600 rounded mr-2"></div>
            <span className="text-sm">Available</span>
          </div>
          <div className="flex items-center">
            <div className="w-6 h-6 bg-red-600 rounded mr-2"></div>
            <span className="text-sm">Selected</span>
          </div>
          <div className="flex items-center">
            <div className="w-6 h-6 bg-gray-800 rounded mr-2"></div>
            <span className="text-sm">Booked/Unavailable</span>
          </div>
        </div>
      </div>
      
      <div className="w-full overflow-x-auto">
        <div className="flex flex-col items-center">
          {/* Screen */}
          <div className="w-full max-w-4xl h-12 bg-gray-600 mb-10 flex items-center justify-center text-white font-medium">
            SCREEN
          </div>
          
          {/* Seats */}
          <div className="mb-8 max-w-4xl w-full">
            {rows.map(row => {
              const rowLabel = String.fromCharCode(64 + row); // A, B, C, etc.
              
              return (
                <div key={row} className="flex items-center mb-2">
                  <div className="w-8 flex items-center justify-center font-bold text-lg">
                    {rowLabel}
                  </div>
                  <div className="flex-1 grid grid-cols-20 gap-1">
                    {columns.map(column => {
                      const key = `${row}-${column}`;
                      const seat = seatsByPosition.get(key);
                      const isAvailable = isSeatAvailable(row, column);
                      
                      // Create a temporary seat object for unavailable seats
                      const seatObj = seat || {
                        id: (row * 100) + column,
                        row_number: row,
                        seat_number: column,
                        theater_id: 1,
                        created_at: '',
                        updated_at: ''
                      };
                      
                      const isSelected = selectedSeats.some(s => 
                        s.row_number === row && s.seat_number === column
                      );
                      const isBooked = bookedSeatIds.includes(seatObj.id);
                      
                      return (
                        <button
                          key={key}
                          onClick={() => isAvailable ? handleSeatClick(seatObj) : null}
                          disabled={!isAvailable || isBooked}
                          className={`w-full h-8 rounded-sm flex items-center justify-center text-xs 
                            ${isSelected ? 'bg-red-600 text-white' : ''}
                            ${isBooked ? 'bg-gray-800 text-gray-600 cursor-not-allowed' : ''}
                            ${!isAvailable ? 'bg-gray-800 text-gray-600 cursor-not-allowed' : ''}
                            ${!isSelected && !isBooked && isAvailable ? 'bg-gray-600 hover:bg-gray-500' : ''}
                          `}
                        >
                          {column}
                        </button>
                      );
                    })}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
      
      {selectedSeats.length > 0 && (
        <div className="mt-6 p-4 bg-gray-800 rounded">
          <p className="font-bold mb-2">Selected Seats: {selectedSeats.length}</p>
          <div className="flex flex-wrap gap-2">
            {selectedSeats.map((seat) => (
              <span key={`${seat.row_number}-${seat.seat_number}`} className="px-2 py-1 bg-red-900 text-white rounded-full text-sm">
                {String.fromCharCode(64 + seat.row_number)}{seat.seat_number}
              </span>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default SeatSelection; 