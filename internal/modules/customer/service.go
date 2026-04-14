package customer

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aalexanderkevin/getstarvio-backend/internal/models"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/shared"
)

type Service struct{ repo *Repo }

func NewService(repo *Repo) *Service { return &Service{repo: repo} }

func (s *Service) List(userID, q, status, sortBy string) ([]map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	cxs, err := s.repo.ListCustomers(biz.ID, q)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(cxs))
	for _, c := range cxs {
		ids = append(ids, c.ID)
	}
	svcs, err := s.repo.ListServices(ids)
	if err != nil {
		return nil, err
	}
	m := map[string][]models.CustomerService{}
	for _, svc := range svcs {
		m[svc.CustomerID] = append(m[svc.CustomerID], svc)
	}

	out := make([]map[string]interface{}, 0, len(cxs))
	for _, c := range cxs {
		services := m[c.ID]
		worst, overdue := worstStatus(services)
		if status != "" && status != "semua" && worst != status {
			continue
		}
		item := map[string]interface{}{
			"id":          c.ID,
			"name":        c.Name,
			"wa":          c.WA,
			"via":         c.Via,
			"status":      worst,
			"overdueDays": overdue,
			"createdAt":   c.CreatedAt.Format(time.RFC3339),
			"services":    toServiceDTO(services),
		}
		out = append(out, item)
	}

	sortCustomers(out, sortBy)
	return out, nil
}

func (s *Service) Create(userID string, req CreateCustomerRequest) error {
	if req.Name == "" || req.WA == "" {
		return fmt.Errorf("name and wa are required")
	}
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	cats, err := s.repo.ListCategories(biz.ID)
	if err != nil {
		return err
	}
	catByID := map[string]models.Category{}
	for _, c := range cats {
		catByID[c.ID] = c
	}

	wa := shared.NormalizePhone(req.WA, "62")
	via := req.Via
	if via == "" {
		via = "manual"
	}

	cx := models.Customer{
		ID:         uuid.NewString(),
		BusinessID: biz.ID,
		Name:       req.Name,
		WA:         wa,
		Via:        via}
	services, err := buildServicesFromInput(cx.ID, req.Services, catByID, time.Now().UTC())
	if err != nil {
		return err
	}

	return s.repo.CreateCustomerWithServices(cx, services)
}

func (s *Service) Update(userID, customerID string, req UpdateCustomerRequest) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	if _, err := s.repo.FindCustomer(biz.ID, customerID); err != nil {
		return err
	}
	cats, err := s.repo.ListCategories(biz.ID)
	if err != nil {
		return err
	}
	catByID := map[string]models.Category{}
	for _, c := range cats {
		catByID[c.ID] = c
	}

	payload := map[string]interface{}{}
	if req.Name != nil {
		payload["name"] = *req.Name
	}
	if req.WA != nil {
		payload["wa"] = shared.NormalizePhone(*req.WA, "62")
	}
	if req.Via != nil {
		payload["via"] = *req.Via
	}

	services := []models.CustomerService{}
	if len(req.Services) > 0 {
		services, err = buildServicesFromInput(customerID, req.Services, catByID, time.Now().UTC())
		if err != nil {
			return err
		}
	}

	return s.repo.UpdateCustomerAndServices(biz.ID, customerID, payload, services)
}

func (s *Service) Delete(userID, customerID string) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteCustomer(biz.ID, customerID)
}

func (s *Service) RecordVisit(userID string, req VisitRequest) error {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}

	visitDate, err := parseVisitDate(req.Date)
	if err != nil {
		return err
	}
	if err := validateBackdate(visitDate); err != nil {
		return err
	}

	cats, err := s.repo.ListCategories(biz.ID)
	if err != nil {
		return err
	}
	catByID := map[string]models.Category{}
	for _, c := range cats {
		catByID[c.ID] = c
	}

	if len(req.CategoryIDs) == 0 {
		return fmt.Errorf("categoryIds is required")
	}

	var cx *models.Customer
	if req.CustomerID != "" {
		cx, err = s.repo.FindCustomer(biz.ID, req.CustomerID)
		if err != nil {
			return err
		}
	} else {
		if req.CustomerName == "" || req.CustomerWA == "" {
			return fmt.Errorf("customerName and customerWa are required for new customer")
		}
		wa := shared.NormalizePhone(req.CustomerWA, "62")
		existing, err := s.repo.FindCustomerByWA(biz.ID, wa)
		if err == nil {
			cx = existing
		} else if err == gorm.ErrRecordNotFound {
			newCx := models.Customer{ID: uuid.NewString(), BusinessID: biz.ID, Name: req.CustomerName, WA: wa, Via: "manual"}
			if err := s.repo.CreateCustomerWithServices(newCx, nil); err != nil {
				return err
			}
			cx = &newCx
		} else {
			return err
		}
	}

	input := make([]ServiceInput, 0, len(req.CategoryIDs))
	for _, cid := range req.CategoryIDs {
		_, ok := catByID[cid]
		if !ok {
			return fmt.Errorf("invalid category id: %s", cid)
		}
		input = append(input, ServiceInput{CategoryID: cid, Date: visitDate.Format(time.RFC3339)})
	}
	services, err := buildServicesFromInput(cx.ID, input, catByID, visitDate)
	if err != nil {
		return err
	}
	return s.repo.UpdateCustomerAndServices(biz.ID, cx.ID, nil, services)
}

