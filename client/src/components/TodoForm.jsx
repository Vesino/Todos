import useTodoForm from "../../hooks/useForm"

function TodoForm() {
    const { todo, description, onInputChange, handleSaveTodo } = useTodoForm()
  return (
    <>
        <h1>Formulario simple</h1>
        <hr />
        <input 
            type="text" 
            name="todo" 
            className="form-control"
            placeholder="Todo"
            onChange={onInputChange}
            value={todo}
        />
        <textarea 
            name="description" 
            className="form-control mt-2"
            placeholder="Description"
            rows="5"
            cols="33"
            onChange={onInputChange}
            value={description}
        />
        <input 
            type="button"
            className="form-control mt-3 btn btn-primary"
            onClick={handleSaveTodo}
            value="Save todo"
        />
    </>
  )
}

export default TodoForm