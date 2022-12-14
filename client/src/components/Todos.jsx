import { useState, useEffect } from "react"

export default function Todos() {
    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [todos, setTodos] = useState([])

    useEffect(() => {
        fetch("http://localhost:3333/v1/todos/")
          .then(res => res.json())
          .then(
            (result) => {
              setIsLoaded(true);
              setTodos(result?.todos);
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
    <h1 className="mt-2">Todos</h1>
  
    { error && todos.length == 0 && <h1>An error occurred </h1>}
    {todos.length > 0 && 
      <table className="table table-striped">
          <thead>
          <tr>
            <th scope="col">#</th>
            <th scope="col">Todo</th>
            <th scope="col">Description</th>
            <th scope="col">Created at</th>
            <th scope="col">Is Done</th>
          </tr>
        </thead>
        <tbody>
            {todos.map(todo => (
              <tr key={todo?.id}>
              <th scope="row">{todo?.id}</th>
              <td>{todo?.todo}</td>
              <td>{todo?.description}</td>
              <td>{todo?.created_at}</td>
              <td>{todo?.is_done}</td>
            </tr>
            ))}
        </tbody>
      </table>
    }
    </>
  )
}
