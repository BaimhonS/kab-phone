import React from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Minus, Plus, Trash2 } from 'lucide-react';
import toast from 'react-hot-toast';

const fetchCart = async () => {
  const token = localStorage.getItem('token');
  const response = await axios.get('http://localhost:8080/api/carts/', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const CartItem = ({ item, onUpdateQuantity, onRemove }) => (
  <Card className="mb-4">
    <CardContent className="flex items-center justify-between p-4">
      <img src={`data:image/jpeg;base64,${item.phone.image}`} alt={item.phone.model_name} className="w-20 h-20 object-cover rounded" />
      <div className="flex-grow ml-4">
        <h3 className="font-semibold">{item.phone.brand_name} {item.phone.model_name}</h3>
        <p className="text-sm text-gray-600">฿{item.phone.price.toLocaleString()}</p>
      </div>
      <div className="flex items-center">
        <Button variant="outline" size="icon" onClick={() => onUpdateQuantity(item.id, item.amount - 1)}>
          <Minus className="h-4 w-4" />
        </Button>
        <span className="mx-2">{item.amount}</span>
        <Button variant="outline" size="icon" onClick={() => onUpdateQuantity(item.id, item.amount + 1)}>
          <Plus className="h-4 w-4" />
        </Button>
        <Button variant="ghost" size="icon" className="ml-2" onClick={() => onRemove(item.id)}>
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    </CardContent>
  </Card>
);

const Cart = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data, isLoading, isError } = useQuery({
    queryKey: ['cart'],
    queryFn: fetchCart,
    onError: (error) => {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        toast.error(error.response?.data?.message || "Failed to fetch cart data");
      }
    },
  });

  const updateItemMutation = useMutation({
    mutationFn: async ({ itemId, newAmount }) => {
      const token = localStorage.getItem('token');
      const response = await axios.patch(`http://localhost:8080/api/carts/items/${itemId}`, 
        { amount: newAmount },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries(['cart']);
      toast.success(data.message);
    },
    onError: (error) => {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        const errorMessage = error.response?.data?.message || 'Failed to update cart';
        toast.error(errorMessage);
      }
    },
  });

  const removeItemMutation = useMutation({
    mutationFn: async (itemId) => {
      const token = localStorage.getItem('token');
      const response = await axios.delete(`http://localhost:8080/api/carts/items/${itemId}`, 
        { headers: { Authorization: `Bearer ${token}` } }
      );
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries(['cart']);
      toast.success(data.message);
    },
    onError: (error) => {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        const errorMessage = error.response?.data?.message || 'Failed to remove item from cart';
        toast.error(errorMessage);
      }
    },
  });

  const confirmOrderMutation = useMutation({
    mutationFn: async () => {
      const token = localStorage.getItem('token');
      const response = await axios.post('http://localhost:8080/api/orders/confirm', {}, {
        headers: { Authorization: `Bearer ${token}` }
      });
      return response.data;
    },
    onSuccess: (data) => {
      queryClient.invalidateQueries(['cart']);
      toast.success(data.message);
      navigate(`/track-order/`);
    },
    onError: (error) => {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        const errorMessage = error.response?.data?.message || 'Failed to confirm order';
        toast.error(errorMessage);
      }
    },
  });

  const handleUpdateQuantity = (itemId, newAmount) => {
    if (newAmount === 0) {
      removeItemMutation.mutate(itemId);
    } else {
      updateItemMutation.mutate({ itemId, newAmount });
    }
  };

  const handleRemoveItem = (itemId) => {
    removeItemMutation.mutate(itemId);
  };

  const handleCheckout = () => {
    confirmOrderMutation.mutate();
  };

  const calculateTotal = () => {
    if (!data || !data.data || !data.data.items) return 0;
    return data.data.items.reduce((total, item) => total + (item.amount * item.phone.price), 0);
  };

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error loading cart data</div>;

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <h1 className="text-2xl font-bold mb-6">Your Cart</h1>
        {data.data.items.map(item => (
          <CartItem 
            key={item.id} 
            item={item} 
            onUpdateQuantity={handleUpdateQuantity}
            onRemove={handleRemoveItem}
          />
        ))}
        <div className="mt-6 flex justify-between items-center">
          <h2 className="text-xl font-semibold">Total: ฿{calculateTotal().toLocaleString()}</h2>
          <Button onClick={handleCheckout} disabled={confirmOrderMutation.isLoading}>
            {confirmOrderMutation.isLoading ? 'Processing...' : 'Checkout'}
          </Button>
        </div>
      </main>
    </div>
  );
};

export default Cart;
