import axios from 'axios';
import React, {useState, useEffect} from 'react'
import {Card,ListGroup,Button, Form, Container} from 'react-bootstrap';
import { GoPencil, GoTrashcan, GoCheck, GoX } from "react-icons/go"
import { Task } from './Task';

export const Project = (props) => {
  const [title, setTitle] = useState('');
  const [taskTitle, setTaskTitle] = useState('');
  const [tasks, setTasks] = useState(null);
  const [edit, setEdit] = useState(0);

  const updateTitle = async () => {
    const response = await axios.put('http://localhost:8000/api/project/update', {
      id: props.id,
      title: title
    },{
      withCredentials: true,
      headers: {'Content-Type': 'application/json'},
    });
    setTitle(title);
    props.handleEdit('');
  }

  const removeProject = async () => {
    const response = await axios.delete('http://localhost:8000/api/project/delete', {
      params: {id: props.id},
      withCredentials: true,
    });
    props.deleted();
  }

  const readTasks = async () => {
    const response = await axios.get('http://localhost:8000/api/task/read', {
      params: {id: props.id},
      withCredentials: true,
      headers: {'Content-Type': 'application/json'},
    });

    const content = await response.data;
    setTasks(content);
  }

  useEffect(() => {
    readTasks();
  }, []);

  const addTask = async () => {
    await axios.post('http://localhost:8000/api/task/create', {
      description: taskTitle,
      project: props.id
    },{
      withCredentials: true,
      headers: {'Content-Type': 'application/json'},
    });
    readTasks();
  }

  const editHandle = (e) => {
    if (e == '') {
      setEdit(0)
    } else {
      setEdit(e)
    }
  }

  return (
    <Card style={{ width: '18rem', margin:25 }}>
      <Card.Body style={{display: 'flex', justifyContent:'space-between', justifyItems: 'center'}}>
        {props.edit == props.id
          ? <Form.Control 
              style={{marginRight: 15}}
              value={title}
              onChange={e => setTitle(e.target.value)}
            ></Form.Control>        
          : <Card.Title style={{flex: 2, fontSize: 20, padding: 5}}>{title!=''?title:props.name}</Card.Title>
        }
        {props.edit == props.id
          ? <Button 
              onClick={updateTitle}
              style={{display: 'flex', alignItems: 'center', justifyContent:'center'}}>
                <GoCheck/>
            </Button>
          : <Button 
              onClick={() => props.handleEdit(props.id)}
              style={{display: 'flex', alignItems: 'center', justifyContent:'center'}}>
                <GoPencil/>
            </Button>
        }
        {props.edit == props.id
          ? <Button 
              style={{display: 'flex', alignItems: 'center', justifyContent:'center', marginLeft: 5}}
              onClick={() => props.handleEdit('')}
            >
              <GoX/>
            </Button>
          : <Button 
              style={{display: 'flex', alignItems: 'center', justifyContent:'center', marginLeft: 5}}
              onClick={removeProject}
            >
              <GoTrashcan/>
            </Button>
        }
      </Card.Body>
      <ListGroup style={{padding: 15}}>
        <Card.Title>
          To Do
        </Card.Title>
        {
          tasks ? tasks.map((item) => {
            return <Task 
                    key={item.ID}
                    item={item}
                    edit={edit}
                    handleEdit={e => {editHandle(e)}}
                    deleted={readTasks}
                   >{item.description}</Task>
          }) : <></>
        }
      </ListGroup>
      <ListGroup style={{padding: 15}}>
        <Card.Title>
          Done
        </Card.Title>
      </ListGroup>
      <Container style={{display: 'flex', marginTop: 15, marginBottom: 15}}>
        <Form.Control
          style={{}}
          value={taskTitle}
          onChange={e => setTaskTitle(e.target.value)}
        ></Form.Control>
        <Button
          style={{marginLeft: 15}}
          onClick={addTask}
        >
          Add
        </Button>
      </Container>
    </Card>
  )
}
