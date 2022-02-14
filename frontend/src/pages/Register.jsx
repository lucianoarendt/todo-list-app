import React, {useState} from 'react';
import { Link } from 'react-router-dom';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';

export const Register = () => {
    const [name, setName] = useState("")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [redirect, setRedirect] = useState(false)

    const submit = async (e) => {
        e.preventDefault();

        const response = await fetch('http://localhost:8000/api/user/create', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({
                name,
                email,
                password
            })
        });

        const content = await response.json();
        if (response.status == 200) {
            setRedirect(true)
        }     
    }

    // if (redirect) {
    //     return (
    //         <Container style={{display: 'flex', flexDirection: 'column', alignItems: 'center', marginTop: 30}}>
    //             <Form.Text style={{marginBottom: 30, fontSize: 40}}> Welcome {name}!</Form.Text>
    //             <Link to="/login">
    //                 <Button variant="dark" style={{width: 200}}>Login</Button>
    //             </Link>
    //         </Container>
    //     )
    // }

    return (
      <div>
        register
      </div>
        // <Container fluid style={{maxWidth: 450, padding: 50}}>
        //     <Form className="justify-content-md-center" onSubmit={submit}>
        //     <Form.Group className="mb-3" controlId="formBasicEmail">
        //             <Form.Label>Your name</Form.Label>
        //             <Form.Control placeholder="Enter your name"
        //                 onChange={e => setName(e.target.value)}
        //             />
        //         </Form.Group>
        //         <Form.Group className="mb-3" controlId="formBasicEmail">
        //             <Form.Label>Email address</Form.Label>
        //             <Form.Control type="email" placeholder="Enter email" 
        //                 onChange={e => setEmail(e.target.value)}
        //             />
        //         </Form.Group>
                
        //         <Form.Group className="mb-3" controlId="formBasicPassword">
        //             <Form.Label>Password</Form.Label>
        //             <Form.Control type="password" placeholder="Password" 
        //                 onChange={e => setPassword(e.target.value)}
        //             />
        //         </Form.Group>
        //         <Container style={{display: 'flex', justifyContent: 'center', marginTop: 40}}>
        //             <Button variant="dark" type="submit" style={{width: 300}}>
        //                 Sign Up
        //             </Button>
        //         </Container>
        //     </Form>
        // </Container>
    )
}