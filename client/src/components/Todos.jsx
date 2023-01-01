import { memo } from "react"
import Table from 'react-bootstrap/Table'
import Todo from './Todo'


const Todos = memo(({todos, updateTodos}) => {

  return (
    <>
    <h1 className="mt-2">Todos</h1>
  
    {todos.length == 0 && <h1>There is not Todos to show</h1>}
    {todos.length > 0 && 
      <Table striped bordered hover >
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
              <Todo key={todo.id} todo={todo} updateTodos={updateTodos}/>
            ))}
        </tbody>
      </Table>
    }
    </>
  )
})

export default Todos