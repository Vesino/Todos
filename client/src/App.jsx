import { useCallback } from "react"
import { useState, useEffect } from "react"
import TodoForm  from './components/TodoForm'
import Todos from './components/Todos'

function App() {
  const [todos, setTodos] = useState([])

  const addTodo = useCallback(
    (todo) => {
      setTodos(todos => [...todos, todo])
    },
    [],
  )
  
  useEffect(() => {
    console.log('Me trigerio el useEffect')
      fetch("http://localhost:3333/v1/todos/")
        .then(res => res.json())
        .then(
          (result) => {
            setTodos(result?.todos);
          },

          // Note: it's important to handle errors here
          // instead of a catch() block so that we don't swallow
          // exceptions from actual bugs in components.
          (error) => {
            console.log(error)
          }
        )
    }, [])

  return (
    <>
    <TodoForm 
      addTodo={addTodo}
      todos={todos}
      />
    <Todos 
      todos={todos}
    />
    </>
  )
}

export default App
