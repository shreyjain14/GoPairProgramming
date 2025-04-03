'use client';

import { useEffect, useState } from 'react';
import { getBookings, cancelBooking } from '../services/api';
import { Booking } from '../types';
import Loading from '../components/Loading';
import { FaTicketAlt, FaCalendarAlt, FaClock, FaTrash } from 'react-icons/fa';
import { format } from 'date-fns';

export default function BookingsPage() {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [cancellingId, setCancellingId] = useState<number | null>(null);

  useEffect(() => {
    const fetchBookings = async () => {
      try {
        setLoading(true);
        const data = await getBookings();
        setBookings(data);
      } catch (err) {
        setError('Failed to load bookings. Please try again later.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchBookings();
  }, []);

  const handleCancelBooking = async (bookingId: number) => {
    try {
      setCancellingId(bookingId);
      await cancelBooking(bookingId);
      setBookings(bookings.filter(booking => booking.id !== bookingId));
    } catch (err) {
      console.error(err);
      setError('Failed to cancel booking. Please try again later.');
    } finally {
      setCancellingId(null);
    }
  };

  if (loading) {
    return <Loading />;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">My Bookings</h1>

      {error && (
        <div className="bg-red-100 text-red-700 p-4 rounded mb-6">
          {error}
        </div>
      )}

      {bookings.length === 0 ? (
        <div className="text-center py-12 bg-gray-50 rounded-lg">
          <FaTicketAlt className="mx-auto text-4xl text-gray-400 mb-4" />
          <h2 className="text-xl font-medium mb-2">No Bookings Found</h2>
          <p className="text-gray-600 mb-6">
            You don't have any movie bookings yet.
          </p>
          <a
            href="/"
            className="inline-block bg-red-600 hover:bg-red-700 text-white py-2 px-6 rounded transition-colors"
          >
            Browse Movies
          </a>
        </div>
      ) : (
        <div className="grid gap-6">
          {bookings.map((booking) => (
            <div
              key={booking.id}
              className="bg-white rounded-lg shadow-md overflow-hidden border border-gray-200"
            >
              <div className="p-6">
                <div className="flex justify-between items-start mb-4">
                  <div className="flex items-center">
                    <FaTicketAlt className="text-red-600 mr-2" />
                    <h2 className="text-xl font-bold">Booking #{booking.id}</h2>
                  </div>
                  <span
                    className={`inline-block px-3 py-1 rounded-full text-xs font-medium ${
                      booking.status === 'confirmed'
                        ? 'bg-green-100 text-green-800'
                        : 'bg-yellow-100 text-yellow-800'
                    }`}
                  >
                    {booking.status.charAt(0).toUpperCase() + booking.status.slice(1)}
                  </span>
                </div>

                <div className="space-y-2 mb-6">
                  <div className="flex items-start">
                    <FaCalendarAlt className="text-gray-500 mt-1 mr-2" />
                    <div>
                      <p className="text-sm text-gray-500">Booking Date</p>
                      <p className="font-medium">
                        {format(new Date(booking.created_at), 'MMMM d, yyyy')}
                      </p>
                    </div>
                  </div>
                  
                  {/* In a real app, we would fetch and show movie and show details here */}
                  <div className="flex items-start">
                    <FaClock className="text-gray-500 mt-1 mr-2" />
                    <div>
                      <p className="text-sm text-gray-500">Show ID</p>
                      <p className="font-medium">{booking.show_id}</p>
                    </div>
                  </div>
                  
                  <div className="flex items-start">
                    <FaTicketAlt className="text-gray-500 mt-1 mr-2" />
                    <div>
                      <p className="text-sm text-gray-500">Seat ID</p>
                      <p className="font-medium">{booking.seat_id}</p>
                    </div>
                  </div>
                </div>

                <button
                  onClick={() => handleCancelBooking(booking.id)}
                  disabled={cancellingId === booking.id}
                  className="flex items-center text-red-600 hover:text-red-800 transition-colors"
                >
                  <FaTrash className="mr-1" />
                  {cancellingId === booking.id ? 'Cancelling...' : 'Cancel Booking'}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
} 