import axios from 'axios';
import Cookies from 'js-cookie';
import React, { createContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    const response = await axios.post('http://localhost:8000/api/login', {
        email,
        password
      },{
      withCredentials: true,
      headers: {'Content-Type': 'application/json'},
    })
    const content = await response.data
    localStorage.setItem('user', JSON.stringify(content));
    setAuthenticated(true);
    navigate("/")
  };

  const logout = async () => {
    const response = await axios.post('http://localhost:8000/api/logout', {
        withCredentials: 'include',
    });
    Cookies.remove('jwt');
    localStorage.removeItem('user');
    setAuthenticated(false);
    navigate('/login');
  };

  return (
    <AuthContext.Provider
      value={{ authenticated, user, login, logout, loading}}
    >
      {children}
    </AuthContext.Provider>
  );
};
