import { useState, useEffect } from "react"

export default function Todos() {
    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [todos, setTodos] = useState([])

    useEffect(() => {
        fetch("http://localhost:3333/todos")
          .then(res => res.json())
          .then(
            (result) => {
              setIsLoaded(true);
              setTodos(result);
            },
            // Note: it's important to handle errors here
            // instead of a catch() block so that we don't swallow
            // exceptions from actual bugs in components.
            (error) => {
              setIsLoaded(true);
              setError(error);
            }
          )
      }, [])
    
  return (
    <>
    { error && todos.length == 0 && <h1>An error occurred </h1>}
    {todos.length > 0 && <ul>
        {todos.map(todo => (
            <li key={todo?.todo}>
                {todo.todo}, {todo?.description}
            </li>
        ))}
        </ul>
    
    }
    </>
  )
}
