import { useState } from 'react'

export default function useTodoForm() {
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
        .then(response => console.log(response.json()));
      setFormState({
        todo: "",
        description: "",}
      )
    }

  return {
    ...formState,
    handleSaveTodo,
    onInputChange
  }
}
