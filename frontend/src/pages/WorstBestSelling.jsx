import React from 'react';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ThumbsDown, ThumbsUp } from 'lucide-react';
import toast from 'react-hot-toast';

const fetchWorstBestSelling = async () => {
  const token = localStorage.getItem('token');
  const response = await axios.get('http://localhost:8080/api/orders/best-worst-phones', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const ProductList = ({ title, icon, data }) => (
  <Card className="w-full">
    <CardHeader>
      <CardTitle className="flex items-center">
        {icon}
        <span className="ml-2">{title}</span>
      </CardTitle>
    </CardHeader>
    <CardContent>
      {['day', 'week', 'month', 'year'].map((period) => (
        <div key={period} className="mb-2">
          <strong className="capitalize">{period}ly:</strong>{' '}
          {data[`${title.toLowerCase()}_phone_${period}`]?.brand_name
            ? `${data[`${title.toLowerCase()}_phone_${period}`].brand_name} ${data[`${title.toLowerCase()}_phone_${period}`].model_name}`
            : 'N/A'}
        </div>
      ))}
    </CardContent>
  </Card>
);

const WorstBestSelling = () => {
  const navigate = useNavigate();
  const { data, isLoading, isError } = useQuery({
    queryKey: ['worstBestSelling'],
    queryFn: fetchWorstBestSelling,
    onError: (error) => {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        toast.error(error.response?.data?.message || 'An error occurred while fetching data');
      }
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error loading data</div>;

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <h1 className="text-2xl font-bold mb-6">Worst-Best Selling Products</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <ProductList title="Worst" icon={<ThumbsDown className="h-6 w-6 text-red-500" />} data={data.data} />
          <ProductList title="Best" icon={<ThumbsUp className="h-6 w-6 text-green-500" />} data={data.data} />
        </div>
      </main>
    </div>
  );
};

export default WorstBestSelling;
