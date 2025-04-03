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

  // Find min/max row numbers from the available seats
  const rowNumbers = seats.map(seat => seat.row_number);
  const minRow = Math.min(...rowNumbers);
  const maxRow = Math.max(...rowNumbers);
  
  // Find min/max seat numbers for each row
  const seatRangesByRow = new Map();
  seats.forEach(seat => {
    if (!seatRangesByRow.has(seat.row_number)) {
      seatRangesByRow.set(seat.row_number, { min: seat.seat_number, max: seat.seat_number });
    } else {
      const range = seatRangesByRow.get(seat.row_number);
      if (seat.seat_number < range.min) range.min = seat.seat_number;
      if (seat.seat_number > range.max) range.max = seat.seat_number;
    }
  });
  
  // Create arrays for rows and columns
  const rows = Array.from({ length: maxRow - minRow + 1 }, (_, i) => minRow + i);
  
  const isSeatSelected = (seatId: number) => {
    return selectedSeats.some(seat => seat.id === seatId);
  };
  
  const isSeatBooked = (seatId: number) => {
    return bookedSeatIds.includes(seatId);
  };
  
  const handleSeatClick = (seat: Seat) => {
    if (isSeatBooked(seat.id)) return;
    toggleSeatSelection(seat);
  };
  
  return (
    <div className="my-8">
      <div className="flex justify-center mb-6">
        <div className="flex items-center space-x-6">
          <div className="flex items-center">
            <div className="w-6 h-6 bg-gray-200 border border-gray-300 rounded mr-2"></div>
            <span className="text-sm">Available</span>
          </div>
          <div className="flex items-center">
            <div className="w-6 h-6 bg-red-500 rounded mr-2"></div>
            <span className="text-sm">Selected</span>
          </div>
          <div className="flex items-center">
            <div className="w-6 h-6 bg-gray-500 rounded mr-2"></div>
            <span className="text-sm">Booked</span>
          </div>
        </div>
      </div>
      
      <div className="w-full overflow-x-auto">
        <div className="flex flex-col items-center">
          {/* Screen */}
          <div className="w-4/5 h-10 bg-gray-300 rounded-t-lg mb-10 flex items-center justify-center text-gray-700">
            SCREEN
          </div>
          
          {/* Seats */}
          <div className="mb-8">
            {rows.map(row => {
              const rowLabel = String.fromCharCode(64 + row); // A, B, C, etc.
              const range = seatRangesByRow.get(row) || { min: 1, max: 20 };
              const columns = Array.from(
                { length: range.max - range.min + 1 }, 
                (_, i) => range.min + i
              );
              
              return (
                <div key={row} className="flex justify-center mb-2">
                  <div className="flex-shrink-0 w-8 flex items-center justify-center font-bold">
                    {rowLabel}
                  </div>
                  <div className="flex gap-1 flex-wrap">
                    {columns.map(column => {
                      const key = `${row}-${column}`;
                      const seat = seatsByPosition.get(key);
                      
                      // Skip rendering if seat doesn't exist from API
                      if (!seat) return null;
                      
                      const isSelected = isSeatSelected(seat.id);
                      const isBooked = isSeatBooked(seat.id);
                      
                      return (
                        <button
                          key={key}
                          onClick={() => handleSeatClick(seat)}
                          disabled={isBooked}
                          className={`w-7 h-7 rounded-t flex items-center justify-center text-xs 
                            ${isSelected ? 'bg-red-500 text-white' : ''}
                            ${isBooked ? 'bg-gray-500 text-white cursor-not-allowed' : ''}
                            ${!isSelected && !isBooked ? 'bg-gray-200 hover:bg-gray-300' : ''}
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
      
      <div className="mt-6 p-4 bg-gray-100 rounded">
        <p className="font-bold mb-2">Selected Seats: {selectedSeats.length}</p>
        <div className="flex flex-wrap gap-2">
          {selectedSeats.map((seat) => (
            <span key={seat.id} className="px-2 py-1 bg-red-100 text-red-800 rounded-full text-sm">
              {String.fromCharCode(64 + seat.row_number)}{seat.seat_number}
            </span>
          ))}
        </div>
      </div>
    </div>
  );
};

export default SeatSelection; 