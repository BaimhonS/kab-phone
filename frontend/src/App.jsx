import { Toaster } from 'react-hot-toast';
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { navItems } from "./nav-items";
import axios from 'axios';
import AddEditProduct from "./pages/AddEditProduct";

const queryClient = new QueryClient();

// Add axios interceptor for global error handling
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response && error.response.status === 401) {
      // Clear all session data
      localStorage.clear();
      sessionStorage.clear();
      // Clear cookies
      document.cookie.split(";").forEach((c) => {
        document.cookie = c
          .replace(/^ +/, "")
          .replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
      });
      // Redirect to login page
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

const ProtectedRoute = ({ children }) => {
  const token = localStorage.getItem('token');
  if (!token) {
    return <Navigate to="/login" replace />;
  }
  return children;
};

const AdminRoute = ({ children }) => {
  const user = JSON.parse(localStorage.getItem('user'));
  if (!user || user.role !== 'admin') {
    return <Navigate to="/" replace />;
  }
  return children;
};

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster position="top-center" reverseOrder={false} />
      <BrowserRouter>
        <Routes>
          {navItems.map(({ to, page, requiresAdmin, requiresAuth }) => (
            <Route 
              key={to} 
              path={to} 
              element={
                requiresAdmin ? (
                  <AdminRoute>{page}</AdminRoute>
                ) : requiresAuth ? (
                  <ProtectedRoute>{page}</ProtectedRoute>
                ) : (
                  page
                )
              } 
            />
          ))}
          <Route 
            path="/add-product" 
            element={
              <AdminRoute>
                <AddEditProduct />
              </AdminRoute>
            } 
          />
          <Route 
            path="/edit-product/:id" 
            element={
              <AdminRoute>
                <AddEditProduct />
              </AdminRoute>
            } 
          />
        </Routes>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
