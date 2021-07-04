# Writing tests for Cloudinary SDK

If you are creating a new feature it is highly recommended adding tests for it.

## Types of tests
There are two types of tests in Cloudinary SDK: E2E and acceptance.
### End-to-End tests (E2E)
These tests are placed in `*_test.go` files.
The main purpose of these tests is to check the main scenarios in the SDK integration with Cloudinary server.
These tests are performing a HTTP queries to the Cloudinary server and requires CLOUDINARY_URL to be given to run.

### Acceptance tests
These tests are placed in `*_acceptance_test.go` files.
The main purpose of these tests is to check whether SDK generates proper HTTP requests or not with mocked HTTP-server.
Although these tests are checking response JSON parsing.

The main acceptance test runner is placed in `admin_test` `testApiByTestCases` method.
This method requires a `[]ApiAcceptanceTestCase` and `*testing.T` to run.

`ApiAcceptanceTestCase` is a test case specification, check it's fields comments for additional information.