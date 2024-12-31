// acts as schema or blueprint for employee data
package schema

//defining a custom struct for employee
type Employee struct {
	//'omitempty' tag is used to omit the field if the value is empty
	//JSON tag is used to map the field name to JSON key
	EID        string `json:"e_id,omitempty"`       //unique employee id
	Name       string `json:"name,omitempty"`       //employee name
	Department string `json:"department,omitempty"` //employee department
	Position   string `json:"position,omitempty"`   //employee position
}
