import datetime
from locust import User, task, between, events
import grpc
import random
import string
import logging
import time
from NFT.grpc import common_pb2 as api_pb2, common_pb2_grpc as api_pb2_grpc

class GrpcUser(User):
    wait_time = between(1, 5)

    username = ""
    jwt = ""
    
    def on_start(self):
        # Establish a gRPC channel
        self.user_channel = grpc.insecure_channel('localhost:7778')
        self.user_stub = api_pb2_grpc.ApiStub(self.user_channel)

        self.expense_channel = grpc.insecure_channel('localhost:7779')
        self.expense_stub = api_pb2_grpc.ApiStub(self.expense_channel)

        self.publisher_channel = grpc.insecure_channel('localhost:7777')
        self.publisher_stub = api_pb2_grpc.ApiStub(self.publisher_channel)

    @task(1)
    def register_user(self):
        # Generate random username to avoid conflicts
        random_suffix = ''.join(random.choices(string.ascii_lowercase + string.digits, k=8))
        username = f"testuser_{random_suffix}"
        
        # Create a RegisterRequest
        request = api_pb2.RegisterRequest(
            requestId=str(random.randint(1, 1000)),
            username=username,
            password="password123"
        )
        
        # Call the Register method and track stats
        start_time = time.time()
        request_type = "register"
        exception = None
        try:
            response = self.user_stub.Register(request)
            self.username = username
            self.jwt = response.jwt
        except grpc.RpcError as e:
            exception = e
        
        # Record request stats
        total_time = int((time.time() - start_time) * 1000)
        events.request.fire(
            request_type=request_type,
            name="Register",
            response_time=total_time,
            response_length=0,
            exception=exception,
        )

    @task(2)
    def login_user(self):
        if not self.username:
            return

        # Create a LoginRequest
        request = api_pb2.LoginRequest(
            requestId=str(random.randint(1, 1000)),
            username=self.username,
            password="password123"
        )

        # Call the Login method and track stats
        start_time = time.time()
        request_type = "login"
        exception = None
        try:
            response = self.user_stub.Login(request)
            self.jwt = response.jwt
        except grpc.RpcError as e:
            exception = e
        
        # Record request stats
        total_time = int((time.time() - start_time) * 1000)
        events.request.fire(
            request_type=request_type,
            name="Login",
            response_time=total_time,
            response_length=0,
            exception=exception,
        )

    @task(3)
    def get_report(self):
        if not self.username:
            return

        # Create a GetReportRequest
        request = api_pb2.GetReportRequest(jwt=self.jwt)

        # Call the GetReport method and track stats
        start_time = time.time()
        request_type = "get_report"
        exception = None
        try:
            response = self.expense_stub.GetReport(request)
            logging.info(f"GetReport response: {response}")
        except grpc.RpcError as e:
            exception = e
            logging.error(f"GetReport failed: {e.details()}")

        # Record request stats
        total_time = int((time.time() - start_time) * 1000)
        events.request.fire(
            request_type=request_type,
            name="GetReport",
            response_time=total_time,
            response_length=0,
            exception=exception,
        )

    @task(10)
    def get_bills(self):
        if not self.username:
            return

        # Create a GetBillsRequest
        request = api_pb2.GetBillsRequest(jwt=self.jwt, category="utilities")

        # Call the GetBills method and track stats
        start_time = time.time()
        request_type = "get_bills"
        exception = None
        try:
            response = self.expense_stub.GetBills(request)
            logging.info(f"GetBills response: {response}")
        except grpc.RpcError as e:
            exception = e
            logging.error(f"GetBills failed: {e.details()}")

        # Record request stats
        total_time = int((time.time() - start_time) * 1000)
        events.request.fire(
            request_type=request_type,
            name="GetBills",
            response_time=total_time,
            response_length=0,
            exception=exception,
        )

    @task(10)
    def create_bill(self):
        if not self.username:
            return

        # Create a BillMessage
        request = api_pb2.BillMessage(
            name="Electricity",
            amount=100.0,
            category="utilities",
            timestamp=datetime.datetime.now(),
            jwt=self.jwt
        )

        # Call the CreateBill method and track stats
        start_time = time.time()
        request_type = "create_bill"
        exception = None
        try:
            response = self.publisher_stub.CreateBill(request)
            logging.info("CreateBill succeeded")
        except grpc.RpcError as e:
            exception = e
            logging.error(f"CreateBill failed: {e.details()}")

        # Record request stats
        total_time = int((time.time() - start_time) * 1000)
        events.request.fire(
            request_type=request_type,
            name="CreateBill",
            response_time=total_time,
            response_length=0,
            exception=exception,
        )

    def on_stop(self):
        # Close the gRPC channel
        self.user_channel.close()
        self.expense_channel.close()
        self.publisher_channel.close()