import axios from 'axios';
import React, { createContext, useEffect, useState } from 'react';

export const UserContext = createContext();

export const UserProvider = ({ children }) => {
  const [userDetails, setUserDetails] = useState(null);
  const [loading, setLoading] = useState(null);


  useEffect(() => {
    const storedUser = localStorage.getItem('userDetails');
    if (storedUser) {
      setUserDetails(JSON.parse(storedUser));
    }
    setLoading(false);
  }, []);

  const readUser = async () => {
    const response = await axios.get('http://localhost:8000/api/user/read', {
      withCredentials: true,
    })
    const content = await response.data;
    localStorage.setItem('userDetails', JSON.stringify(content));
    setUserDetails(content);
  };

  return (
    <UserContext.Provider value={{userDetails, loading, readUser}}>
      {children}
    </UserContext.Provider>
  );
}