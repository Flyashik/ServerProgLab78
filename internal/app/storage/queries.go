package storage

const (
	SelectAllSuperheroes = `SELECT * FROM superhero.superhero`

	SelectSuperheroesWithPublisher = `SELECT s.id, s.superhero_name, p.publisher_name,
       									SUM(ha.attribute_value) AS total_attributes
										FROM superhero.superhero s
										INNER JOIN superhero.publisher p ON s.publisher_id = p.id
										LEFT JOIN superhero.hero_attribute ha ON s.id = ha.hero_id
										LEFT JOIN superhero.attribute a ON ha.attribute_id = a.id
										GROUP BY s.id, s.superhero_name, p.publisher_name
										ORDER BY SUM(ha.attribute_value)  DESC;`

	AddSuperPower = `INSERT INTO superhero.superpower (power_name) VALUES ($1);`

	DeleteSuperPower = `WITH deleted_hero_power AS (
    					DELETE FROM superhero.hero_power
    					WHERE power_id = (SELECT id FROM superhero.superpower WHERE power_name = $1)
    					RETURNING *
						)
						DELETE FROM superhero.superpower
						WHERE power_name = $1;`

	ChangeSuperPower = `UPDATE superhero.hero_power 
						SET power_id = (SELECT id FROM superhero.superpower WHERE power_name = ($3))
						WHERE hero_id = (SELECT id FROM superhero.superhero WHERE superhero_name = ($1)) AND power_id = (SELECT id FROM superhero.superpower WHERE power_name = ($2));`

	AddPowerForHero = `INSERT INTO superhero.hero_power (hero_id, power_id)
						VALUES ((SELECT id FROM superhero.superhero WHERE superhero_name = ($1)), 
						        (SELECT id FROM superhero.superpower WHERE power_name = ($2)))`
)
