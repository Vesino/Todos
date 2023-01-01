
import { useState } from "react"

import {Modal, Button, Form} from 'react-bootstrap'

export default function TodoModalUpdate({todo, show, handleClose, updateTodos, deleteTodo}) {
    const [currentTodo, setTodo] = useState(todo)

    const handleDelete = (e) => {
        e.preventDefault()
        const url = `http://localhost:3333/v1/todos/${currentTodo.id}`
        const requestOptions = {
            method: "DELETE",
        }
        fetch(url, requestOptions)
        .then(res => res)
        .then(
            deleteTodo(currentTodo)
        )
        .catch(error => console.error(error))

        handleClose();
    }

    const handleIsDone = (e) => {
        e.preventDefault()
        const {name, value} = e.target
        if (value === "on") {
            setTodo({
                ...currentTodo,
                is_done: true
            })
        }
        return
    }

    const onInputChange = (e) => {
        e.preventDefault()
        const {name, value} = e.target
        setTodo({
            ...currentTodo,
            [name]: value
        })
    }

    const handleUpdateTodo = (e) => {
        e.preventDefault()
        if (currentTodo.todo === "" || currentTodo.description === "") return;
        const url = `http://localhost:3333/v1/todos/${currentTodo.id}`
        delete currentTodo.id;
        const requestOptions = {
            method: 'PUT',
            body: JSON.stringify(currentTodo)
        }
        fetch(url, requestOptions)
            .then(res => res.json())
                .then(
                ({todo}) => {
                    updateTodos(todo)
                },
                (error) => {
                    console.error(error)
                }
            )
            .catch(error => {
                console.error(error);
            })

        handleClose();
    } 

    return(
        <Modal show={show} onHide={handleClose}>
            <Modal.Header closeButton>
                <Modal.Title>
                    Update Todo
                </Modal.Title>
            </Modal.Header>
            <Modal.Body>
            <Form>
                <Form.Group className="mb-3" controlId="exampleForm.ControlInput1">
                <Form.Label>Todo</Form.Label>
                <Form.Control
                    type="text"
                    name="todo"
                    autoFocus
                    placeholder="Todo"
                    onChange={onInputChange}
                    value={currentTodo.todo}
                />
                </Form.Group>
                <Form.Group
                className="mb-3"
                controlId="exampleForm.ControlTextarea1"
                >
                <Form.Label>Description</Form.Label>
                <Form.Control 
                as="textarea" 
                rows={3} 
                placeholder="Description"
                name="description"
                value={currentTodo.description}
                onChange={onInputChange}
                />
                </Form.Group>
                <Form.Check 
                type="switch"
                id="todo-switch"
                onChange={handleIsDone}
                name="is_done"
                label="Is Done"
                />
            </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="danger" onClick={handleDelete}>Delete</Button>
                <Button variant="secondary" onClick={handleClose}>Close</Button>
                <Button variant="primary" onClick={handleUpdateTodo}>
                    Save Changes
                </Button>
            </Modal.Footer>
        </Modal>
    )
}
