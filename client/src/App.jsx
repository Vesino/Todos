import { useCallback } from "react"
import { useState, useEffect } from "react"
import TodoForm  from './components/TodoForm'
import Todos from './components/Todos'

function App() {
  const [todos, setTodos] = useState([])

  const deleteTodo = useCallback(
    (todo) => {
      // this remove the deleted todo from the todos state
      const newTodos = todos.filter(obj => obj.id != todo.id)
      setTodos(newTodos)
  })

  const addTodo = useCallback(
    (todo) => {
      setTodos(todos => [...todos, todo])
    },
    [],
  )

  const updateTodos = useCallback(
    (todo) => {
      const todosUpdate = todos.filter(obj => obj.id != todo.id)
      todosUpdate.push(todo)
      setTodos(todosUpdate)
    }
  )
  
  useEffect(() => {
      fetch("http://localhost:3333/v1/todos")
        .then(res => res.json())
        .then(
          (result) => {
            setTodos(result?.todos);
          },
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
      updateTodos={updateTodos}
      deleteTodo={deleteTodo}
    />
    </>
  )
}

export default App
