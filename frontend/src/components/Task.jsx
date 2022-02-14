import React from 'react'

export const Task = (props) => {
  return (
    <ListGroup className="list-group-flush">
      <Container style={{display: 'flex', alignItems: 'center'}}>
        <Form.Check aria-label="option 1" />
        <Form.Text style={{margin: 15, fontSize: 18}}>{props.description}</Form.Text>
      </Container>
    </ListGroup>
  )
}
