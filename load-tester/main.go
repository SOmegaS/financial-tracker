package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"load-tester/api" // Adjust this import path based on your project structure
)

const (
	userServiceAddr   = "financial-tracker-user-service:7777"      // Address for Register and Login
	publisherAddr     = "financial-tracker-expense-publisher:7777" // Address for CreateBill
	readerAddr        = "financial-tracker-expense-reader:7777"    // Address for GetReport and GetBills
	numUsers          = 100                                        // Number of simulated users
	operationsPerUser = 10                                         // Number of operations per user
)

// endpointMetrics tracks per-endpoint performance
type endpointMetrics struct {
	successCount prometheus.Counter
	failureCount prometheus.Counter
	latency      prometheus.Histogram
}

// metrics tracks overall load test performance
type metrics struct {
	register   endpointMetrics
	login      endpointMetrics
	createBill endpointMetrics
	getReport  endpointMetrics
	getBills   endpointMetrics
	total      struct {
		successCount uint64
		failureCount uint64
		totalLatency time.Duration
	}
}

var (
	// Global metrics instance
	m = metrics{
		register: endpointMetrics{
			successCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_register_success_total",
				Help: "Total number of successful Register requests",
			}),
			failureCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_register_failure_total",
				Help: "Total number of failed Register requests",
			}),
			latency: prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:    "load_tester_register_latency_seconds",
				Help:    "Register request latency in seconds",
				Buckets: prometheus.DefBuckets,
			}),
		},
		login: endpointMetrics{
			successCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_login_success_total",
				Help: "Total number of successful Login requests",
			}),
			failureCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_login_failure_total",
				Help: "Total number of failed Login requests",
			}),
			latency: prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:    "load_tester_login_latency_seconds",
				Help:    "Login request latency in seconds",
				Buckets: prometheus.DefBuckets,
			}),
		},
		createBill: endpointMetrics{
			successCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_create_bill_success_total",
				Help: "Total number of successful CreateBill requests",
			}),
			failureCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_create_bill_failure_total",
				Help: "Total number of failed CreateBill requests",
			}),
			latency: prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:    "load_tester_create_bill_latency_seconds",
				Help:    "CreateBill request latency in seconds",
				Buckets: prometheus.DefBuckets,
			}),
		},
		getReport: endpointMetrics{
			successCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_get_report_success_total",
				Help: "Total number of successful GetReport requests",
			}),
			failureCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_get_report_failure_total",
				Help: "Total number of failed GetReport requests",
			}),
			latency: prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:    "load_tester_get_report_latency_seconds",
				Help:    "GetReport request latency in seconds",
				Buckets: prometheus.DefBuckets,
			}),
		},
		getBills: endpointMetrics{
			successCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_get_bills_success_total",
				Help: "Total number of successful GetBills requests",
			}),
			failureCount: prometheus.NewCounter(prometheus.CounterOpts{
				Name: "load_tester_get_bills_failure_total",
				Help: "Total number of failed GetBills requests",
			}),
			latency: prometheus.NewHistogram(prometheus.HistogramOpts{
				Name:    "load_tester_get_bills_latency_seconds",
				Help:    "GetBills request latency in seconds",
				Buckets: prometheus.DefBuckets,
			}),
		},
	}
)

func init() {
	// Register Prometheus metrics
	prometheus.MustRegister(
		m.register.successCount,
		m.register.failureCount,
		m.register.latency,
		m.login.successCount,
		m.login.failureCount,
		m.login.latency,
		m.createBill.successCount,
		m.createBill.failureCount,
		m.createBill.latency,
		m.getReport.successCount,
		m.getReport.failureCount,
		m.getReport.latency,
		m.getBills.successCount,
		m.getBills.failureCount,
		m.getBills.latency,
	)

	// Start Prometheus metrics server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func main() {
	// Connect to the user service (Register, Login)
	userConn, err := grpc.Dial(userServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := api.NewApiClient(userConn)

	// Connect to the publisher service (CreateBill)
	publisherConn, err := grpc.Dial(publisherAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to publisher service: %v", err)
	}
	defer publisherConn.Close()
	publisherClient := api.NewApiClient(publisherConn)

	// Connect to the reader service (GetReport, GetBills)
	readerConn, err := grpc.Dial(readerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to reader service: %v", err)
	}
	defer readerConn.Close()
	readerClient := api.NewApiClient(readerConn)

	// Start load test
	var wg sync.WaitGroup
	startTime := time.Now()
	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			userMetrics := runUserOperations(userClient, publisherClient, readerClient, userID)
			atomic.AddUint64(&m.total.successCount, userMetrics.total.successCount)
			atomic.AddUint64(&m.total.failureCount, userMetrics.total.failureCount)
			atomic.AddInt64((*int64)(&m.total.totalLatency), int64(userMetrics.total.totalLatency))
		}(i)
	}

	// Wait for all users to complete
	wg.Wait()
	duration := time.Since(startTime)

	// Report metrics
	avgLatency := time.Duration(0)
	if m.total.successCount > 0 {
		avgLatency = time.Duration(int64(m.total.totalLatency) / int64(m.total.successCount))
	}
	log.Printf("Load test completed in %v", duration)
	log.Printf("Total requests: %d", m.total.successCount+m.total.failureCount)
	log.Printf("Successes: %d", m.total.successCount)
	log.Printf("Failures: %d", m.total.failureCount)
	log.Printf("Average latency per request: %v", avgLatency)
	log.Printf("Register success: %d, failure: %d", m.register.successCount, m.register.failureCount)
	log.Printf("Login success: %d, failure: %d", m.login.successCount, m.login.failureCount)
	log.Printf("CreateBill success: %d, failure: %d", m.createBill.successCount, m.createBill.failureCount)
	log.Printf("GetReport success: %d, failure: %d", m.getReport.successCount, m.getReport.failureCount)
	log.Printf("GetBills success: %d, failure: %d", m.getBills.successCount, m.getBills.failureCount)
	log.Printf("Average latency for Register: %v", m.register.latency.Collect())
	log.Printf("Average latency for Login: %v", m.login.latency.Collect())
	log.Printf("Average latency for CreateBill: %v", m.createBill.latency.Collect())
	log.Printf("Average latency for GetReport: %v", m.getReport.latency.Collect())
	log.Printf("Average latency for GetBills: %v", m.getBills.latency.Collect())
}

