package todo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

type Item struct {
	Text     string
	Priority int
	position int
	Done     bool
}

type ByPri []Item

func (a ByPri) Len() int      { return len(a) }
func (a ByPri) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPri) Less(i, j int) bool {
	if a[i].Done && !a[j].Done {
		return a[i].Done
	}
	if a[i].Priority == a[j].Priority {
		return a[i].position < a[j].position
	}
	return a[i].position < a[j].position
}

func SaveItems(fileName string, items []Item) error {

	b, err := json.Marshal(items)
	if err != nil {
		return err
	}

	slog.Debug("Saving items to file", "fileName", fileName)

	err = os.WriteFile(fileName, b, 0644)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func LoadItems(fileName string) ([]Item, error) {

	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var items []Item
	err = json.Unmarshal(b, &items)
	if err != nil {
		return nil, err
	}

	for i, _ := range items {
		items[i].position = i + 1
	}

	slog.Debug("Loaded items from file", "fileName", fileName)
	return items, nil
}

func (i *Item) SetPriority(pri int) {
	switch pri {
	case 0:
		i.Priority = 0
	case 1:
		i.Priority = 1
	case 2:
		i.Priority = 2
	default:
		i.Priority = 0
	}
}

func (i *Item) PrettyP() string {
	switch i.Priority {
	case 0:
		return "Low"
	case 1:
		return "Medium"
	case 2:
		return "High"
	default:
		return "Low"
	}
}

func (i *Item) Lable() string {
	return strconv.Itoa(i.position) + "."
}
func (i *Item) DoneStatus() string {
	if i.Done {
		return "Done"
	}
	return "Not Done"
}
