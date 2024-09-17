// eslint-disable-next-line no-unused-vars
import React from 'react';
import { ShoppingCart, LogOut, LogIn } from 'lucide-react';
import { Button } from "@/components/ui/button";
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useToast } from "@/components/ui/use-toast";
import { Badge } from "@/components/ui/badge";

const Header = ({ cartItemCount }) => {
  const navigate = useNavigate();
  const { toast } = useToast();
  const isLoggedIn = !!localStorage.getItem('token');

  const callLogoutAPI = async () => {
    const token = localStorage.getItem('token');
    try {
      await axios.post('http://localhost:8080/api/users/logout', {}, {
        headers: { Authorization: `Bearer ${token}` }
      });
    } catch (error) {
      console.error('Logout API call failed', error);
    }
  };

  const handleLogout = async () => {
    try {
      await callLogoutAPI();
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      toast({
        title: "Success",
        description: "Logged out successfully",
      });
      navigate('/');
    } catch (error) {
      console.error('Logout failed', error);
      toast({
        title: "Error",
        description: "Logout failed. Please try again.",
        variant: "destructive",
      });
    }
  };

  const handleCartClick = () => {
    if (isLoggedIn) {
      navigate('/cart');
    } else {
      navigate('/login');
    }
  };

  const handleLoginClick = () => {
    navigate('/login');
  };

  return (
    <div className="bg-white shadow-sm">
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        <div className="flex-grow flex items-center justify-center">
          <img src="/logo.png" alt="Logo" className="h-12 w-12" />
        </div>
        <div className="flex items-center space-x-2">
          <Button variant="ghost" size="icon" onClick={handleCartClick} className="relative">
            <ShoppingCart className="h-6 w-6" />
            {isLoggedIn && cartItemCount > 0 && (
              <Badge className="absolute -top-2 -right-2 px-2 py-1 text-xs bg-red-500 text-white rounded-full">
                {cartItemCount}
              </Badge>
            )}
          </Button>
          {isLoggedIn ? (
            <Button variant="ghost" size="icon" onClick={handleLogout}>
              <LogOut className="h-6 w-6" />
            </Button>
          ) : (
            <Button variant="ghost" size="icon" onClick={handleLoginClick}>
              <LogIn className="h-6 w-6" />
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default Header;
