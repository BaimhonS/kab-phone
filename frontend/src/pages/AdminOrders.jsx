import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";

const API_URL = 'http://localhost:8080/api';

const fetchOrders = async () => {
  const token = localStorage.getItem('token');
  const response = await axios.get(`${API_URL}/orders/check-order`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return response.data;
};

const updateOrderStatus = async ({ orderId, isDelivered }) => {
  const token = localStorage.getItem('token');
  await axios.put(`${API_URL}/orders/${orderId}`, 
    { isDelivered },
    { headers: { Authorization: `Bearer ${token}` } }
  );
};

const addTrackingNumber = async ({ orderId, trackingNumber }) => {
  const token = localStorage.getItem('token');
  const response = await axios.post(`${API_URL}/orders/add-tracking`, 
    { order_id: orderId, tracking_number: trackingNumber },
    { headers: { Authorization: `Bearer ${token}` } }
  );
  return response.data;
};

const AdminOrders = () => {
  const queryClient = useQueryClient();
  const [trackingNumbers, setTrackingNumbers] = useState({});

  const { data, isLoading, isError } = useQuery({
    queryKey: ['adminOrders'],
    queryFn: fetchOrders,
  });

  const updateStatusMutation = useMutation({
    mutationFn: updateOrderStatus,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminOrders']);
      toast.success('Order status updated');
    },
    onError: () => {
      toast.error('Failed to update order status');
    },
  });

  const addTrackingMutation = useMutation({
    mutationFn: addTrackingNumber,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminOrders']);
      toast.success('Tracking number added successfully');
      setTrackingNumbers({});
    },
    onError: () => {
      toast.error('Failed to add tracking number');
    },
  });

  const handleTrackingNumberChange = (orderId, value) => {
    setTrackingNumbers(prev => ({ ...prev, [orderId]: value }));
  };

  const handleAddTrackingNumber = (orderId) => {
    const trackingNumber = trackingNumbers[orderId];
    if (!trackingNumber) {
      toast.error('Please enter a tracking number');
      return;
    }
    addTrackingMutation.mutate({ orderId, trackingNumber });
  };

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error fetching orders</div>;

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <h1 className="text-3xl font-bold mb-6">Admin Orders</h1>
        <Card className="overflow-hidden">
          <CardContent className="p-0">
            <div className="overflow-x-auto">
              <Table className="w-full whitespace-nowrap">
                <TableHeader>
                  <TableRow>
                    <TableHead className="px-4 py-3 text-left">Order ID</TableHead>
                    <TableHead className="px-4 py-3 text-left">First Name</TableHead>
                    <TableHead className="px-4 py-3 text-left">Last Name</TableHead>
                    <TableHead className="px-4 py-3 text-left">Tracking Number</TableHead>
                    <TableHead className="px-4 py-3 text-left">Total Price</TableHead>
                    <TableHead className="px-4 py-3 text-left">Items</TableHead>
                    <TableHead className="px-4 py-3 text-left">User Address</TableHead>
                    {/* <TableHead className="px-4 py-3 text-left">Delivered</TableHead> */}
                    <TableHead className="px-4 py-3 text-left">Add Tracking</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {data.data.map((order) => (
                    <TableRow key={order.id}>
                      <TableCell className="px-4 py-3">{order.id}</TableCell>
                      <TableCell className="px-4 py-3">{order.cart.user.first_name}</TableCell>
                      <TableCell className="px-4 py-3">{order.cart.user.last_name}</TableCell>
                      <TableCell className="px-4 py-3">{order.tracking_number || 'N/A'}</TableCell>
                      <TableCell className="px-4 py-3">฿{order.total_price.toLocaleString()}</TableCell>
                      <TableCell className="px-4 py-3">
                        {order.cart.items.map((item, index) => (
                          <div key={item.id}>
                            {item.amount}x {item.phone.brand_name} {item.phone.model_name}
                            {index < order.cart.items.length - 1 && ', '}
                          </div>
                        ))}
                      </TableCell>
                      <TableCell className="px-4 py-3">{order.cart.user.address}</TableCell>
                      {/* <TableCell className="px-4 py-3">
                        <Checkbox
                          checked={order.isDelivered}
                          onCheckedChange={(checked) =>
                            updateStatusMutation.mutate({ orderId: order.id, isDelivered: checked })
                          }
                        />
                      </TableCell> */}
                      <TableCell className="px-4 py-3">
                        <div className="flex items-center space-x-3">
                          <Input
                            placeholder="Tracking Number"
                            value={trackingNumbers[order.id] || ''}
                            onChange={(e) => handleTrackingNumberChange(order.id, e.target.value)}
                            className="w-33"
                          />
                          <Button onClick={() => handleAddTrackingNumber(order.id)} size="sm">
                            Add
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  );
};

export default AdminOrders;