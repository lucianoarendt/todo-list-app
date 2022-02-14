import React, {useContext, useEffect, useState} from 'react';
import { Project } from '../components/Project';
import { AuthContext } from '../Store/Context';
import {getUser} from '../services/user';

export const Home = () => {
  const {authenticated, user} = useContext(AuthContext);
  const [userProjects, setUserProjects] = useState(null);

    // useEffect(() => {
    //   let mounted = true;
    //   getUser()
    //     .then(items => {
    //       if(mounted) {
    //         console.log(items)
    //         setUserProjects(items)
    //       }
    //     })
    //   return () => mounted = false;
    // }, [userProjects])  

  return (
    <div>
      Hi {authenticated ? user.name : "not logged person"}
      {userProjects.map(e => {
        console.log(e.title)
        // return <Project name={e.title}></Project>
      })}
      <Project name={"projeto 1"}></Project>
    </div>
  )
}
