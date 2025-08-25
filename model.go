package main

import (
	"database/sql"
	"fmt"
	"log"
)

type Product struct {
	ID       int     `json:"ID"`
	Name     string  `json:"Name"`
	Quantity int     `json:"Quantity"`
	Price    float64 `json:"Price"`
}

func allProducts(DB *sql.DB) ([]Product, error) {
	query := "SELECT id,name,quantity,price from products"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	productList := []Product{}
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Quantity, &product.Price)
		if err != nil {
			return productList, err
		}
		productList = append(productList, product)
	}
	return productList, nil
}

func (p *Product) oneProduct(DB *sql.DB) error {
	query := fmt.Sprintf("SELECT name,quantity,price from products where id=%v;", p.ID)
	rows := DB.QueryRow(query)
	//Since id is unique, rows will be 1
	err := rows.Scan(&p.Name, &p.Quantity, &p.Price) //assigning values to the p object
	return err
}

func (p *Product) addProduct(DB *sql.DB) error {
	query := fmt.Sprintf("insert into products(name,quantity,price) values('%v',%v,%v)", p.Name, p.Quantity, p.Price)
	res, err := DB.Exec(query)
	if err != nil {
		return nil
	}
	LastInsertId, _ := res.LastInsertId()
	fmt.Println("Last insertedID:", LastInsertId)
	p.ID = int(LastInsertId)
	return nil
}

func (p *Product) updateProduct(DB *sql.DB) error {
	query := fmt.Sprintf("update products set name=\"%v\", quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.ID)
	res, err := DB.Exec(query)
	if err != nil {
		return err
	}
	log.Println(p)
	// log.Println(res.RowsAffected())
	return nil
}
