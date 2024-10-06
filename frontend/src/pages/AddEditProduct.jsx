import React, { useState, useEffect } from 'react';
import { useNavigate, useParams, useLocation } from 'react-router-dom';
import { useMutation } from '@tanstack/react-query';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { toast } from 'sonner';
import axios from 'axios';

const AddEditProduct = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [formData, setFormData] = useState({
    brand_name: '',
    model_name: '',
    os: '',
    price: '',
    amount: '',
    image: null
  });

  const [imagePreview, setImagePreview] = useState(null);

  useEffect(() => {
    if (id && location.state?.product) {
      const { brand_name, model_name, os, price, amount, image } = location.state.product;
      setFormData({ brand_name, model_name, os, price: price.toString(), amount: amount.toString(), image: null });

      if (image) {
        setImagePreview(`data:image/jpeg;base64,${image}`);
      }
    }
  }, [id, location.state]);

  const mutation = useMutation({
    mutationFn: async (data) => {
      const token = localStorage.getItem('token');
      const formData = new FormData();
      for (const key in data) {
        if (key !== 'image' || (key === 'image' && data[key])) {
          formData.append(key, data[key]);
        }
      }
      if (id) {
        return axios.put(`http://localhost:8080/api/phones/${id}`, formData, {
          headers: { 
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'multipart/form-data'
          }
        });
      } else {
        return axios.post('http://localhost:8080/api/phones', formData, {
          headers: { 
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'multipart/form-data'
          }
        });
      }
    },
    onSuccess: () => {
      toast.success(id ? "Product updated successfully" : "Product added successfully");
      navigate('/');
    },
    onError: (error) => {
      toast.error(error.response?.data?.message || "Failed to save product");
    }
  });

  const handleInputChange = (e) => {
    const { name, value, files } = e.target;
    if (name === 'image') {
      const file = files[0];
      setFormData({ ...formData, [name]: file });
      
      const imageUrl = URL.createObjectURL(file);
      setImagePreview(imageUrl);
    } else {
      setFormData({ ...formData, [name]: value });
    }
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    mutation.mutate(formData);
  };

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <Card>
          <CardHeader>
            <CardTitle>{id ? 'Edit Product' : 'Add New Product'}</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Label htmlFor="brand_name">Brand Name</Label>
                <Input id="brand_name" name="brand_name" value={formData.brand_name} onChange={handleInputChange} required />
              </div>
              <div>
                <Label htmlFor="model_name">Model Name</Label>
                <Input id="model_name" name="model_name" value={formData.model_name} onChange={handleInputChange} required />
              </div>
              <div>
                <Label htmlFor="os">Operating System</Label>
                <Input id="os" name="os" value={formData.os} onChange={handleInputChange} required />
              </div>
              <div>
                <Label htmlFor="price">Price</Label>
                <Input id="price" name="price" type="number" step="0.01" value={formData.price} onChange={handleInputChange} required />
              </div>
              <div>
                <Label htmlFor="amount">Amount</Label>
                <Input id="amount" name="amount" type="number" value={formData.amount} onChange={handleInputChange} required />
              </div>
              <div>
                <Label htmlFor="image">Image</Label>
                <Input id="image" name="image" type="file" onChange={handleInputChange} />
              </div>

              {imagePreview && (
                <div>
                  <img src={imagePreview} alt="Preview" className="max-w-xs max-h-xs object-cover mt-4" />
                </div>
              )}

              <Button type="submit" className="w-full">
                {id ? 'Update Product' : 'Add Product'}
              </Button>
            </form>
          </CardContent>
        </Card>
      </main>
    </div>
  );
};

export default AddEditProduct;