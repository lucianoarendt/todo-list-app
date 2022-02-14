import React from 'react'

export const Task = () => {
  return (
    <ListGroup className="list-group-flush">
        <Container style={{display: 'flex', alignItems: 'center'}}>
          <Form.Check aria-label="option 1" />
          <Form.Text style={{margin: 15, fontSize: 18}}>Cras justo odio</Form.Text>
        </Container>
        <ListGroupItem>Dapibus ac facilisis in</ListGroupItem>
        <ListGroupItem>Vestibulum at eros</ListGroupItem>
      </ListGroup>
  )
}
