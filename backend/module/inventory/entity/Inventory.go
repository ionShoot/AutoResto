package entity

type Inventory struct {
	id       int    `form:"id" json:"id"`
	capacity int    `form:"capacity" json:"capacity"`
	location string `form:"location" json:"location"`
}

type InventoryResponse struct {
	Message string      `form:"message" json:"message"`
	Data    []Inventory `form:"data" json:"data"`
}

func (i *Inventory) GetId() int {
	return i.id
}

func (i *Inventory) SetId(id int) {
	i.id = id
}

func (i *Inventory) GetCapacity() int {
	return i.capacity
}

func (i *Inventory) SetCapacity(capacity int) {
	i.capacity = capacity
}

func (i *Inventory) GetLocation() string {
	return i.location
}

func (i *Inventory) SetLocation(location string) {
	i.location = location
}
