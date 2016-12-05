package service

import "github.com/mklimuk/husar/train"

func (s *timetable) SearchTrain(query string, removeDuplicates bool) (t []*train.Train, err error) {
	if t, err = s.store.Search(query); err != nil {
		return
	}
	if removeDuplicates {
		res := make(map[string]*train.Train)
		for _, tr := range t {
			res[tr.TrainID] = tr
		}
		t = make([]*train.Train, 0, len(res))
		for _, tr := range res {
			t = append(t, tr)
		}
	}
	return
}
