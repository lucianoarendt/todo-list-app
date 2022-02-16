import axios from 'axios';
import React, {useContext, useEffect, useState} from 'react';
import { Container, Button, Form } from 'react-bootstrap';
import { resolveModuleName } from 'typescript';
import { Project } from '../components/Project';
import { AuthContext } from '../Store/AuthContext';
import { UserContext } from '../Store/UserContext';

export const Home = () => {
  const {userDetails, readUser} = useContext(UserContext);
  const [title, setTitle] = useState('');
  const [edit, setEdit] = useState(0);

  const createProject = async (e) => {
    e.preventDefault();

    await axios.post('http://localhost:8000/api/project/create', {
      title: title
    },{
      withCredentials: true,
      headers: {'Content-Type': 'application/json'},
    });
    readUser();
  }

  const editHandle = (e) => {
    if (e == '') {
      setEdit(0)
    } else {
      setEdit(e)
    }
  }

  return (
    <Container style={{display: 'flex'}}>
      <Container style={{display: 'flex', flexWrap: 'wrap', alignItems: 'flex-start'}}>
        {userDetails.projects.map(item => {
          return (
              <Project
                key={item.ID} 
                id={item.ID}
                name={item.title}
                edit={edit}
                handleEdit={e => {editHandle(e)}}
                deleted={readUser}
              ></Project>
          )
        })}
      </Container>
      <Container style={{maxWidth: 300, maxHeight: 220, marginTop: 50, backgroundColor: '#eee', borderRadius: 5, padding: 35}}>
        <Form onSubmit={createProject}>
          <Form.Group className="mb-3" controlId="projectName">
            <Form.Label>Create a new Project</Form.Label>
            <Form.Control 
              placeholder="Project name"
              value={title}
              onChange={e => setTitle(e.target.value)}
            />
          </Form.Group>
          <Button 
            variant="primary" 
            type="submit"
            style={{width: 230}}
          >
            Create
          </Button>
        </Form>
      </Container>
    </Container>
  )
}
