package app

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repository interface {
	GetAll() ([]Todo, error)
	Get(id int) (Todo, error)
	Save(t Todo) (Todo, error)
	DeleteAll() error
	Delete(id int) error
	Update(id int, t Todo) (Todo, error)
	Drop() error
}

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(c *Config) (*postgresRepository, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", c.DatabaseHost, c.DatabasePort, c.DatabaseUsername, c.DatabaseSchema, c.DatabasePassword)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	p := &postgresRepository{db: db}

	db.AutoMigrate(&todo{})

	return p, nil
}

func (p *postgresRepository) GetAll() ([]Todo, error) {
	var t []todo

	err := p.db.Order(`"order" asc`).Find(&t).Error
	if err != nil {
		return nil, err
	}

	var r []Todo
	for _, v := range t {
		r = append(r, toDTO(v))
	}

	return r, nil
}

func (p *postgresRepository) Get(id int) (Todo, error) {
	var t todo

	err := p.db.First(&t, id).Error
	if err != nil {
		return blank, err
	}

	return toDTO(t), nil
}

func (p *postgresRepository) Save(t Todo) (Todo, error) {
	m := toModel(t)

	err := p.db.Create(&m).Error
	if err != nil {
		return blank, err
	}

	return toDTO(m), nil
}

func (p *postgresRepository) DeleteAll() error {
	err := p.db.Delete(&todo{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) Delete(id int) error {
	err := p.db.Delete(&todo{
		Model: gorm.Model{ID: uint(id)},
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *postgresRepository) Update(id int, t Todo) (Todo, error) {
	var m todo

	err := p.db.First(&m, id).Error
	if err != nil {
		return blank, err
	}

	if t.Completed == nil {
		var b bool
		t.Completed = &b
	}

	if t.Order == nil {
		var i int
		t.Order = &i
	}

	if t.Title == nil {
		var s string
		t.Title = &s
	}

	err = p.db.Model(&m).Updates(toModel(t)).Error
	if err != nil {
		return blank, err
	}

	return toDTO(m), nil
}

func (p *postgresRepository) Drop() error {
	err := p.db.Unscoped().Delete(&todo{}).Error
	if err != nil {
		return err
	}

	return nil
}

type todo struct {
	gorm.Model
	Title     string `gorm:"type:text"`
	Completed bool
	Order     int
}

func toModel(t Todo) todo {
	return todo{
		Title:     *t.Title,
		Completed: *t.Completed,
		Order:     *t.Order,
	}
}

func toDTO(t todo) Todo {
	id := int(t.ID)
	return Todo{
		Id:        &id,
		Title:     &t.Title,
		Completed: &t.Completed,
		Order:     &t.Order,
	}
}
