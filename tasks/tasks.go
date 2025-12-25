package tasks

import "github.com/4frag/padalko/utils"

var Registry = []utils.MenuItem[string]{}

func Register(task utils.MenuItem[string]) {
    Registry = append(Registry, task)
}

func GetByID(id string) *utils.MenuItem[string] {
    for _, task := range Registry {
        if task.ID == id {
            return &task
        }
    }
    return nil
}