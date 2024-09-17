import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const OrderDetails = ({ order }) => (
  <Card className="mb-8">
    <CardHeader>
      <CardTitle>Order Details</CardTitle>
    </CardHeader>
    <CardContent>
      <p><strong>Tracking Number:</strong> {order.tracking_number}</p>
      <p><strong>Order Status:</strong> {order.cart.status}</p>
      <p><strong>Total Price:</strong> ${order.total_price}</p>
      <p><strong>Order Date:</strong> {new Date(order.created_at).toLocaleString()}</p>
      <h3 className="font-semibold mt-4 mb-2">Items:</h3>
      <ul>
        {order.cart.items && order.cart.items.map((item) => (
          <li key={item.id} className="flex items-center mb-4">
            <img 
              src={`data:image/jpeg;base64,${item.phone.image}`} 
              alt={`${item.phone.brand_name} ${item.phone.model_name}`} 
              className="w-16 h-16 object-cover rounded-md mr-4"
            />
            <div>
              <p>{item.phone.brand_name} {item.phone.model_name}</p>
              <p>Quantity: {item.amount}</p>
              <p>Price: ฿{item.phone.price.toLocaleString()}</p>
            </div>
          </li>
        ))}
      </ul>
    </CardContent>
  </Card>
);

export default OrderDetails;