func (s *Service) CheckinLookup(userID, wa string) (map[string]interface{}, error) {
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return nil, err
	}
	fwa := shared.NormalizePhone(wa, "62")
	cx, err := s.repo.FindCustomerByWA(biz.ID, fwa)
	if err == gorm.ErrRecordNotFound {
		return map[string]interface{}{"found": false}, nil
	}
	if err != nil {
		return nil, err
	}

	svcs, err := s.repo.ListServices([]string{cx.ID})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"found": true,
		"customer": map[string]interface{}{
			"id":       cx.ID,
			"name":     cx.Name,
			"wa":       cx.WA,
			"via":      cx.Via,
			"services": toServiceDTO(svcs),
		},
	}, nil
}

func (s *Service) CheckinSubmit(userID string, req CheckinSubmitRequest) error {
	if req.WA == "" || len(req.CategoryIDs) == 0 {
		return fmt.Errorf("wa and categoryIds are required")
	}
	biz, err := s.repo.FindBusinessByUser(userID)
	if err != nil {
		return err
	}
	cats, err := s.repo.ListCategories(biz.ID)
	if err != nil {
		return err
	}
	catByID := map[string]models.Category{}
	for _, c := range cats {
		catByID[c.ID] = c
	}

	visitDate, err := parseVisitDate(req.Date)
	if err != nil {
		return err
	}

	wa := shared.NormalizePhone(req.WA, "62")
	cx, err := s.repo.FindCustomerByWA(biz.ID, wa)
	if err == gorm.ErrRecordNotFound {
		if req.Name == "" {
			return fmt.Errorf("name is required for new checkin")
		}
		newCx := models.Customer{ID: uuid.NewString(), BusinessID: biz.ID, Name: req.Name, WA: wa, Via: "qr"}
		if err := s.repo.CreateCustomerWithServices(newCx, nil); err != nil {
			return err
		}
		cx = &newCx
	} else if err != nil {
		return err
	}

	input := make([]ServiceInput, 0, len(req.CategoryIDs))
	for _, cid := range req.CategoryIDs {
		_, ok := catByID[cid]
		if !ok {
			return fmt.Errorf("invalid category id: %s", cid)
		}
		input = append(input, ServiceInput{CategoryID: cid, Date: visitDate.Format(time.RFC3339)})
	}
	services, err := buildServicesFromInput(cx.ID, input, catByID, visitDate)
	if err != nil {
		return err
	}
	return s.repo.UpdateCustomerAndServices(biz.ID, cx.ID, nil, services)
}

func buildServicesFromInput(customerID string, in []ServiceInput, catByID map[string]models.Category, fallbackDate time.Time) ([]models.CustomerService, error) {
	services := make([]models.CustomerService, 0, len(in))
	for _, it := range in {
		cat, ok := catByID[it.CategoryID]
		if !ok {
			return nil, fmt.Errorf("invalid category id: %s", it.CategoryID)
		}
		dt, err := parseVisitDate(it.Date)
		if err != nil {
			dt = fallbackDate
		}
		services = append(services, models.CustomerService{
			ID:           uuid.NewString(),
			CustomerID:   customerID,
			CategoryID:   cat.ID,
			LastVisitAt:  dt,
			IntervalDays: cat.IntervalDays,
		})
	}
	return services, nil
}

func parseVisitDate(s string) (time.Time, error) {
	if s == "" {
		return time.Now().UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.UTC(), nil
	}
	return time.Time{}, fmt.Errorf("invalid date format")
}

func validateBackdate(t time.Time) error {
	now := time.Now().UTC()
	if t.After(now) {
		return fmt.Errorf("visit date cannot be in the future")
	}
	if now.Sub(t) > (7 * 24 * time.Hour) {
		return fmt.Errorf("visit backdate limit is 7 days")
	}
	return nil
}

func getStatus(s models.CustomerService) string {
	diff := int(time.Since(s.LastVisitAt).Hours() / 24)
	if s.IntervalDays <= 0 {
		return "aktif"
	}
	pct := float64(diff) / float64(s.IntervalDays) * 100
	if pct >= 100 {
		return "hilang"
	}
	if pct >= 70 {
		return "mendekati"
	}
	return "aktif"
}

func worstStatus(services []models.CustomerService) (string, int) {
	if len(services) == 0 {
		return "aktif", 0
	}
	order := map[string]int{"hilang": 0, "mendekati": 1, "aktif": 2}
	worst := "aktif"
	worstOverdue := 0
	for _, s := range services {
		st := getStatus(s)
		over := int(time.Since(s.LastVisitAt).Hours()/24) - s.IntervalDays
		if order[st] < order[worst] {
			worst = st
			if over > 0 {
				worstOverdue = over
			}
		}
	}
	return worst, worstOverdue
}

func toServiceDTO(services []models.CustomerService) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(services))
	for _, s := range services {
		out = append(out, map[string]interface{}{
			"name":   s.ServiceName,
			"icon":   s.ServiceIcon,
			"date":   s.LastVisitAt.Format(time.RFC3339),
			"days":   s.IntervalDays,
			"status": getStatus(s),
		})
	}
	return out
}

func sortCustomers(rows []map[string]interface{}, mode string) {
	if mode == "" {
		mode = "urgent"
	}
	sort.SliceStable(rows, func(i, j int) bool {
		a := rows[i]
		b := rows[j]
		switch mode {
		case "name_asc":
			return strings.ToLower(a["name"].(string)) < strings.ToLower(b["name"].(string))
		case "oldest":
			return a["createdAt"].(string) < b["createdAt"].(string)
		case "newest":
			return a["createdAt"].(string) > b["createdAt"].(string)
		default:
			order := map[string]int{"hilang": 0, "mendekati": 1, "aktif": 2}
			ao := order[a["status"].(string)]
			bo := order[b["status"].(string)]
			if ao == bo {
				return a["overdueDays"].(int) > b["overdueDays"].(int)
			}
			return ao < bo
		}
	})
}
