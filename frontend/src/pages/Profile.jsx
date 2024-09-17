import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import Header from '@/components/Header';
import Navigation from '@/components/Navigation';
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import toast from 'react-hot-toast';

const fetchProfile = async () => {
  const token = localStorage.getItem('token');
  const response = await axios.get('http://localhost:8080/api/users/profile', {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const updateProfile = async (profileData) => {
  const token = localStorage.getItem('token');
  const response = await axios.put('http://localhost:8080/api/users', profileData, {
    headers: { Authorization: `Bearer ${token}` }
  });
  return response.data;
};

const Profile = () => {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);
  const [profileData, setProfileData] = useState({});

  const { data, isLoading, isError } = useQuery({
    queryKey: ['profile'],
    queryFn: fetchProfile,
    onError: (error) => {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        toast.error(error.response?.data?.message || 'Failed to fetch profile data');
      }
    },
  });

  const mutation = useMutation({
    mutationFn: updateProfile,
    onSuccess: () => {
      queryClient.invalidateQueries(['profile']);
      toast.success("Profile updated successfully");
      setIsEditing(false);
    },
    onError: (error) => {
      if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        navigate('/login');
      } else {
        toast.error(error.response?.data?.message || "Failed to update profile");
      }
    },
  });

  useEffect(() => {
    if (data && data.data) {
      setProfileData(data.data);
    }
  }, [data]);

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setProfileData(prev => ({ ...prev, [name]: value }));
    setIsEditing(true);
  };

  const handleSave = () => {
    const updatedData = {
      first_name: profileData.first_name,
      last_name: profileData.last_name,
      phone_number: profileData.phone_number,
      line_id: profileData.line_id,
      address: profileData.address,
      age: parseInt(profileData.age),
      birth_date: new Date(profileData.birth_date).toISOString(),
    };
    mutation.mutate(updatedData);
  };

  if (isLoading) return <div>Loading...</div>;
  if (isError) return <div>Error loading profile data</div>;

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col">
      <Header />
      <Navigation />
      <main className="container mx-auto px-4 py-8 flex-grow">
        <h1 className="text-2xl font-bold mb-6">Profile</h1>
        <Card>
          <CardHeader>
            <CardTitle>Your Information</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <label className="font-medium w-1/3">Name</label>
                <div className="flex w-2/3 space-x-2">
                  <Input
                    name="first_name"
                    value={profileData.first_name || ''}
                    onChange={handleInputChange}
                    className="flex-1"
                  />
                  <Input
                    name="last_name"
                    value={profileData.last_name || ''}
                    onChange={handleInputChange}
                    className="flex-1"
                  />
                </div>
              </div>
              <div className="flex items-center justify-between">
                <label className="font-medium w-1/3">Phone Number</label>
                <Input
                  name="phone_number"
                  value={profileData.phone_number || ''}
                  onChange={handleInputChange}
                  className="w-2/3"
                />
              </div>
              <div className="flex items-center justify-between">
                <label className="font-medium w-1/3">Address</label>
                <Input
                  name="address"
                  value={profileData.address || ''}
                  onChange={handleInputChange}
                  className="w-2/3"
                />
              </div>
              <div className="flex items-center justify-between">
                <label className="font-medium w-1/3">Line ID</label>
                <Input
                  name="line_id"
                  value={profileData.line_id || ''}
                  onChange={handleInputChange}
                  className="w-2/3"
                />
              </div>
              <div className="flex items-center justify-between">
                <label className="font-medium w-1/3">Age</label>
                <Input
                  name="age"
                  type="number"
                  value={profileData.age || ''}
                  onChange={handleInputChange}
                  className="w-2/3"
                />
              </div>
              <div className="flex items-center justify-between">
                <label className="font-medium w-1/3">Birth Date</label>
                <Input
                  name="birth_date"
                  type="date"
                  value={profileData.birth_date ? new Date(profileData.birth_date).toISOString().split('T')[0] : ''}
                  onChange={handleInputChange}
                  className="w-2/3"
                />
              </div>
            </div>
            {isEditing && (
              <Button className="mt-4 w-full" onClick={handleSave}>
                Save
              </Button>
            )}
          </CardContent>
        </Card>
      </main>
    </div>
  );
};

export default Profile;
