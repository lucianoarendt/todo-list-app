import React, { createContext, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const navigate = useNavigate();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const storedUser = localStorage.getItem('user');
    if (storedUser) {
      setUser(JSON.parse(storedUser));
    }
    setLoading(false);
  }, []);

  const login = async (email, password) => {
    console.log('login auth', email, password);
    fetch('http://localhost:8000/api/login', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      credentials: 'include',
      body: JSON.stringify({
          email,
          password
      })
    }).then(res => {
      localStorage.setItem('user', JSON.stringify(res));
      navigate("/")
    })
  };

  const logout = async () => {
    console.log('logout');

    await fetch('http://localhost:8000/api/logout', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        credentials: 'include',
    });

    localStorage.removeItem('user');
    navigate('/login');
  };

  return (
    <AuthContext.Provider
      value={{ authenticated: !!user, user, login, logout}}
    >
      {children}
    </AuthContext.Provider>
  );
};
