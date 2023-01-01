import { useState } from "react";

import TodoModalUpdate from './TodoModalUpdate'

export default function Todo({todo, updateTodos}) {
    const [show, setShow] = useState(false);
    const handleClose = () => setShow(false);
    const handleShow = () => setShow(true);
  return (
    <>
    <tr key={todo?.id} onClick={handleShow}>
    <th scope="row">{todo?.id}</th>
    <td>{todo?.todo}</td>
    <td>{todo?.description}</td>
    <td>{todo?.created_at}</td>
    <td>{todo?.is_done?.toString()}</td>
    </tr>
    {show ?
        <TodoModalUpdate 
            show={show}
            todo={todo}
            handleClose={handleClose}
            updateTodos={updateTodos}
        /> : 
    null
    }
    </>
  )
}
