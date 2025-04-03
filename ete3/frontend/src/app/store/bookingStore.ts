import { create } from 'zustand';
import { Movie, Show, Seat } from '../types';

interface BookingState {
  selectedMovie: Movie | null;
  selectedShow: Show | null;
  selectedSeats: Seat[];
  setSelectedMovie: (movie: Movie | null) => void;
  setSelectedShow: (show: Show | null) => void;
  toggleSeatSelection: (seat: Seat) => void;
  clearSelectedSeats: () => void;
  resetBooking: () => void;
}

const useBookingStore = create<BookingState>((set) => ({
  selectedMovie: null,
  selectedShow: null,
  selectedSeats: [],
  
  setSelectedMovie: (movie) => set({ selectedMovie: movie }),
  
  setSelectedShow: (show) => set({ selectedShow: show }),
  
  toggleSeatSelection: (seat) => 
    set((state) => {
      const isSelected = state.selectedSeats.some((s) => s.id === seat.id);
      
      if (isSelected) {
        return {
          selectedSeats: state.selectedSeats.filter((s) => s.id !== seat.id),
        };
      } else {
        return {
          selectedSeats: [...state.selectedSeats, seat],
        };
      }
    }),
    
  clearSelectedSeats: () => set({ selectedSeats: [] }),
  
  resetBooking: () => set({
    selectedMovie: null,
    selectedShow: null,
    selectedSeats: [],
  }),
}));

export default useBookingStore; 