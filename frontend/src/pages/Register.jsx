import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import axios from 'axios';
import toast from 'react-hot-toast';

const Register = () => {
  const navigate = useNavigate();
  const [step, setStep] = useState(1);
  const [formData, setFormData] = useState({
    username: '',
    first_name: '',
    last_name: '',
    password: '',
    confirm_password: '',
    phone_number: '',
    line_id: '',
    address: '',
    age: '',
    birth_date: '',
  });
  const [passwordMatchError, setPasswordMatchError] = useState('');
  const [phoneNumberError, setPhoneNumberError] = useState('');
  const [passwordFormatError, setPasswordFormatError] = useState('');
  const [addressError, setAddressError] = useState('');

  const validatePassword = (password) => {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]).{6,}$/;
    if (!regex.test(password)) {
      setPasswordFormatError('Password must contain at least 1 uppercase letter, 1 lowercase letter, 1 number, 1 special character, and be at least 6 characters long');
    } else {
      setPasswordFormatError('');
    }
  };

  const validatePhoneNumber = (phoneNumber) => {
    const regex = /^0\d{9}$/;
    if (!regex.test(phoneNumber)) {
      setPhoneNumberError('Phone number must be 10 digits and start with 0');
    } else {
      setPhoneNumberError('');
    }
  };

  const validateAddress = (address) => {
    if (address.length <= 10) {
      setAddressError('Address must be more than 10 characters');
    } else {
      setAddressError('');
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
    
    if (name === "password") {
      validatePassword(value);
    } else if (name === "confirm_password") {
      setPasswordMatchError(formData.password !== value ? 'Passwords do not match' : '');
    } else if (name === "phone_number") {
      validatePhoneNumber(value);
    } else if (name === "address") {
      validateAddress(value);
    }
  };

  const validateStep1 = () => {
    const { username, first_name, last_name, password, confirm_password } = formData;
    return username && first_name && last_name && password && confirm_password && !passwordMatchError && !passwordFormatError;
  };

  const validateStep2 = () => {
    const { phone_number, address, age, birth_date } = formData;
    return phone_number && address && age && birth_date && !phoneNumberError && !addressError;
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (validateStep1() && validateStep2()) {
      try {
        const registrationData = {
          ...formData,
          age: parseInt(formData.age),
          birth_date: new Date(formData.birth_date).toISOString(),
        };
        const response = await axios.post('http://localhost:8080/api/users/register', registrationData);
        toast.success('Registration successful!');
        navigate('/login');
      } catch (error) {
        toast.error(error.response?.data?.message || 'Registration failed. Please try again.');
      }
    }
  };

  const nextStep = () => {
    if (step === 1 && validateStep1()) {
      setStep(2);
    }
  };

  const prevStep = () => setStep(1);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <Card className="w-[400px]">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">Register</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit}>
            {step === 1 && (
              <>
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="username">Username <span className="text-red-500">*</span></Label>
                    <Input id="username" name="username" value={formData.username} onChange={handleChange} placeholder="eg.guessuser" required />
                  </div>
                  <div>
                    <Label htmlFor="first_name">First Name <span className="text-red-500">*</span></Label>
                    <Input id="first_name" name="first_name" value={formData.first_name} onChange={handleChange} placeholder="John" required />
                  </div>
                  <div>
                    <Label htmlFor="last_name">Last Name <span className="text-red-500">*</span></Label>
                    <Input id="last_name" name="last_name" value={formData.last_name} onChange={handleChange} placeholder="Doe" required />
                  </div>
                  <div>
                    <Label htmlFor="password">Password <span className="text-red-500">*</span></Label>
                    <Input id="password" name="password" type="password" value={formData.password} onChange={handleChange} placeholder="eg.j_Za1234" required />
                    {passwordFormatError && <p className="text-red-500 text-sm">{passwordFormatError}</p>}
                  </div>
                  <div>
                    <Label htmlFor="confirm_password">Confirm Password <span className="text-red-500">*</span></Label>
                    <Input id="confirm_password" name="confirm_password" type="password" value={formData.confirm_password} onChange={handleChange} placeholder="eg.j_Za1234" required />
                    {passwordMatchError && <p className="text-red-500 text-sm">{passwordMatchError}</p>}
                  </div>
                </div>
                <Button className="w-full mt-4" onClick={nextStep} disabled={!validateStep1()}>Next</Button>
              </>
            )}
            {step === 2 && (
              <>
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="phone_number">Phone Number <span className="text-red-500">*</span></Label>
                    <Input id="phone_number" name="phone_number" value={formData.phone_number} onChange={handleChange} placeholder="09XXXXXXXX" required />
                    {phoneNumberError && <p className="text-red-500 text-sm">{phoneNumberError}</p>}
                  </div>
                  <div>
                    <Label htmlFor="line_id">Line ID</Label>
                    <Input id="line_id" name="line_id" value={formData.line_id} onChange={handleChange} />
                  </div>
                  <div>
                    <Label htmlFor="address">Address <span className="text-red-500">*</span></Label>
                    <Input id="address" name="address" value={formData.address} onChange={handleChange} placeholder="ChiangMai MaeRim RimTai" required />
                    {addressError && <p className="text-red-500 text-sm">{addressError}</p>}
                  </div>
                  <div>
                    <Label htmlFor="age">Age <span className="text-red-500">*</span></Label>
                    <Input id="age" name="age" type="number" value={formData.age} onChange={handleChange} placeholder="12" required />
                  </div>
                  <div>
                    <Label htmlFor="birth_date">Birth Date <span className="text-red-500">*</span></Label>
                    <Input id="birth_date" name="birth_date" type="date" value={formData.birth_date} onChange={handleChange} required />
                  </div>
                </div>
                <div className="flex justify-between mt-4">
                  <Button onClick={prevStep}>Previous</Button>
                  <Button type="submit" disabled={!validateStep2()}>Register</Button>
                </div>
              </>
            )}
          </form>
        </CardContent>
      </Card>
    </div>
  );
};

export default Register;
