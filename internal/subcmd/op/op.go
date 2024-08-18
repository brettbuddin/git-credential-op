package op

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const DefaultLocatorTag = "git-credential-op"

// Runner executes commands against 1Password.
type Runner struct {
	// LocatorTag is a tag value we assign to each 1Password item managed by this helper. Only items
	// tagged with this value are seen by the tool.
	LocatorTag string
	// Account is account URL (e.g. mycompany.1password.com). An empty value indicates the default.
	Account string
	// Vault is name of the 1Password Vault (e.g. Private). An empty value indicates the default.
	Vault string
	// Executor executes commands; usually with "op".
	Executor Executor
	// Stdin is where we'll read input from stdin.
	Stdin io.Reader
	// Stdout is where we'll forward stdout output.
	Stdout io.Writer
	// Stderr is where we'll forward stderr output.
	Stderr io.Writer
}

// ErrNotFound is returned when we can't locate an item in 1Password.
var ErrNotFound = fmt.Errorf("not found")

// FindItem locates the first item in 1Password that has a matching URL and is tagged with
// [LocatorTag]. It returns an [ErrNotFound] error if no item is found.
func (r *Runner) FindItem(targetURL string) (Item, error) {
	items, err := r.listItems()
	if err != nil {
		return Item{}, err
	}
	for _, it := range items {
		for _, url := range it.URLs {
			if url.HRef == targetURL {
				return r.GetItem(it.ID)
			}
		}
	}
	return Item{}, ErrNotFound
}

// Item is a 1Password item.
type Item struct {
	ID       string
	Username string
	Password string
}

// listItem is a list-level view of a 1Password item.
type listItem struct {
	ID   string        `json:"id"`
	URLs []listItemURL `json:"urls"`
}

// listItemURL contains URL details of a 1Password item.
type listItemURL struct {
	HRef string `json:"href,omitempty"`
}

// GetItem fetches details about an item in 1Password.
func (r *Runner) GetItem(id string) (Item, error) {
	var out bytes.Buffer
	err := r.execItemCommand(
		r.execOutput(&out),
		"get",
		"--format", "json",
		"--fields", "username,password",
		id,
	)
	if err != nil {
		return Item{}, err
	}

	var fields []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	}
	if err := json.Unmarshal(out.Bytes(), &fields); err != nil {
		return Item{}, fmt.Errorf("json: %w", err)
	}
	var item Item
	item.ID = id
	for _, f := range fields {
		switch f.Label {
		case "username":
			item.Username = f.Value
		case "password":
			item.Password = f.Value
		}
	}
	return item, nil
}

// listItems lists all items in 1Password with [LocatorTag].
func (r *Runner) listItems() ([]listItem, error) {
	var out bytes.Buffer
	err := r.execItemCommand(
		r.execOutput(&out),
		"list",
		"--categories", "login",
		"--format", "json",
		"--tags", r.LocatorTag,
	)
	if err != nil {
		return nil, err
	}
	var entries []listItem
	if err := json.Unmarshal(out.Bytes(), &entries); err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}
	return entries, nil
}

// CreateItem creates an item in 1Password.
func (r *Runner) CreateItem(req CreateRequest) error {
	tags := []string{r.LocatorTag}
	tags = append(tags, req.AdditionalTags...)

	return r.execItemCommand(
		r.execOutput(io.Discard),
		"create",
		"--category", "login",
		"--title", req.Title,
		"--url", req.URL,
		"--tags", strings.Join(tags, ","),
		"username="+req.Username,
		"password="+req.Password,
	)
}

// CreateRequest are details for creating an item in 1Password.
type CreateRequest struct {
	Title          string
	AdditionalTags []string
	Username       string
	Password       string
	URL            string
}

// UpdateItem updates an item in 1Password.
func (r *Runner) UpdateItem(it Item) error {
	return r.execItemCommand(
		r.execOutput(io.Discard),
		"edit",
		it.ID,
		"username="+it.Username,
		"password="+it.Password,
	)
}

// DeleteItem deletes an item in 1Password.
func (r *Runner) DeleteItem(id string) error {
	return r.execItemCommand(r.execOutput(io.Discard), "delete", id)
}
