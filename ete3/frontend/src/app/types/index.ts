export interface Movie {
  id: number;
  title: string;
  description: string;
  duration: number;
  genre: string;
  poster_url: string;
  created_at: string;
  updated_at: string;
}

export interface Show {
  id: number;
  movie_id: number;
  theater_id: number;
  start_time: string;
  end_time: string;
  price: number;
  created_at: string;
  updated_at: string;
}

export interface Seat {
  id: number;
  theater_id: number;
  row_number: number;
  seat_number: number;
  created_at: string;
  updated_at: string;
  isBooked?: boolean;
}

export interface BookingRequest {
  show_id: number;
  seat_ids: number[];
}

export interface BookingResponse {
  booking_id: number;
  status: string;
  message: string;
}

export interface Booking {
  id: number;
  show_id: number;
  seat_id: number;
  status: string;
  created_at: string;
  updated_at: string;
} 