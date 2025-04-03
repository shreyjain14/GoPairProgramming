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
  const response = await api.get(`/cinema/shows/${showId}/seats`);
  return response.data;
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