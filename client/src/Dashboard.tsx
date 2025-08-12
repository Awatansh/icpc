import React, { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'

const Dashboard: React.FC = () => {
  const [name, setName] = useState<string>('')
  const navigate = useNavigate()

  useEffect(() => {
    const params = new URLSearchParams(window.location.search)
    const n = params.get('name')

    if (n) {
      setName(n)
      localStorage.setItem('name', n)
    } else {
      const storedName = localStorage.getItem('name')
      if (storedName) {
        setName(storedName)
      } else {
        navigate('/') // no name? go to login
      }
    }
  }, [navigate])

  const handleLogout = () => {
  fetch("http://localhost:3000/logout", {
    method: "GET",
    credentials: "include"
  })
    .then(() => {
      localStorage.removeItem('name')
      navigate('/')
    })
    .catch(err => {
      console.error("Logout failed:", err)
    })
}


  return (
    <div>
      <h1>Welcome, {name} 👋</h1>
      <button onClick={handleLogout}>Logout</button>
    </div>
  )
}

export default Dashboard
