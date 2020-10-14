# Load Funds Handler

This application is a transaction handler. It will accept or reject the transaction based on the following rules:

* A maximum of $5,000 can be loaded per day
* A maximum of $20,000 can be loaded per week
* A maximum of 3 loads can be performed per day, regardless of amount 

For this example, the input data is in the [input file](./input.txt).

It will consider the json sample:

```json
{ 
    "id":"15887",
    "customer_id":"528",
    "load_amount":"$3318.47",
    "time":"2000-01-01T00:00:00Z"
}
```

The output expected is:

```json
{ 
    "id": "1234", 
    "customer_id": "1234", 
    "accepted": true
}
```

## Logic implemented 

The idea is to have channels to receive the input and also to send the output. It tries to simulate a queue/event system. 

To control the order and also the end of the output reading, it was used the `sync.WaitGroup` for that.

The busines logic is on the handler package.

## Running and testing

To help with that, this project has a Makefile with several parameters.

### Runing

You can run the program using your golang instaled version (>=1.13) or only using Docker:

```shell
make run
make docker-run
```

### Testing

The same way for testing, you have the two options:

```shell
make tests
make docker-tests
```

You can also see the coverage running the command:

```shell
make coverage
```
