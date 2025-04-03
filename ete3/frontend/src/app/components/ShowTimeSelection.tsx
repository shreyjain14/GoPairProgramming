import { useState } from 'react';
import { format } from 'date-fns';
import { Show } from '../types';
import useBookingStore from '../store/bookingStore';

interface ShowTimeSelectionProps {
  shows: Show[];
}

const ShowTimeSelection = ({ shows }: ShowTimeSelectionProps) => {
  const { selectedShow, setSelectedShow } = useBookingStore();

  // Group shows by date
  const showsByDate = shows.reduce((acc, show) => {
    const date = new Date(show.start_time).toDateString();
    if (!acc[date]) {
      acc[date] = [];
    }
    acc[date].push(show);
    return acc;
  }, {} as Record<string, Show[]>);

  // Sort dates
  const dates = Object.keys(showsByDate).sort(
    (a, b) => new Date(a).getTime() - new Date(b).getTime()
  );

  // Active date (default to first date)
  const [activeDate, setActiveDate] = useState(dates[0] || '');

  const handleShowSelect = (show: Show) => {
    setSelectedShow(show);
  };

  if (shows.length === 0) {
    return <p className="text-center py-4">No show times available</p>;
  }

  return (
    <div className="mt-6">
      <h2 className="text-xl font-bold mb-4">Select Show Time</h2>

      {/* Date selection tabs */}
      <div className="flex overflow-x-auto space-x-2 pb-2 mb-4">
        {dates.map((date) => (
          <button
            key={date}
            onClick={() => setActiveDate(date)}
            className={`px-4 py-2 rounded-full whitespace-nowrap text-sm font-medium
              ${
                activeDate === date
                  ? 'bg-red-600 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }
            `}
          >
            {format(new Date(date), 'EEE, MMM d')}
          </button>
        ))}
      </div>

      {/* Show times */}
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3">
        {activeDate &&
          showsByDate[activeDate].map((show) => (
            <button
              key={show.id}
              onClick={() => handleShowSelect(show)}
              className={`border rounded-md p-3 text-center transition
                ${
                  selectedShow?.id === show.id
                    ? 'border-red-500 bg-red-50'
                    : 'border-gray-300 hover:border-red-300 hover:bg-red-50'
                }
              `}
            >
              <span className="block font-medium">
                {format(new Date(show.start_time), 'h:mm a')}
              </span>
              <span className="text-xs text-gray-500 mt-1 block">
                ${show.price.toFixed(2)}
              </span>
            </button>
          ))}
      </div>
    </div>
  );
};

export default ShowTimeSelection; 