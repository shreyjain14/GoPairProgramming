'use client';

import { useEffect, useState } from 'react';
import { getMovies } from './services/api';
import MovieCard from './components/MovieCard';
import Loading from './components/Loading';
import { Movie } from './types';
import { FaSearch } from 'react-icons/fa';

export default function Home() {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    const fetchMovies = async () => {
      try {
        setLoading(true);
        const data = await getMovies();
        setMovies(data);
      } catch (err) {
        setError('Failed to load movies. Please try again later.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchMovies();
  }, []);

  const filteredMovies = movies.filter(movie => 
    movie.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
    movie.genre.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="container mx-auto px-4 py-8">
      <section className="mb-12">
        <div className="bg-gradient-to-r from-red-800 to-red-600 text-white py-16 px-4 rounded-lg mb-8">
          <div className="max-w-3xl mx-auto text-center">
            <h1 className="text-4xl md:text-5xl font-bold mb-4">
              Book Your Perfect Movie Experience
            </h1>
            <p className="text-lg mb-8">
              Find the latest movies and secure your seats in just a few clicks
            </p>
            <div className="relative max-w-md mx-auto">
              <input
                type="text"
                placeholder="Search for movies..."
                className="w-full px-4 py-3 pl-12 rounded-full bg-white bg-opacity-20 border border-white border-opacity-30 text-white placeholder-gray-200 focus:outline-none focus:ring-2 focus:ring-white"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
              />
              <FaSearch className="absolute left-4 top-1/2 transform -translate-y-1/2 text-white" />
            </div>
          </div>
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-bold mb-6">Available Movies</h2>
        
        {loading ? (
          <Loading />
        ) : error ? (
          <div className="bg-red-100 text-red-700 p-4 rounded text-center">{error}</div>
        ) : filteredMovies.length === 0 ? (
          <div className="text-center py-8">
            <p className="text-gray-500">No movies found matching "{searchTerm}"</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {filteredMovies.map((movie) => (
              <MovieCard key={movie.id} movie={movie} />
            ))}
          </div>
        )}
      </section>
    </div>
  );
}
