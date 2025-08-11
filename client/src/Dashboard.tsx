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
    localStorage.removeItem('name')
    navigate('/')
  }

  return (
    <div>
      <h1>Welcome, {name} 👋</h1>
      <button onClick={handleLogout}>Logout</button>
    </div>
  )
}

export default Dashboard
