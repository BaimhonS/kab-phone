import React from 'react';
import { useQuery } from '@tanstack/react-query';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { DollarSign } from 'lucide-react';
import toast from 'react-hot-toast';

const fetchTotalIncome = async () => {
  const token = localStorage.getItem('token');
  const response = await axios.get('http://localhost:8080/api/orders/total-income', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const IncomeCard = ({ title, amount }) => (
  <Card>
    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
      <CardTitle className="text-sm font-medium">{title}</CardTitle>
      <DollarSign className="h-4 w-4 text-muted-foreground" />
    </CardHeader>
    <CardContent>
      <div className="text-2xl font-bold">฿{amount.toLocaleString()}</div>
    </CardContent>
  </Card>
);

const ShowIncome = () => {
  const navigate = useNavigate();
  const { data, isLoading, isError } = useQuery({
    queryKey: ['totalIncome'],
    queryFn: fetchTotalIncome,
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

  const { total_income_day, total_income_week, total_income_month, total_income_year } = data.data;

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <h1 className="text-2xl font-bold mb-6">Total Income</h1>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <IncomeCard title="Daily Income" amount={total_income_day} />
          <IncomeCard title="Weekly Income" amount={total_income_week} />
          <IncomeCard title="Monthly Income" amount={total_income_month} />
          <IncomeCard title="Yearly Income" amount={total_income_year} />
        </div>
      </main>
    </div>
  );
};

export default ShowIncome;
