import React from 'react'
import {Card,ListGroup} from 'react-bootstrap';
import ListGroupItem from 'react-bootstrap/ListGroupItem';
import Form from 'react-bootstrap/Form';
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import { GoPencil, GoTrashcan } from "react-icons/go"


export const Project = (props) => {
  return (
    <Card style={{ width: '18rem', margin:25 }}>
      <Card.Body style={{display: 'flex', justifyContent:'space-between'}}>
        <Card.Title style={{flex: 2}}>{props.name}</Card.Title>
        <Button style={{display: 'flex', alignItems: 'center', justifyContent:'center'}}>
          <GoPencil/>
        </Button>
        <Button style={{display: 'flex', alignItems: 'center', justifyContent:'center', marginLeft: 5}}>
          <GoTrashcan/>
        </Button>
      </Card.Body>
      <ListGroup>
        <Card.Title>
          TODO
        </Card.Title>
      </ListGroup>
      <ListGroup>
        <Card.Title>
          Done
        </Card.Title>
      </ListGroup>
    </Card>
  )
}
