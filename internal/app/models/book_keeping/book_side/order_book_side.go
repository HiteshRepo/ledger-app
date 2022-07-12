package book_side

import (
	"github.com/hashicorp/go-multierror"
	comparator2 "github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/book_keeping/comparator"
	"github.com/hiteshpattanayak-tw/SupplyDemandLedger/internal/app/models/order"
	"sort"
	"sync"
)

func ProvideOrderBookSide() BookSide {
	return &orderBookSide{orders: make([]*order.Order, 0)}
}

type orderBookSide struct {
	mtx sync.RWMutex

	orders     []*order.Order
	comparator comparator2.Comparator
}

func (a *orderBookSide) GetOrders() []*order.Order {
	a.mtx.RLock()
	defer a.mtx.RUnlock()

	return a.orders
}

func (a *orderBookSide) SetComparator(comparator comparator2.Comparator) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	a.comparator = comparator
}

func (a *orderBookSide) UpdateOrders(newOrders []*order.Order) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	var sortErr error
	sort.Slice(newOrders, func(i int, j int) bool {
		result, err := a.comparator(newOrders[i], newOrders[j])
		if err != nil {
			sortErr = multierror.Append(sortErr, err)
			return false
		}
		return result > 0
	})

	if sortErr != nil {
		return sortErr
	}

	o, err := a.merge(a.orders, newOrders)
	if err != nil {
		return err
	}

	a.orders = o
	return nil
}

func (a *orderBookSide) merge(existOrders, newOrders []*order.Order) ([]*order.Order, error) {
	newLen := len(existOrders) + len(newOrders)
	result := make([]*order.Order, newLen)

	i := 0
	for len(existOrders) > 0 && len(newOrders) > 0 {

		if newOrders[0].Qty.IsZero() {
			id := getExistingOrderById(newOrders[0].Id ,existOrders)
			if id != -1 {
				existOrders = append(existOrders[0:id], existOrders[id+1:]...)
				newOrders = newOrders[1:]
				newLen -= 2
				continue
			}
			newLen -= 1
			continue
		}

		r, err := a.comparator(existOrders[0], newOrders[0])
		if err != nil {
			return nil, err
		}

		switch {
		case r == 0:
			result[i] = newOrders[0]
			result[i+1] = existOrders[0]
			newLen -= 2
			i += 2
			existOrders = existOrders[1:]
			newOrders = newOrders[1:]
		case r > 0:
			result[i] = existOrders[0]
			existOrders = existOrders[1:]
			i++
		default:
			result[i] = newOrders[0]
			i++
			newOrders = newOrders[1:]
		}
	}

	for j := 0; j < len(existOrders); j++ {
		result[i] = existOrders[j]
		i++
	}
	for j := 0; j < len(newOrders); j++ {
		if !newOrders[0].Qty.IsZero() {
			result[i] = newOrders[j]
			i++
		} else {
			newLen -= 1
		}
	}

	result = result[:newLen]

	return result, nil
}

func getExistingOrderById(id string, existOrders []*order.Order) int {
	for i,e := range existOrders {
		if e.Id == id {
			return i
		}
	}
	return -1
}