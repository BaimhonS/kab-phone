import React from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import axios from 'axios';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Card, CardContent } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { toast } from "sonner";

const fetchOrders = async () => {
  const token = localStorage.getItem('token');
  const response = await axios.get('http://localhost:8080/api/orders/check-order', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const updateOrderStatus = async ({ orderId, isDelivered }) => {
  const token = localStorage.getItem('token');
  await axios.put(`http://localhost:8080/api/orders/${orderId}`, 
    { isDelivered },
    { headers: { Authorization: `Bearer ${token}` } }
  );
};

const AdminOrders = () => {
  const queryClient = useQueryClient();
  const { data, isLoading, isError } = useQuery({
    queryKey: ['adminOrders'],
    queryFn: fetchOrders,
  });

  const mutation = useMutation({
    mutationFn: updateOrderStatus,
    onSuccess: () => {
      queryClient.invalidateQueries(['adminOrders']);
      toast.success('Order status updated');
    },
    onError: () => {
      toast.error('Failed to update order status');
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error fetching orders</div>;

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <h1 className="text-2xl font-bold mb-6">Admin Orders</h1>
        <Card>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Order ID</TableHead>
                  <TableHead>Tracking Number</TableHead>
                  <TableHead>Total Price</TableHead>
                  <TableHead>Items</TableHead>
                  <TableHead>User Address</TableHead>
                  <TableHead>Delivered</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.data.map((order) => (
                  <TableRow key={order.id}>
                    <TableCell>{order.id}</TableCell>
                    <TableCell>{order.tracking_number}</TableCell>
                    <TableCell>${order.total_price}</TableCell>
                    <TableCell>
                      {order.cart.items.map((item, index) => (
                        <div key={item.id}>
                          {item.amount}x {item.phone.brand_name} {item.phone.model_name}
                          {index < order.cart.items.length - 1 && ', '}
                        </div>
                      ))}
                    </TableCell>
                    <TableCell>{order.cart.user.address}</TableCell>
                    <TableCell>
                      <Checkbox
                        checked={order.isDelivered}
                        onCheckedChange={(checked) => 
                          mutation.mutate({ orderId: order.id, isDelivered: checked })
                        }
                      />
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </main>
    </div>
  );
};

export default AdminOrders;