package main

import (
	"errors"
	"fmt"
	"time"

	"math/rand"

	"github.com/sony/gobreaker"
)

// callExternalService simula una llamada a un servicio externo que puede fallar
func callExternalService() (string, error) {
	randomDuration := time.Duration(rand.Intn(2001)) * time.Millisecond
	fmt.Println("Esperando:", randomDuration)
	time.Sleep(randomDuration)
	// simulamos un error aleatorio para demostrar el comportamiento del circuit breaker
	if rand.Intn(2) == 0 {
		return "", errors.New("simulated error")
	}
	return "success", nil
}

func main() {
	// configuraci{on del cicuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name: "externalServiceCB", // Es el nombre del CircuitBreaker
		ReadyToTrip: func(counts gobreaker.Counts) bool { // Se llama a ReadyToTrip con una copia de Counts cada vez que falla una solicitud en el estado cerrado. Si ReadyToTrip devuelve verdadero, el CircuitBreaker se colocará en estado abierto. Si ReadyToTrip es nulo, se utiliza ReadyToTrip predeterminado. ReadyToTrip predeterminado devuelve verdadero cuando el número de fallas consecutivas es superior a 5.
			return counts.TotalFailures >= 1
		},
		Timeout: 2 * time.Second, // Es el período del estado abierto, después del cual el estado del CircuitBreaker pasa a estar half-open. Si el tiempo de espera es menor o igual a 0, el valor del tiempo de espera del CircuitBreaker se establece en 60 segundos.
		// Interval: 2 * time.Second, // Es el período cíclico del estado cerrado para que el CircuitBreaker borre los conteos internos. Si el intervalo es menor o igual a 0, el CircuitBreaker no borra los conteos internos durante el estado cerrado.
		// MaxRequests: 0,               // Es la cantidad máxima de solicitudes que pueden pasar cuando el CircuitBreaker está half-open. Si MaxRequests es 0, el CircuitBreaker permite solo 1 solicitud.
	})

	for range [80]int{} {
		// ejecutamos la llamada al servicio protegido por el circuit breaker
		result, err := cb.Execute(func() (interface{}, error) {
			return callExternalService()
		})
		if err != nil {
			fmt.Print(err, "\n\n")
			if cb.State() == gobreaker.StateOpen {
				fmt.Print("circuit breaker is open\n\n")
				time.Sleep(1 * time.Second)
			}
			continue
		}
		fmt.Print(result, "\n\n")
	}
}