func runUserOperations(userClient, publisherClient, readerClient api.ApiClient, userID int) metrics {
	var m metrics
	username := fmt.Sprintf("user%d", userID)
	password := "password123"
	requestID := fmt.Sprintf("req-%d-%d", userID, time.Now().UnixNano())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Register (user service)
	start := time.Now()
	_, err := userClient.Register(ctx, &api.RegisterRequest{
		RequestId: requestID,
		Username:  username,
		Password:  password,
	})
	latency := time.Since(start)
	if err != nil {
		log.Printf("User %d Register failed: %v", userID, err)
		m.register.failureCount.Inc()
		m.total.failureCount++
	} else {
		m.register.successCount.Inc()
		m.register.latency.Observe(latency.Seconds())
		m.total.successCount++
		m.total.totalLatency += latency
	}

	// Login (user service)
	start = time.Now()
	loginResp, err := userClient.Login(ctx, &api.LoginRequest{
		RequestId: requestID,
		Username:  username,
		Password:  password,
	})
	latency = time.Since(start)
	if err != nil {
		log.Printf("User %d Login failed: %v", userID, err)
		m.login.failureCount.Inc()
		m.total.failureCount++
		return m // Stop if login fails
	}
	m.login.successCount.Inc()
	m.login.latency.Observe(latency.Seconds())
	m.total.successCount++
	m.total.totalLatency += latency
	jwt := loginResp.Jwt

	// Perform operations
	for i := 0; i < operationsPerUser; i++ {
		// CreateBill (publisher service)
		start = time.Now()
		_, err = publisherClient.CreateBill(ctx, &api.CreateBillMessage{
			Name:      fmt.Sprintf("Bill %d", i),
			Amount:    float64(100.0 + float64(i)*10.0),
			Category:  "TestCategory",
			UserId:    fmt.Sprintf("user%d", userID),
			Timestamp: timestamppb.Now(),
			Jwt:       jwt,
		})
		latency = time.Since(start)
		if err != nil {
			log.Printf("User %d CreateBill %d failed: %v", userID, i, err)
			m.createBill.failureCount.Inc()
			m.total.failureCount++
		} else {
			m.createBill.successCount.Inc()
			m.createBill.latency.Observe(latency.Seconds())
			m.total.successCount++
			m.total.totalLatency += latency
		}

		// GetReport (reader service)
		start = time.Now()
		_, err = readerClient.GetReport(ctx, &api.GetReportRequest{Jwt: jwt})
		latency = time.Since(start)
		if err != nil {
			log.Printf("User %d GetReport failed: %v", userID, err)
			m.getReport.failureCount.Inc()
			m.total.failureCount++
		} else {
			m.getReport.successCount.Inc()
			m.getReport.latency.Observe(latency.Seconds())
			m.total.successCount++
			m.total.totalLatency += latency
		}

		// GetBills (reader service)
		start = time.Now()
		_, err = readerClient.GetBills(ctx, &api.GetBillsRequest{
			Jwt:      jwt,
			Category: "TestCategory",
		})
		latency = time.Since(start)
		if err != nil {
			log.Printf("User %d GetBills failed: %v", userID, err)
			m.getBills.failureCount.Inc()
			m.total.failureCount++
		} else {
			m.getBills.successCount.Inc()
			m.getBills.latency.Observe(latency.Seconds())
			m.total.successCount++
			m.total.totalLatency += latency
		}
	}

	return m
}
