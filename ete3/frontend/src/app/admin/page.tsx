'use client';

import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { getMovies, createMovie, updateMovie } from '../services/api';
import { Movie } from '../types';
import Loading from '../components/Loading';
import { FaPlus, FaEdit, FaClock, FaTag, FaImage } from 'react-icons/fa';

interface MovieFormData {
  title: string;
  description: string;
  duration: number;
  genre: string;
  poster_url: string;
}

export default function AdminPage() {
  const [movies, setMovies] = useState<Movie[]>([]);
  const [loading, setLoading] = useState(true);
  const [formMode, setFormMode] = useState<'create' | 'edit'>('create');
  const [editingMovie, setEditingMovie] = useState<Movie | null>(null);
  const [formOpen, setFormOpen] = useState(false);
  const [submitLoading, setSubmitLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<MovieFormData>();

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

  const handleEditClick = (movie: Movie) => {
    setEditingMovie(movie);
    setFormMode('edit');
    setFormOpen(true);
    
    setValue('title', movie.title);
    setValue('description', movie.description || '');
    setValue('duration', movie.duration);
    setValue('genre', movie.genre || '');
    setValue('poster_url', movie.poster_url || '');
  };

  const handleAddNewClick = () => {
    setEditingMovie(null);
    setFormMode('create');
    setFormOpen(true);
    reset({
      title: '',
      description: '',
      duration: 0,
      genre: '',
      poster_url: '',
    });
  };

  const onSubmit = async (data: MovieFormData) => {
    try {
      setSubmitLoading(true);
      setError(null);
      setSuccess(null);

      if (formMode === 'create') {
        const newMovie = await createMovie(data);
        setMovies([...movies, newMovie]);
        setSuccess('Movie created successfully!');
      } else if (formMode === 'edit' && editingMovie) {
        const updatedMovie = await updateMovie(editingMovie.id, data);
        setMovies(movies.map(m => m.id === updatedMovie.id ? updatedMovie : m));
        setSuccess('Movie updated successfully!');
      }

      setFormOpen(false);
    } catch (err) {
      setError('Failed to save movie. Please try again later.');
      console.error(err);
    } finally {
      setSubmitLoading(false);
    }
  };

  if (loading && movies.length === 0) {
    return <Loading />;
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Movie Management</h1>
        <button
          onClick={handleAddNewClick}
          className="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-lg flex items-center transition-colors"
        >
          <FaPlus className="mr-2" />
          Add New Movie
        </button>
      </div>

      {error && (
        <div className="bg-red-100 text-red-700 p-4 rounded mb-6">
          {error}
        </div>
      )}

      {success && (
        <div className="bg-green-100 text-green-700 p-4 rounded mb-6">
          {success}
        </div>
      )}

      {/* Movie form */}
      {formOpen && (
        <div className="bg-white shadow-md rounded-lg p-6 mb-8 border border-gray-200">
          <h2 className="text-xl font-bold mb-4">
            {formMode === 'create' ? 'Add New Movie' : 'Edit Movie'}
          </h2>
          <form onSubmit={handleSubmit(onSubmit)}>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-4">
                <div>
                  <label className="block text-gray-700 mb-1">Title</label>
                  <input
                    {...register('title', { required: 'Title is required' })}
                    className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
                  />
                  {errors.title && <p className="text-red-500 text-sm mt-1">{errors.title.message}</p>}
                </div>

                <div>
                  <label className="block text-gray-700 mb-1">Duration (minutes)</label>
                  <input
                    type="number"
                    {...register('duration', { 
                      required: 'Duration is required',
                      min: { value: 1, message: 'Duration must be at least 1 minute' }
                    })}
                    className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
                  />
                  {errors.duration && <p className="text-red-500 text-sm mt-1">{errors.duration.message}</p>}
                </div>

                <div>
                  <label className="block text-gray-700 mb-1">Genre</label>
                  <input
                    {...register('genre')}
                    className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
                  />
                </div>

                <div>
                  <label className="block text-gray-700 mb-1">Poster URL</label>
                  <input
                    {...register('poster_url')}
                    className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
                  />
                </div>
              </div>

              <div>
                <label className="block text-gray-700 mb-1">Description</label>
                <textarea
                  {...register('description')}
                  rows={7}
                  className="w-full px-4 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
                ></textarea>
              </div>
            </div>

            <div className="flex mt-6 gap-4 justify-end">
              <button
                type="button"
                onClick={() => setFormOpen(false)}
                className="px-4 py-2 border border-gray-300 rounded-lg text-gray-700"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={submitLoading}
                className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg transition-colors"
              >
                {submitLoading ? 'Saving...' : 'Save Movie'}
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Movies list */}
      <div className="bg-white shadow-md rounded-lg overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Movie
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Duration
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Genre
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {movies.map((movie) => (
              <tr key={movie.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <div className="flex-shrink-0 h-10 w-10 rounded overflow-hidden bg-gray-100">
                      {movie.poster_url ? (
                        <img
                          src={movie.poster_url}
                          alt={movie.title}
                          className="h-10 w-10 object-cover"
                        />
                      ) : (
                        <div className="h-10 w-10 flex items-center justify-center text-gray-400">
                          <FaImage />
                        </div>
                      )}
                    </div>
                    <div className="ml-4">
                      <div className="text-sm font-medium text-gray-900">
                        {movie.title}
                      </div>
                    </div>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center text-sm text-gray-500">
                    <FaClock className="mr-1" />
                    <span>{movie.duration} mins</span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center text-sm text-gray-500">
                    <FaTag className="mr-1" />
                    <span>{movie.genre || 'N/A'}</span>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <button
                    onClick={() => handleEditClick(movie)}
                    className="text-indigo-600 hover:text-indigo-900 flex items-center"
                  >
                    <FaEdit className="mr-1" />
                    Edit
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
} 