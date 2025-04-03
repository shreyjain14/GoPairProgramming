import { FaHeart } from 'react-icons/fa';

const Footer = () => {
  return (
    <footer className="bg-gray-900 text-white py-6 mt-auto">
      <div className="container mx-auto px-4">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <p className="text-sm mb-4 md:mb-0">
            &copy; {new Date().getFullYear()} CineTicket. All rights reserved.
          </p>
          
          <div className="flex items-center text-sm">
            <span>Made with</span>
            <FaHeart className="text-red-500 mx-1" />
            <span>for movie lovers</span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer; 