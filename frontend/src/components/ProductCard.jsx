import React from 'react';
import { ShoppingCart, Edit, Package } from 'lucide-react';
import { Button } from "@/components/ui/button";
import axios from 'axios';
import toast from 'react-hot-toast';
import { Badge } from "@/components/ui/badge";

const ProductCard = ({ product, onEdit, isAdmin, refetchCartItemCount, onAddToCart, isLoggedIn }) => {
  const handleAddToCart = async () => {
    if (!isLoggedIn) {
      onAddToCart();
      return;
    }

    try {
      const token = localStorage.getItem('token');
      const response = await axios.post('http://localhost:8080/api/carts/items', 
        {
          item: {
            amount: 1,
            phone_id: product.id
          }
        },
        {
          headers: { Authorization: `Bearer ${token}` }
        }
      );
      
      toast.success(response.data.message, {
        duration: 1000,
      });
      refetchCartItemCount();
    } catch (error) {
      console.error('Failed to add item to cart', error);
      const errorMessage = error.response?.data?.message || 'Failed to add item to cart';
      toast.error(errorMessage, {
        duration: 1000,
      });
    }
  };

  const handleEdit = () => {
    onEdit(product);
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden">
      <div className="relative">
        <img 
          src={`data:image/jpeg;base64,${product.image}`} 
          alt={product.model_name} 
          className="w-full h-48 object-cover"
        />
        <Badge className="absolute top-2 right-2 bg-blue-500">
          {product.os}
        </Badge>
      </div>
      <div className="p-4">
        <h3 className="text-lg font-semibold mb-2">{product.model_name}</h3>
        <p className="text-sm text-gray-600 mb-2">{product.brand_name}</p>
        <div className="flex justify-between items-center mb-2">
          <span className="text-xl font-bold">฿{product.price.toLocaleString()}</span>
          <div className="flex items-center">
            <Package className="h-4 w-4 mr-1 text-gray-500" />
            <span className="text-sm text-gray-500">{product.amount} in stock</span>
          </div>
        </div>
        <div className="flex justify-between items-center">
          {isAdmin ? (
            <Button size="sm" onClick={handleEdit} className="w-full">
              <Edit className="h-4 w-4 mr-2" />
              Edit
            </Button>
          ) : (
            <Button size="sm" onClick={handleAddToCart} className="w-full" disabled={product.amount <= 0}>
              <ShoppingCart className="h-4 w-4 mr-2" />
              {product.amount > 0 ? 'Add to Cart' : 'Out of Stock'}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProductCard;
