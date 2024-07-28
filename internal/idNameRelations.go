package internal

import "fmt"

// IdNameRelation is the go representaion of a simple id name relation
type IdNameRelation struct {
	Id int
	Name string
}

// LoadIdNameRelation will load simple id name relations
// the returned map can be used to derive the ids from the strings or vice versa
func (s *Server) LoadIdNameRelation(rName string) (map[string]int, error)  {

	q := fmt.Sprintf("SELECT * FROM %s", rName)
	rows, err := s.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	relMap  := map[string]int{}
	for rows.Next() {
		var inr IdNameRelation
		if err := rows.Scan(&inr.Id, &inr.Name); err != nil {
			return nil, err
		}
		relMap[inr.Name] = inr.Id
	}

	return relMap, nil
}