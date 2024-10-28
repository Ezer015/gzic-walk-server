# GZIC Walk Server API

## Base URL

`https://api.example.com/v1`

## Endpoints

### Image Processing

#### Image Upload

**Endpoint:** `/images`  

**Method:** `POST`  

**Description:** Uploads an image for processing.

**Request:**

- **Headers:**
  - `Content-Type: multipart/form-data`
- **Body:**
  - `file`: The image file to upload (required).

**Response:**

- **Status Code:** `202 Accepted`
- **Body:**
  
  ```json
  {
  "processed_image_id": integer
  }
  ```

**Example Request:**

```http
POST /images
```

#### Image Download

**Endpoint:** `/images/{image_id}`  

**Method:** `GET`  

**Description:** Retrieves the (processed) image based on the provided image_id.

**Path Parameter:**

- `image_id`: The unique identifier for the image (required).

**Response:**

- **Status Code:** `200 OK` (if the image ID exists)
  - **Body:** The processed image file in the response body.
- **Status Code:** `202 Accepted` (if the image has not been processed yet).
- **Status Code:** `404 Not Found` (if the image ID does not exist).

**Example Request:**

```http
GET /images/15
```

### Sights Information

#### Sights Information List

**Endpoint:** `/sights`  

**Method:** `GET`  

**Description:** Retrieves a list of all available sights.

**Response:**

- **Status Code:** `200 OK`
- **Body:**
  
  ```json
  {
    "sights": [
      {
        "sight_id": "string",
        "name": "string",
        "description": "string"
      }
    ]
  }
  ```

**Example Request:**

```http
GET /sights
```

#### Sight Information Retrieval

**Endpoint:** `/sights/{sight_id}`  

**Method:** `GET`  

**Description:** Retrieves information about a specific sight by its ID.

**Path Parameter:**

- `sight_id`: The unique identifier of the sight (required).

**Response:**

- **Status Code:** `200 OK`
- **Body:**
  
  ```json
  {
    "sight_id": "string",
    "name": "string",
    "description": "string"
  }
  ```

- **Status Code:** `404 Not Found` (if the sight ID does not exist).

**Example Request:**

```http
GET /sights/15
```

### AI Copywriting

#### AI Copywriting Generation

**Endpoint:** `/copywriting`  

**Method:** `POST`  

**Description:** Initiates the generation of copywriting based on styles.

**Request:**

- **Headers:**
  
  - `Content-Type: application/json`
- **Body:**
  
  ```json
  {
    "name": "string",
    "description": "string",
    "style": "string"
  }
  ```

**Response:**

- **Status Code:** `202 Accepted`
- **Body:**

  ```json
  {
    "copywriting_id": integer
  }
  ```

**Example Request:**

```http
POST /copywriting
{
  "name": "Sunset",
  "description": "A beautiful sunset over the ocean with waves crashing onto the shore.",
  "style": "tiktok"
}
```

#### AI Copywriting Retrieval

**Endpoint:** `/copywriting/{copywriting_id}`  
**Method:** `GET`  
**Description:** Checks the status of the copywriting job or retrieves the result if complete.

**Path Parameter:**

- `copywriting_id`: The unique identifier for the copy-writing (required).

**Response:**

- **Status Code:** `200 OK` (if the job is complete)

- **Body:**

  ```json
  {
    "copywriting": "string"
  }
  ```

- **Status Code:** `202 Accepted` (if the job is still in progress)

- **Status Code:** `404 Not Found` (if the job ID does not exist).

**Example Request:**

```http
GET /copywriting/15
```

### Record

#### Record Creation

**Endpoint:** `/record`  

**Method:** `POST`  

**Description:** Create random record including an image, sight, and copywriting.

**Request:**

- **Headers:**
  - `Content-Type: application/json`
- **Body:**
  
  ```json
  {
    "image_id": integer,
    "sight_id": integer,
    "sight_name": "string",
    "copywriting": "string"
  }
  ```

**Response:**

- **Status Code:** `201 Created`
- **Body:**
  
  ```json
  {
    "record_id": "string"
  }
  ```

**Example Request:**

```http
POST /random
{
  "image_id": 15,
  "sight_id": null,
  "sight_name": "sunset",
  "copywriting": "A beautiful sight to behold."
}
```

#### Record Retrieval

**Endpoint:** `/record/{record_id}`  

**Method:** `GET`  

**Description:** Retrieves a record by its ID.

**Path Parameter:**

- `record_id`: The unique identifier of the random record (required).

**Response:**

- **Status Code:** `200 OK`
- **Body:**
  
  ```json
  {
    "image_id": integer,
    "sight_id": integer,
    "sight_name": "string",
    "copywriting": "string"
  }
  ```

- **Status Code:** `404 Not Found` (if the record ID does not exist).

**Example Request:**

```http
GET /record/15
```

#### Random Record Retrieval

**Endpoint:** `/record`  

**Method:** `GET`  

**Description:** Retrieves a random record.

**Response:**

- **Status Code:** `200 OK`
- **Body:**
  
  ```json
  {
    "image_id": integer,
    "sight_id": integer,
    "sight_name": "string",
    "copywriting": "string"
  }
  ```

- **Status Code:** `404 Not Found` (if the record ID does not exist).

**Example Request:**

```http
GET /record
```

## Error Handling

The API uses standard HTTP status codes for indicating the success or failure of requests.

- `400 Bad Request`: The request was invalid. Check your input data.
- `404 Not Found`: The requested resource could not be found.
- `500 Internal Server Error`: An unexpected error occurred on the server.