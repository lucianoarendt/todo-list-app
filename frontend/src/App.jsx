import React from 'react';
import { useEffect, useState, useContext } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import './App.css';
import { Login } from './pages/Login';
import { Register } from './pages/Register';
import { Home } from './pages/Home';
import Nav from './components/Navbar';
import { AuthProvider, AuthContext } from './Store/Context';

function App() {

  // useEffect(() => {
  //   async function fetchUser() {
  //     const response = await fetch('http://localhost:8000/api/user/read', {
  //       headers: { 'Content-Type': 'application/json' },
  //       credentials: 'include',
  //     });
  //     const content = await response.json();
  //     setName(content.name);
  //   }
  //   fetchUser();
  //   console.log(name);
  // }, []);

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
          <Nav/>
          <Routes>
            <Route exact path="/" element={
              <Private>
                <Home/>
              </Private>
            }/>
            <Route path="/login" element={<Login/>}/>
            <Route path="/register" element={<Register/>} />
          </Routes>
        </AuthProvider>
      </BrowserRouter>
    </div>
  );
}

export default App;
