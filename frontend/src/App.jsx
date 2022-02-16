import React from 'react';
import { useEffect, useState, useContext } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import './App.css';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Home } from './pages/Home';
import Nav from './components/Navbar';
import { AuthProvider, AuthContext } from './Store/AuthContext';
import { UserProvider, UserContext } from './Store/UserContext';
import ReadUser from './services/users';

function App() {

  const Private = ({ children }) => {
    const { authenticated, loading } = useContext(AuthContext);

    if (loading) {
      return <div className="loading">Carregando...</div>;
    }

    if (!authenticated) {
      return <Navigate to="/login" />;
    }
    return children;
  };

  return (
    <div className="App">
      <BrowserRouter>
        <AuthProvider>
          <UserProvider>
            <Nav/>
            <Routes>
              <Route exact path="/" element={<Home/>}/>
              <Route path="/login" element={<Login/>}/>
              <Route path="/register" element={<Register/>} />
            </Routes>
          </UserProvider>
        </AuthProvider>
      </BrowserRouter>
    </div>
  );
}

export default App;
