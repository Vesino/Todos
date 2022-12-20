import { useState } from 'react'

export default function useTodoForm({addTodo, todos}) {
    const [formState, setFormState] = useState({
        todo: "",
        description: "",
    })

    const onInputChange = ({target}) => {
        const {name, value} = target
        setFormState({
            ...formState,
            [name]: value
        })
    }

    const handleSaveTodo = (e) => {
      e.preventDefault();
      if (formState.todo === "" || formState.description === "") return;
      const requestOptions = {
        method: 'POST',
        body: JSON.stringify({ todo: formState.todo, description: formState.description })
    };
    fetch('http://localhost:3333/v1/todos/', requestOptions)
        .then(response => response.json())
           .then(({ todo }) => {
             setFormState({
               todo: "",
               description: "",}
             )
             addTodo(todo)
        }
      )
    }

  return {
    ...formState,
    handleSaveTodo,
    onInputChange
  }
}
