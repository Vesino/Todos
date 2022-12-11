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
      console.log("Sending TODO to save", {formState});
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
