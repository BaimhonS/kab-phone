import React, { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import Header from '@/components/Header';
import SearchBar from '@/components/SearchBar';
import Navigation from '@/components/Navigation';
import ProductGrid from '@/components/ProductGrid';
import { Skeleton } from "@/components/ui/skeleton"
import { Button } from "@/components/ui/button"
import { ChevronLeft, ChevronRight, Plus } from 'lucide-react';
import ErrorMessage from '@/components/ErrorMessage';
import toast from 'react-hot-toast';

const fetchPhones = async (page, search) => {
  const response = await axios.get(`http://localhost:8080/api/phones?page=${page}&page_size=8${search ? `&search=${search}` : ''}`);
  return response.data;
};

const fetchCartItemCount = async () => {
  const token = localStorage.getItem('token');
  if (!token) return 0;
  try {
    const response = await axios.get('http://localhost:8080/api/carts/', {
      headers: { Authorization: `Bearer ${token}` }
    });
    return response.data.data.items.reduce((total, item) => total + item.amount, 0);
  } catch (error) {
    console.error('Error fetching cart item count:', error);
    return 0;
  }
};

const Index = () => {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const navigate = useNavigate();
  const [isAdmin, setIsAdmin] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ['phones', page, search],
    queryFn: () => fetchPhones(page, search),
    onError: (error) => {
      toast.error(error.response?.data?.message || 'An error occurred while fetching data');
    },
  });

  const { data: cartItemCount, refetch: refetchCartItemCount } = useQuery({
    queryKey: ['cartItemCount'],
    queryFn: fetchCartItemCount,
    enabled: isLoggedIn,
    onError: (error) => {
      console.error('Error fetching cart item count:', error);
    },
  });

  useEffect(() => {
    const token = localStorage.getItem('token');
    setIsLoggedIn(!!token);
    const user = JSON.parse(localStorage.getItem('user'));
    setIsAdmin(user && user.role === 'admin');
  }, []);

  const handleNextPage = () => {
    if (data && data.data.length === 8) {
      setPage(prevPage => prevPage + 1);
    }
  };

  const handlePrevPage = () => {
    if (page > 0) {
      setPage(prevPage => prevPage - 1);
    }
  };

  const handleSearch = (searchTerm) => {
    setSearch(searchTerm);
    setPage(0);
  };

  const handleAddProduct = () => {
    navigate('/add-product');
  };

  const handleEditProduct = (product) => {
    navigate(`/edit-product/${product.id}`, { state: { product } });
  };

  const handleAddToCart = () => {
    if (!isLoggedIn) {
      navigate('/login');
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header cartItemCount={cartItemCount} />
      <SearchBar onSearch={handleSearch} />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <div className="flex justify-between items-center mb-6">
          <h1 className="text-2xl font-bold">Products</h1>
          {isAdmin && (
            <Button onClick={handleAddProduct}>
              <Plus className="h-4 w-4 mr-2" />
              Add Product
            </Button>
          )}
        </div>
        <ErrorMessage message={isError ? (error.response?.data?.message || 'An error occurred while fetching data') : null} />
        {isLoading ? (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {[...Array(8)].map((_, index) => (
              <Skeleton key={index} className="h-[300px] w-full" />
            ))}
          </div>
        ) : isError ? (
          <p className="text-center text-red-500">Error loading products</p>
        ) : (
          <>
            <ProductGrid 
              products={data.data} 
              onEdit={handleEditProduct} 
              isAdmin={isAdmin} 
              refetchCartItemCount={refetchCartItemCount}
              onAddToCart={handleAddToCart}
              isLoggedIn={isLoggedIn}
            />
            <div className="flex justify-end mt-6 space-x-2">
              <Button 
                onClick={handlePrevPage} 
                disabled={page === 0}
                variant="outline"
                size="icon"
              >
                <ChevronLeft className="h-4 w-4" />
              </Button>
              <Button 
                onClick={handleNextPage} 
                disabled={data.data.length < 8}
                variant="outline"
                size="icon"
              >
                <ChevronRight className="h-4 w-4" />
              </Button>
            </div>
          </>
        )}
      </main>
    </div>
  );
};

export default Index;
