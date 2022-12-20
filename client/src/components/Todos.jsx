import { memo } from "react"



const Todos = memo(({todos}) => {
    
  return (
    <>
    <h1 className="mt-2">Todos</h1>
  
    {todos.length == 0 && <h1>There is not Todos to show</h1>}
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
              <td>{todo?.is_done?.toString()}</td>
            </tr>
            ))}
        </tbody>
      </table>
    }
    </>
  )
})

export default Todos