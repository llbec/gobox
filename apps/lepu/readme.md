# Lepu Backend Service

This is a Golang backend service for handling patient monitoring data from MySQL database views.

## Features

- Health check endpoint
- Get patient list by start time
- Get vital signs data merged from multiple views

## Database Views

1. `view_case_info`: Patient case information
   - monitor_case_id
   - pid
   - patient_name
   - n_test_no (nullable)
   - n_bar_no (nullable)
   - n_ad_no (nullable)
   - start_time
   - end_time (nullable)

2. `view_case_ecg`: ECG data
   - monitor_case_id
   - data_time
   - hr (may be "---" for invalid values)

3. `view_case_spo2`: SPO2 data
   - monitor_case_id
   - data_time
   - spo2 (may be "---" for invalid values)
   - pr (may be "---" for invalid values)

4. `view_case_resp`: Respiratory data
   - monitor_case_id
   - data_time
   - rr (may be "---" for invalid values)

5. `view_case_nibp`: NIBP data
   - monitor_case_id
   - data_time
   - sys (may be "---" for invalid values)
   - dia (may be "---" for invalid values)
   - map (may be "---" for invalid values)

## API Endpoints

### Health Check
- **POST** `/health`
- Request Body: None
- Returns: `{"status": "OK"}` if both service and database are healthy, otherwise returns error status

### Get Patients
- **POST** `/patients`
- Request Body:
  ```json
  {
    "start_time": "2023-01-01T00:00:00Z"
  }
  ```
- Returns: Array of patient case info

### Get Vitals
- **POST** `/vitals`
- Request Body:
  ```json
  {
    "case_id": "123",
    "start": "2023-01-01T00:00:00Z",
    "end": "2023-01-01T23:59:59Z"
  }
  ```
- Returns: Array of merged vital data `[data_time, hr, spo2, pr, rr, nIBPs, nIBPd, nIBPm]`

## Configuration

Edit `config.yaml` to set database connection:

```yaml
database:
  url: "monitor_user:monitor_user@123@tcp(192.168.1.242:3307)/cardiot_monitor"
  username: "monitor_user"
  password: "monitor_user@123"
```

## Running

1. Install dependencies: `go mod tidy`
2. Run: `go run main.go`

## Project Structure

- `config/`: Configuration loading
- `logger/`: Logging utilities
- `views/`: Database view queries
- `http/`: HTTP handlers
- `main.go`: Entry point