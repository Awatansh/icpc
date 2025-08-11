import React, { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

const App: React.FC = () => {
  const navigate = useNavigate()

  useEffect(() => {
    const storedName = localStorage.getItem('name')
    if (storedName) {
      navigate('/dashboard')
    }
  }, [navigate])

  const handleLogin = () => {
    window.location.href = "http://localhost:3000/auth/google"
  }

  return (
    <div>
      <h1>Login Page</h1>
      <button onClick={handleLogin}>Login with Google</button>
    </div>
  )
}

export default App
