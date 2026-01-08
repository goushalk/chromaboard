`
Step 1: Define the problem. Write a purpose statement.
Step 2: List user actions, not code. This becomes the workflow.
Step 3: Map actions to responsibilities (modules). Think roles, not filenames.
Step 4: Every file must have a single purpose. If it can’t explain why it exists, remove it.
Step 5: Correct workflow to start a project: README → skeleton → data → happy path → split files.
Step 6: When stuck, write comments explaining inputs, outputs, and callers before writing code
`

# QnA for the Projects

## 1) what problem am i facing ?

- i need to arrange and organize my Projects in TODO, Pending, Done.

## 2&3) what action the user can or should do ?

- Primary action for Prototype :

  - Primary Project action:

    - create/delete a Project description.

      - to create and delete the Projects we need to save the info of it.
      - the project info will be saved in project.json ("$HOME/.local/share/choromaboard/projects/<name_of_the_project>.json").
      - to perform the action we need to read & write the projects.json file.

    - rename Project

      - we need project id to find it and rename the project so we need to read the project file(can be optimized)

  - Primary Task action:

    - create/delete a task with task name and description.
    - move the task along the three category (TODO, Pending, done)
    - edit already existing task and change the category
      - we be will reading and writing the perticular project files to identify the task and perform the above action.
      - To move the task around the category we will pop the needed task name for the list and append is to the required text

# 4) Change the purpose into file name

- main.go : for binding up all the other binary.
- internal/tui/app.go to bind all the view, update, model.
- internal/tui/update.go : for update method in bubble tea.
- internal/tui/view.go : for the view method in bubble tea.
- internal/storage/project_store.go : all the structs for handling the data.
- internal/storage/local.go : all func for creating and providing the path of file and checks for already existing file.
- internal/app/projects.go : all the logic to read and write the data to the file. (i Think both projects.go and structs.go can be combined as  a single file)

# 5) skeleton:-

`
/chromaboard
  /cmd/chromaboard
    main.go

  /internal
    /tui
      app.go
      update.go
      view.go

    /domain
      project.go        // Project, Task, domain rules

    /storage
      local.go          // filesystem paths, dirs
      project_store.go  // load/save Project

---- in future we may add /configs /theme -----
`

example for the structure:
`
{
  "version": 1,
  "project": {
    "id": "uuid",
    "name": "NameOfTheProject",
    "tasks": [
      { "id": "1", "title": "taskOne", "status": "TODO" },
      { "id": "2", "title": "taskTwo", "status": "TODO" },
      { "id": "3", "title": "taskThree", "status": "Pending" },
      { "id": "4", "title": "taskFour", "status": "Done" }
    ]
  }
}

`

Code order

## domain/project.go

- First. Always.
- Define Project
- Define Task
- Define status

### Add pure methods

- AddTask
- MoveTask
- RenameProject

No IO. No Bubble Tea. No JSON.
If this file feels clean, the project will be clean.

---
2️⃣

## storage/local.go

Then storage utilities:

- Base directory
- Project path resolution
- Ensure directories exist
- Still no UI.

## storage/project_store.go

Only now:

- Load project
- Save project

This file should feel boring. That’s a good sign.

## tui/app.go

Define:

- Model
- ActivePane
- ActiveColumn
- Selected IDs
No logic yet.

## tui/update.go

Implement:

- Key handling
- Pane transitions
- Calling domain methods
- Triggering saves

This is where the app “comes alive”.

## tui/view.go

- Last.
- Rendering only.
- No decisions.
- No mutation.

`
  file and the function singnature:

  storage/local_storage.go :

    func DataDir() (string, error),
    func ProjectDir() (string, error),
    func ProjectRegPath(projectName string) (string, error),
    func EnsureStorage() error,
    func ExistsPath(path string) (bool, error)

  storage/project.go :

    func LoadRegistry(projectName string) (domain.Project, error),
    func SaveRegistry(project domain.Project) error

  domain/project.go :
  
    func NewProject(name string) Project
    func (p *Project) AddTask(title string)
    func (p *Project) MoveTask(taskID int, status Status) error
    func (p *Project) DeleteTask(taskID int) error
    func (p *Project) Rename(name string)
`
