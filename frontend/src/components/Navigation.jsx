import React from 'react';
import { Link, useLocation } from 'react-router-dom';

const Navigation = () => {
  const location = useLocation();
  const user = JSON.parse(localStorage.getItem('user'));
  const isAdmin = user && user.role === 'admin';

  return (
    <nav className="container mx-auto px-4 py-2">
      <ul className="flex space-x-4">
        <li>
            <Link to="/profile" className={`${location.pathname === '/profile' ? 'text-blue-600' : 'text-gray-600'}`}>
              Profile
            </Link>
        </li>
        <li>
          <Link to="/" className={`${location.pathname === '/' ? 'text-blue-600' : 'text-gray-600'}`}>
            Products
          </Link>
        </li>
        {isAdmin && (
          <li>
            <Link to="/show-income" className={`${location.pathname === '/show-income' ? 'text-blue-600' : 'text-gray-600'}`}>
              Show Income
            </Link>
          </li>
        )}
        <li>
          <Link to="/worst-best-selling" className={`${location.pathname === '/worst-best-selling' ? 'text-blue-600' : 'text-gray-600'}`}>
            The Worst-Best Product
          </Link>
        </li>
        <li>
          <Link to="/track-order" className={`${location.pathname === '/track-order' ? 'text-blue-600' : 'text-gray-600'}`}>
            Track Order
          </Link>
        </li>
      </ul>
    </nav>
  );
};

export default Navigation;