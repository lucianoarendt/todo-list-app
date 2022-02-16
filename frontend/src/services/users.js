import React, {useEffect,useCallback} from 'react';
import axios from 'axios';

const ReadUser = () => {
  const [user, setUser] = React.useState();

  const readUser = useCallback(async () => {
    const response = await axios.get('http://localhost:8000/api/user/read', {
      withCredentials: true,
    })
    const content = await response.data;
    setUser(content);
  });

  useEffect(() => {
    readUser();
  }, []);

  return (
    <h1>??</h1>
  );
}

export default ReadUser;