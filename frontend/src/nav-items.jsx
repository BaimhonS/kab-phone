import { HomeIcon, DollarSign, TrendingUpIcon, ShoppingCart, Package, UserIcon } from "lucide-react";
import Index from "./pages/Index.jsx";
import Register from "./pages/Register.jsx";
import Login from "./pages/Login.jsx";
import WorstBestSelling from "./pages/WorstBestSelling.jsx";
import ShowIncome from "./pages/ShowIncome.jsx";
import Cart from "./pages/Cart.jsx";
import TrackOrder from "./pages/TrackOrder.jsx";
import Profile from "./pages/Profile.jsx";
import AddEditProduct from "./pages/AddEditProduct.jsx";
import AdminOrders from "./pages/AdminOrders.jsx";

export const navItems = [
  {
    title: "Profile",
    to: "/profile",
    icon: <UserIcon className="h-4 w-4" />,
    page: <Profile />,
    requiresAuth: true,
  },
  {
    title: "Products",
    to: "/",
    icon: <HomeIcon className="h-4 w-4" />,
    page: <Index />,
  },
  {
    title: "Register",
    to: "/register",
    icon: <HomeIcon className="h-4 w-4" />,
    page: <Register />,
  },
  {
    title: "Login",
    to: "/login",
    icon: <HomeIcon className="h-4 w-4" />,
    page: <Login />,
  },
  {
    title: "Worst-Best Selling",
    to: "/worst-best-selling",
    icon: <TrendingUpIcon className="h-4 w-4" />,
    page: <WorstBestSelling />,
    requiresAuth: true,
  },
  {
    title: "Show Income",
    to: "/show-income",
    icon: <DollarSign className="h-4 w-4" />,
    page: <ShowIncome />,
    requiresAdmin: true,
  },
  {
    title: "Cart",
    to: "/cart",
    icon: <ShoppingCart className="h-4 w-4" />,
    page: <Cart />,
    requiresAuth: true,
  },
  {
    title: "Track Order",
    to: "/track-order",
    icon: <Package className="h-4 w-4" />,
    page: <TrackOrder />,
    requiresAuth: true,
  },
  {
    title: "Add/Edit Product",
    to: "/add-product",
    icon: <Package className="h-4 w-4" />,
    page: <AddEditProduct />,
    requiresAdmin: true,
  },
  {
    title: "Check Order",
    to: "/check-order",
    icon: <Package className="h-4 w-4" />,
    page: <AdminOrders />,
    requiresAdmin: true,
  },
];
