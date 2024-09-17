import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useQuery } from '@tanstack/react-query';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search } from 'lucide-react';
import OrderDetails from '@/components/OrderDetails';
import OrderHistory from '@/components/OrderHistory';
import toast from 'react-hot-toast';

const fetchTrackingInfo = async (trackingNumber) => {
  if (!trackingNumber) return null;
  const token = localStorage.getItem('token');
  const response = await axios.get(`http://localhost:8080/api/orders/track-orders/${trackingNumber}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const fetchOrderHistory = async (page) => {
  const token = localStorage.getItem('token');
  const response = await axios.get(`http://localhost:8080/api/orders/track-orders?page=${page}&page_size=4`, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const TrackOrder = () => {
  const [trackingNumber, setTrackingNumber] = useState('');
  const [page, setPage] = useState(0);
  const navigate = useNavigate();

  const { data: orderHistory, isLoading: isHistoryLoading, isError: isHistoryError, refetch: refetchHistory } = useQuery({
    queryKey: ['orderHistory', page],
    queryFn: () => fetchOrderHistory(page),
    onError: handleError,
  });

  const { data, isLoading, isError, refetch } = useQuery({
    queryKey: ['trackOrder', trackingNumber],
    queryFn: () => fetchTrackingInfo(trackingNumber),
    enabled: false,
    onError: handleError,
  });

  function handleError(error) {
    if (error.response && error.response.status === 401) {
      // This part is now handled by the global axios interceptor
      return;
    }
    toast.error(error.response?.data?.message || "Failed to fetch data");
  }

  const handleSearch = () => {
    if (trackingNumber) {
      toast.promise(
        refetch(),
        {
          loading: 'Searching...',
          success: (data) => {
            if (data && data.data) {
              return 'Order found';
            } else {
              throw new Error('Order not found');
            }
          },
          error: (err) => {
            return err.response?.data?.message || 'Failed to find order';
          }
        }
      );
    }
  };

  const handleTrackingClick = async (number) => {
    setTrackingNumber(number);
    try {
      await toast.promise(
        fetchTrackingInfo(number),
        {
          loading: 'Fetching order details...',
          success: (data) => {
            if (data && data.data) {
              return 'Order details loaded';
            } else {
              throw new Error('Order not found');
            }
          },
          error: (err) => {
            return err.response?.data?.message || 'Failed to load order details';
          }
        }
      );
      refetch();
    } catch (error) {
      handleError(error);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <h1 className="text-2xl font-bold mb-6">Track Your Order</h1>
        <div className="flex mb-4">
          <Input
            type="text"
            placeholder="Enter tracking number"
            value={trackingNumber}
            onChange={(e) => setTrackingNumber(e.target.value)}
            className="flex-grow mr-2"
          />
          <Button onClick={handleSearch}>
            <Search className="h-4 w-4 mr-2" />
            Search
          </Button>
        </div>
        {isLoading && <p>Loading...</p>}
        {isError && <p className="text-red-500">Error fetching tracking information</p>}
        {data && data.data && <OrderDetails order={data.data} />}

        <h2 className="text-xl font-bold mb-4 mt-8">Order History</h2>
        <OrderHistory
          orderHistory={orderHistory}
          isLoading={isHistoryLoading}
          isError={isHistoryError}
          page={page}
          setPage={setPage}
          onTrackingClick={handleTrackingClick}
        />
      </main>
    </div>
  );
};

export default TrackOrder;
