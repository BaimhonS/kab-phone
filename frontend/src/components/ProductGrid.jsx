import React from 'react';
import ProductCard from './ProductCard';

const ProductGrid = ({ products, onEdit, isAdmin, refetchCartItemCount, onAddToCart, isLoggedIn }) => {
  return (
    <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
      {products.map((product) => (
        <ProductCard 
          key={product.id} 
          product={product} 
          onEdit={onEdit} 
          isAdmin={isAdmin} 
          refetchCartItemCount={refetchCartItemCount}
          onAddToCart={onAddToCart}
          isLoggedIn={isLoggedIn}
        />
      ))}
    </div>
  );
};

export default ProductGrid;
