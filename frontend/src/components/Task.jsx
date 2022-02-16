import axios from 'axios';
import React, {useState, useEffect, useRef} from 'react'
import {ListGroup,Container,Form,Overlay,Tooltip} from 'react-bootstrap'
import { GoCalendar, GoTrashcan } from "react-icons/go"

export const Task = (props) => {
  const [desc, setDesc] = useState('');
  const [trash, setTrash] = useState(false);
  const [show, setShow] = useState(false);
  const target = useRef(null);

  const updateDesc = async () => {
    await axios.put('http://localhost:8000/api/task/update', {
      description: desc,
    },{
      params: {id: props.item.ID},
      withCredentials: true,
    });
    props.deleted();
  }

  const removeTask = async () => {
    const response = await axios.delete('http://localhost:8000/api/task/delete', {
      params: {id: props.item.ID},
      withCredentials: true,
    });
    props.deleted();
  }

  const updateStatus = async () => {
    const response = await axios.put('http://localhost:8000/api/task/update', {
      status: 1,
    },{
      params: {id: props.item.ID},
      withCredentials: true,
    });
    props.deleted();
  }

  return (
    <ListGroup className="list-group-flush">

      <Container 
        style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}
        onMouseEnter={() => setTrash(true)}
        onMouseLeave={() => setTrash(false)}

        ref={target} onClick={() => setShow(!show)}
      >
        <>
          <GoCalendar/>
          <Overlay target={target.current} show={show} placement="left">
            {(props) => (
              <Tooltip id="overlay-example" {...props}>
                {props.edit.}
              </Tooltip>
            )}
          </Overlay>
        </>
        <Form.Check 
          aria-label="??" 
          onClick={updateStatus}
          onChange={updateStatus}
        />
        {props.edit == props.item.ID
          ? <Form.Control
            style={{marginLeft: 15, marginRight: 15}}
            value={desc}
            onChange={e => setDesc(e.target.value)}
            onMouseOut={updateDesc}
            onDoubleClick={() => props.handleEdit('')}
          >
          </Form.Control>
          : <Form.Text 
          style={{margin: 15, fontSize: 18, flex: 2}}
          onClick={() => props.handleEdit(props.item.ID)}  
        >{props.item.description}</Form.Text>
        }
        {
          trash ? <GoTrashcan onClick={removeTask}/> : <></>
        }
        
      </Container>
    </ListGroup>
  )
}
