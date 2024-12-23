package config

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func Test_Init(t *testing.T) {
	expected := "Seamless Wallet"
	assertEqual(t, Of.App.Name, expected)

	// Print out all
	fmt.Println(Of.App.GetPublicKey())
}

func Test_WatchChange(t *testing.T) {
	expected := "Seamless"

	// Do sleep to delaying program
	// So you can change the value of config manually and viper config will watch the new config value
	time.Sleep(5 * time.Second)

	assertEqual(t, Of.App.Name, expected)
}

func TestKafkaConfig(t *testing.T) {
	if len(Of.Kafka.Servers) <= 0 {
		t.Error("kafka server config not found")
	}
}

func Benchmark_Init(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Test_Init(&testing.T{})
	}
}
