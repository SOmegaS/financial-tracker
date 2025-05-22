package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

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

// metrics tracks load test performance
type metrics struct {
	successCount uint64
	failureCount uint64
	totalLatency time.Duration
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

	// Metrics collection
	var m metrics
	var wg sync.WaitGroup

	// Start load test
	startTime := time.Now()
	for i := 0; i < numUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			userMetrics := runUserOperations(userClient, publisherClient, readerClient, userID)
			atomic.AddUint64(&m.successCount, userMetrics.successCount)
			atomic.AddUint64(&m.failureCount, userMetrics.failureCount)
			atomic.AddInt64((*int64)(&m.totalLatency), int64(userMetrics.totalLatency))
		}(i)
	}

	// Wait for all users to complete
	wg.Wait()
	duration := time.Since(startTime)

	// Report metrics
	avgLatency := time.Duration(0)
	if m.successCount > 0 {
		avgLatency = time.Duration(int64(m.totalLatency) / int64(m.successCount))
	}
	log.Printf("Load test completed in %v", duration)
	log.Printf("Total requests: %d", m.successCount+m.failureCount)
	log.Printf("Successes: %d", m.successCount)
	log.Printf("Failures: %d", m.failureCount)
	log.Printf("Average latency per request: %v", avgLatency)
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
		m.failureCount++
	} else {
		m.successCount++
		m.totalLatency += latency
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
		m.failureCount++
		return m // Stop if login fails
	}
	m.successCount++
	m.totalLatency += latency
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
			m.failureCount++
		} else {
			m.successCount++
			m.totalLatency += latency
		}

		// GetReport (reader service)
		start = time.Now()
		_, err = readerClient.GetReport(ctx, &api.GetReportRequest{Jwt: jwt})
		latency = time.Since(start)
		if err != nil {
			log.Printf("User %d GetReport failed: %v", userID, err)
			m.failureCount++
		} else {
			m.successCount++
			m.totalLatency += latency
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
			m.failureCount++
		} else {
			m.successCount++
			m.totalLatency += latency
		}
	}

	return m
}
