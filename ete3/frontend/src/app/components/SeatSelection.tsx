import { useState, useEffect } from 'react';
import { Seat } from '../types';
import useBookingStore from '../store/bookingStore';

interface SeatSelectionProps {
  seats: Seat[];
  bookedSeatIds?: number[];
}

const SeatSelection = ({ seats, bookedSeatIds = [] }: SeatSelectionProps) => {
  const { selectedSeats, toggleSeatSelection } = useBookingStore();
  
  // Group seats by row for better display
  const seatsByRow = seats.reduce((acc, seat) => {
    const row = seat.row_number;
    if (!acc[row]) {
      acc[row] = [];
    }
    acc[row].push(seat);
    return acc;
  }, {} as Record<number, Seat[]>);
  
  // Sort seats in each row by seat number
  Object.keys(seatsByRow).forEach((row) => {
    seatsByRow[Number(row)].sort((a, b) => a.seat_number - b.seat_number);
  });
  
  // Get sorted row numbers
  const rows = Object.keys(seatsByRow)
    .map(Number)
    .sort((a, b) => a - b);
  
  const isSeatSelected = (seatId: number) => {
    return selectedSeats.some((seat) => seat.id === seatId);
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
            {rows.map((row) => (
              <div key={row} className="flex justify-center mb-2">
                <div className="flex-shrink-0 w-8 flex items-center justify-center font-bold">
                  {String.fromCharCode(64 + row)}
                </div>
                <div className="flex gap-2">
                  {seatsByRow[row].map((seat) => {
                    const isSelected = isSeatSelected(seat.id);
                    const isBooked = isSeatBooked(seat.id);
                    
                    return (
                      <button
                        key={seat.id}
                        onClick={() => handleSeatClick(seat)}
                        disabled={isBooked}
                        className={`w-8 h-8 rounded-t flex items-center justify-center text-xs 
                          ${isSelected ? 'bg-red-500 text-white' : ''}
                          ${isBooked ? 'bg-gray-500 text-white cursor-not-allowed' : ''}
                          ${!isSelected && !isBooked ? 'bg-gray-200 hover:bg-gray-300' : ''}
                        `}
                      >
                        {seat.seat_number}
                      </button>
                    );
                  })}
                </div>
              </div>
            ))}
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