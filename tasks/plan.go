package tasks

import (
	"fmt"
	"sort"
	"strings"
)

// CalculationResult для первой функции
type CalculationResult struct {
	Usage      []float64
	Deficits   map[int]float64
	IsFeasible bool
}

// 1. Рассчет по конкретному плану (проверка вектора B)
func (m *ModelData) CalculatePlan() CalculationResult {
	usage := make([]float64, m.SizeA)
	deficits := make(map[int]float64)
	isFeasible := true

	for i := 0; i < m.SizeA; i++ {
		for j := 0; j < m.SizeB; j++ {
			usage[i] += m.Model[i][j] * m.B[j]
		}
		if usage[i] > m.A[i] {
			deficits[i] = usage[i] - m.A[i]
			isFeasible = false
		}
	}
	return CalculationResult{Usage: usage, Deficits: deficits, IsFeasible: isFeasible}
}

// 2. Расчет оптимального решения по критериям
func (m *ModelData) SolveWithCriteria(criteria []float64) ([]float64, error) {
    // Проверка размеров
    if len(criteria) != m.SizeB {
        return nil, fmt.Errorf("размер критериев (%d) не совпадает с количеством продуктов (%d)", 
            len(criteria), m.SizeB)
    }

    fmt.Println("=== Расчет оптимального плана по критериям ===")
    fmt.Printf("Критерии: %v\n", criteria)

    // 1. Функция дефицита - проверка, можем ли произвести ХОТЯ БЫ 1 каждого продукта
    for j := 0; j < m.SizeB; j++ {
        for i := 0; i < m.SizeA; i++ {
            if m.Model[i][j] > m.A[i] {
                return nil, fmt.Errorf("дефицит: невозможно произвести продукт %d - не хватит ресурса %d", j+1, i+1)
            }
        }
    }
    fmt.Println("✓ Модель жизнеспособна (все продукты можно произвести по одному)")

    // 2. Целочисленный алгоритм распределения
    resultPlan := make([]float64, m.SizeB)
    remainingResources := make([]float64, m.SizeA)
    copy(remainingResources, m.A)

    // Если все критерии нулевые - производим по 1 каждого (если можем)
    totalCriteria := 0.0
    for _, w := range criteria {
        totalCriteria += w
    }

    if totalCriteria == 0 {
        fmt.Println("Все критерии нулевые - производим по 1 каждого продукта")
        for j := 0; j < m.SizeB; j++ {
            // Проверяем, можем ли произвести еще 1 продукт j
            canProduce := true
            for i := 0; i < m.SizeA; i++ {
                if m.Model[i][j] > remainingResources[i] {
                    canProduce = false
                    break
                }
            }
            if canProduce {
                resultPlan[j] = 1
                for i := 0; i < m.SizeA; i++ {
                    remainingResources[i] -= m.Model[i][j]
                }
            }
        }
        return resultPlan, nil
    }

    // 3. Основной алгоритм: жадное распределение по приоритету критериев
    fmt.Println("\nРаспределяем ресурсы по приоритету критериев...")

    // Создаем список продуктов с их критериями и индексами
    type Product struct {
        Index    int
        Criteria float64
    }
    
    products := make([]Product, m.SizeB)
    for j := 0; j < m.SizeB; j++ {
        products[j] = Product{Index: j, Criteria: criteria[j]}
    }
    
    // Сортируем по убыванию критерия (высший приоритет - больший критерий)
    sort.Slice(products, func(i, k int) bool {
        return products[i].Criteria > products[k].Criteria
    })

    // Жадный алгоритм: пытаемся производить максимально возможное количество
    // продуктов с высшим приоритетом, потом переходим к следующим
    for _, p := range products {
        j := p.Index
        if p.Criteria <= 0 {
            continue // Пропускаем продукты с нулевым или отрицательным приоритетом
        }

        // Сколько единиц продукта j можем произвести с оставшимися ресурсами?
        maxPossible := int(1e9) // большое число
        
        for i := 0; i < m.SizeA; i++ {
            if m.Model[i][j] > 0 {
                possible := int(remainingResources[i] / m.Model[i][j])
                if possible < maxPossible {
                    maxPossible = possible
                }
            }
        }

        if maxPossible > 0 {
            // Производим максимально возможное количество
            resultPlan[j] = float64(maxPossible)
            // Вычитаем использованные ресурсы
            for i := 0; i < m.SizeA; i++ {
                remainingResources[i] -= m.Model[i][j] * float64(maxPossible)
            }
            fmt.Printf("  Продукт %d: произведено %d ед. (критерий: %.2f)\n", 
                j+1, maxPossible, p.Criteria)
        }
    }

    // 4. Проверяем, можно ли улучшить план, используя оставшиеся ресурсы
    // (производим продукты с меньшим приоритетом из остатков)
    fmt.Println("\nИспользуем оставшиеся ресурсы для продуктов с меньшим приоритетом...")
    
    // Снова проходим по всем продуктам (теперь в порядке увеличения приоритета)
    sort.Slice(products, func(i, k int) bool {
        return products[i].Criteria < products[k].Criteria
    })
    
    for _, p := range products {
        j := p.Index
        if resultPlan[j] == 0 && p.Criteria > 0 {
            // Для продуктов, которые еще не производились
            // Смотрим, сколько можем произвести из остатков
            maxPossible := int(1e9)
            for i := 0; i < m.SizeA; i++ {
                if m.Model[i][j] > 0 {
                    possible := int(remainingResources[i] / m.Model[i][j])
                    if possible < maxPossible {
                        maxPossible = possible
                    }
                }
            }
            
            if maxPossible > 0 {
                // Можем произвести хотя бы 1 дополнительную единицу
                for maxPossible > 0 {
                    // Проверяем, стоит ли производить (улучшает ли это критерий)
                    // В простейшем случае производим все, что можем
                    resultPlan[j] += float64(maxPossible)
                    for i := 0; i < m.SizeA; i++ {
                        remainingResources[i] -= m.Model[i][j] * float64(maxPossible)
                    }
                    fmt.Printf("  + Продукт %d: дополнительно %d ед. из остатков\n", 
                        j+1, maxPossible)
                    break
                }
            }
        }
    }

    // 5. Расчет статистики
    fmt.Println("\n" + strings.Repeat("=", 50))
    fmt.Println("Итоговый план производства:")
    for j := 0; j < m.SizeB; j++ {
        fmt.Printf("  Продукт %d: %.0f ед. (критерий: %.2f)\n", 
            j+1, resultPlan[j], criteria[j])
    }
    
    // Использование ресурсов
    fmt.Println("\nИспользование ресурсов:")
    for i := 0; i < m.SizeA; i++ {
        used := 0.0
        for j := 0; j < m.SizeB; j++ {
            used += m.Model[i][j] * resultPlan[j]
        }
        utilization := (used / m.A[i]) * 100
        fmt.Printf("  Ресурс %d: %.1f/%.1f (использовано: %.1f%%)\n", 
            i+1, used, m.A[i], utilization)
    }
    
    // Значение целевой функции (сумма критериев * количество)
    totalValue := 0.0
    for j := 0; j < m.SizeB; j++ {
        totalValue += criteria[j] * resultPlan[j]
    }
    fmt.Printf("\nЗначение целевой функции: %.2f\n", totalValue)

    return resultPlan, nil
}