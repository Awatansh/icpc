import React, { useEffect } from 'react'
import { useNavigate } from 'react-router-dom'

const App: React.FC = () => {
  const navigate = useNavigate()

  useEffect(() => {
    fetch("http://localhost:3000/api/auth/status", {
      credentials: "include" // send cookies!
    })
      .then(res => res.json())
      .then(data => {
        if (data.authenticated) {
          // store name if you want for quick access
          localStorage.setItem("name", data.name)
          navigate('/dashboard')
        }
      })
      .catch(err => console.error("Auth check failed:", err))
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
