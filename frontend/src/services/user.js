import {AuthContext} from '../Store/Context'
import { useContext } from 'react'

export function getUser() {
  const {authenticated,user} = useContext(AuthContext);
  if (!authenticated) {
    return
  }
  return fetch('http://localhost:8000/api/user/read', {
    param: {id: user.id}
  })
    .then(data => data.json())
}