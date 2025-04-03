import axios from 'axios';

const API_URL = 'http://localhost:8080/api';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Movie API
export const getMovies = async () => {
  const response = await api.get('/cinema/movies');
  return response.data;
};

export const getMovie = async (id: number) => {
  const response = await api.get(`/cinema/movies/${id}`);
  return response.data;
};

export const createMovie = async (movieData: any) => {
  const response = await api.post('/cinema/movies', movieData);
  return response.data;
};

export const updateMovie = async (id: number, movieData: any) => {
  const response = await api.put(`/cinema/movies/${id}`, movieData);
  return response.data;
};

// Show API
export const getMovieShows = async (movieId: number) => {
  const response = await api.get(`/cinema/movies/${movieId}/shows`);
  return response.data;
};

// Seat API
export const getShowSeats = async (showId: number) => {
  try {
    const response = await api.get(`/cinema/shows/${showId}/seats`);
    return response.data;
  } catch (error) {
    console.error('Error fetching seats from API:', error);
    
    // Mock response for testing - only returning select seats
    // Remove in production and handle error properly
    const mockSeats = [];
    let seatId = 1;
    
    // Only create seats for row A, C, E and G
    for (let row of [1, 3, 5, 7]) {
      // Only create seats 1-5, 10-15 for each row
      for (let seatNum = 1; seatNum <= 20; seatNum++) {
        if ((seatNum <= 5) || (seatNum >= 10 && seatNum <= 15)) {
          mockSeats.push({
            id: seatId++,
            theater_id: 1,
            row_number: row,
            seat_number: seatNum,
            created_at: new Date().toISOString(),
            updated_at: new Date().toISOString()
          });
        }
      }
    }
    
    return mockSeats;
  }
};

// Booking API
export const createBooking = async (bookingData: { show_id: number; seat_ids: number[] }) => {
  const response = await api.post('/cinema/bookings', bookingData);
  return response.data;
};

export const getBookings = async () => {
  const response = await api.get('/cinema/bookings');
  return response.data;
};

export const cancelBooking = async (id: number) => {
  const response = await api.delete(`/cinema/bookings/${id}`);
  return response.data;
}; 