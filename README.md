# Yaml Metadata Exercise

## Configuration

All options below are read from cooresponding environment variables and will default to the stated default value:

```
{
  // The name for the server, used for discovery, logging, etc...
  // DEFAULT: 'application.metadata.exercise'
  // TYPE: STRING
  SERVER_NAME,

  // The address on which the service should listen
  // DEFAULT: '127.0.0.1'
  // TYPE: STRING
  SERVER_HOST,

  // The port on which the service should listen
  // DEFAULT: 8071
  // TYPE: NUMBER
  SERVER_PORT,

  // This setting allows you to specify whether to store 'mutliple' or 'single' instances of metadata on the server
  // DEFAULT: 'multiple'
  // TYPE: STRING
  STORAGE_MODE,

  // The directory name where the bleve indexed data will be stored
  // DEFAULT: applicationMetadata.bleve
  // TYPE: STRING
  INDEX_NAME,

  // The keys of the fields you want to use as the unique identifier for the metadata instance in multiple mode
  // Posting to an existing identifier will overwrite the document at that identifier
  // DEFAULT: 'title,source'
  // TYPE: []STRING
  IDENTIFIER_FIELDS,
}
```

### Prerequisites

- go: https://golang.org/doc/install
- dep: https://github.com/golang/dep

## Build/Run

To build Yaml Metadata Exercise, be sure you have dep installed. Then run `make ensure`. This will initialize all dependencies in the `vendor` directory.

To test the project, run `make test`.

To start the server, run `make start` passing in any arguments from above you would like to set in a key value fashion
like so: `make start SERVER_NAME=test INDEX_NAME=some.test`.

**Write Metadata**
----
  Writes the metadata to the server. If in 'single' mode will always overwrite what is there, if in 'multiple' mode, will write a new document or overwrite a document matching the generated identifier based on the fields the server is using.

**URL:**
  `/application/metadata`

**Method:**
  `POST`

**URL Params:**

  - `format=[string]` -- yaml, json
    - default: yaml
    - returns the data in the selected format

**Data Params**

- yaml:

  - `title`         -- string
  - `version`       -- string
  - `maintainers`   -- array
    - `name`        -- string
    - `email`       -- string
  - `company`       -- string
  - `website`       -- string
  - `source`        -- string
  - `license`       -- string
  - `description`   -- string

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `OK`

* **Error Response:**

  * **Code:** 400 <br />
    **Content:**
    - `yaml || json list of validation errors`

---

**Delete Metadata**
----
  Deletes the metadata matching the identifier key

**URL:**
  `/application/metadata`

**Method:**
  `DELETE`

**Data Params**

- yaml:

  - `title`         -- string
  - `version`       -- string
  - `maintainers`   -- array
    - `name`        -- string
    - `email`       -- string
  - `company`       -- string
  - `website`       -- string
  - `source`        -- string
  - `license`       -- string
  - `description`   -- string

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `OK`

* **Error Response:**

  * **Code:** 400 <br />
    **Content:**
    - `can't parse yaml`
---

**Search Metadata**
----
  Searches the content of the metadata

**URL:**
  `/application/metadata`

**Method:**
  `GET`

**URL Params:**

- optional
  - `query=[string]` -- searches for a match in your metadata, excluding en stopwords. If you leave it blank, it gets all metadata on the server
  - `format=[string]` -- yaml, json
    - default: yaml
    - returns the data in the selected format
  - `field=[string]` -- title, version, maintainers, company, website, source, license, description
    - default: ''
    - adding this parameter will narrow the search down to the selected field
  - `fill=[bool]` -- true, false
    - default: false
    - true will return the whole metadata object

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `OK`

* **Error Response:**

  * **Code:** 400 <br />
    **Content:**
    - `can't parse yaml`