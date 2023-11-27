package main

import "fmt"

// Definición de un tipo llamado 'Person'
type Person struct {
	FirstName string
	LastName  string
	Age       int
}

// Método asociado al tipo 'Person' para obtener el nombre completo
func (p Person) FullName() string {
	return p.FirstName + " " + p.LastName
}

// Método asociado al tipo 'Person' para saludar
func (p Person) Greet() {
	fmt.Printf("Hola, soy %s y tengo %d años.\n", p.FullName(), p.Age)
}

func main() {
	// Creación de una instancia de 'Person'
	person := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}

	// Uso de métodos
	person.Greet()

	// Modificación de atributos directamente (sin encapsulación)
	person.Age = 31

	// Uso de otro método
	person.Greet()
}
