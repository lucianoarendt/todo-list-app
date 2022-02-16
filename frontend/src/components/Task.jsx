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
    if (props.item.status == 1) { return }
    await axios.put('http://localhost:8000/api/task/update', {
      status: 1,
    },{
      params: {id: props.item.ID},
      withCredentials: true,
    });
    props.deleted();
  }

  let timer1 = setTimeout(() => setShow(false), 5000)
  useEffect(
    () => {
      return () => {
        clearTimeout(timer1)
      }
    },
    [show]
  )

  return (
    <ListGroup className="list-group-flush">

      <Container 
        style={{display: 'flex', alignItems: 'center', justifyContent: 'space-between'}}
        onMouseEnter={() => setTrash(true)}
        onMouseLeave={() => setTrash(false)}

        ref={target} onClick={() => setShow(!show)}
      >
        <>
          <GoCalendar style={{marginRight: 15}}/>
          <Overlay 
            target={target.current} 
            show={show} 
            placement="left"
            onClick={timer1}
          >
            {(evs) => (
              <Tooltip id="overlay-example" {...evs}>
                {props.item.status == 1
                  ? `concluded at: ${(new Date(props.item.UpdatedAt).toLocaleString('pt-BR'))}`
                  : `created at: ${(new Date(props.item.CreatedAt).toLocaleString('pt-BR'))}`
                }
                
              </Tooltip>
            )}
          </Overlay>
        </>
        <Form.Check 
          checked={props.item.status == 1}
          aria-label="??" 
          onClick={updateStatus}
        />
        {props.edit == props.item.ID && props.item.status != 1
          ? <Form.Control
              style={{marginLeft: 15, marginRight: 15}}
              value={desc}
              onChange={e => setDesc(e.target.value)}
              onMouseOut={updateDesc}
              onDoubleClick={() => props.handleEdit('')}
            />
          : <Form.Text 
          style={{
            margin: 15, 
            fontSize: 18, 
            flex: 2,
            textDecoration: props.item.status == 1 ? 'line-through' : ''
          }}
          onClick={() => props.handleEdit(props.item.ID)}  
        >{props.item.description}</Form.Text>
        }
        {
          (trash && props.item.status != 1) ? <GoTrashcan onClick={removeTask}/> : <></>
        }
        
      </Container>
    </ListGroup>
  )
}
