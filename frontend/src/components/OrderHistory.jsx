import React from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from 'lucide-react';

const OrderHistory = ({ orderHistory, isLoading, isError, page, setPage, onTrackingClick }) => {
  const handleNextPage = () => {
    if (orderHistory && orderHistory.data.length === 4) {
      setPage(prevPage => prevPage + 1);
    }
  };

  const handlePrevPage = () => {
    if (page > 0) {
      setPage(prevPage => prevPage - 1);
    }
  };

  if (isLoading) return <p>Loading order history...</p>;
  if (isError) return <p className="text-red-500">Error fetching order history</p>;
  if (!orderHistory || !orderHistory.data) return null;

  return (
    <>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Tracking Number</TableHead>
            <TableHead>Total Price</TableHead>
            <TableHead>Order Date</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {orderHistory.data.map((order) => (
            <TableRow 
              key={order.id} 
              className="cursor-pointer hover:bg-gray-100" 
              onClick={() => onTrackingClick(order.tracking_number)}
            >
              <TableCell>{order.tracking_number}</TableCell>
              <TableCell>฿{order.total_price.toLocaleString()}</TableCell>
              <TableCell>{new Date(order.created_at).toLocaleString()}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
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
          disabled={orderHistory.data.length < 4}
          variant="outline"
          size="icon"
        >
          <ChevronRight className="h-4 w-4" />
        </Button>
      </div>
    </>
  );
};

export default OrderHistory;