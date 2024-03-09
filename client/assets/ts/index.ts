interface LoginResponse {
  expire_at: string
  token: string
}

interface Todolist { id: number, checked: boolean, content: string }

function setCookie(cName: string, cValue: any, expDays: number) {
  const date = new Date()
  date.setTime(date.getTime() + (expDays * 24 * 60 * 60 * 1000))
  const expires = `expires=${date.toUTCString()}`
  document.cookie = `${cName}=${JSON.stringify(cValue)}; ${expires}; path=/`
}

window.onload = () => {
  const form = document.getElementById('login-form') as HTMLFormElement | null
  if (form) {
    form.addEventListener('submit', (e: SubmitEvent) => {
      e.preventDefault()
      e.stopPropagation()

      const formData = new FormData(form)
      const jsonData: any = {}

      formData.forEach((value, key) => {
        jsonData[key] = value
      })

      fetch('/api/user/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(jsonData),
      })
        .then(e => e.json())
        .then((e: LoginResponse) => {
          setCookie('token', e.token, (new Date(e.expire_at)).getTime())
        })
    })
  }

  function toggleTodolist(id: number, checked: boolean) {
    fetch(`/api/todolist/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ checked }),
    }).then()
  }

  function onToggleTodolist(e: Event) {
    if (!(e.target instanceof HTMLInputElement))
      return

    const value = e.target.getAttribute('value')
    if (!value)
      return

    toggleTodolist(Number.parseInt(value), e.target.checked)
  }

  const todolists = document.querySelectorAll<HTMLInputElement>('input.todolist-item')
  todolists.forEach((t) => {
    t.addEventListener('change', onToggleTodolist)
  })

  function createTodoListItem(res: Todolist) {
    const label = document.createElement('label')
    label.classList.add('flex', 'items-center', 'gap-4')

    const input = document.createElement('input')
    input.setAttribute('type', 'checkbox')
    input.classList.add('todolist-item')
    input.value = res.id.toString()

    if (res.checked)
      input.setAttribute('checked', 'checked')

    input.addEventListener('change', onToggleTodolist)

    label.appendChild(document.createTextNode(res.content))
    label.appendChild(input)

    const list = document.getElementById('todolist-list')
    if (list)
      list.appendChild(label)
  }

  const todoForm = document.getElementById('todo-form') as HTMLFormElement | null
  if (todoForm) {
    todoForm.addEventListener('submit', (e) => {
      e.preventDefault()
      e.stopPropagation()

      const formData = new FormData(todoForm)
      const jsonData: any = {}

      formData.forEach((value, key) => {
        jsonData[key] = value
      })

      fetch('/api/todolist', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(jsonData),
      })
        .then(e => e.json())
        .then(createTodoListItem)
    })
  }
}
