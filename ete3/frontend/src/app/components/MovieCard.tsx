import Image from 'next/image';
import Link from 'next/link';
import { FaClock, FaTag } from 'react-icons/fa';
import { Movie } from '../types';

interface MovieCardProps {
  movie: Movie;
}

const MovieCard = ({ movie }: MovieCardProps) => {
  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden transition-transform hover:scale-105">
      <div className="relative h-80 w-full">
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
      
      <div className="p-4">
        <h3 className="text-xl font-bold mb-2 truncate">{movie.title}</h3>
        
        <div className="flex items-center mb-2 text-sm text-gray-600">
          <FaClock className="mr-1" />
          <span>{movie.duration} mins</span>
        </div>
        
        <div className="flex items-center mb-3 text-sm text-gray-600">
          <FaTag className="mr-1" />
          <span>{movie.genre}</span>
        </div>
        
        <p className="text-gray-700 text-sm line-clamp-2 mb-4">
          {movie.description || 'No description available.'}
        </p>
        
        <Link 
          href={`/movies/${movie.id}`}
          className="block w-full text-center bg-red-600 hover:bg-red-700 text-white py-2 rounded-md transition-colors"
        >
          Book Now
        </Link>
      </div>
    </div>
  );
};

export default MovieCard; 