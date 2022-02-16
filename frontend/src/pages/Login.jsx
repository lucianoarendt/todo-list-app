import React, {useState, useContext} from 'react';
import { AuthContext } from '../Store/AuthContext';
import { UserContext } from '../Store/UserContext';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';

export const Login = () => {
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const {login} = useContext(AuthContext);
    const {readUser} = useContext(UserContext);

    const submit = async (e) => {
        e.preventDefault();
        login(email, password);
        readUser();
    }
        
    return (
        <Container fluid style={{maxWidth: 450, padding: 50}}>
            <Form className="justify-content-md-center" onSubmit={submit}>
                <Form.Group className="mb-3" controlId="formBasicEmail">
                    <Form.Label>Email address</Form.Label>
                    <Form.Control type="email" placeholder="Enter email" 
                        onChange={e => setEmail(e.target.value)}
                    />
                </Form.Group>
                <Form.Group className="mb-3" controlId="formBasicPassword">
                    <Form.Label>Password</Form.Label>
                    <Form.Control type="password" placeholder="Password" 
                        onChange={e => setPassword(e.target.value)}
                    />
                </Form.Group>
                <Container style={{display: 'flex', justifyContent: 'center', marginTop: 40}}>
                    <Button variant="dark" type="submit" style={{width: 300}}>
                        Sign In
                    </Button>
                </Container>
            </Form>
        </Container>
    )
}