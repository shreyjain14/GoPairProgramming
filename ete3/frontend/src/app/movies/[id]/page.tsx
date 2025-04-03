'use client';

import { useEffect, useState } from 'react';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { FaClock, FaTag, FaCalendarAlt, FaTicketAlt } from 'react-icons/fa';
import { getMovie, getMovieShows, getShowSeats, createBooking } from '../../services/api';
import { Movie, Show, Seat, BookingResponse } from '../../types';
import ShowTimeSelection from '../../components/ShowTimeSelection';
import SeatSelection from '../../components/SeatSelection';
import Loading from '../../components/Loading';
import useBookingStore from '../../store/bookingStore';

export default function MovieDetails({ params }: { params: { id: string } }) {
  const router = useRouter();
  const { 
    selectedMovie, 
    selectedShow, 
    selectedSeats, 
    setSelectedMovie, 
    setSelectedShow, 
    clearSelectedSeats, 
    resetBooking 
  } = useBookingStore();

  const [movie, setMovie] = useState<Movie | null>(null);
  const [shows, setShows] = useState<Show[]>([]);
  const [seats, setSeats] = useState<Seat[]>([]);
  const [bookedSeatIds, setBookedSeatIds] = useState<number[]>([]);
  const [loading, setLoading] = useState(true);
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingResult, setBookingResult] = useState<BookingResponse | null>(null);
  const [bookingError, setBookingError] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const movieId = parseInt(params.id);

  useEffect(() => {
    // Reset booking state when component mounts
    resetBooking();

    const fetchMovieData = async () => {
      try {
        setLoading(true);
        const movieData = await getMovie(movieId);
        setMovie(movieData);
        setSelectedMovie(movieData);

        const showsData = await getMovieShows(movieId);
        setShows(showsData);
      } catch (err) {
        setError('Failed to load movie details. Please try again later.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchMovieData();
  }, [movieId, resetBooking, setSelectedMovie]);

  useEffect(() => {
    if (selectedShow) {
      const fetchSeats = async () => {
        try {
          setLoading(true);
          const seatsData = await getShowSeats(selectedShow.id);
          
          // In a real app, we'd get booked seats from the API
          // For now, let's simulate some random booked seats
          const randomBookedSeats = seatsData
            .filter(() => Math.random() > 0.7)
            .map(seat => seat.id);
          
          setSeats(seatsData);
          setBookedSeatIds(randomBookedSeats);
        } catch (err) {
          console.error(err);
          setError('Failed to load seats. Please try again later.');
        } finally {
          setLoading(false);
        }
      };

      fetchSeats();
      clearSelectedSeats();
    }
  }, [selectedShow, clearSelectedSeats]);

  const handleBooking = async () => {
    if (!selectedShow || selectedSeats.length === 0) return;

    try {
      setBookingLoading(true);
      setBookingError(null);

      const bookingData = {
        show_id: selectedShow.id,
        seat_ids: selectedSeats.map(seat => seat.id)
      };

      const result = await createBooking(bookingData);
      setBookingResult(result);

      // Navigate to booking confirmation after success
      setTimeout(() => {
        router.push('/bookings');
      }, 2000);
    } catch (err: any) {
      setBookingError(err.response?.data?.message || 'Failed to complete booking');
      console.error(err);
    } finally {
      setBookingLoading(false);
    }
  };

  if (loading && !movie) {
    return <Loading />;
  }

  if (error || !movie) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-red-100 text-red-700 p-4 rounded text-center">
          {error || 'Movie not found'}
        </div>
      </div>
    );
  }

  const totalPrice = selectedShow && selectedSeats.length > 0
    ? selectedSeats.length * selectedShow.price
    : 0;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-12">
        <div className="flex flex-col md:flex-row gap-8">
          {/* Movie poster */}
          <div className="md:w-1/3">
            <div className="relative h-96 w-full rounded-lg overflow-hidden shadow-lg">
              {movie.poster_url ? (
                <Image 
                  src={movie.poster_url}
                  alt={movie.title}
                  fill
                  className="object-cover"
                />
              ) : (
                <div className="h-full w-full flex items-center justify-center bg-gray-200">
                  <span className="text-gray-500">No poster available</span>
                </div>
              )}
            </div>
          </div>

          {/* Movie details */}
          <div className="md:w-2/3">
            <h1 className="text-3xl font-bold mb-4">{movie.title}</h1>
            
            <div className="flex flex-wrap gap-4 mb-4 text-sm">
              <div className="flex items-center gap-1">
                <FaClock className="text-red-600" />
                <span>{movie.duration} mins</span>
              </div>
              <div className="flex items-center gap-1">
                <FaTag className="text-red-600" />
                <span>{movie.genre}</span>
              </div>
            </div>
            
            <div className="mb-6">
              <h3 className="text-lg font-medium mb-2">Synopsis</h3>
              <p className="text-gray-700">
                {movie.description || 'No description available.'}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Show times selection */}
      {shows.length > 0 ? (
        <div className="mb-12">
          <ShowTimeSelection shows={shows} />
        </div>
      ) : (
        <div className="bg-yellow-100 text-yellow-800 p-4 rounded mb-12">
          No showtimes available for this movie.
        </div>
      )}

      {/* Seat selection */}
      {selectedShow && seats.length > 0 && (
        <div className="mb-8">
          <h2 className="text-xl font-bold mb-4">Select Your Seats</h2>
          <SeatSelection seats={seats} bookedSeatIds={bookedSeatIds} />
        </div>
      )}

      {/* Booking summary */}
      {selectedShow && selectedSeats.length > 0 && (
        <div className="bg-gray-100 p-6 rounded-lg mb-8">
          <h2 className="text-xl font-bold mb-4">Booking Summary</h2>
          <div className="space-y-3 mb-4">
            <div className="flex justify-between">
              <span>Movie:</span>
              <span className="font-medium">{movie.title}</span>
            </div>
            <div className="flex justify-between">
              <span>Show Time:</span>
              <span className="font-medium">
                {new Date(selectedShow.start_time).toLocaleString()}
              </span>
            </div>
            <div className="flex justify-between">
              <span>Seats:</span>
              <span className="font-medium">
                {selectedSeats.map(seat => 
                  `${String.fromCharCode(64 + seat.row_number)}${seat.seat_number}`
                ).join(', ')}
              </span>
            </div>
            <div className="flex justify-between">
              <span>Price per Ticket:</span>
              <span className="font-medium">${selectedShow.price.toFixed(2)}</span>
            </div>
            <div className="flex justify-between text-lg font-bold">
              <span>Total Price:</span>
              <span>${totalPrice.toFixed(2)}</span>
            </div>
          </div>
        </div>
      )}

      {/* Booking button */}
      <div className="flex justify-center">
        <button
          onClick={handleBooking}
          disabled={!selectedShow || selectedSeats.length === 0 || bookingLoading}
          className={`px-8 py-3 text-white font-medium rounded-lg transition
            ${!selectedShow || selectedSeats.length === 0 || bookingLoading
              ? 'bg-gray-400 cursor-not-allowed'
              : 'bg-red-600 hover:bg-red-700'
            }`
          }
        >
          {bookingLoading ? 'Processing...' : 'Confirm Booking'}
        </button>
      </div>

      {/* Booking result message */}
      {bookingResult && (
        <div className="mt-4 p-4 bg-green-100 text-green-800 rounded text-center">
          Booking successful! Your booking ID is {bookingResult.booking_id}.
          Redirecting to your bookings...
        </div>
      )}

      {/* Booking error message */}
      {bookingError && (
        <div className="mt-4 p-4 bg-red-100 text-red-700 rounded text-center">
          {bookingError}
        </div>
      )}
    </div>
  );
} 