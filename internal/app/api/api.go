package api

import (
	"ApiService/internal/app/config"
	"ApiService/internal/app/model"
	"ApiService/internal/app/storage"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Server struct {
	config  *config.Config
	logger  *logrus.Logger
	router  *mux.Router
	storage *storage.Storage
}

func New(config *config.Config) *Server {
	return &Server{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureStorage(); err != nil {
		return err
	}

	s.logger.Info("Starting server...")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *Server) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/superhero", s.handleSuperhero())
	s.router.HandleFunc("/add_superpower", s.handleAddSuperpower())
	s.router.HandleFunc("/delete_superpower", s.handleDeleteSuperpower())
	s.router.HandleFunc("/change_superpower", s.handleChangeSuperpower())
	s.router.HandleFunc("/add_power_hero", s.handleAddPowerForHero())
}

func (s *Server) configureStorage() error {
	st := storage.New(s.config.Storage)
	if err := st.Open(); err != nil {
		return err
	}

	s.storage = st

	return nil
}

func (s *Server) handleSuperhero() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Response struct {
			Heroes []model.Superhero `json:"superhero"`
		}
		var resp Response

		rows, err := s.storage.Query(storage.SelectAllSuperheroes)
		if err != nil {
			s.logger.Error(err)
		}
		defer rows.Close()

		for rows.Next() {
			var hero model.Superhero

			var (
				fullName     sql.NullString
				genderId     sql.NullInt64
				eyeColourId  sql.NullInt64
				hairColourId sql.NullInt64
				skinColourId sql.NullInt64
				raceId       sql.NullInt64
				publisherId  sql.NullInt64
				alignmentId  sql.NullInt64
				height       sql.NullInt64
				weight       sql.NullInt64
			)

			err = rows.Scan(
				&hero.Id, &hero.Name, &fullName, &genderId, &eyeColourId,
				&hairColourId, &skinColourId, &raceId, &publisherId,
				&alignmentId, &height, &weight,
			)
			if err != nil {
				s.logger.Error(err)
				return
			}

			hero.FullName = fullName.String
			hero.GenderId = int(genderId.Int64)
			hero.EyeColourId = int(genderId.Int64)
			hero.HairColourId = int(hairColourId.Int64)
			hero.SkinColourId = int(skinColourId.Int64)
			hero.RaceId = int(raceId.Int64)
			hero.PublisherId = int(publisherId.Int64)
			hero.AlignmentId = int(alignmentId.Int64)
			hero.HeightCm = int(height.Int64)
			hero.WeightKg = int(weight.Int64)

			resp.Heroes = append(resp.Heroes, hero)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) handleAddSuperpower() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}

		var sp *model.Superpower
		err = json.Unmarshal(body, &sp)
		if err != nil {
			s.logger.Error(err)
			return
		}
		_, err = s.storage.Exec(storage.AddSuperPower, sp.Name)
		if err != nil {
			s.logger.Error(err)
			return
		}
	}
}

func (s *Server) handleDeleteSuperpower() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}

		var sp *model.Superpower
		err = json.Unmarshal(body, &sp)
		if err != nil {
			s.logger.Error(err)
			return
		}
		_, err = s.storage.Exec(storage.DeleteSuperPower, sp.Name)
		if err != nil {
			s.logger.Error(err)
			return
		}
	}
}

func (s *Server) handleChangeSuperpower() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Changer struct {
			Hero     model.Superhero  `json:"superhero"`
			OldPower model.Superpower `json:"old_power"`
			NewPower model.Superpower `json:"new_power"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}

		var changer *Changer
		err = json.Unmarshal(body, &changer)
		if err != nil {
			s.logger.Error(err)
			return
		}
		_, err = s.storage.Exec(storage.ChangeSuperPower, changer.Hero.Name, changer.OldPower.Name, changer.NewPower.Name)
		if err != nil {
			s.logger.Error(err)
			return
		}
	}
}

func (s *Server) handleAddPowerForHero() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type Setter struct {
			Hero  model.Superhero  `json:"superhero"`
			Power model.Superpower `json:"power"`
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.logger.Error(err)
			return
		}

		var setter Setter
		err = json.Unmarshal(body, &setter)
		if err != nil {
			s.logger.Error(err)
			return
		}
		_, err = s.storage.Exec(storage.AddPowerForHero, setter.Hero.Name, setter.Power.Name)
		if err != nil {
			s.logger.Error(err)
			return
		}
	}
}
