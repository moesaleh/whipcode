## API reference

### Headers
- `Content-Type: application/json`
- `X-Master-Key: $MASTER_KEY`

### Body
- **code** [string]\
  *The source code, base64 encoded.*

- **language_id** [integer/string]\
  *Language ID of the submitted code.*

- **args** [string] [optional]\
  *Compiler/interpreter args separated by spaces.*

- **timeout** [integer] [optional]\
  *Timeout in seconds for the code to run. Capped at the timeout set in whipcode's configuration.*

### Response
- **stdout** [string]\
  *All data captured from stdout.*

- **stderr** [string]\
  *All data captured from stderr.*

- **container_age** [float]\
  *Duration the container allocated for your code ran, in seconds.*

- **timeout** [boolean]\
  *Boolean value depending on whether your container lived past the timeout period. A reply from a timed-out request will not have any data in stdout and stderr.*

In the event of an error or an invalid/rejected request, the response will contain only:

- **detail** [string]\
  *Details about why the request failed to complete.*


### Example request
```bash
lang=2
code='console.log("Hello world!");'
timeout=5

curl -s -X POST $ENDPOINT \
    -H 'Content-Type: application/json' \
    -H "X-Master-Key: $MASTER_KEY" \
    -d '{
        "language_id": "'$lang'",
        "code": "'$(echo -n $code | base64)'",
        "timeout": "'$timeout'"
    }' | jq
```


### Example response
```json
{
  "stdout": "Hello world!\n",
  "stderr": "",
  "container_age": 0.335837,
  "timeout": false
}
```
