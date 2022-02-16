import React, {useContext, useEffect} from 'react';
import { AuthContext } from '../Store/AuthContext';
import { UserContext } from '../Store/UserContext';
import { Link } from 'react-router-dom';
import Navbar from 'react-bootstrap/Navbar';
import Container from 'react-bootstrap/Container';

const Nav = () => {
  const {authenticated, user, logout} = useContext(AuthContext);
  const {userDetails} = useContext(UserContext);

  // useEffect(() => {
  //   console.log(userDetails)
  //   console.log(authenticated)
  // }, [userDetails]);
  
  let menu;

  if (!authenticated) {
    menu = (
      <Navbar.Collapse className="justify-content-end">
        <Navbar.Text style={{paddingRight: 30}}>
          <Link to="/login">Login</Link>
        </Navbar.Text>
        <Navbar.Text>
          <Link to="/register">Register</Link>
        </Navbar.Text>
      </Navbar.Collapse>
    )
  } else {
    menu = (
      <Container style={{display: 'flex'}}>
        <Navbar.Collapse className="justify-content-center">
          <Navbar.Text>
            Signed in as: {userDetails.name}
          </Navbar.Text>
        </Navbar.Collapse>
        <Navbar.Collapse className="justify-content-end" onClick={logout}>
          <Navbar.Text>
            Logout
          </Navbar.Text>
        </Navbar.Collapse> 
      </Container> 
    )
  }

  return (
    <Navbar bg="dark" variant="dark">
      <Container>
        <Navbar.Brand>ToDo List App</Navbar.Brand>
        <Navbar.Toggle />
        {menu}
      </Container>
    </Navbar>
  )
}

export default Nav;


