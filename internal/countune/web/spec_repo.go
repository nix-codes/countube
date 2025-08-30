package web

import (
	"countube/internal/common"

	"encoding/json"
	"os"
	"sort"
)

type CountuneSpec struct {
	Id        int    `json:"id"`
	StartNum  int    `json:"start_num"`
	BarCount  int    `json:"bar_count"`
	BarColor1 string `json:"bar_color_1"`
	BarColor2 string `json:"bar_color_2"`
	BgColor   string `json:"bg_color"`
}

type CountuneSpecRepo struct {
	filename string
	specs    map[int]CountuneSpec
}

func NewCountuneSpecRepo() (*CountuneSpecRepo, error) {
	common.EnsurePath(OutputPath)
	filename := OutputPath + CountunePicSpecStorageFile

	repo := &CountuneSpecRepo{
		filename: filename,
		specs:    make(map[int]CountuneSpec),
	}

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return repo, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var loaded []CountuneSpec
	if err := json.NewDecoder(file).Decode(&loaded); err != nil {
		return nil, err
	}

	for _, s := range loaded {
		repo.specs[s.Id] = s
	}

	return repo, nil
}

func (r *CountuneSpecRepo) Upsert(spec CountuneSpec) {
	r.specs[spec.Id] = spec
}

func (r *CountuneSpecRepo) FindAll() []CountuneSpec {
	all := make([]CountuneSpec, 0, len(r.specs))
	for _, s := range r.specs {
		all = append(all, s)
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].Id < all[j].Id
	})

	return all
}

func (r *CountuneSpecRepo) FindMaxId() int {
	maxId := -1
	for id := range r.specs {
		if id > maxId {
			maxId = id
		}
	}
	return maxId
}

func (r *CountuneSpecRepo) GetById(id int) *CountuneSpec {
	spec, ok := r.specs[id]
	if ok {
		return &spec
	}
	return nil
}

func (r *CountuneSpecRepo) Persist() error {
	file, err := os.Create(r.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	all := r.FindAll()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(all)
}
