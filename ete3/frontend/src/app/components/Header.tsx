import Link from 'next/link';
import { FaFilm, FaTicketAlt, FaUser } from 'react-icons/fa';

const Header = () => {
  return (
    <header className="bg-black text-white shadow-md">
      <div className="container mx-auto px-4 py-3">
        <div className="flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold flex items-center gap-2">
            <FaFilm className="text-red-500" />
            <span>CineTicket</span>
          </Link>
          
          <nav>
            <ul className="flex space-x-6">
              <li>
                <Link 
                  href="/" 
                  className="flex items-center gap-1 hover:text-red-400 transition"
                >
                  <FaFilm />
                  <span>Movies</span>
                </Link>
              </li>
              <li>
                <Link 
                  href="/bookings" 
                  className="flex items-center gap-1 hover:text-red-400 transition"
                >
                  <FaTicketAlt />
                  <span>My Bookings</span>
                </Link>
              </li>
            </ul>
          </nav>
        </div>
      </div>
    </header>
  );
};

export default Header;
